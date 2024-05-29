package porcelain

import (
	"bytes"
	gocontext "context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
	"github.com/netlify/open-api/v2/go/porcelain/context"

	apiClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLFSSha(t *testing.T) {
	t.Run("test with not a pointer file", func(t *testing.T) {
		file := strings.NewReader("Not a pointer file")
		data, err := readLFSData(file)
		if err != nil {
			t.Fatal(err)
		}

		if data != nil {
			t.Fatal("expected data to be nil without proper formatting")
		}
	})

	t.Run("test with v1 pointer", func(t *testing.T) {
		content := `version https://git-lfs.github.com/spec/v1
oid sha256:7e56e498ccb4cbb9c672e1aed6710fb91b2fd314394a666c11c33b2059ea3d71
size 1743570
`
		file := strings.NewReader(content)
		data, err := readLFSData(file)
		if err != nil {
			t.Fatal(err)
		}

		if data.SHA != "7e56e498ccb4cbb9c672e1aed6710fb91b2fd314394a666c11c33b2059ea3d71" {
			t.Fatalf("expected `7e56e498ccb4cbb9c672e1aed6710fb91b2fd314394a666c11c33b2059ea3d71`, got `%v`", data.SHA)
		}

		if data.Size != 1743570 {
			t.Fatalf("expected `1743570`, got `%v`", data.Size)
		}
	})
}

func TestAddWithLargeMedia(t *testing.T) {
	files := newDeployFiles()
	tests := []struct {
		rel string
		sum string
	}{
		{"foo.jpg", "sum1"},
		{"bar.jpg", "sum2"},
		{"baz.jpg", "sum3:originalsha"},
	}

	for _, test := range tests {
		file := &FileBundle{}
		file.Sum = test.sum
		files.Add(test.rel, file)
	}

	out := files.Hashed["sum3"]
	if len(out) != 1 {
		t.Fatalf("expected `%d`, got `%d`", 1, len(out))
	}
	out2 := files.Sums["baz.jpg"]
	if out2 != "sum3:originalsha" {
		t.Fatalf("expected `%v`, got `%v`", "sum3:originalsha", out2)
	}
}

func TestOpenAPIClientWithWeirdResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(408)
		rw.Write([]byte(`{ "foo": "bar", "message": "a message", "code": 408 }`))
	}))
	defer server.Close()

	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)

	body := ioutil.NopCloser(bytes.NewReader([]byte("hello world")))
	params := operations.NewUploadDeployFileParams().WithDeployID("1234").WithPath("foo/bar/biz").WithFileBody(body)
	_, operationError := client.Operations.UploadDeployFile(params, nil)
	require.Error(t, operationError)
	require.Equal(t, "[PUT /deploys/{deploy_id}/files/{path}][408] uploadDeployFile default  &{Code:408 Message:a message}", operationError.Error())
}

func TestConcurrentFileUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(408)
		rw.Write([]byte(`{ "foo": "bar", "message": "a message", "code": 408 }`))
	}))
	defer server.Close()

	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	for i := 0; i < 30; i++ {
		go func() {
			body := ioutil.NopCloser(bytes.NewReader([]byte("hello world")))
			params := operations.NewUploadDeployFileParams().WithDeployID("1234").WithPath("foo/bar/biz").WithFileBody(body)
			_, _ = client.Operations.UploadDeployFile(params, nil)
		}()
	}
}

func TestWaitUntilDeployLive_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Write([]byte(`{ "state": "chillin" }`))
	}))
	defer server.Close()

	hu, err := url.Parse(server.URL)
	require.NoError(t, err)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)

	ctx := context.WithAuthInfo(gocontext.Background(), apiClient.BearerToken("token"))
	ctx, _ = gocontext.WithTimeout(ctx, 50*time.Millisecond)
	_, err = client.WaitUntilDeployLive(ctx, &models.Deploy{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}

func TestWaitUntilDeployProcessed_Success(t *testing.T) {
	reqNum := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		reqNum++
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")

		// validate the polling actually works
		if reqNum > 1 {
			rw.Write([]byte(`{ "state": "processed" }`))
		} else {
			rw.Write([]byte(`{ "state": "processing" }`))
		}
	}))
	defer server.Close()

	hu, err := url.Parse(server.URL)
	require.NoError(t, err)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)

	ctx := context.WithAuthInfo(gocontext.Background(), apiClient.BearerToken("token"))
	ctx, _ = gocontext.WithTimeout(ctx, 30*time.Second)
	d, err := client.WaitUntilDeployProcessed(ctx, &models.Deploy{})
	require.NoError(t, err)
	assert.Equal(t, "processed", d.State)
}

