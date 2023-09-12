package routes

import (
	"simple-store-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func PurchaseRoute(app *fiber.App) {
	app.Post("/purcharse", controllers.NewPurchase)
	app.Get("/purcharse", controllers.GetAllPurchases)
	app.Get("/purcharse/:purchaseId/currency/:currency", controllers.GetPurchase)
}
