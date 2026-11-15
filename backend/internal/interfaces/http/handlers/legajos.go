package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/legajos"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// LegajosHandler expone el inventario de carpetas, sus solicitudes y la custodia (CU-19 al CU-24).
type LegajosHandler struct {
	CasoListar            *legajos.ListarLegajos
	CasoObtener           *legajos.ObtenerLegajo
	CasoGuardar           *legajos.GuardarLegajo
	CasoSolicitar         *legajos.Solicitar
	CasoListarSolicitudes *legajos.ListarSolicitudes
	CasoAtender           *legajos.Atender
	CasoDespachar         *legajos.Despachar
	CasoEntregar          *legajos.Entregar
	CasoDevolver          *legajos.Devolver
	CasoListarMovimientos *legajos.ListarMovimientos
}

func paginacion(c *fiber.Ctx) (int, int) {
	pagina, _ := strconv.Atoi(c.Query("pagina", "1"))
	tamano, _ := strconv.Atoi(c.Query("tamano", "10"))
	return max(pagina, 1), min(max(tamano, 1), 100)
}

func idLegajo(c *fiber.Ctx) (int32, error) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return 0, apierror.Nuevo(400, "ID_INVALIDO", "El identificador no es válido")
	}
	return int32(id), nil
}

func aLegajoDTO(l legajo.Legajo) dto.CarpetaDTO {
	return dto.CarpetaDTO{
		ID: l.ID, Codigo: l.Codigo, DocenteID: l.DocenteID, DocenteCodigo: l.DocenteCodigo, DocenteNombre: l.DocenteNombre,
		Gaveta: l.Gaveta, Orden: l.Orden, Anexos: l.Anexos, Observacion: l.Observacion, Estado: string(l.Estado),
	}
}

func fechaOpcional(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(formatoFecha)
	return &s
}

func aSolicitudDTO(s legajo.Solicitud) dto.SolicitudCarpetaDTO {
	return dto.SolicitudCarpetaDTO{
		ID: s.ID, LegajoID: s.LegajoID, LegajoCodigo: s.LegajoCodigo, DocenteNombre: s.DocenteNombre, TramiteID: s.TramiteID, TramiteCodigo: s.TramiteCodigo, Solicitante: s.Solicitante,
		Motivo: s.Motivo, Estado: string(s.Estado), Respuesta: s.Respuesta, CreadoEn: s.CreadoEn, AtendidoEn: s.AtendidoEn,
		FechaLimite: fechaOpcional(s.FechaLimite),
	}
}

