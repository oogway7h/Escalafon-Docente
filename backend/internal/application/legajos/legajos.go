// Package legajos implementa CU-19 (inventario), CU-20 y CU-21 (solicitudes de carpeta) y
// CU-22 al CU-24 (check-out, check-in y bitácora de custodia).
package legajos

import (
	"context"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
)

func acotar(pagina, tamano int) (int, int) {
	if pagina < 1 {
		pagina = 1
	}
	if tamano < 1 || tamano > 100 {
		tamano = 10
	}
	return pagina, tamano
}

func auditar(ctx context.Context, a auditoria.Repository, reloj ports.Clock, actorID, accion, entidad, entidadID string, detalle map[string]any) {
	_ = a.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actorID, Accion: accion, Entidad: entidad, EntidadID: entidadID, Detalle: detalle, Fecha: reloj.Ahora(),
	})
}

// --- CU-19 Inventario -------------------------------------------------------------

type ListarLegajos struct{ Legajos legajo.Repository }

func (c *ListarLegajos) Ejecutar(ctx context.Context, f legajo.FiltroLegajos, pagina, tamano int) ([]legajo.Legajo, int64, error) {
	pagina, tamano = acotar(pagina, tamano)
	return c.Legajos.ListarLegajos(ctx, f, pagina, tamano)
}

type ObtenerLegajo struct{ Legajos legajo.Repository }

func (c *ObtenerLegajo) Ejecutar(ctx context.Context, id int32) (*legajo.Legajo, error) {
	return c.Legajos.BuscarLegajo(ctx, id)
}

type GuardarEntrada struct {
	ActorID     string
	ID          int32  // 0 al crear
	DocenteID   string // solo al crear
	Codigo      string
	Gaveta      string
	Orden       int
	Anexos      string
	Observacion string
}

// GuardarLegajo crea o edita una carpeta. El estado (en archivo / prestado) solo lo cambian los movimientos.
type GuardarLegajo struct {
	Legajos   legajo.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *GuardarLegajo) Ejecutar(ctx context.Context, in GuardarEntrada) (*legajo.Legajo, error) {
	l := &legajo.Legajo{
		ID: in.ID, DocenteID: in.DocenteID, Codigo: strings.ToUpper(strings.TrimSpace(in.Codigo)),
		Gaveta: strings.TrimSpace(in.Gaveta), Orden: in.Orden, Anexos: strings.TrimSpace(in.Anexos),
		Observacion: strings.TrimSpace(in.Observacion),
	}
	accion := auditoria.AccionCrear
	if in.ID == 0 {
		l.Estado = legajo.EnArchivo
		if err := c.Legajos.CrearLegajo(ctx, l); err != nil {
			return nil, err
		}
	} else {
		if err := c.Legajos.ActualizarLegajo(ctx, l); err != nil {
			return nil, err
		}
		accion = auditoria.AccionEditar
	}
	auditar(ctx, c.Auditoria, c.Reloj, in.ActorID, accion, auditoria.EntidadLegajo, l.Codigo, nil)
	return c.Legajos.BuscarLegajo(ctx, l.ID)
}

// --- CU-20 Solicitar legajo -------------------------------------------------------

// Solicitar pide una carpeta física al archivo. El trámite es opcional: quien tiene el permiso puede pedir
// cualquier carpeta. Si se indica uno, debe estar a cargo del solicitante, abierto y ser del docente de la
// carpeta, porque esa solicitud ligada es la que habilita iniciar la revisión del trámite.
type Solicitar struct {
	Legajos   legajo.Repository
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *Solicitar) Ejecutar(ctx context.Context, actorID string, legajoID int32, tramiteID, motivo string) (*legajo.Solicitud, error) {
	motivo = strings.TrimSpace(motivo)
	if motivo == "" {
		return nil, legajo.ErrMotivoRequerido
	}
	l, err := c.Legajos.BuscarLegajo(ctx, legajoID)
	if err != nil {
		return nil, err
	}
	detalle := map[string]any{"legajo": l.Codigo}
	if tramiteID = strings.TrimSpace(tramiteID); tramiteID != "" {
		t, err := c.Tramites.BuscarPorID(ctx, tramiteID)
		if err != nil || t.AnalistaID == nil || *t.AnalistaID != actorID || t.Estado.Cerrado() || t.DocenteID != l.DocenteID {
			return nil, legajo.ErrTramiteInvalido
		}
		detalle["tramite"] = t.Codigo
	}
	if l.Estado != legajo.EnArchivo {
		return nil, legajo.ErrNoDisponible
	}
	if abierta, err := c.Legajos.HaySolicitudAbierta(ctx, legajoID); err != nil {
		return nil, err
	} else if abierta {
		return nil, legajo.ErrSolicitudAbierta
	}
	s := &legajo.Solicitud{LegajoID: legajoID, SolicitanteID: actorID, TramiteID: tramiteID, Motivo: motivo, Estado: legajo.Pendiente, CreadoEn: c.Reloj.Ahora()}
	if err := c.Legajos.CrearSolicitud(ctx, s); err != nil {
		return nil, err
	}
	auditar(ctx, c.Auditoria, c.Reloj, actorID, auditoria.AccionSolicitarCarpeta, auditoria.EntidadSolicitud, s.ID, detalle)
	return c.Legajos.BuscarSolicitud(ctx, s.ID)
}

