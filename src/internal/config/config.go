package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	CanvusServer CanvusServerConfig `mapstructure:"canvus_server"`
	Paths        PathsConfig        `mapstructure:"paths"`
	Logging      LoggingConfig      `mapstructure:"logging"`
	Performance  PerformanceConfig  `mapstructure:"performance"`
}

// CanvusServerConfig contains Canvus Server connection settings
type CanvusServerConfig struct {
	URL      string `mapstructure:"url"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Timeout  int    `mapstructure:"timeout"` // seconds
}

// PathsConfig contains file system paths
type PathsConfig struct {
	AssetsFolder    string `mapstructure:"assets_folder"`
	BackupRootFolder string `mapstructure:"backup_root_folder"`
	OutputFolder    string `mapstructure:"output_folder"`
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level      string `mapstructure:"level"`       // debug, info, warn, error
	Verbose    bool   `mapstructure:"verbose"`
	LogToFile  bool   `mapstructure:"log_to_file"`
	LogFile    string `mapstructure:"log_file"`
}

// PerformanceConfig contains performance tuning settings
type PerformanceConfig struct {
	MaxConcurrentAPI    int `mapstructure:"max_concurrent_api"`
	MaxConcurrentFiles  int `mapstructure:"max_concurrent_files"`
	APIRequestTimeout   int `mapstructure:"api_request_timeout"`   // seconds
	FileOperationTimeout int `mapstructure:"file_operation_timeout"` // seconds
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		CanvusServer: CanvusServerConfig{
			URL:      "",
			Username: "",
			Password: "",
			Timeout:  30,
		},
		Paths: PathsConfig{
			AssetsFolder:     "",
			BackupRootFolder: "",
			OutputFolder:     "./output",
		},
		Logging: LoggingConfig{
			Level:      "info",
			Verbose:    false,
			LogToFile:  false,
			LogFile:    "kpmg-db-solver.log",
		},
		Performance: PerformanceConfig{
			MaxConcurrentAPI:     10,
			MaxConcurrentFiles:   20,
			APIRequestTimeout:    30,
			FileOperationTimeout: 60,
		},
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configFile string) (*Config, error) {
	config := DefaultConfig()

	// Set up viper
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.kpmg-db-solver")

	// Set environment variable prefix
	viper.SetEnvPrefix("KPMG")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Load config file if specified
	if configFile != "" {
		viper.SetConfigFile(configFile)
	}

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is OK, we'll use defaults and env vars
	} else {
		// Config file found, use it
		fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	// Unmarshal into config struct
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

// ValidateConfig validates the configuration
func (c *Config) Validate() error {
	// Validate Canvus Server settings
	if c.CanvusServer.URL == "" {
		return fmt.Errorf("canvus server URL is required")
	}
	if !strings.HasPrefix(c.CanvusServer.URL, "http") {
		return fmt.Errorf("canvus server URL must start with http:// or https://")
	}
	if c.CanvusServer.Username == "" {
		return fmt.Errorf("canvus server username is required")
	}
	if c.CanvusServer.Password == "" {
		return fmt.Errorf("canvus server password is required")
	}

	// Validate paths
	if c.Paths.AssetsFolder == "" {
		return fmt.Errorf("assets folder path is required")
	}
	if c.Paths.BackupRootFolder == "" {
		return fmt.Errorf("backup root folder path is required")
	}

	// Validate that paths exist
	if !pathExists(c.Paths.AssetsFolder) {
		return fmt.Errorf("assets folder does not exist: %s", c.Paths.AssetsFolder)
	}
	if !pathExists(c.Paths.BackupRootFolder) {
		return fmt.Errorf("backup root folder does not exist: %s", c.Paths.BackupRootFolder)
	}

	// Create output folder if it doesn't exist
	if err := os.MkdirAll(c.Paths.OutputFolder, 0755); err != nil {
		return fmt.Errorf("failed to create output folder: %w", err)
	}

	// Validate logging level
	validLevels := []string{"debug", "info", "warn", "error"}
	if !contains(validLevels, c.Logging.Level) {
		return fmt.Errorf("invalid logging level: %s (must be one of: %s)", 
			c.Logging.Level, strings.Join(validLevels, ", "))
	}

	// Validate performance settings
	if c.Performance.MaxConcurrentAPI < 1 {
		return fmt.Errorf("max concurrent API calls must be at least 1")
	}
	if c.Performance.MaxConcurrentFiles < 1 {
		return fmt.Errorf("max concurrent file operations must be at least 1")
	}

	return nil
}

// SaveConfig saves the configuration to a file
func (c *Config) SaveConfig(filename string) error {
	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set viper values
	viper.Set("canvus_server", c.CanvusServer)
	viper.Set("paths", c.Paths)
	viper.Set("logging", c.Logging)
	viper.Set("performance", c.Performance)

	// Write to file
	return viper.WriteConfigAs(filename)
}

// GetCanvusAPIURL returns the full API URL
func (c *Config) GetCanvusAPIURL() string {
	url := strings.TrimSuffix(c.CanvusServer.URL, "/")
	if !strings.HasSuffix(url, "/api/v1") {
		url += "/api/v1"
	}
	return url
}

// Helper functions

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
