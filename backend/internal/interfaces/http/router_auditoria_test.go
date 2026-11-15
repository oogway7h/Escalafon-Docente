package http_test

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/usuario"
)

// --- 7.5 Auditoría (CU-06) ---

const (
	idAna  = "11111111-1111-4111-8111-111111111111"
	idLuis = "22222222-2222-4222-8222-222222222222"
)

var bolivia = time.FixedZone("BOT", -4*60*60)

type listaAuditoria struct {
	Datos []struct {
		ID      int64     `json:"id"`
		Fecha   time.Time `json:"fecha"`
		Usuario *struct {
			ID      string `json:"id"`
			Nombres string `json:"nombres"`
			Correo  string `json:"correo"`
		} `json:"usuario"`
		Accion    string         `json:"accion"`
		Entidad   string         `json:"entidad"`
		EntidadID string         `json:"entidadId"`
		Detalle   map[string]any `json:"detalle"`
	} `json:"datos"`
	Total  int64 `json:"total"`
	Pagina int   `json:"pagina"`
	Tamano int   `json:"tamano"`
}

// sembrarBitacora deja cuatro eventos. El primero ocurre el 1/10 a las 23:30 de Bolivia, que en
// UTC ya es el 2/10: sirve para comprobar que los filtros de fecha usan el día boliviano.
func sembrarBitacora(e *entorno) {
	e.Auditoria.Autores = map[string]auditoria.Autor{
		idAna:  {ID: idAna, Nombres: "Ana", Apellidos: "Pérez", Correo: "ana@uagrm.edu.bo"},
		idLuis: {ID: idLuis, Nombres: "Luis", Apellidos: "Rojas", Correo: "luis@uagrm.edu.bo"},
	}
	ana, luis := idAna, idLuis
	for _, ev := range []auditoria.Evento{
		{UsuarioID: &ana, Accion: auditoria.AccionLogin, Entidad: auditoria.EntidadUsuario, EntidadID: ana, Fecha: time.Date(2026, 10, 1, 23, 30, 0, 0, bolivia)},
		{UsuarioID: &ana, Accion: auditoria.AccionCrear, Entidad: auditoria.EntidadDocente, EntidadID: "doc-1", Fecha: time.Date(2026, 10, 2, 10, 0, 0, 0, bolivia)},
		{UsuarioID: &luis, Accion: auditoria.AccionLogin, Entidad: auditoria.EntidadUsuario, EntidadID: luis, Fecha: time.Date(2026, 10, 3, 8, 0, 0, 0, bolivia)},
		{Accion: auditoria.AccionLoginFallido, Entidad: auditoria.EntidadUsuario, Fecha: time.Date(2026, 10, 4, 8, 0, 0, 0, bolivia)},
	} {
		e.Auditoria.Registrar(reqCtx(), &ev)
	}
}

func auditorToken(t *testing.T, e *entorno, permisos ...string) string {
	t.Helper()
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Correo: "auditor@uagrm.edu.bo", Activo: true})
	return loginToken(t, e.TokenSvc, u, permisos...)
}

func leerLista(t *testing.T, resp *http.Response) listaAuditoria {
	t.Helper()
	var out listaAuditoria
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	return out
}

func TestAuditoria_SinSesion401(t *testing.T) {
	e := construirApp(t)
	resp := peticionJSON(t, e.App, http.MethodGet, "/api/auditoria", nil, "")
	if resp.StatusCode != 401 {
		t.Fatalf("esperaba 401 sin sesión, obtuve %d", resp.StatusCode)
	}
}

func TestAuditoria_SinPermiso403(t *testing.T) {
	e := construirApp(t)
	// Tiene los permisos de gestión de usuarios y roles, pero no AUDITORIA_CONSULTAR.
	tok := auditorToken(t, e, rol.PermisoPerfilGestionar, rol.PermisoUsuarioGestionar, rol.PermisoRolGestionar)
	for _, ruta := range []string{"/api/auditoria", "/api/auditoria/filtros"} {
		resp := peticionJSON(t, e.App, http.MethodGet, ruta, nil, tok)
		if resp.StatusCode != 403 {
			t.Fatalf("%s: esperaba 403 PERMISO_DENEGADO, obtuve %d", ruta, resp.StatusCode)
		}
	}
}

func TestAuditoria_DirectorConsulta200(t *testing.T) {
	e := construirApp(t)
	sembrarBitacora(e)
	tok := auditorToken(t, e, rol.PermisoAuditoriaConsultar)

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/auditoria", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
	lista := leerLista(t, resp)
	if lista.Total != 4 || len(lista.Datos) != 4 || lista.Pagina != 1 || lista.Tamano != 20 {
		t.Fatalf("respuesta inesperada: %+v", lista)
	}
	primero, ultimo := lista.Datos[0], lista.Datos[3]
	if primero.Accion != auditoria.AccionLoginFallido || primero.Usuario != nil {
		t.Fatalf("el más reciente debería ser el login fallido sin usuario: %+v", primero)
	}
	if ultimo.Accion != auditoria.AccionLogin || ultimo.Usuario == nil || ultimo.Usuario.Nombres != "Ana" || ultimo.EntidadID != idAna {
		t.Fatalf("el más antiguo debería ser el login de Ana: %+v", ultimo)
	}
}

