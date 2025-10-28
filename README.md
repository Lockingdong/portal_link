# Portal Link 開發指南

Portal Link 是一個文檔管理平台，使用 Go 後端和 MkDocs 生成文檔。

## 技術架構

**後端框架與工具：**
- **Gin** - Web 框架
- **Viper** - 配置管理

**文檔工具：**
- MkDocs with Material 主題

## 快速開始

### 環境要求
- Go 1.24+
- Python 3.9+ (用於 MkDocs)

### 安裝 MkDocs
```bash
python3 -m venv venv
source venv/bin/activate
python3 -m pip install mkdocs mkdocs-material
```

## 開發指令

```bash
# 啟動後端
go run main.go

# 啟動 MkDocs 文檔服務
cd docs/portal_link

mkdocs serve
```

---