func TestWaitUntilDeployProcessed_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Write([]byte(`{ "state": "processing" }`))
	}))
	defer server.Close()

	hu, err := url.Parse(server.URL)
	require.NoError(t, err)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)

	ctx := context.WithAuthInfo(gocontext.Background(), apiClient.BearerToken("token"))
	ctx, _ = gocontext.WithTimeout(ctx, 50*time.Millisecond)
	_, err = client.WaitUntilDeployProcessed(ctx, &models.Deploy{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}

func TestWalk_IgnoreNodeModulesInRoot(t *testing.T) {
	dir, err := ioutil.TempDir("", "deploy")
	require.Nil(t, err)
	defer os.RemoveAll(dir)

	err = os.Mkdir(filepath.Join(dir, "node_modules"), os.ModePerm)
	require.Nil(t, err)
	err = ioutil.WriteFile(filepath.Join(dir, "node_modules", "root-package"), []byte{}, 0644)
	require.Nil(t, err)

	err = os.MkdirAll(filepath.Join(dir, "more", "node_modules"), os.ModePerm)
	require.Nil(t, err)
	err = ioutil.WriteFile(filepath.Join(dir, "more", "node_modules", "inner-package"), []byte{}, 0644)
	require.Nil(t, err)

	files, err := walk(dir, mockObserver{}, false, false)
	require.Nil(t, err)
	assert.NotNil(t, files.Files["node_modules/root-package"])
	assert.NotNil(t, files.Files["more/node_modules/inner-package"])

	// When deploy directory == build directory, ignore node_modules in deploy directory root.
	files, err = walk(dir, mockObserver{}, false, true)
	require.Nil(t, err)
	assert.Nil(t, files.Files["node_modules/root-package"])
	assert.NotNil(t, files.Files["more/node_modules/inner-package"])
}

func TestWalk_EdgeFunctions(t *testing.T) {
	files := newDeployFiles()

	netlifyDir, err := ioutil.TempDir("", ".netlify")
	require.Nil(t, err)
	defer os.RemoveAll(netlifyDir)

	edgeFunctionsDir, err := ioutil.TempDir(netlifyDir, "edge-functions-dist")
	require.Nil(t, err)
	defer os.RemoveAll(edgeFunctionsDir)

	err = ioutil.WriteFile(filepath.Join(edgeFunctionsDir, "manifest.json"), []byte{}, 0644)
	require.Nil(t, err)
	err = ioutil.WriteFile(filepath.Join(edgeFunctionsDir, "123456789.js"), []byte{}, 0644)
	require.Nil(t, err)

	err = addInternalFilesToDeploy(edgeFunctionsDir, edgeFunctionsInternalPath, files, mockObserver{})
	require.Nil(t, err)

	assert.NotNil(t, files.Files[".netlify/internal/edge-functions/manifest.json"])
	assert.NotNil(t, files.Files[".netlify/internal/edge-functions/123456789.js"])
}

func TestWalk_PublishedFilesAndEdgeFunctions(t *testing.T) {
	files := setupPublishedAssets(t)

	netlifyDir, err := ioutil.TempDir("", ".netlify")
	require.Nil(t, err)
	defer os.RemoveAll(netlifyDir)

	edgeFunctionsDir, err := ioutil.TempDir(netlifyDir, "edge-functions-dist")
	require.Nil(t, err)
	defer os.RemoveAll(edgeFunctionsDir)

	err = ioutil.WriteFile(filepath.Join(edgeFunctionsDir, "manifest.json"), []byte{}, 0644)
	require.Nil(t, err)
	err = ioutil.WriteFile(filepath.Join(edgeFunctionsDir, "123456789.js"), []byte{}, 0644)
	require.Nil(t, err)

	err = addInternalFilesToDeploy(edgeFunctionsDir, edgeFunctionsInternalPath, files, mockObserver{})
	require.Nil(t, err)

	assert.NotNil(t, files.Files["assets/styles.css"])
	assert.NotNil(t, files.Files["index.html"])
	assert.NotNil(t, files.Files[".netlify/internal/edge-functions/manifest.json"])
	assert.NotNil(t, files.Files[".netlify/internal/edge-functions/123456789.js"])
}

func TestWalk_PublishedFilesAndEdgeRedirects(t *testing.T) {
	files := setupPublishedAssets(t)

	netlifyDir, err := ioutil.TempDir("", ".netlify")
	require.Nil(t, err)
	defer os.RemoveAll(netlifyDir)

	edgeRedirectsDir, err := ioutil.TempDir(netlifyDir, "edge-redirects-dist")
	require.Nil(t, err)
	defer os.RemoveAll(edgeRedirectsDir)

	err = ioutil.WriteFile(filepath.Join(edgeRedirectsDir, "redirects.json"), []byte{}, 0644)
	require.Nil(t, err)

	err = addInternalFilesToDeploy(edgeRedirectsDir, edgeRedirectsInternalPath, files, mockObserver{})
	require.Nil(t, err)

	assert.NotNil(t, files.Files["assets/styles.css"])
	assert.NotNil(t, files.Files["index.html"])
	assert.NotNil(t, files.Files[".netlify/deploy-config/redirects.json"])
}

func setupPublishedAssets(t *testing.T) *deployFiles {
	publishDir, err := ioutil.TempDir("", "publish")
	require.Nil(t, err)

	t.Cleanup(func() { os.RemoveAll(publishDir) })

	err = os.Mkdir(filepath.Join(publishDir, "assets"), os.ModePerm)
	require.Nil(t, err)
	err = ioutil.WriteFile(filepath.Join(publishDir, "assets", "styles.css"), []byte{}, 0644)
	require.Nil(t, err)
	err = ioutil.WriteFile(filepath.Join(publishDir, "index.html"), []byte{}, 0644)
	require.Nil(t, err)

	files, err := walk(publishDir, mockObserver{}, false, false)
	require.Nil(t, err)

	return files
}

func TestUploadFiles_Cancelation(t *testing.T) {
	ctx, cancel := gocontext.WithCancel(gocontext.Background())
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		cancel() // Cancel deploy as soon as first file upload is attempted.
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Write([]byte(`{ "state": "canceled" }`))
	}))
	defer server.Close()

	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	client.uploadLimit = 1
	ctx = context.WithAuthInfo(ctx, apiClient.BearerToken("token"))

	// Create some files to deploy
	dir, err := ioutil.TempDir("", "deploy")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "foo.html"), []byte("Hello"), 0644))
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "bar.html"), []byte("World"), 0644))

	files, err := walk(dir, nil, false, false)
	require.NoError(t, err)
	d := &models.Deploy{}
	for _, bundle := range files.Files {
		d.Required = append(d.Required, bundle.Sum)
	}
	err = client.uploadFiles(ctx, d, files, nil, fileUpload, time.Minute, false)
	require.ErrorIs(t, err, gocontext.Canceled)
}

