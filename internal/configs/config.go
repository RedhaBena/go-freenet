package configs

// GlobalConfig is a globally accessible configuration instance, initialized with default values.
var GlobalConfig Config = Config{}

// Config is the structure that contains all the configuration settings.
type Config struct {
	LoggerConfig // Configuration settings for the logger.
	WarehouseConfig
	NetworkConfig
}
