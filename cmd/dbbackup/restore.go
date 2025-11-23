package main

import (
	"log"

	"github.com/antigravity/dbbackup/internal/restore"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore [backup_file]",
	Short: "Restore a database from a backup",
	Long:  `Restores the database from a specified backup file in the storage.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		backupFile := args[0]

		db, st, err := getComponents(appConfig)
		if err != nil {
			log.Fatalf("Error initializing components: %v", err)
		}
		defer db.Close()

		mgr := restore.NewManager(db, st)
		if err := mgr.PerformRestore(backupFile); err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
