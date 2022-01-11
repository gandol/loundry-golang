package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/gofiber/helmet/v2"
	"log"
	"loundry/api/src/database"
	"loundry/api/src/routing"
	"os"
	"time"
)

func main() {
	database.ConnectDb()
	app := fiber.New()
	routing.SetupRoutes(app)
	app.Use(helmet.New())
	csrfConfig := csrf.Config{
		KeyLookup:      "header:X-Csrf-Token", // string in the form of '<source>:<key>' that is used to extract token from the request
		CookieName:     "my_csrf_",            // name of the session cookie
		CookieSameSite: "Strict",              // indicates if CSRF cookie is requested by SameSite
		Expiration:     3 * time.Hour,         // expiration is the duration before CSRF token will expire
		KeyGenerator:   utils.UUID,            // creates a new CSRF token
	}
	app.Use(csrf.New(csrfConfig))
	file, err := os.OpenFile("./access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	// Set config for logger
	loggerConfig := logger.Config{
		Output: file, // add file to save output
	}

	// Use middlewares for each route
	app.Use(
		logger.New(loggerConfig), // add Logger middleware with config
	)
	app.Listen(":8000")
}
