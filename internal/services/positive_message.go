package services

import (
	"freenet/internal/logger"
	"freenet/internal/models"
)

// handlePositiveMessage processes a PositiveMessage
func (client *ServiceClient) handlePositiveMessage(msg models.PositiveMessage, senderID string) {
	// Get the original request from the store using the request ID
	request, exists := client.requestsStore.GetRequest(msg.RequestID)
	if !exists {
		logger.GlobalLogger.Error("Request ID " + msg.RequestID + " not found in the request store.")
		return
	}

	// Forward the PositiveResponse to the node that originally requested the file
	originalRequesterNodeID := request.NodeID
	if originalRequesterNodeID != "local" {
		success, err := client.sendMessageToNeighbor(request.NodeID, "positive", msg)
		if success {
			logger.GlobalLogger.Info("File found for the request  " + msg.RequestID + " of " + request.NodeID + " with key " + request.Key + " by " + msg.NodeID + ", positive message sent to parent node")
		} else {
			logger.GlobalLogger.Error("Failed to send positive message to parent node " + senderID + " for request " + msg.RequestID + " with Node ID " + request.NodeID + " and key " + request.Key + ": " + err.Error())
		}

	} else {
		logger.GlobalLogger.Info("Your request " + msg.RequestID + " for the file with key " + request.Key + " was successfully fulfilled by node " + msg.NodeID)
	}

	// Store the new file location in the warehouse
	err := client.warehouse.StoreFile(request.Key, msg.NodeID)
	if err != nil {
		logger.GlobalLogger.Error("Failed to store file in warehouse: " + err.Error())
		return
	}

	logger.GlobalLogger.Info("File key " + request.Key + " stored in warehouse with node ID " + msg.NodeID)

	Client.warehouseUpdateHook(Client.warehouse)
}
