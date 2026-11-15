package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/auth"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/usuario"
)

func nuevoCaso(t *testing.T, reloj *testdobles.RelojFijo) (*auth.AutenticarUsuario, *testdobles.UsuarioRepoFalso, *testdobles.RolRepoFalso, *testdobles.AuditoriaRepoFalso) {
	t.Helper()
	usuarios := testdobles.NewUsuarioRepoFalso()
	roles := testdobles.NewRolRepoFalso()
	roles.Sembrar(rol.Rol{ID: 1, Nombre: "Director", EsSistema: true}, []string{"USUARIO_GESTIONAR", "ROL_GESTIONAR"})
	auditoria := &testdobles.AuditoriaRepoFalso{}
	caso := &auth.AutenticarUsuario{
		Usuarios: usuarios, Roles: roles, Auditoria: auditoria,
		Hasher: testdobles.HasherFalso{}, Tokens: testdobles.NewTokenServiceFalso(), Reloj: reloj,
		MaxIntentos: 3, MinutosBloqueo: 15, SesionHoras: 8,
	}
	return caso, usuarios, roles, auditoria
}

func TestAutenticarUsuario_CredencialesInvalidasCuandoNoExiste(t *testing.T) {
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	caso, _, _, auditoria := nuevoCaso(t, reloj)

	_, err := caso.Ejecutar(context.Background(), auth.AutenticarEntrada{Correo: "nadie@uagrm.edu.bo", Contrasena: "x"})
	if !errors.Is(err, usuario.ErrCredencialesInvalidas) {
		t.Fatalf("esperaba ErrCredencialesInvalidas, obtuve %v", err)
	}
	if len(auditoria.Eventos) != 1 || auditoria.Eventos[0].Accion != "LOGIN_FALLIDO" {
		t.Fatalf("esperaba auditar LOGIN_FALLIDO, eventos=%v", auditoria.Eventos)
	}
}

func TestAutenticarUsuario_CredencialesInvalidasCuandoInactivo(t *testing.T) {
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	caso, usuarios, _, _ := nuevoCaso(t, reloj)
	usuarios.Sembrar(usuario.Usuario{Correo: "ana@uagrm.edu.bo", HashContrasena: "hash:clave", Activo: false, RolID: 1})

	_, err := caso.Ejecutar(context.Background(), auth.AutenticarEntrada{Correo: "ana@uagrm.edu.bo", Contrasena: "clave"})
	if !errors.Is(err, usuario.ErrCredencialesInvalidas) {
		t.Fatalf("esperaba ErrCredencialesInvalidas para usuario inactivo, obtuve %v", err)
	}
}

func TestAutenticarUsuario_BloqueaAlTercerIntentoFallido(t *testing.T) {
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	caso, usuarios, _, _ := nuevoCaso(t, reloj)
	u := usuarios.Sembrar(usuario.Usuario{Correo: "ana@uagrm.edu.bo", HashContrasena: "hash:clave", Activo: true, RolID: 1})

	for i := 0; i < 2; i++ {
		_, err := caso.Ejecutar(context.Background(), auth.AutenticarEntrada{Correo: u.Correo, Contrasena: "incorrecta"})
		if !errors.Is(err, usuario.ErrCredencialesInvalidas) {
			t.Fatalf("intento %d: esperaba ErrCredencialesInvalidas, obtuve %v", i+1, err)
		}
	}
	actualizado, _ := usuarios.BuscarPorID(context.Background(), u.ID)
	if actualizado.BloqueadoHasta != nil {
		t.Fatalf("no debería bloquearse todavía tras 2 intentos")
	}

	// Tercer intento: debe bloquear.
	_, err := caso.Ejecutar(context.Background(), auth.AutenticarEntrada{Correo: u.Correo, Contrasena: "incorrecta"})
	if !errors.Is(err, usuario.ErrCredencialesInvalidas) {
		t.Fatalf("esperaba ErrCredencialesInvalidas en el 3er intento, obtuve %v", err)
	}
	actualizado, _ = usuarios.BuscarPorID(context.Background(), u.ID)
	if actualizado.BloqueadoHasta == nil {
		t.Fatalf("debería quedar bloqueado tras 3 intentos fallidos")
	}
	esperado := reloj.Momento.Add(15 * time.Minute)
	if !actualizado.BloqueadoHasta.Equal(esperado) {
		t.Fatalf("bloqueadoHasta = %v, esperaba %v", actualizado.BloqueadoHasta, esperado)
	}

	// Un cuarto intento, ya bloqueado, debe devolver el error tipado con la hora.
	_, err = caso.Ejecutar(context.Background(), auth.AutenticarEntrada{Correo: u.Correo, Contrasena: "clave"})
	var errBloqueo *usuario.ErrCuentaBloqueada
	if !errors.As(err, &errBloqueo) {
		t.Fatalf("esperaba *usuario.ErrCuentaBloqueada, obtuve %v", err)
	}
}