func TestAuditoria_FiltrosSeparadosYCombinados(t *testing.T) {
	e := construirApp(t)
	sembrarBitacora(e)
	tok := auditorToken(t, e, rol.PermisoAuditoriaConsultar)

	casos := []struct {
		consulta string
		total    int64
	}{
		{"usuarioId=" + idAna, 2},
		{"accion=LOGIN", 2},
		{"accion=login", 2}, // se normaliza a mayúsculas
		{"desde=2026-10-02", 3},
		{"hasta=2026-10-01", 1}, // 23:30 de Bolivia del 1/10 cuenta como ese día
		{"desde=2026-10-02&hasta=2026-10-03", 2},
		{"desde=2026-10-01&hasta=2026-10-01", 1},
		{"usuarioId=" + idAna + "&accion=CREAR", 1},
		{"usuarioId=" + idLuis + "&accion=LOGIN&desde=2026-10-01&hasta=2026-10-03", 1},
		{"usuarioId=" + idLuis + "&accion=CREAR", 0},
		{"pagina=2&tamano=3", 4},
	}
	for _, c := range casos {
		resp := peticionJSON(t, e.App, http.MethodGet, "/api/auditoria?"+c.consulta, nil, tok)
		if resp.StatusCode != 200 {
			t.Fatalf("%s: esperaba 200, obtuve %d", c.consulta, resp.StatusCode)
		}
		if lista := leerLista(t, resp); lista.Total != c.total {
			t.Fatalf("%s: esperaba %d registros, obtuve %d", c.consulta, c.total, lista.Total)
		}
	}

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/auditoria?pagina=2&tamano=3", nil, tok)
	if lista := leerLista(t, resp); len(lista.Datos) != 1 || lista.Pagina != 2 || lista.Tamano != 3 {
		t.Fatalf("la segunda página de 3 debería tener 1 registro: %+v", lista)
	}
}

func TestAuditoria_ValidacionDeParametros422(t *testing.T) {
	e := construirApp(t)
	sembrarBitacora(e)
	tok := auditorToken(t, e, rol.PermisoAuditoriaConsultar)

	casos := []struct{ consulta, codigo string }{
		{"desde=2026-13-01", "VALIDACION"},
		{"hasta=01/10/2026", "VALIDACION"},
		{"usuarioId=123", "VALIDACION"},
		{"pagina=uno", "VALIDACION"},
		{"accion=BORRAR_TODO", "ACCION_INVALIDA"},
		{"desde=2026-10-05&hasta=2026-10-01", "RANGO_FECHAS_INVALIDO"},
	}
	for _, c := range casos {
		resp := peticionJSON(t, e.App, http.MethodGet, "/api/auditoria?"+c.consulta, nil, tok)
		if resp.StatusCode != 422 {
			t.Fatalf("%s: esperaba 422, obtuve %d", c.consulta, resp.StatusCode)
		}
		var cuerpo struct {
			Codigo string `json:"codigo"`
		}
		json.NewDecoder(resp.Body).Decode(&cuerpo)
		if cuerpo.Codigo != c.codigo {
			t.Fatalf("%s: esperaba el código %s, obtuve %s", c.consulta, c.codigo, cuerpo.Codigo)
		}
	}
}

func TestAuditoria_SoloLectura(t *testing.T) {
	e := construirApp(t)
	sembrarBitacora(e)
	tok := auditorToken(t, e, rol.PermisoAuditoriaConsultar)
	antes := len(e.Auditoria.Eventos)

	for _, metodo := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		for _, ruta := range []string{"/api/auditoria", "/api/auditoria/1"} {
			resp := peticionJSON(t, e.App, metodo, ruta, map[string]any{"accion": "LOGIN"}, tok)
			if resp.StatusCode < 400 {
				t.Fatalf("%s %s no debería existir, obtuve %d", metodo, ruta, resp.StatusCode)
			}
		}
	}
	if len(e.Auditoria.Eventos) != antes {
		t.Fatalf("la bitácora cambió: antes %d eventos, ahora %d", antes, len(e.Auditoria.Eventos))
	}
}

func TestAuditoria_Filtros200(t *testing.T) {
	e := construirApp(t)
	sembrarBitacora(e)
	tok := auditorToken(t, e, rol.PermisoAuditoriaConsultar)

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/auditoria/filtros", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
	var out struct {
		Acciones []string `json:"acciones"`
		Usuarios []struct {
			ID string `json:"id"`
		} `json:"usuarios"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if len(out.Acciones) != len(auditoria.Acciones) || len(out.Usuarios) != 2 || out.Usuarios[0].ID != idAna {
		t.Fatalf("opciones inesperadas: %+v", out)
	}
}

func TestReenviarActivacion_QuedaAuditado(t *testing.T) {
	e := construirApp(t)
	director, tok := directorToken(t, e)
	pendiente := e.Usuarios.Sembrar(usuario.Usuario{RolID: 2, Correo: "pend2@uagrm.edu.bo", Activo: true})

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/usuarios/"+pendiente.ID+"/reenviar-activacion", nil, tok)
	if resp.StatusCode != 202 {
		t.Fatalf("esperaba 202, obtuve %d", resp.StatusCode)
	}
	ultimo := e.Auditoria.Eventos[len(e.Auditoria.Eventos)-1]
	if ultimo.Accion != auditoria.AccionReenviarActivacion || ultimo.UsuarioID == nil || *ultimo.UsuarioID != director.ID || ultimo.EntidadID != pendiente.ID {
		t.Fatalf("evento inesperado: %+v", ultimo)
	}
}
