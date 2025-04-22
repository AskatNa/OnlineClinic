package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Doctor struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Email      string             `bson:"email" json:"email"`
	Specialty  string             `bson:"specialty" json:"specialty"`
	Experience int                `bson:"experience" json:"experience"`
	Bio        string             `bson:"bio" json:"bio"`
}