func TestUploadFiles_Errors(t *testing.T) {
	ctx := gocontext.Background()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	client.uploadLimit = 1
	ctx = context.WithAuthInfo(ctx, apiClient.BearerToken("token"))

	// Create some files to deploy
	dir, err := ioutil.TempDir("", "deploy")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "foo.html"), []byte("Hello"), 0644))

	files, err := walk(dir, nil, false, false)
	require.NoError(t, err)
	d := &models.Deploy{}
	for _, bundle := range files.Files {
		d.Required = append(d.Required, bundle.Sum)
	}
	err = client.uploadFiles(ctx, d, files, nil, fileUpload, time.Minute, false)
	require.Equal(t, err.Error(), "[PUT /deploys/{deploy_id}/files/{path}][500] uploadDeployFile default  &{Code:0 Message:}")
}

func TestUploadFiles422Error_SkipsRetry(t *testing.T) {
	attempts := 0
	ctx := gocontext.Background()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		defer func() {
			attempts++
		}()

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(`{"message": "Unprocessable Entity", "code": 422 }`))
	}))
	defer server.Close()

	// File upload:
	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	client.uploadLimit = 1
	ctx = context.WithAuthInfo(ctx, apiClient.BearerToken("token"))

	// Create some files to deploy
	dir, err := ioutil.TempDir("", "deploy")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "foo.html"), []byte("Hello"), 0644))

	files, err := walk(dir, nil, false, false)
	require.NoError(t, err)
	d := &models.Deploy{}
	for _, bundle := range files.Files {
		d.Required = append(d.Required, bundle.Sum)
	}
	// Set SkipRetry to true
	err = client.uploadFiles(ctx, d, files, nil, fileUpload, time.Minute, true)
	require.ErrorContains(t, err, "Code:422 Message:Unprocessable Entity")
	require.Equal(t, attempts, 1)
}

