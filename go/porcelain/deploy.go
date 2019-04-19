package porcelain

import (
	"archive/zip"
	"bufio"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"debug/elf"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"

	"github.com/go-openapi/errors"
	"github.com/rsc/goversion/version"
)

const (
	jsRuntime = "js"
	goRuntime = "go"

	preProcessingTimeout = time.Minute * 5

	fileUpload uploadType = iota
	functionUpload

	lfsVersionString = "version https://git-lfs.github.com/spec/v1"
)

type uploadType int
type pointerData struct {
	SHA  string
	Size int64
}

type DeployObserver interface {
	OnSetupWalk() error
	OnSuccessfulStep(*FileBundle) error
	OnSuccessfulWalk(*models.DeployFiles) error
	OnFailedWalk()

	OnSetupDelta(*models.DeployFiles) error
	OnSuccessfulDelta(*models.DeployFiles, *models.Deploy) error
	OnFailedDelta(*models.DeployFiles)

	OnSetupUpload(*FileBundle) error
	OnSuccessfulUpload(*FileBundle) error
	OnFailedUpload(*FileBundle)
}

// DeployOptions holds the option for creating a new deploy
type DeployOptions struct {
	SiteID            string
	Dir               string
	FunctionsDir      string
	LargeMediaEnabled bool

	IsDraft bool

	Title             string
	Branch            string
	CommitRef         string
	UploadTimeout     time.Duration
	PreProcessTimeout time.Duration

	Observer DeployObserver

	files     *deployFiles
	functions *deployFiles
}

type uploadError struct {
	err   error
	mutex *sync.Mutex
}

type FileBundle struct {
	Name    string
	Sum     string
	Runtime string
	Size    *int64 `json:"size,omitempty"`

	// Path OR Buffer should be populated
	Path   string
	Buffer io.ReadSeeker
}

func (f *FileBundle) Read(p []byte) (n int, err error) {
	return f.Buffer.Read(p)
}

func (f *FileBundle) Seek(offset int64, whence int) (int64, error) {
	return f.Buffer.Seek(offset, whence)
}

func (f *FileBundle) Close() error {
	return nil
}

type deployFiles struct {
	Files  map[string]*FileBundle
	Sums   map[string]string
	Hashed map[string][]*FileBundle
}

func newDeployFiles() *deployFiles {
	return &deployFiles{
		Files:  make(map[string]*FileBundle),
		Sums:   make(map[string]string),
		Hashed: make(map[string][]*FileBundle),
	}
}

func (d *deployFiles) Add(p string, f *FileBundle) {
	d.Files[p] = f
	d.Sums[p] = f.Sum
	// Remove ":original_sha" part when to save in Hashed (large media)
	sum := f.Sum
	if strings.Contains(sum, ":") {
		sum = strings.Split(sum, ":")[0]
	}
	list, _ := d.Hashed[sum]
	d.Hashed[sum] = append(list, f)
}

func (n *Netlify) overCommitted(d *deployFiles) bool {
	return len(d.Files) > n.syncFileLimit
}

// GetDeploy returns a deploy.
func (n *Netlify) GetDeploy(ctx context.Context, deployID string) (*models.Deploy, error) {
	authInfo := context.GetAuthInfo(ctx)
	resp, err := n.Netlify.Operations.GetDeploy(operations.NewGetDeployParams().WithDeployID(deployID), authInfo)
	if err != nil {
		return nil, err
	}
	return resp.Payload, nil
}

// DeploySite creates a new deploy for a site given a directory in the filesystem.
// It uploads the necessary files that changed between deploys.
func (n *Netlify) DeploySite(ctx context.Context, options DeployOptions) (*models.Deploy, error) {
	return n.DoDeploy(ctx, &options, nil)
}

