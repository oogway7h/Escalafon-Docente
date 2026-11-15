package usuarios_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/application/usuarios"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/usuario"
)

func nuevoEntorno() (*testdobles.UsuarioRepoFalso, *testdobles.RolRepoFalso, *testdobles.TokenRepoFalso, *testdobles.AuditoriaRepoFalso, *testdobles.SessionCacheFalso, *testdobles.RelojFijo) {
	usuariosRepo := testdobles.NewUsuarioRepoFalso()
	roles := testdobles.NewRolRepoFalso()
	roles.Sembrar(rol.Rol{ID: 1, Nombre: "Director", Activo: true}, nil)
	roles.Sembrar(rol.Rol{ID: 2, Nombre: "Analista", Activo: true}, nil)
	tokens := testdobles.NewTokenRepoFalso()
	auditoria := &testdobles.AuditoriaRepoFalso{}
	cache := testdobles.NewSessionCacheFalso()
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	return usuariosRepo, roles, tokens, auditoria, cache, reloj
}

func TestCrear_RolInvalido(t *testing.T) {
	usuariosRepo, roles, tokens, auditoria, _, reloj := nuevoEntorno()
	caso := &usuarios.Crear{
		Usuarios: usuariosRepo, Roles: roles, Tokens: tokens, Auditoria: auditoria,
		Hasher: testdobles.HasherFalso{}, Generador: &testdobles.GeneradorFalso{}, Mailer: &testdobles.MailerFalso{},
		Tx: testdobles.TxManagerFalso{}, Reloj: reloj, HorasTokenActivacion: 24, URLFrontend: "http://localhost:3000",
	}
	_, err := caso.Ejecutar(context.Background(), usuarios.CrearEntrada{RolID: 99, Correo: "a@uagrm.edu.bo"})
	if !errors.Is(err, usuario.ErrRolInvalido) {
		t.Fatalf("esperaba ErrRolInvalido, obtuve %v", err)
	}
}

func TestCrear_ExitoCreaTokenYEnviaCorreo(t *testing.T) {
	usuariosRepo, roles, tokens, auditoria, _, reloj := nuevoEntorno()
	mailer := &testdobles.MailerFalso{}
	caso := &usuarios.Crear{
		Usuarios: usuariosRepo, Roles: roles, Tokens: tokens, Auditoria: auditoria,
		Hasher: testdobles.HasherFalso{}, Generador: &testdobles.GeneradorFalso{}, Mailer: mailer,
		Tx: testdobles.TxManagerFalso{}, Reloj: reloj, HorasTokenActivacion: 24, URLFrontend: "http://localhost:3000",
	}
	u, err := caso.Ejecutar(context.Background(), usuarios.CrearEntrada{ActorID: "director-1", RolID: 1, Nombres: "Ana", Apellidos: "Perez", Correo: "Ana@UAGRM.edu.bo"})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if u.Correo != "ana@uagrm.edu.bo" || !u.Activo {
		t.Fatalf("usuario creado inesperado: %+v", u)
	}
	if len(mailer.Enviados) != 1 || mailer.Enviados[0].Tipo != "ACTIVACION" {
		t.Fatalf("esperaba un correo de activación, obtuve %v", mailer.Enviados)
	}
	if len(auditoria.Eventos) != 1 || auditoria.Eventos[0].Accion != "CREAR" {
		t.Fatalf("esperaba auditar CREAR, obtuve %v", auditoria.Eventos)
	}
}

func TestCambiarRol_RechazaAutocambio(t *testing.T) {
	usuariosRepo, roles, _, auditoria, cache, reloj := nuevoEntorno()
	u := usuariosRepo.Sembrar(usuario.Usuario{RolID: 1, Activo: true})
	caso := &usuarios.CambiarRol{Usuarios: usuariosRepo, Roles: roles, Auditoria: auditoria, Cache: cache, Tx: testdobles.TxManagerFalso{}, Reloj: reloj}
	if err := caso.Ejecutar(context.Background(), u.ID, u.ID, 2); !errors.Is(err, usuario.ErrAutocambioRol) {
		t.Fatalf("esperaba ErrAutocambioRol, obtuve %v", err)
	}
}

func TestCambiarRol_ExitoIncrementaVersionEInvalidaCache(t *testing.T) {
	usuariosRepo, roles, _, auditoria, cache, reloj := nuevoEntorno()
	actor := usuariosRepo.Sembrar(usuario.Usuario{RolID: 1, Activo: true})
	objetivo := usuariosRepo.Sembrar(usuario.Usuario{RolID: 2, Activo: true, VersionSesion: 1})
	cache.GuardarSesion(context.Background(), objetivo.ID, 1, true)

	caso := &usuarios.CambiarRol{Usuarios: usuariosRepo, Roles: roles, Auditoria: auditoria, Cache: cache, Tx: testdobles.TxManagerFalso{}, Reloj: reloj}
	if err := caso.Ejecutar(context.Background(), actor.ID, objetivo.ID, 1); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	actualizado, _ := usuariosRepo.BuscarPorID(context.Background(), objetivo.ID)
	if actualizado.RolID != 1 || actualizado.VersionSesion != 2 {
		t.Fatalf("no se actualizó correctamente: %+v", actualizado)
	}
	if _, _, encontrado, _ := cache.ObtenerSesion(context.Background(), objetivo.ID); encontrado {
		t.Fatalf("la caché del usuario debería invalidarse")
	}
}

