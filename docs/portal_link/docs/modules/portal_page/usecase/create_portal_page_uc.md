# Create Portal Page

## 概述

此用例允許已登入使用者建立個人化的連結整合頁面（Portal Page）。使用者可以設定頁面的基本資訊，包括自訂網址（slug）、標題、個人簡介、頭像和主題風格。

**主要參與者：** 已登入使用者

## 輸入參數

| 參數 | 型態 | 必填 | 說明 | 驗證規則 |
|------|------|------|------|----------|
| user_id | int | 是 | 使用者 ID | 必須為有效的已登入使用者 |
| slug | string | 是 | 頁面的自訂網址 | 長度 3-50 字元，只能包含小寫英文字母、數字和連字號（-），不可以連字號開頭或結尾，必須唯一 |
| title | string | 是 | 頁面標題 | 長度 1-100 字元 |
| bio | string | 否 | 個人簡介 | 最多 500 字元 |
| profile_image_url | string | 否 | 個人頭像 URL | 必須為有效的 URL 格式 |
| theme | string | 否 | 主題風格 | 可選值：`light`、`dark`，預設為 `light` |

## 輸出結果

**成功時：** 返回建立成功的 PortalPage 實體

| 欄位 | 型態 | 說明 |
|------|------|------|
| id | int | Portal Page ID（自動生成） |

## 主要流程

1. 使用者提交建立 Portal Page 的請求（包含 user_id、slug、title、bio 等資訊）
2. 系統驗證輸入參數格式

    - 驗證 slug 格式（長度、字元限制、不可為保留字）
    - 驗證 title 不為空且長度符合規範
    - 驗證 bio 長度（如有提供）
    - 驗證 profile_image_url 格式（如有提供）
    - 驗證 theme 為有效值（如有提供）

3. 系統檢查 slug 是否已被使用
4. 系統建立新的 PortalPage 實體

    - 設定 user_id
    - 設定 slug（轉換為小寫）
    - 設定 title、bio、profile_image_url
    - 設定 theme（未提供時預設為 `light`）
    - 初始化 links 為空陣列
    - 設定 created_at 和 updated_at 為當前時間

5. 系統透過 Repository 將 Portal Page 存入資料庫
6. 系統返回建立成功的 PortalPage 實體

## 錯誤結果

### 輸入參數驗證失敗
- 系統返回錯誤 `ErrInvalidParams`
- 錯誤情境：

    - slug 格式不正確
    - slug 為保留字
    - title 為空或長度超過限制
    - bio 長度超過限制
    - profile_image_url 格式錯誤
    - theme 不是有效值

### Slug 已存在
- 系統返回錯誤 `ErrSlugExists`
- 說明：slug 在系統中必須唯一，若已被其他使用者使用則無法建立

## 業務規則

### Slug 規則
- **唯一性**：每個 slug 在整個系統中必須是唯一的
- **格式限制**：

    - 只能包含小寫英文字母（a-z）
    - 數字（0-9）
    - 連字號（-）作為分隔符
    - 不可以連字號開頭或結尾
    - 不可包含連續的連字號

- **長度限制**：3-50 字元
- **大小寫處理**：系統自動將 slug 轉換為小寫
- **保留字**：系統保留某些 slug 不可使用

    - 例如：`admin`、`api`、`static`、`public`、`auth`、`login`、`signup`、`help`、`about`、`terms`、`privacy`
    - **TODO:** 定義完整的保留字清單並實作驗證


### 頁面限制
- 每個使用者可以建立多個 Portal Pages

    - **TODO:** 後續討論是否需要限制每個使用者可建立的頁面數量（如免費版 1 個，付費版無限制）

### 預設值
- `theme` 未提供時預設為 `light`
- `bio` 和 `profile_image_url` 為選填，可以為空字串
- 新建立的 Portal Page 的 `links` 為空陣列

### 驗證優先順序
1. 參數格式驗證
2. Slug 唯一性檢查
3. 資料庫操作

## 相關物件

- **PortalPage Entity**: Portal Page 領域實體（聚合根）
- **Link Entity**: Link 領域實體（聚合內實體）
- **User Entity**: 使用者領域實體
- **PortalPage Repository**: Portal Page 資料存取介面
