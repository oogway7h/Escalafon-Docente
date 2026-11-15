package legajos_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/legajos"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
)

var ahora = time.Date(2026, 10, 17, 10, 0, 0, 0, time.UTC)

type entorno struct {
	repo     *testdobles.LegajoRepoFalso
	tramites *testdobles.TramiteRepoFalso
	aud   *testdobles.AuditoriaRepoFalso
	reloj *testdobles.RelojFijo
}

func nuevoEntorno() *entorno {
	e := &entorno{repo: testdobles.NewLegajoRepoFalso(), tramites: testdobles.NewTramiteRepoFalso(), aud: &testdobles.AuditoriaRepoFalso{}, reloj: &testdobles.RelojFijo{Momento: ahora}}
	e.repo.Sembrar(legajo.Legajo{ID: 1, Codigo: "LEG-00001", Estado: legajo.EnArchivo})
	analista := "analista"
	e.tramites.Tramites["t1"] = &tramite.Tramite{ID: "t1", Codigo: "ED-2026-AAAAAA", Estado: tramite.EnRevision, AnalistaID: &analista}
	return e
}

func (e *entorno) solicitador() *legajos.Solicitar {
	return &legajos.Solicitar{Legajos: e.repo, Tramites: e.tramites, Auditoria: e.aud, Reloj: e.reloj}
}

func (e *entorno) solicitar(t *testing.T) *legajo.Solicitud {
	t.Helper()
	s, err := e.solicitador().Ejecutar(context.Background(), "analista", 1, "t1", "Revisar antigüedad")
	if err != nil {
		t.Fatalf("solicitar: %v", err)
	}
	return s
}

func (e *entorno) atender(id string, aprobar bool, nota string) error {
	return (&legajos.Atender{Legajos: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), "archivo", id, aprobar, nota)
}

func (e *entorno) despachar(id string, limite time.Time, material legajo.EstadoMaterial) error {
	return (&legajos.Despachar{Legajos: e.repo, Auditoria: e.aud, Tx: testdobles.TxManagerFalso{}, Reloj: e.reloj}).Ejecutar(context.Background(),
		legajos.DespacharEntrada{ActorID: "archivo", SolicitudID: id, FechaLimite: limite, EstadoMaterial: material})
}

func (e *entorno) devolver(material legajo.EstadoMaterial) error {
	return (&legajos.Devolver{Legajos: e.repo, Auditoria: e.aud, Tx: testdobles.TxManagerFalso{}, Reloj: e.reloj}).Ejecutar(context.Background(),
		legajos.DevolverEntrada{ActorID: "archivo", LegajoID: 1, EstadoMaterial: material, Observacion: "novedad"})
}

func TestFlujoCompleto(t *testing.T) {
	e := nuevoEntorno()
	s := e.solicitar(t)
	if err := e.atender(s.ID, true, ""); err != nil {
		t.Fatal(err)
	}
	if err := e.despachar(s.ID, ahora.AddDate(0, 0, 7), legajo.Bueno); err != nil {
		t.Fatal(err)
	}
	if e.repo.Legajos[1].Estado != legajo.Prestado {
		t.Fatalf("el legajo debería estar prestado, está %s", e.repo.Legajos[1].Estado)
	}
	if err := e.devolver(legajo.Regular); err != nil {
		t.Fatal(err)
	}
	if e.repo.Legajos[1].Estado != legajo.EnArchivo || e.repo.Solicitudes[s.ID].Estado != legajo.Devuelta {
		t.Fatal("la devolución debe dejar el legajo en archivo y la solicitud devuelta")
	}
	if len(e.repo.Movimientos) != 2 || e.repo.Movimientos[0].Tipo != legajo.Salida || e.repo.Movimientos[1].Tipo != legajo.Entrada {
		t.Fatalf("bitácora inesperada: %+v", e.repo.Movimientos)
	}
}

