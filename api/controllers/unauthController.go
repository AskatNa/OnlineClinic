package controllers

import (
	"context"
	"fmt"
	"github.com/AskatNa/OnlineClinic/api/models"
	"github.com/AskatNa/OnlineClinic/api/responses"
	"github.com/AskatNa/OnlineClinic/config/configs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"regexp"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, responses.UserResponse{
		Status: http.StatusOK, Message: "Pong",
	})
}

// Registration func
func RegisterUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var data models.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status: http.StatusBadRequest, Message: "Invalid data format",
		})
		return
	}

	fmt.Println("Received Data:", data)
	if len(data.Password) < 8 {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Password must be at least 8 characters long",
		})
		return
	}
	if !emailRegex.MatchString(data.Email) {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid email format",
		})
		return
	}
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": data.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, responses.UserResponse{
			Status:  http.StatusConflict,
			Message: "Email already registered",
		})
		return
	}
	data.Role = "patient"
	result, err := userCollection.InsertOne(ctx, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError,
			Message: "Error registering user"})
		return
	}
	c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "User successfully registered",
		Data:    map[string]interface{}{"id": result.InsertedID, "role": data.Role}})
}

//func RegisterUser(c *gin.Context) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	var data models.User
//	if err := c.ShouldBindJSON(&data); err != nil {
//		c.JSON(http.StatusBadRequest, responses.UserResponse{
//			Status: http.StatusBadRequest, Message: "Invalid data format",
//		})
//		return
//	}
//
//	fmt.Println("Received Data:", data)
//	if len(data.Password) < 8 {
//		c.JSON(http.StatusBadRequest, responses.UserResponse{
//			Status:  http.StatusBadRequest,
//			Message: "Password must be at least 8 characters long",
//		})
//		return
//	}
//	if !emailRegex.MatchString(data.Email) {
//		c.JSON(http.StatusBadRequest, responses.UserResponse{
//			Status:  http.StatusBadRequest,
//			Message: "Invalid email format",
//		})
//		return
//	}
//	var existingUser models.User
//	err := userCollection.FindOne(ctx, bson.M{"email": data.Email}).Decode(&existingUser)
//	if err == nil {
//		c.JSON(http.StatusConflict, responses.UserResponse{
//			Status:  http.StatusConflict,
//			Message: "Email already registered",
//		})
//		return
//	}
//	data.Role = "patient"
//	result, err := userCollection.InsertOne(ctx, data)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError,
//			Message: "Error registering user"})
//		return
//	}
//	c.JSON(http.StatusCreated, responses.UserResponse{
//		Status:  http.StatusCreated,
//		Message: "User successfully registered",
//		Data:    map[string]interface{}{"id": result.InsertedID, "role": data.Role}})
//}

//var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

//func Login(c *gin.Context) {
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	var data map[string]string
//	if err := c.ShouldBindJSON(&data); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid login data format"})
//		return
//	}
//	var user models.User
//	err := userCollection.FindOne(ctx, bson.M{"email": data["email"]}).Decode(&user)
//	if err != nil || user.Password != data["password"] {
//		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
//		return
//	}
//	expirationTime := time.Now().Add(time.Hour * 24)
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"email": user.Email,
//		"role":  user.Role,
//		"exp":   expirationTime.Unix(),
//	})
//	tokenString, err := token.SignedString([]byte(configs.JWTSecret))
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
//		return
//	}
//	c.JSON(http.StatusOK, gin.H{
//		"message":    "Login successful",
//		"token":      "Bearer " + tokenString,
//		"expireTime": expirationTime,
//	})
//}

// Login func
func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest,
			Message: "Invalid login data format"})
		return
	}
	// Validate email
	if !emailRegex.MatchString(data["email"]) {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid email format",
		})
		return
	}
	var user models.User
	if data["email"] == configs.AdminEmail && data["password"] == configs.AdminPassword {
		expirationTime := time.Now().Add(24 * time.Hour)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email":   configs.AdminEmail,
			"isAdmin": true, //user.Email == configs.AdminEmail,
			"exp":     expirationTime.Unix(),
		})
		tokenString, err := token.SignedString([]byte(configs.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error generating token",
			})
			return
		}
		fmt.Println("Generated Token for", data["email"], ":", tokenString)
		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "Admin login successful",
			Data: map[string]interface{}{
				"user":       user,
				"email":      configs.AdminEmail,
				"isAdmin":    true,
				"token":      "Bearer " + tokenString,
				"expireTime": expirationTime,
			},
		})
		return
	}
	err := userCollection.FindOne(ctx, bson.M{"email": data["email"]}).Decode(&user)
	if err != nil || user.Password != data["password"] {
		c.JSON(http.StatusUnauthorized, responses.UserResponse{Status: http.StatusUnauthorized,
			Message: "Invalid credentials"})
		return
	}
	expirationTime := time.Now().Add(time.Hour * 24)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"role":    user.Role,
		"isAdmin": user.Role == "admin" || user.Email == configs.AdminEmail,
		"exp":     expirationTime.Unix(),
	})
	//tokenString, err := token.SignedString(configs.JWTSecret)
	tokenString, err := token.SignedString([]byte(configs.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error generating token",
		})
		return
	}
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization")
	c.Writer.Header().Set("Authorization", "Bearer "+tokenString)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   "Bearer " + tokenString,
		//"isAdmin":    user.Email == configs.AdminEmail,
		"user":       user,
		"expireTime": expirationTime,
	})
	fmt.Println("Generated Token for", data["email"], ":", tokenString)
}
