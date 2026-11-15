package indicadores

import (
	"context"
	"fmt"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/domain/usuario"
)

// Etiquetas en español de los códigos que aparecen en los reportes (las mismas que muestra el panel).
var (
	etiquetaEstado = map[tramite.Estado]string{
		tramite.Ingresado: "Ingresado", tramite.Asignado: "Asignado", tramite.EnRevision: "En revisión",
		tramite.Observado: "Observado", tramite.Dictaminado: "Dictaminado", tramite.EnFacultad: "Fuera de Escalafón",
		tramite.RevisionFinal: "Revisión final", tramite.Remitido: "Remitido a Rectorado", tramite.Archivado: "Archivado",
	}
	etiquetaTipo      = map[tramite.Tipo]string{tramite.Categorizacion: "Categorización", tramite.Regularizacion: "Regularización"}
	etiquetaOrigen    = map[tramite.Origen]string{tramite.OrigenStaff: "Ventanilla", tramite.OrigenPortal: "Portal"}
	etiquetaResultado = map[tramite.Resultado]string{tramite.Favorable: "Favorable", tramite.Desfavorable: "Desfavorable"}
	etiquetaLegajo    = map[legajo.Estado]string{legajo.EnArchivo: "En archivo", legajo.Prestado: "Prestado"}
)

// Archivo es el reporte ya compuesto, listo para descargar.
type Archivo struct {
	Nombre   string
	TipoMime string
	Datos    []byte
	Filas    int
}

type SolicitudReporte struct {
	ActorID string
	Tipo    indicador.TipoReporte
	Formato indicador.Formato
	Filtro  indicador.FiltroReporte
}

// GenerarReporte recupera los datos, verifica que no estén vacíos, compone el archivo en el
// formato pedido y registra la emisión en la auditoría (CU-26).
type GenerarReporte struct {
	Indicadores         indicador.Consulta
	Usuarios            usuario.Repository
	Exportador          ports.ExportadorReportes
	Auditoria           auditoria.Repository
	Reloj               ports.Clock
	DiasMaximosPrestamo int
}

func (c *GenerarReporte) Ejecutar(ctx context.Context, s SolicitudReporte) (*Archivo, error) {
	if !s.Tipo.Valido() {
		return nil, indicador.ErrTipoReporte
	}
	if !s.Formato.Valido() {
		return nil, indicador.ErrFormato
	}
	f := s.Filtro
	if !f.Valido() {
		return nil, indicador.ErrRangoFechas
	}
	if (f.Estado != "" && !estadoValido(f.Estado)) || (f.Tipo != "" && !f.Tipo.Valido()) || f.CategoriaID < 0 ||
		(f.EstadoLegajo != "" && f.EstadoLegajo != legajo.EnArchivo && f.EstadoLegajo != legajo.Prestado) {
		return nil, indicador.ErrFiltroInvalido
	}

	r := &indicador.Reporte{Tipo: s.Tipo, EmitidoEn: c.Reloj.Ahora()}
	var err error
	switch s.Tipo {
	case indicador.ReporteTramites:
		err = c.tramites(ctx, r, f)
	case indicador.ReporteInventario:
		err = c.inventario(ctx, r, f)
	case indicador.ReporteTiempos:
		err = c.tiempos(ctx, r, f)
	}
	if err != nil {
		return nil, err
	}
	if len(r.Filas) == 0 {
		return nil, indicador.ErrReporteVacio
	}

	u, err := c.Usuarios.BuscarPorID(ctx, s.ActorID)
	if err != nil {
		return nil, err
	}
	r.EmitidoPor = strings.TrimSpace(u.Nombres + " " + u.Apellidos)

	a := &Archivo{Filas: len(r.Filas)}
	nombre := fmt.Sprintf("reporte-%s-%s", strings.ToLower(string(s.Tipo)), r.EmitidoEn.In(zonaBolivia).Format("2006-01-02"))
	if s.Formato == indicador.FormatoPDF {
		a.Datos, err = c.Exportador.PDF(r)
		a.Nombre, a.TipoMime = nombre+".pdf", "application/pdf"
	} else {
		a.Datos, err = c.Exportador.Excel(r)
		a.Nombre, a.TipoMime = nombre+".xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	}
	if err != nil {
		return nil, err
	}

	actor := s.ActorID
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actor, Accion: auditoria.AccionGenerarReporte, Entidad: auditoria.EntidadReporte, EntidadID: string(s.Tipo),
		Detalle: map[string]any{"tipo": string(s.Tipo), "formato": string(s.Formato), "filas": len(r.Filas), "filtros": r.Filtros},
		Fecha:   r.EmitidoEn,
	})
	return a, nil
}

// --- Contenido de cada reporte ------------------------------------------------------------

