package roles_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/roles"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/rol"
)

func TestGuardar_PermisoInvalido(t *testing.T) {
	roles_ := testdobles.NewRolRepoFalso()
	roles_.Sembrar(rol.Rol{ID: 1, Nombre: "Director", EsSistema: true}, []string{rol.PermisoUsuarioGestionar, rol.PermisoRolGestionar})
	usuariosRepo := testdobles.NewUsuarioRepoFalso()
	caso := &roles.Guardar{Roles: roles_, Usuarios: usuariosRepo, Auditoria: &testdobles.AuditoriaRepoFalso{}, Cache: testdobles.NewSessionCacheFalso(), Tx: testdobles.TxManagerFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}

	_, err := caso.Ejecutar(context.Background(), roles.GuardarEntrada{Nombre: "Nuevo", PermisoIDs: []int16{999}})
	if !errors.Is(err, rol.ErrPermisoInvalido) {
		t.Fatalf("esperaba ErrPermisoInvalido, obtuve %v", err)
	}
}

func TestGuardar_RolSistemaNoSePuedeRenombrar(t *testing.T) {
	roles_ := testdobles.NewRolRepoFalso()
	roles_.Sembrar(rol.Rol{ID: 1, Nombre: "Director", EsSistema: true}, []string{rol.PermisoUsuarioGestionar, rol.PermisoRolGestionar})
	usuariosRepo := testdobles.NewUsuarioRepoFalso()
	caso := &roles.Guardar{Roles: roles_, Usuarios: usuariosRepo, Auditoria: &testdobles.AuditoriaRepoFalso{}, Cache: testdobles.NewSessionCacheFalso(), Tx: testdobles.TxManagerFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}

	_, err := caso.Ejecutar(context.Background(), roles.GuardarEntrada{ID: 1, Nombre: "Jefe"})
	if !errors.Is(err, rol.ErrSistemaProtegido) {
		t.Fatalf("esperaba ErrSistemaProtegido al renombrar un rol de sistema, obtuve %v", err)
	}
}

func TestGuardar_DirectorNoPierdePermisosProtegidos(t *testing.T) {
	roles_ := testdobles.NewRolRepoFalso()
	roles_.Sembrar(rol.Rol{ID: 1, Nombre: rol.NombreDirector, EsSistema: true}, []string{rol.PermisoUsuarioGestionar, rol.PermisoRolGestionar})
	usuariosRepo := testdobles.NewUsuarioRepoFalso()
	caso := &roles.Guardar{Roles: roles_, Usuarios: usuariosRepo, Auditoria: &testdobles.AuditoriaRepoFalso{}, Cache: testdobles.NewSessionCacheFalso(), Tx: testdobles.TxManagerFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}

	// Intenta guardar al Director sin ROL_GESTIONAR: debe rechazarse.
	_, err := caso.Ejecutar(context.Background(), roles.GuardarEntrada{ID: 1, Nombre: rol.NombreDirector, PermisoIDs: []int16{}})
	if !errors.Is(err, rol.ErrSistemaProtegido) {
		t.Fatalf("esperaba ErrSistemaProtegido al retirar permisos protegidos del Director, obtuve %v", err)
	}
}

func TestGuardar_ExitoIncrementaVersionEInvalidaCache(t *testing.T) {
	roles_ := testdobles.NewRolRepoFalso()
	roles_.Sembrar(rol.Rol{ID: 2, Nombre: "Analista"}, nil)
	usuariosRepo := testdobles.NewUsuarioRepoFalso()
	roles_.AsignarUsuarioARol(2, "usuario-x")
	cache := testdobles.NewSessionCacheFalso()
	cache.GuardarSesion(context.Background(), "usuario-x", 1, true)

	caso := &roles.Guardar{Roles: roles_, Usuarios: usuariosRepo, Auditoria: &testdobles.AuditoriaRepoFalso{}, Cache: cache, Tx: testdobles.TxManagerFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}
	_, err := caso.Ejecutar(context.Background(), roles.GuardarEntrada{ID: 2, Nombre: "Analista", PermisoIDs: nil})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if _, _, encontrado, _ := cache.ObtenerSesion(context.Background(), "usuario-x"); encontrado {
		t.Fatalf("la caché de los usuarios del rol debería invalidarse")
	}
}

func TestEliminar_RechazaRolDeSistema(t *testing.T) {
	roles_ := testdobles.NewRolRepoFalso()
	roles_.Sembrar(rol.Rol{ID: 1, Nombre: "Director", EsSistema: true}, nil)
	caso := &roles.Eliminar{Roles: roles_, Auditoria: &testdobles.AuditoriaRepoFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}
	if err := caso.Ejecutar(context.Background(), "actor-1", 1); !errors.Is(err, rol.ErrSistemaProtegido) {
		t.Fatalf("esperaba ErrSistemaProtegido, obtuve %v", err)
	}
}

func TestEliminar_RechazaRolConUsuarios(t *testing.T) {
	roles_ := testdobles.NewRolRepoFalso()
	roles_.Sembrar(rol.Rol{ID: 3, Nombre: "Temporal"}, nil)
	roles_.AsignarUsuarioARol(3, "usuario-y")
	caso := &roles.Eliminar{Roles: roles_, Auditoria: &testdobles.AuditoriaRepoFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}
	var errConUsuarios *rol.ErrConUsuarios
	err := caso.Ejecutar(context.Background(), "actor-1", 3)
	if !errors.As(err, &errConUsuarios) || errConUsuarios.Cantidad != 1 {
		t.Fatalf("esperaba ErrConUsuarios{Cantidad:1}, obtuve %v", err)
	}
}

func TestEliminar_Exito(t *testing.T) {
	roles_ := testdobles.NewRolRepoFalso()
	roles_.Sembrar(rol.Rol{ID: 4, Nombre: "Temporal2"}, nil)
	auditoria := &testdobles.AuditoriaRepoFalso{}
	caso := &roles.Eliminar{Roles: roles_, Auditoria: auditoria, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}
	if err := caso.Ejecutar(context.Background(), "actor-1", 4); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if len(auditoria.Eventos) != 1 || auditoria.Eventos[0].Accion != "ELIMINAR_ROL" {
		t.Fatalf("esperaba auditar ELIMINAR_ROL, obtuve %v", auditoria.Eventos)
	}
}
