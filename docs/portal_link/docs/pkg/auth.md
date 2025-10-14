# 身份驗證

## Access Token 產生方式

系統在使用者成功登入或註冊後會產生 access token。目前採用簡單的實作方式，未來會進行增強。

### 目前實作
- Token 格式：將使用者 ID 和過期時間 timestamp 進行 base64 編碼
- Token 有效期：產生後 1 天內有效

### 未來優化計畫（TODO）
- 實作 JWT（JSON Web Tokens）以提供更安全和功能豐富的 token 機制
- 新增 token 刷新機制
- 考慮實作 token 撤銷功能

## 方法

### GenerateAccessToken

產生使用者的 access token。

```go
func GenerateAccessToken(userID string) (string, error)
```

**參數：**
- `userID`: 使用者 ID

**回傳：**
- `string`: 產生的 access token
- `error`: 如果產生過程中發生錯誤則回傳錯誤

**處理流程：**
1. 取得當前時間
2. 計算過期時間（當前時間 + 1 天）
3. 組合使用者 ID 和過期時間
4. 進行 base64 編碼
5. 回傳編碼後的 token

### ValidateAccessToken

驗證 access token 的有效性。

```go
func ValidateAccessToken(token string) (string, error)
```

**參數：**
- `token`: 要驗證的 access token

**回傳：**
- `string`: token 對應的使用者 ID
- `error`: 如果驗證失敗則回傳錯誤

**處理流程：**
1. 解碼 base64 token
2. 解析出使用者 ID 和過期時間
3. 檢查是否已過期
4. 回傳使用者 ID

**錯誤類型：**
- `ErrInvalidToken`: token 格式不正確
- `ErrExpiredToken`: token 已過期
- `ErrTokenValidationFailed`: 驗證過程發生錯誤

### AuthMiddleware

Gin 框架的身份驗證中間件，用於保護需要登入的 API 端點。

```go
func AuthMiddleware() gin.HandlerFunc
```

**功能：**
- 從請求標頭獲取並驗證 access token
- 將驗證後的使用者 ID 存入 context
- 處理驗證失敗的情況

**使用方式：**
```go
router := gin.Default()
router.Use(AuthMiddleware()) // 全域使用
// 或
router.GET("/protected", AuthMiddleware(), handleProtected) // 單一路由使用
```

**處理流程：**
1. 從 Authorization 標頭獲取 Bearer token
2. 使用 ValidateAccessToken 驗證 token
3. 如果驗證成功：

    - 將使用者 ID 存入 gin.Context
    - 調用下一個 handler

4. 如果驗證失敗：

    - 中止請求
    - 返回對應的錯誤訊息

**錯誤回應：**
```json
{
  "error": {
    "code": "ErrUnauthorized",
    "message": "Invalid access token"
  }
}
```

**從 Context 取得使用者 ID：**
```go
func GetUserIDFromContext(c *gin.Context) (string, error)
```

**使用範例：**
```go
func handleProtected(c *gin.Context) {
    userID, err := GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // 使用 userID 進行後續處理
}
```

**注意事項：**

- Authorization 標頭格式必須為：`Bearer <token>`
- 未提供 token 或格式錯誤會返回 401 Unauthorized
- token 過期或無效會返回 401 Unauthorized
- Context 中找不到使用者 ID 會返回 500 Internal Server Error