type ListarSolicitudes struct{ Legajos legajo.Repository }

func (c *ListarSolicitudes) Ejecutar(ctx context.Context, f legajo.FiltroSolicitudes, pagina, tamano int) ([]legajo.Solicitud, int64, error) {
	pagina, tamano = acotar(pagina, tamano)
	return c.Legajos.ListarSolicitudes(ctx, f, pagina, tamano)
}

// --- CU-21 Atender solicitud ------------------------------------------------------

// Atender aprueba (la carpeta queda lista para despachar) o rechaza con motivo una solicitud pendiente.
type Atender struct {
	Legajos   legajo.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *Atender) Ejecutar(ctx context.Context, actorID, solicitudID string, aprobar bool, nota string) error {
	nota = strings.TrimSpace(nota)
	hacia, accion := legajo.Atendida, auditoria.AccionAtenderSolicitud
	if !aprobar {
		if nota == "" {
			return legajo.ErrMotivoRequerido
		}
		hacia, accion = legajo.Rechazada, auditoria.AccionRechazarSolicitud
	}
	if err := c.Legajos.ActualizarSolicitud(ctx, solicitudID, legajo.Pendiente, hacia, nota, actorID, c.Reloj.Ahora()); err != nil {
		return err
	}
	auditar(ctx, c.Auditoria, c.Reloj, actorID, accion, auditoria.EntidadSolicitud, solicitudID, nil)
	return nil
}

// --- CU-22 Despachar (check-out) --------------------------------------------------

type DespacharEntrada struct {
	ActorID        string
	SolicitudID    string
	FechaLimite    time.Time
	EstadoMaterial legajo.EstadoMaterial
	Observacion    string
}

// Despachar entrega la carpeta al solicitante: la solicitud debe estar atendida y la carpeta en archivo.
type Despachar struct {
	Legajos   legajo.Repository
	Auditoria auditoria.Repository
	Tx        ports.TxManager
	Reloj     ports.Clock
}

// validarSalida comprueba el estado material y que la fecha límite sea posterior a hoy (compara solo la fecha, UTC).
func validarSalida(in DespacharEntrada, ahora time.Time) error {
	if !in.EstadoMaterial.Valido() {
		return legajo.ErrEstadoMaterial
	}
	if !in.FechaLimite.Truncate(24 * time.Hour).After(ahora.Truncate(24 * time.Hour)) {
		return legajo.ErrFechaLimiteInvalid
	}
	return nil
}

func (c *Despachar) Ejecutar(ctx context.Context, in DespacharEntrada) error {
	ahora := c.Reloj.Ahora()
	if err := validarSalida(in, ahora); err != nil {
		return err
	}
	s, err := c.Legajos.BuscarSolicitud(ctx, in.SolicitudID)
	if err != nil {
		return err
	}
	if s.Estado != legajo.Atendida {
		return legajo.ErrEstadoSolicitud
	}
	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		if ok, err := c.Legajos.CambiarEstadoLegajo(ctx, s.LegajoID, legajo.EnArchivo, legajo.Prestado); err != nil {
			return err
		} else if !ok {
			return legajo.ErrNoDisponible
		}
		if err := c.Legajos.ActualizarSolicitud(ctx, s.ID, legajo.Atendida, legajo.Despachada, s.Respuesta, in.ActorID, ahora); err != nil {
			return err
		}
		return c.Legajos.AnexarMovimiento(ctx, &legajo.Movimiento{
			LegajoID: s.LegajoID, SolicitudID: s.ID, Tipo: legajo.Salida, CustodioID: s.SolicitanteID, FechaLimite: &in.FechaLimite,
			EstadoMaterial: in.EstadoMaterial, Observacion: strings.TrimSpace(in.Observacion), RegistradoPorID: in.ActorID, Fecha: ahora,
		})
	})
	if err != nil {
		return err
	}
	auditar(ctx, c.Auditoria, c.Reloj, in.ActorID, auditoria.AccionDespacharCarpeta, auditoria.EntidadLegajo, s.LegajoCodigo, nil)
	return nil
}

