package controllers

import (
	"context"
	"github.com/AskatNa/OnlineClinic/api/models"
	"net/http"
	"time"

	"github.com/AskatNa/OnlineClinic/api/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection

func SetUserCollection(collection *mongo.Collection) {
	userCollection = collection
}
func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid data format",
		})
		return
	}

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error creating user",
		})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{
		Status:  http.StatusCreated,
		Message: "User successfully created",
		Data:    map[string]interface{}{"id": result.InsertedID},
	})
}

func GetAUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(userId)

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    map[string]interface{}{"user": user},
	})
}

func GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var users []models.User
	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error retrieving users",
		})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    map[string]interface{}{"users": users},
	})
}

func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(userId)

	var updatedData models.User
	if err := c.ShouldBindJSON(&updatedData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid data format",
		})
		return
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": updatedData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error updating user",
		})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data:    map[string]interface{}{"matchedCount": result.MatchedCount},
	})
}

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userId := c.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error deleting user",
		})
		return
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "User deleted successfully",
		Data:    map[string]interface{}{"deletedCount": result.DeletedCount},
	})
}
