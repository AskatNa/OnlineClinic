package models

type User struct {
	Name      string `json:"name,omitempty" validate:"required"`
	Email     string `json:"email,omitempty" validate:"required"`
	Password  string `json:"password,omitempty" validate:"required"`
	Role      string `json:"role,omitempty" validate:"required"`
	Available bool   `json:"availability,omitempty"`
}
