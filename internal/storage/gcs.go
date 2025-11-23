package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/antigravity/dbbackup/internal/config"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type GCSStorage struct {
	Config config.StorageConfig
	client *storage.Client
}

func NewGCSStorage(cfg config.StorageConfig) (*GCSStorage, error) {
	var opts []option.ClientOption
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}

	client, err := storage.NewClient(context.TODO(), opts...)
	if err != nil {
		return nil, err
	}

	return &GCSStorage{Config: cfg, client: client}, nil
}

func (g *GCSStorage) Upload(srcPath string, destPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	wc := g.client.Bucket(g.Config.Path).Object(destPath).NewWriter(context.TODO())
	if _, err = io.Copy(wc, file); err != nil {
		wc.Close()
		return err
	}
	return wc.Close()
}

func (g *GCSStorage) Download(srcPath string, destPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	rc, err := g.client.Bucket(g.Config.Path).Object(srcPath).NewReader(context.TODO())
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(file, rc)
	return err
}

func (g *GCSStorage) List(path string) ([]string, error) {
	it := g.client.Bucket(g.Config.Path).Objects(context.TODO(), &storage.Query{Prefix: path})
	var files []string
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, attrs.Name)
	}
	return files, nil
}

func (g *GCSStorage) Delete(path string) error {
	return g.client.Bucket(g.Config.Path).Object(path).Delete(context.TODO())
}

func (g *GCSStorage) GetReader(path string) (io.ReadCloser, error) {
	return g.client.Bucket(g.Config.Path).Object(path).NewReader(context.TODO())
}
