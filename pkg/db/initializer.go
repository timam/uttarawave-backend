// TODO: if we enable debug mode, gorm log should be printed with hot reload

package db

import (
	"fmt"
	"github.com/spf13/viper"
	models "github.com/timam/uttarawave-backend/internals/models"
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
	gormConfig := &gorm.Config{
		QueryFields: true,
	}

	if viper.GetBool("server.debug") {
		gormConfig.Logger = logger.NewGormLogger()
	} else {
		gormConfig.Logger = logger.NewGormLogger().LogMode(logger.Silent)
	}

	DB, err = gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Auto Migrate the schema
	tables := []interface{}{
		&models.Address{},
		&models.Building{},
		&models.Customer{},
		&models.Package{},
		&models.Device{},
		&models.Subscription{},
		&models.Payment{},
		&models.Expense{},
	}

	for _, table := range tables {
		if err := DB.AutoMigrate(table); err != nil {
			logger.Error("Failed to auto migrate schema", zap.Error(err), zap.String("table", fmt.Sprintf("%T", table)))
			return fmt.Errorf("failed to auto migrate schema for %T: %v", table, err)
		}
		logger.Info("Successfully migrated table", zap.String("table", fmt.Sprintf("%T", table)))
	}

	logger.Info("Successfully connected to PostgreSQL Server")
	return nil
}
