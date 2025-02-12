package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AskatNa/OnlineClinic/api/customsturctures"
	"github.com/AskatNa/OnlineClinic/api/helpers"
	"github.com/AskatNa/OnlineClinic/api/models"
	"github.com/AskatNa/OnlineClinic/api/responses"
	"github.com/AskatNa/OnlineClinic/config/configs"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	bookedAppointmentsCollection *mongo.Collection = configs.GetCollection(configs.DB, "bookedAppointments")
)

type SlotUpdateData struct {
	PatientID string
	Duration  int
	isBooked  bool
}

func BookAppointmentSlot(c *gin.Context) {
	fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	requestData := new(customsturctures.BookSlotRequest)
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "failed to parse request",
			Data:    map[string]interface{}{"error": "couldn't validate request body"},
		})
		return
	}

	if helpers.RoleValidator(requestData.Role, "patient") != "allowed" {
		c.JSON(http.StatusUnauthorized, responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "failed",
			Data:    map[string]interface{}{"problem": "Only patients are allowed to book appointment slots!"},
		})
		return
	}

	requestSlotData := requestData.Slotdata

	doctorObjId, err := primitive.ObjectIDFromHex(requestSlotData.DoctorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid DoctorID",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}

	query := bson.D{{Key: "_id", Value: doctorObjId}}

	var doctorDoc bson.M
	if err := UserCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error finding doctor",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}

	intSlotNo, err := strconv.Atoi(requestSlotData.SlotNo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error parsing SlotNo to integer",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}

	var newSlotData SlotUpdateData
	newSlotData.PatientID = requestSlotData.PatientID
	newSlotData.Duration, err = strconv.Atoi(requestSlotData.Duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error parsing Duration to integer",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}
	newSlotData.isBooked = true

	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, requestSlotData.AppointmentDay, int32(intSlotNo), newSlotData)

	bookedAppointmentItem := models.BookedAppointment{
		PatientId: requestSlotData.PatientID,
		DoctorId:  requestSlotData.DoctorID,
		SlotNo:    intSlotNo,
		Day:       requestSlotData.AppointmentDay,
		Duration:  newSlotData.Duration,

		Date: time.Now(),
	}
	insertItemToBookedAppointmentsCollection(bookedAppointmentItem)

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{"bookedSlot": updatedSlot},
	})
}

func CancelAppointmentSlot(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	requestData := new(customsturctures.BookSlotRequest)
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "failed",
			Data:    map[string]interface{}{"error": "couldn't validate request body"},
		})
		return
	}

	if requestData.Role != "doctor" && requestData.Role != "admin" {
		c.JSON(http.StatusUnauthorized, responses.UserResponse{
			Status:  http.StatusUnauthorized,
			Message: "failed",
			Data:    map[string]interface{}{"problem": "Only Doctor & Admins are allowed to cancel an appointment!"},
		})
		return
	}

	requestSlotData := requestData.Slotdata
	doctorObjId, err := primitive.ObjectIDFromHex(requestSlotData.DoctorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Invalid DoctorID",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}

	query := bson.D{{Key: "_id", Value: doctorObjId}}

	var doctorDoc bson.M
	if err := UserCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println("Error finding doctor:", err)
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed",
			Data:    map[string]interface{}{"problem": err.Error()},
		})
		return
	}

	intSlotNo, err := strconv.Atoi(requestSlotData.SlotNo)
	if err != nil {
		fmt.Println("Error parsing SlotNo to integer:", err)
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "failed",
			Data:    map[string]interface{}{"problem": err.Error()},
		})
		return
	}

	var newSlotData SlotUpdateData
	updatedSlot := UpdateAppointmentSlot(doctorObjId, doctorDoc, requestSlotData.AppointmentDay, int32(intSlotNo), newSlotData)

	deleteItemFromBookedAppointmentsCollection(requestSlotData.DoctorID, requestSlotData.AppointmentDay, int32(intSlotNo))

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{"canceledSlot": updatedSlot},
	})
}

func ExtractAppoinmentSlotFromDoctorProfile(doctorProfile primitive.M, slotDay string, slotNo int32) interface{} {

	ds, ok := doctorProfile["schedule"]
	if !ok {
		return nil
	}
	ws, ok := ds.(primitive.M)["weeklyschedule"]
	if !ok {
		return nil
	}
	day, ok := ws.(primitive.M)[slotDay]
	if !ok {
		return nil
	}
	appointmentsSlots, ok := day.(primitive.M)["appointmentslots"]
	if !ok {
		appointmentsSlots, ok = day.(primitive.M)["appointments_slots"]
		if !ok {
			return nil
		}
	}
	appointmentsArray, ok := appointmentsSlots.(primitive.A)
	if !ok || int(slotNo-1) >= len(appointmentsArray) {
		return nil
	}
	slot := appointmentsArray[slotNo-1]
	return slot
}

