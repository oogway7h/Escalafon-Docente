// Package indicador modela el tablero de indicadores (CU-25) y los reportes exportables (CU-26).
// Todo es de solo lectura: se calcula a partir de los trámites, su historial, los legajos y los
// movimientos de custodia que ya registran los demás módulos.
package indicador

import (
	"context"
	"sort"
	"time"

	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
)

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrRangoFechas    erroDominio = "la fecha inicial es posterior a la final"
	ErrTipoReporte    erroDominio = "el tipo de reporte no existe"
	ErrFormato        erroDominio = "el formato debe ser PDF o XLSX"
	ErrFiltroInvalido erroDominio = "uno de los filtros no tiene un valor válido"
	ErrReporteVacio   erroDominio = "no hay información para los filtros elegidos"

	ErrInterpreteNoDisponible erroDominio = "el intérprete de lenguaje natural no está disponible"
	ErrNoEntendido            erroDominio = "no se pudo entender qué reporte se pide"
)

// DiasMaximosPorDefecto es el umbral de préstamo si DIAS_MAXIMOS_PRESTAMO no está configurado.
const DiasMaximosPorDefecto = 15

// Periodo acota las consultas. Desde es inclusivo y Hasta exclusivo; nil = sin límite.
type Periodo struct {
	Desde *time.Time
	Hasta *time.Time
}

// Valido indica si el rango no está invertido.
func (p Periodo) Valido() bool {
	return p.Desde == nil || p.Hasta == nil || p.Desde.Before(*p.Hasta)
}

// ConteoEstado es la cantidad de trámites radicados en el periodo que hoy están en un estado.
type ConteoEstado struct {
	Estado   tramite.Estado
	Cantidad int64
}

// Resuelto es un trámite que llegó a Remitido: desde ahí el docente conoce el resultado. En datos
// anteriores al flujo actual (sin Remitido) cuenta el paso a Archivado.
type Resuelto struct {
	Codigo        string
	DocenteCodigo string
	DocenteNombre string
	Tipo          tramite.Tipo
	CategoriaID   int16
	Categoria     string
	Analista      string
	Resultado     tramite.Resultado
	RadicadoEn    time.Time
	ResueltoEn    time.Time
}

// Dias es el tiempo de resolución en días (con fracción).
func (r Resuelto) Dias() float64 { return r.ResueltoEn.Sub(r.RadicadoEn).Hours() / 24 }

// TiempoResolucion resume los resueltos de un tipo y una categoría.
type TiempoResolucion struct {
	Tipo         tramite.Tipo
	Categoria    string
	Cantidad     int
	PromedioDias float64
	MinimoDias   float64
	MaximoDias   float64
}

// AgruparTiempos agrupa los resueltos por tipo y categoría, ordenados por tipo y luego por categoría.
func AgruparTiempos(resueltos []Resuelto) []TiempoResolucion {
	type clave struct {
		tipo      tramite.Tipo
		categoria string
	}
	grupos := map[clave]*TiempoResolucion{}
	var orden []clave
	for _, r := range resueltos {
		k := clave{r.Tipo, r.Categoria}
		g, ok := grupos[k]
		d := r.Dias()
		if !ok {
			g = &TiempoResolucion{Tipo: r.Tipo, Categoria: r.Categoria, MinimoDias: d, MaximoDias: d}
			grupos[k] = g
			orden = append(orden, k)
		}
		g.Cantidad++
		g.PromedioDias += d // suma; se divide al final
		g.MinimoDias = min(g.MinimoDias, d)
		g.MaximoDias = max(g.MaximoDias, d)
	}
	sort.Slice(orden, func(i, j int) bool {
		if orden[i].tipo != orden[j].tipo {
			return orden[i].tipo < orden[j].tipo
		}
		return orden[i].categoria < orden[j].categoria
	})
	out := make([]TiempoResolucion, len(orden))
	for i, k := range orden {
		g := grupos[k]
		g.PromedioDias /= float64(g.Cantidad)
		out[i] = *g
	}
	return out
}