func (h *LegajosHandler) Listar(c *fiber.Ctx) error {
	pagina, tamano := paginacion(c)
	f := legajo.FiltroLegajos{Busqueda: c.Query("buscar"), Estado: legajo.Estado(c.Query("estado"))}
	lista, total, err := h.CasoListar.Ejecutar(c.Context(), f, pagina, tamano)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.ListaLegajosDTO{Datos: make([]dto.CarpetaDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
	for i, l := range lista {
		out.Datos[i] = aLegajoDTO(l)
	}
	return c.JSON(out)
}

func (h *LegajosHandler) Obtener(c *fiber.Ctx) error {
	id, err := idLegajo(c)
	if err != nil {
		return err
	}
	l, err := h.CasoObtener.Ejecutar(c.Context(), id)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aLegajoDTO(*l))
}

func (h *LegajosHandler) Guardar(c *fiber.Ctx) error {
	var in dto.GuardarLegajoEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	var id int32
	if c.Params("id") != "" {
		var err error
		if id, err = idLegajo(c); err != nil {
			return err
		}
	} else if strings.TrimSpace(in.DocenteID) == "" {
		return apierror.Validacion(map[string]string{"DocenteID": "Este campo es obligatorio"})
	}
	actor, _ := middleware.DeContexto(c)
	l, err := h.CasoGuardar.Ejecutar(c.Context(), legajos.GuardarEntrada{
		ActorID: actor.UsuarioID, ID: id, DocenteID: in.DocenteID, Codigo: in.Codigo, Gaveta: in.Gaveta, Orden: in.Orden,
		Anexos: in.Anexos, Observacion: in.Observacion,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	status := fiber.StatusOK
	if id == 0 {
		status = fiber.StatusCreated
	}
	return c.Status(status).JSON(aLegajoDTO(*l))
}

func (h *LegajosHandler) Solicitar(c *fiber.Ctx) error {
	id, err := idLegajo(c)
	if err != nil {
		return err
	}
	var in dto.SolicitarCarpetaEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	s, err := h.CasoSolicitar.Ejecutar(c.Context(), actor.UsuarioID, id, in.TramiteID, in.Motivo)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.Status(fiber.StatusCreated).JSON(aSolicitudDTO(*s))
}

// ListarSolicitudes: quien solo puede solicitar (analista) ve únicamente las suyas.
func (h *LegajosHandler) ListarSolicitudes(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	pagina, tamano := paginacion(c)
	f := legajo.FiltroSolicitudes{Estado: legajo.EstadoSolicitud(c.Query("estado"))}
	if !actor.TienePermiso(rol.PermisoLegajoGestionar) && !actor.TienePermiso(rol.PermisoCustodiaRegistrar) {
		f.SolicitanteID = actor.UsuarioID
	}
	lista, total, err := h.CasoListarSolicitudes.Ejecutar(c.Context(), f, pagina, tamano)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.ListaSolicitudesDTO{Datos: make([]dto.SolicitudCarpetaDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
	for i, s := range lista {
		out.Datos[i] = aSolicitudDTO(s)
	}
	return c.JSON(out)
}

func (h *LegajosHandler) Atender(c *fiber.Ctx) error {
	var in dto.AtenderSolicitudEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoAtender.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), in.Aprobar, in.Nota); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LegajosHandler) Despachar(c *fiber.Ctx) error {
	var in dto.DespacharEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	limite, err := time.Parse(formatoFecha, in.FechaLimite)
	if err != nil {
		return apierror.Validacion(map[string]string{"FechaLimite": "Debe tener el formato AAAA-MM-DD"})
	}
	actor, _ := middleware.DeContexto(c)
	err = h.CasoDespachar.Ejecutar(c.Context(), legajos.DespacharEntrada{
		ActorID: actor.UsuarioID, SolicitudID: c.Params("id"), FechaLimite: limite,
		EstadoMaterial: legajo.EstadoMaterial(in.EstadoMaterial), Observacion: in.Observacion,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Entregar aprueba y despacha en un solo paso (pantalla unificada de solicitudes y custodia).
func (h *LegajosHandler) Entregar(c *fiber.Ctx) error {
	var in dto.DespacharEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	limite, err := time.Parse(formatoFecha, in.FechaLimite)
	if err != nil {
		return apierror.Validacion(map[string]string{"FechaLimite": "Debe tener el formato AAAA-MM-DD"})
	}
	actor, _ := middleware.DeContexto(c)
	err = h.CasoEntregar.Ejecutar(c.Context(), legajos.DespacharEntrada{
		ActorID: actor.UsuarioID, SolicitudID: c.Params("id"), FechaLimite: limite,
		EstadoMaterial: legajo.EstadoMaterial(in.EstadoMaterial), Observacion: in.Observacion,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LegajosHandler) Devolver(c *fiber.Ctx) error {
	id, err := idLegajo(c)
	if err != nil {
		return err
	}
	var in dto.DevolverEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	err = h.CasoDevolver.Ejecutar(c.Context(), legajos.DevolverEntrada{
		ActorID: actor.UsuarioID, LegajoID: id, EstadoMaterial: legajo.EstadoMaterial(in.EstadoMaterial), Observacion: in.Observacion,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LegajosHandler) Bitacora(c *fiber.Ctx) error {
	pagina, tamano := paginacion(c)
	f := legajo.FiltroMovimientos{Busqueda: c.Query("buscar"), Tipo: legajo.TipoMovimiento(c.Query("tipo")), Vencidos: c.Query("vencidos") == "true"}
	lista, total, err := h.CasoListarMovimientos.Ejecutar(c.Context(), f, pagina, tamano)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.ListaMovimientosDTO{Datos: make([]dto.MovimientoDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
	for i, m := range lista {
		out.Datos[i] = dto.MovimientoDTO{
			ID: m.ID, LegajoCodigo: m.LegajoCodigo, DocenteNombre: m.DocenteNombre, Tipo: string(m.Tipo), Custodio: m.Custodio,
			FechaLimite: fechaOpcional(m.FechaLimite), EstadoMaterial: string(m.EstadoMaterial), Observacion: m.Observacion,
			RegistradoPor: m.RegistradoPor, Fecha: m.Fecha, Vencido: m.Vencido,
		}
	}
	return c.JSON(out)
}
