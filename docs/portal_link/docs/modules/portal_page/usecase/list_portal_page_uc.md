# List Portal Pages

## 概述

此用例允許已登入使用者查看自己建立的所有 Portal Pages 列表。使用者可以在管理介面查看所有頁面，以便進行編輯、查看 Portal Pages。

**主要參與者：** 已登入使用者

## 輸入參數

| 參數 | 型態 | 必填 | 說明 | 驗證規則 |
|------|------|------|------|----------|
| user_id | int | 是 | 使用者 ID | 必須為有效的已登入使用者 |

## 輸出結果

**成功時：** 返回 Portal Pages 列表

| 欄位 | 型態 | 說明 |
|------|------|------|
| portal_pages | array | PortalPage 實體陣列 |

**PortalPage 欄位：**

| 欄位 | 型態 | 說明 |
|------|------|------|
| id | int | Portal Page ID |
| slug | string | 自訂網址 |
| title | string | 頁面標題 |

## 主要流程

1. 系統透過 Repository 查詢該使用者的所有 Portal Pages

    - 根據建立時間升冪排序

2. 系統返回完整的 Portal Pages 列表

## 錯誤流程

- 錯誤情境：

    - 無

## 業務規則

### 列表查詢規則

- **預設排序**：按建立時間升冪排序
- **空結果處理**：若使用者尚未建立任何 Portal Page，返回空陣列
- **返回所有資料**：一次返回該使用者的所有 Portal Pages，不進行分頁

## 相關物件

- **PortalPage Entity**: Portal Page 領域實體（聚合根）
- **PortalPage Repository**: Portal Page 資料存取介面
