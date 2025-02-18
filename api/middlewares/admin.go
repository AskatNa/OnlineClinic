package middlewares

import (
	"fmt"
	"github.com/AskatNa/OnlineClinic/config/configs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Received Headers:", c.Request.Header)

		authHeader := c.GetHeader("Authorization")
		fmt.Println("Authorization Header:", c.GetHeader("Authorization"))

		fmt.Println("Extracted Authorization Header:", authHeader)

		fmt.Println("JWT Secret Key:", configs.JWTSecret)

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			fmt.Println("No valid token found")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized - No Token"})
			c.Abort()
			return
		}
		fmt.Println("Token Found:", authHeader)

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		fmt.Println("Extracted Token:", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.JWTSecret), nil
		})
		fmt.Println("JWT Secret Key:", configs.JWTSecret)

		if err != nil || !token.Valid {
			fmt.Println("JWT Error:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			c.Abort()
			return
		}

		userEmail, ok := claims["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Email missing in token"})
			c.Abort()
			return
		}
		fmt.Println("Extracted Email:", userEmail)
		fmt.Println("Extracted Token:", tokenString)
		fmt.Println("Claims:", claims)

		isAdmin := false
		if adminFlag, exists := claims["isAdmin"].(bool); exists && adminFlag {
			isAdmin = true
		}
		fmt.Println("Is Admin:", isAdmin)

		c.Set("email", userEmail)
		c.Set("isAdmin", isAdmin)
		fmt.Println("Middleware - Extracted Email:", claims["email"], "Is Admin:", claims["isAdmin"])

		c.Next()
	}
}
