package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"Dash/pkg/env"
	"Dash/pkg/logger"
)

var DB *sql.DB

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func LoadConfig() *Config {
	host := env.GetString("DB_HOST")
	if host == "" {
		logger.Warn("LoadConfig:DEFAULT_DB_HOST", map[string]interface{}{
			"warning": "DB_HOST not set, using default 'localhost'",
		})
		host = "localhost"
	}

	port := env.GetString("DB_PORT")
	if port == "" {
		logger.Warn("LoadConfig:DEFAULT_DB_PORT", map[string]interface{}{
			"warning": "DB_PORT not set, using default '5432'",
		})
		port = "5432"
	}

	user := env.GetString("DB_USER")
	if user == "" {
		logger.Warn("LoadConfig:DEFAULT_DB_USER", map[string]interface{}{
			"warning": "DB_USER not set, using default 'dash'",
		})
		user = "dash"
	}

	password := env.GetString("DB_PASSWORD")
	if password == "" {
		logger.Warn("LoadConfig:DEFAULT_DB_PASSWORD", map[string]interface{}{
			"warning": "DB_PASSWORD not set, using default 'dash'",
		})
		password = "dash"
	}

	dbName := env.GetString("DB_NAME")
	if dbName == "" {
		logger.Warn("LoadConfig:DEFAULT_DB_NAME", map[string]interface{}{
			"warning": "DB_NAME not set, using default 'dash'",
		})
		dbName = "dash"
	}

	sslMode := env.GetString("DB_SSLMODE")
	if sslMode == "" {
		logger.Warn("LoadConfig:DEFAULT_DB_SSLMODE", map[string]interface{}{
			"warning": "DB_SSLMODE not set, using default 'disable'",
		})
		sslMode = "disable"
	}

	return &Config{
		Host:            host,
		Port:            port,
		User:            user,
		Password:        password,
		DBName:          dbName,
		SSLMode:         sslMode,
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

func Connect(config *Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	DB.SetMaxOpenConns(config.MaxOpenConns)
	DB.SetMaxIdleConns(config.MaxIdleConns)
	DB.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return DB.Ping()
}
