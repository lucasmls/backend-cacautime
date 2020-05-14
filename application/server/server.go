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
	DutiesRepo    domain.DutiesRepository
	CandiesRepo   domain.CandiesRepository
}

// Service ...
type Service struct {
	in ServiceInput
}

// NewService ...
func NewService(in ServiceInput) (*Service, *infra.Error) {
	const opName infra.OpName = "server.NewService"

	if in.CustomersRepo == nil {
		err := infra.MissingDependencyError{DependencyName: "CustomersRepo"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.DutiesRepo == nil {
		err := infra.MissingDependencyError{DependencyName: "DutiesRepo"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in: in,
	}, nil
}

// Engine ...
func (s Service) Engine(app *fiber.App) {
	app.Get("/ping", s.pingEndpoint)

	app.Post("/customer", s.registerCustomerEndpoint)

	app.Get("/duty", s.listDutiesEndpoint)
	app.Post("/duty", s.registerDutyEndpoint)
	app.Get("/duty/sales", s.listDutiesSales)

	app.Post("/candy", s.registerCandyEndpoint)
}

// Run ...
func (s Service) Run() {
	app := fiber.New()

	app.Use(requestLogger.New())

	s.Engine(app)

	app.Listen(3000)
}
