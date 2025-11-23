package database

// Database interface defines the methods that any database provider must implement
type Database interface {
	// Connect establishes a connection to the database
	Connect() error

	// TestConnection verifies that the database is reachable and credentials are valid
	TestConnection() error

	// Backup performs a database backup and returns the path to the backup file
	// The backupType can be "full", "incremental", or "differential"
	Backup(backupType string) (string, error)

	// Restore restores the database from the given backup file
	Restore(backupFile string) error

	// Close closes the database connection
	Close() error
}
