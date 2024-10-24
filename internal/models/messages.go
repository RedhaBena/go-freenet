package models

import "encoding/json"

// Message is a wrapper structure to detect the type of incoming messages.
type Message struct {
	Type     string          `json:"type"`      // Type of the message: "request", "positive", or "negative"
	Data     json.RawMessage `json:"data"`      // Raw data for the actual message
	SenderID string          `json:"sender_id"` // ID of the node that sent the message
}

// RequestMessage represents a message to request a file from the network.
type RequestMessage struct {
	RequestID string `json:"request_id"` // RequestID is the unique identifier for this request.
	Key       string `json:"key"`        // Key is the unique identifier of the file being requested.
}

// PositiveMessage represents a message indicating that the requested file was found at a specific node.
type PositiveMessage struct {
	RequestID string `json:"request_id"` // RequestID is the unique identifier of the original request.
	NodeID    string `json:"node_id"`    // NodeID is the identifier of the node that contains the requested file.
}

// NegativeMessage represents a message indicating that the requested file was not found.
type NegativeMessage struct {
	RequestID string `json:"request_id"` // RequestID is the unique identifier of the original request.
}
