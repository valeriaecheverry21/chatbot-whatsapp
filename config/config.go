package config

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	WhatsApp WhatsAppConfig
}

type ServerConfig struct {
	Port        string
	Environment string
	LogLevel    string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type WhatsAppConfig struct {
	AccessToken    string
	PhoneNumberID  string
	BusinessID     string
	VerifyToken    string
	APIVersion     string
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("PORT", "8080")
	viper.SetDefault("ENVIRONMENT", "development")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_DB", 0)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	viper.BindEnv("WHATSAPP_ACCESS_TOKEN")
	viper.BindEnv("WHATSAPP_PHONE_NUMBER_ID")
	viper.BindEnv("WHATSAPP_BUSINESS_ID")
	viper.BindEnv("WHATSAPP_VERIFY_TOKEN")
	viper.BindEnv("WHATSAPP_API_VERSION")
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_SSLMODE")
	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("REDIS_DB")
	viper.BindEnv("PORT")
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("LOG_LEVEL")

	cfg := &Config{
		Server: ServerConfig{
			Port:        viper.GetString("PORT"),
			Environment: viper.GetString("ENVIRONMENT"),
			LogLevel:    viper.GetString("LOG_LEVEL"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		WhatsApp: WhatsAppConfig{
			AccessToken:    viper.GetString("WHATSAPP_ACCESS_TOKEN"),
			PhoneNumberID:  viper.GetString("WHATSAPP_PHONE_NUMBER_ID"),
			BusinessID:     viper.GetString("WHATSAPP_BUSINESS_ID"),
			VerifyToken:    viper.GetString("WHATSAPP_VERIFY_TOKEN"),
			APIVersion:     viper.GetString("WHATSAPP_API_VERSION"),
		},
	}

	return cfg, nil
}

func NewLogger(cfg *Config) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.Set(cfg.Server.LogLevel); err != nil {
		level = zapcore.InfoLevel
	}

	zapCfg := zap.NewProductionConfig()
	zapCfg.Level = zap.NewAtomicLevelAt(level)
	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Server.Environment == "development" {
		zapCfg.Encoding = "console"
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return zapCfg.Build()
}
