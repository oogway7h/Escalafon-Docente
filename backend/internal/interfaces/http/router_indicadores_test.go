package http_test

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/domain/usuario"
)

// --- CU-25 tablero y CU-26 reportes ---

func TestIndicadores_SinSesion401YSinPermiso403(t *testing.T) {
	e := construirApp(t)
	for _, r := range []struct{ metodo, ruta string }{{http.MethodGet, "/api/tablero"}, {http.MethodPost, "/api/reportes"}} {
		if resp := peticionJSON(t, e.App, r.metodo, r.ruta, nil, ""); resp.StatusCode != 401 {
			t.Fatalf("%s sin sesión: esperaba 401, obtuve %d", r.ruta, resp.StatusCode)
		}
	}
	// Tiene la auditoría, pero no TABLERO_VER ni REPORTE_GENERAR.
	tok := auditorToken(t, e, rol.PermisoAuditoriaConsultar)
	if resp := peticionJSON(t, e.App, http.MethodGet, "/api/tablero", nil, tok); resp.StatusCode != 403 {
		t.Fatalf("tablero sin permiso: esperaba 403, obtuve %d", resp.StatusCode)
	}
	cuerpo := map[string]any{"tipo": "TRAMITES", "formato": "PDF"}
	if resp := peticionJSON(t, e.App, http.MethodPost, "/api/reportes", cuerpo, tok); resp.StatusCode != 403 {
		t.Fatalf("reporte sin permiso: esperaba 403, obtuve %d", resp.StatusCode)
	}
}

