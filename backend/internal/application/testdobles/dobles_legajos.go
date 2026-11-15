package testdobles

import (
	"context"
	"strconv"
	"time"

	"escalafon/backend/internal/domain/legajo"
)

// LegajoRepoFalso guarda legajos, solicitudes y movimientos en memoria.
type LegajoRepoFalso struct {
	Legajos     map[int32]*legajo.Legajo
	Solicitudes map[string]*legajo.Solicitud
	Movimientos []legajo.Movimiento
	contador    int
}

func NewLegajoRepoFalso() *LegajoRepoFalso {
	return &LegajoRepoFalso{Legajos: map[int32]*legajo.Legajo{}, Solicitudes: map[string]*legajo.Solicitud{}}
}

func (r *LegajoRepoFalso) Sembrar(l legajo.Legajo) { r.Legajos[l.ID] = &l }

func (r *LegajoRepoFalso) ListarLegajos(context.Context, legajo.FiltroLegajos, int, int) ([]legajo.Legajo, int64, error) {
	var out []legajo.Legajo
	for _, l := range r.Legajos {
		out = append(out, *l)
	}
	return out, int64(len(out)), nil
}

func (r *LegajoRepoFalso) BuscarLegajo(_ context.Context, id int32) (*legajo.Legajo, error) {
	if l, ok := r.Legajos[id]; ok {
		c := *l
		return &c, nil
	}
	return nil, legajo.ErrNoEncontrado
}

func (r *LegajoRepoFalso) ActualizarLegajo(_ context.Context, l *legajo.Legajo) error {
	e, ok := r.Legajos[l.ID]
	if !ok {
		return legajo.ErrNoEncontrado
	}
	e.Gaveta, e.Orden, e.Anexos, e.Observacion = l.Gaveta, l.Orden, l.Anexos, l.Observacion
	if l.Codigo != "" {
		e.Codigo = l.Codigo
	}
	return nil
}

func (r *LegajoRepoFalso) CrearLegajo(_ context.Context, l *legajo.Legajo) error {
	r.contador++
	l.ID = int32(r.contador) + 1000
	if l.Codigo == "" {
		l.Codigo = "LEG-" + strconv.Itoa(int(l.ID))
	}
	c := *l
	r.Legajos[l.ID] = &c
	return nil
}

func (r *LegajoRepoFalso) CrearSolicitud(_ context.Context, s *legajo.Solicitud) error {
	r.contador++
	s.ID = "sol-" + strconv.Itoa(r.contador)
	c := *s
	c.LegajoCodigo = r.Legajos[s.LegajoID].Codigo
	r.Solicitudes[s.ID] = &c
	return nil
}

func (r *LegajoRepoFalso) BuscarSolicitud(_ context.Context, id string) (*legajo.Solicitud, error) {
	if s, ok := r.Solicitudes[id]; ok {
		c := *s
		return &c, nil
	}
	return nil, legajo.ErrSolicitudNoExiste
}

func (r *LegajoRepoFalso) ListarSolicitudes(_ context.Context, f legajo.FiltroSolicitudes, _, _ int) ([]legajo.Solicitud, int64, error) {
	var out []legajo.Solicitud
	for _, s := range r.Solicitudes {
		if (f.Estado == "" || s.Estado == f.Estado) && (f.SolicitanteID == "" || s.SolicitanteID == f.SolicitanteID) {
			out = append(out, *s)
		}
	}
	return out, int64(len(out)), nil
}

func (r *LegajoRepoFalso) HaySolicitudAbierta(_ context.Context, legajoID int32) (bool, error) {
	for _, s := range r.Solicitudes {
		if s.LegajoID == legajoID && (s.Estado == legajo.Pendiente || s.Estado == legajo.Atendida || s.Estado == legajo.Despachada) {
			return true, nil
		}
	}
	return false, nil
}

func (r *LegajoRepoFalso) ActualizarSolicitud(_ context.Context, id string, desde, hacia legajo.EstadoSolicitud, respuesta, _ string, _ time.Time) error {
	s, ok := r.Solicitudes[id]
	if !ok {
		return legajo.ErrSolicitudNoExiste
	}
	if s.Estado != desde {
		return legajo.ErrEstadoSolicitud
	}
	s.Estado, s.Respuesta = hacia, respuesta
	return nil
}

func (r *LegajoRepoFalso) SolicitudDespachada(_ context.Context, legajoID int32) (*legajo.Solicitud, error) {
	for _, s := range r.Solicitudes {
		if s.LegajoID == legajoID && s.Estado == legajo.Despachada {
			c := *s
			return &c, nil
		}
	}
	return nil, legajo.ErrNoPrestado
}

func (r *LegajoRepoFalso) CambiarEstadoLegajo(_ context.Context, id int32, desde, hacia legajo.Estado) (bool, error) {
	l, ok := r.Legajos[id]
	if !ok || l.Estado != desde {
		return false, nil
	}
	l.Estado = hacia
	return true, nil
}

func (r *LegajoRepoFalso) AnexarMovimiento(_ context.Context, m *legajo.Movimiento) error {
	r.Movimientos = append(r.Movimientos, *m)
	return nil
}

func (r *LegajoRepoFalso) ListarMovimientos(context.Context, legajo.FiltroMovimientos, int, int) ([]legajo.Movimiento, int64, error) {
	return r.Movimientos, int64(len(r.Movimientos)), nil
}
