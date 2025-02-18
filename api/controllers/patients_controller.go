package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/AskatNa/OnlineClinic/api/models"
	"github.com/AskatNa/OnlineClinic/api/responses"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllPatients(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var patients []models.User

	query := bson.M{"role": "patient"}
	opts := options.Find().SetProjection(bson.D{{Key: "_id", Value: 0}})

	cursor, err := userCollection.Find(ctx, query, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching patients",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var patient models.User
		if err := cursor.Decode(&patient); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error decoding patient data",
				Data:    map[string]interface{}{"error": err.Error()},
			})
			return
		}
		patients = append(patients, patient)
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    map[string]interface{}{"patients": patients},
	})
}

func GetPatientById(c *gin.Context) {
	idParam := c.Param("id")
	patientDoc := getDPatientProfileByStringId(idParam)
	if len(patientDoc) == 0 {
		c.JSON(http.StatusNotFound, responses.UserResponse{
			Status:  http.StatusNotFound,
			Message: "Patient not found",
		})
		return
	}
	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    map[string]interface{}{"patient": patientDoc},
	})
}

func getDPatientProfileByStringId(strPatientId string) bson.M {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	patientId, err := primitive.ObjectIDFromHex(strPatientId)
	if err != nil {
		fmt.Println("Error converting id from hex:", err)
		return bson.M{}
	}

	query := bson.D{{Key: "_id", Value: patientId}}

	var patientDoc bson.M
	if err := UserCollection.FindOne(ctx, query).Decode(&patientDoc); err != nil {
		fmt.Println("Error fetching patient:", err)
		return bson.M{}
	}

	return patientDoc
}
