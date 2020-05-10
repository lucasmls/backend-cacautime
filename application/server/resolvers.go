package server

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber"
	"github.com/lucasmls/backend-cacautime/domain"
)

func (s Service) pingEndpoint(c *fiber.Ctx) {
	c.Send("pong")
}

func (s Service) registerCustomerEndpoint(c *fiber.Ctx) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	customer := domain.Customer{}
	if err := c.BodyParser(&customer); err != nil {
		fmt.Println(err)
		return
	}

	err := s.in.CustomersRepo.Register(ctx, customer)
	if err != nil {
		fmt.Println(err)
		return
	}

	c.Status(200).JSON(
		map[string]string{
			"message": "Customer registered successfully.",
		},
	)
}
