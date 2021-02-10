package plumbing

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func Example() {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.Write([]byte(`{ "name": "Pets of Netlify" }`))
	}))
	defer server.Close()
	host := strings.ReplaceAll(server.URL, "http://", "")

	// Create the API client
	// For Netlify's production API use DefaultHost, DefaultBasePath, DefaultSchemes
	transport := httptransport.New(host, DefaultBasePath, []string{"http"})
	client := New(transport, strfmt.Default)

	// Prepare the API token
	authInfo := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("User-Agent", "Your app")
		r.SetHeaderParam("Authorization", "Bearer your_netlify_api_token")
		return nil
	})

	// Make a request
	params := operations.NewGetSiteParams()
	params.SiteID = "123"
	res, err := client.Operations.GetSite(params, authInfo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Payload.Name)
	// Output: Pets of Netlify
}
