package database

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/antigravity/dbbackup/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Config config.DatabaseConfig
	client *mongo.Client
}

func NewMongoDB(cfg config.DatabaseConfig) *MongoDB {
	return &MongoDB{Config: cfg}
}

func (m *MongoDB) Connect() error {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", m.Config.User, m.Config.Password, m.Config.Host, m.Config.Port)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}
	m.client = client
	return nil
}

func (m *MongoDB) TestConnection() error {
	if m.client == nil {
		if err := m.Connect(); err != nil {
			return err
		}
	}
	return m.client.Ping(context.TODO(), nil)
}

func (m *MongoDB) Backup(backupType string) (string, error) {
	// mongodump creates a directory by default, we should probably zip it or just use --archive
	filename := fmt.Sprintf("backup_mongo_%s_%s.archive", m.Config.DBName, time.Now().Format("20060102_150405"))
	
	args := []string{
		fmt.Sprintf("--host=%s", m.Config.Host),
		fmt.Sprintf("--port=%d", m.Config.Port),
		fmt.Sprintf("--username=%s", m.Config.User),
		fmt.Sprintf("--password=%s", m.Config.Password),
		fmt.Sprintf("--db=%s", m.Config.DBName),
		fmt.Sprintf("--archive=%s", filename),
	}

	cmdName := "mongodump"
	if m.Config.ToolPath != "" {
		cmdName = m.Config.ToolPath
	}

	cmd := exec.Command(cmdName, args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mongodump failed: %v, output: %s", err, string(output))
	}

	return filename, nil
}

func (m *MongoDB) Restore(backupFile string) error {
	args := []string{
		fmt.Sprintf("--host=%s", m.Config.Host),
		fmt.Sprintf("--port=%d", m.Config.Port),
		fmt.Sprintf("--username=%s", m.Config.User),
		fmt.Sprintf("--password=%s", m.Config.Password),
		fmt.Sprintf("--archive=%s", backupFile),
	}

	cmdName := "mongorestore"
	if m.Config.ToolPath != "" {
		dir := filepath.Dir(m.Config.ToolPath)
		cmdName = filepath.Join(dir, "mongorestore")
	}

	cmd := exec.Command(cmdName, args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("mongorestore failed: %v, output: %s", err, string(output))
	}

	return nil
}

func (m *MongoDB) Close() error {
	if m.client != nil {
		return m.client.Disconnect(context.TODO())
	}
	return nil
}
