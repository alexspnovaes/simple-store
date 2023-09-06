package routes

import (
	"simple-store-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func PurchaseRoute(app *fiber.App) {
	app.Post("/user", controllers.NewPurchase) //add thi
}
