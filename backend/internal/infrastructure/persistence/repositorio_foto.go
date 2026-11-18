package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"escalafon/backend/internal/domain/usuario"
)

// RepositorioFoto implementa usuario.FotoRepository sobre GORM.
type RepositorioFoto struct {
	DB *gorm.DB
}

func NewRepositorioFoto(dbase *gorm.DB) *RepositorioFoto { return &RepositorioFoto{DB: dbase} }

func (r *RepositorioFoto) Obtener(ctx context.Context, usuarioID string) (*usuario.Foto, error) {
	var m UsuarioFotoModel
	if err := db(ctx, r.DB).Where("usuario_id = ?", usuarioID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usuario.ErrFotoNoEncontrada
		}
		return nil, err
	}
	return &usuario.Foto{UsuarioID: m.UsuarioID, Clave: m.Clave, TipoMime: m.TipoMime, ActualizadoEn: m.ActualizadoEn}, nil
}

func (r *RepositorioFoto) Guardar(ctx context.Context, f *usuario.Foto) error {
	m := UsuarioFotoModel{UsuarioID: f.UsuarioID, Clave: f.Clave, TipoMime: f.TipoMime, ActualizadoEn: f.ActualizadoEn}
	return db(ctx, r.DB).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "usuario_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"clave", "tipo_mime", "actualizado_en"}),
	}).Create(&m).Error
}

func (r *RepositorioFoto) Eliminar(ctx context.Context, usuarioID string) error {
	return db(ctx, r.DB).Where("usuario_id = ?", usuarioID).Delete(&UsuarioFotoModel{}).Error
}
