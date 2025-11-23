package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/antigravity/dbbackup/internal/config"
)

type LocalStorage struct {
	Config config.StorageConfig
}

func NewLocalStorage(cfg config.StorageConfig) *LocalStorage {
	return &LocalStorage{Config: cfg}
}

func (l *LocalStorage) Upload(srcPath string, destPath string) error {
	// For local storage, upload is just a copy to the target directory
	targetPath := filepath.Join(l.Config.Path, destPath)
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func (l *LocalStorage) Download(srcPath string, destPath string) error {
	// For local storage, download is just a copy from the target directory
	sourcePath := filepath.Join(l.Config.Path, srcPath)

	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func (l *LocalStorage) List(path string) ([]string, error) {
	targetPath := filepath.Join(l.Config.Path, path)
	entries, err := os.ReadDir(targetPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func (l *LocalStorage) Delete(path string) error {
	targetPath := filepath.Join(l.Config.Path, path)
	return os.Remove(targetPath)
}

func (l *LocalStorage) GetReader(path string) (io.ReadCloser, error) {
	targetPath := filepath.Join(l.Config.Path, path)
	return os.Open(targetPath)
}
