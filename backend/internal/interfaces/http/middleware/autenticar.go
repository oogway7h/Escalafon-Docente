package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/usuario"
	"escalafon/backend/internal/interfaces/http/apierror"
)

// Autenticar implementa el orden exacto de la sección 5.2: verificar el JWT localmente (sin
// red), una consulta a la caché de sesión, y solo si hay fallo de caché, una consulta a la BD.
func Autenticar(tokens ports.TokenService, cache ports.SessionCache, usuarios usuario.Repository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return apierror.Nuevo(401, "NO_AUTENTICADO", "Falta el encabezado de autorización")
		}
		claims, err := tokens.Verificar(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			return apierror.Nuevo(401, "TOKEN_INVALIDO", "La sesión no es válida")
		}

		ctx := c.Context()
		revocado, err := cache.EstaRevocado(ctx, claims.JTI)
		if err != nil { // falla cerrado: sin la lista de revocación no se puede saber si la sesión fue cerrada
			return apierror.Nuevo(503, "SESION_NO_VERIFICABLE", "No se pudo verificar la sesión, intenta de nuevo en unos segundos")
		}
		if revocado {
			return apierror.Nuevo(401, "SESION_REVOCADA", "La sesión fue cerrada")
		}

		version, activo, encontrado, _ := cache.ObtenerSesion(ctx, claims.UsuarioID)
		if !encontrado {
			u, err := usuarios.BuscarPorID(ctx, claims.UsuarioID)
			if err != nil {
				return apierror.Nuevo(401, "TOKEN_INVALIDO", "La sesión no es válida")
			}
			version, activo = u.VersionSesion, u.Activo
			_ = cache.GuardarSesion(ctx, claims.UsuarioID, version, activo)
		}
		if claims.Version != version || !activo {
			return apierror.Nuevo(401, "TOKEN_INVALIDO", "La sesión ya no es válida")
		}

		permisos := map[string]bool{}
		for _, p := range claims.Permisos {
			permisos[p] = true
		}
		GuardarIdentidad(c, Identidad{UsuarioID: claims.UsuarioID, RolID: claims.RolID, Permisos: permisos, JTI: claims.JTI, Expira: claims.Expira})
		return c.Next()
	}
}
