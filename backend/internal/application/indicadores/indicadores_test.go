package indicadores_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"escalafon/backend/internal/application/indicadores"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/domain/usuario"
)

var (
	bolivia = time.FixedZone("BOT", -4*60*60)
	ahora   = time.Date(2026, 10, 25, 10, 0, 0, 0, bolivia)
)

func ptr[T any](v T) *T { return &v }

func dia(y int, m time.Month, d int) *time.Time { return ptr(time.Date(y, m, d, 0, 0, 0, 0, time.UTC)) }

// --- CU-25 ---------------------------------------------------------------------------

func TestTablero_DistribucionCompletaEnOrdenDelFlujo(t *testing.T) {
	repo := &testdobles.IndicadorRepoFalso{Conteos: []indicador.ConteoEstado{
		{Estado: tramite.Archivado, Cantidad: 4}, {Estado: tramite.EnRevision, Cantidad: 3},
		{Estado: tramite.Ingresado, Cantidad: 2}, {Estado: tramite.Remitido, Cantidad: 1},
	}}
	caso := &indicadores.ConsultarTablero{Indicadores: repo, Reloj: &testdobles.RelojFijo{Momento: ahora}}
	tab, err := caso.Ejecutar(context.Background(), indicador.Periodo{})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if len(tab.PorEstado) != 9 || tab.PorEstado[0].Estado != tramite.Ingresado || tab.PorEstado[8].Estado != tramite.Archivado {
		t.Fatalf("esperaba los nueve estados en el orden del flujo: %+v", tab.PorEstado)
	}
	if tab.PorEstado[1].Cantidad != 0 || tab.PorEstado[2].Cantidad != 3 {
		t.Fatalf("los estados sin trámites aparecen con 0: %+v", tab.PorEstado)
	}
	if tab.Radicados != 10 || tab.EnCurso != 5 {
		t.Fatalf("remitidos y archivados no están en curso: radicados=%d enCurso=%d", tab.Radicados, tab.EnCurso)
	}
	if tab.DiasMaximosPrestamo != indicador.DiasMaximosPorDefecto || !tab.GeneradoEn.Equal(ahora) {
		t.Fatalf("sin umbral configurado se usan %d días: %+v", indicador.DiasMaximosPorDefecto, tab)
	}
}

func TestTablero_TiemposDeResolucion(t *testing.T) {
	radicado := time.Date(2026, 9, 1, 9, 0, 0, 0, bolivia)
	repo := &testdobles.IndicadorRepoFalso{Resueltos: []indicador.Resuelto{
		{Tipo: tramite.Categorizacion, Categoria: "Titular A", RadicadoEn: radicado, ResueltoEn: radicado.AddDate(0, 0, 10)},
		{Tipo: tramite.Categorizacion, Categoria: "Titular A", RadicadoEn: radicado, ResueltoEn: radicado.AddDate(0, 0, 20)},
		{Tipo: tramite.Regularizacion, Categoria: "Adjunto", RadicadoEn: radicado, ResueltoEn: radicado.AddDate(0, 0, 30)},
	}}
	caso := &indicadores.ConsultarTablero{Indicadores: repo, Reloj: &testdobles.RelojFijo{Momento: ahora}}
	tab, _ := caso.Ejecutar(context.Background(), indicador.Periodo{})
	if tab.Resueltos != 3 || tab.PromedioDias != 20 || len(tab.Tiempos) != 2 || tab.Tiempos[0].PromedioDias != 15 {
		t.Fatalf("tiempos inesperados: resueltos=%d promedio=%v tiempos=%+v", tab.Resueltos, tab.PromedioDias, tab.Tiempos)
	}
}