func UpdateAppointmentSlot(doctorObjId primitive.ObjectID, doctorProfile primitive.M,
	slotDay string, slotNo int32, newSlotData SlotUpdateData) interface{} {

	slot := ExtractAppoinmentSlotFromDoctorProfile(doctorProfile, slotDay, slotNo)
	if slot == nil {
		fmt.Printf("Slot not found for day %s, slot number %d\n", slotDay, slotNo)
		return nil
	}

	slotMap, ok := slot.(primitive.M)
	if !ok {
		fmt.Println("Slot format error")
		return nil
	}

	slotMap["patientid"] = newSlotData.PatientID
	slotMap["duration"] = newSlotData.Duration
	slotMap["isbooked"] = newSlotData.isBooked

	newSchedule := doctorProfile["schedule"]

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	updatedResult, err := UserCollection.UpdateOne(
		ctx,
		bson.M{"_id": doctorObjId},
		bson.D{{Key: "$set", Value: bson.D{{Key: "schedule", Value: newSchedule}}}},
	)
	if err != nil {
		fmt.Println("Error updating doctor schedule:", err)
		return nil
	}
	fmt.Printf("Updated %v Documents!\n", updatedResult.ModifiedCount)

	return slot
}

func ViewAppointmentDetails(c *gin.Context) {
	var requestData map[string]interface{}
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "failed",
			Data:    map[string]interface{}{"error": "couldn't validate request body"},
		})
		return
	}

	doctorId, ok := requestData["doctorId"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid doctorId format",
		})
		return
	}
	doctorProfile := GetDoctorProfileByStringId(doctorId)
	appointmentDay, ok := requestData["appointmentDay"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid appointmentDay format",
		})
		return
	}
	slotNoFloat, ok := requestData["slotNo"].(float64)
	if !ok {
		c.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid slotNo format",
		})
		return
	}
	slotNo := int32(slotNoFloat)

	slotInterface := ExtractAppoinmentSlotFromDoctorProfile(doctorProfile, appointmentDay, slotNo)
	slot, ok := slotInterface.(primitive.M)
	if !ok || slot == nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error extracting appointment slot",
		})
		return
	}

	if requestData["role"] == "patient" {
		patientId, ok := requestData["patientId"].(string)
		if !ok || patientId != slot["patientid"] {
			c.JSON(http.StatusUnauthorized, responses.UserResponse{
				Status:  http.StatusUnauthorized,
				Message: "failed",
				Data:    map[string]interface{}{"message": "you are not authorized!"},
			})
			return
		}
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{"appointment_details": slot},
	})
}

func GetDoctorProfileByStringId(strDoctorId string) bson.M {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doctorId, err := primitive.ObjectIDFromHex(strDoctorId)
	if err != nil {
		fmt.Println("Error converting doctorId from hex:", err)
		return bson.M{}
	}

	query := bson.D{{Key: "_id", Value: doctorId}}

	var doctorDoc bson.M
	if err := UserCollection.FindOne(ctx, query).Decode(&doctorDoc); err != nil {
		fmt.Println("Error finding doctor:", err)
		return bson.M{}
	}

	return doctorDoc
}

func insertItemToBookedAppointmentsCollection(ba models.BookedAppointment) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := bookedAppointmentsCollection.InsertOne(ctx, ba)
	if err != nil {
		fmt.Println("Error inserting booked appointment:", err)
	} else {
		fmt.Println("Booked appointment inserted, result:", result)
	}
}

func deleteItemFromBookedAppointmentsCollection(doctorId string, slotDay string, slotNo int32) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"$and": []bson.M{
			{"doctorid": doctorId},
			{"day": slotDay},
			{"slotno": slotNo},
		},
	}
	result, err := bookedAppointmentsCollection.DeleteOne(ctx, query)
	if err != nil {
		fmt.Println("Error deleting booked appointment:", err)
	} else {
		fmt.Println("Booked appointment deleted, result:", result)
	}
}

func ViewPatientAppointmentsHistory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	idParam := c.Param("id")
	query := bson.M{"patientid": idParam}
	results, err := bookedAppointmentsCollection.Find(ctx, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error fetching appointment history",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}
	defer results.Close(ctx)

	var patientBookedAppointments []models.BookedAppointment
	for results.Next(ctx) {
		var singleAppointment models.BookedAppointment
		if err := results.Decode(&singleAppointment); err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error decoding appointment data",
				Data:    map[string]interface{}{"error": err.Error()},
			})
			return
		}
		patientBookedAppointments = append(patientBookedAppointments, singleAppointment)
	}

	c.JSON(http.StatusOK, responses.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    map[string]interface{}{"patient_appointments": patientBookedAppointments},
	})
}