func TestUploadFunctions422Error_SkipsRetry(t *testing.T) {
	attempts := 0
	ctx := gocontext.Background()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		defer func() {
			attempts++
		}()

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		rw.Write([]byte(`{"message": "Unprocessable Entity", "code": 422 }`))
	}))
	defer server.Close()

	// Function upload:
	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	client.uploadLimit = 1
	apiCtx := context.WithAuthInfo(ctx, apiClient.BearerToken("token"))

	dir, err := ioutil.TempDir("", "deploy")
	functionsPath := filepath.Join(dir, ".netlify", "functions")
	os.MkdirAll(functionsPath, os.ModePerm)
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, ioutil.WriteFile(filepath.Join(functionsPath, "foo.js"), []byte("module.exports = () => {}"), 0644))

	files, _, _, err := bundle(ctx, functionsPath, mockObserver{})
	require.NoError(t, err)
	d := &models.Deploy{}
	for _, bundle := range files.Files {
		d.RequiredFunctions = append(d.RequiredFunctions, bundle.Sum)
	}
	// Set SkipRetry to true
	err = client.uploadFiles(apiCtx, d, files, nil, functionUpload, time.Minute, true)
	require.ErrorContains(t, err, "Code:422 Message:Unprocessable Entity")
	require.Equal(t, attempts, 1)
}

func TestUploadFiles400Error_NoSkipRetry(t *testing.T) {
	attempts := 0
	ctx := gocontext.Background()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		defer func() {
			attempts++
		}()

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"message": "Bad Request", "code": 400 }`))
		return
	}))
	defer server.Close()

	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	client.uploadLimit = 1
	ctx = context.WithAuthInfo(ctx, apiClient.BearerToken("token"))

	// Create some files to deploy
	dir, err := ioutil.TempDir("", "deploy")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "foo.html"), []byte("Hello"), 0644))

	files, err := walk(dir, nil, false, false)
	require.NoError(t, err)
	d := &models.Deploy{}
	for _, bundle := range files.Files {
		d.Required = append(d.Required, bundle.Sum)
	}
	// Set SkipRetry to false
	err = client.uploadFiles(ctx, d, files, nil, fileUpload, time.Minute, false)
	require.ErrorContains(t, err, "Code:400 Message:Bad Request")
	require.Greater(t, attempts, 1)
}

func TestUploadFiles_SkipEqualFiles(t *testing.T) {
	ctx := gocontext.Background()

	serverRequests := 0

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		serverRequests++

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Write([]byte(`{}`))
	}))
	defer server.Close()

	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	client.uploadLimit = 1
	ctx = context.WithAuthInfo(ctx, apiClient.BearerToken("token"))

	// Create some files to deploy
	dir, err := ioutil.TempDir("", "deploy")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	fileBody := []byte("Hello")

	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "a.html"), fileBody, 0644))
	require.NoError(t, ioutil.WriteFile(filepath.Join(dir, "b.html"), fileBody, 0644))

	files, err := walk(dir, nil, false, false)
	require.NoError(t, err)

	// Create some fake function bundles to deploy
	functionsDir, err := ioutil.TempDir("", "deploy-functions")
	require.NoError(t, err)
	defer os.RemoveAll(functionsDir)

	// Get the JS function bundle
	cwd, _ := os.Getwd()
	basePath := path.Join(filepath.Dir(cwd), "internal", "data")
	jsFunctionPath := strings.Replace(filepath.Join(basePath, "hello-js-function-test.zip"), "\\", "/", -1)
	bundleBody, err := ioutil.ReadFile(jsFunctionPath)
	require.NoError(t, err)

	require.NoError(t, ioutil.WriteFile(filepath.Join(functionsDir, "a.zip"), bundleBody, 0644))
	require.NoError(t, ioutil.WriteFile(filepath.Join(functionsDir, "b.zip"), bundleBody, 0644))

	functions, _, _, err := bundle(ctx, functionsDir, mockObserver{})
	require.NoError(t, err)

	d := &models.Deploy{}
	// uploadFiles relies on the fact that the list of sums is an array of unique values, as both
	// the files and bundles have the same SHA we only need one of them for the Required array
	d.Required = []string{files.Sums["a.html"]}
	d.RequiredFunctions = []string{functions.Sums["a"]}

	err = client.uploadFiles(ctx, d, files, nil, fileUpload, time.Minute, false)
	require.NoError(t, err)
	assert.Equal(t, 1, serverRequests)

	err = client.uploadFiles(ctx, d, functions, nil, functionUpload, time.Minute, false)
	require.NoError(t, err)
	assert.Equal(t, 2, serverRequests)
}

func TestUploadFunctions_RetryCountHeader(t *testing.T) {
	attempts := 0
	ctx, cancel := gocontext.WithCancel(gocontext.Background())
	t.Cleanup(cancel)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		defer func() {
			attempts++
		}()

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")

		retryCount := req.Header.Get("X-Nf-Retry-Count")

		if attempts == 0 {
			require.Empty(t, retryCount)
		} else {
			require.Equal(t, fmt.Sprint(attempts), retryCount)
		}

		if attempts <= 2 {
			rw.WriteHeader(http.StatusInternalServerError)

			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{ "name": "foo" }`))
	}))
	defer server.Close()

	hu, _ := url.Parse(server.URL)
	tr := apiClient.NewWithClient(hu.Host, "/api/v1", []string{"http"}, http.DefaultClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	client.uploadLimit = 1
	apiCtx := context.WithAuthInfo(ctx, apiClient.BearerToken("token"))

	dir, err := ioutil.TempDir("", "deploy")
	functionsPath := filepath.Join(dir, ".netlify", "functions")
	os.MkdirAll(functionsPath, os.ModePerm)
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	require.NoError(t, ioutil.WriteFile(filepath.Join(functionsPath, "foo.js"), []byte("module.exports = () => {}"), 0644))

	files, _, _, err := bundle(ctx, functionsPath, mockObserver{})
	require.NoError(t, err)
	d := &models.Deploy{}
	for _, bundle := range files.Files {
		d.RequiredFunctions = append(d.RequiredFunctions, bundle.Sum)
	}

	require.NoError(t, client.uploadFiles(apiCtx, d, files, nil, functionUpload, time.Minute, false))
}

