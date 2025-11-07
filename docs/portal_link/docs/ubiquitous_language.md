# 通用語言（Ubiquitous Language）

本表整理自 `User` 與 `Portal Page` 領域內的實體與錯誤、枚舉，供業務、設計與工程在討論與實作時共用一致語彙。

## 實體（Entities）

| 名稱 | 定義 | 關鍵屬性 | 備註 |
|------|------|----------|------|
| User | 使用 Portal Link 的使用者。系統核心領域物件，代表一個帳號的擁有者。 | id, name, email(唯一), password, created_at, updated_at | email 必須唯一，用於登入與識別 |
| Portal Page | 使用者在平台上的個人化連結整合頁面；系統中的核心物件，也是聚合根（Aggregate Root）。 | id, user_id, slug(唯一), title, bio, profile_image_url, theme, created_at, updated_at | 管理其內部的 Link 實體；所有 Link 變更需透過此聚合根 |
| Link | Portal Page 中的個別連結項目（如社群、網站、商店等）。 | id, portal_page_id, title, url, description?, icon_url?, display_order, created_at, updated_at | 隸屬於 Portal Page 的聚合內，不能獨立存在 |

## 值物件／枚舉（Value Objects / Enums）

| 名稱 | 定義 | 可能值 | 備註 |
|------|------|--------|------|
| Theme | Portal Page 的主題風格，控制頁面視覺與色彩配置。 | `light`（預設）、`dark` | 使用者可依偏好或品牌形象選擇 |

## 錯誤（Domain Errors）

| 錯誤碼 | 訊息 | 說明 | 所屬領域 |
|--------|------|------|----------|
| ErrInvalidParams | invalid parameters | 參數驗證失敗（格式錯誤、長度不符、必填欄位為空等） | User / Portal Page |
| ErrEmailExists | email already exists | Email 已存在於系統 | User |
| ErrInvalidCredentials | invalid credentials | 登入憑證錯誤（帳號或密碼錯誤） | User |
| ErrSlugExists | slug already exists | Slug 已被使用，無法建立或更新 | Portal Page |
| ErrPortalPageNotFound | portal page not found | 找不到指定的 Portal Page | Portal Page |
| ErrLinkNotFound | link not found | 找不到指定的 Link | Portal Page |

## 聚合與不變量

- 聚合根：`Portal Page`

  - 管理 `Link` 實體的新增、修改、刪除
  - 聚合邊界：所有 `Link` 變更必須透過 `Portal Page`

- 重要不變量：

  - `Portal Page.slug` 在系統中必須唯一
  - `Portal Page` 必須屬於有效的 `User`
  - `Link` 必須隸屬於一個 `Portal Page`，不可獨立存在

## 術語對照

| 術語 | 英文 | 定義 |
|------|------|------|
| 使用者 | User | 使用平台的人，擁有帳號的主體 |
| 入口頁/個人頁 | Portal Page | 個人化的連結聚合頁，對外公開展示 |
| 連結項目 | Link | 在 Portal Page 上展示的單一連結資料 |
| 主題 | Theme | Portal Page 的視覺主題（light/dark） |
| 識別代稱 | Slug | 用於組成 URL 的唯一識別字串（例如 `example.com/username`） |
