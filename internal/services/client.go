package services

import (
	"context"
	"fmt"
	"freenet/internal/configs"
	"freenet/internal/models"
)

var Client ServiceClient

type ServiceClient struct {
	warehouse           *models.Warehouse
	requestsStore       *models.RequestsStore
	listeningAddress    string                  // Address of this node for listening to requests
	warehouseUpdateHook func(*models.Warehouse) // Callback for notifying UI of warehouse changes
}

// InitServiceClient initializes the ServiceClient with necessary configurations and connections.
func InitServiceClient(ctx context.Context, updateHook func(*models.Warehouse)) error {
	Client = *new(ServiceClient) // Initializes Client as a new instance of ServiceClient.

	// Set the warehouse update hook
	Client.warehouseUpdateHook = updateHook

	// Créer un entrepôt en chargeant les données depuis le fichier
	warehouse, err := models.NewWarehouse(configs.GlobalConfig.WarehouseConfig.Path)
	if err != nil {
		return fmt.Errorf("failed to create warehouse: %v", err)
	}
	Client.warehouse = warehouse

	Client.requestsStore = models.NewRequestsStore()

	Client.listeningAddress = fmt.Sprintf("%s:%d", configs.GlobalConfig.NetworkConfig.Address, configs.GlobalConfig.NetworkConfig.Port)

	// Trigger an update to UI to load initial warehouse data
	Client.warehouseUpdateHook(Client.warehouse)

	return nil
}

func (client *ServiceClient) Start(ctx context.Context) error {
	return Client.startListening(ctx)
}
