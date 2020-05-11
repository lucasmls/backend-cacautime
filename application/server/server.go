package server

import (
	"github.com/gofiber/fiber"
	requestLogger "github.com/gofiber/logger"
	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
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
func NewService(in ServiceInput) (*Service, *infra.Error) {
	const opName infra.OpName = "server.NewService"

	if in.CustomersRepo == nil {
		// @TODO => Create the missing dependency error
		return nil, errors.New(opName, "Customers Repo is required.", infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
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
