package health

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/pkg/response"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(router fiber.Router) {
	router.Get("/health", h.health)
	log.Println("[HEALTH] routes registered")
}

// @Summary Check API Health
// @Description Returns the health status of the API Gateway
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=object{status=string,version=string}}
// @Router /health [get]
func (h *HealthHandler) health(c *fiber.Ctx) error {
	res := fiber.Map{
		"status":  "ok",
		"version": "1.0.0",
	}
	return response.SuccessWithMessage(c, fiber.StatusOK, "Health check successful", res)
}
