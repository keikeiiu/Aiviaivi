package ws

import (
	"encoding/json"
)

// ParseClientMessage parses a raw WebSocket message from a client.
// Expected format: {"type":"danmaku","content":"...","video_time":12.5,...}
func ParseClientMessage(data []byte) (Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return Message{}, err
	}
	return msg, nil
}

// MakeError creates an error message for a client.
func MakeError(err string) Message {
	return Message{Type: "error", Error: err}
}
