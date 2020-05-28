package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	requestLogger "github.com/gofiber/logger"
	"github.com/lucasmls/backend-cacautime/domain"
	"github.com/lucasmls/backend-cacautime/infra"
	"github.com/lucasmls/backend-cacautime/infra/errors"
)

// ServiceInput ...
type ServiceInput struct {
	Log           infra.LogProvider
	CustomersRepo domain.CustomersRepository
	DutiesRepo    domain.DutiesRepository
	CandiesRepo   domain.CandiesRepository
	SalesRepo     domain.SalesRepository
	Validator     *validator.Validate
}

// Service ...
type Service struct {
	in    ServiceInput
	errCh chan *infra.Error
}

// NewService ...
func NewService(in ServiceInput) (*Service, *infra.Error) {
	const opName infra.OpName = "server.NewService"

	if in.Log == nil {
		err := infra.MissingDependencyError{DependencyName: "Log"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.CustomersRepo == nil {
		err := infra.MissingDependencyError{DependencyName: "CustomersRepo"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	if in.DutiesRepo == nil {
		err := infra.MissingDependencyError{DependencyName: "DutiesRepo"}
		return nil, errors.New(err, opName, infra.KindBadRequest)
	}

	return &Service{
		in:    in,
		errCh: make(chan *infra.Error),
	}, nil
}

// Engine ...
func (s Service) Engine(app *fiber.App) {
	app.Get("/ping", s.pingEndpoint)

	app.Get("/customer", s.listCustomersEndpoint)
	app.Post("/customer", s.registerCustomerEndpoint)
	app.Put("/customer/:id", s.updateCustomerEndpoint)

	app.Get("/duty", s.listDutiesEndpoint)
	app.Post("/duty", s.registerDutyEndpoint)
	app.Put("/duty/:id", s.updateDutyEndpoint)
	app.Get("/duty/:id/sales", s.listDutySales)

	app.Get("/candy", s.listCandiesEndpoint)
	app.Post("/candy", s.registerCandyEndpoint)

	app.Post("/sale", s.registerSaleEndpoint)
}

// Run ...
func (s Service) Run(ctx context.Context) <-chan *infra.Error {
	const opName infra.OpName = "server.Run"

	app := fiber.New()

	app.Use(requestLogger.New())
	app.Use(cors.New())

	s.Engine(app)

	go func() {
		if err := app.Listen(3000); err != nil {
			s.errCh <- errors.New(opName, err, infra.KindUnexpected)
		}

		close(s.errCh)
	}()

	s.in.Log.Info(ctx, opName, "Server up and running...")
	return s.errCh
}
