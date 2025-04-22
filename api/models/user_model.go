package models

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Email     string             `json:"email,omitempty" validate:"required"`
	Password  string             `json:"password,omitempty" validate:"required"`
	Role      string             `json:"role,omitempty"`
	Available bool               `json:"availability,omitempty"`
}

func (u *User) ValidateRole() error {
	validRoles := []string{"doctor", "patient", "admin"}
	roleLower := strings.ToLower(u.Role)
	for _, valid := range validRoles {
		if roleLower == valid {
			return nil
		}
	}
	return errors.New("invalid role: must be 'doctor', 'patient', or 'admin'")
}
