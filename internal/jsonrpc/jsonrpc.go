// Package jsonrpc provides JSON-RPC 2.0 message types for the Thalex WebSocket API.
package jsonrpc

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

// Version is the JSON-RPC protocol version.
const Version = "2.0"

// Request represents a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      uint64      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *uint64         `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

// IsError returns true if the response contains an error.
func (r *Response) IsError() bool {
	return r.Error != nil
}

// Notification represents a JSON-RPC 2.0 notification (no ID).
type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Error represents a JSON-RPC 2.0 error object.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("jsonrpc error %d: %s", e.Code, e.Message)
}

// Message represents a parsed JSON-RPC message that can be either a Response or Notification.
type Message struct {
	// Response is set if this message is a response (has an ID).
	Response *Response
	// Notification is set if this message is a notification (has a method, no ID).
	Notification *Notification
}

// ParseMessage parses a raw JSON message into either a Response or Notification.
func ParseMessage(data []byte) (*Message, error) {
	// First, try to determine the type by looking for key fields.
	var raw struct {
		JSONRPC string           `json:"jsonrpc"`
		ID      *json.RawMessage `json:"id"`
		Method  string           `json:"method"`
		Result  json.RawMessage  `json:"result"`
		Error   json.RawMessage  `json:"error"`
		Params  json.RawMessage  `json:"params"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse JSON-RPC message: %w", err)
	}

	// If it has a method and no ID, it's a notification.
	if raw.Method != "" && raw.ID == nil {
		var notif Notification
		if err := json.Unmarshal(data, &notif); err != nil {
			return nil, fmt.Errorf("failed to parse notification: %w", err)
		}
		return &Message{Notification: &notif}, nil
	}

	// Otherwise it's a response.
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &Message{Response: &resp}, nil
}

// IDGenerator generates unique request IDs.
type IDGenerator struct {
	counter atomic.Uint64
}

// Next returns the next unique request ID.
func (g *IDGenerator) Next() uint64 {
	return g.counter.Add(1)
}

// NewRequest creates a new JSON-RPC request with an auto-generated ID.
func NewRequest(id uint64, method string, params interface{}) *Request {
	return &Request{
		JSONRPC: Version,
		ID:      id,
		Method:  method,
		Params:  params,
	}
}
