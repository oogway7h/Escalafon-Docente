// Package handlers traduce HTTP a casos de uso (sección 7): no contienen lógica de negocio ni
// acceden a repositorios directamente.
package handlers

import (
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/auth"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

type AuthHandler struct {
	Autenticar              *auth.AutenticarUsuario
	CerrarSesionCaso        *auth.CerrarSesion
	ObtenerSesionCaso       *auth.ObtenerSesion
	SolicitarRecuperacion   *auth.SolicitarRecuperacion
	RestablecerRecuperacion *auth.RestablecerContrasena // PropositoEsperado = RECUPERACION
	RestablecerActivacion   *auth.RestablecerContrasena // PropositoEsperado = ACTIVACION
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var in dto.LoginEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	salida, err := h.Autenticar.Ejecutar(c.Context(), auth.AutenticarEntrada{
		Correo: in.Correo, Contrasena: in.Contrasena, DireccionIP: c.IP(),
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(dto.LoginSalida{
		Token: salida.Token, ExpiraEn: salida.ExpiraEn, Permisos: salida.Permisos,
		Usuario: dto.UsuarioResumen{
			ID: salida.Usuario.ID, Nombres: salida.Usuario.Nombres, Apellidos: salida.Usuario.Apellidos,
			Correo: salida.Usuario.Correo, RolID: salida.Usuario.RolID,
		},
	})
}

func (h *AuthHandler) Recuperacion(c *fiber.Ctx) error {
	var in dto.RecuperacionEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}
	if err := h.SolicitarRecuperacion.Ejecutar(c.Context(), in.Correo); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func (h *AuthHandler) Restablecer(c *fiber.Ctx) error {
	return h.restablecerConCaso(c, h.RestablecerRecuperacion)
}

func (h *AuthHandler) Activar(c *fiber.Ctx) error {
	return h.restablecerConCaso(c, h.RestablecerActivacion)
}

func (h *AuthHandler) restablecerConCaso(c *fiber.Ctx, caso *auth.RestablecerContrasena) error {
	var in dto.RestablecerEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}
	if err := caso.Ejecutar(c.Context(), in.Token, in.Nueva); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	id, ok := middleware.DeContexto(c)
	if !ok {
		return c.SendStatus(fiber.StatusNoContent)
	}
	restante := time.Until(id.Expira)
	if err := h.CerrarSesionCaso.Ejecutar(c.Context(), id.JTI, restante); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) Sesion(c *fiber.Ctx) error {
	id, _ := middleware.DeContexto(c)
	permisos := make([]string, 0, len(id.Permisos))
	for p := range id.Permisos {
		permisos = append(permisos, p)
	}
	sort.Strings(permisos)
	salida, err := h.ObtenerSesionCaso.Ejecutar(c.Context(), id.UsuarioID, permisos)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(dto.SesionSalida{
		Rol: salida.Rol, Permisos: salida.Permisos,
		Usuario: dto.UsuarioResumen{
			ID: salida.Usuario.ID, Nombres: salida.Usuario.Nombres, Apellidos: salida.Usuario.Apellidos,
			Correo: salida.Usuario.Correo, RolID: salida.Usuario.RolID,
		},
	})
}
