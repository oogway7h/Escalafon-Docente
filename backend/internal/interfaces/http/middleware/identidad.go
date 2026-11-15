package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

const claveIdentidad = "identidad"

// Identidad es lo que Autenticar deja en c.Locals tras validar la sesión (sección 5.2).
type Identidad struct {
	UsuarioID string
	RolID     int16
	Permisos  map[string]bool
	JTI       string
	Expira    time.Time // para calcular cuánto falta al cerrar sesión
}

func (i Identidad) TienePermiso(codigo string) bool { return i.Permisos[codigo] }

func GuardarIdentidad(c *fiber.Ctx, id Identidad) {
	c.Locals(claveIdentidad, id)
}

// DeContexto lee la identidad guardada por Autenticar. El segundo valor es false si no hay
// ninguna (la ruta no pasó por el middleware).
func DeContexto(c *fiber.Ctx) (Identidad, bool) {
	id, ok := c.Locals(claveIdentidad).(Identidad)
	return id, ok
}
