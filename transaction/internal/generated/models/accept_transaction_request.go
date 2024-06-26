// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AcceptTransactionRequest accept transaction request
//
// swagger:model AcceptTransactionRequest
type AcceptTransactionRequest struct {

	// sender
	// Required: true
	Sender *AcceptTransactionUserRequest `json:"sender"`
}

// Validate validates this accept transaction request
func (m *AcceptTransactionRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateSender(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AcceptTransactionRequest) validateSender(formats strfmt.Registry) error {

	if err := validate.Required("sender", "body", m.Sender); err != nil {
		return err
	}

	if m.Sender != nil {
		if err := m.Sender.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sender")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("sender")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this accept transaction request based on the context it is used
func (m *AcceptTransactionRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateSender(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AcceptTransactionRequest) contextValidateSender(ctx context.Context, formats strfmt.Registry) error {

	if m.Sender != nil {

		if err := m.Sender.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("sender")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("sender")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AcceptTransactionRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AcceptTransactionRequest) UnmarshalBinary(b []byte) error {
	var res AcceptTransactionRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
