package dto

import "encoding/json"

type Process struct {
	From    string `json:"from"`
	ToTopic string `json:"to_topic"`
	Payload any    `json:"payload"`
}

func (p *Process) Decode(data []byte) error {
	return json.Unmarshal(data, &p)
}
