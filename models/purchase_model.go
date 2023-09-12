package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	Id              primitive.ObjectID `json:"id,omitempty"`
	Description     string             `json:"description"`
	Date            time.Time          `json:"date"`
	Amount          float64            `json:"amount"`
	ExchangeRate    float64            `json:"exchangeRate"`
	ConvertedAmount float64            `json:"convertedAmount"`
}
