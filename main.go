package main

import (
	"simple-store-api/configs"
	"simple-store-api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	configs.ConnectDB()
	routes.PurchaseRoute(app) //add this

	app.Listen(":8001")
}