func TestBundle(t *testing.T) {
	functions, schedules, functionsConfig, err := bundle(gocontext.Background(), "../internal/data", mockObserver{})

	assert.Nil(t, err)
	assert.Equal(t, 5, len(functions.Files))
	assert.Empty(t, schedules)
	assert.Nil(t, functionsConfig)

	jsFunction := functions.Files["hello-js-function-test"]
	pyFunction := functions.Files["hello-py-function-test"]
	rsFunction := functions.Files["hello-rs-function-test"]
	goFunction := functions.Files["hello-go-binary-function"]
	goBackgroundFunction := functions.Files["hello-go-binary-function-background"]

	assert.Equal(t, "js", jsFunction.Runtime)
	assert.Equal(t, "py", pyFunction.Runtime)
	assert.Equal(t, "rs", rsFunction.Runtime)
	assert.Equal(t, "provided.al2", goFunction.Runtime)
	assert.Equal(t, "provided.al2", goBackgroundFunction.Runtime)

	assert.NotEqual(t, goFunction.Sum, goBackgroundFunction.Sum)
}

func TestBundleWithManifest(t *testing.T) {
	cwd, _ := os.Getwd()
	basePath := path.Join(filepath.Dir(cwd), "internal", "data")
	jsFunctionPath := strings.Replace(filepath.Join(basePath, "hello-js-function-test.zip"), "\\", "/", -1)
	pyFunctionPath := strings.Replace(filepath.Join(basePath, "hello-py-function-test.zip"), "\\", "/", -1)
	goFunctionPath := strings.Replace(filepath.Join(basePath, "hello-go-binary-function"), "\\", "/", -1)
	manifestPath := path.Join(basePath, "manifest.json")
	manifestFile := fmt.Sprintf(`{
		"functions": [
			{
				"path": "%s",
				"runtime": "a-runtime",
				"mainFile": "/some/path/hello-js-function-test.js",
				"displayName": "Hello Javascript Function",
				"generator": "@netlify/fake-plugin@1.0.0",
				"timeout": 60,
				"buildData": { "runtimeAPIVersion": 2 },
				"name": "hello-js-function-test",
				"schedule": "* * * * *",
				"routes": [
					{
						"pattern": "/products",
						"literal": "/products",
						"prefer_static": true
					},
					{
						"pattern": "/products/:id",
						"expression": "^/products/(.*)$",
						"methods": ["GET", "POST"]
					}
				]
			},
			{
				"path": "%s",
				"runtime": "some-other-runtime",
				"mainFile": "/some/path/hello-py-function-test",
				"name": "hello-py-function-test",
				"invocationMode": "stream"
			},	
			{
				"path": "%s",
				"runtime": "go",
				"runtimeVersion": "provided.al2",
				"name": "hello-go-binary-function"
			}
		],
		"version": 1
	}`, jsFunctionPath, pyFunctionPath, goFunctionPath)

	err := ioutil.WriteFile(manifestPath, []byte(manifestFile), 0644)
	defer os.Remove(manifestPath)
	assert.Nil(t, err)

	functions, schedules, functionsConfig, err := bundle(gocontext.Background(), "../internal/data", mockObserver{})
	assert.Nil(t, err)

	assert.Equal(t, 1, len(schedules))
	assert.Equal(t, "hello-js-function-test", schedules[0].Name)
	assert.Equal(t, "* * * * *", schedules[0].Cron)

	assert.Equal(t, 3, len(functions.Files))
	assert.Equal(t, "a-runtime", functions.Files["hello-js-function-test"].Runtime)
	assert.Empty(t, functions.Files["hello-js-function-test"].FunctionMetadata.InvocationMode)
	assert.Equal(t, int64(60), functions.Files["hello-js-function-test"].FunctionMetadata.Timeout)
	assert.Equal(t, "some-other-runtime", functions.Files["hello-py-function-test"].Runtime)
	assert.Equal(t, "stream", functions.Files["hello-py-function-test"].FunctionMetadata.InvocationMode)
	assert.Equal(t, "provided.al2", functions.Files["hello-go-binary-function"].Runtime)
	assert.Empty(t, functions.Files["hello-go-binary-function"].FunctionMetadata.InvocationMode)

	helloJSConfig := functionsConfig["hello-js-function-test"]

	assert.Equal(t, 1, len(functionsConfig))
	assert.Equal(t, "Hello Javascript Function", helloJSConfig.DisplayName)
	assert.Equal(t, "@netlify/fake-plugin@1.0.0", helloJSConfig.Generator)
	assert.EqualValues(t, 2, helloJSConfig.BuildData.(map[string]interface{})["runtimeAPIVersion"])

	assert.Equal(t, "/products", helloJSConfig.Routes[0].Pattern)
	assert.Equal(t, "/products", helloJSConfig.Routes[0].Literal)
	assert.Empty(t, helloJSConfig.Routes[0].Expression)
	assert.True(t, helloJSConfig.Routes[0].PreferStatic)

	assert.Equal(t, "/products/:id", helloJSConfig.Routes[1].Pattern)
	assert.Empty(t, helloJSConfig.Routes[1].Literal)
	assert.False(t, helloJSConfig.Routes[1].PreferStatic)
	assert.Equal(t, "^/products/(.*)$", helloJSConfig.Routes[1].Expression)
	assert.Equal(t, []string{"GET", "POST"}, helloJSConfig.Routes[1].Methods)
}

func TestReadZipRuntime(t *testing.T) {
	runtime, err := readZipRuntime("../internal/data/hello-rs-function-test.zip")

	assert.Nil(t, err)
	assert.Equal(t, "rs", runtime)
}

type mockObserver struct{}

func (m mockObserver) OnSetupWalk() error                         { return nil }
func (m mockObserver) OnSuccessfulStep(*FileBundle) error         { return nil }
func (m mockObserver) OnSuccessfulWalk(*models.DeployFiles) error { return nil }
func (m mockObserver) OnFailedWalk()                              {}

func (m mockObserver) OnSetupDelta(*models.DeployFiles) error                      { return nil }
func (m mockObserver) OnSuccessfulDelta(*models.DeployFiles, *models.Deploy) error { return nil }
func (m mockObserver) OnFailedDelta(*models.DeployFiles)                           {}

func (m mockObserver) OnSetupUpload(*FileBundle) error      { return nil }
func (m mockObserver) OnSuccessfulUpload(*FileBundle) error { return nil }
func (m mockObserver) OnFailedUpload(*FileBundle)           {}
