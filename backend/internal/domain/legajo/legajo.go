// Package legajo modela el inventario de carpetas físicas (CU-19), sus solicitudes (CU-20, CU-21)
// y los movimientos de custodia (CU-22 al CU-24).
package legajo

import (
	"context"
	"time"
)

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrNoEncontrado       erroDominio = "legajo no encontrado"
	ErrSolicitudNoExiste  erroDominio = "solicitud no encontrada"
	ErrCodigoDuplicado    erroDominio = "ya existe un legajo con ese código"
	ErrDocenteInvalido    erroDominio = "el docente no existe"
	ErrNoDisponible       erroDominio = "el legajo no está disponible en archivo"
	ErrNoPrestado         erroDominio = "el legajo no está prestado"
	ErrSolicitudAbierta   erroDominio = "ya hay una solicitud en curso para ese legajo"
	ErrEstadoSolicitud    erroDominio = "la solicitud no está en el estado requerido"
	ErrMotivoRequerido    erroDominio = "el motivo es obligatorio"
	ErrTramiteRequerido   erroDominio = "indica el trámite para el que se solicita la carpeta"
	ErrTramiteInvalido    erroDominio = "el trámite no existe, no está a tu cargo o ya está cerrado"
	ErrFechaLimiteInvalid erroDominio = "la fecha límite debe ser posterior a hoy"
	ErrEstadoMaterial     erroDominio = "el estado material no es válido"
	ErrObservacionRequerida erroDominio = "la observación es obligatoria si la carpeta no está en buen estado"
)

// Estado es la ubicación actual de la carpeta.
type Estado string

const (
	EnArchivo Estado = "EN_ARCHIVO"
	Prestado  Estado = "PRESTADO"
)

// EstadoSolicitud recorre PENDIENTE → ATENDIDA → DESPACHADA → DEVUELTA, o PENDIENTE → RECHAZADA.
type EstadoSolicitud string

const (
	Pendiente  EstadoSolicitud = "PENDIENTE"
	Atendida   EstadoSolicitud = "ATENDIDA"
	Rechazada  EstadoSolicitud = "RECHAZADA"
	Despachada EstadoSolicitud = "DESPACHADA"
	Devuelta   EstadoSolicitud = "DEVUELTA"
)

// EstadoMaterial describe la condición física de la carpeta al salir o volver.
type EstadoMaterial string

const (
	Bueno       EstadoMaterial = "BUENO"
	Regular     EstadoMaterial = "REGULAR"
	Deteriorado EstadoMaterial = "DETERIORADO"
)

func (e EstadoMaterial) Valido() bool { return e == Bueno || e == Regular || e == Deteriorado }

type TipoMovimiento string

const (
	Salida  TipoMovimiento = "SALIDA"
	Entrada TipoMovimiento = "ENTRADA"
)

type Legajo struct {
	ID            int32
	Codigo        string
	DocenteID     string
	DocenteCodigo string
	DocenteNombre string
	Gaveta        string
	Orden         int
	Anexos        string
	Observacion   string
	Estado        Estado
}

type Solicitud struct {
	ID            string
	LegajoID      int32
	LegajoCodigo  string
	DocenteNombre string
	TramiteID     string // trámite para el que se pide la carpeta
	TramiteCodigo string
	SolicitanteID string
	Solicitante   string
	Motivo        string
	Estado        EstadoSolicitud
	Respuesta     string // motivo del rechazo, o nota al atender
	AtendidoPor   string
	CreadoEn      time.Time
	AtendidoEn    *time.Time
	FechaLimite   *time.Time // solo desde que se despacha
}

type Movimiento struct {
	ID              int64
	LegajoID        int32
	LegajoCodigo    string
	DocenteNombre   string
	SolicitudID     string
	Tipo            TipoMovimiento
	Custodio        string
	CustodioID      string
	FechaLimite     *time.Time
	EstadoMaterial  EstadoMaterial
	Observacion     string
	RegistradoPor   string
	RegistradoPorID string
	Fecha           time.Time
	Vencido         bool // salida cuya carpeta sigue prestada pasada la fecha límite
}

type FiltroLegajos struct {
	Busqueda string // código de legajo o de docente (exactos), cédula o nombre
	Estado   Estado
}

type FiltroSolicitudes struct {
	Estado        EstadoSolicitud
	SolicitanteID string // vacío = todas
}

type FiltroMovimientos struct {
	Busqueda string
	Tipo     TipoMovimiento
	Vencidos bool
}

type Repository interface {
	ListarLegajos(ctx context.Context, f FiltroLegajos, pagina, tamano int) ([]Legajo, int64, error)
	BuscarLegajo(ctx context.Context, id int32) (*Legajo, error)
	// ActualizarLegajo cambia código, gaveta, orden, anexos y observación.
	ActualizarLegajo(ctx context.Context, l *Legajo) error
	CrearLegajo(ctx context.Context, l *Legajo) error

	CrearSolicitud(ctx context.Context, s *Solicitud) error
	BuscarSolicitud(ctx context.Context, id string) (*Solicitud, error)
	ListarSolicitudes(ctx context.Context, f FiltroSolicitudes, pagina, tamano int) ([]Solicitud, int64, error)
	// HaySolicitudAbierta indica si el legajo tiene una solicitud PENDIENTE, ATENDIDA o DESPACHADA.
	HaySolicitudAbierta(ctx context.Context, legajoID int32) (bool, error)
	// ActualizarSolicitud pasa la solicitud de un estado a otro (condicional: ErrEstadoSolicitud si ya cambió).
	ActualizarSolicitud(ctx context.Context, id string, desde, hacia EstadoSolicitud, respuesta, usuarioID string, en time.Time) error
	// SolicitudDespachada devuelve la solicitud DESPACHADA vigente del legajo.
	SolicitudDespachada(ctx context.Context, legajoID int32) (*Solicitud, error)

	// CambiarEstadoLegajo es condicional (desde → hacia): devuelve false si el legajo no estaba en "desde".
	CambiarEstadoLegajo(ctx context.Context, id int32, desde, hacia Estado) (bool, error)
	AnexarMovimiento(ctx context.Context, m *Movimiento) error
	ListarMovimientos(ctx context.Context, f FiltroMovimientos, pagina, tamano int) ([]Movimiento, int64, error)
}
