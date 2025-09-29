# User

## 介紹

User 實體代表使用 Portal Link 的使用者，是系統中的核心領域物件。此實體封裝了使用者的基本屬性，包括唯一識別碼、姓名、電子郵件地址以及密碼等重要資訊。

## 屬性

| 屬性 | 型態 | 說明 |
|------|------|------|
| id | int | 使用者的唯一標識符 |
| name | string | 使用者的全名 |
| email | string | 使用者的電子郵件地址，必須是唯一的 |
| password | string | 使用者的密碼 |
| created_at | timestamp | 使用者建立時間 |
| updated_at | timestamp | 使用者資料更新時間 |
