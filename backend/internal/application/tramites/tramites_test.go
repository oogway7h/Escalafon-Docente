package tramites_test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/categoria"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/tramite"
)

var ahora = time.Date(2026, 10, 5, 10, 0, 0, 0, time.UTC)

const (
	director = "u-director"
	analista = "u-analista"
	otro     = "u-otro"
	archivo  = "u-archivo"
)

type entorno struct {
	repo     *testdobles.TramiteRepoFalso
	docentes *testdobles.DocenteRepoFalso
	cats     *testdobles.CategoriaRepoFalso
	almacen  *testdobles.AlmacenFalso
	aud      *testdobles.AuditoriaRepoFalso
	reloj    *testdobles.RelojFijo
}

func nuevoEntorno() *entorno {
	e := &entorno{
		repo: testdobles.NewTramiteRepoFalso(), docentes: testdobles.NewDocenteRepoFalso(), cats: testdobles.NewCategoriaRepoFalso(),
		almacen: testdobles.NewAlmacenFalso(), aud: &testdobles.AuditoriaRepoFalso{}, reloj: &testdobles.RelojFijo{Momento: ahora},
	}
	e.repo.Analistas = []tramite.Analista{{ID: analista, Nombre: "Ana Analista"}}
	e.docentes.Sembrar(docente.Docente{ID: "d1", Codigo: "123", Cedula: "7654321", Nombres: "María", Apellidos: "Vargas", Activo: true})
	e.cats.Sembrar(categoria.Categoria{ID: 1, Nombre: "Titular A", Activa: true,
		Requisitos: []categoria.Requisito{{ID: 10, Nombre: "Título"}, {ID: 11, Nombre: "Cédula"}}})
	return e
}

func (e *entorno) radicar(t *testing.T) *tramite.Tramite {
	t.Helper()
	caso := &tramites.Radicar{Tramites: e.repo, Docentes: e.docentes, Categorias: e.cats, Auditoria: e.aud, Reloj: e.reloj}
	tr, err := caso.Ejecutar(context.Background(), tramites.RadicarEntrada{ActorID: director, DocenteID: "d1", Tipo: tramite.Categorizacion, CategoriaID: 1})
	if err != nil {
		t.Fatalf("radicar: %v", err)
	}
	return tr
}

