# User Handler

## 概述

`UserHandler` 是用戶模組的 REST API 適配器層，負責處理用戶相關的 HTTP 請求，包括用戶註冊和登入功能。它作為 HTTP 層和業務邏輯層（Use Case）之間的橋樑。

### Use Cases

- `signUpUC`: 用戶註冊用例，處理註冊業務邏輯
- `signInUC`: 用戶登入用例，處理登入業務邏輯

## API Endpoints

### 用戶註冊 (SignUp)

**Endpoint:** `POST /api/v1/user/signup`

**Method:** SignUp

**處理流程:**

1. 綁定並驗證 Request 為 `SignUpParams`
2. 調用 `signUpUC.Execute()` 執行註冊邏輯 （請參考 [sign_up_uc.md](../../usecase/sign_up_uc.md)）

** Request:**
```json
{
    "name": "用戶名稱",
    "email": "user@example.com",
    "password": "用戶密碼"
}
```

**Response:** 
`200 OK`
```json
{
    "access_token": "access token..."
}
```

**Error Response:**

`400 Bad Request` - 參數驗證失敗
```json
{
    "error": "ErrInvalidParams",
    "message": "Invalid request parameters"
}
```

`400 Bad Request` - 電子郵件已存在
```json
{
    "error": "ErrEmailExists",
    "message": "Email already exists"
}
```

`500 Internal Server Error` - 伺服器內部錯誤
```json
{
    "error": "ErrInternal",
    "message": "Internal server error"
}
```

### 用戶登入 (SignIn)

**Endpoint:** `POST /api/v1/user/signin`

**Method:** SignIn

**處理流程:**

1. 綁定並驗證 Request 為 `SignInParams`
2. 調用 `signInUC.Execute()` 執行登入邏輯（請參考 [sign_in_uc.md](../../usecase/sign_in_uc.md)）

**Request:**
```json
{
    "email": "user@example.com",
    "password": "用戶密碼"
}
```

**Response:** 
`200 OK`
```json
{
    "access_token": "access token..."
}
```

**Error Response:**

`400 Bad Request` - 參數驗證失敗
```json
{
    "error": "ErrInvalidParams",
    "message": "Invalid request parameters"
}
```

`401 Unauthorized` - 認證失敗
```json
{
    "error": "ErrInvalidCredentials",
    "message": "Invalid email or password"
}
```

`500 Internal Server Error` - 伺服器內部錯誤
```json
{
    "error": "ErrInternal",
    "message": "Internal server error"
}
```
