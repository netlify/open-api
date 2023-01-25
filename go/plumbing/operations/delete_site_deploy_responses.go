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

// DeleteSiteDeployReader is a Reader for the DeleteSiteDeploy structure.
type DeleteSiteDeployReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteSiteDeployReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 204:
		result := NewDeleteSiteDeployNoContent()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewDeleteSiteDeployDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteSiteDeployNoContent creates a DeleteSiteDeployNoContent with default headers values
func NewDeleteSiteDeployNoContent() *DeleteSiteDeployNoContent {
	return &DeleteSiteDeployNoContent{}
}

/*
DeleteSiteDeployNoContent handles this case with default header values.

OK
*/
type DeleteSiteDeployNoContent struct {
	Payload *models.Deploy
}

func (o *DeleteSiteDeployNoContent) Error() string {
	return fmt.Sprintf("[DELETE /sites/{site_id}/deploys/{deploy_id}][%d] deleteSiteDeployNoContent  %+v", 204, o.Payload)
}

func (o *DeleteSiteDeployNoContent) GetPayload() *models.Deploy {
	return o.Payload
}

func (o *DeleteSiteDeployNoContent) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Deploy)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDeleteSiteDeployDefault creates a DeleteSiteDeployDefault with default headers values
func NewDeleteSiteDeployDefault(code int) *DeleteSiteDeployDefault {
	return &DeleteSiteDeployDefault{
		_statusCode: code,
	}
}

/*
DeleteSiteDeployDefault handles this case with default header values.

error
*/
type DeleteSiteDeployDefault struct {
	_statusCode int

	Payload *models.Error
}

// Code gets the status code for the delete site deploy default response
func (o *DeleteSiteDeployDefault) Code() int {
	return o._statusCode
}

func (o *DeleteSiteDeployDefault) Error() string {
	return fmt.Sprintf("[DELETE /sites/{site_id}/deploys/{deploy_id}][%d] deleteSiteDeploy default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteSiteDeployDefault) GetPayload() *models.Error {
	return o.Payload
}

func (o *DeleteSiteDeployDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}