func (e *entorno) entregar(id string, limite time.Time) error {
	atender := &legajos.Atender{Legajos: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	despachar := &legajos.Despachar{Legajos: e.repo, Auditoria: e.aud, Tx: testdobles.TxManagerFalso{}, Reloj: e.reloj}
	return (&legajos.Entregar{Legajos: e.repo, Atender: atender, Despachar: despachar}).Ejecutar(context.Background(),
		legajos.DespacharEntrada{ActorID: "archivo", SolicitudID: id, FechaLimite: limite, EstadoMaterial: legajo.Bueno})
}

func TestEntregar_AprueboYDespachaEnUnPaso(t *testing.T) {
	e := nuevoEntorno()
	s := e.solicitar(t)
	if err := e.entregar(s.ID, ahora); !errors.Is(err, legajo.ErrFechaLimiteInvalid) {
		t.Fatalf("fecha de hoy: %v", err)
	}
	if e.repo.Solicitudes[s.ID].Estado != legajo.Pendiente {
		t.Fatal("un dato inválido no debe dejar la solicitud aprobada a medias")
	}
	if err := e.entregar(s.ID, ahora.AddDate(0, 0, 7)); err != nil {
		t.Fatal(err)
	}
	if e.repo.Solicitudes[s.ID].Estado != legajo.Despachada || e.repo.Legajos[1].Estado != legajo.Prestado || len(e.repo.Movimientos) != 1 {
		t.Fatalf("debe quedar despachada y prestada con su salida en la bitácora: %+v", e.repo.Solicitudes[s.ID])
	}
	if err := e.entregar(s.ID, ahora.AddDate(0, 0, 7)); !errors.Is(err, legajo.ErrEstadoSolicitud) {
		t.Fatalf("entregar dos veces: %v", err)
	}
}

func TestEntregar_SolicitudYaAprobadaSeEntregaDirecto(t *testing.T) {
	e := nuevoEntorno()
	s := e.solicitar(t)
	if err := e.atender(s.ID, true, ""); err != nil {
		t.Fatal(err)
	}
	if err := e.entregar(s.ID, ahora.AddDate(0, 0, 3)); err != nil {
		t.Fatal(err)
	}
	if e.repo.Solicitudes[s.ID].Estado != legajo.Despachada {
		t.Fatal("la solicitud aprobada debe despacharse")
	}
}

func TestDevolver_ExigeObservacionSiNoEstaEnBuenEstado(t *testing.T) {
	e := nuevoEntorno()
	s := e.solicitar(t)
	if err := e.entregar(s.ID, ahora.AddDate(0, 0, 3)); err != nil {
		t.Fatal(err)
	}
	devolver := &legajos.Devolver{Legajos: e.repo, Auditoria: e.aud, Tx: testdobles.TxManagerFalso{}, Reloj: e.reloj}
	if err := devolver.Ejecutar(context.Background(), legajos.DevolverEntrada{ActorID: "archivo", LegajoID: 1, EstadoMaterial: legajo.Deteriorado}); !errors.Is(err, legajo.ErrObservacionRequerida) {
		t.Fatalf("deteriorada sin observación: %v", err)
	}
	if err := devolver.Ejecutar(context.Background(), legajos.DevolverEntrada{ActorID: "archivo", LegajoID: 1, EstadoMaterial: legajo.Bueno}); err != nil {
		t.Fatalf("en buen estado no la exige: %v", err)
	}
}

func TestReglas(t *testing.T) {
	e := nuevoEntorno()
	s := e.solicitar(t)

	_, err := e.solicitador().Ejecutar(context.Background(), "analista", 1, "t1", "Otra vez")
	if !errors.Is(err, legajo.ErrSolicitudAbierta) {
		t.Errorf("segunda solicitud: %v", err)
	}
	if err := e.despachar(s.ID, ahora.AddDate(0, 0, 7), legajo.Bueno); !errors.Is(err, legajo.ErrEstadoSolicitud) {
		t.Errorf("despachar sin atender: %v", err)
	}
	if err := e.atender(s.ID, false, ""); !errors.Is(err, legajo.ErrMotivoRequerido) {
		t.Errorf("rechazar sin motivo: %v", err)
	}
	if err := e.atender(s.ID, true, ""); err != nil {
		t.Fatal(err)
	}
	if err := e.atender(s.ID, true, ""); !errors.Is(err, legajo.ErrEstadoSolicitud) {
		t.Errorf("atender dos veces: %v", err)
	}
	if err := e.despachar(s.ID, ahora, legajo.Bueno); !errors.Is(err, legajo.ErrFechaLimiteInvalid) {
		t.Errorf("fecha límite de hoy: %v", err)
	}
	if err := e.despachar(s.ID, ahora.AddDate(0, 0, 3), "ROTO"); !errors.Is(err, legajo.ErrEstadoMaterial) {
		t.Errorf("estado material inválido: %v", err)
	}
	if err := e.devolver(legajo.Bueno); !errors.Is(err, legajo.ErrNoPrestado) {
		t.Errorf("devolver lo que no salió: %v", err)
	}
}

func TestSolicitarPrestadoYSinMotivo(t *testing.T) {
	e := nuevoEntorno()
	caso := e.solicitador()
	if _, err := caso.Ejecutar(context.Background(), "analista", 1, "t1", "  "); !errors.Is(err, legajo.ErrMotivoRequerido) {
		t.Errorf("sin motivo: %v", err)
	}
	e.repo.Legajos[1].Estado = legajo.Prestado
	if _, err := caso.Ejecutar(context.Background(), "analista", 1, "t1", "x"); !errors.Is(err, legajo.ErrNoDisponible) {
		t.Errorf("legajo prestado: %v", err)
	}
	if _, err := caso.Ejecutar(context.Background(), "analista", 99, "t1", "x"); !errors.Is(err, legajo.ErrNoEncontrado) {
		t.Errorf("legajo inexistente: %v", err)
	}
}

func TestSolicitar_SinTramiteSeCreaSinLigar(t *testing.T) {
	e := nuevoEntorno()
	s, err := e.solicitador().Ejecutar(context.Background(), "analista", 1, "", "Consulta de antecedentes")
	if err != nil || s.TramiteID != "" || s.Estado != legajo.Pendiente {
		t.Fatalf("sin trámite la solicitud debe crearse sin ligar: %+v %v", s, err)
	}
}

func TestSolicitar_SiIndicaTramiteDebeSerAbiertoYDelSolicitante(t *testing.T) {
	e := nuevoEntorno()
	caso := e.solicitador()
	if _, err := caso.Ejecutar(context.Background(), "otro", 1, "t1", "x"); !errors.Is(err, legajo.ErrTramiteInvalido) {
		t.Errorf("trámite de otro analista: %v", err)
	}
	if _, err := caso.Ejecutar(context.Background(), "analista", 1, "no-existe", "x"); !errors.Is(err, legajo.ErrTramiteInvalido) {
		t.Errorf("trámite inexistente: %v", err)
	}
	e.tramites.Tramites["t1"].DocenteID = "otro-docente"
	if _, err := caso.Ejecutar(context.Background(), "analista", 1, "t1", "x"); !errors.Is(err, legajo.ErrTramiteInvalido) {
		t.Errorf("carpeta de otro docente: %v", err)
	}
	e.tramites.Tramites["t1"].DocenteID = ""
	e.tramites.Tramites["t1"].Estado = tramite.Remitido
	if _, err := caso.Ejecutar(context.Background(), "analista", 1, "t1", "x"); !errors.Is(err, legajo.ErrTramiteInvalido) {
		t.Errorf("trámite cerrado: %v", err)
	}
	e.tramites.Tramites["t1"].Estado = tramite.EnRevision
	s, err := caso.Ejecutar(context.Background(), "analista", 1, "t1", "Revisar antigüedad")
	if err != nil || s.TramiteID != "t1" {
		t.Fatalf("la solicitud debe quedar ligada al trámite: %+v %v", s, err)
	}
}
