package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/indicadores"
	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// IndicadoresHandler expone el tablero de indicadores (CU-25) y la generación de reportes (CU-26).
type IndicadoresHandler struct {
	CasoTablero *indicadores.ConsultarTablero
	CasoReporte *indicadores.GenerarReporte
	// CasoInterpretar es opcional (nil = ruta no registrada).
	CasoInterpretar *indicadores.InterpretarReporte
}

// leerPeriodo interpreta desde y hasta (AAAA-MM-DD, días de Bolivia, ambos inclusive); los errores
// de formato quedan en campos con las claves indicadas.
func leerPeriodo(desde, hasta, claveDesde, claveHasta string, campos map[string]string) indicador.Periodo {
	fecha := func(clave, v string) *time.Time {
		if v = strings.TrimSpace(v); v == "" {
			return nil
		}
		t, err := time.ParseInLocation(formatoFecha, v, zonaBolivia)
		if err != nil {
			campos[clave] = "Debe ser una fecha AAAA-MM-DD"
			return nil
		}
		return &t
	}
	p := indicador.Periodo{Desde: fecha(claveDesde, desde)}
	if h := fecha(claveHasta, hasta); h != nil {
		fin := h.AddDate(0, 0, 1) // el día "hasta" completo
		p.Hasta = &fin
	}
	return p
}

func textoOpcional(s string) *string {
	if s = strings.TrimSpace(s); s == "" {
		return nil
	}
	return &s
}

// Tablero: GET /api/tablero?desde=&hasta= (el periodo acota los trámites radicados y resueltos;
// la carga vigente y las alertas de mora son siempre las de hoy).
func (h *IndicadoresHandler) Tablero(c *fiber.Ctx) error {
	campos := map[string]string{}
	p := leerPeriodo(c.Query("desde"), c.Query("hasta"), "desde", "hasta", campos)
	if len(campos) > 0 {
		return apierror.Validacion(campos)
	}
	t, err := h.CasoTablero.Ejecutar(c.Context(), p)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.TableroDTO{
		GeneradoEn: t.GeneradoEn, DiasMaximosPrestamo: t.DiasMaximosPrestamo, PrestamosVigentes: t.PrestamosVigentes,
		Periodo:   dto.PeriodoDTO{Desde: textoOpcional(c.Query("desde")), Hasta: textoOpcional(c.Query("hasta"))},
		Tramites:  dto.ResumenTramitesDTO{Radicados: t.Radicados, EnCurso: t.EnCurso, Resueltos: t.Resueltos, PromedioDias: t.PromedioDias},
		PorEstado: make([]dto.ConteoEstadoDTO, len(t.PorEstado)), Tiempos: make([]dto.TiempoResolucionDTO, len(t.Tiempos)),
		Analistas: make([]dto.CargaAnalistaDTO, len(t.Analistas)), Alertas: make([]dto.AlertaMoraDTO, len(t.Alertas)),
	}
	for i, e := range t.PorEstado {
		out.PorEstado[i] = dto.ConteoEstadoDTO{Estado: string(e.Estado), Cantidad: e.Cantidad}
	}
	for i, g := range t.Tiempos {
		out.Tiempos[i] = dto.TiempoResolucionDTO{
			Tipo: string(g.Tipo), Categoria: g.Categoria, Cantidad: g.Cantidad,
			PromedioDias: g.PromedioDias, MinimoDias: g.MinimoDias, MaximoDias: g.MaximoDias,
		}
	}
	for i, a := range t.Analistas {
		out.Analistas[i] = dto.CargaAnalistaDTO{ID: a.ID, Nombre: a.Nombre, Activos: a.Activos, EnFacultad: a.EnFacultad, Resueltos: a.Resueltos}
	}
	for i, a := range t.Alertas {
		out.Alertas[i] = dto.AlertaMoraDTO{
			LegajoCodigo: a.LegajoCodigo, DocenteNombre: a.DocenteNombre, Custodio: a.Custodio, TramiteCodigo: a.TramiteCodigo,
			Salida: a.Salida, FechaLimite: fechaOpcional(a.FechaLimite), DiasPrestado: a.DiasPrestado, DiasMora: a.DiasMora,
			Motivo: string(a.Motivo),
		}
	}
	return c.JSON(out)
}

// InterpretarReporte: POST /api/reportes/interpretar. Devuelve los filtros que se entendieron de una frase
// (con los mismos nombres que GenerarReporteEntrada); no genera ni audita nada.
func (h *IndicadoresHandler) InterpretarReporte(c *fiber.Ctx) error {
	var in dto.InterpretarReporteEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	cats := make([]ports.CategoriaRef, len(in.Categorias))
	for i, x := range in.Categorias {
		cats[i] = ports.CategoriaRef{ID: x.ID, Nombre: x.Nombre}
	}
	p, err := h.CasoInterpretar.Ejecutar(c.Context(), in.Frase, cats)
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(fiber.Map{
		"tipo": p.Tipo, "formato": p.Formato, "desde": p.Desde, "hasta": p.Hasta, "estado": p.Estado,
		"tipoTramite": p.TipoTramite, "categoriaId": p.CategoriaID, "estadoLegajo": p.EstadoLegajo, "explicacion": p.Explicacion,
	})
}

// Reporte: POST /api/reportes. Responde el archivo (PDF o XLSX) como adjunto; la emisión queda en auditoría.
func (h *IndicadoresHandler) Reporte(c *fiber.Ctx) error {
	var in dto.GenerarReporteEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	campos := map[string]string{}
	p := leerPeriodo(in.Desde, in.Hasta, "Desde", "Hasta", campos)
	if len(campos) > 0 {
		return apierror.Validacion(campos)
	}
	actor, _ := middleware.DeContexto(c)
	a, err := h.CasoReporte.Ejecutar(c.Context(), indicadores.SolicitudReporte{
		ActorID: actor.UsuarioID, Tipo: indicador.TipoReporte(in.Tipo), Formato: indicador.Formato(in.Formato),
		Filtro: indicador.FiltroReporte{
			Periodo: p, Estado: tramite.Estado(strings.ToUpper(strings.TrimSpace(in.Estado))),
			Tipo: tramite.Tipo(strings.ToUpper(strings.TrimSpace(in.TipoTramite))), CategoriaID: in.CategoriaID,
			EstadoLegajo: legajo.Estado(strings.ToUpper(strings.TrimSpace(in.EstadoLegajo))),
		},
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	c.Set(fiber.HeaderContentType, a.TipoMime)
	c.Set(fiber.HeaderContentDisposition, `attachment; filename="`+a.Nombre+`"`)
	c.Set("X-Reporte-Filas", strconv.Itoa(a.Filas))
	return c.Send(a.Datos)
}