// Entregar aprueba (si sigue pendiente) y despacha la carpeta en un solo paso: es lo que hace el personal de
// archivo en la pantalla unificada. Valida todo antes de aprobar, para no dejar la solicitud a medias por un
// dato inválido; una solicitud ya aprobada se entrega directamente.
// ponytail: no usa una transacción envolvente (Despachar abre la suya); si la entrega falla tras aprobar
// (carrera con otro despacho), la solicitud queda «aprobada» y se puede entregar de nuevo.
type Entregar struct {
	Legajos   legajo.Repository
	Atender   *Atender
	Despachar *Despachar
}

func (c *Entregar) Ejecutar(ctx context.Context, in DespacharEntrada) error {
	if err := validarSalida(in, c.Despachar.Reloj.Ahora()); err != nil {
		return err
	}
	s, err := c.Legajos.BuscarSolicitud(ctx, in.SolicitudID)
	if err != nil {
		return err
	}
	switch s.Estado {
	case legajo.Pendiente:
		l, err := c.Legajos.BuscarLegajo(ctx, s.LegajoID)
		if err != nil {
			return err
		}
		if l.Estado != legajo.EnArchivo {
			return legajo.ErrNoDisponible
		}
		if err := c.Atender.Ejecutar(ctx, in.ActorID, s.ID, true, ""); err != nil {
			return err
		}
	case legajo.Atendida:
	default:
		return legajo.ErrEstadoSolicitud
	}
	return c.Despachar.Ejecutar(ctx, in)
}

// --- CU-23 Devolución (check-in) --------------------------------------------------

type DevolverEntrada struct {
	ActorID        string
	LegajoID       int32
	EstadoMaterial legajo.EstadoMaterial
	Observacion    string
}

type Devolver struct {
	Legajos   legajo.Repository
	Auditoria auditoria.Repository
	Tx        ports.TxManager
	Reloj     ports.Clock
}

func (c *Devolver) Ejecutar(ctx context.Context, in DevolverEntrada) error {
	if !in.EstadoMaterial.Valido() {
		return legajo.ErrEstadoMaterial
	}
	if in.EstadoMaterial != legajo.Bueno && strings.TrimSpace(in.Observacion) == "" {
		return legajo.ErrObservacionRequerida // la novedad queda registrada para asignar responsabilidad
	}
	s, err := c.Legajos.SolicitudDespachada(ctx, in.LegajoID)
	if err != nil {
		return err
	}
	ahora := c.Reloj.Ahora()
	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		if ok, err := c.Legajos.CambiarEstadoLegajo(ctx, in.LegajoID, legajo.Prestado, legajo.EnArchivo); err != nil {
			return err
		} else if !ok {
			return legajo.ErrNoPrestado
		}
		if err := c.Legajos.ActualizarSolicitud(ctx, s.ID, legajo.Despachada, legajo.Devuelta, s.Respuesta, in.ActorID, ahora); err != nil {
			return err
		}
		return c.Legajos.AnexarMovimiento(ctx, &legajo.Movimiento{
			LegajoID: in.LegajoID, SolicitudID: s.ID, Tipo: legajo.Entrada, CustodioID: s.SolicitanteID,
			EstadoMaterial: in.EstadoMaterial, Observacion: strings.TrimSpace(in.Observacion), RegistradoPorID: in.ActorID, Fecha: ahora,
		})
	})
	if err != nil {
		return err
	}
	auditar(ctx, c.Auditoria, c.Reloj, in.ActorID, auditoria.AccionDevolverCarpeta, auditoria.EntidadLegajo, s.LegajoCodigo, nil)
	return nil
}

// --- CU-24 Bitácora ---------------------------------------------------------------

type ListarMovimientos struct{ Legajos legajo.Repository }

func (c *ListarMovimientos) Ejecutar(ctx context.Context, f legajo.FiltroMovimientos, pagina, tamano int) ([]legajo.Movimiento, int64, error) {
	pagina, tamano = acotar(pagina, tamano)
	return c.Legajos.ListarMovimientos(ctx, f, pagina, tamano)
}
