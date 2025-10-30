#!/bin/sh
set -e

# Allow override by environment variables
DB_PATH="${DB_PATH:-/data/mangahub.db}"
SQL_FILE="${SQL_FILE:-/data/init.sql}"

echo "DB_PATH=$DB_PATH"
echo "SQL_FILE=$SQL_FILE"

if [ ! -f "$DB_PATH" ]; then
    echo "Initializing SQLite database..."
    sqlite3 "$DB_PATH" < "$SQL_FILE"
    echo "Database initialized successfully."
else
    echo "Database already exists, skipping init."
fi
