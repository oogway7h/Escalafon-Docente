package notificaciones_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/notificaciones"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/notificacion"
	"escalafon/backend/internal/domain/tramite"
)

func armar(correo string) (*notificaciones.Notificar, *testdobles.MailerFalso, *testdobles.NotificacionRepoFalso) {
	docentes := testdobles.NewDocenteRepoFalso()
	docentes.Sembrar(docente.Docente{ID: "d1", Nombres: "Ana", Apellidos: "Paz", Correo: correo, Activo: true})
	mailer, repo := &testdobles.MailerFalso{}, &testdobles.NotificacionRepoFalso{}
	n := &notificaciones.Notificar{
		Docentes: docentes, Envios: repo, Mailer: mailer, Reloj: &testdobles.RelojFijo{Momento: time.Now()}, URLFrontend: "http://app/",
	}
	return n, mailer, repo
}

func tr(estado tramite.Estado) *tramite.Tramite {
	return &tramite.Tramite{ID: "t1", Codigo: "TR-2026-ABC", DocenteID: "d1", Estado: estado}
}

func TestEnviaAvisoEnLosEstadosNotificables(t *testing.T) {
	for _, e := range []tramite.Estado{tramite.Observado, tramite.EnFacultad, tramite.Remitido, tramite.Archivado} {
		n, mailer, repo := armar("ana@x.bo")
		n.Transicion(context.Background(), tr(e))
		if len(mailer.Enviados) != 1 || repo.Envios[0].Estado != notificacion.Enviado {
			t.Fatalf("%s: debía enviarse y quedar ENVIADO: %+v", e, repo.Envios)
		}
		if got := mailer.Enviados[0].Enlace; got != "http://app/seguimiento?codigo=TR-2026-ABC" {
			t.Errorf("enlace inesperado %q", got)
		}
	}
}

func TestIgnoraEstadosSinAviso(t *testing.T) {
	n, mailer, repo := armar("ana@x.bo")
	for _, e := range []tramite.Estado{tramite.EnRevision, tramite.Dictaminado, tramite.RevisionFinal} {
		n.Transicion(context.Background(), tr(e))
	}
	if len(mailer.Enviados) != 0 || len(repo.Envios) != 0 {
		t.Fatal("En revisión, Dictaminado y Revisión final no deben notificar")
	}
}

func TestSinCorreoQuedaFallido(t *testing.T) {
	n, mailer, repo := armar("")
	n.Transicion(context.Background(), tr(tramite.EnFacultad))
	if len(mailer.Enviados) != 0 || repo.Envios[0].Estado != notificacion.Fallido || repo.Envios[0].Error == "" {
		t.Fatalf("sin correo debe registrarse FALLIDO con motivo: %+v", repo.Envios[0])
	}
}

func TestFalloSeReintentaHastaAgotarIntentos(t *testing.T) {
	n, mailer, repo := armar("ana@x.bo")
	mailer.ErrCambioEstado = errors.New("smtp caído")
	n.Transicion(context.Background(), tr(tramite.Observado))
	if e := repo.Envios[0]; e.Estado != notificacion.Pendiente || e.Intentos != 1 {
		t.Fatalf("tras el primer fallo debe seguir PENDIENTE con 1 intento: %+v", e)
	}
	_ = n.ReintentarPendientes(context.Background())
	_ = n.ReintentarPendientes(context.Background())
	if e := repo.Envios[0]; e.Estado != notificacion.Fallido || e.Intentos != notificacion.MaxIntentos {
		t.Fatalf("agotados los intentos debe quedar FALLIDO: %+v", e)
	}

	n2, mailer2, repo2 := armar("ana@x.bo")
	mailer2.ErrCambioEstado = errors.New("smtp caído")
	n2.Transicion(context.Background(), tr(tramite.Archivado))
	mailer2.ErrCambioEstado = nil
	_ = n2.ReintentarPendientes(context.Background())
	if repo2.Envios[0].Estado != notificacion.Enviado || len(mailer2.Enviados) != 1 {
		t.Fatalf("al volver el SMTP el reintento debe enviar: %+v", repo2.Envios[0])
	}
}
