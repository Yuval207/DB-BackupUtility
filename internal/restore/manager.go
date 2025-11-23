package restore

import (
	"fmt"
	"os"
	"strings"

	"github.com/antigravity/dbbackup/internal/backup"
	"github.com/antigravity/dbbackup/internal/database"
	"github.com/antigravity/dbbackup/internal/logger"
	"github.com/antigravity/dbbackup/internal/storage"
)

type Manager struct {
	DB      database.Database
	Storage storage.Storage
}

func NewManager(db database.Database, st storage.Storage) *Manager {
	return &Manager{
		DB:      db,
		Storage: st,
	}
}

func (m *Manager) PerformRestore(backupFile string) error {
	logger.Info.Printf("Starting restore from %s...", backupFile)

	// 1. Download from Storage
	localFile := backupFile
	// If path contains directories, we might want to flatten it or ensure dirs exist.
	// For now, let's just download to current dir with same name.
	
	err := m.Storage.Download(backupFile, localFile)
	if err != nil {
		return fmt.Errorf("download from storage failed: %v", err)
	}
	defer os.Remove(localFile)

	// 2. Decompress if needed
	restoreFile := localFile
	if strings.HasSuffix(localFile, ".gz") {
		decompressedFile, err := backup.DecompressFile(localFile)
		if err != nil {
			return fmt.Errorf("decompression failed: %v", err)
		}
		restoreFile = decompressedFile
		defer os.Remove(restoreFile)
		logger.Info.Printf("Decompressed to: %s", restoreFile)
	}

	// 3. Restore to DB
	if err := m.DB.Restore(restoreFile); err != nil {
		return fmt.Errorf("database restore failed: %v", err)
	}

	logger.Info.Println("Restore completed successfully")
	return nil
}
