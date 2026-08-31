package porcelain

import (
	"archive/zip"
	"bufio"
	"bytes"
	gocontext "context"
	"crypto/sha1"
	"crypto/sha256"
	"debug/elf"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/pkg/errors"
	"github.com/rsc/goversion/version"
	"github.com/sirupsen/logrus"

	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
	"github.com/netlify/open-api/v2/go/porcelain/context"
)

const (
	jsRuntime    = "js"
	goRuntime    = "go"
	amazonLinux2 = "provided.al2"

	preProcessingTimeout = time.Minute * 5

	fileUpload uploadType = iota
	functionUpload
	edgeFunctionUpload

	lfsVersionString = "version https://git-lfs.github.com/spec/v1"

	edgeFunctionsInternalPath = ".netlify/internal/edge-functions/"
	edgeRedirectsInternalPath = ".netlify/deploy-config/"
	dbMigrationsInternalPath  = ".netlify/internal/db/migrations/"
)

var installDirs = []string{"node_modules/", "bower_components/"}

type (
	uploadType  int
	pointerData struct {
		SHA  string
		Size int64
	}
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

type DeployWarner interface {
	OnWalkWarning(path, msg string)
}

// DeployOptions holds the option for creating a new deploy
type DeployOptions struct {
	SiteID            string
	Dir               string
	FunctionsDir      string
	EdgeFunctionsDir  string
	EdgeRedirectsDir  string
	DbMigrationsDir   string
	BuildDir          string
	LargeMediaEnabled bool
	Environment       []*models.DeployEnvironmentVariable

	// DirRoot and friends are optional pre-opened handles for the corresponding
	// *Dir path fields. When set, all filesystem access for that directory goes
	// through the handle; when nil, the path field is opened with [os.OpenRoot]
	// once at the start of the deploy. Caller-provided handles must stay open for
	// the duration of the deploy and are not closed by this package. The path
	// fields should still be set: they are used for logging and to resolve paths
	// read from manifest files.
	DirRoot              *os.Root
	FunctionsDirRoot     *os.Root
	EdgeFunctionsDirRoot *os.Root
	EdgeRedirectsDirRoot *os.Root
	DbMigrationsDirRoot  *os.Root

	IsDraft   bool
	SkipRetry bool

	Title             string
	Branch            string
	CommitRef         string
	Framework         string
	FrameworkVersion  string
	UploadTimeout     time.Duration
	PreProcessTimeout time.Duration

	Observer DeployObserver

	files             *deployFiles
	functions         *deployFiles
	edgeFunctions     *deployFiles
	functionSchedules []*models.FunctionSchedule
	functionsConfig   map[string]models.FunctionConfig
}

type deployApiError interface {
	error
	Code() int
}

type uploadError struct {
	err   error
	mutex *sync.Mutex
}

type FileBundle struct {
	Name             string
	Sum              string
	Runtime          string
	Size             *int64 `json:"size,omitempty"`
	FunctionMetadata *FunctionMetadata

	// Path is the location of the file on disk. Uploads always stream from Path.
	Path string

	// Deprecated: uploads always stream from Path; this package no longer reads Buffer. It is retained
	// only for backwards compatibility with external callers and may be removed in a future release.
	// Leave it nil to have the (also deprecated) Read/Seek/Close methods stream from Path instead.
	Buffer io.ReadSeeker

	// pathReader is lazily opened from Path when Buffer is nil, so the deprecated Read/Seek/Close
	// methods keep working for external callers that treat a FileBundle as an io.ReadSeekCloser.
	pathReader *os.File

	// root is the directory handle the file lives in and rel its path within that
	// handle; every read of the file's contents goes through root.
	root *os.Root
	rel  string
}

// open returns a reader for the bundle's contents, always through the root
// handle the bundle was created with.
func (f *FileBundle) open() (*os.File, error) {
	if f.root == nil {
		return nil, fmt.Errorf("file bundle %s has no root handle", f.Name)
	}
	return openRegularFileInRoot(f.root, f.rel)
}

// legacyOpen backs the deprecated Read/Seek methods only: a FileBundle
// constructed by hand by an external caller has no root handle and keeps the
// historical direct open of its Path.
func (f *FileBundle) legacyOpen() (*os.File, error) {
	if f.root == nil {
		return os.Open(f.Path)
	}
	return f.open()
}

type FunctionMetadata struct {
	InvocationMode string
	Timeout        int64
}

type toolchainSpec struct {
	Runtime string `json:"runtime"`
}

// Deprecated: read directly from Path (e.g. via os.Open) instead. When Buffer is set, Read reads
// from it; otherwise it streams from Path. Retained for backwards compatibility with external
// callers and may be removed in a future release.
func (f *FileBundle) Read(p []byte) (n int, err error) {
	if f.Buffer != nil {
		return f.Buffer.Read(p)
	}
	if f.pathReader == nil {
		if f.pathReader, err = f.legacyOpen(); err != nil {
			return 0, err
		}
	}
	return f.pathReader.Read(p)
}

// Deprecated: read directly from Path (e.g. via os.Open) instead. When Buffer is set, Seek seeks
// it; otherwise it seeks the stream opened from Path. Retained for backwards compatibility with
// external callers and may be removed in a future release.
func (f *FileBundle) Seek(offset int64, whence int) (int64, error) {
	if f.Buffer != nil {
		return f.Buffer.Seek(offset, whence)
	}
	if f.pathReader == nil {
		var err error
		if f.pathReader, err = f.legacyOpen(); err != nil {
			return 0, err
		}
	}
	return f.pathReader.Seek(offset, whence)
}

// Deprecated: retained for backwards compatibility with external callers and may be removed in a
// future release. It closes the stream lazily opened from Path by Read/Seek; it never closes a
// caller-supplied Buffer.
func (f *FileBundle) Close() error {
	if f.pathReader != nil {
		err := f.pathReader.Close()
		f.pathReader = nil
		return err
	}
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

// dirHandle pairs an open directory handle with its configured path, which is
// used for messages, FileBundle.Path, and resolving absolute manifest paths.
type dirHandle struct {
	root *os.Root
	name string
}

func (h dirHandle) valid() bool {
	return h.root != nil
}

// deployRoots holds the resolved directory handles for one deploy. Handles
// opened here (rather than provided by the caller) are recorded in owned and
// closed when the deploy returns.
type deployRoots struct {
	dir, functions, edgeFunctions, edgeRedirects, dbMigrations dirHandle

	owned []*os.Root
}

func (r *deployRoots) close() {
	for _, root := range r.owned {
		_ = root.Close()
	}
}

// resolveRoot returns the caller-provided handle when set, otherwise opens one
// for path; an empty path means the directory takes no part in the deploy.
func (r *deployRoots) resolveRoot(handle *os.Root, path string) (dirHandle, error) {
	if handle != nil {
		name := path
		if name == "" {
			name = handle.Name()
		}
		return dirHandle{root: handle, name: name}, nil
	}
	if path == "" {
		return dirHandle{}, nil
	}
	if fi, err := os.Lstat(path); err != nil {
		return dirHandle{}, err
	} else if fi.Mode()&os.ModeSymlink != 0 {
		return dirHandle{}, fmt.Errorf("%s is a symbolic link", path)
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return dirHandle{}, err
	}
	r.owned = append(r.owned, root)
	return dirHandle{root: root, name: path}, nil
}

func resolveDeployRoots(options *DeployOptions) (*deployRoots, error) {
	roots := &deployRoots{}

	resolve := func(dst *dirHandle, handle *os.Root, path string) error {
		h, err := roots.resolveRoot(handle, path)
		if err != nil {
			roots.close()
			return err
		}
		*dst = h
		return nil
	}

	// Keep the historical error message for a path that is not a directory.
	if options.DirRoot == nil {
		f, err := os.Stat(options.Dir)
		if err != nil {
			return nil, err
		}
		if !f.IsDir() {
			return nil, fmt.Errorf("%s is not a directory", options.Dir)
		}
	}

	if err := resolve(&roots.dir, options.DirRoot, options.Dir); err != nil {
		return nil, err
	}
	if !roots.dir.valid() {
		return nil, fmt.Errorf("no deploy directory provided")
	}
	if err := resolve(&roots.functions, options.FunctionsDirRoot, options.FunctionsDir); err != nil {
		return nil, err
	}
	if err := resolve(&roots.edgeFunctions, options.EdgeFunctionsDirRoot, options.EdgeFunctionsDir); err != nil {
		return nil, err
	}
	if err := resolve(&roots.edgeRedirects, options.EdgeRedirectsDirRoot, options.EdgeRedirectsDir); err != nil {
		return nil, err
	}
	if err := resolve(&roots.dbMigrations, options.DbMigrationsDirRoot, options.DbMigrationsDir); err != nil {
		return nil, err
	}

	return roots, nil
}

// DoDeploy deploys the changes for a site given a directory in the filesystem.
// It uploads the necessary files that changed between deploys.
func (n *Netlify) DoDeploy(ctx context.Context, options *DeployOptions, deploy *models.Deploy) (*models.Deploy, error) {
	roots, err := resolveDeployRoots(options)
	if err != nil {
		return nil, err
	}
	// The upload phase re-reads every required file through these handles.
	defer roots.close()

	if options.Observer != nil {
		if err := options.Observer.OnSetupWalk(); err != nil {
			return nil, err
		}
	}

	largeMediaEnabled := options.LargeMediaEnabled
	ignoreInstallDirs := options.Dir != "" && options.Dir == options.BuildDir

	context.GetLogger(ctx).Infof("Getting files info with large media flag: %v", largeMediaEnabled)
	files, err := walk(roots.dir, options.Observer, largeMediaEnabled, ignoreInstallDirs)
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

	if roots.edgeFunctions.valid() {
		err = addInternalFilesToDeploy(roots.edgeFunctions, edgeFunctionsInternalPath, files, options.Observer)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedWalk()
			}
			return nil, err
		}
	}

	if roots.edgeRedirects.valid() {
		err = addInternalFilesToDeploy(roots.edgeRedirects, edgeRedirectsInternalPath, files, options.Observer)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedWalk()
			}
			return nil, err
		}
	}

	if roots.dbMigrations.valid() {
		err = addInternalFilesToDeploy(roots.dbMigrations, dbMigrationsInternalPath, files, options.Observer)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedWalk()
			}
			return nil, err
		}
	}

	options.files = files

	// The temp dir is created lazily, only if a function actually needs to be zipped. Pre-bundled
	// .zip/.tar functions stream from their original path and never touch it, so a deploy with no
	// unbundled functions creates no temp dir at all.
	functionsTmpDir := &lazyTempDir{}
	defer functionsTmpDir.remove()

	functions, schedules, functionsConfig, err := bundle(ctx, roots.functions, functionsTmpDir, options.Observer)
	if err != nil {
		if options.Observer != nil {
			options.Observer.OnFailedWalk()
		}
		return nil, err
	}
	options.functions = functions
	options.functionSchedules = schedules
	options.functionsConfig = functionsConfig

	edgeFunctions, err := bundleEdgeFunctions(ctx, roots.edgeFunctions, options.Observer)
	if err != nil {
		if options.Observer != nil {
			options.Observer.OnFailedWalk()
		}
		return nil, err
	}
	options.edgeFunctions = edgeFunctions

	deployFiles := &models.DeployFiles{
		Files:            options.files.Sums,
		Draft:            options.IsDraft,
		Async:            n.overCommitted(options.files),
		Framework:        options.Framework,
		FrameworkVersion: options.FrameworkVersion,
	}
	if options.functions != nil {
		deployFiles.Functions = options.functions.Sums
	}
	if options.edgeFunctions != nil {
		deployFiles.EdgeFunctions = options.edgeFunctions.Sums
	}

	if len(options.Environment) > 0 {
		deployFiles.Environment = options.Environment
	}

	if options.Observer != nil {
		if err := options.Observer.OnSuccessfulWalk(deployFiles); err != nil {
			return nil, err
		}
	}

	if len(schedules) > 0 {
		deployFiles.FunctionSchedules = schedules
	}

	if options.functionsConfig != nil {
		deployFiles.FunctionsConfig = options.functionsConfig
	}

	l := context.GetLogger(ctx)
	l.WithFields(logrus.Fields{
		"site_id":             options.SiteID,
		"deploy_files":        len(options.files.Sums),
		"scheduled_functions": len(schedules),
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

		timeout := options.PreProcessTimeout
		if timeout <= 0 {
			timeout = preProcessingTimeout
		}
		deployReadyCtx, _ := gocontext.WithTimeout(ctx, timeout)
		deploy, err = n.WaitUntilDeployReady(deployReadyCtx, deploy)
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

	if len(deploy.Required) == 0 && len(deploy.RequiredFunctions) == 0 && len(deploy.RequiredEdgeFunctions) == 0 {
		return deploy, nil
	}

	skipRetry := options.SkipRetry

	if err := n.uploadFiles(ctx, deploy, options.files, options.Observer, fileUpload, options.UploadTimeout, skipRetry); err != nil {
		return nil, err
	}

	if options.functions != nil {
		if err := n.uploadFiles(ctx, deploy, options.functions, options.Observer, functionUpload, options.UploadTimeout, skipRetry); err != nil {
			return nil, err
		}
	}

	if options.edgeFunctions != nil {
		if err := n.uploadFiles(ctx, deploy, options.edgeFunctions, options.Observer, edgeFunctionUpload, options.UploadTimeout, skipRetry); err != nil {
			return nil, err
		}
	}

	return deploy, nil
}

func (n *Netlify) waitForState(ctx context.Context, d *models.Deploy, states ...string) (*models.Deploy, error) {
	authInfo := context.GetAuthInfo(ctx)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	params := operations.NewGetSiteDeployParams().WithSiteID(d.SiteID).WithDeployID(d.ID)
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timed out while waiting to enter states [%s]", strings.Join(states, ", "))
		case <-ticker.C:
			resp, err := n.Operations.GetSiteDeploy(params, authInfo)
			if err != nil {
				time.Sleep(3 * time.Second)
				continue
			}
			context.GetLogger(ctx).WithFields(logrus.Fields{
				"deploy_id": d.ID,
				"state":     resp.Payload.State,
			}).Debugf("Waiting until deploy state in %s", states)

			for _, state := range states {
				if resp.Payload.State == state {
					return resp.Payload, nil
				}
			}

			if resp.Payload.State == "error" {
				return nil, fmt.Errorf("entered error state while waiting to enter states [%s]", strings.Join(states, ", "))
			}
		}
	}
}

