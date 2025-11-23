package database

import (
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/antigravity/dbbackup/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Config config.DatabaseConfig
	conn   *sql.DB
}

func NewMySQL(cfg config.DatabaseConfig) *MySQL {
	return &MySQL{Config: cfg}
}

func (m *MySQL) Connect() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.Config.User, m.Config.Password, m.Config.Host, m.Config.Port, m.Config.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	m.conn = db
	return nil
}

func (m *MySQL) TestConnection() error {
	if m.conn == nil {
		if err := m.Connect(); err != nil {
			return err
		}
	}
	return m.conn.Ping()
}

func (m *MySQL) Backup(backupType string) (string, error) {
	// Note: mysqldump typically performs a full backup. 
	// Incremental backups in MySQL usually require binary logs, which is complex for a CLI tool.
	// We will stick to full backups for now unless 'incremental' logic is strictly required via binlogs.
	
	filename := fmt.Sprintf("backup_mysql_%s_%s.sql", m.Config.DBName, time.Now().Format("20060102_150405"))
	
	args := []string{
		fmt.Sprintf("-h%s", m.Config.Host),
		fmt.Sprintf("-P%d", m.Config.Port),
		fmt.Sprintf("-u%s", m.Config.User),
		fmt.Sprintf("-p%s", m.Config.Password),
		m.Config.DBName,
		"--result-file=" + filename,
	}

	cmdName := "mysqldump"
	if m.Config.ToolPath != "" {
		cmdName = m.Config.ToolPath
	}

	cmd := exec.Command(cmdName, args...)
	// Hide password from process list if possible, but passing as arg is standard for mysqldump in simple scripts.
	// A better way is using a config file or env var, but for now this is direct.
	// WARNING: -p with password directly can be insecure in shared environments. 
	// Ideally we write a temporary .my.cnf file.
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("mysqldump failed: %v, output: %s", err, string(output))
	}

	return filename, nil
}

func (m *MySQL) Restore(backupFile string) error {
	// mysql -h... -u... -p... dbname < backupFile
	
	args := []string{
		fmt.Sprintf("-h%s", m.Config.Host),
		fmt.Sprintf("-P%d", m.Config.Port),
		fmt.Sprintf("-u%s", m.Config.User),
		fmt.Sprintf("-p%s", m.Config.Password),
		m.Config.DBName,
	}

	cmdName := "mysql"
	if m.Config.ToolPath != "" {
		dir := filepath.Dir(m.Config.ToolPath)
		cmdName = filepath.Join(dir, "mysql")
	}

	cmd := exec.Command(cmdName, args...)
	
	file, err := os.Open(backupFile)
	if err != nil {
		return err
	}
	defer file.Close()
	
	cmd.Stdin = file
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("mysql restore failed: %v, output: %s", err, string(output))
	}

	return nil
}

func (m *MySQL) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}
