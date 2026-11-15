package tramites_test

import (
	"context"
	"errors"
	"testing"

	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/tramite"
)

func (e *entorno) remitirFacultad(actor, id string) error {
	return (&tramites.RemitirAFacultad{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), actor, id)
}

func (e *entorno) regreso(actor string, ventanilla bool, id, numero string) error {
	return (&tramites.RegistrarRegreso{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), tramites.RegresoEntrada{
		ActorID: actor, EsVentanilla: ventanilla, TramiteID: id, ResolucionNumero: numero, ResolucionFecha: ahora.AddDate(0, 0, -2),
	})
}

func (e *entorno) remitirRectorado(actor, id, nota string) error {
	return (&tramites.RemitirARectorado{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), actor, id, nota)
}

// avanzarHastaRemitido lleva un trámite ya dictaminado por la facultad y el Rectorado.
func (e *entorno) avanzarHastaRemitido(t *testing.T, id string) {
	t.Helper()
	if err := e.remitirFacultad(analista, id); err != nil {
		t.Fatalf("remitir a la facultad: %v", err)
	}
	if err := e.regreso(analista, false, id, "RES-045/2026"); err != nil {
		t.Fatalf("registrar regreso: %v", err)
	}
	if err := e.remitirRectorado(analista, id, "Nota 112/2026 con ejemplares"); err != nil {
		t.Fatalf("remitir al Rectorado: %v", err)
	}
}

// dictaminado deja un trámite listo para salir a la facultad.
func (e *entorno) dictaminado(t *testing.T) *tramite.Tramite {
	t.Helper()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	_ = e.iniciarRevision(analista, tr.ID)
	if err := e.dictaminar(tr.ID); err != nil {
		t.Fatalf("dictaminar: %v", err)
	}
	return tr
}

