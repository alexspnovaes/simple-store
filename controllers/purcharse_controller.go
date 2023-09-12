package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"simple-store-api/configs"
	"simple-store-api/models"
	"simple-store-api/responses"
	"strconv"
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

func GetPurchase(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	purchaseId := c.Params("purchaseId")
	currency := c.Params("currency")

	var purcharse models.Purchase
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(purchaseId)

	err := purchaseCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&purcharse)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	purchaseDate := purcharse.Date.Format("2006-01-02")
	dateBefore := purcharse.Date.AddDate(0, -6, 0).Format("2006-01-02")

	baseUrl := "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange"
	filterDate := fmt.Sprintf("?filter=record_date:gte:%s,record_date:lte:%s", dateBefore, purchaseDate)
	endpoint := fmt.Sprintf("%s%s,currency:eq:%s&sort=-record_date&page[number]=1&page[size]=1", baseUrl, filterDate, currency)

	response, err := http.Get(endpoint)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	var currencyResponse responses.CurrencyResponse
	json.Unmarshal(responseData, &currencyResponse)

	if len(currencyResponse.Data) == 0 {
		return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "stating the purchase cannot be converted to the target currency."})
	}

	exchangeRate, err := strconv.ParseFloat(currencyResponse.Data[0].ExchangeRate, 64)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.DataResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	purcharse.ExchangeRate = math.Round(exchangeRate*100) / 100
	purcharse.ConvertedAmount = math.Round(exchangeRate*purcharse.Amount*100) / 100

	return c.Status(http.StatusOK).JSON(responses.DataResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": purcharse}})
}
