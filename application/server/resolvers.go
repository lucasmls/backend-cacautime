package server

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

func handleValidationError(payload interface{}, err error) map[string]string {
	errorsMap := make(map[string]string)

	reflected := reflect.ValueOf(payload)

	for _, e := range err.(validator.ValidationErrors) {
		field, _ := reflected.Type().FieldByName(e.StructField())

		var key string
		if key = field.Tag.Get("json"); key == "" {
			key = strings.ToLower(e.StructField())
		}

		switch e.Tag() {
		case "required":
			errorsMap[key] = "The " + key + " is required."
		case "max":
			errorsMap[key] = "The " + key + " is bigger than the maximum expected value."
		case "min":
			errorsMap[key] = "The " + key + " is smaller than the minimum expected value."
		default:
			errorsMap[key] = "The " + key + " is invalid."
		}
	}

	return errorsMap
}

func (s Service) pingEndpoint(c *fiber.Ctx) {
	c.Send("pong")
}

func (s Service) login(c *fiber.Ctx) {
	const opName infra.OpName = "server.login"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	payload := loginPayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(422).JSON(
			map[string]string{
				"message": "Invalid payload.",
			},
		)

		return
	}

	if err := s.in.Validator.Struct(payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		response := handleValidationError(payload, err)

		c.Status(422).JSON(response)

		return
	}

	token, err := s.in.AuthRepo.Login(ctx, payload.Email, payload.Password)
	if err != nil && errors.Kind(err) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"email": payload.Email,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "User not found",
		})

		return
	}

	if err != nil && errors.Kind(err) == infra.KindUnauthorized {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"email": payload.Email,
		})

		c.Status(401).JSON(map[string]interface{}{
			"message": "Wrong e-mail or password",
		})

		return
	}

	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"email": payload.Email,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(
		map[string]string{
			"token": token,
		},
	)
}

func (s Service) registerCustomerEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerCustomerEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	payload := customerPayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(422).JSON(
			map[string]string{
				"message": "Invalid payload.",
			},
		)

		return
	}

	if err := s.in.Validator.Struct(payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		response := handleValidationError(payload, err)

		c.Status(422).JSON(response)

		return
	}

	customerDTO := domain.Customer{
		Name:  payload.Name,
		Phone: payload.Phone,
	}

	customer, cErr := s.in.CustomersRepo.Register(ctx, customerDTO)
	if cErr != nil {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": customerDTO,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(customer)
}

func (s Service) updateCustomerEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.updateCustomerEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customerIDParam := c.Params("id")
	customerID, err := strconv.Atoi(customerIDParam)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": customerIDParam,
		})

		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid customer id.",
		})

		return
	}

	payload := customerPayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(422).JSON(
			map[string]string{
				"message": "Invalid payload.",
			},
		)

		return
	}

	if err := s.in.Validator.Struct(payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		response := handleValidationError(payload, err)

		c.Status(422).JSON(response)

		return
	}

	customerDTO := domain.Customer{
		Name:  payload.Name,
		Phone: payload.Phone,
	}

	customer, cErr := s.in.CustomersRepo.Update(ctx, infra.ObjectID(customerID), customerDTO)
	if cErr != nil && errors.Kind(cErr) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified customer was not found",
		})

		return
	}

	if cErr != nil {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": customerDTO,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(customer)
}

func (s Service) deleteCustomerEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.deleteCustomerEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customerIDParam := c.Params("id")
	customerID, err := strconv.Atoi(customerIDParam)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": customerIDParam,
		})

		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid customer id.",
		})

		return
	}

	cErr := s.in.CustomersRepo.Delete(ctx, infra.ObjectID(customerID))
	if cErr != nil && errors.Kind(cErr) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"param": customerIDParam,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified customer was not found",
		})

		return
	}

	if cErr != nil {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"param": customerIDParam,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(map[string]string{"Message": "Customer deleted successfully!"})
}

func (s Service) listCustomersEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listCustomersEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customers, err := s.in.CustomersRepo.List(ctx)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName)

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(customers)
}

func (s Service) registerCandyEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerCandyEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	payload := candyPayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(422).JSON(
			map[string]string{
				"message": "Invalid payload.",
			},
		)

		return
	}

	if err := s.in.Validator.Struct(payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		response := handleValidationError(payload, err)

		c.Status(422).JSON(response)

		return
	}

	candyDto := domain.Candy{
		Name:  payload.Name,
		Price: payload.Price,
	}

	candy, cErr := s.in.CandiesRepo.Register(ctx, candyDto)
	if cErr != nil {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": candyDto,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(candy)
}

