package http

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
)

func TestRetryableTransport(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		attempts++

		if attempts == 1 {
			reset := fmt.Sprintf("%d", time.Now().Add(1*time.Second).Unix())
			rw.Header().Set("X-RateLimit-Reset", reset)
			rw.WriteHeader(http.StatusTooManyRequests)
			rw.Write([]byte("rate limited"))
		} else {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("ok"))
		}
	}))
	defer server.Close()

	rwrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
		return nil
	})

	hu, _ := url.Parse(server.URL)
	rt := NewRetryableTransport(httptransport.New(hu.Host, "/", []string{"http"}), 2)

	res, err := rt.Submit(&runtime.ClientOperation{
		ID:          "getSites",
		Method:      "GET",
		PathPattern: "/",
		Params:      rwrtr,
		Reader: runtime.ClientResponseReaderFunc(func(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
			if response.Code() == 200 {
				var result string
				if err := consumer.Consume(response.Body(), &result); err != nil {
					return nil, err
				}
				return result, nil
			}
			return nil, errors.New("Generic error")
		}),
	})

	require.NoError(t, err)
	require.Equal(t, 2, attempts)

	require.IsType(t, "", res)
	actual := res.(string)
	require.EqualValues(t, "ok", actual)
}

func TestRetryableTransportExceedsMaxAttempts(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		attempts++
		reset := fmt.Sprintf("%d", time.Now().Add(1*time.Second).Unix())
		rw.Header().Set("X-RateLimit-Reset", reset)
		rw.WriteHeader(http.StatusTooManyRequests)
		rw.Write([]byte("rate limited"))
	}))
	defer server.Close()

	rwrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
		return nil
	})

	hu, _ := url.Parse(server.URL)
	rt := NewRetryableTransport(httptransport.New(hu.Host, "/", []string{"http"}), 2)

	_, err := rt.Submit(&runtime.ClientOperation{
		ID:          "getSites",
		Method:      "GET",
		PathPattern: "/",
		Params:      rwrtr,
		Reader: runtime.ClientResponseReaderFunc(func(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
			if response.Code() == 200 {
				var result string
				if err := consumer.Consume(response.Body(), &result); err != nil {
					return nil, err
				}
				return result, nil
			}
			return nil, errors.New("Generic error")
		}),
	})

	require.Error(t, err)
	require.Equal(t, 2, attempts)
}

func TestRetryableWithDifferentError(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		attempts++

		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("not found"))
	}))
	defer server.Close()

	rwrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
		return nil
	})

	hu, _ := url.Parse(server.URL)
	rt := NewRetryableTransport(httptransport.New(hu.Host, "/", []string{"http"}), 2)

	_, err := rt.Submit(&runtime.ClientOperation{
		ID:          "getSites",
		Method:      "GET",
		PathPattern: "/",
		Params:      rwrtr,
		Reader: runtime.ClientResponseReaderFunc(func(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
			if response.Code() == 200 {
				var result string
				if err := consumer.Consume(response.Body(), &result); err != nil {
					return nil, err
				}
				return result, nil
			}
			return nil, errors.New("Generic error")
		}),
	})

	require.Error(t, err)
	require.Equal(t, 1, attempts)
}

func TestRetryableTransport_POST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte("test result"))
	}))
	defer server.Close()

	rwrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
		return req.SetBodyParam("test body")
	})

	hu, _ := url.Parse(server.URL)
	rt := NewRetryableTransport(httptransport.New(hu.Host, "/", []string{"http"}), 2)

	result, err := rt.Submit(&runtime.ClientOperation{
		ID:          "createSite",
		Method:      "POST",
		PathPattern: "/",
		Params:      rwrtr,
		Reader: runtime.ClientResponseReaderFunc(func(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
			if response.Code() == 201 {
				var result string
				if err := consumer.Consume(response.Body(), &result); err != nil {
					return nil, err
				}
				return result, nil
			}
			return nil, errors.New("Generic error")
		}),
	})

	require.NoError(t, err)
	actual := result.(string)
	require.EqualValues(t, "test result", actual)
}

func TestRetryableTransportWithRetry_POST(t *testing.T) {
	attempts := 0

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		attempts++

		if attempts == 1 {
			reset := fmt.Sprintf("%d", time.Now().Add(1*time.Second).Unix())
			body, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, "\"test body\"\n", string(body))
			rw.Header().Set("X-RateLimit-Reset", reset)
			rw.WriteHeader(http.StatusTooManyRequests)
			rw.Write([]byte("rate limited"))
		} else {
			body, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, "\"test body\"\n", string(body))
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("ok"))
		}
	}))
	defer server.Close()

	rwrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
		return req.SetBodyParam("test body")
	})

	hu, _ := url.Parse(server.URL)
	rt := NewRetryableTransport(httptransport.New(hu.Host, "/", []string{"http"}), 2)

	res, err := rt.Submit(&runtime.ClientOperation{
		ID:          "getSites",
		Method:      "POST",
		PathPattern: "/",
		Params:      rwrtr,
		Reader: runtime.ClientResponseReaderFunc(func(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
			if response.Code() == 200 {
				var result string
				if err := consumer.Consume(response.Body(), &result); err != nil {
					return nil, err
				}
				return result, nil
			}
			return nil, errors.New("Generic error")
		}),
	})

	require.NoError(t, err)
	require.Equal(t, 2, attempts)

	require.IsType(t, "", res)
	actual := res.(string)
	require.EqualValues(t, "ok", actual)
}

func TestRetryableTransportHTMLReply(t *testing.T) {
	responseBody := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Some random HTML response</title>
	</head>
	<body>

	</body>
	</html>`
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.Write([]byte(responseBody))
	}))
	defer server.Close()

	rwrtr := runtime.ClientRequestWriterFunc(func(req runtime.ClientRequest, _ strfmt.Registry) error {
		return nil
	})

	hu, _ := url.Parse(server.URL)
	rt := NewRetryableTransport(httptransport.New(hu.Host, "/", []string{"http"}), 2)

	result, err := rt.Submit(&runtime.ClientOperation{
		ID:          "getSite",
		Method:      "GET",
		PathPattern: "/",
		Params:      rwrtr,
		Reader: runtime.ClientResponseReaderFunc(func(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
			if response.Code() == 500 {
				var result string
				if err := consumer.Consume(response.Body(), &result); err != nil {
					return nil, err
				}
				return result, nil
			}
			return nil, errors.New("Generic error")
		}),
	})

	require.NoError(t, err)
	actual := result.(string)
	require.EqualValues(t, responseBody, actual)
}