func (c *GenerarReporte) tramites(ctx context.Context, r *indicador.Reporte, f indicador.FiltroReporte) error {
	filas, err := c.Indicadores.ListarTramites(ctx, f)
	if err != nil {
		return err
	}
	r.Titulo = "Reporte de trámites"
	r.Columnas = []indicador.Columna{
		{Titulo: "Código", Tipo: indicador.ColTexto, Ancho: 15}, {Titulo: "Cód. docente", Tipo: indicador.ColTexto, Ancho: 9},
		{Titulo: "Docente", Tipo: indicador.ColTexto, Ancho: 26}, {Titulo: "Tipo", Tipo: indicador.ColTexto, Ancho: 13},
		{Titulo: "Categoría", Tipo: indicador.ColTexto, Ancho: 14}, {Titulo: "Estado", Tipo: indicador.ColTexto, Ancho: 16},
		{Titulo: "Origen", Tipo: indicador.ColTexto, Ancho: 9}, {Titulo: "Analista", Tipo: indicador.ColTexto, Ancho: 20},
		{Titulo: "Resultado", Tipo: indicador.ColTexto, Ancho: 11}, {Titulo: "Radicado", Tipo: indicador.ColFecha, Ancho: 10},
	}
	var enCurso, concluidos int
	for _, t := range filas {
		r.Filas = append(r.Filas, []any{
			t.Codigo, t.DocenteCodigo, t.DocenteNombre, etiquetaTipo[t.Tipo], t.Categoria, etiquetaEstado[t.Estado],
			etiquetaOrigen[t.Origen], t.Analista, etiquetaResultado[t.Resultado], t.RadicadoEn.In(zonaBolivia),
		})
		if t.Estado.Cerrado() {
			concluidos++
		} else {
			enCurso++
		}
	}
	r.Resumen = []indicador.Dato{
		{Etiqueta: "Total de trámites", Valor: fmt.Sprint(len(filas))},
		{Etiqueta: "En curso", Valor: fmt.Sprint(enCurso)},
		{Etiqueta: "Remitidos o archivados", Valor: fmt.Sprint(concluidos)},
	}
	r.Filtros = describirPeriodo("Radicados", f.Periodo)
	if f.Estado != "" {
		r.Filtros = append(r.Filtros, "Estado: "+etiquetaEstado[f.Estado])
	}
	r.Filtros = append(r.Filtros, describirTipoCategoria(f, primeraCategoria(len(filas), func() string { return filas[0].Categoria }))...)
	return nil
}

func (c *GenerarReporte) inventario(ctx context.Context, r *indicador.Reporte, f indicador.FiltroReporte) error {
	filas, err := c.Indicadores.ListarInventario(ctx, f)
	if err != nil {
		return err
	}
	maximo := diasMaximos(c.DiasMaximosPrestamo)
	hoy := r.EmitidoEn.In(zonaBolivia)
	r.Titulo = "Inventario de legajos"
	r.Columnas = []indicador.Columna{
		{Titulo: "Legajo", Tipo: indicador.ColTexto, Ancho: 11}, {Titulo: "Cód. docente", Tipo: indicador.ColTexto, Ancho: 9},
		{Titulo: "Docente", Tipo: indicador.ColTexto, Ancho: 26}, {Titulo: "Gaveta", Tipo: indicador.ColTexto, Ancho: 7},
		{Titulo: "Orden", Tipo: indicador.ColEntero, Ancho: 6}, {Titulo: "Anexos", Tipo: indicador.ColTexto, Ancho: 12},
		{Titulo: "Estado", Tipo: indicador.ColTexto, Ancho: 10}, {Titulo: "Custodio", Tipo: indicador.ColTexto, Ancho: 20},
		{Titulo: "Salida", Tipo: indicador.ColFecha, Ancho: 10}, {Titulo: "Fecha límite", Tipo: indicador.ColFecha, Ancho: 10},
		{Titulo: "Días de mora", Tipo: indicador.ColEntero, Ancho: 8},
	}
	var prestados, enMora int
	for _, l := range filas {
		var salida, limite, mora any
		if l.Estado == legajo.Prestado {
			prestados++
		}
		if l.Salida != nil {
			salida = l.Salida.In(zonaBolivia)
			if l.FechaLimite != nil {
				limite = *l.FechaLimite // fecha sin hora: se muestra tal cual
			}
			if a, ok := indicador.EvaluarMora(indicador.Prestamo{Salida: *l.Salida, FechaLimite: l.FechaLimite}, hoy, maximo); ok {
				mora = a.DiasMora
				enMora++
			}
		}
		r.Filas = append(r.Filas, []any{
			l.LegajoCodigo, l.DocenteCodigo, l.DocenteNombre, l.Gaveta, l.Orden, l.Anexos, etiquetaLegajo[l.Estado],
			l.Custodio, salida, limite, mora,
		})
	}
	r.Resumen = []indicador.Dato{
		{Etiqueta: "Total de legajos", Valor: fmt.Sprint(len(filas))},
		{Etiqueta: "En archivo", Valor: fmt.Sprint(len(filas) - prestados)},
		{Etiqueta: "Prestados", Valor: fmt.Sprint(prestados)},
		{Etiqueta: "En mora", Valor: fmt.Sprintf("%d (fecha límite vencida o más de %d días fuera)", enMora, maximo)},
	}
	if f.EstadoLegajo != "" {
		r.Filtros = append(r.Filtros, "Estado: "+etiquetaLegajo[f.EstadoLegajo])
	}
	return nil
}

