# Portal Page

## 介紹

Portal Page 實體代表使用者在 Portal Link 平台上的個人化連結整合頁面，是系統中的核心領域物件，同時也是 **聚合根（Aggregate Root）**。此實體封裝了使用者專屬頁面的基本屬性，包括唯一識別碼、頁面識別名稱、個人簡介、主題設定等重要資訊。每個使用者可以擁有一個專屬的 Portal Page，用於展示所有重要的社群、網站、商店連結給朋友、粉絲或客戶。

作為聚合根，Portal Page 負責管理其內部的所有 Link 實體，確保聚合內部的一致性和業務規則。所有對 Link 的操作都必須透過 Portal Page 來進行。

## 屬性

| 屬性 | 型態 | 說明 |
|------|------|------|
| id | int | Portal Page 的唯一標識符 |
| user_id | int | 擁有此頁面的使用者 ID（外鍵關聯至 User） |
| slug | string | 頁面的 URL 識別名稱，必須是唯一的（例如：example.com/username） |
| title | string | 頁面標題或顯示名稱 |
| bio | string | 使用者的個人簡介或描述 |
| profile_image_url | string | 個人頭像圖片的 URL |
| theme | string | 頁面主題設定（例如：顏色、樣式等） |
| created_at | timestamp | 頁面建立時間 UTC |
| updated_at | timestamp | 頁面資料更新時間 UTC |

## 聚合設計

Portal Page 作為聚合根，負責管理以下實體：

- **Link**：Portal Page 中的連結項目（請參考 [link_entity](link_entity.md)）

### 業務規則

- 一個 Portal Page 的 `slug` 在系統中必須是唯一的
- Portal Page 必須屬於一個有效的使用者（User）
