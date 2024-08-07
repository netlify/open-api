// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/netlify/open-api/v2/go/models"
)

// GetSiteDevServerHookReader is a Reader for the GetSiteDevServerHook structure.
type GetSiteDevServerHookReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetSiteDevServerHookReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetSiteDevServerHookOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetSiteDevServerHookDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetSiteDevServerHookOK creates a GetSiteDevServerHookOK with default headers values
func NewGetSiteDevServerHookOK() *GetSiteDevServerHookOK {
	return &GetSiteDevServerHookOK{}
}

/*
GetSiteDevServerHookOK handles this case with default header values.

OK
*/
type GetSiteDevServerHookOK struct {
	Payload *models.DevServerHook
}

func (o *GetSiteDevServerHookOK) Error() string {
	return fmt.Sprintf("[GET /sites/{site_id}/dev_server_hooks/{id}][%d] getSiteDevServerHookOK  %+v", 200, o.Payload)
}

func (o *GetSiteDevServerHookOK) GetPayload() *models.DevServerHook {
	return o.Payload
}

func (o *GetSiteDevServerHookOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.DevServerHook)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetSiteDevServerHookDefault creates a GetSiteDevServerHookDefault with default headers values
func NewGetSiteDevServerHookDefault(code int) *GetSiteDevServerHookDefault {
	return &GetSiteDevServerHookDefault{
		_statusCode: code,
	}
}

/*
GetSiteDevServerHookDefault handles this case with default header values.

error
*/
type GetSiteDevServerHookDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the get site dev server hook default response
func (o *GetSiteDevServerHookDefault) Code() int {
	return o._statusCode
}

func (o *GetSiteDevServerHookDefault) Error() string {
	return fmt.Sprintf("[GET /sites/{site_id}/dev_server_hooks/{id}][%d] getSiteDevServerHook default  %+v", o._statusCode, o.Payload)
}

func (o *GetSiteDevServerHookDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *GetSiteDevServerHookDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
