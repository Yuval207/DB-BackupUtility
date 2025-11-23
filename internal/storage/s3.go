package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	internalConfig "github.com/antigravity/dbbackup/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	Config internalConfig.StorageConfig
	client *s3.Client
}

func NewS3Storage(cfg internalConfig.StorageConfig) (*S3Storage, error) {
	// Load AWS config
	// This will automatically pick up AWS_ACCESS_KEY_ID etc from env if not specified
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), 
		config.WithRegion(cfg.Region),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	return &S3Storage{Config: cfg, client: client}, nil
}

func (s *S3Storage) Upload(srcPath string, destPath string) error {
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.Config.Path), // Path is used as Bucket name for S3
		Key:    aws.String(destPath),
		Body:   file,
	})
	return err
}

func (s *S3Storage) Download(srcPath string, destPath string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	file, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Config.Path),
		Key:    aws.String(srcPath),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (s *S3Storage) List(path string) ([]string, error) {
	resp, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.Config.Path),
		Prefix: aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	var files []string
	for _, item := range resp.Contents {
		files = append(files, *item.Key)
	}
	return files, nil
}

func (s *S3Storage) Delete(path string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.Config.Path),
		Key:    aws.String(path),
	})
	return err
}

func (s *S3Storage) GetReader(path string) (io.ReadCloser, error) {
	resp, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.Config.Path),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}
