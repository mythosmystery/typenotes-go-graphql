package model

import (
	"context"

	"github.com/mythosmystery/typenotes-go-graphql/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserById(userId string, db *database.DB) (*User, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	var user User
	res := db.User.FindOne(context.Background(), bson.M{"_id": id})
	e := res.Decode(&user)
	if e != nil {
		return nil, err
	}
	return &user, nil
}
