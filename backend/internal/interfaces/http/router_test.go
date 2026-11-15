package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	auditoriasApp "escalafon/backend/internal/application/auditorias"
	"escalafon/backend/internal/application/auth"
	indicadoresApp "escalafon/backend/internal/application/indicadores"
	"escalafon/backend/internal/application/perfil"
	"escalafon/backend/internal/application/ports"
	rolesApp "escalafon/backend/internal/application/roles"
	"escalafon/backend/internal/application/testdobles"
	tramitesApp "escalafon/backend/internal/application/tramites"
	usuariosApp "escalafon/backend/internal/application/usuarios"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/usuario"
	httpapi "escalafon/backend/internal/interfaces/http"
	"escalafon/backend/internal/interfaces/http/handlers"
	"escalafon/backend/internal/interfaces/http/middleware"
)

type entorno struct {
	App       *fiber.App
	Usuarios  *testdobles.UsuarioRepoFalso
	Roles     *testdobles.RolRepoFalso
	Tokens    *testdobles.TokenRepoFalso
	Cache     *testdobles.SessionCacheFalso
	TokenSvc  *testdobles.TokenServiceFalso
	Reloj     *testdobles.RelojFijo
	Generador *testdobles.GeneradorFalso
	Mailer    *testdobles.MailerFalso
	Auditoria *testdobles.AuditoriaRepoFalso

	Indicadores *testdobles.IndicadorRepoFalso
	Tramites    *testdobles.TramiteRepoFalso
}

