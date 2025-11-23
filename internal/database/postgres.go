package database

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/antigravity/dbbackup/internal/config"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Config config.DatabaseConfig
	conn   *sql.DB
}

func NewPostgres(cfg config.DatabaseConfig) *Postgres {
	return &Postgres{Config: cfg}
}

func (p *Postgres) Connect() error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", 
		p.Config.Host, p.Config.Port, p.Config.User, p.Config.Password, p.Config.DBName)
	
	if p.Config.ExtraParams != "" {
		dsn = fmt.Sprintf("%s %s", dsn, p.Config.ExtraParams)
	} else {
		dsn = fmt.Sprintf("%s sslmode=disable", dsn)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	p.conn = db
	return nil
}

func (p *Postgres) TestConnection() error {
	if p.conn == nil {
		if err := p.Connect(); err != nil {
			return err
		}
	}
	return p.conn.Ping()
}

func (p *Postgres) Backup(backupType string) (string, error) {
	filename := fmt.Sprintf("backup_pg_%s_%s.sql", p.Config.DBName, time.Now().Format("20060102_150405"))
	
	// PGPASSWORD env var is safer than command line arg
	os.Setenv("PGPASSWORD", p.Config.Password)
	defer os.Unsetenv("PGPASSWORD")

	args := []string{
		"-h", p.Config.Host,
		"-p", fmt.Sprintf("%d", p.Config.Port),
		"-U", p.Config.User,
		"-f", filename,
		p.Config.DBName,
	}

	cmdName := "pg_dump"
	if p.Config.ToolPath != "" {
		cmdName = p.Config.ToolPath
	}

	cmd := exec.Command(cmdName, args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("pg_dump failed: %v, output: %s", err, string(output))
	}

	return filename, nil
}

func (p *Postgres) Restore(backupFile string) error {
	os.Setenv("PGPASSWORD", p.Config.Password)
	defer os.Unsetenv("PGPASSWORD")

	args := []string{
		"-h", p.Config.Host,
		"-p", fmt.Sprintf("%d", p.Config.Port),
		"-U", p.Config.User,
		"-d", p.Config.DBName,
		"-f", backupFile,
	}



	// For restore we use psql, but we might want a separate config for it?
	// For now, let's assume if ToolPath is set, it points to the *dump* tool.
	// We might need a separate RestoreToolPath or just rely on psql being in PATH.
	// Actually, usually if pg_dump is missing, psql is too.
	// Let's try to infer psql path from pg_dump path if possible, or just use "psql" default.
	// A better approach for the user is to add the bin dir to PATH.
	// But let's stick to "psql" default for now as the user specifically failed on backup.
	
	cmdName := "psql"
	if p.Config.ToolPath != "" {
		// If ToolPath is set (e.g. /path/to/pg_dump), try to find psql in the same dir
		// or just replace "pg_dump" with "psql"
		// Simple heuristic: replace base name
		// This is a bit hacky but works for standard setups
		// A better way would be a separate RestoreToolPath config
		dir := filepath.Dir(p.Config.ToolPath)
		cmdName = filepath.Join(dir, "psql")
	}

	cmd := exec.Command(cmdName, args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("psql restore failed: %v, output: %s", err, string(output))
	}

	return nil
}

func (p *Postgres) Close() error {
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
