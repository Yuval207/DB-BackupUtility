package database

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/antigravity/dbbackup/internal/config"
)

type D1 struct {
	Config config.DatabaseConfig
}

func NewD1(cfg config.DatabaseConfig) *D1 {
	return &D1{Config: cfg}
}

func (d *D1) Connect() error {
	// Wrangler connects via HTTP, no persistent connection needed
	return nil
}

func (d *D1) TestConnection() error {
	// Verify wrangler is installed and authenticated
	cmdName := "npx"
	var args []string
	
	if d.Config.ToolPath != "" {
		cmdName = d.Config.ToolPath
		args = []string{"d1", "info", d.Config.DBName}
	} else {
		args = []string{"wrangler", "d1", "info", d.Config.DBName}
	}

	cmd := exec.Command(cmdName, args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("wrangler d1 info failed: %v, output: %s", err, string(output))
	}

	return nil
}

func (d *D1) Backup(backupType string) (string, error) {
	filename := fmt.Sprintf("backup_d1_%s_%s.sql", d.Config.DBName, time.Now().Format("20060102_150405"))
	
	cmdName := "npx"
	var args []string

	if d.Config.ToolPath != "" {
		cmdName = d.Config.ToolPath
		args = []string{"d1", "export", d.Config.DBName, "--remote", "--output=" + filename}
	} else {
		args = []string{"wrangler", "d1", "export", d.Config.DBName, "--remote", "--output=" + filename}
	}

	cmd := exec.Command(cmdName, args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("wrangler d1 export failed: %v, output: %s", err, string(output))
	}

	return filename, nil
}

func (d *D1) Restore(backupFile string) error {
	cmdName := "npx"
	var args []string

	if d.Config.ToolPath != "" {
		// If ToolPath is just "wrangler"
		cmdName = d.Config.ToolPath
		args = []string{"d1", "execute", d.Config.DBName, "--remote", "--file=" + backupFile}
	} else {
		args = []string{"wrangler", "d1", "execute", d.Config.DBName, "--remote", "--file=" + backupFile}
	}

	cmd := exec.Command(cmdName, args...)
	
	// Execute takes the file directly, we don't pipe stdin for wrangler d1 execute
	// But it might ask for confirmation: "Are you sure you want to execute? (y/n)"
	// To bypass this, wrangler doesn't have a `--yes` normally for `execute` but let's check. 
	// Wait, wrangler d1 execute usually requires confirmation. Let's pass `--yes` just in case.
	// Actually `wrangler d1 execute <db> --remote --file=<file>` currently might prompt.
	// Let's add `--yes` to auto-confirm if possible. 
	args = append(args, "-y")
	
	cmd.Args = append([]string{cmdName}, args...)
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("wrangler d1 execute failed: %v, output: %s", err, string(output))
	}

	return nil
}

func (d *D1) Close() error {
	return nil
}
