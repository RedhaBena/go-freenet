package services

import (
	"freenet/internal/logger"
	"freenet/internal/models"
)

// handleRequestMessage processes a RequestMessage
func (client *ServiceClient) handleRequestMessage(msg models.RequestMessage, senderID string) {
	// Check if the request has already been processed
	_, exists := client.requestsStore.GetRequest(msg.RequestID)
	if exists {
		// If the request has already been processed, send a NegativeMessage
		refusalMessage := models.NegativeMessage{
			RequestID: msg.RequestID,
		}

		success, err := client.sendMessageToNeighbor(senderID, "negative", refusalMessage)
		if success {
			logger.GlobalLogger.Warn("Request " + msg.RequestID + " has already been processed, negative message sent to parent node " + senderID)
		} else {
			logger.GlobalLogger.Error("Failed to send negative message to parent node " + senderID + " for already processed request " + refusalMessage.RequestID + ": " + err.Error())
		}
		return
	}

	// Add the request to the RequestsStore
	client.requestsStore.AddRequest(msg.RequestID, msg.Key, senderID, []string{senderID}) // visited neighbors: [senderID]

	// Check if the requested key (file) exists in the Warehouse
	fileLocation, found := client.warehouse.GetFileLocation(msg.Key)
	if found {
		// If the file location is "local", use the client's listening address, otherwise use the fileLocation
		nodeID := client.listeningAddress
		if fileLocation != "local" {
			nodeID = fileLocation
		}

		// If the file is found locally, send a PositiveMessage
		positiveResponse := models.PositiveMessage{
			RequestID: msg.RequestID,
			NodeID:    nodeID, // Use the determined node ID (either local address or file location)
		}

		logger.GlobalLogger.Info("File found in our warehouse: Key = " + msg.Key + ", NodeID = " + fileLocation)

		// Send the refusal message to the parent node (NodeID is the parent node)
		success, err := client.sendMessageToNeighbor(senderID, "positive", positiveResponse)
		if success {
			logger.GlobalLogger.Info("Positive message sent to parent node " + senderID + " for request " + positiveResponse.RequestID + " with Node ID " + positiveResponse.NodeID)
		} else {
			logger.GlobalLogger.Error("Failed to send positive message to parent node " + senderID + " for request " + positiveResponse.RequestID + " with Node ID " + positiveResponse.NodeID + ": " + err.Error())
		}
		return
	}
	logger.GlobalLogger.Warn("File " + msg.Key + " searched by " + senderID + " not found in our warehouse")

	client.handleRequest(msg.RequestID)
}
