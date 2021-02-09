package porcelain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-openapi/runtime"
	apiClient "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiContext "github.com/netlify/open-api/v2/go/porcelain/context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFormsBySiteId(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Write([]byte(`
			[
				{
					"id": "1",
					"site_id": "123",
					"name": "contact",
					"paths": [],
					"submission_count": 0,
					"fields": [],
					"created_at": ""
				}
			]`))
		assert.Equal(t, "/api/v1/sites/123/forms", req.URL.String())
	}))
	defer server.Close()

	httpClient := http.DefaultClient

	authInfo := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "buildbot")
		r.SetHeaderParam("Authorization", "Bearer 1234")
		return nil
	})

	parsedURL, err := url.Parse(server.URL)
	require.NoError(t, err)
	tr := apiClient.NewWithClient(parsedURL.Host, "/api/v1", []string{"http"}, httpClient)
	client := NewRetryable(tr, strfmt.Default, 1)
	forms, err := client.ListFormsBySiteId(apiContext.WithAuthInfo(context.Background(), authInfo), "123")
	require.NoError(t, err)
	assert.Equal(t, len(forms), 1)
}
