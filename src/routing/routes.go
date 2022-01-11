package routing

import (
	"github.com/gofiber/fiber/v2"
	"loundry/api/src/http"
	"loundry/api/src/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	apiV1 := api.Group("/v1")
	apiV1.Post("/login", http.AuthCheck)

	SecuredV1 := apiV1.Use(middleware.ProtectedArea())
	//SecuredV1.Use(middleware.RestrictedArea())
	CustomerRoute := SecuredV1.Group("/customer")
	CustomerRoute.Get("/", http.GetCustomer)
	CustomerRoute.Get("/:id", http.CustomerDetail)
	CustomerRoute.Post("/", http.CreateCustomer)
	CustomerRoute.Post("search", http.SearchCustomer)

	//item route
	ItemRoute := SecuredV1.Group("/item")
	ItemRoute.Post("/", http.CreateNewItems)
	ItemRoute.Get("/", http.GetAllItems)
	ItemRoute.Get("/:id", http.GetItemById)
	ItemRoute.Put("/:id", http.UpdateItemById)
	ItemRoute.Delete("/:id", http.DeleteItemById)

	//transaction route
	TransactionRoute := SecuredV1.Group("/transaction")
	TransactionRoute.Post("/", http.CreateTransaction)
	TransactionRoute.Get("/", http.GetAllTransaction)
	TransactionRoute.Get("/:id", http.GetTransactionById)

	//transaction items route
	TransactionItemRoute := SecuredV1.Group("/transaction_item")
	TransactionItemRoute.Get("/:id", http.GetTransaksiitemsbyIdTransaction)
	TransactionItemRoute.Post("/:id", http.AddItemToTransaction)
	TransactionItemRoute.Put("/:id", http.UpdateItemsTransaction)
	TransactionItemRoute.Delete("/:id", http.DeleteTransactionItems)

}
