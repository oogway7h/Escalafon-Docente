package perfil_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/perfil"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/usuario"
)

func TestActualizarContacto_SoloCambiaTelefono(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	auditoria := &testdobles.AuditoriaRepoFalso{}
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	u := usuarios.Sembrar(usuario.Usuario{Correo: "ana@uagrm.edu.bo", Telefono: "70000000", Nombres: "Ana"})

	caso := &perfil.ActualizarContacto{Usuarios: usuarios, Auditoria: auditoria, Reloj: reloj}
	actualizado, err := caso.Ejecutar(context.Background(), u.ID, "71111111")
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if actualizado.Telefono != "71111111" || actualizado.Nombres != "Ana" {
		t.Fatalf("actualización inesperada: %+v", actualizado)
	}
	if len(auditoria.Eventos) != 1 || auditoria.Eventos[0].Accion != "ACTUALIZAR_CONTACTO" {
		t.Fatalf("esperaba auditar ACTUALIZAR_CONTACTO, obtuve %v", auditoria.Eventos)
	}
}

func TestCambiarContrasena_ExigeContrasenaActual(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	u := usuarios.Sembrar(usuario.Usuario{HashContrasena: "hash:actual", VersionSesion: 1})
	roles := testdobles.NewRolRepoFalso()
	roles.Sembrar(rol.Rol{ID: 1, Nombre: "Director"}, []string{"USUARIO_GESTIONAR"})

	caso := &perfil.CambiarContrasena{
		Usuarios: usuarios, Roles: roles, Auditoria: &testdobles.AuditoriaRepoFalso{},
		Hasher: testdobles.HasherFalso{}, Tokens: testdobles.NewTokenServiceFalso(),
		Cache: testdobles.NewSessionCacheFalso(), Tx: testdobles.TxManagerFalso{},
		Reloj: &testdobles.RelojFijo{Momento: time.Now()}, SesionHoras: 8,
	}
	_, err := caso.Ejecutar(context.Background(), u.ID, "incorrecta", "NuevaClave1")
	if !errors.Is(err, usuario.ErrContrasenaActualIncorrecta) {
		t.Fatalf("esperaba ErrContrasenaActualIncorrecta, obtuve %v", err)
	}
}

func TestCambiarContrasena_ExitoIncrementaVersionYEmiteToken(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	u := usuarios.Sembrar(usuario.Usuario{HashContrasena: "hash:actual", VersionSesion: 1, RolID: 1})
	roles := testdobles.NewRolRepoFalso()
	roles.Sembrar(rol.Rol{ID: 1, Nombre: "Director"}, []string{"USUARIO_GESTIONAR"})
	cache := testdobles.NewSessionCacheFalso()
	cache.GuardarSesion(context.Background(), u.ID, 1, true)

	caso := &perfil.CambiarContrasena{
		Usuarios: usuarios, Roles: roles, Auditoria: &testdobles.AuditoriaRepoFalso{},
		Hasher: testdobles.HasherFalso{}, Tokens: testdobles.NewTokenServiceFalso(),
		Cache: cache, Tx: testdobles.TxManagerFalso{},
		Reloj: &testdobles.RelojFijo{Momento: time.Now()}, SesionHoras: 8,
	}
	salida, err := caso.Ejecutar(context.Background(), u.ID, "actual", "NuevaClave1")
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if salida.Token == "" {
		t.Fatalf("esperaba un token nuevo")
	}
	actualizado, _ := usuarios.BuscarPorID(context.Background(), u.ID)
	if actualizado.VersionSesion != 2 || actualizado.HashContrasena != "hash:NuevaClave1" {
		t.Fatalf("no se actualizó correctamente: %+v", actualizado)
	}
	if _, _, encontrado, _ := cache.ObtenerSesion(context.Background(), u.ID); encontrado {
		t.Fatalf("la caché debería invalidarse tras el cambio de contraseña")
	}
}

func TestCambiarContrasena_RechazaPoliticaDebil(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	u := usuarios.Sembrar(usuario.Usuario{HashContrasena: "hash:actual", RolID: 1})
	roles := testdobles.NewRolRepoFalso()
	roles.Sembrar(rol.Rol{ID: 1, Nombre: "Director"}, []string{"USUARIO_GESTIONAR"})

	caso := &perfil.CambiarContrasena{
		Usuarios: usuarios, Roles: roles, Auditoria: &testdobles.AuditoriaRepoFalso{},
		Hasher: testdobles.HasherFalso{}, Tokens: testdobles.NewTokenServiceFalso(),
		Cache: testdobles.NewSessionCacheFalso(), Tx: testdobles.TxManagerFalso{},
		Reloj: &testdobles.RelojFijo{Momento: time.Now()}, SesionHoras: 8,
	}
	_, err := caso.Ejecutar(context.Background(), u.ID, "actual", "debil")
	if !errors.Is(err, usuario.ErrPoliticaContrasena) {
		t.Fatalf("esperaba ErrPoliticaContrasena, obtuve %v", err)
	}
}
