// Package notificacion modela los avisos por correo al docente cuando su trámite cambia de estado (CU-14).
package notificacion

import (
	"context"
	"time"
)

// Estado del envío: PENDIENTE se reintenta hasta agotar los intentos; después queda FALLIDO.
type Estado string

const (
	Pendiente Estado = "PENDIENTE"
	Enviado   Estado = "ENVIADO"
	Fallido   Estado = "FALLIDO"
)

// MaxIntentos es el tope de envíos por aviso antes de darlo por fallido.
const MaxIntentos = 3

type Envio struct {
	ID             string
	TramiteID      string
	TramiteCodigo  string
	EstadoTramite  string // estado al que pasó el trámite (OBSERVADO, DICTAMINADO, ARCHIVADO)
	Destino        string
	Nombres        string
	Estado         Estado
	Intentos       int
	Error          string
	CreadoEn       time.Time
	EnviadoEn      *time.Time
}

type Filtro struct {
	Estado Estado // vacío = todos
}

type Repository interface {
	Crear(ctx context.Context, e *Envio) error
	Actualizar(ctx context.Context, e *Envio) error
	// ListarPendientes devuelve los envíos por reintentar, los más antiguos primero.
	ListarPendientes(ctx context.Context, limite int) ([]Envio, error)
	Listar(ctx context.Context, f Filtro, pagina, tamano int) ([]Envio, int64, error)
}
