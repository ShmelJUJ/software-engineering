package dto

import "encoding/json"

// SucceededTransaction represents a successful transaction with a transaction ID.
type SucceededTransaction struct {
	TransactionID string `json:"transaction_id"`
}

// Encode converts the SucceededTransaction struct to JSON bytes.
func (t *SucceededTransaction) Encode() ([]byte, error) {
	data, err := json.Marshal(&t)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// FailedTransaction represents a failed transaction with a transaction ID and a reason for failure.
type FailedTransaction struct {
	TransactionID string `json:"transaction_id"`
	Reason        string `json:"reason"`
}

// Encode converts the FailedTransaction struct to JSON bytes.
func (t *FailedTransaction) Encode() ([]byte, error) {
	data, err := json.Marshal(&t)
	if err != nil {
		return nil, err
	}

	return data, nil
}
