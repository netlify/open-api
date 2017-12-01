package porcelain

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"

	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"

	"github.com/go-openapi/errors"
)

const (
	preProcessingTimeout = time.Minute * 5
)

type uploadType int

const (
	fileUpload uploadType = iota
	functionUpload
)

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
	SiteID       string
	Dir          string
	FunctionsDir string

	IsDraft bool

	Title         string
	Branch        string
	CommitRef     string
	UploadTimeout time.Duration

	Observer DeployObserver

	files     *deployFiles
	functions *deployFiles
}

type uploadError struct {
	err   error
	mutex *sync.Mutex
}

type FileBundle struct {
	Name string
	SHA  hash.Hash

	Buffer io.ReadSeeker
}

func (f *FileBundle) Sum() string {
	return hex.EncodeToString(f.SHA.Sum(nil))
}

func (f *FileBundle) Read(p []byte) (n int, err error) {
	return f.Buffer.Read(p)
}

func (f *FileBundle) Close() error {
	return nil
}

// We're mocking up a closer, to make sure the underlying file handle
// doesn't get closed during an upload, but can be rewinded for retries
// This method closes the file handle for real.
func (f *FileBundle) CloseForReal() error {
	closer, ok := f.Buffer.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}

func (f *FileBundle) Rewind() error {
	_, err := f.Buffer.Seek(0, 0)
	return err
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
	sum := f.Sum()

	d.Files[p] = f
	d.Sums[p] = sum
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

	files, err := walk(options.Dir, options.Observer)
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
		deploy, err = n.WaitUntilDeployReady(ctx, deploy)
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

func (n *Netlify) WaitUntilDeployReady(ctx context.Context, d *models.Deploy) (*models.Deploy, error) {
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

		if resp.Payload.State == "prepared" || resp.Payload.State == "ready" {
			return resp.Payload, nil
		}

		if resp.Payload.State == "error" {
			return nil, fmt.Errorf("Error: preprocessing deploy failed")
		}

		if t.Sub(start) > preProcessingTimeout {
			return nil, fmt.Errorf("Error: preprocessing deploy timed out")
		}
	}

	return d, nil
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
		"file_sum":  f.Sum(),
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
			name := f.Name
			params := operations.NewUploadDeployFileParams().WithDeployID(d.ID).WithPath(name).WithFileBody(f)
			if timeout != 0 {
				params.SetTimeout(timeout)
			}
			_, operationError = n.Operations.UploadDeployFile(params, authInfo)
		case functionUpload:
			params := operations.NewUploadDeployFunctionParams().WithDeployID(d.ID).WithName(f.Name).WithFileBody(f)
			if timeout != 0 {
				params.SetTimeout(timeout)
			}
			_, operationError = n.Operations.UploadDeployFunction(params, authInfo)
		}

		if operationError != nil {
			f.Rewind()
			context.GetLogger(ctx).WithError(operationError).Errorf("Failed to upload file %v", f.Name)
			apiErr, ok := operationError.(errors.Error)

			if ok && apiErr.Code() == 401 {
				sharedErr.mutex.Lock()
				sharedErr.err = operationError
				sharedErr.mutex.Unlock()
			}
		} else {
			f.CloseForReal()
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

func walk(dir string, observer DeployObserver) (*deployFiles, error) {
	files := newDeployFiles()

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Mode().IsRegular() {
			rel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			if ignoreFile(rel) {
				return nil
			}

			o, err := os.Open(path)
			if err != nil {
				return err
			}

			file := &FileBundle{
				Name: rel,
				SHA:  sha1.New(),
			}

			if _, err := io.Copy(file.SHA, o); err != nil {
				return err
			}

			o.Seek(0, 0)
			file.Buffer = o
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
		switch filepath.Ext(i.Name()) {
		case ".js":
			file := &FileBundle{
				Name: strings.TrimSuffix(i.Name(), filepath.Ext(i.Name())),
				SHA:  sha256.New(),
			}

			buf := new(bytes.Buffer)
			archive := zip.NewWriter(buf)
			fileHeader, err := archive.Create(i.Name())
			if err != nil {
				return nil, err
			}
			fileEntry, err := os.Open(filepath.Join(functionDir, i.Name()))
			if err != nil {
				return nil, err
			}
			if _, err = io.Copy(fileHeader, fileEntry); err != nil {
				return nil, err
			}

			if err := archive.Close(); err != nil {
				return nil, err
			}

			fileBuffer := new(bytes.Buffer)
			m := io.MultiWriter(file.SHA, fileBuffer)

			if _, err := io.Copy(m, buf); err != nil {
				return nil, err
			}

			fileEntry.Seek(0, 0)
			file.Buffer = bytes.NewReader(fileBuffer.Bytes())
			functions.Add(file.Name, file)

			if observer != nil {
				if err := observer.OnSuccessfulStep(file); err != nil {
					return nil, err
				}
			}
		default:
			// Ignore this file
		}
	}

	return functions, nil
}

func ignoreFile(rel string) bool {
	if strings.HasPrefix(rel, ".") || strings.Contains(rel, "/.") || strings.HasPrefix(rel, "__MACOS") {
		return !strings.HasPrefix(rel, ".well-known/")
	}
	return false
}
