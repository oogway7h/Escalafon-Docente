// Package roles implementa los casos de uso de HU-02: Roles y Permisos.
package roles

import (
	"context"
	"strconv"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/usuario"
)

type Listar struct{ Roles rol.Repository }

func (c *Listar) Ejecutar(ctx context.Context) ([]rol.ConConteo, error) { return c.Roles.Listar(ctx) }

type Obtener struct{ Roles rol.Repository }

func (c *Obtener) Ejecutar(ctx context.Context, id int16) (*rol.Rol, error) {
	return c.Roles.BuscarPorID(ctx, id)
}

type ListarPermisos struct{ Roles rol.Repository }

func (c *ListarPermisos) Ejecutar(ctx context.Context) ([]rol.Permiso, error) {
	return c.Roles.ListarPermisos(ctx)
}

// Guardar implementa GuardarRol (sección 8): crea o actualiza un rol y reemplaza sus permisos.
// Si el rol es de sistema no se puede renombrar, y al Director no se le pueden retirar
// USUARIO_GESTIONAR ni ROL_GESTIONAR (sección 4.4). Incrementa version_sesion de los usuarios
// del rol para que el cambio de permisos aplique en su siguiente inicio de sesión.
type Guardar struct {
	Roles     rol.Repository
	Usuarios  usuario.Repository
	Auditoria auditoria.Repository
	Cache     ports.SessionCache
	Tx        ports.TxManager
	Reloj     ports.Clock
}

type GuardarEntrada struct {
	ActorID     string
	ID          int16 // 0 = crear
	Nombre      string
	Descripcion string
	PermisoIDs  []int16
}

func (c *Guardar) Ejecutar(ctx context.Context, in GuardarEntrada) (*rol.Rol, error) {
	todosPermisos, err := c.Roles.ListarPermisos(ctx)
	if err != nil {
		return nil, err
	}
	codigoPorID := map[int16]string{}
	idPorCodigo := map[string]int16{}
	for _, p := range todosPermisos {
		codigoPorID[p.ID] = p.Codigo
		idPorCodigo[p.Codigo] = p.ID
	}
	for _, id := range in.PermisoIDs {
		if _, ok := codigoPorID[id]; !ok {
			return nil, rol.ErrPermisoInvalido
		}
	}

	r := &rol.Rol{ID: in.ID, Nombre: in.Nombre, Descripcion: in.Descripcion}
	var codigosAnteriores map[string]bool
	var idUsuariosDelRol []string

	if in.ID != 0 {
		existente, err := c.Roles.BuscarPorID(ctx, in.ID)
		if err != nil {
			return nil, err
		}
		r.EsSistema = existente.EsSistema
		if existente.EsSistema {
			if existente.Nombre != in.Nombre {
				return nil, rol.ErrSistemaProtegido
			}
			if existente.Nombre == rol.NombreDirector {
				presentes := map[int16]bool{}
				for _, id := range in.PermisoIDs {
					presentes[id] = true
				}
				for _, codigo := range rol.PermisosProtegidosDirector {
					if !presentes[idPorCodigo[codigo]] {
						return nil, rol.ErrSistemaProtegido
					}
				}
			}
		}

		codigosAnteriores = map[string]bool{}
		for _, p := range existente.Permisos {
			codigosAnteriores[p.Codigo] = true
		}
		idUsuariosDelRol, err = c.Roles.IDsUsuariosDelRol(ctx, in.ID)
		if err != nil {
			return nil, err
		}
	}

	codigosNuevos := map[string]bool{}
	for _, id := range in.PermisoIDs {
		codigosNuevos[codigoPorID[id]] = true
	}
	var agregados, retirados []string
	for cod := range codigosNuevos {
		if !codigosAnteriores[cod] {
			agregados = append(agregados, cod)
		}
	}
	for cod := range codigosAnteriores {
		if !codigosNuevos[cod] {
			retirados = append(retirados, cod)
		}
	}

	ahora := c.Reloj.Ahora()
	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := c.Roles.Guardar(ctx, r, in.PermisoIDs); err != nil {
			return err
		}
		if in.ID != 0 && len(idUsuariosDelRol) > 0 {
			if err := c.Usuarios.IncrementarVersionSesionPorRol(ctx, in.ID); err != nil {
				return err
			}
		}
		return c.Auditoria.Registrar(ctx, &auditoria.Evento{
			UsuarioID: &in.ActorID, Accion: auditoria.AccionGuardarRol, Entidad: auditoria.EntidadRol,
			EntidadID: strconv.Itoa(int(r.ID)), Detalle: map[string]any{"agregados": agregados, "retirados": retirados}, Fecha: ahora,
		})
	})
	if err != nil {
		return nil, err
	}

	for _, uid := range idUsuariosDelRol {
		_ = c.Cache.InvalidarSesion(ctx, uid)
	}
	return r, nil
}

// Eliminar rechaza roles de sistema o con usuarios asignados (incluidos inactivos); la FK
// RESTRICT en la base de datos es la segunda barrera (sección 8).
type Eliminar struct {
	Roles     rol.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *Eliminar) Ejecutar(ctx context.Context, actorID string, id int16) error {
	r, err := c.Roles.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if r.EsSistema {
		return rol.ErrSistemaProtegido
	}
	if err := c.Roles.Eliminar(ctx, id); err != nil {
		return err
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actorID, Accion: auditoria.AccionEliminarRol, Entidad: auditoria.EntidadRol,
		EntidadID: strconv.Itoa(int(id)), Fecha: c.Reloj.Ahora(),
	})
	return nil
}