func construirApp(t *testing.T) *entorno {
	t.Helper()
	usuarios := testdobles.NewUsuarioRepoFalso()
	roles := testdobles.NewRolRepoFalso()
	roles.Sembrar(rol.Rol{ID: 1, Nombre: "Director", EsSistema: true, Activo: true}, []string{rol.PermisoPerfilGestionar, rol.PermisoUsuarioGestionar, rol.PermisoRolGestionar})
	tokens := testdobles.NewTokenRepoFalso()
	cache := testdobles.NewSessionCacheFalso()
	tokenSvc := testdobles.NewTokenServiceFalso()
	reloj := &testdobles.RelojFijo{Momento: time.Now()}
	generador := &testdobles.GeneradorFalso{}
	mailer := &testdobles.MailerFalso{}
	auditoria := &testdobles.AuditoriaRepoFalso{}
	tx := testdobles.TxManagerFalso{}
	hasher := testdobles.HasherFalso{}

	autenticarCaso := &auth.AutenticarUsuario{
		Usuarios: usuarios, Roles: roles, Auditoria: auditoria, Hasher: hasher, Tokens: tokenSvc, Reloj: reloj,
		MaxIntentos: 3, MinutosBloqueo: 15, SesionHoras: 8,
	}
	authHandler := &handlers.AuthHandler{
		Autenticar:        autenticarCaso,
		CerrarSesionCaso:  &auth.CerrarSesion{Cache: cache},
		ObtenerSesionCaso: &auth.ObtenerSesion{Usuarios: usuarios, Roles: roles},
		SolicitarRecuperacion: &auth.SolicitarRecuperacion{
			Usuarios: usuarios, Tokens: tokens, Generador: generador, Mailer: mailer, Reloj: reloj,
			MinutosVigencia: 30, URLFrontend: "http://localhost:3000",
		},
		RestablecerRecuperacion: &auth.RestablecerContrasena{
			Usuarios: usuarios, Tokens: tokens, Auditoria: auditoria, Hasher: hasher, Generador: generador,
			Cache: cache, Tx: tx, Reloj: reloj, PropositoEsperado: token.PropositoRecuperacion,
		},
		RestablecerActivacion: &auth.RestablecerContrasena{
			Usuarios: usuarios, Tokens: tokens, Auditoria: auditoria, Hasher: hasher, Generador: generador,
			Cache: cache, Tx: tx, Reloj: reloj, PropositoEsperado: token.PropositoActivacion,
		},
	}
	perfilHandler := &handlers.PerfilHandler{
		CasoObtener:            &perfil.ObtenerPerfil{Usuarios: usuarios},
		CasoActualizarContacto: &perfil.ActualizarContacto{Usuarios: usuarios, Auditoria: auditoria, Reloj: reloj},
		CasoCambiarContrasena: &perfil.CambiarContrasena{
			Usuarios: usuarios, Roles: roles, Auditoria: auditoria, Hasher: hasher, Tokens: tokenSvc,
			Cache: cache, Tx: tx, Reloj: reloj, SesionHoras: 8,
		},
	}

	usuariosHandler := &handlers.UsuariosHandler{
		CasoListar:  &usuariosApp.Listar{Usuarios: usuarios},
		CasoObtener: &usuariosApp.Obtener{Usuarios: usuarios},
		CasoCrear: &usuariosApp.Crear{
			Usuarios: usuarios, Roles: roles, Tokens: tokens, Auditoria: auditoria, Hasher: hasher,
			Generador: generador, Mailer: mailer, Tx: tx, Reloj: reloj, HorasTokenActivacion: 24, URLFrontend: "http://localhost:3000",
		},
		CasoEditar: &usuariosApp.Editar{Usuarios: usuarios, Auditoria: auditoria, Reloj: reloj},
		CasoCambiarRol: &usuariosApp.CambiarRol{
			Usuarios: usuarios, Roles: roles, Auditoria: auditoria, Cache: cache, Tx: tx, Reloj: reloj,
		},
		CasoCambiarEstado: &usuariosApp.CambiarEstado{
			Usuarios: usuarios, Auditoria: auditoria, Cache: cache, Tx: tx, Reloj: reloj,
		},
		CasoReenviarActivacion: &usuariosApp.ReenviarActivacion{
			Usuarios: usuarios, Tokens: tokens, Auditoria: auditoria, Generador: generador, Mailer: mailer, Reloj: reloj, HorasTokenActivacion: 24, URLFrontend: "http://localhost:3000",
		},
	}
	rolesHandler := &handlers.RolesHandler{
		CasoListar:         &rolesApp.Listar{Roles: roles},
		CasoObtener:        &rolesApp.Obtener{Roles: roles},
		CasoListarPermisos: &rolesApp.ListarPermisos{Roles: roles},
		CasoGuardar:        &rolesApp.Guardar{Roles: roles, Usuarios: usuarios, Auditoria: auditoria, Cache: cache, Tx: tx, Reloj: reloj},
		CasoEliminar:       &rolesApp.Eliminar{Roles: roles, Auditoria: auditoria, Reloj: reloj},
	}

	auditoriaHandler := &handlers.AuditoriaHandler{
		CasoConsultar: &auditoriasApp.Consultar{Auditoria: auditoria},
		CasoOpciones:  &auditoriasApp.ListarOpciones{Auditoria: auditoria},
	}

	indicadores := &testdobles.IndicadorRepoFalso{}
	indicadoresHandler := &handlers.IndicadoresHandler{
		CasoTablero: &indicadoresApp.ConsultarTablero{Indicadores: indicadores, Reloj: reloj, DiasMaximosPrestamo: 15},
		CasoReporte: &indicadoresApp.GenerarReporte{
			Indicadores: indicadores, Usuarios: usuarios, Exportador: &testdobles.ExportadorFalso{}, Auditoria: auditoria,
			Reloj: reloj, DiasMaximosPrestamo: 15,
		},
	}

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ManejadorErrores})
	tramitesRepo := testdobles.NewTramiteRepoFalso()
	documentosHandler := &handlers.DocumentosHandler{
		CasoBuscar:    &tramitesApp.BuscarDocumentos{Tramites: tramitesRepo},
		CasoBajas:     &tramitesApp.ListarDocumentosDeBaja{Tramites: tramitesRepo},
		CasoRestituir: &tramitesApp.RestituirDocumento{Tramites: tramitesRepo, Auditoria: auditoria, Reloj: reloj},
	}

	httpapi.Registrar(app, httpapi.Dependencias{
		Auth: authHandler, Perfil: perfilHandler,
		Usuarios: usuariosHandler, Roles: rolesHandler, Auditoria: auditoriaHandler, Indicadores: indicadoresHandler, Documentos: documentosHandler,
		Tokens: tokenSvc, Cache: cache, UsuarioRepo: usuarios,
	})

	return &entorno{App: app, Usuarios: usuarios, Roles: roles, Tokens: tokens, Cache: cache, TokenSvc: tokenSvc, Reloj: reloj, Generador: generador, Mailer: mailer, Auditoria: auditoria, Indicadores: indicadores, Tramites: tramitesRepo}
}

