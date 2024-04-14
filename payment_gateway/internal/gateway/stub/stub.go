package stub

import (
	"context"
	"time"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway"
)

const (
	defaultTimeout = 3 * time.Second
	defaultRetries = 10
)

type gatewayStub struct {
	transactionInfo *gateway.TransactionInfo
}

func New(transactionInfo *gateway.TransactionInfo) gateway.PaymentGateway {
	return &gatewayStub{
		transactionInfo: transactionInfo,
	}
}

func (g *gatewayStub) CreatePayment(_ context.Context) (string, error) {
	return "test", nil
}

func (g *gatewayStub) CheckStatus(context.Context, string) (gateway.PaymentStatus, error) {
	return gateway.Succeeded, nil
}

func (g *gatewayStub) TransactionID() string {
	return g.transactionInfo.TransactionID
}

func (g *gatewayStub) Timeout() time.Duration {
	return defaultTimeout
}

func (g *gatewayStub) Retries() int {
	return defaultRetries
}
