package persistence

import (
	"context"

	"gorm.io/gorm"
)

type claveTx struct{}

// TxManager implementa ports.TxManager guardando la transacción activa en el contexto
// (sección 3): los repositorios la recuperan de ahí si existe.
type TxManager struct {
	DB *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager { return &TxManager{DB: db} }

func (t *TxManager) WithinTx(ctx context.Context, fn func(context.Context) error) error {
	return t.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(context.WithValue(ctx, claveTx{}, tx))
	})
}

// db devuelve la transacción activa en ctx, o una sesión nueva sobre fallback si no hay ninguna.
func db(ctx context.Context, fallback *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(claveTx{}).(*gorm.DB); ok {
		return tx
	}
	return fallback.WithContext(ctx)
}
