package persistence

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/auditoria"
)

type RepositorioAuditoria struct {
	DB *gorm.DB
}

func NewRepositorioAuditoria(dbase *gorm.DB) *RepositorioAuditoria { return &RepositorioAuditoria{DB: dbase} }

func (r *RepositorioAuditoria) Registrar(ctx context.Context, e *auditoria.Evento) error {
	var detalle *string
	if len(e.Detalle) > 0 {
		b, err := json.Marshal(e.Detalle)
		if err != nil {
			return err
		}
		s := string(b)
		detalle = &s
	}
	var ip *string
	if e.DireccionIP != "" {
		ip = &e.DireccionIP
	}
	m := AuditoriaModel{
		UsuarioID: e.UsuarioID, Accion: e.Accion, Entidad: e.Entidad, EntidadID: e.EntidadID,
		Detalle: detalle, DireccionIP: ip, Fecha: e.Fecha,
	}
	if err := db(ctx, r.DB).Create(&m).Error; err != nil {
		return err
	}
	e.ID = m.ID
	return nil
}

// Listar solo lee: cuenta las coincidencias y trae una página ordenada por fecha descendente.
// El autor se carga con un único Preload por página.
func (r *RepositorioAuditoria) Listar(ctx context.Context, f auditoria.Filtro) ([]auditoria.Registro, int64, error) {
	q := db(ctx, r.DB).Model(&AuditoriaModel{})
	if f.UsuarioID != "" {
		q = q.Where("usuario_id = ?", f.UsuarioID)
	}
	if f.Accion != "" {
		q = q.Where("accion = ?", f.Accion)
	}
	if f.Desde != nil {
		q = q.Where("fecha >= ?", *f.Desde)
	}
	if f.Hasta != nil {
		q = q.Where("fecha < ?", *f.Hasta)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var modelos []AuditoriaModel
	if err := q.Preload("Usuario").Order("fecha DESC, id DESC").
		Offset((f.Pagina - 1) * f.Tamano).Limit(f.Tamano).Find(&modelos).Error; err != nil {
		return nil, 0, err
	}
	out := make([]auditoria.Registro, len(modelos))
	for i := range modelos {
		reg, err := aDominioRegistro(&modelos[i])
		if err != nil {
			return nil, 0, err
		}
		out[i] = reg
	}
	return out, total, nil
}

func aDominioRegistro(m *AuditoriaModel) (auditoria.Registro, error) {
	reg := auditoria.Registro{Evento: auditoria.Evento{
		ID: m.ID, UsuarioID: m.UsuarioID, Accion: m.Accion, Entidad: m.Entidad, EntidadID: m.EntidadID,
		DireccionIP: val(m.DireccionIP), Fecha: m.Fecha,
	}}
	if m.Detalle != nil {
		if err := json.Unmarshal([]byte(*m.Detalle), &reg.Detalle); err != nil {
			return reg, err
		}
	}
	if m.Usuario != nil {
		reg.Autor = &auditoria.Autor{ID: m.Usuario.ID, Nombres: m.Usuario.Nombres, Apellidos: m.Usuario.Apellidos, Correo: m.Usuario.Correo}
	}
	return reg, nil
}

// ListarAutores devuelve los usuarios con al menos un evento (EXISTS usa idx_auditoria_usuario_fecha).
func (r *RepositorioAuditoria) ListarAutores(ctx context.Context) ([]auditoria.Autor, error) {
	var modelos []UsuarioModel
	if err := db(ctx, r.DB).
		Where("EXISTS (SELECT 1 FROM auditoria a WHERE a.usuario_id = usuario.id)").
		Order("apellidos, nombres").Find(&modelos).Error; err != nil {
		return nil, err
	}
	out := make([]auditoria.Autor, len(modelos))
	for i, m := range modelos {
		out[i] = auditoria.Autor{ID: m.ID, Nombres: m.Nombres, Apellidos: m.Apellidos, Correo: m.Correo}
	}
	return out, nil
}
