package ui

import (
	"freenet/internal/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UpdateWarehouseView refreshes the warehouse table by fetching the latest warehouse data.
func (client *UI) UpdateWarehouseView(warehouse *models.Warehouse) {
	// Clear the current table content
	GlobalUI.WarehouseView.Clear()

	// Fetch the latest warehouse data
	warehouseData := warehouse.ListFiles()

	// Set table headers for better readability
	GlobalUI.WarehouseView.SetCell(0, 0, tview.NewTableCell("File ID").SetSelectable(true).SetTextColor(tcell.ColorYellow))
	GlobalUI.WarehouseView.SetCell(0, 1, tview.NewTableCell("Location").SetSelectable(false).SetTextColor(tcell.ColorYellow))

	// Populate the table with the new data
	row := 1
	for fileID, location := range warehouseData {
		GlobalUI.WarehouseView.SetCell(row, 0, tview.NewTableCell(fileID).SetSelectable(false))   // File ID
		GlobalUI.WarehouseView.SetCell(row, 1, tview.NewTableCell(location).SetSelectable(false)) // Location
		row++
	}
}
