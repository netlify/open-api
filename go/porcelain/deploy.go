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
)

const (
	preProcessingTimeout = time.Minute * 5
)

type uploadType int

const (
	fileUpload uploadType = iota
	functionUpload
)

// DeployOptions holds the option for creating a new deploy
type DeployOptions struct {
	SiteID       string
	Dir          string
	FunctionsDir string

	IsDraft bool

	Title     string
	Branch    string
	CommitRef string

	files     *deployFiles
	functions *deployFiles
}

type uploadError struct {
	err   error
	mutex *sync.Mutex
}

type file struct {
	Name   string
	SHA    hash.Hash
	Buffer *bytes.Reader
}

func (f *file) Sum() string {
	return hex.EncodeToString(f.SHA.Sum(nil))
}

func (f *file) Read(p []byte) (n int, err error) {
	return f.Buffer.Read(p)
}

func (f *file) Close() error {
	return nil
}

func (f *file) Rewind() error {
	_, err := f.Buffer.Seek(0, 0)
	return err
}

type deployFiles struct {
	Files  map[string]*file
	Sums   map[string]string
	Hashed map[string][]*file
}

func newDeployFiles() *deployFiles {
	return &deployFiles{
		Files:  make(map[string]*file),
		Sums:   make(map[string]string),
		Hashed: make(map[string][]*file),
	}
}

func (d *deployFiles) Add(p string, f *file) {
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

	files, err := walk(options.Dir)
	if err != nil {
		return nil, err
	}
	options.files = files

	functions, err := bundle(options.FunctionsDir)
	if err != nil {
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

	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id":      options.SiteID,
		"deploy_files": len(options.files.Sums),
	}).Debug("Starting to deploy files")
	authInfo := context.GetAuthInfo(ctx)

	if deploy == nil {
		params := operations.NewCreateSiteDeployParams().WithSiteID(options.SiteID).WithDeploy(deployFiles)
		if options.Title != "" {
			params = params.WithTitle(&options.Title)
		}
		resp, err := n.Operations.CreateSiteDeploy(params, authInfo)
		if err != nil {
			return nil, err
		}
		deploy = resp.Payload
	} else {
		params := operations.NewUpdateSiteDeployParams().WithSiteID(options.SiteID).WithDeployID(deploy.ID).WithDeploy(deployFiles)
		resp, err := n.Operations.UpdateSiteDeploy(params, authInfo)
		if err != nil {
			return nil, err
		}
		deploy = resp.Payload
	}

	if n.overCommitted(options.files) {
		var err error
		deploy, err = n.WaitUntilDeployReady(ctx, deploy)
		if err != nil {
			return nil, err
		}
	}

	if len(deploy.Required) == 0 && len(deploy.RequiredFunctions) == 0 {
		return deploy, nil
	}

	if err := n.uploadFiles(ctx, deploy, options.files, fileUpload); err != nil {
		return nil, err
	}

	if options.functions != nil {
		if err := n.uploadFiles(ctx, deploy, options.functions, functionUpload); err != nil {
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

func (n *Netlify) uploadFiles(ctx context.Context, d *models.Deploy, files *deployFiles, t uploadType) error {
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

	for _, sha := range required {
		if files, exist := files.Hashed[sha]; exist {
			for _, file := range files {
				sem <- 1
				wg.Add(1)

				go n.uploadFile(ctx, d, file, t, wg, sem, sharedErr)
			}
		}
	}

	wg.Wait()

	return sharedErr.err
}

func (n *Netlify) uploadFile(ctx context.Context, d *models.Deploy, f *file, t uploadType, wg *sync.WaitGroup, sem chan int, sharedErr *uploadError) {
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

	err := backoff.Retry(func() error {
		sharedErr.mutex.Lock()

		if sharedErr.err != nil {
			sharedErr.mutex.Unlock()
			return fmt.Errorf("Upload cancelled: %s", f.Name)
		}
		sharedErr.mutex.Unlock()

		var operationError error

		switch t {
		case fileUpload:
			params := operations.NewUploadDeployFileParams().WithDeployID(d.ID).WithPath(f.Name).WithFileBody(f)
			_, operationError = n.Operations.UploadDeployFile(params, authInfo)
		case functionUpload:
			params := operations.NewUploadDeployFunctionParams().WithDeployID(d.ID).WithName(f.Name).WithFileBody(f)
			_, operationError = n.Operations.UploadDeployFunction(params, authInfo)
		}

		if operationError != nil {
			f.Rewind()
			context.GetLogger(ctx).WithError(operationError).Error("Failed to upload file")
		}

		return operationError
	}, b)

	if err != nil {
		sharedErr.mutex.Lock()
		sharedErr.err = err
		sharedErr.mutex.Unlock()
	}
}

func walk(dir string) (*deployFiles, error) {
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

			buf := new(bytes.Buffer)
			file := &file{
				Name: rel,
				SHA:  sha1.New(),
			}
			m := io.MultiWriter(file.SHA, buf)

			if _, err := io.Copy(m, o); err != nil {
				return err
			}

			file.Buffer = bytes.NewReader(buf.Bytes())
			files.Add(rel, file)
		}

		return nil
	})
	return files, err
}

func bundle(functionDir string) (*deployFiles, error) {
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
			file := &file{
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

			file.Buffer = bytes.NewReader(fileBuffer.Bytes())
			functions.Add(file.Name, file)
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
