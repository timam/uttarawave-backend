// pkg/db/initializer.go

package db

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/timam/uttarawave-backend/models"
	"github.com/timam/uttarawave-backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitializePostgreSQL() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Dhaka",
		viper.GetString("database.postgres.host"),
		viper.GetString("database.postgres.user"),
		viper.GetString("database.postgres.password"),
		viper.GetString("database.postgres.dbname"),
		viper.GetString("database.postgres.port"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto Migrate the schema
	err = DB.AutoMigrate(&models.Building{}, &models.Device{}, &models.Customer{})
	if err != nil {
		logger.Error("Failed to auto migrate schema", zap.Error(err))
		return fmt.Errorf("failed to auto migrate schema: %v", err)
	}

	logger.Info("Successfully connected to PostgreSQL Server")
	return nil
}
