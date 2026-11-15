package handlers

import (
	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/perfil"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

type PerfilHandler struct {
	CasoObtener            *perfil.ObtenerPerfil
	CasoActualizarContacto *perfil.ActualizarContacto
	CasoCambiarContrasena  *perfil.CambiarContrasena
}

func (h *PerfilHandler) Obtener(c *fiber.Ctx) error {
	id, _ := middleware.DeContexto(c)
	u, err := h.CasoObtener.Ejecutar(c.Context(), id.UsuarioID)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(dto.PerfilSalida{ID: u.ID, Nombres: u.Nombres, Apellidos: u.Apellidos, Correo: u.Correo, Telefono: u.Telefono, RolID: u.RolID})
}

func (h *PerfilHandler) ActualizarContacto(c *fiber.Ctx) error {
	var crudo map[string]any
	if err := c.BodyParser(&crudo); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if dto.CamposProhibidosPresentes(crudo) {
		return apierror.Nuevo(403, "CAMPO_NO_EDITABLE", "Solo puedes editar tu teléfono desde aquí")
	}
	var in dto.ActualizarContactoEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	id, _ := middleware.DeContexto(c)
	u, err := h.CasoActualizarContacto.Ejecutar(c.Context(), id.UsuarioID, in.Telefono)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(dto.PerfilSalida{ID: u.ID, Nombres: u.Nombres, Apellidos: u.Apellidos, Correo: u.Correo, Telefono: u.Telefono, RolID: u.RolID})
}

func (h *PerfilHandler) CambiarContrasena(c *fiber.Ctx) error {
	var in dto.CambiarContrasenaEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	id, _ := middleware.DeContexto(c)
	salida, err := h.CasoCambiarContrasena.Ejecutar(c.Context(), id.UsuarioID, in.Actual, in.Nueva)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(dto.CambiarContrasenaSalida{Token: salida.Token, ExpiraEn: salida.ExpiraEn})
}
