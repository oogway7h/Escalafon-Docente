package http_test

import (
	"net/http"
	"testing"

	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/tramite"
)

// CU-17 y CU-18: las rutas de /api/documentos exigen sesión y su permiso; la restitución es solo del Director.
func TestDocumentos_SinSesion401YSinPermiso403(t *testing.T) {
	e := construirApp(t)
	rutas := []struct{ metodo, ruta string }{
		{http.MethodGet, "/api/documentos/buscar?q=titulo"},
		{http.MethodGet, "/api/documentos/bajas"},
		{http.MethodPost, "/api/documentos/doc-1/restituir"},
	}
	for _, r := range rutas {
		if resp := peticionJSON(t, e.App, r.metodo, r.ruta, map[string]any{"motivo": "Fue un error"}, ""); resp.StatusCode != 401 {
			t.Fatalf("%s sin sesión: esperaba 401, obtuve %d", r.ruta, resp.StatusCode)
		}
	}
	// Quien busca (analista, archivo) no puede ver la bandeja de bajas ni restituir.
	tok := auditorToken(t, e, rol.PermisoDocumentoBuscar)
	if resp := peticionJSON(t, e.App, http.MethodGet, rutas[0].ruta, nil, tok); resp.StatusCode != 200 {
		t.Fatalf("buscar con DOCUMENTO_BUSCAR: esperaba 200, obtuve %d", resp.StatusCode)
	}
	for _, r := range rutas[1:] {
		if resp := peticionJSON(t, e.App, r.metodo, r.ruta, map[string]any{"motivo": "Fue un error"}, tok); resp.StatusCode != 403 {
			t.Fatalf("%s sin DOCUMENTO_RESTITUIR: esperaba 403, obtuve %d", r.ruta, resp.StatusCode)
		}
	}
}

func TestDocumentos_RestituirYValidaciones(t *testing.T) {
	e := construirApp(t)
	e.Tramites.Tramites["t1"] = &tramite.Tramite{ID: "t1", Estado: tramite.EnRevision, Documentos: []tramite.Documento{{ID: "doc-1", NombreOriginal: "titulo.pdf", MotivoBaja: "Duplicado"}}}
	tok := auditorToken(t, e, rol.PermisoDocumentoRestituir, rol.PermisoDocumentoBuscar)

	if resp := peticionJSON(t, e.App, http.MethodGet, "/api/documentos/buscar?q=a", nil, tok); resp.StatusCode != 422 {
		t.Fatalf("consulta de un carácter: esperaba 422, obtuve %d", resp.StatusCode)
	}
	if resp := peticionJSON(t, e.App, http.MethodPost, "/api/documentos/doc-1/restituir", map[string]any{"motivo": "no"}, tok); resp.StatusCode != 422 {
		t.Fatalf("motivo corto: esperaba 422, obtuve %d", resp.StatusCode)
	}
	if resp := peticionJSON(t, e.App, http.MethodPost, "/api/documentos/doc-1/restituir", map[string]any{"motivo": "La baja fue un error"}, tok); resp.StatusCode != 204 {
		t.Fatalf("restituir: esperaba 204, obtuve %d", resp.StatusCode)
	}
	if resp := peticionJSON(t, e.App, http.MethodPost, "/api/documentos/doc-1/restituir", map[string]any{"motivo": "La baja fue un error"}, tok); resp.StatusCode != 409 {
		t.Fatalf("restituir dos veces: esperaba 409, obtuve %d", resp.StatusCode)
	}
}
