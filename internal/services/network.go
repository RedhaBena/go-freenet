package services

import (
	"context"
	"encoding/json"
	"fmt"
	"freenet/internal/logger"
	"freenet/internal/models"
	"io"
	"net"
)

// startListening starts a TCP server to listen for incoming requests from other nodes.
func (client *ServiceClient) startListening(ctx context.Context) error {
	// Listen on the configured address and port
	listener, err := net.Listen("tcp", client.listeningAddress)
	if err != nil {
		return fmt.Errorf("failed to start listening on %s: %v", client.listeningAddress, err)
	}

	logger.GlobalLogger.Info("Listening for incoming requests on " + client.listeningAddress + "...")

	go func() {
		defer listener.Close()

		for {
			select {
			case <-ctx.Done():
				// Handle context cancellation (e.g., graceful shutdown)
				logger.GlobalLogger.Debug("Shutting down listener...")
				return
			default:
				// Accept incoming connections
				conn, err := listener.Accept()
				if err != nil {
					logger.GlobalLogger.Error("Error accepting connection: " + err.Error())
					continue
				}

				// Handle each connection in a separate goroutine
				logger.GlobalLogger.Debug("New connection from: " + conn.RemoteAddr().String())
				go client.handleIncomingConnection(conn)
			}
		}
	}()

	return nil
}

// handleIncomingConnection processes incoming messages from other nodes.
func (client *ServiceClient) handleIncomingConnection(conn net.Conn) {
	defer conn.Close()

	for {
		// Read and process the incoming message
		data := make([]byte, 4096) // Buffer to hold incoming data
		n, err := conn.Read(data)
		if err != nil {
			if err == io.EOF {
				// Connection was closed by the sender; this is expected.
				logger.GlobalLogger.Debug("Connection closed by " + conn.RemoteAddr().String())
				return
			}
			logger.GlobalLogger.Error("Failed to read data from connection " + conn.RemoteAddr().String() + ": " + err.Error())
			return
		}
		// Trim the data to the number of bytes actually read
		trimmedData := data[:n]

		// Unmarshal the wrapper message
		var msg models.Message
		err = json.Unmarshal(trimmedData, &msg)
		if err != nil {
			logger.GlobalLogger.Error("Failed to unmarshal message from " + conn.RemoteAddr().String() + ": " + err.Error())
			return
		}

		// Switch based on the message type
		switch msg.Type {
		case "request":
			var requestMsg models.RequestMessage
			err := json.Unmarshal(msg.Data, &requestMsg)
			if err == nil {
				logger.GlobalLogger.Info("Receive a request message for file " + requestMsg.Key + " from " + msg.SenderID + " with request ID " + requestMsg.RequestID)
				client.handleRequestMessage(requestMsg, msg.SenderID)
			} else {
				logger.GlobalLogger.Error("Failed to read request message from " + msg.SenderID + ": " + err.Error())
			}
		case "positive":
			var positiveMsg models.PositiveMessage
			err := json.Unmarshal(msg.Data, &positiveMsg)
			if err == nil {
				logger.GlobalLogger.Info("Receive a positive message for request " + positiveMsg.RequestID + " from " + msg.SenderID + " with node ID " + positiveMsg.NodeID)
				client.handlePositiveMessage(positiveMsg, msg.SenderID)
			} else {
				logger.GlobalLogger.Error("Failed to read positive message from " + msg.SenderID + ": " + err.Error())
			}
		case "negative":
			var negativeMsg models.NegativeMessage
			err := json.Unmarshal(msg.Data, &negativeMsg)
			if err == nil {
				logger.GlobalLogger.Warn("Receive a negative message for request " + negativeMsg.RequestID + " from " + msg.SenderID)
				client.handleNegativeMessage(negativeMsg, msg.SenderID)
			} else {
				logger.GlobalLogger.Error("Failed to read negative message from " + msg.SenderID + ": " + err.Error())
			}
		default:
			logger.GlobalLogger.Error("Unknown message type received from " + msg.SenderID + ": " + msg.Type)
		}

	}
}

// sendMessageToNeighbor sends any message to a neighbor and returns a boolean indicating success or failure.
// It wraps the message in a Message struct with the given message type and sender ID.
func (client *ServiceClient) sendMessageToNeighbor(neighborID string, messageType string, messagePayload interface{}) (bool, error) {
	// Marshal the actual message (e.g., RequestMessage, PositiveMessage, NegativeMessage)
	messageData, err := json.Marshal(messagePayload)
	if err != nil {
		return false, fmt.Errorf("Failed to marshal message payload: %v", err)
	}

	// Create the wrapper Message struct
	wrappedMessage := models.Message{
		Type:     messageType,
		Data:     messageData,
		SenderID: client.listeningAddress,
	}

	// Marshal the wrapped message into JSON
	wrappedMessageData, err := json.Marshal(wrappedMessage)
	if err != nil {
		return false, fmt.Errorf("Failed to marshal wrapped message: %v", err)
	}

	// Establish a TCP connection to the neighbor
	conn, err := net.Dial("tcp", neighborID)
	if err != nil {
		return false, fmt.Errorf("Failed to connect to neighbor %s: %v", neighborID, err)
	}
	defer conn.Close()

	// Send the wrapped message to the neighbor
	_, err = conn.Write(wrappedMessageData)
	if err != nil {
		return false, fmt.Errorf("Failed to send message to neighbor %s: %v", neighborID, err)

	}

	// Successfully sent the message
	logger.GlobalLogger.Debug("Message of type '" + messageType + "' successfully sent to neighbor " + neighborID)
	return true, nil
}
