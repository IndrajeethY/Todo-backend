package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterAuthRoutes(rg *gin.RouterGroup, jwtExpiry time.Duration) {
	rg.POST("/login", func(c *gin.Context) {
		var in LoginInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		adminUser := os.Getenv("ADMIN_USER")
		adminPass := os.Getenv("ADMIN_PASS")
		adminID := os.Getenv("ADMIN_ID")
		if adminUser == "" || adminPass == "" || adminID == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server not configured"})
			return
		}
		if in.Username != adminUser || in.Password != adminPass {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		secret := os.Getenv("JWT_SECRET")
		now := time.Now()
		claims := jwt.RegisteredClaims{
			Subject:   adminID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtExpiry)),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": signed, "expires_in": int(jwtExpiry.Seconds())})
	})
}
