package http

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"loundry/api/src/database"
	"loundry/api/src/helper"
	"loundry/api/src/models"
)

type Customer struct {
	ID          uint   `json:"id"`
	Nama        string `json:"nama"`
	PhoneNumber string `json:"phone_number"`
}

func customerResponse(customer models.Customer) Customer {
	return Customer{
		ID:          customer.ID,
		Nama:        customer.Nama,
		PhoneNumber: customer.PhoneNumber,
	}
}

func CreateCustomer(ctx *fiber.Ctx) error {
	type CustomerPost struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	}
	userId := helper.GetUserIdFromToken(ctx)
	customer := []models.Customer{}
	dataPost := new(CustomerPost)
	if err := ctx.BodyParser(dataPost); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var count int64
	database.DBConn.Model(customer).Where("user_id=? AND phone_number=?", userId, dataPost.PhoneNumber).Count(&count)

	if count > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Customer already exist",
		})
	}

	newCustomer := models.Customer{
		Nama:        dataPost.Name,
		PhoneNumber: dataPost.PhoneNumber,
	}
	newCustomer.UserID = userId

	database.DBConn.Create(&newCustomer)
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Customer created",
		"data":    customerResponse(newCustomer),
	})

}

func SearchCustomer(ctx *fiber.Ctx) error {
	type SearchCust struct {
		Search string `json:"search"`
	}

	userId := helper.GetUserIdFromToken(ctx)
	customer := []models.Customer{}
	postData := new(SearchCust)
	if err := ctx.BodyParser(postData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := database.DBConn.Where("user_id=? AND nama LIKE ?", userId, "%"+postData.Search+"%").Find(&customer).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	searCustomers := []Customer{}
	for _, cust := range customer {
		searCustomers = append(searCustomers, customerResponse(cust))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Customer found",
		"data":    searCustomers,
	})
}

func GetCustomer(ctx *fiber.Ctx) error {
	userId := helper.GetUserIdFromToken(ctx)
	customer := []models.Customer{}
	if err := database.DBConn.Where("user_id=?", userId).Find(&customer).Error; err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	responseCustomers := []Customer{}
	for _, value := range customer {
		responseCustomers = append(responseCustomers, customerResponse(value))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Customer found",
		"data":    responseCustomers,
	})
}

func findCustomer(id int, userId uint, customer *models.Customer) error {
	database.DBConn.Find(&customer, "id=? AND user_id=?", id, userId)
	if customer.ID == 0 {
		return errors.New("Customer not found")
	}
	return nil
}

func CustomerDetail(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")

	var customer models.Customer

	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userId := helper.GetUserIdFromToken(ctx)
	if err := findCustomer(id, userId, &customer); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Customer found",
		"data":    customerResponse(customer),
	})
}
