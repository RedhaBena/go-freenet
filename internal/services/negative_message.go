package services

import (
	"freenet/internal/models"
)

// handleNegativeMessage processes a NegativeMessage
func (client *ServiceClient) handleNegativeMessage(msg models.NegativeMessage, senderID string) { // unused senderID, maybe remove
	client.handleRequest(msg.RequestID)
}