func TestAutenticarUsuario_ExitoReiniciaIntentosYEmiteToken(t *testing.T) {
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	caso, usuarios, _, auditoriaRepo := nuevoCaso(t, reloj)
	u := usuarios.Sembrar(usuario.Usuario{
		Correo: "ana@uagrm.edu.bo", HashContrasena: "hash:clave", Activo: true, RolID: 1, IntentosFallidos: 2,
	})

	salida, err := caso.Ejecutar(context.Background(), auth.AutenticarEntrada{Correo: u.Correo, Contrasena: "clave"})
	if err != nil {
		t.Fatalf("no esperaba error, obtuve %v", err)
	}
	if salida.Token == "" {
		t.Fatalf("esperaba un token emitido")
	}
	if salida.ExpiraEn != 8*3600 {
		t.Fatalf("expiraEn = %d, esperaba %d", salida.ExpiraEn, 8*3600)
	}

	actualizado, _ := usuarios.BuscarPorID(context.Background(), u.ID)
	if actualizado.IntentosFallidos != 0 || actualizado.BloqueadoHasta != nil {
		t.Fatalf("debería reiniciar intentos y bloqueo tras un login exitoso")
	}
	if actualizado.UltimoAcceso == nil || !actualizado.UltimoAcceso.Equal(reloj.Momento) {
		t.Fatalf("debería registrar ultimoAcceso")
	}

	var vioLogin bool
	for _, e := range auditoriaRepo.Eventos {
		if e.Accion == "LOGIN" {
			vioLogin = true
		}
	}
	if !vioLogin {
		t.Fatalf("debería auditar LOGIN")
	}
}

func TestSolicitarRecuperacion_SiempreExitosaAunqueNoExistaElCorreo(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	tokens := testdobles.NewTokenRepoFalso()
	mailer := &testdobles.MailerFalso{}
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	caso := &auth.SolicitarRecuperacion{
		Usuarios: usuarios, Tokens: tokens, Generador: &testdobles.GeneradorFalso{}, Mailer: mailer,
		Reloj: reloj, MinutosVigencia: 30, URLFrontend: "http://localhost:3000",
	}

	if err := caso.Ejecutar(context.Background(), "nadie@uagrm.edu.bo"); err != nil {
		t.Fatalf("no debería fallar aunque el correo no exista: %v", err)
	}
	if len(mailer.Enviados) != 0 {
		t.Fatalf("no debería enviar correo si el usuario no existe")
	}

	u := usuarios.Sembrar(usuario.Usuario{Correo: "ana@uagrm.edu.bo", Activo: true})
	if err := caso.Ejecutar(context.Background(), u.Correo); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if len(mailer.Enviados) != 1 || mailer.Enviados[0].Tipo != "RECUPERACION" {
		t.Fatalf("esperaba un correo de recuperación enviado, obtuve %v", mailer.Enviados)
	}
}

func TestRestablecerContrasena_RechazaPoliticaDebil(t *testing.T) {
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	caso := &auth.RestablecerContrasena{
		Usuarios: testdobles.NewUsuarioRepoFalso(), Tokens: testdobles.NewTokenRepoFalso(),
		Auditoria: &testdobles.AuditoriaRepoFalso{}, Hasher: testdobles.HasherFalso{},
		Generador: &testdobles.GeneradorFalso{}, Cache: testdobles.NewSessionCacheFalso(),
		Tx: testdobles.TxManagerFalso{}, Reloj: reloj, PropositoEsperado: token.PropositoRecuperacion,
	}
	err := caso.Ejecutar(context.Background(), "cualquiera", "debil")
	if !errors.Is(err, usuario.ErrPoliticaContrasena) {
		t.Fatalf("esperaba ErrPoliticaContrasena, obtuve %v", err)
	}
}