// PromedioDias es el tiempo medio de resolución de todos los resueltos (0 si no hay ninguno).
func PromedioDias(resueltos []Resuelto) float64 {
	if len(resueltos) == 0 {
		return 0
	}
	var suma float64
	for _, r := range resueltos {
		suma += r.Dias()
	}
	return suma / float64(len(resueltos))
}

// CargaAnalista es el trabajo de un analista: lo que tiene a cargo hoy y lo que resolvió en el periodo.
type CargaAnalista struct {
	ID         string
	Nombre     string
	Activos    int64 // en etapas que cuentan como carga vigente (tramite.Estado.Activo)
	EnFacultad int64 // fuera del Departamento: no cuentan como carga
	Resueltos  int64 // remitidos dentro del periodo
}

// Prestamo es una carpeta que hoy está fuera del archivo.
type Prestamo struct {
	LegajoID      int32
	LegajoCodigo  string
	DocenteNombre string
	Custodio      string
	TramiteCodigo string
	Salida        time.Time
	FechaLimite   *time.Time // fecha (sin hora) acordada al despachar
}

// MotivoMora explica por qué un préstamo está en mora.
type MotivoMora string

const (
	MoraFechaLimite  MotivoMora = "FECHA_LIMITE_VENCIDA" // pasó la fecha límite acordada al despachar
	MoraExcedeMaximo MotivoMora = "EXCEDE_MAXIMO"        // lleva más días fuera que el máximo configurado
)

// AlertaMora es un préstamo en mora. Se calcula al consultar: al registrarse la devolución la carpeta
// deja de estar prestada y la alerta desaparece sola (CU-25, 4a).
type AlertaMora struct {
	Prestamo
	DiasPrestado int
	DiasMora     int
	Motivo       MotivoMora
}

// EvaluarMora compara un préstamo con su fecha límite y con el máximo de días de préstamo. hoy fija
// el día y la zona horaria del cálculo: solo cuentan días calendario completos.
func EvaluarMora(p Prestamo, hoy time.Time, diasMaximos int) (AlertaMora, bool) {
	dia := func(t time.Time) time.Time {
		y, m, d := t.Date()
		return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	}
	dias := func(desde, hasta time.Time) int { return int(dia(hasta).Sub(dia(desde)).Hours() / 24) }

	a := AlertaMora{Prestamo: p, DiasPrestado: dias(p.Salida.In(hoy.Location()), hoy)}
	if p.FechaLimite != nil {
		if vencida := dias(*p.FechaLimite, hoy); vencida > 0 {
			a.DiasMora, a.Motivo = vencida, MoraFechaLimite
		}
	}
	if diasMaximos > 0 {
		if excede := a.DiasPrestado - diasMaximos; excede > a.DiasMora {
			a.DiasMora, a.Motivo = excede, MoraExcedeMaximo
		}
	}
	return a, a.DiasMora > 0
}

// OrdenarAlertas deja primero las de más días de mora.
func OrdenarAlertas(alertas []AlertaMora) {
	sort.SliceStable(alertas, func(i, j int) bool {
		if alertas[i].DiasMora != alertas[j].DiasMora {
			return alertas[i].DiasMora > alertas[j].DiasMora
		}
		return alertas[i].LegajoCodigo < alertas[j].LegajoCodigo
	})
}

// --- CU-26: reportes ----------------------------------------------------------------

type TipoReporte string

const (
	ReporteTramites   TipoReporte = "TRAMITES"   // trámites radicados, con su estado actual
	ReporteInventario TipoReporte = "INVENTARIO" // legajos físicos con su ubicación y préstamo vigente
	ReporteTiempos    TipoReporte = "TIEMPOS"    // tiempos de atención de los trámites resueltos
)

