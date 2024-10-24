package services

import (
	"context"
	"freenet/internal/logger"
	"freenet/internal/models"

	"github.com/google/uuid"
)

// Search creates a new request message, stores it in the request store, and returns the request ID.
func (client *ServiceClient) Search(ctx context.Context, key string) {
	// Check if the requested key (file) exists in the Warehouse
	fileLocation, found := client.warehouse.GetFileLocation(key)
	if found {
		logger.GlobalLogger.Info("File found in our warehouse: Key = " + key + ", NodeID = " + fileLocation)
		return
	}
	logger.GlobalLogger.Warn("File not found in our warehouse: Key = " + key)

	// Step 1: Generate a new UUID for the RequestID
	requestID := uuid.New().String()

	// Step 2: Store the request in the RequestsStore
	client.requestsStore.AddRequest(requestID, key, "local", []string{})

	// Step 3: Log the new request
	logger.GlobalLogger.Info("New search request created for file " + key + ": " + requestID)

	client.handleRequest(requestID)
}

// handleRequest takes a request ID, searches for a neighbor, and forwards the request.
func (client *ServiceClient) handleRequest(requestID string) {
	// Get the request from the RequestsStore
	request, exists := client.requestsStore.GetRequest(requestID)
	if !exists {
		logger.GlobalLogger.Error("handleRequest: Request ID " + requestID + " not found in the request store.")
		return
	}

	// Keep searching for neighbors until no more are found or the request is successfully sent
	for {
		// Step 1: Find the next neighbor to forward the request to, excluding already visited neighbors
		neighborID, err := client.warehouse.NearestNeighborByFileID(request.Key, request.VisitedNeighbors)
		if err != nil {
			// If no more neighbors are available, send a refusal to the parent node
			if request.NodeID == "local" {
				// If the request originated locally, just print the message
				logger.GlobalLogger.Error("Your request " + requestID + " has no more neighbors to contact, file " + request.Key + " not found.")
			} else {
				// Send a refusal (negative message) to the parent node
				refusalMessage := models.NegativeMessage{
					RequestID: requestID,
				}

				// Send the refusal message to the parent node (NodeID is the parent node)
				success, err := client.sendMessageToNeighbor(request.NodeID, "negative", refusalMessage)
				if success {
					logger.GlobalLogger.Info("Negative message sent to parent node " + request.NodeID + " for request " + requestID)
				} else {
					logger.GlobalLogger.Error("Failed to send refusal message to parent node " + request.NodeID + " for request " + requestID + ": " + err.Error())
				}
			}
			return
		}

		// Step 2: Create a new RequestMessage to send to the neighbor
		requestMessage := models.RequestMessage{
			RequestID: requestID,
			Key:       request.Key,
		}

		// Step 3: Attempt to forward the request to the neighbor
		success, err := client.sendMessageToNeighbor(neighborID, "request", requestMessage)
		// Mark the neighbor as visited
		request.VisitedNeighbors = append(request.VisitedNeighbors, neighborID)
		client.requestsStore.UpdateRequest(requestID, request)
		if success {
			logger.GlobalLogger.Info("Successfully sent request " + requestID + " to neighbor " + neighborID)
			break // Exit the loop as the request was successfully sent
		} else {
			logger.GlobalLogger.Error("Failed to send request " + requestID + " to neighbor " + neighborID + ": " + err.Error())
			// continue to the next neighbor
		}
	}
}
