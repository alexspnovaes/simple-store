package controllers

import (
	"context"
	"net/http"
	"simple-store-api/configs"
	"simple-store-api/models"
	"simple-store-api/responses"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var purchaseCollection *mongo.Collection = configs.GetCollection(configs.DB, "purchase")
var validate = validator.New()

func NewPurchase(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var purchase models.Purchase
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&purchase); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&purchase); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.DataResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newPurchase := models.Purchase{
		Id:          primitive.NewObjectID(),
		Description: purchase.Description,
		Date:        purchase.Date,
		Amount:      purchase.Amount,
	}

	result, err := purchaseCollection.InsertOne(ctx, newPurchase)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.DataResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetAllPurchases(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var purchases []models.Purchase
	defer cancel()

	results, err := purchaseCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singlePurchase models.Purchase
		if err = results.Decode(&singlePurchase); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		purchases = append(purchases, singlePurchase)
	}

	return c.Status(http.StatusOK).JSON(
		responses.DataResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": purchases}},
	)
}
