# Database Backup CLI

A comprehensive CLI tool for backing up and restoring databases (MySQL, PostgreSQL, MongoDB) to local or cloud storage (AWS S3, Google Cloud Storage, Azure Blob Storage).

## Features

- **Multiple Database Support**: MySQL, PostgreSQL, MongoDB, Cloudflare D1.
- **Flexible Storage**: Local filesystem, AWS S3, Google Cloud Storage, Azure Blob Storage.
- **Compression**: Gzip compression support to save space.
- **Notifications**: Slack integration for backup status updates.
- **Easy to Use**: Simple CLI interface with configuration file.

> **[Read the Detailed Documentation](DOCUMENTATION.md)** for architecture, workflows, and file responsibilities.

## Installation

### Method 1: Docker (Recommended)
The easiest way to run the CLI is using the pre-built Docker image, which already includes all necessary database client tools (`mysql`, `pg_dump`, `mongodump`, etc.).

```bash
docker pull yuval207/databasebackup:latest
```

### Method 2: Build from Source
**Prerequisites**
- Go 1.25 or higher
- Database tools installed on your host system (`mysqldump`, `pg_dump`, `mongodump`, etc.)

**Build**
```bash
make build
```

## Configuration

Create a `config.yaml` file:

```yaml
database:
  type: d1
  dbname: testing-db
  
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

### 🚀 The Easy Way (Wrapper Script)
We've included a handy `dbbackup.sh` script so you don't have to remember the long Docker commands!

**1. Make the script executable (One-time only):**
```bash
chmod +x dbbackup.sh
```

**2. Export Cloudflare D1 Credentials (If using D1):**
If you are backing up a Cloudflare D1 database, your session needs authentication:
```bash
export CLOUDFLARE_API_TOKEN="your_cloudflare_token_here"
export CLOUDFLARE_ACCOUNT_ID="your_cloudflare_account_id_here"
```

**3. Run a Backup:**
By default, it looks for `config.yaml` in the current folder and saves to `./backups`.
```bash
./dbbackup.sh backup
```

**4. Run a Restore:**
```bash
./dbbackup.sh restore <backup_file_name.sql.gz>
```

**Advanced Usage (Custom Configs/Folders):**
If you are using this in a separate project folder, just copy the `dbbackup.sh` script there and use it! You can override the config file and backup directory:
```bash
CONFIG_FILE=./my-project-config.yaml BACKUP_DIR=./my-backups ./dbbackup.sh backup
```

---

### Method 2: Raw Docker Commands
If you prefer running Docker manually or are writing CI/CD pipelines:

**Backup**
Mount your local directory containing `config.yaml` into the container:
```bash
docker run --rm \
  -v $(pwd)/config.yaml:/app/config.yaml \
  yuval207/databasebackup:latest backup --config /app/config.yaml
```

**Restore**
If you are restoring from a default local directory, you also need to mount that backups directory:
```bash
docker run --rm \
  -v $(pwd)/config.yaml:/app/config.yaml \
  yuval207/databasebackup:latest restore <backup_file_name> --config /app/config.yaml
```

### Using Local Binary

**Backup**
```bash
./dbbackup backup --config config.yaml
```

**Restore**
```bash
./dbbackup restore <backup_file_name> --config config.yaml
```

## Project URL
```bash
https://roadmap.sh/projects/database-backup-utility
```
