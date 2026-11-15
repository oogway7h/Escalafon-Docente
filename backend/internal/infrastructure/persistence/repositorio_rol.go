package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"escalafon/backend/internal/domain/rol"
)

type RepositorioRol struct {
	DB *gorm.DB
}

func NewRepositorioRol(dbase *gorm.DB) *RepositorioRol { return &RepositorioRol{DB: dbase} }

func aDominioRol(m *RolModel) *rol.Rol {
	r := &rol.Rol{ID: m.ID, Nombre: m.Nombre, Descripcion: m.Descripcion, EsSistema: m.EsSistema, Activo: m.Activo, CreadoEn: m.CreadoEn}
	return r
}

func (r *RepositorioRol) Listar(ctx context.Context) ([]rol.ConConteo, error) {
	var modelos []RolModel
	if err := db(ctx, r.DB).Order("nombre").Find(&modelos).Error; err != nil {
		return nil, err
	}
	out := make([]rol.ConConteo, len(modelos))
	for i, m := range modelos {
		var conteo int64
		db(ctx, r.DB).Model(&UsuarioModel{}).Where("rol_id = ?", m.ID).Count(&conteo)
		out[i] = rol.ConConteo{Rol: *aDominioRol(&m), UsuariosCount: conteo}
	}
	return out, nil
}

func (r *RepositorioRol) BuscarPorID(ctx context.Context, id int16) (*rol.Rol, error) {
	var m RolModel
	if err := db(ctx, r.DB).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rol.ErrNoEncontrado
		}
		return nil, err
	}
	out := aDominioRol(&m)
	permisos, err := r.permisosDeRol(ctx, m.ID)
	if err != nil {
		return nil, err
	}
	out.Permisos = permisos
	return out, nil
}

func (r *RepositorioRol) BuscarPorNombre(ctx context.Context, nombre string) (*rol.Rol, error) {
	var m RolModel
	if err := db(ctx, r.DB).First(&m, "nombre = ?", nombre).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rol.ErrNoEncontrado
		}
		return nil, err
	}
	return aDominioRol(&m), nil
}

func (r *RepositorioRol) permisosDeRol(ctx context.Context, rolID int16) ([]rol.Permiso, error) {
	var permisos []PermisoModel
	if err := db(ctx, r.DB).
		Joins("JOIN rol_permiso rp ON rp.permiso_id = permiso.id").
		Where("rp.rol_id = ?", rolID).Order("permiso.modulo, permiso.nombre").Find(&permisos).Error; err != nil {
		return nil, err
	}
	out := make([]rol.Permiso, len(permisos))
	for i, p := range permisos {
		out[i] = rol.Permiso{ID: p.ID, Codigo: p.Codigo, Nombre: p.Nombre, Modulo: p.Modulo, CasoUso: p.CasoUso}
	}
	return out, nil
}

func (r *RepositorioRol) Guardar(ctx context.Context, rl *rol.Rol, permisoIDs []int16) error {
	return db(ctx, r.DB).Transaction(func(tx *gorm.DB) error {
		m := RolModel{ID: rl.ID, Nombre: rl.Nombre, Descripcion: rl.Descripcion, EsSistema: rl.EsSistema, Activo: rl.Activo}
		if rl.ID == 0 {
			if err := tx.Create(&m).Error; err != nil {
				if esViolacionUnicidad(err) {
					return rol.ErrDuplicado
				}
				return err
			}
			rl.ID = m.ID
			rl.CreadoEn = m.CreadoEn
		} else {
			if err := tx.Model(&RolModel{}).Where("id = ?", rl.ID).
				Updates(map[string]any{"nombre": m.Nombre, "descripcion": m.Descripcion}).Error; err != nil {
				if esViolacionUnicidad(err) {
					return rol.ErrDuplicado
				}
				return err
			}
			if err := tx.Where("rol_id = ?", rl.ID).Delete(&RolPermisoModel{}).Error; err != nil {
				return err
			}
		}
		for _, pid := range permisoIDs {
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).
				Create(&RolPermisoModel{RolID: rl.ID, PermisoID: pid}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *RepositorioRol) Eliminar(ctx context.Context, id int16) error {
	n, err := r.ContarUsuarios(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return &rol.ErrConUsuarios{Cantidad: n}
	}
	res := db(ctx, r.DB).Delete(&RolModel{}, "id = ?", id)
	if res.Error != nil {
		if esViolacionFK(res.Error) {
			return &rol.ErrConUsuarios{Cantidad: n}
		}
		return res.Error
	}
	return nil
}

func (r *RepositorioRol) ContarUsuarios(ctx context.Context, id int16) (int64, error) {
	var n int64
	err := db(ctx, r.DB).Model(&UsuarioModel{}).Where("rol_id = ?", id).Count(&n).Error
	return n, err
}

func (r *RepositorioRol) IDsUsuariosDelRol(ctx context.Context, id int16) ([]string, error) {
	var ids []string
	err := db(ctx, r.DB).Model(&UsuarioModel{}).Where("rol_id = ?", id).Pluck("id", &ids).Error
	return ids, err
}

func (r *RepositorioRol) ListarPermisos(ctx context.Context) ([]rol.Permiso, error) {
	var permisos []PermisoModel
	if err := db(ctx, r.DB).Order("modulo, nombre").Find(&permisos).Error; err != nil {
		return nil, err
	}
	out := make([]rol.Permiso, len(permisos))
	for i, p := range permisos {
		out[i] = rol.Permiso{ID: p.ID, Codigo: p.Codigo, Nombre: p.Nombre, Modulo: p.Modulo, CasoUso: p.CasoUso}
	}
	return out, nil
}

func (r *RepositorioRol) PermisosExisten(ctx context.Context, ids []int16) (bool, error) {
	if len(ids) == 0 {
		return true, nil
	}
	var n int64
	if err := db(ctx, r.DB).Model(&PermisoModel{}).Where("id IN ?", ids).Count(&n).Error; err != nil {
		return false, err
	}
	return int(n) == len(uniqueInt16(ids)), nil
}

func (r *RepositorioRol) PermisosPorRol(ctx context.Context, rolID int16) ([]string, error) {
	var codigos []string
	err := db(ctx, r.DB).Model(&PermisoModel{}).
		Joins("JOIN rol_permiso rp ON rp.permiso_id = permiso.id").
		Where("rp.rol_id = ?", rolID).Pluck("codigo", &codigos).Error
	return codigos, err
}

func uniqueInt16(in []int16) []int16 {
	visto := map[int16]bool{}
	var out []int16
	for _, v := range in {
		if !visto[v] {
			visto[v] = true
			out = append(out, v)
		}
	}
	return out
}
