package backup

import (
	"fmt"
	"os"
	"time"

	"github.com/antigravity/dbbackup/internal/config"
	"github.com/antigravity/dbbackup/internal/database"
	"github.com/antigravity/dbbackup/internal/logger"
	"github.com/antigravity/dbbackup/internal/notifier"
	"github.com/antigravity/dbbackup/internal/storage"
)

type Manager struct {
	DB       database.Database
	Storage  storage.Storage
	Config   config.BackupConfig
	Notifier notifier.Notifier
}

func NewManager(db database.Database, st storage.Storage, cfg config.BackupConfig, notif notifier.Notifier) *Manager {
	return &Manager{
		DB:       db,
		Storage:  st,
		Config:   cfg,
		Notifier: notif,
	}
}

func (m *Manager) PerformBackup() error {
	startTime := time.Now()
	logger.Info.Println("Starting backup...")

	// 1. Test DB Connection
	if err := m.DB.TestConnection(); err != nil {
		errMsg := fmt.Sprintf("Backup failed: database connection failed: %v", err)
		if m.Notifier != nil {
			m.Notifier.Notify(errMsg)
		}
		return fmt.Errorf("database connection failed: %v", err)
	}

	// 2. Perform DB Backup
	backupFile, err := m.DB.Backup(m.Config.Type)
	if err != nil {
		errMsg := fmt.Sprintf("Backup failed: database backup failed: %v", err)
		if m.Notifier != nil {
			m.Notifier.Notify(errMsg)
		}
		return fmt.Errorf("database backup failed: %v", err)
	}
	defer os.Remove(backupFile) // Clean up local file after upload

	logger.Info.Printf("Database backup created: %s", backupFile)

	// 3. Compress if enabled
	finalFile := backupFile
	if m.Config.Compression {
		compressedFile, err := CompressFile(backupFile)
		if err != nil {
			return fmt.Errorf("compression failed: %v", err)
		}
		// Remove original uncompressed file
		os.Remove(backupFile)
		finalFile = compressedFile
		defer os.Remove(finalFile) // Clean up compressed file too
		logger.Info.Printf("Backup compressed: %s", finalFile)
	}

	// 4. Upload to Storage
	// Use the filename as the destination path
	err = m.Storage.Upload(finalFile, finalFile)
	if err != nil {
		errMsg := fmt.Sprintf("Backup failed: upload to storage failed: %v", err)
		if m.Notifier != nil {
			m.Notifier.Notify(errMsg)
		}
		return fmt.Errorf("upload to storage failed: %v", err)
	}

	duration := time.Since(startTime)
	msg := fmt.Sprintf("Backup completed successfully in %s. File: %s", duration, finalFile)
	logger.Info.Println(msg)
	if m.Notifier != nil {
		m.Notifier.Notify(msg)
	}
	return nil
}
