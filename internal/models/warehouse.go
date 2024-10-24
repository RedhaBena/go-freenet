package models

import (
	"fmt"
	"freenet/internal/logger"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v2"
)

// Structure pour représenter l'entrepôt dans YAML
type WarehouseData struct {
	Files map[string]string `yaml:"files"` // clé : ID du fichier, valeur : emplacement
}

// Warehouse structure pour stocker les fichiers et leurs emplacements
type Warehouse struct {
	mu      sync.RWMutex
	storage WarehouseData
	file    string // chemin vers le fichier warehouse.yaml
}

// NewWarehouse initialise le warehouse en chargeant les données depuis warehouse.yaml
func NewWarehouse(file string) (*Warehouse, error) {
	warehouse := &Warehouse{
		storage: WarehouseData{Files: make(map[string]string)},
		file:    file,
	}

	// Vérifie si le fichier existe déjà, sinon crée un nouveau fichier
	if _, err := os.Stat(file); err == nil {
		// Le fichier existe, on le lit
		err = warehouse.loadFromFile()
		if err != nil {
			return nil, err
		}
	} else {
		// Crée un fichier vide
		err := warehouse.saveToFile()
		if err != nil {
			return nil, err
		}
	}
	return warehouse, nil
}

// loadFromFile charge les données depuis warehouse.yaml
func (w *Warehouse) loadFromFile() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := ioutil.ReadFile(w.file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &w.storage)
	if err != nil {
		return err
	}
	logger.GlobalLogger.Debug("Entrepôt chargé depuis le fichier: " + w.file)
	return nil
}

// saveToFile enregistre les données dans warehouse.yaml
// NEED TO LOCK BEFORE
func (w *Warehouse) saveToFile() error {
	data, err := yaml.Marshal(&w.storage)
	if err != nil {
		return err
	}

	err = os.WriteFile(w.file, data, 0644)
	if err != nil {
		return err
	}
	logger.GlobalLogger.Debug("Entrepôt sauvegardé dans le fichier : " + w.file)
	return nil
}

// StoreFile stocke un fichier avec son ID et son emplacement, et met à jour warehouse.yaml
func (w *Warehouse) StoreFile(fileID, location string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.storage.Files[fileID] = location

	// Sauvegarder les modifications dans le fichier
	err := w.saveToFile()
	if err != nil {
		return err
	}
	logger.GlobalLogger.Debug("Fichier " + fileID + " stocké à l'emplacement : " + location)
	return nil
}

// GetFileLocation récupère l'emplacement d'un fichier
func (w *Warehouse) GetFileLocation(fileID string) (string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	location, exists := w.storage.Files[fileID]
	return location, exists
}

// RemoveFile supprime un fichier de l'entrepôt et met à jour warehouse.yaml
func (w *Warehouse) RemoveFile(fileID string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.storage.Files, fileID)

	// Sauvegarder les modifications dans le fichier
	err := w.saveToFile()
	if err != nil {
		return err
	}
	logger.GlobalLogger.Debug("Fichier " + fileID + " supprimé de l'entrepôt.")
	return nil
}

// ListFiles affiche tous les fichiers dans l'entrepôt
func (w *Warehouse) ListFiles() map[string]string {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.storage.Files
}

// asciiSum calculates the sum of ASCII values of each character in a string.
func asciiSum(s string) int {
	sum := 0
	for _, char := range s {
		sum += int(char)
	}
	return sum
}

// NearestNeighborByFileID returns the neighbor (node) corresponding to the file ID
// whose ASCII value sum is nearest to the given target key, excluding visited neighbors and local files.
func (w *Warehouse) NearestNeighborByFileID(targetKey string, visitedNeighbors []string) (string, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if len(w.storage.Files) == 0 {
		return "", fmt.Errorf("warehouse is empty")
	}

	// Calculate the ASCII sum of the target key
	targetKeySum := asciiSum(targetKey)

	// Prepare a set for fast lookup of visited neighbors (nodes)
	visitedSet := make(map[string]struct{}, len(visitedNeighbors))
	for _, node := range visitedNeighbors {
		visitedSet[node] = struct{}{}
	}

	var nearestNeighbor string
	var nearestDistance = math.MaxInt64

	// Iterate over all files in the warehouse, excluding visited nodes and local files
	for fileID, node := range w.storage.Files {
		// Skip if the node is "local" or has been visited
		if node == "local" {
			continue
		}
		if _, visited := visitedSet[node]; visited {
			continue // Skip visited nodes
		}

		// Calculate the ASCII sum of the file ID
		fileIDSum := asciiSum(fileID)

		// Calculate the absolute difference between the target key sum and the file ID sum
		distance := int(math.Abs(float64(targetKeySum - fileIDSum)))

		// If this distance is smaller, update the nearest node
		if distance < nearestDistance {
			nearestDistance = distance
			nearestNeighbor = node
			logger.GlobalLogger.Debug("Nearest neighbor updated to: " + nearestNeighbor + " with distance : " + strconv.Itoa(nearestDistance))
		}
	}

	if nearestNeighbor == "" {
		return "", fmt.Errorf("no more neighbors to contact")
	}

	return nearestNeighbor, nil
}
