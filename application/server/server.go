package server

import (
	"github.com/gofiber/cors"
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
	SalesRepo     domain.SalesRepository
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
	app.Get("/customer", s.listCustomersEndpoint)

	app.Get("/duty", s.listDutiesEndpoint)
	app.Post("/duty", s.registerDutyEndpoint)
	app.Get("/duty/sales", s.listDutiesSales)

	app.Get("/candy", s.listCandiesEndpoint)
	app.Post("/candy", s.registerCandyEndpoint)

	app.Post("/sale", s.registerSaleEndpoint)
}

// Run ...
func (s Service) Run() {
	app := fiber.New()

	app.Use(requestLogger.New())
	app.Use(cors.New())

	s.Engine(app)

	app.Listen(3000)
}