func TestTablero_AlertasDeMora(t *testing.T) {
	repo := &testdobles.IndicadorRepoFalso{Prestamos: []indicador.Prestamo{
		{LegajoCodigo: "LEG-1", Salida: time.Date(2026, 10, 20, 9, 0, 0, 0, bolivia), FechaLimite: dia(2026, 10, 27)}, // al día
		{LegajoCodigo: "LEG-2", Salida: time.Date(2026, 10, 18, 9, 0, 0, 0, bolivia), FechaLimite: dia(2026, 10, 22)}, // 3 días vencida
		{LegajoCodigo: "LEG-3", Salida: time.Date(2026, 10, 1, 9, 0, 0, 0, bolivia), FechaLimite: dia(2026, 11, 30)},  // 24 días fuera
	}}
	caso := &indicadores.ConsultarTablero{Indicadores: repo, Reloj: &testdobles.RelojFijo{Momento: ahora}, DiasMaximosPrestamo: 10}
	tab, _ := caso.Ejecutar(context.Background(), indicador.Periodo{})
	if tab.PrestamosVigentes != 3 || len(tab.Alertas) != 2 {
		t.Fatalf("esperaba 3 préstamos y 2 alertas: %+v", tab.Alertas)
	}
	if a := tab.Alertas[0]; a.LegajoCodigo != "LEG-3" || a.DiasMora != 14 || a.Motivo != indicador.MoraExcedeMaximo {
		t.Fatalf("primero la de más mora (24 días con máximo 10): %+v", a)
	}
	if a := tab.Alertas[1]; a.LegajoCodigo != "LEG-2" || a.DiasMora != 3 || a.Motivo != indicador.MoraFechaLimite {
		t.Fatalf("la de fecha límite vencida: %+v", a)
	}

	// Al registrarse la devolución la carpeta ya no está prestada: la alerta desaparece sola (4a).
	repo.Prestamos = repo.Prestamos[:2]
	tab, _ = caso.Ejecutar(context.Background(), indicador.Periodo{})
	if len(tab.Alertas) != 1 || tab.Alertas[0].LegajoCodigo != "LEG-2" {
		t.Fatalf("la alerta de la carpeta devuelta debe cerrarse: %+v", tab.Alertas)
	}
}

func TestTablero_RechazaRangoInvertido(t *testing.T) {
	caso := &indicadores.ConsultarTablero{Indicadores: &testdobles.IndicadorRepoFalso{}, Reloj: &testdobles.RelojFijo{Momento: ahora}}
	_, err := caso.Ejecutar(context.Background(), indicador.Periodo{Desde: dia(2026, 10, 5), Hasta: dia(2026, 10, 1)})
	if !errors.Is(err, indicador.ErrRangoFechas) {
		t.Fatalf("esperaba ErrRangoFechas, obtuve %v", err)
	}
}

// --- CU-26 ---------------------------------------------------------------------------

type entorno struct {
	caso       *indicadores.GenerarReporte
	repo       *testdobles.IndicadorRepoFalso
	exportador *testdobles.ExportadorFalso
	auditoria  *testdobles.AuditoriaRepoFalso
	director   *usuario.Usuario
}

func nuevoEntorno() *entorno {
	usuarios := testdobles.NewUsuarioRepoFalso()
	e := &entorno{
		repo: &testdobles.IndicadorRepoFalso{}, exportador: &testdobles.ExportadorFalso{}, auditoria: &testdobles.AuditoriaRepoFalso{},
		director: usuarios.Sembrar(usuario.Usuario{Nombres: "Rosa", Apellidos: "Director", Correo: "director@uagrm.edu.bo", Activo: true}),
	}
	e.caso = &indicadores.GenerarReporte{
		Indicadores: e.repo, Usuarios: usuarios, Exportador: e.exportador, Auditoria: e.auditoria,
		Reloj: &testdobles.RelojFijo{Momento: ahora}, DiasMaximosPrestamo: 15,
	}
	return e
}

func (e *entorno) sembrarTramites() {
	radicado := time.Date(2026, 10, 1, 9, 0, 0, 0, bolivia)
	e.repo.Tramites = []indicador.FilaTramite{
		{Codigo: "ED-2026-AAAAAA", DocenteNombre: "Ana Pérez", Tipo: tramite.Categorizacion, Categoria: "Titular A", Estado: tramite.EnRevision, Origen: tramite.OrigenPortal, RadicadoEn: radicado},
		{Codigo: "ED-2026-BBBBBB", DocenteNombre: "Luis Rojas", Tipo: tramite.Regularizacion, Categoria: "Adjunto", Estado: tramite.Archivado, Origen: tramite.OrigenStaff, Resultado: tramite.Favorable, RadicadoEn: radicado},
	}
}