func (s Service) updateCandyEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.updateCandyEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	candyIDParam := c.Params("id")
	candyID, err := strconv.Atoi(candyIDParam)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": candyIDParam,
		})

		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid candy id.",
		})

		return
	}

	payload := candyPayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(422).JSON(
			map[string]string{
				"message": "Invalid payload.",
			},
		)

		return
	}

	if err := s.in.Validator.Struct(payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		response := handleValidationError(payload, err)

		c.Status(422).JSON(response)

		return
	}

	candyDTO := domain.Candy{
		Name:  payload.Name,
		Price: payload.Price,
	}

	candy, cErr := s.in.CandiesRepo.Update(ctx, infra.ObjectID(candyID), candyDTO)
	if cErr != nil && errors.Kind(cErr) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified candy was not found",
		})

		return
	}

	if cErr != nil {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": candyDTO,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(candy)
}

func (s Service) deleteCandyEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.deleteCandyEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	candyIDParam := c.Params("id")
	candyID, err := strconv.Atoi(candyIDParam)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": candyIDParam,
		})

		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid candy id.",
		})

		return
	}

	cErr := s.in.CandiesRepo.Delete(ctx, infra.ObjectID(candyID))
	if cErr != nil && errors.Kind(cErr) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"param": candyIDParam,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified candy was not found",
		})

		return
	}

	if cErr != nil {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"param": candyIDParam,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(map[string]string{"Message": "Candy deleted successfully!"})
}

func (s Service) listCandiesEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listCandiesEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	candies, err := s.in.CandiesRepo.List(ctx)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName)

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(candies)
}

func (s Service) registerSaleEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerSaleEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	payload := salePayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(422).JSON(
			map[string]string{
				"message": "Invalid payload.",
			},
		)

		return
	}

	if err := s.in.Validator.Struct(payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		response := handleValidationError(payload, err)

		c.Status(422).JSON(response)

		return
	}

	saleDTO := domain.Sale{
		CustomerID:    infra.ObjectID(payload.CustomerID),
		CandyID:       infra.ObjectID(payload.CandyID),
		Status:        domain.Status(payload.Status),
		PaymentMethod: domain.PaymentMethod(payload.PaymentMethod),
	}

	sale, sErr := s.in.SalesRepo.Register(ctx, saleDTO)
	if sErr != nil {
		s.errCh <- errors.New(ctx, sErr, opName, infra.Metadata{
			"payload": saleDTO,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(sale)
}

func (s Service) updateSaleEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.updateSaleEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	saleIDParam := c.Params("id")
	saleID, err := strconv.Atoi(saleIDParam)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": saleIDParam,
		})

		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid sale id.",
		})

		return
	}

	payload := updateSalePayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(422).JSON(
			map[string]string{
				"message": "Invalid payload.",
			},
		)

		return
	}

	if err := s.in.Validator.Struct(payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		response := handleValidationError(payload, err)

		c.Status(422).JSON(response)

		return
	}

	saleDTO := domain.Sale{
		PaymentMethod: domain.PaymentMethod(payload.PaymentMethod),
		Status:        domain.Status(payload.Status),
	}

	sale, cErr := s.in.SalesRepo.Update(ctx, infra.ObjectID(saleID), saleDTO)
	if cErr != nil && errors.Kind(cErr) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": payload,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified sale was not found",
		})

		return
	}

	if cErr != nil {
		s.errCh <- errors.New(ctx, cErr, opName, infra.Metadata{
			"payload": saleDTO,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(sale)
}

func (s Service) deleteSaleEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.deleteSaleEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	saleIDParam := c.Params("id")
	saleID, err := strconv.Atoi(saleIDParam)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": saleIDParam,
		})

		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid sale id.",
		})

		return
	}

	dErr := s.in.SalesRepo.Delete(ctx, infra.ObjectID(saleID))
	if dErr != nil && errors.Kind(dErr) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, dErr, opName, infra.Metadata{
			"param": saleIDParam,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified sale was not found",
		})

		return
	}

	if dErr != nil {
		s.errCh <- errors.New(ctx, dErr, opName, infra.Metadata{
			"param": saleIDParam,
		})

		c.Status(500).JSON(
			map[string]string{
				"message": "Internal server error.",
			},
		)

		return
	}

	c.Status(200).JSON(map[string]string{"Message": "Sale deleted successfully!"})
}
