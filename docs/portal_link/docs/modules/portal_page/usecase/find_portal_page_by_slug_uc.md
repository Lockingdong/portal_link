# Find Portal Page By Slug

## 概述

此用例允許訪客（未註冊使用者）根據 `slug` 取得公開的單一 Portal Page 完整資料（含 Links）。不需登入，無擁有權限制。

**主要參與者：** 訪客（未註冊使用者）

## 輸入參數

| 參數 | 型態 | 必填 | 說明 | 驗證規則 |
|------|------|------|------|----------|
| slug | string | 是 | Portal Page 自訂網址 | 必須為非空字串 |

## 輸出結果

**成功時：** 返回完整的 PortalPage 實體

**PortalPage 欄位：**

| 欄位 | 型態 | 說明 |
|------|------|------|
| id | int | Portal Page ID |
| slug | string | 自訂網址 |
| title | string | 頁面標題 |
| bio | string | 個人簡介（可為空字串） |
| profile_image_url | string | 個人頭像 URL（可為空字串） |
| theme | string | 主題風格，可選值：`light`、`dark` |
| links | array | Link 實體陣列 |

**Link 欄位：**

| 欄位 | 型態 | 說明 |
|------|------|------|
| id | int | Link ID |
| title | string | 連結標題 |
| url | string | 連結網址 |
| description | string | 連結描述（可為空字串） |
| icon_url | string | 圖示 URL（可為空字串） |
| display_order | int | 顯示順序（正整數） |

## 主要流程

1. 驗證輸入參數：

    - 檢查 `slug` 是否為非空字串

2. 透過 Repository `FindBySlug` 取得指定 `slug` 的 Portal Page（含 Links）

    - 載入 Links 並依 `display_order` 升冪排序

3. 系統返回完整的 Portal Page 實體（含 Links）

## 錯誤流程

### 輸入參數驗證失敗
- 系統返回錯誤 `ErrInvalidParams`
- 錯誤情境：

    - `slug` 未提供或為空字串

### 找不到頁面
- 系統返回錯誤 `ErrPortalPageNotFound`
- 說明：查無符合 `slug` 的 Portal Page

## 業務規則

### 查詢與權限規則
- 僅支援以 `slug` 查詢
- 公開可讀，不需登入，無擁有權限制

### 返回資料
- 返回完整的 Portal Page（含 Links）
- 若頁面尚未設定任何 Link，`links` 為空陣列
- Links 以 `display_order` 升冪排序返回

## 相關物件

- **PortalPage Entity**: Portal Page 領域實體（聚合根）
- **Link Entity**: Link 領域實體（聚合內實體）
- **PortalPage Repository**: Portal Page 資料存取介面（支援 `FindBySlug`）