func TestReporte_TramitesEnPDFQuedaAuditado(t *testing.T) {
	e := nuevoEntorno()
	e.sembrarTramites()
	a, err := e.caso.Ejecutar(context.Background(), indicadores.SolicitudReporte{
		ActorID: e.director.ID, Tipo: indicador.ReporteTramites, Formato: indicador.FormatoPDF,
		Filtro: indicador.FiltroReporte{Periodo: indicador.Periodo{Desde: ptr(time.Date(2026, 10, 1, 0, 0, 0, 0, bolivia)), Hasta: ptr(time.Date(2026, 11, 1, 0, 0, 0, 0, bolivia))}},
	})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if a.Nombre != "reporte-tramites-2026-10-25.pdf" || a.TipoMime != "application/pdf" || a.Filas != 2 || string(a.Datos) != "%PDF-falso" {
		t.Fatalf("archivo inesperado: %+v", a)
	}
	r := e.exportador.Ultimo
	if r.Titulo != "Reporte de trámites" || r.EmitidoPor != "Rosa Director" || len(r.Columnas) != len(r.Filas[0]) {
		t.Fatalf("reporte inesperado: %+v", r)
	}
	if r.Filas[0][5] != "En revisión" || r.Filas[1][6] != "Ventanilla" || r.Filas[1][8] != "Favorable" {
		t.Fatalf("las celdas usan las etiquetas en español: %+v", r.Filas)
	}
	if len(r.Filtros) != 1 || r.Filtros[0] != "Radicados del 01/10/2026 al 31/10/2026" {
		t.Fatalf("el periodo se describe con el hasta inclusive: %+v", r.Filtros)
	}
	if r.Resumen[0].Valor != "2" || r.Resumen[1].Valor != "1" || r.Resumen[2].Valor != "1" {
		t.Fatalf("resumen inesperado: %+v", r.Resumen)
	}

	if len(e.auditoria.Eventos) != 1 {
		t.Fatalf("la emisión debe quedar en auditoría")
	}
	ev := e.auditoria.Eventos[0]
	if ev.Accion != auditoria.AccionGenerarReporte || *ev.UsuarioID != e.director.ID || ev.Detalle["formato"] != "PDF" || ev.Detalle["filas"] != 2 {
		t.Fatalf("evento inesperado: %+v", ev)
	}
}

func TestReporte_InventarioEnExcelMarcaLaMora(t *testing.T) {
	e := nuevoEntorno()
	e.repo.Inventario = []indicador.FilaInventario{
		{LegajoCodigo: "LEG-00001", Estado: legajo.EnArchivo, Gaveta: "G1", Orden: 1},
		{LegajoCodigo: "LEG-00002", Estado: legajo.Prestado, Custodio: "Ana Analista", Salida: ptr(time.Date(2026, 10, 1, 9, 0, 0, 0, bolivia)), FechaLimite: dia(2026, 10, 8)},
	}
	a, err := e.caso.Ejecutar(context.Background(), indicadores.SolicitudReporte{ActorID: e.director.ID, Tipo: indicador.ReporteInventario, Formato: indicador.FormatoExcel})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if !strings.HasSuffix(a.Nombre, ".xlsx") || !strings.Contains(a.TipoMime, "spreadsheetml") {
		t.Fatalf("esperaba un Excel: %+v", a)
	}
	r := e.exportador.Ultimo
	if r.Filas[0][10] != nil || r.Filas[1][10] != 17 {
		t.Fatalf("solo la carpeta prestada y vencida tiene días de mora (17): %v / %v", r.Filas[0][10], r.Filas[1][10])
	}
	if r.Resumen[2].Valor != "1" || !strings.HasPrefix(r.Resumen[3].Valor, "1 ") {
		t.Fatalf("resumen inesperado: %+v", r.Resumen)
	}
}

