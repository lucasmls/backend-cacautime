package server

import (
	"context"
	"fmt"
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

func (s Service) registerCustomerEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerCustomerEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	payload := registerCustomerPayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

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

	customer, err := s.in.CustomersRepo.Register(ctx, customerDTO)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": customerDTO,
		})

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(customer)
}

func (s Service) listCustomersEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listCustomersEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customers, err := s.in.CustomersRepo.List(ctx)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName)

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(customers)
}

func (s Service) registerDutyEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerDutyEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	payload := registerDutyPayload{}
	if err := c.BodyParser(&payload); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": payload,
		})

		// @TODO => Retornar o erro de dominio...
		fmt.Println(err)
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

	dutyDTO := domain.Duty{
		Date:          payload.Date,
		CandyQuantity: payload.CandyQuantity,
	}

	duty, err := s.in.DutiesRepo.Register(ctx, dutyDTO)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": dutyDTO,
		})

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(duty)
}

func (s Service) listDutiesEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listDutiesEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	duties, err := s.in.DutiesRepo.List(ctx)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName)

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(duties)
}

func (s Service) listDutySales(c *fiber.Ctx) {
	const opName infra.OpName = "server.listDutySales"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	dutyIDParam := c.Params("id")
	dutyID, err := strconv.Atoi(dutyIDParam)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": dutyIDParam,
		})

		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid duty id.",
		})

		return
	}

	dutiesSales, sErr := s.in.DutiesRepo.Sales(ctx, infra.ObjectID(dutyID))
	if sErr != nil && errors.Kind(sErr) == infra.KindNotFound {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"param": dutyIDParam,
		})

		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified duty was not found",
		})

		return
	}

	if sErr != nil {
		s.errCh <- errors.New(ctx, sErr, opName, infra.Metadata{
			"param": dutyIDParam,
		})

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(dutiesSales)
}

func (s Service) registerCandyEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerCandyEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	candyDto := domain.Candy{}
	if err := c.BodyParser(&candyDto); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": candyDto,
		})

		// @TODO => Retornar o erro de dominio...
		return
	}

	candy, err := s.in.CandiesRepo.Register(ctx, candyDto)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": candyDto,
		})

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(candy)
}

func (s Service) listCandiesEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listCandiesEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	candies, err := s.in.CandiesRepo.List(ctx)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName)

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(candies)
}

func (s Service) registerSaleEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerSaleEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	saleDTO := domain.Sale{}
	if err := c.BodyParser(&saleDTO); err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": saleDTO,
		})

		// @TODO => Retornar o erro de dominio...
		return
	}

	sale, err := s.in.SalesRepo.Register(ctx, saleDTO)
	if err != nil {
		s.errCh <- errors.New(ctx, err, opName, infra.Metadata{
			"payload": saleDTO,
		})

		// @TODO => Retornar o erro de dominio...
		return
	}

	c.Status(200).JSON(sale)
}
