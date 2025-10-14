# Sign Up

## 概述

此用例允許新使用者註冊 Portal Link 帳號。使用者需要提供基本資訊（姓名、電子郵件、密碼），系統將驗證資料並建立新的使用者帳號。

**主要參與者：** 訪客（未註冊使用者）

## 輸入參數

| 參數 | 型態 | 必填 | 說明 | 驗證規則 |
|------|------|------|------|----------|
| name | string | 是 | 使用者全名 | 長度 1-255 字元 |
| email | string | 是 | 電子郵件地址 | 必須符合 email 格式，長度 1-255 字元 |
| password | string | 是 | 使用者密碼 | 最少 8 字元，需包含英文和數字 |

## 輸出結果

**成功時：**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**注意：** 註冊成功後會自動登入，返回 access token

## 主要流程

1. 使用者提交註冊資訊（姓名、電子郵件、密碼）
2. 系統驗證輸入參數格式
3. 系統檢查電子郵件地址是否已被註冊
4. 系統建立新的 User 實體
5. 系統將使用者資訊存入資料庫
6. 系統使用 `GenerateAccessToken` 方法產生該 User 的 access_token（詳見 [Authentication](../../../pkg/auth.md)）
7. 系統返回 access_token

## 錯誤流程

### 輸入參數驗證失敗
- 系統返回錯誤 `ErrInvalidParams`

### 電子郵件已被註冊
- 系統返回錯誤 `ErrEmailExists`

## 業務規則

- 每個電子郵件地址只能註冊一個帳號
- 電子郵件地址不區分大小寫
- 使用者姓名和電子郵件不可為空
- 密碼暫時以明文存儲
  - **TODO:** 後續討論密碼加密方式（如 bcrypt）
- Access token 產生方式：請參考 [Authentication](../../../pkg/auth.md)

## 相關實體

- **User Entity**: 使用者領域實體
- **User Repository**: 使用者資料存取介面
