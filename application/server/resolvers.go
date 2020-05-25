package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber"
	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

func (s Service) pingEndpoint(c *fiber.Ctx) {
	c.Send("pong")
}

func (s Service) registerCustomerEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerCustomerEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customerDto := domain.Customer{}
	if err := c.BodyParser(&customerDto); err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	customer, err := s.in.CustomersRepo.Register(ctx, customerDto)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(
		map[string]interface{}{
			"id":    customer.ID,
			"name":  customer.Name,
			"phone": customer.Phone,
		},
	)
}

func (s Service) listCustomersEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listCustomersEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customers, err := s.in.CustomersRepo.List(ctx)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(customers)
}

func (s Service) registerDutyEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerDutyEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	duty := domain.Duty{}
	if err := c.BodyParser(&duty); err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	err := s.in.DutiesRepo.Register(ctx, duty)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(
		map[string]string{
			"message": "Duty registered successfully.",
		},
	)
}

func (s Service) listDutiesEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listDutiesEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	duties, err := s.in.DutiesRepo.List(ctx)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(duties)
}

func (s Service) listDutiesSales(c *fiber.Ctx) {
	const opName infra.OpName = "server.listDutiesSales"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	dutyIDParam := c.Params("id")
	dutyID, err := strconv.Atoi(dutyIDParam)
	if err != nil {
		c.Status(422).JSON(map[string]interface{}{
			"message": "Invalid duty id.",
		})

		return
	}

	dutiesSales, sErr := s.in.DutiesRepo.Sales(ctx, infra.ObjectID(dutyID))
	if sErr != nil && errors.Kind(sErr) == infra.KindNotFound {
		c.Status(404).JSON(map[string]interface{}{
			"message": "The specified duty was not found",
		})

		return
	}

	if sErr != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(sErr)
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
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	candy, err := s.in.CandiesRepo.Register(ctx, candyDto)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(
		map[string]interface{}{
			"id":    candy.ID,
			"name":  candy.Name,
			"price": candy.Price,
		},
	)
}

func (s Service) listCandiesEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.listCandiesEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	candies, err := s.in.CandiesRepo.List(ctx)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(candies)
}

func (s Service) registerSaleEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerSaleEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	sale := domain.Sale{}
	if err := c.BodyParser(&sale); err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	err := s.in.SalesRepo.Register(ctx, sale)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(
		map[string]string{
			"message": "Sale registered successfully.",
		},
	)
}
