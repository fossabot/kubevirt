// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Manifest manifest
// swagger:model Manifest
type Manifest struct {

	// content
	Content *OciDescriptor `json:"content,omitempty"`

	// created-at
	//
	// creation data
	// Pattern: \d{4}-\d{1,2}-\d{1,2}T\d{2}:\d{2}:\d{2}\.\d+
	CreatedAt string `json:"created_at,omitempty"`

	// media-type
	//
	// manifest-type
	MediaType string `json:"mediaType,omitempty"`

	// metadata
	//
	// KeyValue object to add complementary and format specific information
	Metadata interface{} `json:"metadata,omitempty"`

	// package-name
	//
	// package name
	Package string `json:"package,omitempty"`

	// release-name
	//
	// release name
	Release string `json:"release,omitempty"`
}

// Validate validates this manifest
func (m *Manifest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateContent(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateCreatedAt(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Manifest) validateContent(formats strfmt.Registry) error {

	if swag.IsZero(m.Content) { // not required
		return nil
	}

	if m.Content != nil {
		if err := m.Content.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("content")
			}
			return err
		}
	}

	return nil
}

func (m *Manifest) validateCreatedAt(formats strfmt.Registry) error {

	if swag.IsZero(m.CreatedAt) { // not required
		return nil
	}

	if err := validate.Pattern("created_at", "body", string(m.CreatedAt), `\d{4}-\d{1,2}-\d{1,2}T\d{2}:\d{2}:\d{2}\.\d+`); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Manifest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Manifest) UnmarshalBinary(b []byte) error {
	var res Manifest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
