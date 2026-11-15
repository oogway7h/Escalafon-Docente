package persistence

import (
	"context"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/notificacion"
)

type RepositorioNotificacion struct {
	DB *gorm.DB
}

func NewRepositorioNotificacion(dbase *gorm.DB) *RepositorioNotificacion {
	return &RepositorioNotificacion{DB: dbase}
}

func aModelo(e *notificacion.Envio) NotificacionEnvioModel {
	return NotificacionEnvioModel{
		ID: e.ID, TramiteID: e.TramiteID, TramiteCodigo: e.TramiteCodigo, EstadoTramite: e.EstadoTramite,
		Destino: e.Destino, Nombres: e.Nombres, Estado: string(e.Estado), Intentos: int16(e.Intentos),
		Error: ptr(e.Error), CreadoEn: e.CreadoEn, EnviadoEn: e.EnviadoEn,
	}
}

func (m NotificacionEnvioModel) dominio() notificacion.Envio {
	return notificacion.Envio{
		ID: m.ID, TramiteID: m.TramiteID, TramiteCodigo: m.TramiteCodigo, EstadoTramite: m.EstadoTramite,
		Destino: m.Destino, Nombres: m.Nombres, Estado: notificacion.Estado(m.Estado), Intentos: int(m.Intentos),
		Error: val(m.Error), CreadoEn: m.CreadoEn, EnviadoEn: m.EnviadoEn,
	}
}

func (r *RepositorioNotificacion) Crear(ctx context.Context, e *notificacion.Envio) error {
	m := aModelo(e)
	m.ID = "" // lo asigna la BD
	if err := db(ctx, r.DB).Omit("Tramite").Create(&m).Error; err != nil {
		return err
	}
	e.ID = m.ID
	return nil
}

func (r *RepositorioNotificacion) Actualizar(ctx context.Context, e *notificacion.Envio) error {
	return db(ctx, r.DB).Model(&NotificacionEnvioModel{}).Where("id = ?", e.ID).Updates(map[string]any{
		"estado": string(e.Estado), "intentos": e.Intentos, "error": ptr(e.Error), "enviado_en": e.EnviadoEn,
	}).Error
}

func (r *RepositorioNotificacion) ListarPendientes(ctx context.Context, limite int) ([]notificacion.Envio, error) {
	var filas []NotificacionEnvioModel
	err := db(ctx, r.DB).Where("estado = ?", string(notificacion.Pendiente)).Order("creado_en").Limit(limite).Find(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]notificacion.Envio, len(filas))
	for i, f := range filas {
		out[i] = f.dominio()
	}
	return out, nil
}

func (r *RepositorioNotificacion) Listar(ctx context.Context, f notificacion.Filtro, pagina, tamano int) ([]notificacion.Envio, int64, error) {
	q := db(ctx, r.DB).Model(&NotificacionEnvioModel{})
	if f.Estado != "" {
		q = q.Where("estado = ?", string(f.Estado))
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var filas []NotificacionEnvioModel
	if err := q.Order("creado_en DESC").Offset((pagina - 1) * tamano).Limit(tamano).Find(&filas).Error; err != nil {
		return nil, 0, err
	}
	out := make([]notificacion.Envio, len(filas))
	for i, fila := range filas {
		out[i] = fila.dominio()
	}
	return out, total, nil
}
