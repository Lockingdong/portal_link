# Portal Link 開發指南

Portal Link 是一個文檔管理平台，使用 Go 後端和 MkDocs 生成文檔。

## 技術架構

**後端框架與工具：**
- **Gin** - Web 框架
- **SQLBoiler** - ORM 工具
- **Viper** - 配置管理
- **Goose** - 資料庫遷移工具

**資料庫：**
- PostgreSQL

**文檔工具：**
- MkDocs with Material 主題

## 快速開始

### 1. 環境要求
- Go 1.24+
- Docker & Docker Compose
- Python 3.9+ (用於 MkDocs)

### 2. 啟動資料庫
```bash
docker-compose up -d
```

### 3. 安裝 Goose（資料庫遷移工具）
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### 4. 執行資料庫遷移
```bash
# 執行遷移
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=portal_link sslmode=disable search_path=portal_link" up

# 回滾遷移
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=portal_link sslmode=disable search_path=portal_link" down

# 查看狀態
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=portal_link sslmode=disable search_path=portal_link" status
```

### 5. 安裝 MkDocs
```bash
python3 -m pip install mkdocs mkdocs-material
```

## 開發指令

```bash
# 啟動後端
go run main.go

# 生成 SQLBoiler 模型
sqlboiler psql

# 啟動 MkDocs 文檔服務
mkdocs serve
```

## 使用 Taskfile

Taskfile 讓常用開發指令更一致。請先安裝 Task：

```bash
brew install go-task/tap/go-task
```

列出可用任務：

```bash
task --list
```

常用任務對照：

```bash
# 執行資料庫遷移（向上）
task migrate-up

# 回滾資料庫遷移（向下）
task migrate-down

# 生成 SQLBoiler 模型
task generate-models

# 啟動 MkDocs 文檔服務
task mkdocs-serve
```

---
