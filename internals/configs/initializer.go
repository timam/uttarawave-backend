package configs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/cmd"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var serverInstance *cmd.Server

func InitializeConfig() error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("unable to get the current filename")
	}
	configDir := filepath.Dir(filename)

	files, err := os.ReadDir(configDir)
	if err != nil {
		logger.Error("Error reading config directory", zap.String("directory", configDir), zap.Error(err))
		return err
	}

	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
			if err := loadIndividualConfig(filepath.Join(configDir, file.Name())); err != nil {
				return err
			}
		}
	}

	if err := checkRequiredEnvs(); err != nil {
		return err
	}

	if viper.GetString("server.debug") == "true" {
		logger.Info("Debug mode enabled, logger will be reinitializing ")
		err := logger.InitializeLogger()
		if err != nil {
			return err
		}
		logger.Info("Logger reinitialized successfully")
	}

	// Set server and tracing service names dynamically
	configureServerAndTracingNames()

	return nil
}

// Add a global variable to store the previous server state
var previousServerConfig string

func loadIndividualConfig(path string) error {
	vip := viper.New()
	vip.SetConfigFile(path)

	if err := vip.ReadInConfig(); err != nil {
		logger.Error("Error reading config file", zap.String("file", path), zap.Error(err))
		return err
	}

	for _, key := range vip.AllKeys() {
		viper.Set(key, vip.Get(key)) // Merge into main viper instance
		logger.Info("Loaded config", zap.String("key", key), zap.String("value", vip.GetString(key)))
	}

	// Serialize the current server configuration
	currentServerConfig := viper.GetString("server.debug") + viper.GetString("server.port")

	vip.AutomaticEnv() // To read from environment variables

	vip.WatchConfig()
	vip.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Config file changed", zap.String("file", path))

		if err := vip.ReadInConfig(); err != nil {
			logger.Error("Error re-reading config file", zap.String("file", path), zap.Error(err))
		} else {
			// Re-serialize the new server configuration
			newServerConfig := vip.GetString("server.debug") + vip.GetString("server.port")

			if newServerConfig != currentServerConfig {
				for _, key := range vip.AllKeys() {
					viper.Set(key, vip.Get(key)) // Re-merge config changes into main viper instance
					logger.Info("Reloaded config", zap.String("key", key), zap.String("value", vip.GetString(key)))
				}
				logger.Info("Server configuration changed, reloading server")

				if serverInstance == nil {
					serverInstance, err = cmd.InitializeServer()
					if err != nil {
						logger.Error("Error initializing server", zap.String("file", path), zap.Error(err))
					}
					go serverInstance.RunServer()
				} else {
					if err := serverInstance.ReloadServer(); err != nil {
						logger.Error("Failed to reload server", zap.Error(err))
					}
				}
				currentServerConfig = newServerConfig // Update the stored config
			} else {
				logger.Info("No change in server-relevant configuration, no need to reload server")
			}
		}
	})

	return nil
}

func checkRequiredEnvs() error {
	requiredEnvs := []string{"ENV"}
	var missingEnvs []string

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			missingEnvs = append(missingEnvs, env)
		}
	}

	if len(missingEnvs) > 0 {
		err := fmt.Errorf("missing required environment variables: %s", strings.Join(missingEnvs, ", "))
		logger.Error("Configuration error", zap.Error(err))
		return err
	}

	// Validate ENV value
	env := os.Getenv("ENV")
	if env != "dev" && env != "prod" {
		err := fmt.Errorf("invalid ENV value: %s. It must be either 'dev' or 'prod'", env)
		logger.Error("Configuration error", zap.Error(err))
		return err
	}

	return nil
}

func configureServerAndTracingNames() {
	env := os.Getenv("ENV")

	baseName := viper.GetString("server.name")
	if baseName == "" {
		baseName = "uttarawave-backend"
	}

	var serverName string
	if env == "prod" {
		serverName = "prod-" + baseName
	} else {
		serverName = "dev-" + baseName
	}

	viper.Set("server.name", serverName)
}
