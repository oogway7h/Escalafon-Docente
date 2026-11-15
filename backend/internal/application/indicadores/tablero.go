// Package indicadores implementa CU-25 (visualizar tablero de indicadores) y CU-26 (generar y
// exportar reportes). Ambos leen lo que registran los demás módulos; el reporte deja constancia
// en la auditoría.
package indicadores

import (
	"context"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/tramite"
)

// zonaBolivia fija el día de los cálculos de mora y las fechas de los reportes (UTC-4, sin horario de verano).
var zonaBolivia = time.FixedZone("BOT", -4*60*60)

// estadosTramite son los estados en el orden del flujo, para presentar la distribución completa.
var estadosTramite = []tramite.Estado{
	tramite.Ingresado, tramite.Asignado, tramite.EnRevision, tramite.Observado, tramite.Dictaminado,
	tramite.EnFacultad, tramite.RevisionFinal, tramite.Remitido, tramite.Archivado,
}

func estadoValido(e tramite.Estado) bool {
	for _, v := range estadosTramite {
		if v == e {
			return true
		}
	}
	return false
}

// Tablero es la vista consolidada del Director: trámites del periodo, tiempos de resolución,
// carga por analista y alertas de carpetas en mora.
type Tablero struct {
	GeneradoEn          time.Time
	DiasMaximosPrestamo int

	PorEstado    []indicador.ConteoEstado // los nueve estados en el orden del flujo, con 0 si no hay
	Radicados    int64                    // radicados en el periodo
	EnCurso      int64                    // de esos, los que aún no se remitieron ni archivaron
	Resueltos    int                      // resueltos (remitidos) en el periodo
	PromedioDias float64                  // tiempo medio de resolución de los resueltos
	Tiempos      []indicador.TiempoResolucion
	Analistas    []indicador.CargaAnalista

	PrestamosVigentes int
	Alertas           []indicador.AlertaMora // más días de mora primero
}

// ConsultarTablero calcula los indicadores en cada consulta: siempre están al día y una alerta de
// mora se cierra sola cuando se registra la devolución de la carpeta.
type ConsultarTablero struct {
	Indicadores         indicador.Consulta
	Reloj               ports.Clock
	DiasMaximosPrestamo int // umbral de préstamo (DIAS_MAXIMOS_PRESTAMO)
}

func diasMaximos(n int) int {
	if n < 1 {
		return indicador.DiasMaximosPorDefecto
	}
	return n
}

func (c *ConsultarTablero) Ejecutar(ctx context.Context, p indicador.Periodo) (*Tablero, error) {
	if !p.Valido() {
		return nil, indicador.ErrRangoFechas
	}
	conteos, err := c.Indicadores.ContarPorEstado(ctx, p)
	if err != nil {
		return nil, err
	}
	resueltos, err := c.Indicadores.ListarResueltos(ctx, indicador.FiltroReporte{Periodo: p})
	if err != nil {
		return nil, err
	}
	analistas, err := c.Indicadores.CargaAnalistas(ctx, p)
	if err != nil {
		return nil, err
	}
	prestamos, err := c.Indicadores.ListarPrestamos(ctx)
	if err != nil {
		return nil, err
	}

	t := &Tablero{
		GeneradoEn: c.Reloj.Ahora(), DiasMaximosPrestamo: diasMaximos(c.DiasMaximosPrestamo),
		Resueltos: len(resueltos), PromedioDias: indicador.PromedioDias(resueltos),
		Tiempos: indicador.AgruparTiempos(resueltos), Analistas: analistas, PrestamosVigentes: len(prestamos),
	}
	porEstado := map[tramite.Estado]int64{}
	for _, ce := range conteos {
		porEstado[ce.Estado] += ce.Cantidad
	}
	for _, e := range estadosTramite {
		n := porEstado[e]
		t.PorEstado = append(t.PorEstado, indicador.ConteoEstado{Estado: e, Cantidad: n})
		t.Radicados += n
		if !e.Cerrado() {
			t.EnCurso += n
		}
	}

	hoy := t.GeneradoEn.In(zonaBolivia)
	for _, pr := range prestamos {
		if a, enMora := indicador.EvaluarMora(pr, hoy, t.DiasMaximosPrestamo); enMora {
			t.Alertas = append(t.Alertas, a)
		}
	}
	indicador.OrdenarAlertas(t.Alertas)
	return t, nil
}