// DoDeploy deploys the changes for a site given a directory in the filesystem.
// It uploads the necessary files that changed between deploys.
func (n *Netlify) DoDeploy(ctx context.Context, options *DeployOptions, deploy *models.Deploy) (*models.Deploy, error) {
	f, err := os.Stat(options.Dir)
	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", options.Dir)
	}

	if options.Observer != nil {
		if err := options.Observer.OnSetupWalk(); err != nil {
			return nil, err
		}
	}

	largeMediaEnabled := options.LargeMediaEnabled
	if !largeMediaEnabled && deploy != nil {
		largeMediaEnabled = deploy.SiteCapabilities.LargeMediaEnabled
	}

	context.GetLogger(ctx).Infof("Getting files info with large media flag: %v", largeMediaEnabled)
	files, err := walk(options.Dir, options.Observer, largeMediaEnabled)
	if err != nil {
		if options.Observer != nil {
			options.Observer.OnFailedWalk()
		}
		return nil, err
	}
	for name := range files.Files {
		if strings.ContainsAny(name, "#?") {
			return nil, fmt.Errorf("Invalid filename '%s'. Deployed filenames cannot contain # or ? characters", name)
		}
	}

	options.files = files

	functions, err := bundle(options.FunctionsDir, options.Observer)
	if err != nil {
		if options.Observer != nil {
			options.Observer.OnFailedWalk()
		}
		return nil, err
	}
	options.functions = functions

	deployFiles := &models.DeployFiles{
		Files: options.files.Sums,
		Draft: options.IsDraft,
		Async: n.overCommitted(options.files),
	}
	if options.functions != nil {
		deployFiles.Functions = options.functions.Sums
	}

	if options.Observer != nil {
		if err := options.Observer.OnSuccessfulWalk(deployFiles); err != nil {
			return nil, err
		}
	}

	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id":      options.SiteID,
		"deploy_files": len(options.files.Sums),
	}).Debug("Starting to deploy files")
	authInfo := context.GetAuthInfo(ctx)

	if options.Observer != nil {
		if err := options.Observer.OnSetupDelta(deployFiles); err != nil {
			return nil, err
		}
	}

	if deploy == nil {
		params := operations.NewCreateSiteDeployParams().WithSiteID(options.SiteID).WithDeploy(deployFiles)
		if options.Title != "" {
			params = params.WithTitle(&options.Title)
		}
		resp, err := n.Operations.CreateSiteDeploy(params, authInfo)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedDelta(deployFiles)
			}
			return nil, err
		}
		deploy = resp.Payload
	} else {
		params := operations.NewUpdateSiteDeployParams().WithSiteID(options.SiteID).WithDeployID(deploy.ID).WithDeploy(deployFiles)
		resp, err := n.Operations.UpdateSiteDeploy(params, authInfo)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedDelta(deployFiles)
			}
			return nil, err
		}
		deploy = resp.Payload
	}

	if n.overCommitted(options.files) {
		var err error
		deploy, err = n.WaitUntilDeployReady(ctx, deploy, options.PreProcessTimeout)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedDelta(deployFiles)
			}
			return nil, err
		}
	}

	if options.Observer != nil {
		if err := options.Observer.OnSuccessfulDelta(deployFiles, deploy); err != nil {
			return nil, err
		}
	}

	if len(deploy.Required) == 0 && len(deploy.RequiredFunctions) == 0 {
		return deploy, nil
	}

	if err := n.uploadFiles(ctx, deploy, options.files, options.Observer, fileUpload, options.UploadTimeout); err != nil {
		return nil, err
	}

	if options.functions != nil {
		if err := n.uploadFiles(ctx, deploy, options.functions, options.Observer, functionUpload, options.UploadTimeout); err != nil {
			return nil, err
		}
	}

	return deploy, nil
}

func (n *Netlify) waitForState(ctx context.Context, d *models.Deploy, timeout time.Duration, states ...string) (*models.Deploy, error) {
	authInfo := context.GetAuthInfo(ctx)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	params := operations.NewGetSiteDeployParams().WithSiteID(d.SiteID).WithDeployID(d.ID)
	start := time.Now()
	for t := range ticker.C {
		resp, err := n.Operations.GetSiteDeploy(params, authInfo)
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		context.GetLogger(ctx).WithFields(logrus.Fields{
			"deploy_id": d.ID,
			"state":     resp.Payload.State,
		}).Debug("Waiting until deploy ready")

		for _, state := range states {
			if resp.Payload.State == state {
				return resp.Payload, nil
			}
		}

		if resp.Payload.State == "error" {
			return nil, fmt.Errorf("Error: entered errors state while waiting to enter states: %s", strings.Join(states, ","))
		}

		if t.Sub(start) > timeout {
			return nil, fmt.Errorf("Error: deploy timed out while waiting to enter states: %s", strings.Join(states, ","))
		}
	}

	return d, nil
}

func (n *Netlify) WaitUntilDeployReady(ctx context.Context, d *models.Deploy, timeout time.Duration) (*models.Deploy, error) {
	if timeout <= 0 {
		timeout = preProcessingTimeout
	}

	return n.waitForState(ctx, d, timeout, "prepared", "ready")
}

func (n *Netlify) uploadFiles(ctx context.Context, d *models.Deploy, files *deployFiles, observer DeployObserver, t uploadType, timeout time.Duration) error {
	sharedErr := &uploadError{err: nil, mutex: &sync.Mutex{}}
	sem := make(chan int, n.uploadLimit)
	wg := &sync.WaitGroup{}

	var required []string
	switch t {
	case fileUpload:
		required = d.Required
	case functionUpload:
		required = d.RequiredFunctions
	}

	count := 0
	for _, sha := range required {
		if files, exist := files.Hashed[sha]; exist {
			for range files {
				count++
			}
		}
	}

	context.GetLogger(ctx).Infof("Uploading %v files", count)

	for _, sha := range required {
		if files, exist := files.Hashed[sha]; exist {
			for _, file := range files {
				sem <- 1
				wg.Add(1)

				go n.uploadFile(ctx, d, file, observer, t, timeout, wg, sem, sharedErr)
			}
		}
	}

	wg.Wait()

	return sharedErr.err
}

