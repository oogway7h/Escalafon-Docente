package handlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/categorias"
	"escalafon/backend/internal/application/docentes"
	"escalafon/backend/internal/domain/categoria"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

func aDocenteDTO(d docente.Docente) dto.DocenteDTO {
	legajos := make([]dto.LegajoDTO, len(d.Legajos))
	for i, l := range d.Legajos {
		legajos[i] = dto.LegajoDTO{Gaveta: l.Gaveta, Orden: l.Orden, Anexos: l.Anexos, Observacion: l.Observacion}
	}
	return dto.DocenteDTO{
		ID: d.ID, Codigo: d.Codigo, Cedula: d.Cedula, Nombres: d.Nombres, Apellidos: d.Apellidos, Correo: d.Correo,
		Telefono: d.Telefono, CategoriaID: d.CategoriaID, Categoria: d.Categoria, FacultadID: d.FacultadID,
		Facultad: d.Facultad, Legajos: legajos, Activo: d.Activo,
	}
}

func aCategoriaDTO(c categoria.Categoria) dto.CategoriaDTO {
	reqs := make([]dto.RequisitoDTO, len(c.Requisitos))
	for i, r := range c.Requisitos {
		reqs[i] = dto.RequisitoDTO{ID: r.ID, Nombre: r.Nombre}
	}
	return dto.CategoriaDTO{
		ID: c.ID, Nombre: c.Nombre, AniosAntiguedad: c.AniosAntiguedad, Activa: c.Activa, Requisitos: reqs,
		Docentes: c.Docentes, Tramites: c.Tramites,
	}
}

func leerCuerpo(c *fiber.Ctx, destino any) error {
	if err := c.BodyParser(destino); err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "El cuerpo de la solicitud no es JSON válido")
	}
	if campos := dto.Validar(destino); campos != nil {
		return apierror.Validacion(campos)
	}
	return nil
}

// DocentesHandler expone /api/docentes (CU-07).
type DocentesHandler struct {
	CasoListar        *docentes.Listar
	CasoListarPagina  *docentes.ListarPagina
	CasoSugerir       *docentes.Sugerir
	CasoObtener       *docentes.Obtener
	CasoGuardar       *docentes.Guardar
	CasoCambiarEstado *docentes.CambiarEstado
	CasoFacultades    *docentes.ListarFacultades
}

func (h *DocentesHandler) ListarFacultades(c *fiber.Ctx) error {
	lista, err := h.CasoFacultades.Ejecutar(c.Context())
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.FacultadDTO, len(lista))
	for i, f := range lista {
		out[i] = dto.FacultadDTO{ID: f.ID, Nombre: f.Nombre}
	}
	return c.JSON(out)
}

// Sugerir devuelve hasta 8 docentes activos para el buscador con sugerencias.
func (h *DocentesHandler) Sugerir(c *fiber.Ctx) error {
	lista, err := h.CasoSugerir.Ejecutar(c.Context(), c.Query("q"))
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.DocenteDTO, len(lista))
	for i, d := range lista {
		out[i] = aDocenteDTO(d)
	}
	return c.JSON(out)
}

// Listar devuelve el arreglo completo; con ?pagina=N devuelve una página (?tamano=, 10 por defecto)
// con la forma {datos, total, pagina, tamano}.
func (h *DocentesHandler) Listar(c *fiber.Ctx) error {
	if c.Query("pagina") != "" {
		pagina, _ := strconv.Atoi(c.Query("pagina"))
		tamano, _ := strconv.Atoi(c.Query("tamano", "10"))
		if pagina < 1 {
			pagina = 1
		}
		if tamano < 1 || tamano > 100 {
			tamano = 10
		}
		lista, total, err := h.CasoListarPagina.Ejecutar(c.Context(), c.Query("buscar"), c.Query("activo") == "true", pagina, tamano)
		if err != nil {
			return apierror.Mapear(err)
		}
		out := dto.ListaDocentesDTO{Datos: make([]dto.DocenteDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
		for i, d := range lista {
			out.Datos[i] = aDocenteDTO(d)
		}
		return c.JSON(out)
	}
	lista, err := h.CasoListar.Ejecutar(c.Context(), c.Query("buscar"), c.Query("activo") == "true")
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.DocenteDTO, len(lista))
	for i, d := range lista {
		out[i] = aDocenteDTO(d)
	}
	return c.JSON(out)
}

func (h *DocentesHandler) Obtener(c *fiber.Ctx) error {
	d, err := h.CasoObtener.Ejecutar(c.Context(), c.Params("id"))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aDocenteDTO(*d))
}

