package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"freenet/internal/configs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// GlobalLogger is a globally accessible logger instance.
var GlobalLogger *zap.Logger

// PipeReader is the reader end of the pipe to capture logs for the UI.
var PipeReader *os.File

// PipeWriter is the writer end of the pipe.
var PipeWriter *os.File

// InitGlobalLogger initializes the global logger with a specified log level and configuration.
func InitGlobalLogger(ctx context.Context, config configs.LoggerConfig) error {
	// Create the pipe for redirecting logs to the UI.
	var err error
	PipeReader, PipeWriter, err = os.Pipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %v", err)
	}

	// Redirect stdout and stderr to the PipeWriter
	os.Stdout = PipeWriter
	os.Stderr = PipeWriter

	// Set default log level to "info"
	level := "info"
	if config.Debug {
		// If debug mode is enabled, set log level to "debug"
		level = "debug"
	}

	var lvl zapcore.Level
	// Set the zapcore level based on the specified log level
	if err := lvl.Set(level); err != nil {
		return fmt.Errorf("invalid log-level: %v", err)
	}

	// Create an encoder configuration for formatting log messages
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",                           // Key for the log entry time
		LevelKey:       "level",                          // Key for the log entry level
		NameKey:        "logger",                         // Key for the logger name
		MessageKey:     "msg",                            // Key for the log message
		LineEnding:     zapcore.DefaultLineEnding,        // Line ending character
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // Function to encode the level in capital letters with color
		EncodeTime:     humanReadableTimeEncoder,         // Function to encode the time in a human-readable format
		EncodeDuration: zapcore.StringDurationEncoder,    // Function to encode the duration as a string
		EncodeCaller:   zapcore.ShortCallerEncoder,       // Function to encode the caller information in a short format
	}

	// Create a zap configuration
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl), // Set the atomic level to the specified log level
		Development: false,                     // Set development mode to false
		Sampling: &zap.SamplingConfig{ // Sampling configuration for the logger
			Initial:    100, // Initial number of logs to sample
			Thereafter: 100, // Number of logs to sample thereafter
		},
		Encoding:         "console",          // Set the encoding to console format
		EncoderConfig:    encoderCfg,         // Use the specified encoder configuration
		OutputPaths:      []string{"stdout"}, // Output logs to stdout
		ErrorOutputPaths: []string{"stderr"}, // Output error logs to stderr
	}

	// Build the logger with the specified configuration
	logger, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}
	GlobalLogger = logger

	// Ensure that any buffered log entries are flushed before the program exits
	defer GlobalLogger.Sync()

	return nil
}

// humanReadableTimeEncoder formats the log entry time in a human-readable format.
func humanReadableTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
