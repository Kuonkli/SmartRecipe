package handlers

import (
	"SmartRecipe/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			log.Println("Missing access token, attempting to refresh")
			// Передаем исходный URL в параметре redirect_uri
			redirectURL := c.Request.URL.Path
			c.Redirect(http.StatusFound, "/o/refresh?redirect_uri="+redirectURL)
			c.Abort() // Остановить дальнейшую обработку
			return
		} else {
			log.Println("Access token found")
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			log.Println("Invalid token or parsing error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		log.Println("Token valid, proceeding")
		c.Set("userId", claims.UserID)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			c.Abort()
			return
		}
		userId, err := strconv.Atoi(userIdStr.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		user, err := database.Database.Users.GetUserById(userId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
			return
		}

		if user.RoleId != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this API"})
			c.Abort()
			return
		}
	}
}
