package token

import (
	"context"
	"time"
)

type Proposito string

const (
	PropositoRecuperacion Proposito = "RECUPERACION"
	PropositoActivacion   Proposito = "ACTIVACION"
)

// TokenRecuperacion respalda tanto la recuperación de contraseña como la activación de cuenta.
type TokenRecuperacion struct {
	ID        string
	UsuarioID string
	Proposito Proposito
	TokenHash string
	ExpiraEn  time.Time
	UsadoEn   *time.Time
	CreadoEn  time.Time
}

// EsValido indica si el token, al instante ahora, no fue usado y no expiró.
func (t *TokenRecuperacion) EsValido(ahora time.Time) bool {
	return t.UsadoEn == nil && t.ExpiraEn.After(ahora)
}

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const ErrInvalido erroDominio = "el token es inválido o expiró"

// Repository es el puerto de persistencia de tokens de un solo uso.
type Repository interface {
	Crear(ctx context.Context, t *TokenRecuperacion) error
	BuscarPorHash(ctx context.Context, hash string) (*TokenRecuperacion, error)
	MarcarUsado(ctx context.Context, id string, cuando time.Time) error
	// InvalidarPendientes marca como usados los tokens no usados de ese propósito (ReenviarActivacion).
	InvalidarPendientes(ctx context.Context, usuarioID string, proposito Proposito, cuando time.Time) error
	// ExisteTokenUsado indica si el usuario ya consumió alguna vez un token de ese propósito
	// (así se sabe si una cuenta ya fue activada, sección 8 HU-02).
	ExisteTokenUsado(ctx context.Context, usuarioID string, proposito Proposito) (bool, error)
	// ContarDesde cuenta los tokens de ese propósito creados para el usuario a partir de desde (tope de solicitudes).
	ContarDesde(ctx context.Context, usuarioID string, proposito Proposito, desde time.Time) (int64, error)
}
