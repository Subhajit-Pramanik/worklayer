package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/vyolayer/vyolayer/internal/gateway/middleware"
	m "github.com/vyolayer/vyolayer/internal/gateway/middleware"
	gm "github.com/vyolayer/vyolayer/internal/shared/middleware"
)

// Server represents the API Gateway instance
type Server struct {
	app  *fiber.App
	port string
}

// Router interface for registering routes (ISP)
type RouteRegistrar interface {
	RegisterRoutes(router fiber.Router)
}

// New creates and configures a new Server instance
func New(port string) *Server {
	app := fiber.New(fiber.Config{AppName: "VyoLayer Gateway"})

	app.Use(m.RequestContext())
	app.Use(m.ErrorHandler())
	app.Use(m.RequestLogger())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
	}))

	// Inject custom middleware to propagate headers to gRPC requests
	app.Use(m.GRPCMetadataMiddleware())

	app.Get("/swagger/*", swagger.HandlerDefault)

	return &Server{
		app:  app,
		port: port,
	}
}

// RegisterRegistrars allows appending groups of routes (OCP)
func (s *Server) RegisterRegistrars(registrars ...RouteRegistrar) {
	v1 := s.app.Group("/v1")
	rateLimiter := middleware.NewRateLimiter(100, 5*time.Minute, "global")
	v1.Use(rateLimiter.Handler())

	for _, registrar := range registrars {
		registrar.RegisterRoutes(v1)
	}

	v1.Use(gm.NotFoundMiddleware)
}

// Start runs the HTTP server with graceful shutdown handling
func (s *Server) Start() error {
	// Graceful shutdown handling
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		log.Println("Shutting down HTTP server...")
		s.app.Shutdown()
	}()

	log.Printf("API Gateway listening on :%s", s.port)
	return s.app.Listen(fmt.Sprintf(":%s", s.port))
}
