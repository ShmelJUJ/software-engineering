// Code generated by go-swagger; DO NOT EDIT.

package transaction

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// CancelTransactionHandlerFunc turns a function with the right signature into a cancel transaction handler
type CancelTransactionHandlerFunc func(CancelTransactionParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn CancelTransactionHandlerFunc) Handle(params CancelTransactionParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// CancelTransactionHandler interface for that can handle valid cancel transaction params
type CancelTransactionHandler interface {
	Handle(CancelTransactionParams, interface{}) middleware.Responder
}

// NewCancelTransaction creates a new http.Handler for the cancel transaction operation
func NewCancelTransaction(ctx *middleware.Context, handler CancelTransactionHandler) *CancelTransaction {
	return &CancelTransaction{Context: ctx, Handler: handler}
}

/*
	CancelTransaction swagger:route POST /transaction/{id}/cancel transaction cancelTransaction

The method is used to cancel a transaction.
*/
type CancelTransaction struct {
	Context *middleware.Context
	Handler CancelTransactionHandler
}

func (o *CancelTransaction) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewCancelTransactionParams()
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
