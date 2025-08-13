#!/bin/sh
set -e

echo "Backup job started..."

while true; do
  # Jalankan backup tepat pada jam 00:00 (tengah malam)
  sleep_time=$(($(date -d 'tomorrow 00:00:00' '+%s') - $(date '+%s')))
  echo "Sleeping for $sleep_time seconds until midnight..."
  sleep $sleep_time

  # Format nama file: backup-srikandisehat-2025-08-14.sql.gz
  FILE_NAME="backup-srikandisehat-$(date +'%Y-%m-%d').sql.gz"
  BACKUP_PATH="/backups/$FILE_NAME"

  echo "Creating backup: $BACKUP_PATH"

  # Lakukan dump database dan kompres langsung dengan gzip
  mariadb-dump -h db -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" | gzip > "$BACKUP_PATH"

  echo "Backup created successfully."

  # Hapus backup yang lebih tua dari 7 hari
  echo "Cleaning up old backups (older than 7 days)..."
  find /backups -name "backup-*.sql.gz" -type f -mtime +7 -delete
  echo "Cleanup complete."
done