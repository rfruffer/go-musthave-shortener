package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

const (
	cookieName = "user_id"
)

func signUserID(userID, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(userID))
	return hex.EncodeToString(h.Sum(nil))
}

func validateCookie(userID, signature, secretKey string) bool {
	return hmac.Equal([]byte(signUserID(userID, secretKey)), []byte(signature))
}

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie(cookieName)
		if err != nil {
			userID := uuid.New().String()
			signature := signUserID(userID, secretKey)
			c.Set("user_id", userID)
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     cookieName,
				Value:    userID + "|" + signature,
				Path:     "/",
				HttpOnly: true,
				Expires:  time.Now().Add(365 * 24 * time.Hour),
			})
		} else {
			parts := []rune(cookie.Value)
			splitIndex := -1
			for i := len(parts) - 1; i >= 0; i-- {
				if parts[i] == '|' {
					splitIndex = i
					break
				}
			}
			if splitIndex == -1 {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			userID := string(parts[:splitIndex])
			signature := string(parts[splitIndex+1:])

			if !validateCookie(userID, signature, secretKey) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.Set("user_id", userID)
		}
		c.Next()
	}
}
