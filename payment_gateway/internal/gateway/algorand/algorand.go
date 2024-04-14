package algorand

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/ShmelJUJ/software-engineering/payment_gateway/internal/gateway"
	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/mnemonic"
	"github.com/algorand/go-algorand-sdk/v2/transaction"
)

// UserData represents user-specific data like wallet address and mnemonic.
type UserData struct {
	WalletAddress string
	Mnemonic      string
}

// AlgorandGateway provides methods for interacting with the Algorand blockchain.
type Gateway struct {
	client           *algod.Client
	cfg              *Config
	sender, receiver *UserData
	transactionInfo  *gateway.TransactionInfo
}

// New creates a new instance of AlgorandGateway.
func New(
	cfg *Config,
	transactionInfo *gateway.TransactionInfo,
	sender, receiver *UserData,
) (*Gateway, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, gateway.NewCreationGatewayError("failed to merge with default config", err)
	}

	algodClient, err := algod.MakeClient(cfg.AlgodAddress, cfg.AlgodToken)
	if err != nil {
		return nil, gateway.NewCreationGatewayError("failed to make algod client", err)
	}

	return &Gateway{
		client:          algodClient,
		cfg:             cfg,
		sender:          sender,
		receiver:        receiver,
		transactionInfo: transactionInfo,
	}, nil
}

// CreatePayment initiates a payment transaction on the Algorand blockchain.
func (g *Gateway) CreatePayment(ctx context.Context) (string, error) {
	sp, err := g.client.SuggestedParams().Do(ctx)
	if err != nil {
		return "", gateway.NewCreatePaymentError("failed to get suggested params", err)
	}

	amount, err := strconv.ParseInt(g.transactionInfo.Value, 10, 64)
	if err != nil {
		return "", gateway.NewCreationGatewayError("failed to parse transaction value", err)
	}

	ptxn, err := transaction.MakePaymentTxn(
		g.sender.WalletAddress,
		g.receiver.WalletAddress,
		uint64(amount),
		nil,
		"",
		sp,
	)
	if err != nil {
		return "", gateway.NewCreatePaymentError("failed to make payment txn", err)
	}

	privateKey, err := mnemonic.ToPrivateKey(g.sender.Mnemonic)
	if err != nil {
		return "", gateway.NewCreatePaymentError("failed convert mnemonic to private key", err)
	}

	_, sptxn, err := crypto.SignTransaction(privateKey, ptxn)
	if err != nil {
		return "", gateway.NewCreatePaymentError("failed to sign transaction", err)
	}

	pendingTxID, err := g.client.SendRawTransaction(sptxn).Do(ctx)
	if err != nil {
		return "", gateway.NewCreatePaymentError("failed to send raw transaction", err)
	}

	return pendingTxID, nil
}

// CheckStatus checks the status of a payment transaction on the Algorand blockchain.
func (g *Gateway) CheckStatus(ctx context.Context, paymentID string) (gateway.PaymentStatus, error) {
	if _, err := transaction.WaitForConfirmation(
		g.client,
		paymentID,
		g.cfg.ConfirmationWaitRounds,
		ctx,
	); err != nil {
		if strings.Contains(err.Error(), "timed out") {
			return gateway.Pending, nil
		} else if strings.Contains(err.Error(), "Transaction rejected") {
			return gateway.Cancelled, nil
		}

		return gateway.Undefined, gateway.NewCheckStatusError("failed to wait for confirmation", err)
	}

	return gateway.Succeeded, nil
}

// TransactionID returns the ID of the current transaction.
func (g *Gateway) TransactionID() string {
	return g.transactionInfo.TransactionID
}

// Timeout returns the configured timeout duration for the gateway operations.
func (g *Gateway) Timeout() time.Duration {
	return g.cfg.Timeout
}

// Retries returns the number of retries configured for gateway operations.
func (g *Gateway) Retries() int {
	return g.cfg.Retries
}
