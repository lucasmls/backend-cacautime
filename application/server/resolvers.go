package server

import "github.com/gofiber/fiber"

func (s Service) pingEndpoint(c *fiber.Ctx) {
	c.Send("pong")
}
