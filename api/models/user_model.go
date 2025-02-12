package models

import (
	"errors"
	"strings"
)

type User struct {
	Name      string `json:"name,omitempty" validate:"required"`
	Email     string `json:"email,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
	Role      string `json:"role,omitempty"`
	Available bool   `json:"availability,omitempty"`
}

func (u *User) ValidateRole() error {
	validRoles := []string{"doctor", "patient"}

	roleLower := strings.ToLower(u.Role)

	for _, valid := range validRoles {
		if roleLower == valid {
			return nil
		}
	}
	return errors.New("invalid role: must be 'doctor' or 'patient'")
}
