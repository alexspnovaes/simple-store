package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Description string             `json:"description"`
	Date        time.Time          `json:"date"`
	Amount      int                `json:"amount"`
}