func (c *GenerarReporte) tiempos(ctx context.Context, r *indicador.Reporte, f indicador.FiltroReporte) error {
	resueltos, err := c.Indicadores.ListarResueltos(ctx, f)
	if err != nil {
		return err
	}
	r.Titulo = "Tiempos de atención"
	r.Columnas = []indicador.Columna{
		{Titulo: "Código", Tipo: indicador.ColTexto, Ancho: 15}, {Titulo: "Docente", Tipo: indicador.ColTexto, Ancho: 26},
		{Titulo: "Tipo", Tipo: indicador.ColTexto, Ancho: 13}, {Titulo: "Categoría", Tipo: indicador.ColTexto, Ancho: 14},
		{Titulo: "Analista", Tipo: indicador.ColTexto, Ancho: 20}, {Titulo: "Resultado", Tipo: indicador.ColTexto, Ancho: 11},
		{Titulo: "Radicado", Tipo: indicador.ColFecha, Ancho: 10}, {Titulo: "Resuelto", Tipo: indicador.ColFecha, Ancho: 10},
		{Titulo: "Días", Tipo: indicador.ColDecimal, Ancho: 7},
	}
	for _, t := range resueltos {
		r.Filas = append(r.Filas, []any{
			t.Codigo, t.DocenteNombre, etiquetaTipo[t.Tipo], t.Categoria, t.Analista, etiquetaResultado[t.Resultado],
			t.RadicadoEn.In(zonaBolivia), t.ResueltoEn.In(zonaBolivia), t.Dias(),
		})
	}
	if len(resueltos) > 0 {
		grupos := indicador.AgruparTiempos(resueltos)
		minimo, maximo := grupos[0].MinimoDias, grupos[0].MaximoDias
		for _, g := range grupos {
			minimo, maximo = min(minimo, g.MinimoDias), max(maximo, g.MaximoDias)
		}
		r.Resumen = []indicador.Dato{
			{Etiqueta: "Trámites resueltos", Valor: fmt.Sprint(len(resueltos))},
			{Etiqueta: "Tiempo promedio", Valor: dias(indicador.PromedioDias(resueltos))},
			{Etiqueta: "Más rápido / más lento", Valor: dias(minimo) + " / " + dias(maximo)},
		}
		for _, g := range grupos {
			r.Resumen = append(r.Resumen, indicador.Dato{
				Etiqueta: etiquetaTipo[g.Tipo] + " · " + g.Categoria,
				Valor:    fmt.Sprintf("%s en promedio (%d %s)", dias(g.PromedioDias), g.Cantidad, plural(g.Cantidad, "trámite", "trámites")),
			})
		}
	}
	r.Filtros = describirPeriodo("Resueltos", f.Periodo)
	r.Filtros = append(r.Filtros, describirTipoCategoria(f, primeraCategoria(len(resueltos), func() string { return resueltos[0].Categoria }))...)
	return nil
}

// --- Descripción legible de los filtros y formatos -----------------------------------------

func fecha(t time.Time) string { return t.In(zonaBolivia).Format("02/01/2006") }

// describirPeriodo traduce el periodo (Hasta exclusivo) a "Radicados del 01/01/2026 al 31/01/2026".
func describirPeriodo(sujeto string, p indicador.Periodo) []string {
	switch {
	case p.Desde != nil && p.Hasta != nil:
		return []string{fmt.Sprintf("%s del %s al %s", sujeto, fecha(*p.Desde), fecha(p.Hasta.Add(-time.Second)))}
	case p.Desde != nil:
		return []string{fmt.Sprintf("%s desde el %s", sujeto, fecha(*p.Desde))}
	case p.Hasta != nil:
		return []string{fmt.Sprintf("%s hasta el %s", sujeto, fecha(p.Hasta.Add(-time.Second)))}
	}
	return []string{sujeto + " en todo el historial"}
}

// primeraCategoria da el nombre de la categoría filtrada a partir de la primera fila (todas coinciden).
func primeraCategoria(n int, nombre func() string) string {
	if n == 0 {
		return ""
	}
	return nombre()
}

func describirTipoCategoria(f indicador.FiltroReporte, categoria string) []string {
	var out []string
	if f.Tipo != "" {
		out = append(out, "Tipo: "+etiquetaTipo[f.Tipo])
	}
	if f.CategoriaID > 0 && categoria != "" {
		out = append(out, "Categoría: "+categoria)
	}
	return out
}

// dias formatea con un decimal y coma decimal: 12,5 días.
func dias(d float64) string {
	return strings.Replace(fmt.Sprintf("%.1f", d), ".", ",", 1) + " días"
}

func plural(n int, uno, varios string) string {
	if n == 1 {
		return uno
	}
	return varios
}
