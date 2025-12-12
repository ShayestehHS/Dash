package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"Dash/internal/database"
	"Dash/pkg/logger"
)

func findMigrationsPath() (string, error) {
	basePath := filepath.Join("db", "repository")
	
	standardPath := filepath.Join(basePath, "postgresql", "migrations")
	if _, err := os.Stat(standardPath); err == nil {
		return standardPath, nil
	}

	var foundPath string
	walkErr := filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == "migrations" {
			foundPath = path
			return filepath.SkipAll
		}
		return nil
	})

	if walkErr != nil {
		return "", walkErr
	}

	if foundPath != "" {
		return foundPath, nil
	}

	return "", fmt.Errorf("no migrations directory found in %s", basePath)
}

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up, down, force, version")
		steps   = flag.Int("steps", 0, "Number of migration steps (0 = all)")
		version = flag.Int("version", 0, "Version for force command")
	)
	flag.Parse()

	cfg := database.LoadConfig()

	migrationsPath, err := findMigrationsPath()
	if err != nil {
		logger.Error("Migrate:MIGRATIONS_DIR_NOT_FOUND", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		logger.Error("Migrate:FAILED_TO_GET_ABS_PATH", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	if err := database.Connect(cfg); err != nil {
		logger.Error("Migrate:DB_CONNECTION_FAILED", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}
	defer database.Close()

	driver, err := postgres.WithInstance(database.DB, &postgres.Config{})
	if err != nil {
		logger.Error("Migrate:FAILED_TO_CREATE_DRIVER", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	migrationURL := fmt.Sprintf("file://%s", absPath)
	m, err := migrate.NewWithDatabaseInstance(migrationURL, "postgres", driver)
	if err != nil {
		logger.Error("Migrate:FAILED_TO_CREATE_MIGRATOR", map[string]interface{}{
			"error": err.Error(),
		})
		os.Exit(1)
	}

	switch *command {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			logger.Error("Migrate:UP_FAILED", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
		if err == migrate.ErrNoChange {
			logger.Info("Migrate:NO_CHANGES", map[string]interface{}{
				"message": "No migrations to apply",
			})
		} else {
			logger.Info("Migrate:UP_SUCCESS", map[string]interface{}{
				"message": "Migrations applied successfully",
			})
		}

	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			logger.Error("Migrate:DOWN_FAILED", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
		if err == migrate.ErrNoChange {
			logger.Info("Migrate:NO_CHANGES", map[string]interface{}{
				"message": "No migrations to rollback",
			})
		} else {
			logger.Info("Migrate:DOWN_SUCCESS", map[string]interface{}{
				"message": "Migrations rolled back successfully",
			})
		}

	case "force":
		if *version == 0 {
			logger.Error("Migrate:VERSION_REQUIRED", map[string]interface{}{
				"error": "version flag is required for force command",
			})
			os.Exit(1)
		}
		err = m.Force(*version)
		if err != nil {
			logger.Error("Migrate:FORCE_FAILED", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
		logger.Info("Migrate:FORCE_SUCCESS", map[string]interface{}{
			"message": fmt.Sprintf("Migration forced to version %d", *version),
		})

	case "version":
		version, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			logger.Error("Migrate:VERSION_FAILED", map[string]interface{}{
				"error": err.Error(),
			})
			os.Exit(1)
		}
		if err == migrate.ErrNilVersion {
			logger.Info("Migrate:NO_VERSION", map[string]interface{}{
				"message": "No migrations have been applied",
			})
		} else {
			logger.Info("Migrate:VERSION_SUCCESS", map[string]interface{}{
				"version": version,
				"dirty":   dirty,
			})
		}

	default:
		logger.Error("Migrate:INVALID_COMMAND", map[string]interface{}{
			"command": *command,
			"error":   "invalid command. Use: up, down, force, version",
		})
		os.Exit(1)
	}
}

