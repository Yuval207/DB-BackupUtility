package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/antigravity/dbbackup/internal/config"
)

type AzureStorage struct {
	Config config.StorageConfig
	client *azblob.Client
}

func NewAzureStorage(cfg config.StorageConfig) (*AzureStorage, error) {
	// Assuming credentials are provided via connection string or similar mechanism
	// For simplicity, let's assume CredentialsFile contains the connection string
	// OR we can use environment variables.
	// The azblob SDK supports connection strings.
	
	// If CredentialsFile is set, read it. If not, try env var AZURE_STORAGE_CONNECTION_STRING
	connStr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	if cfg.CredentialsFile != "" {
		content, err := os.ReadFile(cfg.CredentialsFile)
		if err == nil {
			connStr = string(content)
		}
	}

	if connStr == "" {
		return nil, fmt.Errorf("azure connection string not found")
	}

	client, err := azblob.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, err
	}

	return &AzureStorage{Config: cfg, client: client}, nil
}

func (a *AzureStorage) Upload(srcPath string, destPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = a.client.UploadFile(context.TODO(), a.Config.Path, destPath, file, nil)
	return err
}

func (a *AzureStorage) Download(srcPath string, destPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = a.client.DownloadFile(context.TODO(), a.Config.Path, srcPath, file, nil)
	return err
}

func (a *AzureStorage) List(path string) ([]string, error) {
	pager := a.client.NewListBlobsFlatPager(a.Config.Path, &azblob.ListBlobsFlatOptions{
		Prefix: &path,
	})

	var files []string
	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}
		for _, blob := range resp.Segment.BlobItems {
			files = append(files, *blob.Name)
		}
	}
	return files, nil
}

func (a *AzureStorage) Delete(path string) error {
	_, err := a.client.DeleteBlob(context.TODO(), a.Config.Path, path, nil)
	return err
}

func (a *AzureStorage) GetReader(path string) (io.ReadCloser, error) {
	// DownloadStream is the method for streaming
	resp, err := a.client.DownloadStream(context.TODO(), a.Config.Path, path, nil)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