func peticionJSON(t *testing.T, app *fiber.App, metodo, ruta string, cuerpo any, token string) *http.Response {
	t.Helper()
	var lector io.Reader
	if cuerpo != nil {
		b, _ := json.Marshal(cuerpo)
		lector = bytes.NewReader(b)
	}
	req := httptest.NewRequest(metodo, ruta, lector)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("no esperaba error al hacer la petición: %v", err)
	}
	return resp
}

func reqCtx() context.Context { return context.Background() }

func loginToken(t *testing.T, tokenSvc *testdobles.TokenServiceFalso, u *usuario.Usuario, permisos ...string) string {
	t.Helper()
	tok, err := tokenSvc.Emitir(ports.ClaimsSesion{UsuarioID: u.ID, RolID: u.RolID, Version: u.VersionSesion, Permisos: permisos, JTI: "jti-" + u.ID}, time.Hour)
	if err != nil {
		t.Fatalf("no se pudo emitir el token: %v", err)
	}
	return tok
}

// --- 7.1 Autenticación ---

func TestLogin_Exito200(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "ana@uagrm.edu.bo", HashContrasena: "hash:Clave123", Activo: true})

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/login", map[string]string{"correo": u.Correo, "contrasena": "Clave123"}, "")
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

func TestLogin_CredencialesInvalidas401(t *testing.T) {
	e := construirApp(t)
	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/login", map[string]string{"correo": "nadie@uagrm.edu.bo", "contrasena": "x"}, "")
	if resp.StatusCode != 401 {
		t.Fatalf("esperaba 401, obtuve %d", resp.StatusCode)
	}
}

func TestLogin_CuentaBloqueada423(t *testing.T) {
	e := construirApp(t)
	futuro := e.Reloj.Momento.Add(10 * time.Minute)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "bloq@uagrm.edu.bo", HashContrasena: "hash:Clave123", Activo: true, BloqueadoHasta: &futuro})

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/login", map[string]string{"correo": u.Correo, "contrasena": "Clave123"}, "")
	if resp.StatusCode != 423 {
		t.Fatalf("esperaba 423, obtuve %d", resp.StatusCode)
	}
}

func TestLogin_ValidacionDeEntrada422(t *testing.T) {
	e := construirApp(t)
	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/login", map[string]string{"correo": "no-es-correo"}, "")
	if resp.StatusCode != 422 {
		t.Fatalf("esperaba 422, obtuve %d", resp.StatusCode)
	}
}

func TestRecuperacion_Siempre202(t *testing.T) {
	e := construirApp(t)
	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/recuperacion", map[string]string{"correo": "nadie@uagrm.edu.bo"}, "")
	if resp.StatusCode != 202 {
		t.Fatalf("esperaba 202, obtuve %d", resp.StatusCode)
	}
}

func TestRestablecer_Exito200(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "res@uagrm.edu.bo", VersionSesion: 1})
	plano, hash, _ := e.Generador.Generar()
	e.Tokens.Crear(reqCtx(), &token.TokenRecuperacion{UsuarioID: u.ID, Proposito: token.PropositoRecuperacion, TokenHash: hash, ExpiraEn: e.Reloj.Momento.Add(30 * time.Minute)})

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/restablecer", map[string]string{"token": plano, "nueva": "NuevaClave1"}, "")
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

func TestRestablecer_TokenInvalido400(t *testing.T) {
	e := construirApp(t)
	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/restablecer", map[string]string{"token": "no-existe", "nueva": "NuevaClave1"}, "")
	if resp.StatusCode != 400 {
		t.Fatalf("esperaba 400, obtuve %d", resp.StatusCode)
	}
}

