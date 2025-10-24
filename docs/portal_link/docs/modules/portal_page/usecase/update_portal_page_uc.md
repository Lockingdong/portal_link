# Update Portal Page

## 概述

此用例允許已登入使用者更新其現有的連結整合頁面（Portal Page）。使用者可以修改頁面的基本資訊，包括自訂網址（slug）、標題、個人簡介、頭像和主題風格，以及管理頁面中的連結（Links）。

**主要參與者：** 已登入使用者

## 輸入參數

| 參數 | 型態 | 必填 | 說明 | 驗證規則 |
|------|------|------|------|----------|
| user_id | int | 是 | 使用者 ID | 必須為有效的已登入使用者，且必須是頁面擁有者 |
| portal_page_id | int | 是 | Portal Page ID | 必須存在且屬於該使用者 |
| slug | string | 否 | 頁面的自訂網址 | 長度 3-50 字元，只能包含小寫英文字母、數字和連字號（-），不可以連字號開頭或結尾，必須唯一（若更新） |
| title | string | 否 | 頁面標題 | 長度 1-100 字元（若更新） |
| bio | string | 否 | 個人簡介 | 最多 500 字元 |
| profile_image_url | string | 否 | 個人頭像 URL | 必須為有效的 URL 格式 |
| theme | string | 否 | 主題風格 | 可選值：`light`、`dark` |
| links | array | 是 | 頁面中的連結清單 | 連結陣列，每個連結包含以下欄位 |

### Link 欄位結構

| 欄位 | 型態 | 必填 | 說明 | 驗證規則 |
|------|------|------|------|----------|
| id | int | 否 | Link ID | 若提供則為更新現有連結，若不提供則為新增連結 |
| title | string | 是 | 連結的顯示標題 | 長度 1-100 字元 |
| url | string | 是 | 連結的目標 URL | 必須為有效的 URL 格式 |
| description | string | 否 | 連結的描述或說明 | 最多 500 字元 |
| icon_url | string | 否 | 連結的圖示 URL | 必須為有效的 URL 格式（若提供） |
| display_order | int | 是 | 連結在頁面上的顯示順序 | 必須為正整數 |

## 輸出結果

**成功時：** 返回更新後的 PortalPage 實體

| 欄位 | 型態 | 說明 |
|------|------|------|
| id | int | Portal Page ID |

## 主要流程

1. 使用者提交更新 Portal Page 的請求
2. 系統驗證使用者權限

    - 確認使用者是頁面擁有者

3. 系統驗證輸入參數格式

    - 驗證 slug 格式（如有更新）
    - 驗證 title 長度（如有更新）
    - 驗證 bio 長度（如有更新）
    - 驗證 profile_image_url 格式（如有更新）
    - 驗證 theme 為有效值（如有更新）
    - 驗證 links 陣列（必填）
        - 驗證 title 長度
        - 驗證 URL 格式
        - 驗證 description 長度
        - 驗證 icon_url 格式
        - 驗證 link id 存在（若為更新）
        - 驗證 display_order 為正整數

4. 如果要更新 slug，系統檢查新的 slug 是否已被使用（排除當前頁面）
5. 系統從資料庫讀取現有的 PortalPage 實體
6. 系統更新 PortalPage 實體

    - 更新提供的欄位值
    - 保持未提供欄位的原有值不變
    - 處理 links 的更新

        - 新增新的連結
        - 更新現有連結
        - 儲存前端傳來的 display_order
    
    - 更新 updated_at 為當前時間

7. 系統透過 Repository 將更新後的 Portal Page 及其連結存入資料庫
8. 系統返回更新後的 PortalPage 實體

## 錯誤流程

### 權限驗證失敗

- 系統返回錯誤 `ErrUnauthorized`
- 錯誤情境：

    - 使用者不是頁面擁有者

### 頁面不存在

- 系統返回錯誤 `ErrPortalPageNotFound`
- 說明：指定的 portal_page_id 不存在

### 輸入參數驗證失敗

- 系統返回錯誤 `ErrInvalidParams`
- Portal Page 欄位驗證規則請參考 [Create Portal Page Use Case](create_portal_page_uc.md#輸入參數驗證失敗)
- Link 欄位錯誤情境：
    - link title 為空或長度超過限制
    - link URL 格式錯誤
    - link description 長度超過限制
    - link icon_url 格式錯誤
    - link display_order 不是正整數
    - link id 不存在或不屬於該頁面（若為更新）

### Slug 已存在

- 系統返回錯誤 `ErrSlugExists`
- 說明：新的 slug 在系統中已被其他使用者使用

## 業務規則

### 更新規則

- 只更新有提供的欄位，未提供的欄位保持原值
- 不允許將必填欄位（如 title）更新為空值
- 更新時會自動更新 updated_at 時間戳

### Links 更新規則

- Link 必須隸屬於一個 Portal Page，不能獨立存在
- Link 的新增、修改、刪除操作必須透過 Portal Page 聚合根來執行，以維護聚合的一致性
- 可以同時進行新增、更新和刪除連結操作
- 新增連結時：

    - 不需提供 id，系統會自動生成
    - 連結會自動關聯到當前更新的 Portal Page
    - created_at 和 updated_at 由系統自動設定為當前時間

- 更新連結時：

    - 必須提供有效的 id
    - 必須確認連結屬於當前 Portal Page
    - updated_at 會自動更新為當前時間
    - 其他欄位保持不變（除非有提供新值）

- 刪除連結時：

    - 直接從資料庫中刪除該連結

- 排序規則：

    - display_order 必須為正整數
    - 系統直接儲存前端傳來的 display_order，不做自動排序處理


### Slug 規則

請參考 [Create Portal Page Use Case](create_portal_page_uc.md#slug-規則) 的 Slug 規則。
唯一的差異是在檢查 slug 唯一性時，需要排除當前頁面本身。

### 驗證優先順序

1. 使用者權限驗證
2. 頁面存在性驗證
3. 參數格式驗證
4. Slug 唯一性檢查（如有更新 slug）
5. 資料庫操作

## 相關物件

- **PortalPage Entity**: Portal Page 領域實體（聚合根）
- **Link Entity**: Link 領域實體（聚合內實體）
- **User Entity**: 使用者領域實體
- **PortalPage Repository**: Portal Page 資料存取介面