func TestTablero_200(t *testing.T) {
	e := construirApp(t)
	e.Indicadores.Conteos = []indicador.ConteoEstado{{Estado: tramite.EnRevision, Cantidad: 2}, {Estado: tramite.Archivado, Cantidad: 1}}
	e.Indicadores.Prestamos = []indicador.Prestamo{{LegajoCodigo: "LEG-00007", Salida: time.Now().AddDate(0, 0, -30)}}
	tok := auditorToken(t, e, rol.PermisoTableroVer)

	resp := peticionJSON(t, e.App, http.MethodGet, "/api/tablero?desde=2026-01-01&hasta=2026-12-31", nil, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
	var out struct {
		Periodo struct {
			Desde *string `json:"desde"`
			Hasta *string `json:"hasta"`
		} `json:"periodo"`
		Tramites struct {
			Radicados int64 `json:"radicados"`
			EnCurso   int64 `json:"enCurso"`
		} `json:"tramites"`
		PorEstado []struct {
			Estado   string `json:"estado"`
			Cantidad int64  `json:"cantidad"`
		} `json:"porEstado"`
		Tiempos []any `json:"tiempos"`
		Alertas []struct {
			LegajoCodigo string `json:"legajoCodigo"`
			DiasMora     int    `json:"diasMora"`
			Motivo       string `json:"motivo"`
		} `json:"alertas"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("respuesta JSON inválida: %v", err)
	}
	if out.Periodo.Desde == nil || *out.Periodo.Desde != "2026-01-01" || out.Tramites.Radicados != 3 || out.Tramites.EnCurso != 2 || len(out.PorEstado) != 9 {
		t.Fatalf("tablero inesperado: %+v", out)
	}
	if out.Tiempos == nil {
		t.Fatal("las listas vacías se envían como [] y no como null")
	}
	if len(out.Alertas) != 1 || out.Alertas[0].LegajoCodigo != "LEG-00007" || out.Alertas[0].DiasMora != 15 || out.Alertas[0].Motivo != "EXCEDE_MAXIMO" {
		t.Fatalf("alertas inesperadas: %+v", out.Alertas)
	}
}

func TestTablero_ValidaFechas422(t *testing.T) {
	e := construirApp(t)
	tok := auditorToken(t, e, rol.PermisoTableroVer)
	for consulta, codigo := range map[string]string{
		"desde=01/10/2026":                  "VALIDACION",
		"desde=2026-10-05&hasta=2026-10-01": "RANGO_FECHAS_INVALIDO",
	} {
		resp := peticionJSON(t, e.App, http.MethodGet, "/api/tablero?"+consulta, nil, tok)
		var cuerpo struct {
			Codigo string `json:"codigo"`
		}
		json.NewDecoder(resp.Body).Decode(&cuerpo)
		if resp.StatusCode != 422 || cuerpo.Codigo != codigo {
			t.Fatalf("%s: esperaba 422 %s, obtuve %d %s", consulta, codigo, resp.StatusCode, cuerpo.Codigo)
		}
	}
}

func TestReporte_DescargaYAuditoria(t *testing.T) {
	e := construirApp(t)
	e.Indicadores.Tramites = []indicador.FilaTramite{
		{Codigo: "ED-2026-AAAAAA", Tipo: tramite.Categorizacion, Estado: tramite.Ingresado, RadicadoEn: time.Now()},
	}
	u := e.Usuarios.Sembrar(usuario.Usuario{RolID: 1, Nombres: "Rosa", Apellidos: "Director", Correo: "rosa@uagrm.edu.bo", Activo: true})
	tok := loginToken(t, e.TokenSvc, u, rol.PermisoReporteGenerar)

	resp := peticionJSON(t, e.App, http.MethodPost, "/api/reportes", map[string]any{"tipo": "TRAMITES", "formato": "PDF"}, tok)
	if resp.StatusCode != 200 {
		t.Fatalf("esperaba 200, obtuve %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/pdf" {
		t.Fatalf("Content-Type inesperado: %q", ct)
	}
	if cd := resp.Header.Get("Content-Disposition"); !strings.HasPrefix(cd, `attachment; filename="reporte-tramites-`) || resp.Header.Get("X-Reporte-Filas") != "1" {
		t.Fatalf("cabeceras de descarga inesperadas: %q filas=%q", cd, resp.Header.Get("X-Reporte-Filas"))
	}
	if b, _ := io.ReadAll(resp.Body); string(b) != "%PDF-falso" {
		t.Fatalf("cuerpo inesperado: %q", b)
	}
	ultimo := e.Auditoria.Eventos[len(e.Auditoria.Eventos)-1]
	if ultimo.Accion != auditoria.AccionGenerarReporte || ultimo.UsuarioID == nil || *ultimo.UsuarioID != u.ID || ultimo.EntidadID != "TRAMITES" {
		t.Fatalf("la emisión debe quedar auditada: %+v", ultimo)
	}
}

func TestReporte_Errores422(t *testing.T) {
	e := construirApp(t)
	tok := auditorToken(t, e, rol.PermisoReporteGenerar)
	casos := []struct {
		cuerpo map[string]any
		codigo string
	}{
		{map[string]any{"tipo": "DOCENTES", "formato": "PDF"}, "VALIDACION"},
		{map[string]any{"tipo": "TRAMITES", "formato": "CSV"}, "VALIDACION"},
		{map[string]any{"tipo": "TRAMITES", "formato": "PDF", "desde": "ayer"}, "VALIDACION"},
		{map[string]any{"tipo": "TRAMITES", "formato": "PDF", "estado": "BORRADO"}, "FILTRO_INVALIDO"},
		{map[string]any{"tipo": "TRAMITES", "formato": "XLSX"}, "REPORTE_VACIO"}, // no hay trámites sembrados
	}
	for _, c := range casos {
		resp := peticionJSON(t, e.App, http.MethodPost, "/api/reportes", c.cuerpo, tok)
		var cuerpo struct {
			Codigo string `json:"codigo"`
		}
		json.NewDecoder(resp.Body).Decode(&cuerpo)
		if resp.StatusCode != 422 || cuerpo.Codigo != c.codigo {
			t.Fatalf("%v: esperaba 422 %s, obtuve %d %s", c.cuerpo, c.codigo, resp.StatusCode, cuerpo.Codigo)
		}
	}
}