func (n *Netlify) uploadFile(ctx context.Context, d *models.Deploy, f *FileBundle, c DeployObserver, t uploadType, timeout time.Duration, wg *sync.WaitGroup, sem chan int, sharedErr *uploadError) {
	defer func() {
		wg.Done()
		<-sem
	}()

	sharedErr.mutex.Lock()
	if sharedErr.err != nil {
		sharedErr.mutex.Unlock()
		return
	}
	sharedErr.mutex.Unlock()

	authInfo := context.GetAuthInfo(ctx)

	context.GetLogger(ctx).WithFields(logrus.Fields{
		"deploy_id": d.ID,
		"file_path": f.Name,
		"file_sum":  f.Sum,
	}).Debug("Uploading file")

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute

	if c != nil {
		if err := c.OnSetupUpload(f); err != nil {
			sharedErr.mutex.Lock()
			sharedErr.err = err
			sharedErr.mutex.Unlock()
			return
		}
	}

	err := backoff.Retry(func() error {
		sharedErr.mutex.Lock()

		if sharedErr.err != nil {
			sharedErr.mutex.Unlock()
			return fmt.Errorf("Upload cancelled: %s", f.Name)
		}
		sharedErr.mutex.Unlock()

		var operationError error

		context.GetLogger(ctx).Infof("Uploading file %v", f.Name)

		switch t {
		case fileUpload:
			var body io.ReadCloser
			body, operationError = os.Open(f.Path)
			if operationError == nil {
				defer body.Close()
				params := operations.NewUploadDeployFileParams().WithDeployID(d.ID).WithPath(f.Name).WithFileBody(body)
				if f.Size != nil {
					params.WithSize(f.Size)
				}
				if timeout != 0 {
					params.SetTimeout(timeout)
				}
				_, operationError = n.Operations.UploadDeployFile(params, authInfo)
			}
		case functionUpload:
			params := operations.NewUploadDeployFunctionParams().WithDeployID(d.ID).WithName(f.Name).WithFileBody(f).WithRuntime(&f.Runtime)
			if timeout != 0 {
				params.SetTimeout(timeout)
			}
			_, operationError = n.Operations.UploadDeployFunction(params, authInfo)
			if operationError != nil {
				f.Buffer.Seek(0, 0)
			}
		}

		if operationError != nil {
			context.GetLogger(ctx).WithError(operationError).Errorf("Failed to upload file %v", f.Name)
			apiErr, ok := operationError.(errors.Error)

			if ok && apiErr.Code() == 401 {
				sharedErr.mutex.Lock()
				sharedErr.err = operationError
				sharedErr.mutex.Unlock()
			}
		}
		return operationError
	}, b)

	if err != nil {
		if c != nil {
			c.OnFailedUpload(f)
		}

		sharedErr.mutex.Lock()
		sharedErr.err = err
		sharedErr.mutex.Unlock()
	} else {
		if c != nil {
			if err := c.OnSuccessfulUpload(f); err != nil {
				sharedErr.mutex.Lock()
				sharedErr.err = err
				sharedErr.mutex.Unlock()
			}
		}
	}
}

func walk(dir string, observer DeployObserver, useLargeMedia bool) (*deployFiles, error) {
	files := newDeployFiles()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Mode().IsRegular() {
			osRel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			rel := forceSlashSeparators(osRel)

			if ignoreFile(rel) {
				return nil
			}

			o, err := os.Open(path)
			if err != nil {
				return err
			}
			defer o.Close()

			file := &FileBundle{
				Name: rel,
				Path: path,
			}

			s := sha1.New()
			if _, err := io.Copy(s, o); err != nil {
				return err
			}
			file.Sum = hex.EncodeToString(s.Sum(nil))

			if useLargeMedia {
				o, err := os.Open(path)
				if err != nil {
					return err
				}
				defer o.Close()

				data, err := readLFSData(o)
				if err != nil {
					return err
				}

				if data != nil {
					if data.SHA != "" {
						file.Sum += ":" + data.SHA
					}
					if data.Size > 0 {
						file.Size = &data.Size
					}
				}
			}

			files.Add(rel, file)

			if observer != nil {
				if err := observer.OnSuccessfulStep(file); err != nil {
					return err
				}
			}
		}

		return nil
	})
	return files, err
}