func (e *entorno) asignar(t *testing.T, id, motivo string) error {
	t.Helper()
	caso := &tramites.Asignar{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	_, err := caso.Ejecutar(context.Background(), tramites.AsignarEntrada{ActorID: director, TramiteID: id, AnalistaID: analista, Motivo: motivo})
	return err
}

// iniciarRevision asume la carpeta ya solicitada; la regla se prueba en TestIniciarRevision_ExigeCarpetaSolicitada.
func (e *entorno) iniciarRevision(actor, id string) error {
	e.repo.Tramites[id].CarpetasLigadas = 1
	return (&tramites.IniciarRevision{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), actor, id)
}

func (e *entorno) observar(id string) error {
	return (&tramites.RegistrarObservacion{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(),
		tramites.ObservacionEntrada{ActorID: analista, TramiteID: id, Documento: "Título", Motivo: "Ilegible", Plazo: ahora.AddDate(0, 0, 10)})
}

// dictaminar asume los años a evaluar ya registrados.
func (e *entorno) dictaminar(id string) error {
	e.repo.Tramites[id].AnioDesde, e.repo.Tramites[id].AnioHasta = 2016, 2021
	return (&tramites.Dictaminar{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(),
		tramites.DictamenEntrada{ActorID: analista, TramiteID: id, Resultado: tramite.Favorable, Texto: "Cumple los requisitos"})
}

func TestFlujoCompleto_HastaQueArchivoArchiva(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	if tr.Estado != tramite.Ingresado || !strings.HasPrefix(tr.Codigo, "ED-2026-") {
		t.Fatalf("estado/código inesperados: %s %s", tr.Estado, tr.Codigo)
	}
	if err := e.asignar(t, tr.ID, ""); err != nil {
		t.Fatalf("asignar: %v", err)
	}
	if err := e.iniciarRevision(analista, tr.ID); err != nil {
		t.Fatalf("iniciar revisión: %v", err)
	}
	if err := e.dictaminar(tr.ID); err != nil {
		t.Fatalf("dictaminar: %v", err)
	}
	e.avanzarHastaRemitido(t, tr.ID)
	if err := (&tramites.Archivar{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), archivo, tr.ID, "Estante 3"); err != nil {
		t.Fatalf("archivar: %v", err)
	}
	final, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	if final.Estado != tramite.Archivado || final.UbicacionLegajo != "Estante 3" || final.ArchivadoEn == nil {
		t.Fatalf("trámite mal archivado: %+v", final)
	}
	// Archivado = solo lectura.
	if err := e.observar(tr.ID); !errors.Is(err, tramite.ErrArchivado) {
		t.Fatalf("esperaba ErrArchivado, obtuve %v", err)
	}
}

func TestRadicar_DocenteDeBajaYCategoriaInactiva(t *testing.T) {
	e := nuevoEntorno()
	caso := &tramites.Radicar{Tramites: e.repo, Docentes: e.docentes, Categorias: e.cats, Auditoria: e.aud, Reloj: e.reloj}
	e.docentes.Docentes["d1"].Activo = false
	if _, err := caso.Ejecutar(context.Background(), tramites.RadicarEntrada{DocenteID: "d1", Tipo: tramite.Categorizacion, CategoriaID: 1}); !errors.Is(err, tramite.ErrDocenteInactivo) {
		t.Fatalf("esperaba ErrDocenteInactivo, obtuve %v", err)
	}
	e.docentes.Docentes["d1"].Activo = true
	e.cats.Categorias[1].Activa = false
	if _, err := caso.Ejecutar(context.Background(), tramites.RadicarEntrada{DocenteID: "d1", Tipo: tramite.Categorizacion, CategoriaID: 1}); !errors.Is(err, tramite.ErrCategoriaInvalida) {
		t.Fatalf("esperaba ErrCategoriaInvalida, obtuve %v", err)
	}
}

func TestAsignar_ReasignarExigeMotivo(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	if err := e.asignar(t, tr.ID, ""); err != nil {
		t.Fatal(err)
	}
	if err := e.asignar(t, tr.ID, "  "); !errors.Is(err, tramite.ErrMotivoRequerido) {
		t.Fatalf("esperaba ErrMotivoRequerido, obtuve %v", err)
	}
	if err := e.asignar(t, tr.ID, "Carga desbalanceada"); err != nil {
		t.Fatalf("reasignar con motivo: %v", err)
	}
	final, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	ultimo := final.Historial[len(final.Historial)-1]
	if !strings.Contains(ultimo.Detalle, "Carga desbalanceada") || final.Estado != tramite.Asignado {
		t.Fatalf("el motivo debe quedar en el historial sin cambiar de estado: %+v", ultimo)
	}
}

func TestAsignar_AnalistaInvalido(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	caso := &tramites.Asignar{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	if _, err := caso.Ejecutar(context.Background(), tramites.AsignarEntrada{ActorID: director, TramiteID: tr.ID, AnalistaID: "fantasma"}); !errors.Is(err, tramite.ErrAnalistaInvalido) {
		t.Fatalf("esperaba ErrAnalistaInvalido, obtuve %v", err)
	}
}

func TestRevision_SoloElAnalistaAsignado(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	if err := e.iniciarRevision(otro, tr.ID); !errors.Is(err, tramite.ErrNoAsignado) {
		t.Fatalf("esperaba ErrNoAsignado, obtuve %v", err)
	}
}

func TestDictamen_BloqueadoConObservacionesPendientes(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	_ = e.iniciarRevision(analista, tr.ID)
	if err := e.observar(tr.ID); err != nil {
		t.Fatalf("observar: %v", err)
	}
	if err := e.dictaminar(tr.ID); !errors.Is(err, tramite.ErrObservacionesPendientes) {
		t.Fatalf("esperaba ErrObservacionesPendientes, obtuve %v", err)
	}

	obs := e.repo.Observaciones[tr.ID][0]
	resolver := &tramites.ResolverObservacion{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	if err := resolver.Ejecutar(context.Background(), analista, tr.ID, obs.ID); err != nil {
		t.Fatalf("resolver: %v", err)
	}
	actual, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	if actual.Estado != tramite.EnRevision {
		t.Fatalf("al resolver la última observación debe volver a EN_REVISION, está en %s", actual.Estado)
	}
	if err := e.dictaminar(tr.ID); err != nil {
		t.Fatalf("dictaminar tras resolver: %v", err)
	}
}

func TestObservacion_PlazoEnElPasado(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	_ = e.iniciarRevision(analista, tr.ID)
	err := (&tramites.RegistrarObservacion{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(),
		tramites.ObservacionEntrada{ActorID: analista, TramiteID: tr.ID, Documento: "Título", Motivo: "Ilegible", Plazo: ahora.AddDate(0, 0, -3)})
	if !errors.Is(err, tramite.ErrPlazoInvalido) {
		t.Fatalf("esperaba ErrPlazoInvalido, obtuve %v", err)
	}
}

func TestArchivar_SoloDespuesDeRemitirAlRectorado(t *testing.T) {
	e := nuevoEntorno()
	tr := e.radicar(t)
	_ = e.asignar(t, tr.ID, "")
	_ = e.iniciarRevision(analista, tr.ID)
	_ = e.dictaminar(tr.ID)
	err := (&tramites.Archivar{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), archivo, tr.ID, "Estante")
	if !errors.Is(err, tramite.ErrTransicionInvalida) {
		t.Fatalf("esperaba ErrTransicionInvalida, obtuve %v", err)
	}
}

// pdf mínimo: la validación detecta el tipo por el contenido, no por la extensión.
var pdf = append([]byte("%PDF-1.4\n"), bytes.Repeat([]byte("x"), 50)...)

func (e *entorno) publico() *tramites.RadicarPublico {
	return &tramites.RadicarPublico{Tramites: e.repo, Docentes: e.docentes, Categorias: e.cats, Almacen: e.almacen, Auditoria: e.aud, Reloj: e.reloj}
}

func TestRadicarPublico_ExigeDocumentosDeLaCategoria(t *testing.T) {
	e := nuevoEntorno()
	_, err := e.publico().Ejecutar(context.Background(), tramites.RadicarPublicoEntrada{
		Codigo: "123", Cedula: "7654321", CategoriaID: 1, Archivos: []tramites.Archivo{{RequisitoID: 10, Nombre: "titulo.pdf", Datos: pdf}},
	})
	var faltan *tramite.ErrDocumentosFaltantes
	if !errors.As(err, &faltan) || len(faltan.Faltan) != 1 || faltan.Faltan[0] != "Cédula" {
		t.Fatalf("esperaba faltar «Cédula», obtuve %v", err)
	}
}

func TestRadicarPublico_CedulaDesconocidaYArchivoInvalido(t *testing.T) {
	e := nuevoEntorno()
	archivos := []tramites.Archivo{{RequisitoID: 10, Nombre: "a.pdf", Datos: pdf}, {RequisitoID: 11, Nombre: "b.pdf", Datos: pdf}}
	if _, err := e.publico().Ejecutar(context.Background(), tramites.RadicarPublicoEntrada{Codigo: "999", Cedula: "7654321", CategoriaID: 1, Archivos: archivos}); !errors.Is(err, tramite.ErrDocenteNoVerificado) {
		t.Fatalf("código desconocido: esperaba ErrDocenteNoVerificado, obtuve %v", err)
	}
	// Con el código correcto pero otra cédula la respuesta es la misma: no se puede radicar a nombre de otro.
	if _, err := e.publico().Ejecutar(context.Background(), tramites.RadicarPublicoEntrada{Codigo: "123", Cedula: "1111111", CategoriaID: 1, Archivos: archivos}); !errors.Is(err, tramite.ErrDocenteNoVerificado) {
		t.Fatalf("cédula ajena: esperaba ErrDocenteNoVerificado, obtuve %v", err)
	}
	malos := []tramites.Archivo{{RequisitoID: 10, Nombre: "a.exe", Datos: []byte("MZ\x90\x00binario")}, {RequisitoID: 11, Nombre: "b.pdf", Datos: pdf}}
	if _, err := e.publico().Ejecutar(context.Background(), tramites.RadicarPublicoEntrada{Codigo: "123", Cedula: "7654321", CategoriaID: 1, Archivos: malos}); !errors.Is(err, tramite.ErrArchivoInvalido) {
		t.Fatalf("esperaba ErrArchivoInvalido, obtuve %v", err)
	}
}

func TestPortal_RadicarYConsultarSinExponerDatosPersonales(t *testing.T) {
	e := nuevoEntorno()
	codigo, err := e.publico().Ejecutar(context.Background(), tramites.RadicarPublicoEntrada{
		Codigo: "123", Cedula: "7654321", CategoriaID: 1, Archivos: []tramites.Archivo{{RequisitoID: 10, Nombre: "a.pdf", Datos: pdf}, {RequisitoID: 11, Nombre: "b.pdf", Datos: pdf}},
	})
	if err != nil {
		t.Fatalf("radicar público: %v", err)
	}
	if len(e.almacen.Guardados) != 2 {
		t.Fatalf("debían subirse 2 adjuntos, se subieron %d", len(e.almacen.Guardados))
	}

	// Un trámite del portal no se puede asignar hasta validarlo en ventanilla.
	var id string
	for k := range e.repo.Tramites {
		id = k
	}
	if err := e.asignar(t, id, ""); !errors.Is(err, tramite.ErrPendienteValidacion) {
		t.Fatalf("esperaba ErrPendienteValidacion, obtuve %v", err)
	}
	if err := (&tramites.ValidarSolicitud{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(), director, id); err != nil {
		t.Fatalf("validar: %v", err)
	}
	if err := e.asignar(t, id, ""); err != nil {
		t.Fatalf("asignar tras validar: %v", err)
	}

	vista, err := (&tramites.ConsultarPublico{Tramites: e.repo}).Ejecutar(context.Background(), " "+strings.ToLower(codigo)+" ")
	if err != nil {
		t.Fatalf("consultar público (el código debe normalizarse): %v", err)
	}
	if vista.Estado != tramite.Asignado || len(vista.Etapas) != 2 {
		t.Fatalf("vista inesperada: %+v", vista)
	}
	// La validación en ventanilla no es una fase: no aparece como etapa.
	for _, et := range vista.Etapas {
		if et.Estado != tramite.Ingresado && et.Estado != tramite.Asignado {
			t.Fatalf("etapa inesperada %s", et.Estado)
		}
	}

	if _, err := (&tramites.ConsultarPublico{Tramites: e.repo}).Ejecutar(context.Background(), "ED-2026-NOEXISTE"); !errors.Is(err, tramite.ErrNoEncontrado) {
		t.Fatalf("esperaba ErrNoEncontrado, obtuve %v", err)
	}
}

// --- CU-16: dar de baja un documento ---

// conDocumentos radica desde el portal para tener un expediente con adjuntos.
func (e *entorno) conDocumentos(t *testing.T) *tramite.Tramite {
	t.Helper()
	codigo, err := e.publico().Ejecutar(context.Background(), tramites.RadicarPublicoEntrada{
		Codigo: "123", Cedula: "7654321", CategoriaID: 1,
		Archivos: []tramites.Archivo{{RequisitoID: 10, Nombre: "titulo.pdf", Datos: pdf}, {RequisitoID: 11, Nombre: "cedula.pdf", Datos: pdf}},
	})
	if err != nil {
		t.Fatalf("radicar público: %v", err)
	}
	tr, err := e.repo.BuscarPorCodigo(context.Background(), codigo)
	if err != nil {
		t.Fatalf("buscar por código: %v", err)
	}
	a := analista
	e.repo.Tramites[tr.ID].AnalistaID = &a // la baja la hace el analista asignado
	return tr
}

func (e *entorno) darDeBaja(id, docID, motivo string) error {
	return (&tramites.DarDeBajaDocumento{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(),
		tramites.BajaDocumentoEntrada{ActorID: analista, TramiteID: id, DocumentoID: docID, Motivo: motivo})
}

func TestDarDeBajaDocumento_ExigeMotivoYSoloSePuedeUnaVez(t *testing.T) {
	e := nuevoEntorno()
	tr := e.conDocumentos(t)
	doc := tr.Documentos[0]

	if err := e.darDeBaja(tr.ID, doc.ID, "   "); !errors.Is(err, tramite.ErrMotivoBajaRequerido) {
		t.Fatalf("esperaba ErrMotivoBajaRequerido, obtuve %v", err)
	}
	if err := e.darDeBaja(tr.ID, doc.ID, "Cargado por error"); err != nil {
		t.Fatalf("dar de baja: %v", err)
	}

	actual, _ := e.repo.BuscarPorID(context.Background(), tr.ID)
	baja := actual.Documentos[0]
	if baja.Activo || baja.MotivoBaja != "Cargado por error" || baja.BajaEn == nil || baja.ClaveObjeto == "" {
		t.Fatalf("la baja es lógica: conserva el objeto con su motivo y fecha, obtuve %+v", baja)
	}
	if !actual.Documentos[1].Activo {
		t.Fatalf("la baja no debe afectar a los demás adjuntos del expediente")
	}
	ultimo := e.aud.Eventos[len(e.aud.Eventos)-1]
	if ultimo.Accion != auditoria.AccionBajaDocumento || ultimo.Detalle["documentoId"] != doc.ID {
		t.Fatalf("la baja debe quedar en la bitácora: %+v", ultimo)
	}
	if err := e.darDeBaja(tr.ID, doc.ID, "Cargado por error otra vez"); !errors.Is(err, tramite.ErrDocumentoDeBaja) {
		t.Fatalf("esperaba ErrDocumentoDeBaja, obtuve %v", err)
	}
}

func TestDarDeBajaDocumento_SoloElAnalistaAsignado(t *testing.T) {
	e := nuevoEntorno()
	tr := e.conDocumentos(t)
	err := (&tramites.DarDeBajaDocumento{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}).Ejecutar(context.Background(),
		tramites.BajaDocumentoEntrada{ActorID: otro, TramiteID: tr.ID, DocumentoID: tr.Documentos[0].ID, Motivo: "Cargado por error"})
	if !errors.Is(err, tramite.ErrNoAsignado) {
		t.Fatalf("esperaba ErrNoAsignado, obtuve %v", err)
	}
}

func TestUrlDocumento_LaBajaSoloLaVeQuienPuedeRestituir(t *testing.T) {
	e := nuevoEntorno()
	tr := e.conDocumentos(t)
	doc := tr.Documentos[0]
	if err := e.darDeBaja(tr.ID, doc.ID, "Cargado por error"); err != nil {
		t.Fatalf("dar de baja: %v", err)
	}
	url := &tramites.UrlDocumento{Tramites: e.repo, Almacen: e.almacen}
	if _, err := url.Ejecutar(context.Background(), tr.ID, doc.ID, false); !errors.Is(err, tramite.ErrDocumentoDeBaja) {
		t.Fatalf("esperaba ErrDocumentoDeBaja, obtuve %v", err)
	}
	if _, err := url.Ejecutar(context.Background(), tr.ID, doc.ID, true); err != nil {
		t.Fatalf("el Director debe poder verlo: %v", err)
	}
}

func TestDarDeBajaDocumento_ExpedienteArchivadoEsSoloLectura(t *testing.T) {
	e := nuevoEntorno()
	tr := e.conDocumentos(t)
	e.repo.Tramites[tr.ID].Estado = tramite.Archivado
	if err := e.darDeBaja(tr.ID, tr.Documentos[0].ID, "Cargado por error"); !errors.Is(err, tramite.ErrArchivado) {
		t.Fatalf("esperaba ErrArchivado, obtuve %v", err)
	}
}

func TestNuevoCodigo_FormatoYAlfabeto(t *testing.T) {
	for i := 0; i < 200; i++ {
		c, err := tramite.NuevoCodigo(2026)
		if err != nil || len(c) != len("ED-2026-ABCDEF") || strings.ContainsAny(c[8:], "0O1I") {
			t.Fatalf("código inválido %q (%v)", c, err)
		}
	}
}

// --- Regularización: sin categoría ---

func TestRadicar_RegularizacionNoLlevaCategoria(t *testing.T) {
	e := nuevoEntorno()
	caso := &tramites.Radicar{Tramites: e.repo, Docentes: e.docentes, Categorias: e.cats, Auditoria: e.aud, Reloj: e.reloj}
	// Aunque llegue una categoría (incluso inexistente) se ignora.
	tr, err := caso.Ejecutar(context.Background(), tramites.RadicarEntrada{ActorID: director, DocenteID: "d1", Tipo: tramite.Regularizacion, CategoriaID: 99})
	if err != nil || tr.CategoriaID != 0 {
		t.Fatalf("la regularización debe radicarse sin categoría: %+v %v", tr, err)
	}
	if _, err := caso.Ejecutar(context.Background(), tramites.RadicarEntrada{ActorID: director, DocenteID: "d1", Tipo: tramite.Categorizacion}); !errors.Is(err, tramite.ErrCategoriaInvalida) {
		t.Fatalf("la categorización sin categoría debe fallar, obtuve %v", err)
	}
}

func TestRadicarPublico_RegularizacionConAdjuntosLibres(t *testing.T) {
	e := nuevoEntorno()
	entrada := tramites.RadicarPublicoEntrada{Codigo: "123", Cedula: "7654321", Tipo: tramite.Regularizacion}
	if _, err := e.publico().Ejecutar(context.Background(), entrada); !errors.Is(err, tramite.ErrAdjuntosLibres) {
		t.Fatalf("sin adjuntos: esperaba ErrAdjuntosLibres, obtuve %v", err)
	}
	entrada.Archivos = []tramites.Archivo{{Nombre: "nota.pdf", Datos: pdf}, {Nombre: "respaldo.pdf", Datos: pdf}}
	codigo, err := e.publico().Ejecutar(context.Background(), entrada)
	if err != nil {
		t.Fatalf("radicar regularización: %v", err)
	}
	tr, _ := e.repo.BuscarPorCodigo(context.Background(), codigo)
	if tr.CategoriaID != 0 || len(tr.Documentos) != 2 || tr.Documentos[0].RequisitoID != nil {
		t.Fatalf("esperaba 2 adjuntos libres y ninguna categoría: %+v", tr)
	}
}

func TestRadicarPublico_TopeDeSolicitudesPendientes(t *testing.T) {
	e := nuevoEntorno()
	entrada := tramites.RadicarPublicoEntrada{Codigo: "123", Cedula: "7654321", Tipo: tramite.Regularizacion, Archivos: []tramites.Archivo{{Nombre: "nota.pdf", Datos: pdf}}}
	for i := 0; i < 2; i++ {
		if _, err := e.publico().Ejecutar(context.Background(), entrada); err != nil {
			t.Fatalf("solicitud %d: %v", i+1, err)
		}
	}
	if _, err := e.publico().Ejecutar(context.Background(), entrada); !errors.Is(err, tramite.ErrDemasiadasPendientes) {
		t.Fatalf("la tercera sin validar debe rechazarse, obtuve %v", err)
	}
}