// WaitUntilDeployReady blocks until the deploy is in the "prepared" or "ready" state.
func (n *Netlify) WaitUntilDeployReady(ctx context.Context, d *models.Deploy) (*models.Deploy, error) {
	return n.waitForState(ctx, d, "prepared", "ready")
}

// WaitUntilDeployLive blocks until the deploy is in the "ready" state. At this point, the deploy is ready to receive traffic to all of its URLs.
func (n *Netlify) WaitUntilDeployLive(ctx context.Context, d *models.Deploy) (*models.Deploy, error) {
	return n.waitForState(ctx, d, "ready")
}

// WaitUntilDeployProcessed blocks until the deploy is in the "processed" state. At this point, the deploy is ready to receive traffic via its permalink.
func (n *Netlify) WaitUntilDeployProcessed(ctx context.Context, d *models.Deploy) (*models.Deploy, error) {
	return n.waitForState(ctx, d, "processed")
}

func (n *Netlify) uploadFiles(ctx context.Context, d *models.Deploy, files *deployFiles, observer DeployObserver, t uploadType, timeout time.Duration, skipRetry bool) error {
	sharedErr := &uploadError{err: nil, mutex: &sync.Mutex{}}
	sem := make(chan int, n.uploadLimit)
	wg := &sync.WaitGroup{}

	var required []string
	switch t {
	case fileUpload:
		required = d.Required
	case functionUpload:
		required = d.RequiredFunctions
	case edgeFunctionUpload:
		required = d.RequiredEdgeFunctions
	}

	count := 0
	for _, sha := range required {
		if files, exist := files.Hashed[sha]; exist {
			count += len(files)
		}
	}

	log := context.GetLogger(ctx)
	log.Infof("Uploading %v files", count)

	var abortErr error
	for _, sha := range required {
		if files, exist := files.Hashed[sha]; exist {
			file := files[0]

			select {
			case sem <- 1:
				wg.Add(1)
				go n.uploadFile(ctx, d, file, observer, t, timeout, wg, sem, sharedErr, skipRetry)
			case <-ctx.Done():
				log.Info("Context terminated, aborting file upload")
				abortErr = errors.Wrap(ctx.Err(), "aborted file upload early")
			}

			if abortErr != nil {
				break
			}

			if len(files) > 1 {
				skippedFiles := files[1:]
				for _, file := range skippedFiles {
					log.Infof("Skipping file with content already uploaded: %s", file.Name)
				}
			}
		}
	}

	// Always wait for in-flight uploads to finish before returning. On the ctx.Done()
	// path this prevents orphaned uploadFile goroutines from racing against the caller's
	// deferred temp-dir cleanup (os.RemoveAll), which would otherwise open files that are
	// being deleted and surface spurious "no such file or directory" errors.
	wg.Wait()

	if abortErr != nil {
		return abortErr
	}

	return sharedErr.err
}