func bundle(functionDir string, observer DeployObserver) (*deployFiles, error) {
	if functionDir == "" {
		return nil, nil
	}

	functions := newDeployFiles()

	info, err := ioutil.ReadDir(functionDir)
	if err != nil {
		return nil, err
	}

	for _, i := range info {
		filePath := filepath.Join(functionDir, i.Name())

		switch {
		case zipFile(i):
			file, err := newFunctionFile(filePath, i, jsRuntime, observer)
			if err != nil {
				return nil, err
			}
			functions.Add(file.Name, file)
		case jsFile(i):
			file, err := newFunctionFile(filePath, i, jsRuntime, observer)
			if err != nil {
				return nil, err
			}
			functions.Add(file.Name, file)
		case goFile(filePath, i):
			file, err := newFunctionFile(filePath, i, goRuntime, observer)
			if err != nil {
				return nil, err
			}
			functions.Add(file.Name, file)
		}
	}

	return functions, nil
}

func newFunctionFile(filePath string, i os.FileInfo, runtime string, observer DeployObserver) (*FileBundle, error) {
	file := &FileBundle{
		Name:    strings.TrimSuffix(i.Name(), filepath.Ext(i.Name())),
		Runtime: runtime,
	}

	s := sha256.New()

	fileEntry, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer fileEntry.Close()

	var buf io.ReadWriter

	if zipFile(i) {
		buf = fileEntry
	} else {
		buf = new(bytes.Buffer)
		archive := zip.NewWriter(buf)

		fileHeader, err := createHeader(archive, i, runtime)
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(fileHeader, fileEntry); err != nil {
			return nil, err
		}

		if err := archive.Close(); err != nil {
			return nil, err
		}
	}

	fileBuffer := new(bytes.Buffer)
	m := io.MultiWriter(s, fileBuffer)

	if _, err := io.Copy(m, buf); err != nil {
		return nil, err
	}
	file.Sum = hex.EncodeToString(s.Sum(nil))
	file.Buffer = bytes.NewReader(fileBuffer.Bytes())

	if observer != nil {
		if err := observer.OnSuccessfulStep(file); err != nil {
			return nil, err
		}
	}

	return file, nil
}

func zipFile(i os.FileInfo) bool {
	return filepath.Ext(i.Name()) == ".zip"
}

func jsFile(i os.FileInfo) bool {
	return filepath.Ext(i.Name()) == ".js"
}

func goFile(filePath string, i os.FileInfo) bool {
	if m := i.Mode(); m&0111 == 0 { // check if it's an executable file
		return false
	}

	if _, err := elf.Open(filePath); err != nil { // check if it's a linux executable
		return false
	}

	v, err := version.ReadExe(filePath)
	return err == nil && strings.HasPrefix(v.Release, "go1.")
}

func ignoreFile(rel string) bool {
	if strings.HasPrefix(rel, ".") || strings.Contains(rel, "/.") || strings.HasPrefix(rel, "__MACOS") {
		return !strings.HasPrefix(rel, ".well-known/")
	}
	return false
}

func createHeader(archive *zip.Writer, i os.FileInfo, runtime string) (io.Writer, error) {
	if runtime == goRuntime {
		return archive.CreateHeader(&zip.FileHeader{
			CreatorVersion: 3 << 8,     // indicates Unix
			ExternalAttrs:  0777 << 16, // -rwxrwxrwx file permissions
			Name:           i.Name(),
			Method:         zip.Deflate,
		})
	}
	return archive.Create(i.Name())
}

func readLFSData(file io.Reader) (*pointerData, error) {
	// currently this only supports certain type of git lfs pointer files
	// version [version]\noid sha256:[oid]\nsize [size]
	data := make([]byte, len(lfsVersionString))
	count, err := file.Read(data)
	if err != nil {
		// ignore file if it's not an LFS pointer with the expected header
		return nil, nil
	}
	if count != len(lfsVersionString) || string(data) != lfsVersionString {
		// ignore file if it's not an LFS pointer with the expected header
		return nil, nil
	}

	scanner := bufio.NewScanner(file)
	values := map[string]string{}
	for scanner.Scan() {
		keyAndValue := bytes.SplitN(scanner.Bytes(), []byte(" "), 2)
		if len(keyAndValue) > 1 {
			values[string(keyAndValue[0])] = string(keyAndValue[1])
		}
	}

	var sha string
	oid, ok := values["oid"]
	if !ok {
		return nil, fmt.Errorf("missing LFS OID")
	}

	sha = strings.SplitN(oid, ":", 2)[1]

	size, err := strconv.ParseInt(values["size"], 10, 0)
	if err != nil {
		return nil, err
	}

	return &pointerData{
		SHA:  sha,
		Size: size,
	}, nil
}
