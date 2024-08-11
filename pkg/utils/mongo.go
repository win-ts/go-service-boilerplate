package utils

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ConvertToObjectID converts a string to a primitive.ObjectID
func ConvertToObjectID(id string) primitive.ObjectID {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return objectID
}
