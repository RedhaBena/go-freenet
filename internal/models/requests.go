package models

import (
	"freenet/internal/logger"
	"sync"
)

// Structure pour une requête
type Request struct {
	Key              string
	NodeID           string
	VisitedNeighbors []string
}

// RequestsStore : Dictionnaire pour stocker les requêtes déjà traitées
type RequestsStore struct {
	mu       sync.RWMutex
	Requests map[string]Request
}

// NewRequestsStore initialise le store de requêtes
func NewRequestsStore() *RequestsStore {
	return &RequestsStore{
		Requests: make(map[string]Request),
	}
}

// AddRequest ajoute une nouvelle requête au store
func (store *RequestsStore) AddRequest(requestID, key, nodeID string, visitedNeighbors []string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	// Crée une nouvelle requête
	request := Request{
		Key:              key,
		NodeID:           nodeID,
		VisitedNeighbors: visitedNeighbors,
	}

	// Ajoute la requête dans le dictionnaire
	store.Requests[requestID] = request
	logger.GlobalLogger.Debug("Requête ajoutée dans le RequestsStore: ID = " + requestID + ", NodeID = " + nodeID + ", Key = " + key)
}

// GetRequest récupère une requête du store
func (store *RequestsStore) GetRequest(requestID string) (Request, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	request, exists := store.Requests[requestID]
	return request, exists
}

// RemoveRequest supprime une requête du store
func (store *RequestsStore) RemoveRequest(requestID string) {
	store.mu.Lock()
	defer store.mu.Unlock()

	delete(store.Requests, requestID)
	logger.GlobalLogger.Debug("Requête supprimée : ID = " + requestID)
}

// UpdateRequest met à jour une requête existante dans le store
func (store *RequestsStore) UpdateRequest(requestID string, updatedRequest Request) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.Requests[requestID]; exists {
		store.Requests[requestID] = updatedRequest
		logger.GlobalLogger.Debug("Requête mise à jour dans le RequestsStore: ID = " + requestID + ", Key = " + updatedRequest.Key + ", NodeID = " + updatedRequest.NodeID)
	} else {
		logger.GlobalLogger.Error("Requête non trouvée dans le RequestsStore: ID = " + requestID)
	}
}
