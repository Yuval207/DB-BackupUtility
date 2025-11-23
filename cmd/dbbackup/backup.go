package main

import (
	"log"

	"github.com/antigravity/dbbackup/internal/backup"
	"github.com/antigravity/dbbackup/internal/notifier"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Perform a database backup",
	Long:  `Initiates a backup of the configured database and uploads it to the configured storage.`,
	Run: func(cmd *cobra.Command, args []string) {
		db, st, err := getComponents(appConfig)
		if err != nil {
			log.Fatalf("Error initializing components: %v", err)
		}
		defer db.Close()

		notif := notifier.NewSlackNotifier(appConfig.Notify)
		mgr := backup.NewManager(db, st, appConfig.Backup, notif)
		if err := mgr.PerformBackup(); err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
