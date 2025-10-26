#!/bin/bash
# 1. 複制 .env
cp .env.example .env

### 2. 啟動資料庫
docker-compose up -d

### 3. 安裝 Goose（資料庫遷移工具）
go install github.com/pressly/goose/v3/cmd/goose@latest

### 5. 安裝 MkDocs
python3 -m venv venv
source venv/bin/activate
python3 -m pip install mkdocs mkdocs-material

## 使用 Taskfile
npm install -g @go-task/cli

### 4. 執行資料庫遷移
task migrate-up

