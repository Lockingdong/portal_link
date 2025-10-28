package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"portal_link/modules/user/domain"

	"github.com/gin-gonic/gin"
)

// TODO: 將常數移至配置文件，方便根據環境調整
const (
	// TokenExpiration 定義 token 的有效期為 1 天
	TokenExpiration = 24 * time.Hour
	// ContextUserIDKey 用於在 gin.Context 中存儲使用者 ID 的鍵值
	ContextUserIDKey = "userID"
)

// TODO: 擴充錯誤類型以支援更多驗證場景
var (
	ErrInvalidToken  = errors.New("invalid token format")
	ErrExpiredToken  = errors.New("token has expired")
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user ID in token")
)

// TODO: 考慮加入額外的安全相關欄位，如 token 版本、裝置識別碼等
type tokenData struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GenerateAccessToken 產生使用者的 access token
func GenerateAccessToken(userID string) (string, error) {
	// TODO: 改用 JWT 實作，提供更安全和標準的 token 機制
	// TODO: 加入 token 簽名機制確保資料完整性
	// TODO: 考慮使用 Redis 等快取系統儲存 token 狀態

	// 建立 token 資料
	data := tokenData{
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(TokenExpiration),
	}

	// 序列化資料
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token data: %w", err)
	}

	// Base64 編碼
	token := base64.StdEncoding.EncodeToString(jsonData)
	return token, nil
}

// ValidateAccessToken 驗證 access token 的有效性，並檢查使用者是否存在
func ValidateAccessToken(ctx context.Context, token string, userRepo domain.UserRepository) (string, error) {
	// TODO: 實作 token 黑名單機制，支援 token 撤銷功能
	// TODO: 加入 token 使用紀錄，以便追蹤可疑活動
	// TODO: 實作 rate limiting 機制防止暴力破解

	// Base64 解碼
	jsonData, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", ErrInvalidToken
	}

	// 解析 token 資料
	var data tokenData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return "", ErrInvalidToken
	}

	// 檢查是否過期
	if time.Now().UTC().After(data.ExpiresAt) {
		return "", ErrExpiredToken
	}

	// 將 userID 從字串轉換為整數
	userID, err := strconv.Atoi(data.UserID)
	if err != nil {
		return "", ErrInvalidUserID
	}

	// 檢查使用者是否存在於資料庫中
	_, err = userRepo.Find(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		log.Printf("Error finding user: %v", err)
		return "", ErrUserNotFound
	}

	return data.UserID, nil
}

// AuthMiddleware Gin 框架的身份驗證中間件
func AuthMiddleware(userRepo domain.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 支援多種身份驗證方式（如 API Key、OAuth 等）
		// TODO: 加入請求來源驗證（CORS 設定）
		// TODO: 實作 token 刷新機制

		// 從 Authorization 標頭獲取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "ErrUnauthorized",
				"message": "Invalid access token",
			})
			return
		}

		// 檢查 Bearer token 格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "ErrUnauthorized",
				"message": "Invalid access token",
			})
			return
		}

		// 驗證 token 並檢查使用者是否存在
		userID, err := ValidateAccessToken(c.Request.Context(), parts[1], userRepo)
		if err != nil {
			log.Println("ValidateAccessToken error:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "ErrUnauthorized",
				"message": "Invalid access token",
			})
			return
		}

		// 將使用者 ID 存入 context
		c.Set(ContextUserIDKey, userID)
		c.Next()
	}
}

// GetUserIDFromContext 從 gin.Context 中取得使用者 ID
func GetUserIDFromContext(c *gin.Context) (string, error) {
	// TODO: 考慮加入使用者角色和權限的快取機制
	// TODO: 提供更豐富的使用者資訊（如 email、角色等）

	userID, exists := c.Get(ContextUserIDKey)
	if !exists {
		return "", errors.New("user ID not found in context")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID type in context")
	}

	return userIDStr, nil
}
