package server

import "github.com/gofiber/fiber"

// ServiceInput ...
type ServiceInput struct{}

// Service ...
type Service struct {
	in ServiceInput
}

// NewService ...
func NewService(in ServiceInput) *Service {
	return &Service{
		in: in,
	}
}

// Engine ...
func (s Service) Engine(app *fiber.App) {
	app.Get("/ping", s.pingEndpoint)
}

// Run ...
func (s Service) Run() {
	app := fiber.New()

	s.Engine(app)

	app.Listen(3000)
}