func TestRestablecer_PoliticaContrasena422(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "res2@uagrm.edu.bo"})
	plano, hash, _ := e.Generador.Generar()
	e.Tokens.Crear(reqCtx(), &token.TokenRecuperacion{UsuarioID: u.ID, Proposito: token.PropositoRecuperacion, TokenHash: hash, ExpiraEn: e.Reloj.Momento.Add(30 * time.Minute)})

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/restablecer", map[string]string{"token": plano, "nueva": "debil"}, "")
	if resp.StatusCode != 422 {
		t.Fatalf("esperaba 422, obtuve %d", resp.StatusCode)
	}
}

func TestActivar_MismoContratoConPropositoActivacion(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "act@uagrm.edu.bo"})
	plano, hash, _ := e.Generador.Generar()
	e.Tokens.Crear(reqCtx(), &token.TokenRecuperacion{UsuarioID: u.ID, Proposito: token.PropositoActivacion, TokenHash: hash, ExpiraEn: e.Reloj.Momento.Add(24 * time.Hour)})

	// Un token de RECUPERACION no debe servir para /activar.
	otroPlano, otroHash, _ := e.Generador.Generar()
	e.Tokens.Crear(reqCtx(), &token.TokenRecuperacion{UsuarioID: u.ID, Proposito: token.PropositoRecuperacion, TokenHash: otroHash, ExpiraEn: e.Reloj.Momento.Add(time.Hour)})
	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/activar", map[string]string{"token": otroPlano, "nueva": "NuevaClave1"}, "")
	if resp.StatusCode != 400 {
		t.Fatalf("esperaba 400 al usar un token de recuperación en /activar, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPost, "/api/auth/activar", map[string]string{"token": plano, "nueva": "NuevaClave1"}, "")
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

// --- 7.2 Sesión y perfil ---

func TestLogout_204(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Activo: true})
	tok := loginToken(t, e.TokenSvc, u)

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/auth/logout", nil, tok)
	if resp.StatusCode != 204 {
		t.Fatalf("esperaba 204, obtuve %d", resp.StatusCode)
	}
}

func TestSesion_200(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "ses@uagrm.edu.bo", Activo: true})
	tok := loginToken(t, e.TokenSvc, u)

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/auth/sesion", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

func TestPerfil_ObtenerRequierePermiso(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "sinpermiso@uagrm.edu.bo", Activo: true})
	tok := loginToken(t, e.TokenSvc, u) // sin el permiso PERFIL_GESTIONAR en el JWT

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/perfil/", nil, tok)
	if resp.StatusCode != 403 {
		t.Fatalf("esperaba 403 sin el permiso, obtuve %d", resp.StatusCode)
	}

	tokConPermiso := loginToken(t, e.TokenSvc, u, rol.PermisoPerfilGestionar)
	resp = peticionJSON(t, e.App, http.MethodGet, "/api/perfil/", nil, tokConPermiso)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200 con el permiso, obtuve %d", resp.StatusCode)
	}
}

func TestPerfil_ActualizarContactoRechazaCambioDeRol(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "contacto@uagrm.edu.bo", Activo: true})
	tok := loginToken(t, e.TokenSvc, u, rol.PermisoPerfilGestionar)

	resp := peticionJSON(t, e.App, http.MethodPut, "/api/perfil/", map[string]any{"telefono": "7000", "rol": "Director"}, tok)
	if resp.StatusCode != 403 {
		t.Fatalf("esperaba 403 CAMPO_NO_EDITABLE, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPut, "/api/perfil/", map[string]any{"telefono": "70001111"}, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200 al editar solo el teléfono, obtuve %d", resp.StatusCode)
	}
}