func TestReporte_TiemposConResumenPorCategoria(t *testing.T) {
	e := nuevoEntorno()
	radicado := time.Date(2026, 9, 1, 9, 0, 0, 0, bolivia)
	e.repo.Resueltos = []indicador.Resuelto{
		{Codigo: "ED-1", Tipo: tramite.Categorizacion, CategoriaID: 3, Categoria: "Titular A", RadicadoEn: radicado, ResueltoEn: radicado.AddDate(0, 0, 12)},
		{Codigo: "ED-2", Tipo: tramite.Regularizacion, CategoriaID: 5, Categoria: "Adjunto", RadicadoEn: radicado, ResueltoEn: radicado.AddDate(0, 0, 40)},
	}
	_, err := e.caso.Ejecutar(context.Background(), indicadores.SolicitudReporte{
		ActorID: e.director.ID, Tipo: indicador.ReporteTiempos, Formato: indicador.FormatoPDF,
		Filtro: indicador.FiltroReporte{CategoriaID: 3},
	})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	r := e.exportador.Ultimo
	if len(r.Filas) != 1 || r.Filas[0][8] != 12.0 {
		t.Fatalf("el filtro por categoría deja un trámite de 12 días: %+v", r.Filas)
	}
	if r.Resumen[1].Valor != "12,0 días" || r.Resumen[3].Etiqueta != "Categorización · Titular A" {
		t.Fatalf("resumen inesperado: %+v", r.Resumen)
	}
	if r.Filtros[len(r.Filtros)-1] != "Categoría: Titular A" {
		t.Fatalf("la categoría filtrada se nombra en los filtros: %+v", r.Filtros)
	}
}

func TestReporte_SinDatosNoSeEmite(t *testing.T) {
	e := nuevoEntorno()
	e.sembrarTramites()
	_, err := e.caso.Ejecutar(context.Background(), indicadores.SolicitudReporte{
		ActorID: e.director.ID, Tipo: indicador.ReporteTramites, Formato: indicador.FormatoExcel,
		Filtro: indicador.FiltroReporte{Estado: tramite.Observado},
	})
	if !errors.Is(err, indicador.ErrReporteVacio) {
		t.Fatalf("esperaba ErrReporteVacio, obtuve %v", err)
	}
	if e.exportador.Ultimo != nil || len(e.auditoria.Eventos) != 0 {
		t.Fatal("un reporte vacío no se compone ni se audita")
	}
}

func TestReporte_ValidaLaSolicitud(t *testing.T) {
	e := nuevoEntorno()
	e.sembrarTramites()
	casos := []struct {
		nombre string
		s      indicadores.SolicitudReporte
		err    error
	}{
		{"tipo desconocido", indicadores.SolicitudReporte{Tipo: "DOCENTES", Formato: indicador.FormatoPDF}, indicador.ErrTipoReporte},
		{"formato desconocido", indicadores.SolicitudReporte{Tipo: indicador.ReporteTramites, Formato: "CSV"}, indicador.ErrFormato},
		{"rango invertido", indicadores.SolicitudReporte{Tipo: indicador.ReporteTramites, Formato: indicador.FormatoPDF,
			Filtro: indicador.FiltroReporte{Periodo: indicador.Periodo{Desde: dia(2026, 10, 9), Hasta: dia(2026, 10, 1)}}}, indicador.ErrRangoFechas},
		{"estado inexistente", indicadores.SolicitudReporte{Tipo: indicador.ReporteTramites, Formato: indicador.FormatoPDF,
			Filtro: indicador.FiltroReporte{Estado: "BORRADO"}}, indicador.ErrFiltroInvalido},
		{"estado de legajo inexistente", indicadores.SolicitudReporte{Tipo: indicador.ReporteInventario, Formato: indicador.FormatoPDF,
			Filtro: indicador.FiltroReporte{EstadoLegajo: "PERDIDO"}}, indicador.ErrFiltroInvalido},
	}
	for _, c := range casos {
		c.s.ActorID = e.director.ID
		if _, err := e.caso.Ejecutar(context.Background(), c.s); !errors.Is(err, c.err) {
			t.Fatalf("%s: esperaba %v, obtuve %v", c.nombre, c.err, err)
		}
	}
}
