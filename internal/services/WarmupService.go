package services

import "github.com/gofiber/fiber/v2"

type WarmupService struct{}

func (w WarmupService) GetResponse(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "healthy"})
}
