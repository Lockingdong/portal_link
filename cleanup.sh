#!/bin/bash
# Cleanup script to remove resources created by initialize.sh

# 1. Remove .env file
if [ -f .env ]; then
  echo "Removing .env file..."
  rm .env
fi

# 2. Run database migration rollback
echo "Rolling back database migrations..."
task migrate-down

# 3. Stop and remove docker containers and volumes
echo "Stopping and removing docker containers and volumes..."
docker-compose down -v

# 4. Remove Python virtual environment directory
if [ -d venv ]; then
  echo "Removing Python virtual environment directory..."
  rm -rf venv
fi

# 5. Uninstall Taskfile CLI globally
if command -v task >/dev/null 2>&1; then
  echo "Uninstalling Taskfile CLI globally..."
  npm uninstall -g @go-task/cli
fi

echo "Cleanup completed."