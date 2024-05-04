package dto

import "encoding/json"

type Process struct {
	From    string `json:"from"`
	ToTopic string `json:"to_topic"`
	Payload any    `json:"payload"`
}

// Encode serializes a ProcessedTransaction into a JSON-encoded byte slice.
func (t *Process) Encode() ([]byte, error) {
	data, err := json.Marshal(&t)
	if err != nil {
		return nil, err
	}

	return data, nil
}
