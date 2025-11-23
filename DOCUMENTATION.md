# Database Backup CLI - Technical Documentation

## 1. Project Overview
The **Database Backup CLI** is a robust tool designed to automate the backup and restoration of databases. It supports multiple database engines (MySQL, PostgreSQL, MongoDB) and storage providers (Local, AWS S3, Google Cloud Storage, Azure Blob Storage). It also includes features for compression and Slack notifications.

## 2. Architecture
The project follows a modular architecture using Go's standard layout.

### High-Level Design
1.  **CLI Layer (`cmd/`)**: Handles user input, flags, and configuration loading. Uses `Cobra` for command management.
2.  **Manager Layer (`internal/backup`, `internal/restore`)**: Orchestrates the business logic. It connects the Database and Storage components.
3.  **Provider Layer (`internal/database`, `internal/storage`)**: Implements specific logic for each database and storage type behind common interfaces.
4.  **Support Layer (`internal/config`, `internal/logger`, `internal/notifier`)**: Provides utilities for configuration, logging, and notifications.

## 3. Directory Structure & File Responsibilities

### `cmd/dbbackup/`
Entry point of the application.
-   `main.go`: Initializes the application and executes the root command.
-   `root.go`: Defines the root command and global flags (like `--config`). Initializes the configuration system (`Viper`).
-   `backup.go`: Implements the `backup` command. Initializes the `BackupManager`.
-   `restore.go`: Implements the `restore` command. Initializes the `RestoreManager`.
-   `utils.go`: Factory functions to instantiate the correct Database and Storage providers based on configuration.

### `internal/database/`
Contains database implementations.
-   `interface.go`: Defines the `Database` interface (`Connect`, `Backup`, `Restore`, `Close`).
-   `mysql.go`: MySQL implementation. Uses `mysqldump` and `mysql` binaries.
-   `postgres.go`: PostgreSQL implementation. Uses `pg_dump` and `psql` binaries. Handles `sslmode` and custom tool paths.
-   `mongodb.go`: MongoDB implementation. Uses `mongodump` and `mongorestore` binaries.

### `internal/storage/`
Contains storage implementations.
-   `interface.go`: Defines the `Storage` interface (`Upload`, `Download`, `List`, `Delete`).
-   `local.go`: Local filesystem storage.
-   `s3.go`: AWS S3 implementation using the AWS SDK.
-   `gcs.go`: Google Cloud Storage implementation.
-   `azure.go`: Azure Blob Storage implementation.

### `internal/backup/`
-   `manager.go`: The `BackupManager`. It coordinates the backup process:
    1.  Tests DB connection.
    2.  Calls `DB.Backup()` to generate a dump file.
    3.  Calls `CompressFile()` (if enabled).
    4.  Calls `Storage.Upload()` to save the file.
    5.  Sends Slack notifications on success/failure.
-   `compression.go`: Helper functions for Gzip compression and decompression.

### `internal/restore/`
-   `manager.go`: The `RestoreManager`. It coordinates the restore process:
    1.  Calls `Storage.Download()` to retrieve the backup.
    2.  Calls `DecompressFile()` (if needed).
    3.  Calls `DB.Restore()` to apply the dump to the database.

### `internal/config/`
-   `config.go`: Defines the configuration structs (`Config`, `DatabaseConfig`, `StorageConfig`, etc.) that map to `config.yaml`.

### `internal/notifier/`
-   `slack.go`: Implements Slack webhook notifications.

## 4. Workflows

### 4.1 Backup Workflow
1.  **Start**: User runs `./dbbackup backup --config config.yaml`.
2.  **Init**: `cmd/dbbackup/main.go` loads the config file.
3.  **Factory**: `cmd/dbbackup/utils.go` creates instances of the specific Database (e.g., `Postgres`) and Storage (e.g., `S3`) providers.
4.  **Execution**: `internal/backup/manager.go` takes control.
    *   **Connect**: Verifies database connectivity.
    *   **Dump**: Executes the external tool (e.g., `pg_dump`) to create a local `.sql` or `.archive` file.
    *   **Compress**: If `compression: true`, gzips the file.
    *   **Upload**: Uploads the file to the configured storage destination.
    *   **Notify**: Sends a "Backup successful" message to Slack.
5.  **Cleanup**: Deletes the temporary local files.

### 4.2 Restore Workflow
1.  **Start**: User runs `./dbbackup restore <filename> --config config.yaml`.
2.  **Init**: Config is loaded.
3.  **Factory**: Database and Storage providers are instantiated.
4.  **Execution**: `internal/restore/manager.go` takes control.
    *   **Download**: Downloads the specified file from storage to a local temporary path.
    *   **Decompress**: If the file ends in `.gz`, it is decompressed.
    *   **Restore**: Executes the external tool (e.g., `psql`) to feed the file back into the database.
        *   *Note*: The tool attempts to find the restore binary in the same directory as the configured backup binary.
5.  **Cleanup**: Deletes the temporary local files.

## 5. Configuration Guide
The `config.yaml` file drives the behavior of the tool.

```yaml
database:
  type: postgres          # Options: mysql, postgres, mongodb
  host: localhost
  port: 5432
  user: myuser
  password: mypassword
  dbname: mydb
  extra_params: "sslmode=require"  # Optional: Extra connection params
  tool_path: "/usr/bin/pg_dump"    # Optional: Path to the dump binary

storage:
  type: s3                # Options: local, s3, gcs, azure
  path: my-bucket-name    # Bucket name (cloud) or directory path (local)
  region: us-east-1       # Required for S3
  credentials_file: ""    # Optional: Path to cloud credentials file

backup:
  type: full              # Currently only 'full' is supported
  compression: true       # Enable Gzip compression

notify:
  slack_webhook_url: "..." # Optional: Slack Webhook URL
```
