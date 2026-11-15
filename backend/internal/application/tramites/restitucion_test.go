package tramites_test

import (
	"context"
	"errors"
	"testing"

	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/tramite"
)

func TestRestituirDocumento_VuelveALasBusquedasYQuedaEnAuditoria(t *testing.T) {
	e := nuevoEntorno()
	tr, doc := e.procesado(t, "Resolución rectoral de ascenso")
	restituir := &tramites.RestituirDocumento{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	buscar := &tramites.BuscarDocumentos{Tramites: e.repo}
	ctx := context.Background()

	if err := restituir.Ejecutar(ctx, director, doc, "Fue un error"); !errors.Is(err, tramite.ErrDocumentoActivo) {
		t.Fatalf("un documento activo no se restituye, obtuve %v", err)
	}
	if err := e.darDeBaja(tr.ID, doc, "Cargado por error"); err != nil {
		t.Fatalf("dar de baja: %v", err)
	}
	if r, _, _ := buscar.Ejecutar(ctx, tramite.FiltroBusqueda{Consulta: "rectoral"}, 1, 10); len(r) != 0 {
		t.Fatalf("un documento de baja no sale en la búsqueda: %+v", r)
	}
	if err := restituir.Ejecutar(ctx, director, doc, "  "); !errors.Is(err, tramite.ErrMotivoRestitucion) {
		t.Fatalf("esperaba ErrMotivoRestitucion, obtuve %v", err)
	}
	if err := restituir.Ejecutar(ctx, director, doc, "La baja fue un error"); err != nil {
		t.Fatalf("restituir: %v", err)
	}
	d, _ := e.repo.BuscarDocumento(ctx, tr.ID, doc)
	if !d.Activo || d.MotivoBaja != "" || d.BajaEn != nil {
		t.Fatalf("el documento debe quedar activo y sin datos de baja: %+v", d)
	}
	if r, _, _ := buscar.Ejecutar(ctx, tramite.FiltroBusqueda{Consulta: "rectoral"}, 1, 10); len(r) != 1 || r[0].Paginas[0] != 1 {
		t.Fatalf("restituido vuelve a la búsqueda: %+v", r)
	}
	ultimo := e.aud.Eventos[len(e.aud.Eventos)-1]
	if ultimo.Accion != auditoria.AccionRestituirDocumento || ultimo.EntidadID != tr.ID || ultimo.Detalle["motivo"] != "La baja fue un error" {
		t.Fatalf("la restitución queda en auditoría: %+v", ultimo)
	}
}

func TestRestituirDocumento_TramiteArchivadoEsSoloLectura(t *testing.T) {
	e := nuevoEntorno()
	tr, doc := e.procesado(t, "Título")
	if err := e.darDeBaja(tr.ID, doc, "Duplicado"); err != nil {
		t.Fatalf("dar de baja: %v", err)
	}
	e.repo.Tramites[tr.ID].Estado = tramite.Archivado
	err := (&tramites.RestituirDocumento{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), director, doc, "Corrección")
	if !errors.Is(err, tramite.ErrArchivado) {
		t.Fatalf("esperaba ErrArchivado, obtuve %v", err)
	}
}

func TestBuscarDocumentos_ValidaLaConsulta(t *testing.T) {
	buscar := &tramites.BuscarDocumentos{Tramites: nuevoEntorno().repo}
	for _, f := range []tramite.FiltroBusqueda{{Consulta: "a"}, {Consulta: "título", Campo: "OTRO"}, {Consulta: "título", Tipo: "X"}} {
		if _, _, err := buscar.Ejecutar(context.Background(), f, 1, 10); err == nil {
			t.Fatalf("la búsqueda %+v debía rechazarse", f)
		}
	}
}
