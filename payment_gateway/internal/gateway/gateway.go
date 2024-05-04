package gateway

import (
	"context"
	"time"
)

//go:generate mockgen -package mocks -destination mocks/gateway_mocks.go github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway PaymentGateway

// PaymentStatus represents the status of a payment.
type PaymentStatus int

const (
	Undefined PaymentStatus = iota
	Pending
	Succeeded
	Cancelled
)

// TransactionInfo holds information about a transaction.
type TransactionInfo struct {
	TransactionID string
	Value         string
	Currency      string
}

// PaymentGateway is an interface that defines operations for a payment gateway.
type PaymentGateway interface {
	CreatePayment(context.Context) (string, error)
	CheckStatus(context.Context, string) (PaymentStatus, error)
	TransactionID() string
	Timeout() time.Duration
	Retries() int
}
