package configs

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/timam/uttaracloud-finance-backend/pkg/logger"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
)

func InitializeConfig() error {
	configDir := "./config/"

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

	return nil
}

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

	vip.AutomaticEnv() // To read from environment variables

	vip.WatchConfig()
	vip.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Config file changed", zap.String("file", path))

		if err := vip.ReadInConfig(); err != nil {
			logger.Error("Error re-reading config file", zap.String("file", path), zap.Error(err))
		} else {
			for _, key := range vip.AllKeys() {
				viper.Set(key, vip.Get(key)) // Re-merge config changes into main viper instance
				logger.Info("Reloaded config", zap.String("key", key), zap.String("value", vip.GetString(key)))
			}
		}
	})

	return nil
}
