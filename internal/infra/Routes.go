package infra

import (
	"AuthAPI/internal/services"

	"github.com/gofiber/fiber/v2"
)

func Routes(
	app *fiber.App,
	warmUpService services.WarmupService,
	googleAuthService services.GoogleAuthService,
) {
	/* warm up routes */
	app.Get("/health", warmUpService.GetResponse)

	/* auth routes */
	app.Post("/auth/google", googleAuthService.GetAuthUrl)
	app.Get("/auth/google/callback", googleAuthService.HandleCallBack)
}
