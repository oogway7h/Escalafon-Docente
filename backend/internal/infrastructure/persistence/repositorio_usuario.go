package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/usuario"
)

// RepositorioUsuario implementa usuario.Repository sobre GORM.
type RepositorioUsuario struct {
	DB *gorm.DB
}

func NewRepositorioUsuario(dbase *gorm.DB) *RepositorioUsuario { return &RepositorioUsuario{DB: dbase} }

func aDominioUsuario(m *UsuarioModel) *usuario.Usuario {
	return &usuario.Usuario{
		ID: m.ID, RolID: m.RolID, Nombres: m.Nombres, Apellidos: m.Apellidos, Correo: m.Correo,
		Telefono: m.Telefono, HashContrasena: m.HashContrasena, IntentosFallidos: m.IntentosFallidos,
		BloqueadoHasta: m.BloqueadoHasta, VersionSesion: m.VersionSesion, Activo: m.Activo,
		UltimoAcceso: m.UltimoAcceso, CreadoEn: m.CreadoEn, ActualizadoEn: m.ActualizadoEn,
		TerminosVersion: m.TerminosVersion, TerminosAceptadosEn: m.TerminosAceptadosEn,
	}
}

func aModeloUsuario(u *usuario.Usuario) *UsuarioModel {
	return &UsuarioModel{
		ID: u.ID, RolID: u.RolID, Nombres: u.Nombres, Apellidos: u.Apellidos, Correo: u.Correo,
		Telefono: u.Telefono, HashContrasena: u.HashContrasena, IntentosFallidos: u.IntentosFallidos,
		BloqueadoHasta: u.BloqueadoHasta, VersionSesion: u.VersionSesion, Activo: u.Activo,
		UltimoAcceso: u.UltimoAcceso, CreadoEn: u.CreadoEn, ActualizadoEn: u.ActualizadoEn,
		TerminosVersion: u.TerminosVersion, TerminosAceptadosEn: u.TerminosAceptadosEn,
	}
}

func (r *RepositorioUsuario) BuscarPorID(ctx context.Context, id string) (*usuario.Usuario, error) {
	var m UsuarioModel
	if err := db(ctx, r.DB).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usuario.ErrNoEncontrado
		}
		return nil, err
	}
	return aDominioUsuario(&m), nil
}

func (r *RepositorioUsuario) BuscarPorCorreo(ctx context.Context, correo string) (*usuario.Usuario, error) {
	var m UsuarioModel
	if err := db(ctx, r.DB).First(&m, "lower(correo) = lower(?)", correo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, usuario.ErrNoEncontrado
		}
		return nil, err
	}
	return aDominioUsuario(&m), nil
}

func (r *RepositorioUsuario) Listar(ctx context.Context, filtro usuario.ListarFiltro) ([]usuario.Usuario, int64, error) {
	q := db(ctx, r.DB).Model(&UsuarioModel{})
	if filtro.Buscar != "" {
		patron := "%" + escaparLike(filtro.Buscar) + "%"
		q = q.Where("nombres ILIKE ? OR apellidos ILIKE ? OR correo ILIKE ?", patron, patron, patron)
	}
	if filtro.RolID != nil {
		q = q.Where("rol_id = ?", *filtro.RolID)
	}
	if filtro.Activo != nil {
		q = q.Where("activo = ?", *filtro.Activo)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	pagina, tamano := filtro.Pagina, filtro.Tamano
	if pagina < 1 {
		pagina = 1
	}
	if tamano < 1 || tamano > 100 {
		tamano = 20
	}

	var modelos []UsuarioModel
	if err := q.Order("creado_en DESC").Offset((pagina - 1) * tamano).Limit(tamano).Find(&modelos).Error; err != nil {
		return nil, 0, err
	}
	out := make([]usuario.Usuario, len(modelos))
	for i := range modelos {
		out[i] = *aDominioUsuario(&modelos[i])
	}
	return out, total, nil
}

func (r *RepositorioUsuario) Crear(ctx context.Context, u *usuario.Usuario) error {
	m := aModeloUsuario(u)
	if err := db(ctx, r.DB).Create(m).Error; err != nil {
		if esViolacionUnicidad(err) {
			return usuario.ErrCorreoDuplicado
		}
		return err
	}
	*u = *aDominioUsuario(m)
	return nil
}

func (r *RepositorioUsuario) IncrementarVersionSesionPorRol(ctx context.Context, rolID int16) error {
	return db(ctx, r.DB).Model(&UsuarioModel{}).Where("rol_id = ?", rolID).
		UpdateColumn("version_sesion", gorm.Expr("version_sesion + 1")).Error
}

func (r *RepositorioUsuario) Actualizar(ctx context.Context, u *usuario.Usuario) error {
	m := aModeloUsuario(u)
	if err := db(ctx, r.DB).Model(&UsuarioModel{}).Where("id = ?", u.ID).Updates(map[string]any{
		"rol_id": m.RolID, "nombres": m.Nombres, "apellidos": m.Apellidos, "correo": m.Correo,
		"telefono": m.Telefono, "hash_contrasena": m.HashContrasena, "intentos_fallidos": m.IntentosFallidos,
		"bloqueado_hasta": m.BloqueadoHasta, "version_sesion": m.VersionSesion, "activo": m.Activo,
		"ultimo_acceso": m.UltimoAcceso, "terminos_version": m.TerminosVersion, "terminos_aceptados_en": m.TerminosAceptadosEn,
	}).Error; err != nil {
		if esViolacionUnicidad(err) {
			return usuario.ErrCorreoDuplicado
		}
		return err
	}
	return nil
}