func TestCambiarEstado_RechazaAutodesactivacion(t *testing.T) {
	usuariosRepo, _, _, auditoria, cache, reloj := nuevoEntorno()
	u := usuariosRepo.Sembrar(usuario.Usuario{RolID: 1, Activo: true})
	caso := &usuarios.CambiarEstado{Usuarios: usuariosRepo, Auditoria: auditoria, Cache: cache, Tx: testdobles.TxManagerFalso{}, Reloj: reloj}
	if err := caso.Ejecutar(context.Background(), u.ID, u.ID, false); !errors.Is(err, usuario.ErrAutodesactivacion) {
		t.Fatalf("esperaba ErrAutodesactivacion, obtuve %v", err)
	}
}

func TestCambiarEstado_BajaEsLogica(t *testing.T) {
	usuariosRepo, _, _, auditoria, cache, reloj := nuevoEntorno()
	actor := usuariosRepo.Sembrar(usuario.Usuario{RolID: 1, Activo: true})
	objetivo := usuariosRepo.Sembrar(usuario.Usuario{RolID: 2, Activo: true})
	caso := &usuarios.CambiarEstado{Usuarios: usuariosRepo, Auditoria: auditoria, Cache: cache, Tx: testdobles.TxManagerFalso{}, Reloj: reloj}
	if err := caso.Ejecutar(context.Background(), actor.ID, objetivo.ID, false); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	actualizado, _ := usuariosRepo.BuscarPorID(context.Background(), objetivo.ID)
	if actualizado.Activo {
		t.Fatalf("debería quedar inactivo (baja lógica), no eliminado")
	}
	if auditoria.Eventos[len(auditoria.Eventos)-1].Accion != "BAJA" {
		t.Fatalf("esperaba auditar BAJA")
	}
}

func TestReenviarActivacion_RechazaSiYaActivada(t *testing.T) {
	usuariosRepo, _, tokens, auditoria, _, reloj := nuevoEntorno()
	u := usuariosRepo.Sembrar(usuario.Usuario{RolID: 1, Activo: true, Correo: "act@uagrm.edu.bo"})
	usado := reloj.Momento
	tokens.Crear(context.Background(), &token.TokenRecuperacion{UsuarioID: u.ID, Proposito: token.PropositoActivacion, TokenHash: "h1", ExpiraEn: reloj.Momento.Add(time.Hour), UsadoEn: &usado})

	caso := &usuarios.ReenviarActivacion{Usuarios: usuariosRepo, Tokens: tokens, Auditoria: auditoria, Generador: &testdobles.GeneradorFalso{}, Mailer: &testdobles.MailerFalso{}, Reloj: reloj, HorasTokenActivacion: 24, URLFrontend: "http://localhost:3000"}
	if err := caso.Ejecutar(context.Background(), "director-1", u.ID); !errors.Is(err, usuario.ErrCuentaYaActivada) {
		t.Fatalf("esperaba ErrCuentaYaActivada, obtuve %v", err)
	}
	if len(auditoria.Eventos) != 0 {
		t.Fatalf("un reenvío rechazado no debe auditarse, obtuve %v", auditoria.Eventos)
	}
}

func TestReenviarActivacion_ExitoEnviaNuevoCorreo(t *testing.T) {
	usuariosRepo, _, tokens, auditoria, _, reloj := nuevoEntorno()
	u := usuariosRepo.Sembrar(usuario.Usuario{RolID: 1, Activo: true, Correo: "pend@uagrm.edu.bo"})
	tokens.Crear(context.Background(), &token.TokenRecuperacion{UsuarioID: u.ID, Proposito: token.PropositoActivacion, TokenHash: "h2", ExpiraEn: reloj.Momento.Add(time.Hour)})

	mailer := &testdobles.MailerFalso{}
	caso := &usuarios.ReenviarActivacion{Usuarios: usuariosRepo, Tokens: tokens, Auditoria: auditoria, Generador: &testdobles.GeneradorFalso{}, Mailer: mailer, Reloj: reloj, HorasTokenActivacion: 24, URLFrontend: "http://localhost:3000"}
	if err := caso.Ejecutar(context.Background(), "director-1", u.ID); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if len(mailer.Enviados) != 1 {
		t.Fatalf("esperaba un correo de activación enviado")
	}
	// Se audita el reenvío con el actor y el usuario afectado, sin el token ni el enlace.
	if len(auditoria.Eventos) != 1 {
		t.Fatalf("esperaba un evento de auditoría, obtuve %v", auditoria.Eventos)
	}
	ev := auditoria.Eventos[0]
	if ev.Accion != "REENVIAR_ACTIVACION" || ev.UsuarioID == nil || *ev.UsuarioID != "director-1" || ev.EntidadID != u.ID || ev.Detalle != nil {
		t.Fatalf("evento de auditoría inesperado: %+v", ev)
	}
	// El token pendiente anterior debe quedar invalidado (marcado como usado), no eliminado.
	anterior, err := tokens.BuscarPorHash(context.Background(), "h2")
	if err != nil {
		t.Fatalf("no esperaba error al releer el token anterior: %v", err)
	}
	if anterior.UsadoEn == nil {
		t.Fatalf("el token de activación anterior debería quedar marcado como usado")
	}
}
