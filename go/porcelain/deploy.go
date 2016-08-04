package porcelain

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	logContext "github.com/docker/distribution/context"
	"github.com/netlify/open-api/go/models"
	"github.com/netlify/open-api/go/plumbing/operations"
	"github.com/netlify/open-api/go/porcelain/context"
)

const (
	maxFilesForSyncDeploy = 7000
	preProcessingTimeout  = time.Minute * 5
)

type uploadError struct {
	err   error
	mutex *sync.Mutex
}

type file struct {
	Name   string
	SHA1   hash.Hash
	Buffer *bytes.Buffer
}

func (f *file) Sum() string {
	return hex.EncodeToString(f.SHA1.Sum(nil))
}

func (f *file) Read(p []byte) (n int, err error) {
	return f.Buffer.Read(p)
}

func (f *file) Close() error {
	return nil
}

type deployFiles struct {
	Files  map[string]*file
	Sums   map[string]string
	Hashed map[string]*file
}

func newDeployFiles() *deployFiles {
	return &deployFiles{
		Files:  make(map[string]*file),
		Sums:   make(map[string]string),
		Hashed: make(map[string]*file),
	}
}

func (d *deployFiles) Add(p string, f *file) {
	sum := f.Sum()

	d.Files[p] = f
	d.Sums[p] = sum
	d.Hashed[sum] = f
}

func (d *deployFiles) OverCommitted() bool {
	return len(d.Files) > maxFilesForSyncDeploy
}

// DeploySite creates a new deploy for a site given a directory in the filesystem.
// It uploads the necessary files that changed between deploys.
func (n *Netlify) DeploySite(ctx context.Context, siteID, dir string) (*models.Deploy, error) {
	f, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	files, err := walk(dir)
	if err != nil {
		return nil, err
	}

	return n.createDeploy(ctx, siteID, files)
}

func (n *Netlify) createDeploy(ctx context.Context, siteID string, files *deployFiles) (*models.Deploy, error) {
	deployFiles := &models.DeployFiles{
		Files: files.Sums,
		Async: files.OverCommitted(),
	}

	logContext.GetLoggerWithFields(ctx, context.Fields{"site_id": siteID, "deploy_files": len(files.Sums)}).Debug("Deploy files")
	authInfo := context.GetAuthWriter(ctx)

	params := operations.NewCreateSiteDeployParams().WithSiteID(siteID).WithDeploy(deployFiles)
	resp, err := n.Operations.CreateSiteDeploy(params, authInfo)
	if err != nil {
		return nil, err
	}

	deploy := resp.Payload
	if files.OverCommitted() {
		var err error
		deploy, err = n.WaitUntilDeployReady(ctx, deploy)
		if err != nil {
			return nil, err
		}
	}

	if err := n.uploadFiles(ctx, deploy, files); err != nil {
		return nil, err
	}

	return deploy, nil
}

func (n *Netlify) WaitUntilDeployReady(ctx context.Context, d *models.Deploy) (*models.Deploy, error) {
	authInfo := context.GetAuthWriter(ctx)
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

		logContext.GetLoggerWithFields(ctx, context.Fields{"deploy_id": d.ID, "state": resp.Payload.State}).Debug("Wait until deploy ready")

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

func (n *Netlify) uploadFiles(ctx context.Context, d *models.Deploy, files *deployFiles) error {
	sharedErr := &uploadError{err: nil, mutex: &sync.Mutex{}}
	sem := make(chan int, 10) // FIXME(david): make max concurrent uploads configurable.
	wg := &sync.WaitGroup{}

	for _, sha := range d.Required {
		if file, exist := files.Hashed[sha]; exist {
			sem <- 1
			wg.Add(1)

			go n.uploadFile(ctx, d, file, wg, sem, sharedErr)
		}
	}

	wg.Wait()

	return sharedErr.err
}

func (n *Netlify) uploadFile(ctx context.Context, d *models.Deploy, f *file, wg *sync.WaitGroup, sem chan int, sharedErr *uploadError) {
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

	authInfo := context.GetAuthWriter(ctx)

	logContext.GetLoggerWithFields(ctx, context.Fields{"deploy_id": d.ID, "file_path": f.Name, "file_sum": f.Sum()}).Debug("Upload file")

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 2 * time.Minute

	err := backoff.Retry(func() error {
		sharedErr.mutex.Lock()

		if sharedErr.err != nil {
			sharedErr.mutex.Unlock()
			return fmt.Errorf("Upload cancelled: %s", f.Name)
		}
		sharedErr.mutex.Unlock()

		params := operations.NewUploadDeployFileParams().WithDeployID(d.ID).WithPath(f.Name).WithFileBody(f)
		_, err := n.Operations.UploadDeployFile(params, authInfo)

		if err != nil {
			logContext.GetLogger(ctx).Error(err)
		}

		return err
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

			o, err := os.Open(rel)
			if err != nil {
				return err
			}

			file := &file{
				Name:   rel,
				SHA1:   sha1.New(),
				Buffer: new(bytes.Buffer),
			}
			m := io.MultiWriter(file.SHA1, file.Buffer)

			if _, err := io.Copy(m, o); err != nil {
				return err
			}

			files.Add(rel, file)
		}

		return nil
	})
	return files, err
}

func ignoreFile(rel string) bool {
	if strings.HasPrefix(rel, ".") || strings.Contains(rel, "/.") || strings.HasPrefix(rel, "__MACOS") {
		return !strings.HasPrefix(rel, ".well-known/")
	}
	return false
}