func TestIniciarRevision_ExigeCarpetaSolicitada(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	caso := &tramites.IniciarRevision{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	if err := caso.Ejecutar(context.Background(), analista, tr.ID); !errors.Is(err, tramite.ErrCarpetaNoSolicitada) {
		t.Fatalf("sin carpeta solicitada esperaba ErrCarpetaNoSolicitada, obtuve %v", err)
	}
	e.repo.Tramites[tr.ID].CarpetasLigadas = 1
	if err := caso.Ejecutar(context.Background(), analista, tr.ID); err != nil {
		t.Fatalf("con la carpeta solicitada debe iniciar: %v", err)
	}
}

func TestDictaminar_ExigeLosAniosAEvaluar(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	_ = e.iniciarRevision(analista, tr.ID)
	dictamen := &tramites.Dictaminar{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	entrada := tramites.DictamenEntrada{ActorID: analista, TramiteID: tr.ID, Resultado: tramite.Favorable, Texto: "Cumple los requisitos"}
	if err := dictamen.Ejecutar(context.Background(), entrada); !errors.Is(err, tramite.ErrPeriodoRequerido) {
		t.Fatalf("esperaba ErrPeriodoRequerido, obtuve %v", err)
	}

	periodo := &tramites.RegistrarPeriodo{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	for _, p := range [][2]int16{{2021, 2016}, {1900, 2020}, {2016, 2099}} {
		if err := periodo.Ejecutar(context.Background(), analista, tr.ID, p[0], p[1]); !errors.Is(err, tramite.ErrPeriodoInvalido) {
			t.Fatalf("%v: esperaba ErrPeriodoInvalido, obtuve %v", p, err)
		}
	}
	if err := periodo.Ejecutar(context.Background(), otro, tr.ID, 2016, 2021); !errors.Is(err, tramite.ErrNoAsignado) {
		t.Fatalf("esperaba ErrNoAsignado, obtuve %v", err)
	}
	if err := periodo.Ejecutar(context.Background(), analista, tr.ID, 2016, 2021); err != nil {
		t.Fatalf("registrar periodo: %v", err)
	}
	if err := dictamen.Ejecutar(context.Background(), entrada); err != nil {
		t.Fatalf("dictaminar con periodo: %v", err)
	}
}

func TestRemitirAFacultad_SoloElAnalistaYDespuesDelDictamen(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	_ = e.iniciarRevision(analista, tr.ID)
	if err := e.remitirFacultad(analista, tr.ID); !errors.Is(err, tramite.ErrTransicionInvalida) {
		t.Fatalf("sin dictamen esperaba ErrTransicionInvalida, obtuve %v", err)
	}
	_ = e.dictaminar(tr.ID)
	if err := e.remitirFacultad(otro, tr.ID); !errors.Is(err, tramite.ErrNoAsignado) {
		t.Fatalf("esperaba ErrNoAsignado, obtuve %v", err)
	}
	if err := e.remitirFacultad(analista, tr.ID); err != nil {
		t.Fatalf("remitir: %v", err)
	}
	actual, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	if actual.Estado != tramite.EnFacultad || actual.Estado.Activo() {
		t.Fatalf("en facultad no cuenta como carga del analista: %s", actual.Estado)
	}
}

func TestRegreso_AnalistaOVentanillaYVuelveAlMismoAnalista(t *testing.T) {
	e := nuevoEntorno()
	tr := e.dictaminado(t)
	if err := e.regreso(analista, false, tr.ID, "RES-1"); !errors.Is(err, tramite.ErrTransicionInvalida) {
		t.Fatalf("no ha salido: esperaba ErrTransicionInvalida, obtuve %v", err)
	}
	_ = e.remitirFacultad(analista, tr.ID)

	if err := e.regreso(otro, false, tr.ID, "RES-1"); !errors.Is(err, tramite.ErrNoAsignado) {
		t.Fatalf("otro analista esperaba ErrNoAsignado, obtuve %v", err)
	}
	if err := e.regreso(analista, false, tr.ID, "  "); !errors.Is(err, tramite.ErrResolucionRequerida) {
		t.Fatalf("sin número esperaba ErrResolucionRequerida, obtuve %v", err)
	}
	futuro := &tramites.RegistrarRegreso{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	if err := futuro.Ejecutar(context.Background(), tramites.RegresoEntrada{ActorID: otro, EsVentanilla: true, TramiteID: tr.ID, ResolucionNumero: "RES-1", ResolucionFecha: ahora.AddDate(0, 0, 3)}); !errors.Is(err, tramite.ErrResolucionRequerida) {
		t.Fatalf("fecha futura esperaba ErrResolucionRequerida, obtuve %v", err)
	}
	if err := e.regreso(otro, true, tr.ID, "RES-045/2026"); err != nil { // ventanilla
		t.Fatalf("ventanilla registra el regreso: %v", err)
	}
	actual, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	if actual.Estado != tramite.RevisionFinal || actual.ResolucionNumero != "RES-045/2026" || actual.ResolucionFecha == nil {
		t.Fatalf("regreso mal registrado: %+v", actual)
	}
	if actual.AnalistaID == nil || *actual.AnalistaID != analista {
		t.Fatalf("debe volver al mismo analista, está con %v", actual.AnalistaID)
	}
}

func TestRemitirARectorado_ExigeNotaYCierraParaElAnalista(t *testing.T) {
	e := nuevoEntorno()
	tr := e.dictaminado(t)
	_ = e.remitirFacultad(analista, tr.ID)
	_ = e.regreso(analista, false, tr.ID, "RES-045/2026")

	if err := e.remitirRectorado(analista, tr.ID, " "); !errors.Is(err, tramite.ErrNotaRequerida) {
		t.Fatalf("esperaba ErrNotaRequerida, obtuve %v", err)
	}
	if err := e.remitirRectorado(analista, tr.ID, "Nota 112/2026"); err != nil {
		t.Fatalf("remitir: %v", err)
	}
	if err := e.observar(tr.ID); !errors.Is(err, tramite.ErrRemitido) {
		t.Fatalf("remitido es de solo lectura para el analista, obtuve %v", err)
	}
	if err := e.asignar(t, tr.ID, "Cambio"); !errors.Is(err, tramite.ErrRemitido) {
		t.Fatalf("tampoco se reasigna, obtuve %v", err)
	}
}

func TestArchivar_ExigeCarpetasDevueltas(t *testing.T) {
	e := nuevoEntorno()
	tr := e.dictaminado(t)
	e.avanzarHastaRemitido(t, tr.ID)
	archivar := &tramites.Archivar{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}

	e.repo.Tramites[tr.ID].CarpetasAbiertas = 1
	if err := archivar.Ejecutar(context.Background(), archivo, tr.ID, "Estante 3"); !errors.Is(err, tramite.ErrCarpetasAbiertas) {
		t.Fatalf("con la carpeta fuera esperaba ErrCarpetasAbiertas, obtuve %v", err)
	}
	e.repo.Tramites[tr.ID].CarpetasAbiertas = 0
	if err := archivar.Ejecutar(context.Background(), archivo, tr.ID, "Estante 3"); err != nil {
		t.Fatalf("archivar: %v", err)
	}
	actual, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	if actual.Estado != tramite.Archivado {
		t.Fatalf("esperaba ARCHIVADO, está en %s", actual.Estado)
	}
}

func TestReasignar_EnFacultadConMotivo(t *testing.T) {
	e := nuevoEntorno()
	tr := e.dictaminado(t)
	_ = e.remitirFacultad(analista, tr.ID)
	if err := e.asignar(t, tr.ID, "El analista dejó el cargo"); err != nil {
		t.Fatalf("reasignar mientras está en la facultad: %v", err)
	}
	actual, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	if actual.Estado != tramite.EnFacultad {
		t.Fatalf("la reasignación no cambia el estado, está en %s", actual.Estado)
	}
}

func TestPublico_ElResultadoSeMuestraSoloAlRemitir(t *testing.T) {
	e := nuevoEntorno()
	tr := e.dictaminado(t)
	consultar := func() *tramites.VistaPublica {
		v, err := (&tramites.ConsultarPublico{Tramites: e.repo}).Ejecutar(context.Background(), tr.Codigo)
		if err != nil {
			t.Fatalf("consultar: %v", err)
		}
		return v
	}
	if consultar().Resultado != "" {
		t.Fatalf("el dictamen es un informe interno: no debe verse")
	}
	_ = e.remitirFacultad(analista, tr.ID)
	if v := consultar(); v.Estado != tramite.EnFacultad || v.Resultado != "" {
		t.Fatalf("en facultad: %+v", v)
	}
	_ = e.regreso(analista, false, tr.ID, "RES-1")
	_ = e.remitirRectorado(analista, tr.ID, "Nota 1")
	if consultar().Resultado != tramite.Favorable {
		t.Fatalf("remitido: el resultado ya debe verse")
	}
}
