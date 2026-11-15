package persistence

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Conectar abre una conexión GORM a Postgres. Usa el modo de protocolo simple de pgx porque
// la conexión pooled de Neon (PgBouncer en modo transacción) no soporta sentencias preparadas
// (sección 4.1).
func Conectar(dsn string) (*gorm.DB, error) {
	dbase, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	// Pool acotado: Neon cierra las conexiones inactivas y limita las simultáneas. Sin estos valores database/sql
	// abre conexiones sin tope y reutiliza algunas que el servidor ya cerró.
	pool, err := dbase.DB()
	if err != nil {
		return nil, err
	}
	pool.SetMaxOpenConns(20)
	pool.SetMaxIdleConns(5)
	pool.SetConnMaxIdleTime(4 * time.Minute)
	pool.SetConnMaxLifetime(30 * time.Minute)
	return dbase, nil
}
