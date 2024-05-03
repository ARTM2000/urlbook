package controller

import "github.com/gofiber/fiber/v2"

type HttpController interface {
	InitRoutes(fiber.Router)
	GetPrefix() string
}