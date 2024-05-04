package dto

import "encoding/json"

// SucceededTransaction represents a successful transaction.
type SucceededTransaction struct {
	TransactionID string `json:"transaction_id"`
}

// Decode decodes JSON data into a SucceededTransaction object.
func (t *SucceededTransaction) Decode(data []byte) error {
	return json.Unmarshal(data, &t)
}

// FailedTransaction represents a failed transaction.
type FailedTransaction struct {
	TransactionID string `json:"transaction_id"`
	Reason        string `json:"reason"`
}

// Decode decodes JSON data into a FailedTransaction object.
func (t *FailedTransaction) Decode(data []byte) error {
	return json.Unmarshal(data, &t)
}
