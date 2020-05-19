package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber"
	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
)

func (s Service) pingEndpoint(c *fiber.Ctx) {
	c.Send("pong")
}

func (s Service) registerCustomerEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerCustomerEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customer := domain.Customer{}
	if err := c.BodyParser(&customer); err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	err := s.in.CustomersRepo.Register(ctx, customer)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(
		map[string]string{
			"message": "Customer registered successfully.",
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

	dutiesSales, err := s.in.DutiesRepo.Consolidate(ctx)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(dutiesSales)
}

func (s Service) registerCandyEndpoint(c *fiber.Ctx) {
	const opName infra.OpName = "server.registerCandyEndpoint"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	candy := domain.Candy{}
	if err := c.BodyParser(&candy); err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	err := s.in.CandiesRepo.Register(ctx, candy)
	if err != nil {
		// @TODO => Criar o canal de error e inserir o erro lá...
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(
		map[string]string{
			"message": "Candy registered successfully.",
		},
	)
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
