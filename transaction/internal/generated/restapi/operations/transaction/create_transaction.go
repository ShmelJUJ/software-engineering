// Code generated by go-swagger; DO NOT EDIT.

package transaction

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// CreateTransactionHandlerFunc turns a function with the right signature into a create transaction handler
type CreateTransactionHandlerFunc func(CreateTransactionParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn CreateTransactionHandlerFunc) Handle(params CreateTransactionParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// CreateTransactionHandler interface for that can handle valid create transaction params
type CreateTransactionHandler interface {
	Handle(CreateTransactionParams, interface{}) middleware.Responder
}

// NewCreateTransaction creates a new http.Handler for the create transaction operation
func NewCreateTransaction(ctx *middleware.Context, handler CreateTransactionHandler) *CreateTransaction {
	return &CreateTransaction{Context: ctx, Handler: handler}
}

/*
	CreateTransaction swagger:route POST /transaction/create transaction createTransaction

The method is used to create transactions.
*/
type CreateTransaction struct {
	Context *middleware.Context
	Handler CreateTransactionHandler
}

func (o *CreateTransaction) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCreateTransactionParams()
	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		*r = *aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc.(interface{}) // this is really a interface{}, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
