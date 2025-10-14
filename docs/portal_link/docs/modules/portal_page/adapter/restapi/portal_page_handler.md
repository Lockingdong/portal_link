# Portal Page Handler

## 概述

`PortalPageHandler` 是個人頁面模組的 REST API 適配器層，負責處理個人頁面相關的 HTTP 請求，包括創建個人頁面功能。它作為 HTTP 層和業務邏輯層（Use Case）之間的橋樑。

### Use Cases

- `createPortalPageUC`: 創建個人頁面用例，處理頁面創建業務邏輯

### 身份驗證

所有需要身份驗證的端點都使用 [auth](../../../../pkg/auth.md) 套件提供的功能：

- 使用 `AuthMiddleware` 進行 token 驗證
- 使用 `GetUserIDFromContext` 取得已驗證的使用者 ID

## API Endpoints

### 創建個人頁面 (CreatePortalPage)

**Endpoint:** `POST /api/v1/portal-pages`

**Method:** CreatePortalPage

**路由設定:**

```go
router.POST("/api/v1/portal-pages", pkg.AuthMiddleware(), handler.CreatePortalPage)
```

**處理流程:**
1. 透過 AuthMiddleware 驗證用戶認證（需要 Bearer Token）
2. 使用 GetUserIDFromContext 取得當前登入的使用者 ID
3. 綁定並驗證 Request 為 `CreatePortalPageParams`
4. 調用 `createPortalPageUC.Execute()` 執行頁面創建邏輯（請參考 [create_portal_page_uc.md](../../usecase/create_portal_page_uc.md)）

**Request:**
```json
{
    "slug": "john-doe",
    "title": "John's Page",
    "bio": "Welcome to my personal page!",
    "profile_image_url": "https://example.com/images/john.jpg",
    "theme": "light"
}
```

**Response:** 
`201 Created`
```json
{
    "id": 1
}
```

**Error Response:**

`400 Bad Request` - 參數驗證失敗
```json
{
    "code": "ErrInvalidParams",
    "message": "Invalid request parameters"
}
```

`401 Unauthorized` - 未認證或 token 無效
```json
{
    "code": "ErrUnauthorized",
    "message": "Invalid access token"
}
```

`500 Internal Server Error` - 伺服器內部錯誤
```json
{
    "code": "ErrInternal",
    "message": "Internal server error"
}
```

**注意事項：**

- 所有請求都需要在 Authorization 標頭提供有效的 Bearer Token
- Token 格式必須為：`Bearer <token>`
- 無效的 token 會由 AuthMiddleware 自動返回 401 錯誤
- 使用者 ID 獲取失敗會返回 500 錯誤
- 所有錯誤響應都使用 http_error 套件統一處理
- Domain errors 會自動映射到對應的 HTTP 錯誤響應