func TestPerfil_CambiarContrasenaCodigos(t *testing.T) {
	e := construirApp(t)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "clave@uagrm.edu.bo", HashContrasena: "hash:Actual123", Activo: true})
	tok := loginToken(t, e.TokenSvc, u, rol.PermisoPerfilGestionar)

	resp := peticionJSON(t, e.App, http.MethodPut, "/api/perfil/contrasena", map[string]string{"actual": "Incorrecta1", "nueva": "NuevaClave1"}, tok)
	if resp.StatusCode != 422 {
		t.Fatalf("esperaba 422 CONTRASENA_ACTUAL_INCORRECTA, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPut, "/api/perfil/contrasena", map[string]string{"actual": "Actual123", "nueva": "debil"}, tok)
	if resp.StatusCode != 422 {
		t.Fatalf("esperaba 422 POLITICA_CONTRASENA, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPut, "/api/perfil/contrasena", map[string]string{"actual": "Actual123", "nueva": "NuevaClave1"}, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

// --- 7.3 Usuarios ---

func directorToken(t *testing.T, e *entorno) (*usuario.Usuario, string) {
	t.Helper()
	director := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "director@uagrm.edu.bo", Activo: true})
	tok := loginToken(t, e.TokenSvc, director, rol.PermisoUsuarioGestionar, rol.PermisoRolGestionar)
	return director, tok
}

func TestUsuarios_ListarYObtener(t *testing.T) {
	e := construirApp(t)
	_, tok := directorToken(t, e)
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 2, Correo: "ana@uagrm.edu.bo", Activo: true})

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/usuarios/", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200 al listar, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodGet, "/api/usuarios/"+u.ID, nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200 al obtener, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodGet, "/api/usuarios/no-existe", nil, tok)
	if resp.StatusCode != 404 {
		t.Fatalf("esperaba 404, obtuve %d", resp.StatusCode)
	}
}

func TestUsuarios_CrearCodigos(t *testing.T) {
	e := construirApp(t)
	_, tok := directorToken(t, e)

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/usuarios/", map[string]any{
		"nombres": "Ana", "apellidos": "Perez", "correo": "nueva@uagrm.edu.bo", "rolId": 99,
	}, tok)
	if resp.StatusCode != 422 {
		t.Fatalf("esperaba 422 ROL_INVALIDO, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPost, "/api/usuarios/", map[string]any{
		"nombres": "Ana", "apellidos": "Perez", "correo": "nueva@uagrm.edu.bo", "rolId": 1,
	}, tok)
	if resp.StatusCode != 201 {
		t.Fatalf("esperaba 201, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPost, "/api/usuarios/", map[string]any{
		"nombres": "Otra", "apellidos": "Persona", "correo": "NUEVA@UAGRM.EDU.BO", "rolId": 1,
	}, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 CORREO_DUPLICADO, obtuve %d", resp.StatusCode)
	}
}

