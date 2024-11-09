package infra

import (
	"AuthAPI/cfg"
	"AuthAPI/internal/services"

	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"golang.org/x/sync/errgroup"
)

func Run(cfg *cfg.AppConfig) {

	app := fiber.New(fiber.Config{
		Prefork:               true,
		CaseSensitive:         true,
		StrictRouting:         true,
		DisableStartupMessage: true,
	})

	/* setup middleware */
	app.Use(recover.New())
	app.Use(logger.New())

	/* initiate services */
	warmupService := services.WarmupService{}
	googleAuthService := services.NewGoogleAuthService(cfg)

	/* routing */
	Routes(
		app,
		warmupService,
		*googleAuthService,
	)

	/* Start server */
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.Port)); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	shutdownGroup := new(errgroup.Group)

	shutdownGroup.Go(func() error {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

		sig := <-shutdownSignal
		log.Printf("Received signal: %v. Shutting down gracefully...", sig)
		return app.Shutdown()
	})

	if err := shutdownGroup.Wait(); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}

	log.Println("Server shut down gracefully.")
	os.Exit(0)
}
