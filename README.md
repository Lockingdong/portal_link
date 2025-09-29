# Portal Link 開發指南

歡迎來到 Portal Link 專案！這個專案旨在提供一個強大的文檔平台，使用 MkDocs 來生成和管理文檔。

## 專案概覽

Portal Link 是一個用於創建和管理文檔的靜態網站生成器，旨在幫助開發者和團隊輕鬆地維護和分享項目文檔。

## Golang 開發

要開始 Golang 開發，請確保您的系統上已安裝 Golang。您可以選擇直接安裝特定版本或使用 GVM 來管理多個版本。

### 直接安裝 Golang 1.24

1. **下載 Golang 1.24 安裝包：**
   - 前往 [Golang 官方網站](https://golang.org/dl/) 下載 Golang 1.24 版本的安裝包，適合您操作系統的版本。

2. **安裝 Golang：**
   - 根據您的操作系統，運行下載的安裝包並按照指示完成安裝。

3. **驗證安裝：**
   - 打開終端機並輸入以下命令以確認 Golang 是否安裝成功：
   ```bash
   go version
   ```
   - 如果安裝成功，您應該會看到 Golang 的版本號。

### 使用 GVM 管理 Golang 版本

1. **安裝 GVM：**
   - GVM（Go Version Manager）允許您輕鬆地安裝和切換不同的 Golang 版本。
   ```bash
   bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
   source ~/.gvm/scripts/gvm
   ```

2. **使用 GVM 安裝 Golang 1.24：**
   ```bash
   gvm install go1.24
   gvm use go1.24 --default
   ```

3. **驗證安裝：**
   ```bash
   go version
   ```
   - 如果安裝成功，您應該會看到 Golang 1.24 的版本號。

---

### Golang 使用的套件

在本專案中，我們使用以下 Golang 套件來構建應用程式：

1. **Gin**
   - Gin 是一個用於構建高效能 Web 應用的 Go Web 框架。
   - 提供快速的路由和中間件支持。

2. **SQLBoiler**
   - SQLBoiler 是一個 ORM 工具，用於生成 Go 的 SQL 資料庫模型。
   - 支持多種資料庫，並提供強大的查詢生成功能。

3. **Viper**
   - Viper 是一個完整的配置解決方案，用於 Go 應用程式。
   - 支持多種配置格式（如 JSON、TOML、YAML 等）和環境變數。

---

## 安裝 PostgreSQL

要使用 Docker Compose 設定您的本地開發環境以支援 PostgreSQL，請按照以下步驟操作：

1. **確保您的系統上已安裝 Docker 和 Docker Compose。**

2. **導航到您的專案目錄：**
   ```bash
   cd /Users/lidongying/Documents/Projects/portal_link
   ```

3. **運行 Docker Compose：**
   ```bash
   docker-compose up -d
   ```
   此命令將啟動在 `docker-compose.yml` 文件中定義的 Golang 和 PostgreSQL 服務。

4. **驗證服務是否正在運行：**
   - PostgreSQL 服務應可在端口 5432 上訪問。

---

## SQL Migrate 安裝指南

SQL Migrate 是一個用於管理資料庫結構變更的工具，讓您可以輕鬆地執行和回滾資料庫遷移。

### 安裝 SQL Migrate

1. **使用 Go 安裝 SQL Migrate：**
   ```bash
   go install -tags 'postgres' github.com/rubenv/sql-migrate/...@latest
   ```
   這將安裝支援 PostgreSQL 的 SQL Migrate 工具。

2. **驗證安裝：**
   ```bash
   sql-migrate --version
   ```
   如果安裝成功，您應該會看到 SQL Migrate 的版本號。

### 配置 SQL Migrate

1. **創建配置文件：**
   在專案根目錄創建 `dbconfig.yml` 文件：
   ```yaml
   development:
     dialect: postgres
     datasource: host=localhost port=5432 user=postgres password=password dbname=portal_link sslmode=disable
     dir: migrations
     table: migrations
   
   production:
     dialect: postgres
     datasource: $DATABASE_URL
     dir: migrations
     table: migrations
   ```

2. **設定環境變數（可選）：**
   ```bash
   export DATABASE_URL="postgres://postgres:password@localhost:5432/portal_link?sslmode=disable"
   ```

### 使用 SQL Migrate

1. **執行遷移（向上）：**
   ```bash
   sql-migrate up
   ```
   這將執行所有待執行的遷移文件。

2. **回滾遷移（向下）：**
   ```bash
   sql-migrate down
   ```
   這將回滾最後一次遷移。

3. **查看遷移狀態：**
   ```bash
   sql-migrate status
   ```
   這將顯示所有遷移的執行狀態。

4. **執行特定數量的遷移：**
   ```bash
   sql-migrate up -limit=1
   sql-migrate down -limit=1
   ```

### 遷移文件格式

遷移文件應使用以下格式：
```sql
-- +migrate Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS users;
```

---

## MkDocs 安裝指南

### 先決條件
- 系統上已安裝 Python
- pip，Python 的套件安裝工具

### 安裝步驟

1. **打開終端機。**
   
2. **導航到您的專案目錄：**
   ```bash
   cd /Users/lidongying/Documents/Projects/portal_link
   ```

3. **使用 pip 安裝 MkDocs：**
   ```bash
   python3 -m pip install mkdocs
   ```
   此命令將下載並安裝 MkDocs 及其依賴項。

4. **驗證安裝：**
   執行以下命令以確保 MkDocs 安裝正確：
   ```bash
   mkdocs --version
   ```
   如果安裝成功，您應該會看到 MkDocs 的版本號。

5. **更新 PATH**

   為了能夠從任何地方運行 MkDocs 命令，您需要將 MkDocs 的安裝路徑添加到您的 PATH 中。

   **打開終端機並運行以下命令：**
   ```bash
   echo 'export PATH="$PATH:/Users/lidongying/Library/Python/3.9/bin"' >> ~/.zshrc
   source ~/.zshrc
   ```
   這將更新您的 PATH 並使更改生效。

6. **安裝 MkDocs Material 主題：**
   要使用 MkDocs Material 主題，請使用以下命令安裝：
   ```bash
   python3 -m pip install mkdocs-material
   ```
   這將安裝 MkDocs Material 主題，您可以在 `mkdocs.yaml` 中設置主題為 `material`。

---

本文檔提供了安裝 MkDocs 的逐步指南，MkDocs 是一個用於專案文件的靜態網站生成器。在進行安裝之前，請確保您已安裝 Python 和 pip。

---
