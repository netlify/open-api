// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewListHookTypesParams creates a new ListHookTypesParams object
// with the default values initialized.
func NewListHookTypesParams() *ListHookTypesParams {

	return &ListHookTypesParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewListHookTypesParamsWithTimeout creates a new ListHookTypesParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewListHookTypesParamsWithTimeout(timeout time.Duration) *ListHookTypesParams {

	return &ListHookTypesParams{

		timeout: timeout,
	}
}

// NewListHookTypesParamsWithContext creates a new ListHookTypesParams object
// with the default values initialized, and the ability to set a context for a request
func NewListHookTypesParamsWithContext(ctx context.Context) *ListHookTypesParams {

	return &ListHookTypesParams{

		Context: ctx,
	}
}

// NewListHookTypesParamsWithHTTPClient creates a new ListHookTypesParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewListHookTypesParamsWithHTTPClient(client *http.Client) *ListHookTypesParams {

	return &ListHookTypesParams{
		HTTPClient: client,
	}
}

/*
ListHookTypesParams contains all the parameters to send to the API endpoint
for the list hook types operation typically these are written to a http.Request
*/
type ListHookTypesParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the list hook types params
func (o *ListHookTypesParams) WithTimeout(timeout time.Duration) *ListHookTypesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list hook types params
func (o *ListHookTypesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list hook types params
func (o *ListHookTypesParams) WithContext(ctx context.Context) *ListHookTypesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list hook types params
func (o *ListHookTypesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list hook types params
func (o *ListHookTypesParams) WithHTTPClient(client *http.Client) *ListHookTypesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list hook types params
func (o *ListHookTypesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *ListHookTypesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
