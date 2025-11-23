package main

import (
	"fmt"

	"github.com/antigravity/dbbackup/internal/config"
	"github.com/antigravity/dbbackup/internal/database"
	"github.com/antigravity/dbbackup/internal/storage"
)

func getComponents(cfg config.Config) (database.Database, storage.Storage, error) {
	var db database.Database
	var st storage.Storage
	var err error

	// Initialize Database
	switch cfg.Database.Type {
	case "mysql":
		db = database.NewMySQL(cfg.Database)
	case "postgres":
		db = database.NewPostgres(cfg.Database)
	case "mongodb":
		db = database.NewMongoDB(cfg.Database)
	default:
		return nil, nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}

	// Initialize Storage
	switch cfg.Storage.Type {
	case "local":
		st = storage.NewLocalStorage(cfg.Storage)
	case "s3":
		st, err = storage.NewS3Storage(cfg.Storage)
	case "gcs":
		st, err = storage.NewGCSStorage(cfg.Storage)
	case "azure":
		st, err = storage.NewAzureStorage(cfg.Storage)
	default:
		return nil, nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize storage: %v", err)
	}

	return db, st, nil
}
