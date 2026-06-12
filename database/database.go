package database

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/user/chatbot-whatsapp/config"
	"github.com/user/chatbot-whatsapp/models"
)

func Connect(cfg *config.DatabaseConfig, log *zap.Logger) (*gorm.DB, error) {
	gormLogLevel := logger.Warn
	if cfg.SSLMode == "disable" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Info("database connected",
		zap.String("host", cfg.Host),
		zap.String("db", cfg.DBName),
	)

	return db, nil
}

func Migrate(db *gorm.DB, log *zap.Logger) error {
	if err := db.AutoMigrate(
		&models.Customer{},
		&models.Message{},
		&models.Template{},
	); err != nil {
		return err
	}

	log.Info("database migration completed")
	return nil
}
