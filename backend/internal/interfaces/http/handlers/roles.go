package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/roles"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// RolesHandler expone los endpoints de 7.4 (HU-02).
type RolesHandler struct {
	CasoListar         *roles.Listar
	CasoObtener        *roles.Obtener
	CasoListarPermisos *roles.ListarPermisos
	CasoGuardar        *roles.Guardar
	CasoEliminar       *roles.Eliminar
}

func aPermisoDTO(p rol.Permiso) dto.PermisoDTO {
	return dto.PermisoDTO{ID: p.ID, Codigo: p.Codigo, Nombre: p.Nombre, Modulo: p.Modulo, CasoUso: p.CasoUso}
}

func (h *RolesHandler) ListarPermisos(c *fiber.Ctx) error {
	permisos, err := h.CasoListarPermisos.Ejecutar(c.Context())
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.PermisoDTO, len(permisos))
	for i, p := range permisos {
		out[i] = aPermisoDTO(p)
	}
	return c.JSON(out)
}

func (h *RolesHandler) Listar(c *fiber.Ctx) error {
	lista, err := h.CasoListar.Ejecutar(c.Context())
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.RolConConteoDTO, len(lista))
	for i, r := range lista {
		out[i] = dto.RolConConteoDTO{
			ID: r.Rol.ID, Nombre: r.Rol.Nombre, Descripcion: r.Rol.Descripcion,
			EsSistema: r.Rol.EsSistema, Activo: r.Rol.Activo, UsuariosCount: r.UsuariosCount,
		}
	}
	return c.JSON(out)
}

func (h *RolesHandler) Obtener(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apierror.Nuevo(400, "ID_INVALIDO", "El identificador no es válido")
	}
	r, err := h.CasoObtener.Ejecutar(c.Context(), int16(id))
	if err != nil {
		return apierror.Mapear(err)
	}
	permisos := make([]dto.PermisoDTO, len(r.Permisos))
	for i, p := range r.Permisos {
		permisos[i] = aPermisoDTO(p)
	}
	return c.JSON(dto.RolConPermisosDTO{
		ID: r.ID, Nombre: r.Nombre, Descripcion: r.Descripcion, EsSistema: r.EsSistema, Activo: r.Activo, Permisos: permisos,
	})
}

func (h *RolesHandler) Guardar(c *fiber.Ctx) error {
	var in dto.GuardarRolEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	var id int16
	if idStr := c.Params("id"); idStr != "" {
		n, err := strconv.Atoi(idStr)
		if err != nil {
			return apierror.Nuevo(400, "ID_INVALIDO", "El identificador no es válido")
		}
		id = int16(n)
	}

	actor, _ := middleware.DeContexto(c)
	r, err := h.CasoGuardar.Ejecutar(c.Context(), roles.GuardarEntrada{
		ActorID: actor.UsuarioID, ID: id, Nombre: in.Nombre, Descripcion: in.Descripcion, PermisoIDs: in.PermisoIDs,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	status := fiber.StatusOK
	if id == 0 {
		status = fiber.StatusCreated
	}
	return c.Status(status).JSON(dto.RolConConteoDTO{ID: r.ID, Nombre: r.Nombre, Descripcion: r.Descripcion, EsSistema: r.EsSistema, Activo: r.Activo})
}

func (h *RolesHandler) Eliminar(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apierror.Nuevo(400, "ID_INVALIDO", "El identificador no es válido")
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoEliminar.Ejecutar(c.Context(), actor.UsuarioID, int16(id)); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