func TestRestablecerContrasena_ExitoInvalidaTokenYCache(t *testing.T) {
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	usuarios := testdobles.NewUsuarioRepoFalso()
	tokens := testdobles.NewTokenRepoFalso()
	cache := testdobles.NewSessionCacheFalso()
	generador := &testdobles.GeneradorFalso{}

	u := usuarios.Sembrar(usuario.Usuario{Correo: "ana@uagrm.edu.bo", HashContrasena: "hash:vieja", VersionSesion: 1, Activo: true})
	plano, hash, _ := generador.Generar()
	tokens.Crear(context.Background(), &token.TokenRecuperacion{
		UsuarioID: u.ID, Proposito: token.PropositoRecuperacion, TokenHash: hash, ExpiraEn: reloj.Momento.Add(30 * time.Minute),
	})
	cache.GuardarSesion(context.Background(), u.ID, 1, true)

	caso := &auth.RestablecerContrasena{
		Usuarios: usuarios, Tokens: tokens, Auditoria: &testdobles.AuditoriaRepoFalso{}, Hasher: testdobles.HasherFalso{},
		Generador: generador, Cache: cache, Tx: testdobles.TxManagerFalso{}, Reloj: reloj,
		PropositoEsperado: token.PropositoRecuperacion,
	}
	if err := caso.Ejecutar(context.Background(), plano, "NuevaClave1"); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}

	actualizado, _ := usuarios.BuscarPorID(context.Background(), u.ID)
	if actualizado.HashContrasena != "hash:NuevaClave1" || actualizado.VersionSesion != 2 {
		t.Fatalf("no se actualizó correctamente el usuario: %+v", actualizado)
	}
	if _, _, encontrado, _ := cache.ObtenerSesion(context.Background(), u.ID); encontrado {
		t.Fatalf("la caché de sesión debería invalidarse tras restablecer la contraseña")
	}

	// Reutilizar el mismo token debe fallar (ya está usado).
	if err := caso.Ejecutar(context.Background(), plano, "OtraClave2"); !errors.Is(err, token.ErrInvalido) {
		t.Fatalf("esperaba ErrInvalido al reutilizar el token, obtuve %v", err)
	}
}

func TestRestablecerContrasena_RechazaPropositoIncorrecto(t *testing.T) {
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	usuarios := testdobles.NewUsuarioRepoFalso()
	tokens := testdobles.NewTokenRepoFalso()
	generador := &testdobles.GeneradorFalso{}
	u := usuarios.Sembrar(usuario.Usuario{Correo: "ana@uagrm.edu.bo", Activo: true})
	plano, hash, _ := generador.Generar()
	tokens.Crear(context.Background(), &token.TokenRecuperacion{
		UsuarioID: u.ID, Proposito: token.PropositoActivacion, TokenHash: hash, ExpiraEn: reloj.Momento.Add(time.Hour),
	})

	caso := &auth.RestablecerContrasena{
		Usuarios: usuarios, Tokens: tokens, Auditoria: &testdobles.AuditoriaRepoFalso{}, Hasher: testdobles.HasherFalso{},
		Generador: generador, Cache: testdobles.NewSessionCacheFalso(), Tx: testdobles.TxManagerFalso{}, Reloj: reloj,
		PropositoEsperado: token.PropositoRecuperacion, // se espera RECUPERACION, el token es de ACTIVACION
	}
	if err := caso.Ejecutar(context.Background(), plano, "NuevaClave1"); !errors.Is(err, token.ErrInvalido) {
		t.Fatalf("esperaba ErrInvalido por propósito incorrecto, obtuve %v", err)
	}
}

func TestCerrarSesion_RevocaJTI(t *testing.T) {
	cache := testdobles.NewSessionCacheFalso()
	caso := &auth.CerrarSesion{Cache: cache}
	if err := caso.Ejecutar(context.Background(), "jti-1", time.Hour); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	revocado, _ := cache.EstaRevocado(context.Background(), "jti-1")
	if !revocado {
		t.Fatalf("esperaba que el jti quedara revocado")
	}
}

// Con la cuenta bloqueada, una contraseña incorrecta responde igual que un correo inexistente: el bloqueo solo
// se le informa a quien conoce la contraseña.
func TestAutenticarUsuario_BloqueadaConClaveIncorrectaNoRevelaElBloqueo(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	hasta := time.Now().Add(10 * time.Minute)
	u := usuarios.Sembrar(usuario.Usuario{Correo: "bloq@uagrm.edu.bo", HashContrasena: "hash:clave", Activo: true, IntentosFallidos: 3, BloqueadoHasta: &hasta})
	caso := &auth.AutenticarUsuario{
		Usuarios: usuarios, Roles: testdobles.NewRolRepoFalso(), Auditoria: &testdobles.AuditoriaRepoFalso{},
		Hasher: &testdobles.HasherFalso{}, Tokens: testdobles.NewTokenServiceFalso(), Reloj: &testdobles.RelojFijo{Momento: time.Now()},
		MaxIntentos: 3, MinutosBloqueo: 15, SesionHoras: 8,
	}
	_, err := caso.Ejecutar(context.Background(), auth.AutenticarEntrada{Correo: u.Correo, Contrasena: "otra"})
	if !errors.Is(err, usuario.ErrCredencialesInvalidas) {
		t.Fatalf("esperaba ErrCredencialesInvalidas, obtuve %v", err)
	}
}
