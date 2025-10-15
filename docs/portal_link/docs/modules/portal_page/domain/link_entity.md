# Link

## 介紹

Link 實體代表使用者在 Portal Page 中展示的個別連結項目，包含社群媒體、個人網站、商店或任何外部連結。Link 是 Portal Page 聚合（Aggregate）內的實體，必須透過 Portal Page（聚合根）來管理，不能獨立存在。每個 Link 包含連結的標題、URL、排序順序等資訊。

## 屬性

| 屬性 | 型態 | 說明 |
|------|------|------|
| id | int | Link 的唯一標識符 |
| portal_page_id | int | 所屬的 Portal Page ID（外鍵關聯至 Portal Page，聚合根） |
| title | string | 連結的顯示標題 |
| url | string | 連結的目標 URL |
| description | string | 連結的描述或說明（選填） |
| icon_url | string | 連結的圖示 URL（選填） |
| display_order | int | 連結在頁面上的顯示順序 |
| created_at | timestamp | 連結建立時間 UTC |
| updated_at | timestamp | 連結資料更新時間 UTC |

## 業務規則

- Link 必須隸屬於一個 Portal Page，不能獨立存在
- Link 的新增、修改、刪除操作必須透過 Portal Page 聚合根來執行，以維護聚合的一致性