func TestUsuarios_CambiarRolYEstadoCodigos(t *testing.T) {
	e := construirApp(t)
	director, tok := directorToken(t, e)
	objetivo := e.Usuarios.Sembrar(usuario.Usuario{RolID: 2, Correo: "obj@uagrm.edu.bo", Activo: true})

	// AUTOCAMBIO_ROL: el propio Director no puede cambiarse el rol.
	resp := peticionJSON(t, e.App, http.MethodPatch, "/api/usuarios/"+director.ID+"/rol", map[string]any{"rolId": 2}, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 AUTOCAMBIO_ROL, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPatch, "/api/usuarios/"+objetivo.ID+"/rol", map[string]any{"rolId": 1}, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}

	// AUTODESACTIVACION: el propio Director no puede desactivarse.
	activo := false
	resp = peticionJSON(t, e.App, http.MethodPatch, "/api/usuarios/"+director.ID+"/estado", map[string]any{"activo": activo}, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 AUTODESACTIVACION, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPatch, "/api/usuarios/"+objetivo.ID+"/estado", map[string]any{"activo": activo}, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
}

func TestUsuarios_ReenviarActivacionCodigos(t *testing.T) {
	e := construirApp(t)
	_, tok := directorToken(t, e)
	pendiente := e.Usuarios.Sembrar(usuario.Usuario{RolID: 2, Correo: "pend@uagrm.edu.bo", Activo: true})
	e.Tokens.Crear(reqCtx(), &token.TokenRecuperacion{UsuarioID: pendiente.ID, Proposito: token.PropositoActivacion, TokenHash: "hp1", ExpiraEn: e.Reloj.Momento.Add(time.Hour)})

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/usuarios/"+pendiente.ID+"/reenviar-activacion", nil, tok)
	if resp.StatusCode != 202 {
		t.Fatalf("esperaba 202, obtuve %d", resp.StatusCode)
	}

	activado := e.Usuarios.Sembrar(usuario.Usuario{RolID: 2, Correo: "ya-act@uagrm.edu.bo", Activo: true})
	usada := e.Reloj.Momento
	e.Tokens.Crear(reqCtx(), &token.TokenRecuperacion{UsuarioID: activado.ID, Proposito: token.PropositoActivacion, TokenHash: "hp2", ExpiraEn: e.Reloj.Momento.Add(time.Hour), UsadoEn: &usada})

	resp = peticionJSON(t, e.App, http.MethodPost, "/api/usuarios/"+activado.ID+"/reenviar-activacion", nil, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 CUENTA_YA_ACTIVADA, obtuve %d", resp.StatusCode)
	}
}

// --- 7.4 Roles y permisos ---

func TestRoles_ListarPermisosYRoles(t *testing.T) {
	e := construirApp(t)
	_, tok := directorToken(t, e)

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/permisos", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200 al listar permisos, obtuve %d", resp.StatusCode)
	}
	resp = peticionJSON(t, e.App, http.MethodGet, "/api/roles/", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200 al listar roles, obtuve %d", resp.StatusCode)
	}
	resp = peticionJSON(t, e.App, http.MethodGet, "/api/roles/1", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200 al obtener el rol con sus permisos, obtuve %d", resp.StatusCode)
	}
}

func TestRoles_GuardarCodigos(t *testing.T) {
	e := construirApp(t)
	_, tok := directorToken(t, e)

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/roles/", map[string]any{"nombre": "Nuevo Rol", "permisoIds": []int16{999}}, tok)
	if resp.StatusCode != 422 {
		t.Fatalf("esperaba 422 PERMISO_INVALIDO, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPost, "/api/roles/", map[string]any{"nombre": "Nuevo Rol"}, tok)
	if resp.StatusCode != 201 {
		t.Fatalf("esperaba 201, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPost, "/api/roles/", map[string]any{"nombre": "Nuevo Rol"}, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 ROL_DUPLICADO, obtuve %d", resp.StatusCode)
	}

	resp = peticionJSON(t, e.App, http.MethodPut, "/api/roles/1", map[string]any{"nombre": "Jefe Supremo"}, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 ROL_SISTEMA_PROTEGIDO al renombrar el rol de sistema, obtuve %d", resp.StatusCode)
	}
}

func TestRoles_EliminarCodigos(t *testing.T) {
	e := construirApp(t)
	_, tok := directorToken(t, e)

	// Rol de sistema: protegido.
	resp := peticionJSON(t, e.App, http.MethodDelete, "/api/roles/1", nil, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 ROL_SISTEMA_PROTEGIDO, obtuve %d", resp.StatusCode)
	}

	// Rol con usuarios asignados.
	e.Roles.Sembrar(rol.Rol{ID: 5, Nombre: "Con Usuarios"}, nil)
	e.Roles.AsignarUsuarioARol(5, "usuario-z")
	resp = peticionJSON(t, e.App, http.MethodDelete, "/api/roles/5", nil, tok)
	if resp.StatusCode != 409 {
		t.Fatalf("esperaba 409 ROL_CON_USUARIOS, obtuve %d", resp.StatusCode)
	}

	// Rol sin usuarios: se elimina.
	e.Roles.Sembrar(rol.Rol{ID: 6, Nombre: "Sin Usuarios"}, nil)
	resp = peticionJSON(t, e.App, http.MethodDelete, "/api/roles/6", nil, tok)
	if resp.StatusCode != 204 {
		t.Fatalf("esperaba 204, obtuve %d", resp.StatusCode)
	}
}
