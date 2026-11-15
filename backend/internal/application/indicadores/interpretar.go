package indicadores

import (
	"context"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
)

// InterpretarReporte convierte una frase ("trámites observados de octubre en Excel") en los filtros
// del formulario de reportes. La salida del modelo es solo una propuesta: aquí se sanea contra los
// valores reales y el usuario la revisa antes de generar nada; GenerarReporte la vuelve a validar.
type InterpretarReporte struct {
	Interprete ports.InterpreteReportes
	Reloj      ports.Clock
}

const formatoFecha = "2006-01-02"

func (c *InterpretarReporte) Ejecutar(ctx context.Context, frase string, categorias []ports.CategoriaRef) (*ports.PeticionReporte, error) {
	if c.Interprete == nil {
		return nil, indicador.ErrInterpreteNoDisponible
	}
	frase = strings.TrimSpace(frase)
	if frase == "" {
		return nil, indicador.ErrNoEntendido
	}
	p, err := c.Interprete.Interpretar(ctx, frase, c.Reloj.Ahora().In(zonaBolivia), categorias)
	if err != nil {
		return nil, err
	}
	if p == nil || !indicador.TipoReporte(p.Tipo).Valido() {
		return nil, indicador.ErrNoEntendido
	}

	out := &ports.PeticionReporte{Tipo: p.Tipo, Formato: p.Formato, Explicacion: strings.TrimSpace(p.Explicacion)}
	if !indicador.Formato(out.Formato).Valido() {
		out.Formato = string(indicador.FormatoPDF)
	}
	// Los filtros que no aplican al tipo o que no existen se descartan en vez de rechazar todo.
	if out.Tipo != string(indicador.ReporteInventario) {
		out.Desde, out.Hasta = fechaValida(p.Desde), fechaValida(p.Hasta)
		if out.Desde != "" && out.Hasta != "" && out.Desde > out.Hasta {
			out.Desde, out.Hasta = "", ""
		}
		if t := tramite.Tipo(p.TipoTramite); t.Valido() {
			out.TipoTramite = p.TipoTramite
		}
		for _, cat := range categorias {
			if cat.ID == p.CategoriaID {
				out.CategoriaID = cat.ID
			}
		}
		if out.Tipo == string(indicador.ReporteTramites) && estadoValido(tramite.Estado(p.Estado)) {
			out.Estado = p.Estado
		}
	} else if e := legajo.Estado(p.EstadoLegajo); e == legajo.EnArchivo || e == legajo.Prestado {
		out.EstadoLegajo = p.EstadoLegajo
	}
	return out, nil
}

func fechaValida(s string) string {
	if _, err := time.Parse(formatoFecha, s); err != nil {
		return ""
	}
	return s
}
