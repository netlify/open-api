package http

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/go-openapi/runtime"
)

var DefaultTransport = httpTransport()

type RetryableTransport struct {
	tr       runtime.ClientTransport
	attempts int
}

type retryableRoundTripper struct {
	tr       http.RoundTripper
	attempts int
}

func NewRetryableTransport(tr runtime.ClientTransport, attempts int) *RetryableTransport {
	return &RetryableTransport{
		tr:       tr,
		attempts: attempts,
	}
}

func (t *RetryableTransport) Submit(op *runtime.ClientOperation) (interface{}, error) {
	var client http.Client
	if op.Client == nil {
		client = *http.DefaultClient
	} else {
		client = *op.Client
	}

	transport := client.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	client.Transport = &retryableRoundTripper{
		tr:       transport,
		attempts: t.attempts,
	}

	op.Client = &client
	res, err := t.tr.Submit(op)

	return res, err
}

func (t *retryableRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	rr := autorest.NewRetriableRequest(req)

	for attempt := 0; attempt < t.attempts; attempt++ {
		err = rr.Prepare()
		if err != nil {
			return resp, err
		}

		resp, err = t.tr.RoundTrip(rr.Request())

		if err != nil || resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}

		if attempt+1 < t.attempts { // ignore delay check in the last request attempt
			if !delayWithRateLimit(resp, req.Cancel) {
				return resp, err
			}
		}
	}

	return resp, err
}

func delayWithRateLimit(resp *http.Response, cancel <-chan struct{}) bool {
	r := resp.Header.Get("X-RateLimit-Reset")
	if r == "" {
		return false
	}
	retryReset, err := strconv.ParseInt(r, 10, 0)
	if err != nil {
		return false
	}

	t := time.Unix(retryReset, 0)
	select {
	case <-time.After(t.Sub(time.Now())):
		return true
	case <-cancel:
		return false
	}
}

func httpTransport() *http.Transport {
	protoUpgrade := map[string]func(string, *tls.Conn) http.RoundTripper{
		"ignore-h2": func(string, *tls.Conn) http.RoundTripper { return nil },
	}

	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSNextProto:          protoUpgrade,
	}
}
