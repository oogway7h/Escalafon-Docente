package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/usuarios"
	"escalafon/backend/internal/domain/usuario"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// UsuariosHandler expone los endpoints de 7.3 (HU-02).
type UsuariosHandler struct {
	CasoListar             *usuarios.Listar
	CasoObtener            *usuarios.Obtener
	CasoCrear              *usuarios.Crear
	CasoEditar             *usuarios.Editar
	CasoCambiarRol         *usuarios.CambiarRol
	CasoCambiarEstado      *usuarios.CambiarEstado
	CasoReenviarActivacion *usuarios.ReenviarActivacion
}

func aUsuarioListado(u usuario.Usuario) dto.UsuarioListado {
	return dto.UsuarioListado{
		ID: u.ID, Nombres: u.Nombres, Apellidos: u.Apellidos, Correo: u.Correo,
		Telefono: u.Telefono, RolID: u.RolID, Activo: u.Activo,
	}
}

func (h *UsuariosHandler) Listar(c *fiber.Ctx) error {
	pagina, _ := strconv.Atoi(c.Query("pagina", "1"))
	tamano, _ := strconv.Atoi(c.Query("tamano", "20"))
	filtro := usuario.ListarFiltro{Buscar: c.Query("buscar"), Pagina: pagina, Tamano: tamano}
	if v := c.Query("rolId"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			r := int16(n)
			filtro.RolID = &r
		}
	}
	if v := c.Query("activo"); v != "" {
		a := v == "true"
		filtro.Activo = &a
	}

	lista, total, err := h.CasoListar.Ejecutar(c.Context(), filtro)
	if err != nil {
		return apierror.Mapear(err)
	}
	datos := make([]dto.UsuarioListado, len(lista))
	for i, u := range lista {
		datos[i] = aUsuarioListado(u)
	}
	return c.JSON(dto.ListarUsuariosSalida{Datos: datos, Total: total, Pagina: filtro.Pagina, Tamano: filtro.Tamano})
}

func (h *UsuariosHandler) Obtener(c *fiber.Ctx) error {
	u, err := h.CasoObtener.Ejecutar(c.Context(), c.Params("id"))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aUsuarioListado(*u))
}

func (h *UsuariosHandler) Crear(c *fiber.Ctx) error {
	var in dto.CrearUsuarioEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	actor, _ := middleware.DeContexto(c)
	u, err := h.CasoCrear.Ejecutar(c.Context(), usuarios.CrearEntrada{
		ActorID: actor.UsuarioID, Nombres: in.Nombres, Apellidos: in.Apellidos, Correo: in.Correo,
		Telefono: in.Telefono, RolID: in.RolID,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.Status(fiber.StatusCreated).JSON(aUsuarioListado(*u))
}

func (h *UsuariosHandler) Editar(c *fiber.Ctx) error {
	var in dto.EditarUsuarioEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	actor, _ := middleware.DeContexto(c)
	u, err := h.CasoEditar.Ejecutar(c.Context(), c.Params("id"), usuarios.EditarEntrada{
		ActorID: actor.UsuarioID, Nombres: in.Nombres, Apellidos: in.Apellidos, Correo: in.Correo, Telefono: in.Telefono,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aUsuarioListado(*u))
}

func (h *UsuariosHandler) CambiarRol(c *fiber.Ctx) error {
	var in dto.CambiarRolEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	actor, _ := middleware.DeContexto(c)
	if err := h.CasoCambiarRol.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), in.RolID); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsuariosHandler) CambiarEstado(c *fiber.Ctx) error {
	var in dto.CambiarEstadoEntrada
	if err := c.BodyParser(&in); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(in); campos != nil {
		return apierror.Validacion(campos)
	}

	actor, _ := middleware.DeContexto(c)
	if err := h.CasoCambiarEstado.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), *in.Activo); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *UsuariosHandler) ReenviarActivacion(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoReenviarActivacion.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id")); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusAccepted)
}
