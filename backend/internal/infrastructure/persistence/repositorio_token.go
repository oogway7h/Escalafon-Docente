package persistence

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/token"
)

type RepositorioToken struct {
	DB *gorm.DB
}

func NewRepositorioToken(dbase *gorm.DB) *RepositorioToken { return &RepositorioToken{DB: dbase} }

func aDominioToken(m *TokenRecuperacionModel) *token.TokenRecuperacion {
	return &token.TokenRecuperacion{
		ID: m.ID, UsuarioID: m.UsuarioID, Proposito: token.Proposito(m.Proposito),
		TokenHash: m.TokenHash, ExpiraEn: m.ExpiraEn, UsadoEn: m.UsadoEn, CreadoEn: m.CreadoEn,
	}
}

func (r *RepositorioToken) Crear(ctx context.Context, t *token.TokenRecuperacion) error {
	m := TokenRecuperacionModel{
		UsuarioID: t.UsuarioID, Proposito: string(t.Proposito), TokenHash: t.TokenHash, ExpiraEn: t.ExpiraEn,
	}
	if err := db(ctx, r.DB).Create(&m).Error; err != nil {
		return err
	}
	t.ID = m.ID
	t.CreadoEn = m.CreadoEn
	return nil
}

func (r *RepositorioToken) BuscarPorHash(ctx context.Context, hash string) (*token.TokenRecuperacion, error) {
	var m TokenRecuperacionModel
	if err := db(ctx, r.DB).First(&m, "token_hash = ?", hash).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, token.ErrInvalido
		}
		return nil, err
	}
	return aDominioToken(&m), nil
}

func (r *RepositorioToken) MarcarUsado(ctx context.Context, id string, cuando time.Time) error {
	return db(ctx, r.DB).Model(&TokenRecuperacionModel{}).Where("id = ?", id).Update("usado_en", cuando).Error
}

func (r *RepositorioToken) InvalidarPendientes(ctx context.Context, usuarioID string, proposito token.Proposito, cuando time.Time) error {
	return db(ctx, r.DB).Model(&TokenRecuperacionModel{}).
		Where("usuario_id = ? AND proposito = ? AND usado_en IS NULL", usuarioID, string(proposito)).
		Update("usado_en", cuando).Error
}

func (r *RepositorioToken) ExisteTokenUsado(ctx context.Context, usuarioID string, proposito token.Proposito) (bool, error) {
	var n int64
	err := db(ctx, r.DB).Model(&TokenRecuperacionModel{}).
		Where("usuario_id = ? AND proposito = ? AND usado_en IS NOT NULL", usuarioID, string(proposito)).
		Count(&n).Error
	return n > 0, err
}

func (r *RepositorioToken) ContarDesde(ctx context.Context, usuarioID string, proposito token.Proposito, desde time.Time) (int64, error) {
	var n int64
	err := db(ctx, r.DB).Model(&TokenRecuperacionModel{}).
		Where("usuario_id = ? AND proposito = ? AND creado_en >= ?", usuarioID, string(proposito), desde).
		Count(&n).Error
	return n, err
}
