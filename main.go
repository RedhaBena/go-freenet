package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"freenet/internal/configs"
	"freenet/internal/logger"
	"freenet/internal/services"
	"freenet/internal/ui"

	"github.com/urfave/cli/v2"

	"go.uber.org/zap"
)

func main() {
	// Save the original stderr before it's piped by logger
	originalStderr := os.Stderr

	// Defer a function to recover from any panics and log the error.
	defer func() {
		if rerr := recover(); rerr != nil {
			if logger.GlobalLogger != nil {
				logger.GlobalLogger.Error("fatal panic | ", zap.Any("error", rerr))
			} else {
				fmt.Fprintf(originalStderr, "fatal panic: %v", rerr)
			}
		}
	}()

	// Create a context that can be cancelled.
	ctx, cancelFn := context.WithCancel(context.Background())
	// Create a channel to listen for OS signals.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to handle cancellation of the context when an OS signal is received.
	go func() {
		for range sigCh {
			logger.GlobalLogger.Warn("Cancel signal triggered")
			cancelFn()
		}
	}()

	// Define the CLI application.
	app := &cli.App{
		Name:  "freenet",
		Usage: "Freenet client for practical work MA_ParaDis",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Value:       "127.0.0.1",
				Usage:       "network address",
				Category:    "NETWORK",
				EnvVars:     []string{"ADDRESS"},
				Destination: &configs.GlobalConfig.NetworkConfig.Address,
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       43210,
				Usage:       "network port",
				Category:    "NETWORK",
				EnvVars:     []string{"PORT"},
				Destination: &configs.GlobalConfig.NetworkConfig.Port,
			},
			&cli.StringFlag{
				Name:        "warehouse",
				Value:       "warehouse.yaml",
				Usage:       "warehouse file path",
				Category:    "WAREHOUSE",
				EnvVars:     []string{"WAREHOUSE"},
				Destination: &configs.GlobalConfig.WarehouseConfig.Path,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Usage:       "debug logs",
				Category:    "LOGS",
				EnvVars:     []string{"DEBUG"},
				Destination: &configs.GlobalConfig.LoggerConfig.Debug,
			},
		},
		// Before function runs before any other actions.
		Before: func(cCtx *cli.Context) error {
			// Initialize the global logger.
			if err := logger.InitGlobalLogger(cCtx.Context, configs.GlobalConfig.LoggerConfig); err != nil {
				fmt.Fprintf(originalStderr, "Error: failed to create global logger: %v\n", err)

				return fmt.Errorf("failed to create global logger: %v", err)
			}

			// Initialize the ui.
			ui.InitUI(cCtx.Context)

			// Initialize the service client.
			if err := services.InitServiceClient(cCtx.Context, ui.GlobalUI.UpdateWarehouseView); err != nil {
				fmt.Fprintf(originalStderr, "Error: failed to create service client: %v\n", err)

				return fmt.Errorf("failed to create service client: %v", err)
			}

			return nil
		},
		Action: func(cCtx *cli.Context) error {
			// Initialize the service client.
			if err := services.Client.Start(cCtx.Context); err != nil {
				// Log the error to the console before exiting
				fmt.Fprintf(originalStderr, "Error: %v\n", err)

				return fmt.Errorf("failed to start listening: %v", err)
			}

			// Start the application with the layout.
			return ui.GlobalUI.Start()
		},
	}

	// Run the CLI application with the provided context and command-line arguments.
	if err := app.RunContext(ctx, os.Args); err != nil {
		if logger.GlobalLogger != nil {
			logger.GlobalLogger.Fatal("Fatal error | " + err.Error())
		} else {
			fmt.Fprintf(originalStderr, "Fatal error | %s\n", err.Error())
		}
	}
}