func (t TipoReporte) Valido() bool {
	return t == ReporteTramites || t == ReporteInventario || t == ReporteTiempos
}

type Formato string

const (
	FormatoPDF   Formato = "PDF"
	FormatoExcel Formato = "XLSX"
)

func (f Formato) Valido() bool { return f == FormatoPDF || f == FormatoExcel }

// FiltroReporte reúne los filtros de los tres reportes; cada uno usa los que le corresponden.
type FiltroReporte struct {
	Periodo                     // TRAMITES: fecha de radicación; TIEMPOS: fecha de resolución
	Estado       tramite.Estado // TRAMITES
	Tipo         tramite.Tipo   // TRAMITES y TIEMPOS
	CategoriaID  int16          // TRAMITES y TIEMPOS
	EstadoLegajo legajo.Estado  // INVENTARIO
}

// FilaTramite es un trámite del reporte de trámites.
type FilaTramite struct {
	Codigo        string
	DocenteCodigo string
	DocenteNombre string
	Tipo          tramite.Tipo
	Categoria     string
	Estado        tramite.Estado
	Origen        tramite.Origen
	Analista      string
	Resultado     tramite.Resultado
	RadicadoEn    time.Time
}

// FilaInventario es un legajo del reporte de inventario; Salida y Custodio solo si está prestado.
type FilaInventario struct {
	LegajoCodigo  string
	DocenteCodigo string
	DocenteNombre string
	Gaveta        string
	Orden         int
	Anexos        string
	Estado        legajo.Estado
	Custodio      string
	Salida        *time.Time
	FechaLimite   *time.Time
}

// TipoColumna indica cómo mostrar y exportar los valores de una columna.
type TipoColumna string

const (
	ColTexto   TipoColumna = "TEXTO"
	ColEntero  TipoColumna = "ENTERO"
	ColDecimal TipoColumna = "DECIMAL" // un decimal
	ColFecha   TipoColumna = "FECHA"   // time.Time, solo el día (hora de Bolivia)
)

type Columna struct {
	Titulo string
	Tipo   TipoColumna
	Ancho  float64 // peso relativo del ancho en el PDF y caracteres aproximados en Excel
}

// Dato es una línea de resumen del reporte ("Total de trámites: 42").
type Dato struct{ Etiqueta, Valor string }

// Reporte es el contenido ya calculado, independiente del formato de salida. Cada celda es string,
// int, int64, float64, time.Time o nil (vacío), según el tipo de su columna.
type Reporte struct {
	Tipo       TipoReporte
	Titulo     string
	Filtros    []string // descripción legible de los filtros aplicados
	Resumen    []Dato
	Columnas   []Columna
	Filas      [][]any
	EmitidoPor string
	EmitidoEn  time.Time
}

// Consulta es el puerto de lectura de los indicadores: no expone escritura.
type Consulta interface {
	// ContarPorEstado agrupa por estado actual los trámites radicados en el periodo.
	ContarPorEstado(ctx context.Context, p Periodo) ([]ConteoEstado, error)
	// ListarResueltos devuelve los trámites resueltos dentro del periodo (por fecha de resolución),
	// del más reciente al más antiguo; Tipo y CategoriaID filtran si no están vacíos.
	ListarResueltos(ctx context.Context, f FiltroReporte) ([]Resuelto, error)
	// CargaAnalistas devuelve a los usuarios activos con TRAMITE_REVISAR, con su carga vigente y
	// los trámites que resolvieron en el periodo.
	CargaAnalistas(ctx context.Context, p Periodo) ([]CargaAnalista, error)
	// ListarPrestamos devuelve las carpetas que hoy están fuera del archivo.
	ListarPrestamos(ctx context.Context) ([]Prestamo, error)

	ListarTramites(ctx context.Context, f FiltroReporte) ([]FilaTramite, error)
	ListarInventario(ctx context.Context, f FiltroReporte) ([]FilaInventario, error)
}
