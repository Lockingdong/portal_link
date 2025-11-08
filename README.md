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
- Node.js 20+ (用於前端開發)
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

## 快速設定腳本

為了簡化開發環境的設定和清理，提供了兩個便利腳本：

### 初始化腳本 (`initialize.sh`)

自動化完成所有開發環境設定：

```bash
./initialize.sh
```

這個腳本會執行以下操作：
- 複製 `.env.example` 為 `.env` 配置檔案
- 啟動 Docker 資料庫容器
- 安裝 Goose 資料庫遷移工具
- 建立 Python 虛擬環境並安裝 MkDocs
- 安裝 Taskfile CLI 工具
- 執行資料庫遷移設定

### 清理腳本 (`cleanup.sh`)

完全清除開發環境和相關資源：

```bash
./cleanup.sh
```

這個腳本會執行以下清理操作：
- 刪除 `.env` 配置檔案
- 回滾所有資料庫遷移
- 停止並移除 Docker 容器及資料卷
- 刪除 Python 虛擬環境目錄
- 卸載全域安裝的 Taskfile CLI

**注意：** 執行清理腳本會完全移除開發環境，包括資料庫資料。請在執行前確保已備份重要資料。

## 開發指令

### 後端

```bash
# 啟動後端
go run main.go

# 生成 SQLBoiler 模型
sqlboiler psql
```

### 前端

```bash
# 進入前端目錄
cd web

# 安裝依賴（首次運行）
npm install

# 啟動前端開發服務器
npm run dev
```

### 文檔

```bash
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