func (n *Netlify) uploadFile(ctx context.Context, d *models.Deploy, f *FileBundle, c DeployObserver, t uploadType, timeout time.Duration, wg *sync.WaitGroup, sem chan int, sharedErr *uploadError, skipRetry bool) {
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

	var retryCount int64 = 0

	err := backoff.Retry(func() error {
		sharedErr.mutex.Lock()

		if sharedErr.err != nil {
			sharedErr.mutex.Unlock()
			return fmt.Errorf("aborting upload of file %s due to failed upload of another file", f.Name)
		}
		sharedErr.mutex.Unlock()

		// Opening the file cannot start succeeding on a retry, so fail permanently
		// rather than backing off for the full retry window.
		body, openErr := f.open()
		if openErr != nil {
			context.GetLogger(ctx).WithError(openErr).Errorf("Failed to open %v for upload", f.Name)
			return backoff.Permanent(openErr)
		}
		defer func() { _ = body.Close() }()

		var operationError error
		switch t {
		case fileUpload:
			params := operations.NewUploadDeployFileParams().WithDeployID(d.ID).WithPath(f.Name).WithFileBody(body)
			if f.Size != nil {
				params.WithSize(f.Size)
			}
			if timeout != 0 {
				params.SetTimeout(timeout)
			}
			_, operationError = n.Operations.UploadDeployFile(params, authInfo)
		case functionUpload:
			params := operations.NewUploadDeployFunctionParams().WithDeployID(d.ID).WithName(f.Name).WithFileBody(body).WithRuntime(&f.Runtime)

			if retryCount > 0 {
				params = params.WithXNfRetryCount(&retryCount)
			}

			if f.FunctionMetadata != nil {
				params = params.WithInvocationMode(&f.FunctionMetadata.InvocationMode)
				params = params.WithTimeout(&f.FunctionMetadata.Timeout)
			}

			if timeout != 0 {
				params.SetRequestTimeout(timeout)
			}
			_, operationError = n.Operations.UploadDeployFunction(params, authInfo)
		case edgeFunctionUpload:
			params := operations.NewUploadDeployEdgeFunctionParams().WithDeployID(d.ID).WithCodeSha(f.Sum).WithFileBody(body)
			if retryCount > 0 {
				params = params.WithXNfRetryCount(&retryCount)
			}
			if timeout != 0 {
				params.SetTimeout(timeout)
			}
			_, operationError = n.Operations.UploadDeployEdgeFunction(params, authInfo)
		}

		if operationError != nil {
			context.GetLogger(ctx).WithError(operationError).Errorf("Failed to upload file %v", f.Name)
			apiErr, ok := operationError.(deployApiError)

			if ok {
				if apiErr.Code() == 401 {
					sharedErr.mutex.Lock()
					sharedErr.err = operationError
					sharedErr.mutex.Unlock()
				}

				if skipRetry && (apiErr.Code() == 400 || apiErr.Code() == 422) {
					operationError = backoff.Permanent(operationError)
				}
			}
		}

		retryCount++

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

func createFileBundle(rel string, dir dirHandle, relPath string) (*FileBundle, error) {
	return createFileBundleWithHasher(rel, dir, relPath, sha1.New())
}

func createFunctionFileBundle(rel string, dir dirHandle, relPath string) (*FileBundle, error) {
	return createFileBundleWithHasher(rel, dir, relPath, sha256.New())
}

// createFileBundleWithHasher builds the bundle for the file at relPath inside
// dir; rel is the name it deploys as.
func createFileBundleWithHasher(rel string, dir dirHandle, relPath string, s hash.Hash) (*FileBundle, error) {
	o, err := openRegularFileInRoot(dir.root, relPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = o.Close() }()

	file := &FileBundle{
		Name: rel,
		Path: filepath.Join(dir.name, relPath),
		root: dir.root,
		rel:  relPath,
	}

	if _, err := io.Copy(s, o); err != nil {
		return nil, err
	}

	file.Sum = hex.EncodeToString(s.Sum(nil))

	return file, nil
}

// openRegularFileInRoot opens relPath inside root for reading, rejecting
// anything but a regular file. O_NONBLOCK avoids blocking the open if relPath
// is a FIFO.
func openRegularFileInRoot(root *os.Root, relPath string) (*os.File, error) {
	f, err := root.OpenFile(filepath.FromSlash(relPath), os.O_RDONLY|openNonblock, 0)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if !fi.Mode().IsRegular() {
		_ = f.Close()
		return nil, fmt.Errorf("%s is not a regular file", relPath)
	}
	return f, nil
}

func walk(dir dirHandle, observer DeployObserver, useLargeMedia, ignoreInstallDirs bool) (*deployFiles, error) {
	files := newDeployFiles()

	err := fs.WalkDir(dir.root.FS(), ".", func(rel string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && d.Type().IsRegular() {
			if ignoreFile(rel, ignoreInstallDirs) {
				return nil
			}

			file, err := createFileBundle(rel, dir, rel)
			if err != nil {
				return err
			}

			if useLargeMedia {
				o, err := openRegularFileInRoot(dir.root, rel)
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

func addInternalFilesToDeploy(dir dirHandle, internalPath string, files *deployFiles, observer DeployObserver) error {
	return fs.WalkDir(dir.root.FS(), ".", func(osRel string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && d.Type().IsRegular() {
			rel := internalPath + osRel

			file, err := createFileBundle(rel, dir, osRel)
			if err != nil {
				return err
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
}

type lazyTempDir struct {
	root    string
	path    string
	handle  *os.Root
	created bool
}

func (l *lazyTempDir) get() (string, *os.Root, error) {
	if !l.created {
		path, err := os.MkdirTemp(l.root, "netlify-deploy-functions-")
		if err != nil {
			return "", nil, err
		}
		handle, err := os.OpenRoot(path)
		if err != nil {
			os.RemoveAll(path)
			return "", nil, err
		}
		l.path, l.handle, l.created = path, handle, true
	}
	return l.path, l.handle, nil
}

func (l *lazyTempDir) remove() {
	if l.created {
		_ = l.handle.Close()
		os.RemoveAll(l.path)
	}
}

func bundle(ctx context.Context, functionsDir dirHandle, tmpDir *lazyTempDir, observer DeployObserver) (*deployFiles, []*models.FunctionSchedule, map[string]models.FunctionConfig, error) {
	if !functionsDir.valid() {
		return nil, nil, nil, nil
	}

	manifestFile, err := openRegularFileInRoot(functionsDir.root, "manifest.json")

	// If a `manifest.json` file is found, we extract the functions and their
	// metadata from it.
	if err == nil {
		defer manifestFile.Close()

		return bundleFromManifest(ctx, functionsDir, manifestFile, tmpDir, observer)
	}

	functions := newDeployFiles()

	info, err := fs.ReadDir(functionsDir.root.FS(), ".")
	if err != nil {
		return nil, nil, nil, err
	}

	for _, entry := range info {
		i, err := entry.Info()
		if err != nil {
			return nil, nil, nil, err
		}

		// filePath is only used for warnings and the go-binary classification.
		filePath := filepath.Join(functionsDir.name, i.Name())

		switch {
		case zipFile(i):
			runtime, err := readZipRuntime(functionsDir.root, i.Name())
			if err != nil {
				return nil, nil, nil, err
			}
			file, err := newFunctionFile(functionsDir, i.Name(), i, runtime, nil, tmpDir, observer)
			if err != nil {
				return nil, nil, nil, err
			}
			functions.Add(file.Name, file)
		case jsFile(i):
			file, err := newFunctionFile(functionsDir, i.Name(), i, jsRuntime, nil, tmpDir, observer)
			if err != nil {
				return nil, nil, nil, err
			}
			functions.Add(file.Name, file)
		case goFile(filePath, i, observer):
			file, err := newFunctionFile(functionsDir, i.Name(), i, amazonLinux2, nil, tmpDir, observer)
			if err != nil {
				return nil, nil, nil, err
			}
			functions.Add(file.Name, file)
		default:
			if warner, ok := observer.(DeployWarner); ok {
				warner.OnWalkWarning(filePath, "Function is not valid for deployment. Please check that it matches the format for the runtime.")
			}
		}
	}

	return functions, nil, nil, nil
}

func bundleFromManifest(ctx context.Context, functionsDir dirHandle, manifestFile *os.File, tmpDir *lazyTempDir, observer DeployObserver) (*deployFiles, []*models.FunctionSchedule, map[string]models.FunctionConfig, error) {
	manifestBytes, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return nil, nil, nil, err
	}

	logger := context.GetLogger(ctx)
	logger.Debug("Found functions manifest file")

	var manifest functionsManifest

	err = json.Unmarshal(manifestBytes, &manifest)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("malformed functions manifest file: %w", err)
	}

	schedules := make([]*models.FunctionSchedule, 0, len(manifest.Functions))
	functions := newDeployFiles()
	functionsConfig := make(map[string]models.FunctionConfig)

	for _, function := range manifest.Functions {
		// The manifest is untrusted input: the paths it names must resolve
		// inside the functions directory.
		relPath, err := manifestFunctionRel(functionsDir.name, function.Path)
		if err != nil {
			return nil, nil, nil, err
		}

		fileInfo, err := functionsDir.root.Stat(relPath)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("manifest file specifies a function path that cannot be found: %s", function.Path)
		}

		var runtime string
		if function.RuntimeVersion != "" {
			runtime = function.RuntimeVersion
		} else {
			runtime = function.Runtime
		}

		meta := FunctionMetadata{
			InvocationMode: function.InvocationMode,
			Timeout:        function.Timeout,
		}
		file, err := newFunctionFile(functionsDir, relPath, fileInfo, runtime, &meta, tmpDir, observer)
		if err != nil {
			return nil, nil, nil, err
		}

		if function.Schedule != "" {
			schedules = append(schedules, &models.FunctionSchedule{
				Cron: function.Schedule,
				Name: function.Name,
			})
		}

		routes := make([]*models.FunctionRoute, len(function.Routes))
		for i, route := range function.Routes {
			routes[i] = &models.FunctionRoute{
				Pattern:      route.Pattern,
				Literal:      route.Literal,
				Expression:   route.Expression,
				Methods:      route.Methods,
				PreferStatic: route.PreferStatic,
			}
		}

		excludedRoutes := make([]*models.ExcludedFunctionRoute, len(function.ExcludedRoutes))
		for i, route := range function.ExcludedRoutes {
			excludedRoutes[i] = &models.ExcludedFunctionRoute{
				Pattern:    route.Pattern,
				Literal:    route.Literal,
				Expression: route.Expression,
			}
		}

		hasConfig := function.DisplayName != "" || function.Generator != "" || len(routes) > 0 || len(excludedRoutes) > 0 || len(function.BuildData) > 0 || function.Priority != 0 || function.TrafficRules != nil || function.Timeout != 0 || len(function.EventSubscriptions) > 0 || function.Region != "" || function.Memory != 0 || function.Vcpu != 0
		if hasConfig {
			cfg := models.FunctionConfig{
				DisplayName:        function.DisplayName,
				Generator:          function.Generator,
				Memory:             function.Memory,
				Region:             function.Region,
				Routes:             routes,
				ExcludedRoutes:     excludedRoutes,
				BuildData:          function.BuildData,
				Priority:           int64(function.Priority),
				EventSubscriptions: function.EventSubscriptions,
				Vcpu:               function.Vcpu,
			}

			if function.TrafficRules != nil {
				cfg.TrafficRules = &models.TrafficRulesConfig{
					Action: &models.TrafficRulesConfigAction{
						Type: function.TrafficRules.Action.Type,
						Config: &models.TrafficRulesConfigActionConfig{
							Aggregate: function.TrafficRules.Action.Config.Aggregate,
							RateLimitConfig: &models.TrafficRulesRateLimitConfig{
								Algorithm:   function.TrafficRules.Action.Config.RateLimitConfig.Algorithm,
								WindowSize:  int64(function.TrafficRules.Action.Config.RateLimitConfig.WindowSize),
								WindowLimit: int64(function.TrafficRules.Action.Config.RateLimitConfig.WindowLimit),
							},
							To: function.TrafficRules.Action.Config.To,
						},
					},
				}
			}

			functionsConfig[file.Name] = cfg
		}

		functions.Add(file.Name, file)
	}

	return functions, schedules, functionsConfig, nil
}

// manifestFunctionRel converts a manifest function path into a path relative to
// the functions directory, rejecting anything that points outside it.
func manifestFunctionRel(rootName, path string) (string, error) {
	rel := path
	if filepath.IsAbs(path) {
		absRoot, err := filepath.Abs(rootName)
		if err != nil {
			return "", err
		}
		rel, err = filepath.Rel(absRoot, path)
		if err != nil {
			return "", fmt.Errorf("manifest file specifies a function path outside the functions directory: %s", path)
		}
	}
	rel = filepath.Clean(rel)
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("manifest file specifies a function path outside the functions directory: %s", path)
	}
	return rel, nil
}

func readZipRuntime(root *os.Root, relPath string) (string, error) {
	f, err := openRegularFileInRoot(root, relPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	zf, err := zip.NewReader(f, info.Size())
	if err != nil {
		return "", err
	}

	for _, file := range zf.File {
		if file.Name == "netlify-toolchain" {
			fc, err := file.Open()
			if err != nil {
				// Ignore any errors and choose the default runtime.
				// This preserves the current behavior in this library.
				return jsRuntime, nil
			}
			defer func() { _ = fc.Close() }()

			var tc toolchainSpec
			if err := json.NewDecoder(fc).Decode(&tc); err != nil {
				// Ignore any errors and choose the default runtime.
				// This preserves the current behavior in this library.
				return jsRuntime, nil
			}
			return tc.Runtime, nil
		}
	}

	return jsRuntime, nil
}

func newFunctionFile(dir dirHandle, relPath string, i os.FileInfo, runtime string, metadata *FunctionMetadata, tmpDir *lazyTempDir, observer DeployObserver) (*FileBundle, error) {
	var file *FileBundle
	var err error

	if zipFile(i) || tarFile(i) {
		name := strings.TrimSuffix(i.Name(), filepath.Ext(i.Name()))
		file, err = createFunctionFileBundle(name, dir, relPath)
	} else {
		file, err = zipFunctionFile(dir, relPath, i, runtime, tmpDir)
	}
	if err != nil {
		return nil, err
	}

	file.Runtime = runtime
	file.FunctionMetadata = metadata

	if observer != nil {
		if err := observer.OnSuccessfulStep(file); err != nil {
			return nil, err
		}
	}

	return file, nil
}

func zipFunctionFile(dir dirHandle, relPath string, i os.FileInfo, runtime string, tmpDir *lazyTempDir) (*FileBundle, error) {
	src, err := openRegularFileInRoot(dir.root, relPath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = src.Close() }()

	tmpPath, tmpRoot, err := tmpDir.get()
	if err != nil {
		return nil, err
	}

	tmp, err := os.CreateTemp(tmpPath, "function-*.zip")
	if err != nil {
		return nil, err
	}
	defer func() {
		if tmp != nil {
			_ = tmp.Close()
		}
	}()

	s := sha256.New()
	archive := zip.NewWriter(io.MultiWriter(tmp, s))

	fileHeader, err := createHeader(archive, i, runtime)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(fileHeader, src); err != nil {
		return nil, err
	}
	if err := archive.Close(); err != nil {
		return nil, err
	}

	tmpName := tmp.Name()
	if err := tmp.Close(); err != nil {
		return nil, err
	}
	tmp = nil

	return &FileBundle{
		Name: strings.TrimSuffix(i.Name(), filepath.Ext(i.Name())),
		Sum:  hex.EncodeToString(s.Sum(nil)),
		Path: tmpName,
		root: tmpRoot,
		rel:  filepath.Base(tmpName),
	}, nil
}

// bundleEdgeFunctions reads the edge-bundler manifest from edgeFunctionsDir and turns each bundle it
// lists into an uploadable FileBundle. The deploy declares these as its edge_functions map
// ({format => code_sha}); the server replies with the subset (required_edge_functions) not already
// stored, and only those are streamed up. A missing manifest means no edge functions to upload.
func bundleEdgeFunctions(ctx context.Context, edgeFunctionsDir dirHandle, observer DeployObserver) (*deployFiles, error) {
	if !edgeFunctionsDir.valid() {
		return nil, nil
	}

	manifestFile, err := openRegularFileInRoot(edgeFunctionsDir.root, "manifest.json")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	manifestBytes, err := io.ReadAll(manifestFile)
	manifestFile.Close()
	if err != nil {
		return nil, err
	}

	context.GetLogger(ctx).Debug("Found edge functions manifest file")

	var manifest edgeFunctionsManifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("malformed edge functions manifest file: %w", err)
	}

	if len(manifest.Bundles) == 0 {
		return nil, nil
	}

	files := newDeployFiles()
	for _, bundle := range manifest.Bundles {
		file, err := newEdgeFunctionFile(edgeFunctionsDir, bundle)
		if err != nil {
			return nil, err
		}
		files.Add(file.Name, file)

		if observer != nil {
			if err := observer.OnSuccessfulStep(file); err != nil {
				return nil, err
			}
		}
	}

	return files, nil
}

func newEdgeFunctionFile(edgeFunctionsDir dirHandle, bundle edgeFunctionsManifestBundle) (*FileBundle, error) {
	// code_sha is the dedup key in the deployer<->functions-origin contract, so we compute it from the
	// bundle's bytes rather than trusting the edge-bundler's asset filename (which currently also happens
	// to be the sha256, but that's a bundler implementation detail). createFileBundleWithHasher streams
	// the bytes through the hasher, so the bundle is never held in memory.
	file, err := createFileBundleWithHasher(bundle.Format, edgeFunctionsDir, bundle.Asset, sha256.New())
	if err != nil {
		return nil, fmt.Errorf("edge functions manifest specifies a bundle that cannot be read: %s: %w", bundle.Asset, err)
	}

	return file, nil
}

func zipFile(i os.FileInfo) bool {
	return filepath.Ext(i.Name()) == ".zip"
}

func tarFile(i os.FileInfo) bool {
	name := i.Name()
	ext := filepath.Ext(name)
	return ext == ".tar" || ext == ".tgz" || strings.HasSuffix(name, ".tar.gz")
}

func jsFile(i os.FileInfo) bool {
	return filepath.Ext(i.Name()) == ".js"
}

func goFile(filePath string, i os.FileInfo, observer DeployObserver) bool {
	warner, hasWarner := observer.(DeployWarner)

	if m := i.Mode(); m&0o111 == 0 && runtime.GOOS != "windows" { // check if it's an executable file. skip on windows, since it doesn't have that mode
		if hasWarner {
			warner.OnWalkWarning(filePath, "Go binary does not have executable permissions")
		}
		return false
	}

	if _, err := elf.Open(filePath); err != nil { // check if it's a linux executable
		if hasWarner {
			warner.OnWalkWarning(filePath, "Go binary is not a linux executable")
		}
		return false
	}

	v, err := version.ReadExe(filePath)
	if err != nil || !strings.HasPrefix(v.Release, "go1.") {
		if hasWarner {
			warner.OnWalkWarning(filePath, "Unable to detect Go version 1.x")
		}
	}

	return true
}

func ignoreFile(rel string, ignoreInstallDirs bool) bool {
	if strings.HasPrefix(rel, ".") || strings.Contains(rel, "/.") || strings.HasPrefix(rel, "__MACOS") {
		return !strings.HasPrefix(rel, ".well-known/")
	}
	if !ignoreInstallDirs {
		return false
	}
	for _, ignorePath := range installDirs {
		if strings.HasPrefix(rel, ignorePath) {
			return true
		}
	}
	return false
}

func createHeader(archive *zip.Writer, i os.FileInfo, runtime string) (io.Writer, error) {
	if runtime == goRuntime || runtime == amazonLinux2 {
		return archive.CreateHeader(&zip.FileHeader{
			CreatorVersion: 3 << 8,      // indicates Unix
			ExternalAttrs:  0o777 << 16, // -rwxrwxrwx file permissions

			// we need to make sure we don't have two ZIP files with the exact same contents - otherwise, our upload deduplication mechanism will do weird things.
			// adding in the function name as a comment ensures that every function ZIP is unique
			Comment: i.Name(),

			Name:   "bootstrap",
			Method: zip.Deflate,
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
