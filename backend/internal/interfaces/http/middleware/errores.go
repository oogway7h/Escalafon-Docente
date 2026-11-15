package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
)

// ManejadorErrores es el manejador de errores central de Fiber (sección 9, Fase 4): traduce
// cualquier error devuelto por un handler al formato uniforme de la sección 7.
func ManejadorErrores(c *fiber.Ctx, err error) error {
	apiErr := apierror.Mapear(err)
	if apiErr.Status >= 500 {
		slog.Error("error interno", "ruta", c.Path(), "metodo", c.Method(), "requestId", c.GetRespHeader("X-Request-Id"), "error", err)
	}
	return c.Status(apiErr.Status).JSON(dto.ErrorRespuesta{Codigo: apiErr.Codigo, Mensaje: apiErr.Mensaje, Detalles: apiErr.Detalles})
}
