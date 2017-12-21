package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/go-openapi/runtime"
)

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

func (tr *RetryableTransport) Submit(op *runtime.ClientOperation) (interface{}, error) {
	op.Client.Transport = &retryableRoundTripper{
		tr:       op.Client.Transport,
		attempts: tr.attempts,
	}

	return tr.Submit(op)
}

func (tr *retryableRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	rr := autorest.NewRetriableRequest(req)

	// Increment to add the first call (attempts denotes number of retries)
	tr.attempts++
	for attempt := 0; attempt < tr.attempts; attempt++ {
		err = rr.Prepare()
		if err != nil {
			return resp, err
		}
		resp, err = tr.RoundTrip(rr.Request())
		if err != nil || resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}
		delayWithRateLimit(resp, req.Cancel)
	}
	return resp, err
}

func delayWithRateLimit(resp *http.Response, cancel <-chan struct{}) {
	retryReset, err := strconv.ParseInt(resp.Header.Get("X-RateLimit-Reset"), 10, 0)
	if err != nil {
		return
	}

	t := time.Unix(retryReset, 0)
	select {
	case <-time.After(t.Sub(time.Now())):
		return
	case <-cancel:
		return
	}
}
