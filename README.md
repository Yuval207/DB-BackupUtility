# Database Backup CLI

A comprehensive CLI tool for backing up and restoring databases (MySQL, PostgreSQL, MongoDB) to local or cloud storage (AWS S3, Google Cloud Storage, Azure Blob Storage).

## Features

- **Multiple Database Support**: MySQL, PostgreSQL, MongoDB.
- **Flexible Storage**: Local filesystem, AWS S3, Google Cloud Storage, Azure Blob Storage.
- **Compression**: Gzip compression support to save space.
- **Notifications**: Slack integration for backup status updates.
- **Easy to Use**: Simple CLI interface with configuration file.

> **[Read the Detailed Documentation](DOCUMENTATION.md)** for architecture, workflows, and file responsibilities.

## Installation

### Prerequisites
- Go 1.21 or higher
- Database tools installed (mysqldump, pg_dump, mongodump, etc.)

### Build
```bash
make build
```

## Configuration

Create a `config.yaml` file:

```yaml
database:
  type: mysql
  host: localhost
  port: 3306
  user: root
  password: password
  dbname: testdb

storage:
  type: s3
  path: my-backup-bucket
  region: us-east-1

backup:
  type: full
  compression: true

notify:
  slack_webhook_url: "https://hooks.slack.com/..."
```

## Usage

### Backup
```bash
./dbbackup backup --config config.yaml
```

### Restore
```bash
./dbbackup restore <backup_file_name> --config config.yaml
```

## License
MIT
