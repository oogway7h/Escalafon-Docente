package middleware

import (
	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/interfaces/http/apierror"
)

// RequierePermiso comprueba el permiso en la identidad (viene del JWT, sección 5.2): no accede
// a Redis ni a la base de datos.
func RequierePermiso(codigo string) fiber.Handler {
	return RequiereAlguno(codigo)
}

// RequiereAlguno exige al menos uno de los permisos dados (p. ej. para lecturas compartidas).
func RequiereAlguno(codigos ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, ok := DeContexto(c)
		if ok {
			for _, codigo := range codigos {
				if id.TienePermiso(codigo) {
					return c.Next()
				}
			}
		}
		return apierror.Nuevo(403, "PERMISO_DENEGADO", "No tienes permiso para esta acción")
	}
}
