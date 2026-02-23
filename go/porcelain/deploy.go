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
	"io"
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

	// Path OR Buffer should be populated
	Path   string
	Buffer io.ReadSeeker
}

type FunctionMetadata struct {
	InvocationMode string
	Timeout        int64
}

type toolchainSpec struct {
	Runtime string `json:"runtime"`
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
	ignoreInstallDirs := options.Dir == options.BuildDir

	context.GetLogger(ctx).Infof("Getting files info with large media flag: %v", largeMediaEnabled)
	files, err := walk(options.Dir, options.Observer, largeMediaEnabled, ignoreInstallDirs)
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

	if options.EdgeFunctionsDir != "" {
		err = addInternalFilesToDeploy(options.EdgeFunctionsDir, edgeFunctionsInternalPath, files, options.Observer)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedWalk()
			}
			return nil, err
		}
	}

	if options.EdgeRedirectsDir != "" {
		err = addInternalFilesToDeploy(options.EdgeRedirectsDir, edgeRedirectsInternalPath, files, options.Observer)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedWalk()
			}
			return nil, err
		}
	}

	if options.DbMigrationsDir != "" {
		err = addInternalFilesToDeploy(options.DbMigrationsDir, dbMigrationsInternalPath, files, options.Observer)
		if err != nil {
			if options.Observer != nil {
				options.Observer.OnFailedWalk()
			}
			return nil, err
		}
	}

	options.files = files

	functions, schedules, functionsConfig, err := bundle(ctx, options.FunctionsDir, options.Observer)
	if err != nil {
		if options.Observer != nil {
			options.Observer.OnFailedWalk()
		}
		return nil, err
	}
	options.functions = functions
	options.functionSchedules = schedules
	options.functionsConfig = functionsConfig

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

	if len(deploy.Required) == 0 && len(deploy.RequiredFunctions) == 0 {
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
	}

	count := 0
	for _, sha := range required {
		if files, exist := files.Hashed[sha]; exist {
			count += len(files)
		}
	}

	log := context.GetLogger(ctx)
	log.Infof("Uploading %v files", count)

	for _, sha := range required {
		if files, exist := files.Hashed[sha]; exist {
			file := files[0]

			select {
			case sem <- 1:
				wg.Add(1)
				go n.uploadFile(ctx, d, file, observer, t, timeout, wg, sem, sharedErr, skipRetry)
			case <-ctx.Done():
				log.Info("Context terminated, aborting file upload")
				return errors.Wrap(ctx.Err(), "aborted file upload early")
			}

			if len(files) > 1 {
				skippedFiles := files[1:]
				for _, file := range skippedFiles {
					log.Infof("Skipping file with content already uploaded: %s", file.Name)
				}
			}
		}
	}

	wg.Wait()

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

		var operationError error

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
			if operationError != nil {
				f.Buffer.Seek(0, 0)
			}
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

func createFileBundle(rel, path string) (*FileBundle, error) {
	o, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer o.Close()

	file := &FileBundle{
		Name: rel,
		Path: path,
	}

	s := sha1.New()
	if _, err := io.Copy(s, o); err != nil {
		return nil, err
	}

	file.Sum = hex.EncodeToString(s.Sum(nil))

	return file, nil
}

func walk(dir string, observer DeployObserver, useLargeMedia, ignoreInstallDirs bool) (*deployFiles, error) {
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

			if ignoreFile(rel, ignoreInstallDirs) {
				return nil
			}

			file, err := createFileBundle(rel, path)
			if err != nil {
				return err
			}

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

func addInternalFilesToDeploy(dir, internalPath string, files *deployFiles, observer DeployObserver) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Mode().IsRegular() {
			osRel, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			rel := internalPath + forceSlashSeparators(osRel)

			file, err := createFileBundle(rel, path)
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

func bundle(ctx context.Context, functionDir string, observer DeployObserver) (*deployFiles, []*models.FunctionSchedule, map[string]models.FunctionConfig, error) {
	if functionDir == "" {
		return nil, nil, nil, nil
	}

	manifestFile, err := os.Open(filepath.Join(functionDir, "manifest.json"))

	// If a `manifest.json` file is found, we extract the functions and their
	// metadata from it.
	if err == nil {
		defer manifestFile.Close()

		return bundleFromManifest(ctx, manifestFile, observer)
	}

	functions := newDeployFiles()

	info, err := ioutil.ReadDir(functionDir)
	if err != nil {
		return nil, nil, nil, err
	}

	for _, i := range info {
		filePath := filepath.Join(functionDir, i.Name())

		switch {
		case zipFile(i):
			runtime, err := readZipRuntime(filePath)
			if err != nil {
				return nil, nil, nil, err
			}
			file, err := newFunctionFile(filePath, i, runtime, nil, observer)
			if err != nil {
				return nil, nil, nil, err
			}
			functions.Add(file.Name, file)
		case jsFile(i):
			file, err := newFunctionFile(filePath, i, jsRuntime, nil, observer)
			if err != nil {
				return nil, nil, nil, err
			}
			functions.Add(file.Name, file)
		case goFile(filePath, i, observer):
			file, err := newFunctionFile(filePath, i, amazonLinux2, nil, observer)
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

func bundleFromManifest(ctx context.Context, manifestFile *os.File, observer DeployObserver) (*deployFiles, []*models.FunctionSchedule, map[string]models.FunctionConfig, error) {
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
		fileInfo, err := os.Stat(function.Path)
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
		file, err := newFunctionFile(function.Path, fileInfo, runtime, &meta, observer)
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

		hasConfig := function.DisplayName != "" || function.Generator != "" || len(routes) > 0 || len(excludedRoutes) > 0 || len(function.BuildData) > 0 || function.Priority != 0 || function.TrafficRules != nil || function.Timeout != 0
		if hasConfig {
			cfg := models.FunctionConfig{
				DisplayName:    function.DisplayName,
				Generator:      function.Generator,
				Routes:         routes,
				ExcludedRoutes: excludedRoutes,
				BuildData:      function.BuildData,
				Priority:       int64(function.Priority),
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

func readZipRuntime(filePath string) (string, error) {
	zf, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer zf.Close()

	for _, file := range zf.File {
		if file.Name == "netlify-toolchain" {
			fc, err := file.Open()
			if err != nil {
				// Ignore any errors and choose the default runtime.
				// This preserves the current behavior in this library.
				return jsRuntime, nil
			}
			defer fc.Close()

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

func newFunctionFile(filePath string, i os.FileInfo, runtime string, metadata *FunctionMetadata, observer DeployObserver) (*FileBundle, error) {
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

	if zipFile(i) || tarFile(i) {
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

	file.FunctionMetadata = metadata

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
