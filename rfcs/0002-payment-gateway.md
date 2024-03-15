# Payment gateway service

## Problem
Solves the problem with the integration of third-party payment gateways, for example, with the [Yookassa Client](https://github.com/rvinnie/yookassa-sdk-go) or the [Algorand SDK](https://developer.algorand.org/docs/sdks/go/). The service provides a single contract for all payment gateways and performs a strictly defined sequence of actions

## Solution
1. The event service reads in Kafka from the topic `transaction.processing`, the possible structure of the event:
```json
{
	"transaction_id": "test-uuid",
	"sender_id": "test-uuid",
	"receiver_id": "test-uuid",
	"payment_method": "algorand",
	"payment_payload": "MIIBOgIBAAJBAKj34GkxFhD90vcNLYLInFEX6Ppy"
}
```

`payment_payload` - saves secret user payment data in encrypted form

from `transaction.canceled` topic, a possible event structure:
```json
{
	"transaction_id": "test-uuid"
}
```

Such events are written to the `cancel` channel associated with each worker

2. Next, for each transaction, we launch our own worker, which can be described by the structure:
```go
type Worker interface {
	Start()
	Stop() error
}

type PaymentGateway interface {
	CreateTransaction() error
	CheckTransactionStatus() (Status, error)
	CancelTransaction() error
	Timeout() time.Time
	Retries() int
}

type transactionWorker struct {
	transactionID uuid.UUID
	gateway PaymentGateway
	cancel <- chan struct{}
	stop <- chan struct{}
}
```

In separate files, we describe the implementation of the `PaymentGateway` interface, for example, `YookassaGateway` or `AlgorandGateway`

Each worker must have an event channel associated with it in order, for example, to stop processing a transaction if the client clicked cancel the operation

3. Next, we will describe in general terms the pipeline of working with `transactionWorker`
   
Creating a transaction using the `CreateTransaction` method -> expect the switch to branch either the successful status of the transaction, or the cancellation of the transaction, or an error with ending retries

4. In case of any problems with the transaction, the following event occurs:
```json
{
	"transaction_id": "test-uuid",
	"reason": "transaction retries have ended"
}
```

in the topic `transaction.failed`

If successful, the created event is:
```json
{
	"transaction_id": "test-uuid"
}
```

in the topic `transaction.finished`

## Algorand Payment Gateway
[Algorand](https://www.algorand.foundation/) - cryptocurrency protocol providing pure proof-of-stake on a blockchain. Algorand's native cryptocurrency is called ALGO

**Golang SDK**: https://developer.algorand.org/docs/sdks/go/

**Possibilities:**
1. Create account with `GenerateAccount()` method
2. Check balance with `algodClient.AccountInformation(acct.Address.String()).Do(context.Background())` method
3. Create transaction with `transaction.MakePaymentTxn(acct.Address.String(), acct.Address.String(), 100000, nil, "", sp)`, sign transaction with `crypto.SignTransaction(acct.PrivateKey, ptxn)` and send transaction with `algodClient.SendRawTransaction(sptxn).Do(context.Background())` methods
4. Check status with `transaction.WaitForConfirmation(algodClient, pendingTxID, 4, context.Background())` method

**Required data:**
1. Sender's wallet address
2. Sender's wallet private key
3. Receiver's wallet address
4. Amount
5. Token

## Yookassa Payment Gateway
[Yookassa](https://yookassa.ru/) - service for accepting payments on the Internet

**Not official Golang SDK:** https://github.com/rvinnie/yookassa-sdk-go

**Possibilities:**
1. Create transaction with `paymentHandler.CreatePayment`
2. Accept transaction with `paymentHandler.CapturePayment`
3. Cancel transaction with `paymentHandler.CancelPayment`
4. Check transaction status with `paymentHandler.FindPayment`

**Required data:**
1. Reciever's PAN
2. Sender's yookassa secret key
3. Payment method (ex. `bank_card`)
4. Amount
5. Currency
6. Confirmation type (ex. `redirect`)

## Advantages
1. Due to the interfaces, we avoid duplication of the logic of the main pipeline of the worker
2. The service does not store secret user payment data, it retrieves them once from the event message `transaction.processing` and this data is valid within the processing of a single transaction

## Disadvantages
1. In the current implementation, it will not be possible to achieve horizontal scaling, but this is a solvable problem
