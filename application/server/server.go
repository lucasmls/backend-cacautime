package server

import (
	"github.com/gofiber/fiber"
	requestLogger "github.com/gofiber/logger"
	"github.com/lucasmls/backend-cacautime/domain"
)

// ServiceInput ...
type ServiceInput struct {
	CustomersRepo domain.CustomersRepository
}

// Service ...
type Service struct {
	in ServiceInput
}

// NewService ...
func NewService(in ServiceInput) *Service {
	if in.CustomersRepo == nil {
		panic("Customers Repo is required.")
	}

	return &Service{
		in: in,
	}
}

// Engine ...
func (s Service) Engine(app *fiber.App) {
	app.Get("/ping", s.pingEndpoint)

	app.Post("/customer", s.registerCustomerEndpoint)
}

// Run ...
func (s Service) Run() {
	app := fiber.New()

	app.Use(requestLogger.New())

	s.Engine(app)

	app.Listen(3000)
}
