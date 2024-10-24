package ui

import (
	"context"
	"fmt"
	"freenet/internal/configs"
	"freenet/internal/logger"
	"freenet/internal/services"
	"io"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// UI is a struct to hold the application's UI components.
type UI struct {
	App           *tview.Application
	LogView       *tview.TextView
	WarehouseView *tview.Table
	SettingsView  *tview.TextView
	FooterView    *tview.TextView   // FooterView to show the footer text
	SearchInput   *tview.InputField // InputField for search functionality
	layout        *tview.Flex       // The layout containing all UI components
	SearchVisible bool              // Track if the SearchInput is visible
}

// GlobalUI is a global instance of the UI struct.
var GlobalUI *UI

// SetupTUI initializes the UI components.
func InitUI(context context.Context) error {
	// Initialize the UI components and store them in the UI struct.
	GlobalUI = &UI{
		App:           tview.NewApplication(),
		LogView:       tview.NewTextView(),
		WarehouseView: tview.NewTable(),
		SettingsView:  tview.NewTextView(),
		FooterView:    tview.NewTextView(),   // Initialize FooterView
		SearchInput:   tview.NewInputField(), // Initialize SearchInput
		SearchVisible: false,                 // Initially, the search input is not visible
	}

	// Set up logView
	GlobalUI.LogView.
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			// Auto-scroll to the end whenever content changes
			GlobalUI.LogView.ScrollToEnd()
			GlobalUI.App.Draw()
		})

	// Set up warehouseView
	GlobalUI.WarehouseView.SetBorders(true)

	// Set up settingsView
	GlobalUI.SettingsView.
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetScrollable(true).
		SetText(fmt.Sprintf(
			"[yellow]Warehouse Path: [white]%s\n[yellow]Address: [white]%s\n[yellow]Port: [white]%d\n",
			configs.GlobalConfig.WarehouseConfig.Path,
			configs.GlobalConfig.NetworkConfig.Address,
			configs.GlobalConfig.NetworkConfig.Port,
		)).SetBorder(true).SetTitle("Settings")

	// Set up footerView
	GlobalUI.FooterView.
		SetDynamicColors(true).
		SetText("[yellow]Press [white]S[yellow] for search").
		SetTextAlign(tview.AlignCenter).
		SetBorder(false)

	// Set up searchInput (hidden at first)
	GlobalUI.SearchInput.
		SetLabel("Search: ").
		SetFieldWidth(0). // Set to 0 to allow full width
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				searchTerm := GlobalUI.SearchInput.GetText()

				if searchTerm != "" {
					logger.GlobalLogger.Debug(fmt.Sprintf("Searching for: %s", searchTerm))
					services.Client.Search(context, searchTerm)
				} else {
					logger.GlobalLogger.Debug("Closing Search Input")
				}

				// Restore the footer text after search
				GlobalUI.layout.RemoveItem(GlobalUI.SearchInput)
				GlobalUI.layout.AddItem(GlobalUI.FooterView, 1, 1, false)
				GlobalUI.FooterView.SetText("[yellow]Press [white]S[yellow] for search")
				GlobalUI.App.SetFocus(GlobalUI.FooterView)

				// Set search as no longer visible
				GlobalUI.SearchVisible = false
			}
		})

	// Layout the UI components
	GlobalUI.layout = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(
			tview.NewFlex().
				AddItem(GlobalUI.LogView, 0, 3, false).
				AddItem(
					tview.NewFlex().
						SetDirection(tview.FlexRow).
						AddItem(GlobalUI.WarehouseView, 0, 2, false).
						AddItem(GlobalUI.SettingsView, 0, 1, false),
					0, 1, true),
								0, 1, true).
		AddItem(GlobalUI.FooterView, 1, 1, false) // Add footer at the bottom

	// Capture 'S' key press for search, and only if the SearchInput is not already visible
	GlobalUI.App.SetRoot(GlobalUI.layout, true).SetFocus(GlobalUI.layout).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		logger.GlobalLogger.Debug("Pressed key: " + strconv.QuoteRune(event.Rune()))
		if (event.Rune() == 's' || event.Rune() == 'S') && !GlobalUI.SearchVisible {
			// Clear the SearchInput field before displaying it
			GlobalUI.SearchInput.SetText("")

			// Replace footer with search input (set to full width)
			GlobalUI.layout.RemoveItem(GlobalUI.FooterView)
			GlobalUI.layout.AddItem(GlobalUI.SearchInput, 1, 1, true) // Full width

			GlobalUI.App.SetFocus(GlobalUI.SearchInput)

			// Set search as visible and prevent 'S' from being entered
			GlobalUI.SearchVisible = true
			return nil // Return nil to discard the first 'S'
		}
		return event
	})

	// Initialize the WriterWrapper to capture stdout and stderr.
	writerWrapper := &WriterWrapper{
		logView: GlobalUI.LogView,
		app:     GlobalUI.App,
	}

	// Redirect log output to the writerWrapper.
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := logger.PipeReader.Read(buf) // Use the global PipeReader here
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Fprintln(os.Stderr, "Error reading from pipe:", err)
				return
			}
			if n > 0 {
				writerWrapper.Write(buf[:n])
			}
		}
	}()

	return nil
}

// Start runs the application with the layout. It encapsulates the SetRoot and Run logic.
func (ui *UI) Start() error {
	if err := ui.App.SetRoot(ui.layout, true).Run(); err != nil {
		return err
	}
	return nil
}
