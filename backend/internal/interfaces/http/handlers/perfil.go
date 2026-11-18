package handlers

import (
	"io"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/perfil"
	"escalafon/backend/internal/domain/usuario"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

type PerfilHandler struct {
	CasoObtener            *perfil.ObtenerPerfil
	CasoActualizarContacto *perfil.ActualizarContacto
	CasoCambiarContrasena  *perfil.CambiarContrasena
	CasoObtenerFoto        *perfil.ObtenerFoto
	CasoSubirFoto          *perfil.SubirFoto
	CasoQuitarFoto         *perfil.QuitarFoto
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

// ObtenerFoto devuelve un enlace firmado y temporal a la foto del usuario autenticado (404 si no tiene).
func (h *PerfilHandler) ObtenerFoto(c *fiber.Ctx) error {
	id, _ := middleware.DeContexto(c)
	f, err := h.CasoObtenerFoto.Ejecutar(c.Context(), id.UsuarioID)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(dto.FotoSalida{URL: f.URL, TipoMime: f.TipoMime, Version: f.Version, VigenciaSegundos: f.VigenciaS})
}

// SubirFoto recibe la imagen como multipart (campo "foto"); el contenido se valida en el caso de uso.
func (h *PerfilHandler) SubirFoto(c *fiber.Ctx) error {
	cabecera, err := c.FormFile("foto")
	if err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "Envía la imagen en el campo «foto» de un formulario multipart")
	}
	if cabecera.Size > usuario.FotoMaxBytes {
		return apierror.Mapear(usuario.ErrFotoGrande)
	}
	archivo, err := cabecera.Open()
	if err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "No se pudo leer la imagen")
	}
	defer archivo.Close()
	datos, err := io.ReadAll(io.LimitReader(archivo, usuario.FotoMaxBytes+1))
	if err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "No se pudo leer la imagen")
	}

	id, _ := middleware.DeContexto(c)
	f, err := h.CasoSubirFoto.Ejecutar(c.Context(), id.UsuarioID, datos)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(fiber.Map{"version": f.ActualizadoEn.Unix()})
}

func (h *PerfilHandler) QuitarFoto(c *fiber.Ctx) error {
	id, _ := middleware.DeContexto(c)
	if err := h.CasoQuitarFoto.Ejecutar(c.Context(), id.UsuarioID); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(204)
}
