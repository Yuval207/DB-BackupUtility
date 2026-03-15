#!/bin/bash

# Configuration defaults
CONFIG_FILE=${CONFIG_FILE:-"$(pwd)/config.yaml"}
BACKUP_DIR=${BACKUP_DIR:-"$(pwd)/backups"}
IMAGE="yuval207/databasebackup:latest"

# Ensure the backup directory exists locally
mkdir -p "$BACKUP_DIR"

if [ "$1" == "backup" ]; then
    echo "🚀 Starting Database Backup..."
    docker run --rm \
      -v "$CONFIG_FILE":/app/config.yaml \
      -v "$BACKUP_DIR":/app/backups \
      -e CLOUDFLARE_API_TOKEN="$CLOUDFLARE_API_TOKEN" \
      -e CLOUDFLARE_ACCOUNT_ID="$CLOUDFLARE_ACCOUNT_ID" \
      "$IMAGE" backup --config /app/config.yaml

elif [ "$1" == "restore" ]; then
    if [ -z "$2" ]; then
        echo "❌ Error: Please provide the backup file name to restore."
        echo "Usage: ./dbbackup.sh restore <filename.sql.gz>"
        exit 1
    fi
    echo "🚀 Starting Database Restore for $2..."
    docker run --rm \
      -v "$CONFIG_FILE":/app/config.yaml \
      -v "$BACKUP_DIR":/app/backups \
      -e CLOUDFLARE_API_TOKEN="$CLOUDFLARE_API_TOKEN" \
      -e CLOUDFLARE_ACCOUNT_ID="$CLOUDFLARE_ACCOUNT_ID" \
      "$IMAGE" restore "$2" --config /app/config.yaml

else
    echo "❌ Unknown command: $1"
    echo "Usage: ./dbbackup.sh [backup|restore] [filename]"
    exit 1
fi
