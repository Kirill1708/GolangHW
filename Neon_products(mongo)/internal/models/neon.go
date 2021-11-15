package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Neon struct {
	ID     primitive.ObjectID `bson:"_id"`
	Title  string             `bson:"title"`
	Sizes  string             `bson:"sizes"`
	Colors string             `bson:"colors"`
	Cost   int                `bson:"cost"`
	Theme  string             `bson:"theme"`
}