func (h *DocentesHandler) Guardar(c *fiber.Ctx) error {
	var in dto.GuardarDocenteEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	id := c.Params("id")
	if id == "" && strings.TrimSpace(in.Codigo) == "" {
		return apierror.Validacion(map[string]string{"Codigo": "Este campo es obligatorio"})
	}
	actor, _ := middleware.DeContexto(c)
	d, err := h.CasoGuardar.Ejecutar(c.Context(), docentes.Entrada{
		ActorID: actor.UsuarioID, ID: id, Codigo: in.Codigo, Cedula: in.Cedula, Nombres: in.Nombres, Apellidos: in.Apellidos,
		Correo: in.Correo, Telefono: in.Telefono, CategoriaID: in.CategoriaID, FacultadID: in.FacultadID,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	status := fiber.StatusOK
	if id == "" {
		status = fiber.StatusCreated
	}
	return c.Status(status).JSON(aDocenteDTO(*d))
}

func (h *DocentesHandler) CambiarEstado(c *fiber.Ctx) error {
	var in dto.CambiarEstadoEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoCambiarEstado.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), *in.Activo); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// CategoriasHandler expone /api/categorias (CU-07) y el catálogo público de categorías.
type CategoriasHandler struct {
	CasoListar        *categorias.Listar
	CasoObtener       *categorias.Obtener
	CasoGuardar       *categorias.Guardar
	CasoCambiarEstado *categorias.CambiarEstado
	CasoEliminar      *categorias.Eliminar
}

func (h *CategoriasHandler) Eliminar(c *fiber.Ctx) error {
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

func (h *CategoriasHandler) Listar(c *fiber.Ctx) error {
	lista, err := h.CasoListar.Ejecutar(c.Context(), c.Query("activa") == "true")
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.CategoriaDTO, len(lista))
	for i, cat := range lista {
		out[i] = aCategoriaDTO(cat)
	}
	return c.JSON(out)
}

// ListarPublicas devuelve solo las categorías activas, sin autenticación.
func (h *CategoriasHandler) ListarPublicas(c *fiber.Ctx) error {
	lista, err := h.CasoListar.Ejecutar(c.Context(), true)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.CategoriaDTO, len(lista))
	for i, cat := range lista {
		out[i] = aCategoriaDTO(cat)
	}
	return c.JSON(out)
}

func (h *CategoriasHandler) Obtener(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apierror.Nuevo(400, "ID_INVALIDO", "El identificador no es válido")
	}
	cat, err := h.CasoObtener.Ejecutar(c.Context(), int16(id))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aCategoriaDTO(*cat))
}

func (h *CategoriasHandler) Guardar(c *fiber.Ctx) error {
	var in dto.GuardarCategoriaEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	var id int16
	if s := c.Params("id"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return apierror.Nuevo(400, "ID_INVALIDO", "El identificador no es válido")
		}
		id = int16(n)
	}
	actor, _ := middleware.DeContexto(c)
	cat, err := h.CasoGuardar.Ejecutar(c.Context(), categorias.Entrada{
		ActorID: actor.UsuarioID, ID: id, Nombre: in.Nombre, AniosAntiguedad: in.AniosAntiguedad, Requisitos: in.Requisitos,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	status := fiber.StatusOK
	if id == 0 {
		status = fiber.StatusCreated
	}
	return c.Status(status).JSON(aCategoriaDTO(*cat))
}

func (h *CategoriasHandler) CambiarEstado(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apierror.Nuevo(400, "ID_INVALIDO", "El identificador no es válido")
	}
	var in dto.CambiarEstadoEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoCambiarEstado.Ejecutar(c.Context(), actor.UsuarioID, int16(id), *in.Activo); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
