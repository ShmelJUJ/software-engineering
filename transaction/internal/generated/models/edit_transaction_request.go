// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// EditTransactionRequest edit transaction request
//
// swagger:model EditTransactionRequest
type EditTransactionRequest struct {

	// money info
	MoneyInfo *EditMoneyInfo `json:"money_info,omitempty"`
}

// Validate validates this edit transaction request
func (m *EditTransactionRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMoneyInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EditTransactionRequest) validateMoneyInfo(formats strfmt.Registry) error {
	if swag.IsZero(m.MoneyInfo) { // not required
		return nil
	}

	if m.MoneyInfo != nil {
		if err := m.MoneyInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("money_info")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("money_info")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this edit transaction request based on the context it is used
func (m *EditTransactionRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateMoneyInfo(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *EditTransactionRequest) contextValidateMoneyInfo(ctx context.Context, formats strfmt.Registry) error {

	if m.MoneyInfo != nil {

		if swag.IsZero(m.MoneyInfo) { // not required
			return nil
		}

		if err := m.MoneyInfo.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("money_info")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("money_info")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *EditTransactionRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *EditTransactionRequest) UnmarshalBinary(b []byte) error {
	var res EditTransactionRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
