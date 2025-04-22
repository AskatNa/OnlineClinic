package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Patient struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name"`
	Email string             `bson:"email" json:"email"`
	Age   int                `bson:"age" json:"age"`     // Default: 0
	Notes string             `bson:"notes" json:"notes"` // Default: ""
}
