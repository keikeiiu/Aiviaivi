#!/bin/bash
# PostgreSQL backup script for AiliVili
# Usage: ./scripts/backup.sh [backup_dir]
#
# Recommended cron: 0 2 * * * /app/scripts/backup.sh /backups

set -euo pipefail

BACKUP_DIR="${1:-./backups}"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
DB_NAME="${DB_NAME:-ailivili}"
DB_USER="${DB_USER:-ailivili}"
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
RETENTION_DAYS="${RETENTION_DAYS:-7}"

mkdir -p "$BACKUP_DIR"

BACKUP_FILE="$BACKUP_DIR/ailivili_${TIMESTAMP}.sql.gz"

echo "[$(date)] Starting backup of $DB_NAME to $BACKUP_FILE"

PGPASSWORD="${PGPASSWORD:-}" pg_dump \
    -h "$DB_HOST" \
    -p "$DB_PORT" \
    -U "$DB_USER" \
    -d "$DB_NAME" \
    --no-owner \
    --no-acl \
    | gzip > "$BACKUP_FILE"

if [ $? -eq 0 ]; then
    echo "[$(date)] Backup complete: $(du -h "$BACKUP_FILE" | cut -f1)"
else
    echo "[$(date)] Backup FAILED" >&2
    exit 1
fi

# Remove backups older than retention period
find "$BACKUP_DIR" -name "ailivili_*.sql.gz" -mtime "+$RETENTION_DAYS" -delete 2>/dev/null

echo "[$(date)] Cleaned backups older than $RETENTION_DAYS days"
echo "[$(date)] Current backups: $(ls "$BACKUP_DIR" | wc -l) files"
