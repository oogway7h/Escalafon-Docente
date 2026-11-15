package tramites_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/tramite"
)

// asignado radica un trámite y lo asigna al analista.
func (e *entorno) asignado(t *testing.T) *tramite.Tramite {
	t.Helper()
	tr := e.radicar(t)
	if err := e.asignar(t, tr.ID, ""); err != nil {
		t.Fatalf("asignar: %v", err)
	}
	return tr
}

func (e *entorno) preparar(actor, id, mime string, tamano int64) (*tramites.SubidaPreparada, error) {
	return (&tramites.PrepararSubida{Tramites: e.repo, Almacen: e.almacen, Reloj: e.reloj}).Ejecutar(context.Background(),
		tramites.SubidaEntrada{ActorID: actor, TramiteID: id, Nombre: "resolución 12.pdf", TipoMime: mime, Tamano: tamano})
}

func (e *entorno) confirmador(p *testdobles.ProcesadorFalso) *tramites.ConfirmarSubida {
	return &tramites.ConfirmarSubida{Tramites: e.repo, Almacen: e.almacen, OCR: &tramites.ProcesoOCR{Tramites: e.repo, Procesador: p}, Auditoria: e.aud, Reloj: e.reloj}
}

func TestPrepararSubida_ValidaTipoTamanoYAnalista(t *testing.T) {
	e := nuevoEntorno()
	tr := e.asignado(t)

	if _, err := e.preparar(otro, tr.ID, "application/pdf", 10); !errors.Is(err, tramite.ErrNoAsignado) {
		t.Fatalf("esperaba ErrNoAsignado, obtuve %v", err)
	}
	if _, err := e.preparar(analista, tr.ID, "text/html", 10); !errors.Is(err, tramite.ErrTipoNoPermitido) {
		t.Fatalf("esperaba ErrTipoNoPermitido, obtuve %v", err)
	}
	if _, err := e.preparar(analista, tr.ID, "image/png", 21<<20); !errors.Is(err, tramite.ErrArchivoGrande) {
		t.Fatalf("una imagen de 21 MB debe rechazarse, obtuve %v", err)
	}
	if _, err := e.preparar(analista, tr.ID, "application/pdf", 400<<20+1); !errors.Is(err, tramite.ErrArchivoGrande) {
		t.Fatalf("un PDF de más de 400 MB debe rechazarse, obtuve %v", err)
	}
	r, err := e.preparar(analista, tr.ID, "application/pdf", 400<<20)
	if err != nil {
		t.Fatalf("un PDF de 400 MB es válido: %v", err)
	}
	if r.Documento.EstadoOCR != tramite.OCRSubiendo || !strings.HasPrefix(r.Documento.ClaveObjeto, "docentes/123/"+tr.ID+"/") || r.URL == "" {
		t.Fatalf("el documento va al directorio del docente y queda SUBIENDO: %+v", r)
	}
}

func TestConfirmarSubida_ProcesaYGuardaTextoYCompresion(t *testing.T) {
	e := nuevoEntorno()
	tr := e.asignado(t)
	r, _ := e.preparar(analista, tr.ID, "application/pdf", 100)
	p := &testdobles.ProcesadorFalso{Resultado: ports.ResultadoProceso{
		Estado: "LISTO", Texto: "Resolución 12/2026", HashSHA256: strings.Repeat("a", 64), TamanoOriginal: 100, TamanoFinal: 30,
		ClaveComprimida: "docentes/123/x-comprimido.pdf",
	}}
	c := e.confirmador(p)

	if _, err := c.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID); !errors.Is(err, tramite.ErrSubidaIncompleta) {
		t.Fatalf("sin objeto en el almacén esperaba ErrSubidaIncompleta, obtuve %v", err)
	}
	e.almacen.Guardados[r.Documento.ClaveObjeto] = make([]byte, 100)
	if _, err := c.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID); err != nil {
		t.Fatalf("confirmar: %v", err)
	}

	d, _ := e.repo.BuscarDocumento(context.Background(), tr.ID, r.Documento.ID)
	if d.EstadoOCR != tramite.OCRProcesado || d.TextoOCR != "Resolución 12/2026" || d.TamanoComprimido != 30 || d.ClaveDescarga() != "docentes/123/x-comprimido.pdf" {
		t.Fatalf("resultado del OCR no guardado: %+v", d)
	}
	if len(p.Recibidos) != 1 || !strings.HasSuffix(p.Recibidos[0].ClaveSalida, "-comprimido.pdf") {
		t.Fatalf("el PDF debe pedir su versión comprimida: %+v", p.Recibidos)
	}
	if e.aud.Eventos[len(e.aud.Eventos)-1].Accion != auditoria.AccionCargarDocumento {
		t.Fatalf("la carga debe quedar en la bitácora")
	}
	if _, err := c.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID); !errors.Is(err, tramite.ErrDocumentoNoSubido) {
		t.Fatalf("un documento ya procesado no se confirma otra vez, obtuve %v", err)
	}
}

func TestConfirmarSubida_ErrorDelOCRSePuedeReintentar(t *testing.T) {
	e := nuevoEntorno()
	tr := e.asignado(t)
	r, _ := e.preparar(analista, tr.ID, "image/png", 10)
	e.almacen.Guardados[r.Documento.ClaveObjeto] = make([]byte, 10)
	p := &testdobles.ProcesadorFalso{Resultado: ports.ResultadoProceso{Estado: "ERROR", Error: "imagen ilegible"}}
	c := e.confirmador(p)

	if _, err := c.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID); err != nil {
		t.Fatalf("confirmar: %v", err)
	}
	d, _ := e.repo.BuscarDocumento(context.Background(), tr.ID, r.Documento.ID)
	if d.EstadoOCR != tramite.OCRError || d.ErrorOCR != "imagen ilegible" {
		t.Fatalf("esperaba ERROR con su motivo: %+v", d)
	}
	if p.Recibidos[0].ClaveSalida != "" {
		t.Fatalf("las imágenes no se comprimen")
	}

	p.Resultado = ports.ResultadoProceso{Estado: "LISTO", Texto: "ok"}
	if _, err := c.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID); err != nil {
		t.Fatalf("reintentar: %v", err)
	}
	d, _ = e.repo.BuscarDocumento(context.Background(), tr.ID, r.Documento.ID)
	if d.EstadoOCR != tramite.OCRProcesado || d.ErrorOCR != "" {
		t.Fatalf("el reintento debe limpiar el error: %+v", d)
	}
}

// procesado deja un documento con OCR listo y sus páginas indexadas.
func (e *entorno) procesado(t *testing.T, paginas ...string) (*tramite.Tramite, string) {
	t.Helper()
	tr := e.asignado(t)
	r, _ := e.preparar(analista, tr.ID, "application/pdf", 10)
	e.almacen.Guardados[r.Documento.ClaveObjeto] = make([]byte, 10)
	p := &testdobles.ProcesadorFalso{Resultado: ports.ResultadoProceso{Estado: "LISTO", Paginas: paginas}}
	if _, err := e.confirmador(p).Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID); err != nil {
		t.Fatalf("confirmar: %v", err)
	}
	return tr, r.Documento.ID
}

func TestProcesar_GuardaElTextoPorPagina(t *testing.T) {
	e := nuevoEntorno()
	tr, doc := e.procesado(t, "Resolución rectoral número 12", "Título de licenciatura", "")
	paginas, err := (&tramites.TextoDocumento{Tramites: e.repo}).Ejecutar(context.Background(), tr.ID, doc, false)
	if err != nil || len(paginas) != 3 || paginas[1] != "Título de licenciatura" {
		t.Fatalf("debe guardar una entrada por página, incluida la vacía: %q %v", paginas, err)
	}
	d, _ := e.repo.BuscarDocumento(context.Background(), tr.ID, doc)
	if d.EstadoOCR != tramite.OCRProcesado || !strings.Contains(d.TextoOCR, "licenciatura") {
		t.Fatalf("texto_ocr debe ser la concatenación: %+v", d)
	}
}

func TestBuscarEnDocumento_PaginasYValidaciones(t *testing.T) {
	e := nuevoEntorno()
	tr, doc := e.procesado(t, "Resolución rectoral número 12", "Título de licenciatura", "Otra resolución rectoral")
	buscar := &tramites.BuscarEnDocumento{Tramites: e.repo}

	lista, err := buscar.Ejecutar(context.Background(), tr.ID, doc, " rectoral ", false)
	if err != nil || len(lista) != 2 || lista[0].Pagina != 1 || lista[1].Pagina != 3 {
		t.Fatalf("esperaba las páginas 1 y 3 en orden: %+v %v", lista, err)
	}
	if lista, _ := buscar.Ejecutar(context.Background(), tr.ID, doc, "inexistente", false); len(lista) != 0 {
		t.Fatalf("sin coincidencias no debe devolver páginas: %+v", lista)
	}
	for _, q := range []string{"", " a ", strings.Repeat("x", 101)} {
		if _, err := buscar.Ejecutar(context.Background(), tr.ID, doc, q, false); !errors.Is(err, tramite.ErrConsultaInvalida) {
			t.Fatalf("%q: esperaba ErrConsultaInvalida, obtuve %v", q, err)
		}
	}
	if _, err := buscar.Ejecutar(context.Background(), tr.ID, "no-existe", "rectoral", false); !errors.Is(err, tramite.ErrDocumentoNoEncontrado) {
		t.Fatalf("esperaba ErrDocumentoNoEncontrado, obtuve %v", err)
	}
}

func TestBuscarEnDocumento_SinOCRYDocumentoDeBaja(t *testing.T) {
	e := nuevoEntorno()
	tr := e.asignado(t)
	r, _ := e.preparar(analista, tr.ID, "application/pdf", 10) // sigue SUBIENDO
	buscar := &tramites.BuscarEnDocumento{Tramites: e.repo}
	if _, err := buscar.Ejecutar(context.Background(), tr.ID, r.Documento.ID, "rectoral", false); !errors.Is(err, tramite.ErrSinOCR) {
		t.Fatalf("sin OCR esperaba ErrSinOCR, obtuve %v", err)
	}

	tr2, doc := e.procesado(t, "Resolución rectoral")
	if err := e.darDeBaja(tr2.ID, doc, "Cargado por error"); err != nil {
		t.Fatalf("dar de baja: %v", err)
	}
	if _, err := buscar.Ejecutar(context.Background(), tr2.ID, doc, "rectoral", false); !errors.Is(err, tramite.ErrDocumentoDeBaja) {
		t.Fatalf("un documento de baja no se busca, obtuve %v", err)
	}
	if lista, err := buscar.Ejecutar(context.Background(), tr2.ID, doc, "rectoral", true); err != nil || len(lista) != 1 {
		t.Fatalf("quien puede restituir sí: %+v %v", lista, err)
	}
}

func TestCorregirPaginas_ConservaLasPaginasYQuedaValidado(t *testing.T) {
	e := nuevoEntorno()
	tr := e.asignado(t)
	r, _ := e.preparar(analista, tr.ID, "application/pdf", 10)
	e.almacen.Guardados[r.Documento.ClaveObjeto] = make([]byte, 10)
	corregir := &tramites.CorregirPaginas{Tramites: e.repo, Auditoria: e.aud, Tx: testdobles.TxManagerFalso{}, Reloj: e.reloj}

	if err := corregir.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID, []string{"x"}); !errors.Is(err, tramite.ErrDocumentoNoSubido) {
		t.Fatalf("sin OCR terminado esperaba ErrDocumentoNoSubido, obtuve %v", err)
	}
	p := &testdobles.ProcesadorFalso{Resultado: ports.ResultadoProceso{Estado: "LISTO", Paginas: []string{"Hola Mundo", "Segunda"}}}
	if _, err := e.confirmador(p).Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID); err != nil {
		t.Fatalf("confirmar: %v", err)
	}
	if err := corregir.Ejecutar(context.Background(), otro, tr.ID, r.Documento.ID, []string{"a", "b"}); !errors.Is(err, tramite.ErrNoAsignado) {
		t.Fatalf("esperaba ErrNoAsignado, obtuve %v", err)
	}
	if err := corregir.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID, []string{"solo una"}); !errors.Is(err, tramite.ErrPaginasInvalidas) {
		t.Fatalf("cambiar el número de páginas esperaba ErrPaginasInvalidas, obtuve %v", err)
	}
	if err := corregir.Ejecutar(context.Background(), analista, tr.ID, r.Documento.ID, []string{"  Hola, mundo  ", "Segunda página"}); err != nil {
		t.Fatalf("corregir: %v", err)
	}
	paginas, _ := (&tramites.TextoDocumento{Tramites: e.repo}).Ejecutar(context.Background(), tr.ID, r.Documento.ID, false)
	d, _ := e.repo.BuscarDocumento(context.Background(), tr.ID, r.Documento.ID)
	if len(paginas) != 2 || paginas[0] != "Hola, mundo" || !d.TextoValidado {
		t.Fatalf("el texto corregido debe quedar validado y por página: %q %+v", paginas, d)
	}
}

func TestValidarSolicitud_LanzaElOCRDeLosAdjuntosDelPortal(t *testing.T) {
	e := nuevoEntorno()
	p := &testdobles.ProcesadorFalso{Resultado: ports.ResultadoProceso{Estado: "LISTO", Paginas: []string{"Cédula de identidad"}}}
	codigo, err := e.publico().Ejecutar(context.Background(), tramites.RadicarPublicoEntrada{
		Codigo: "123", Cedula: "7654321", CategoriaID: 1,
		Archivos: []tramites.Archivo{{RequisitoID: 10, Nombre: "titulo.pdf", Datos: pdf}, {RequisitoID: 11, Nombre: "cedula.pdf", Datos: pdf}},
	})
	if err != nil {
		t.Fatalf("radicar: %v", err)
	}
	tr, _ := e.repo.BuscarPorCodigo(context.Background(), codigo)
	if len(p.Recibidos) != 0 || tr.Documentos[0].EstadoOCR != tramite.OCRNoAplica {
		t.Fatalf("una solicitud sin validar no debe ocupar el OCR: %d %+v", len(p.Recibidos), tr.Documentos[0])
	}
	validar := &tramites.ValidarSolicitud{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj, OCR: &tramites.ProcesoOCR{Tramites: e.repo, Procesador: p}}
	if err := validar.Ejecutar(context.Background(), director, tr.ID); err != nil {
		t.Fatalf("validar: %v", err)
	}
	tr, _ = e.repo.BuscarPorCodigo(context.Background(), codigo)
	if len(p.Recibidos) != 2 {
		t.Fatalf("al validar, los dos adjuntos deben ir al OCR: %d", len(p.Recibidos))
	}
	for _, d := range tr.Documentos {
		if d.EstadoOCR != tramite.OCRProcesado || len(e.repo.Paginas[d.ID]) != 1 {
			t.Fatalf("el adjunto del portal debe quedar procesado e indexado: %+v", d)
		}
	}
}

func TestRadicarPublico_GuardaEnElDirectorioDelDocente(t *testing.T) {
	e := nuevoEntorno()
	tr := e.conDocumentos(t)
	for _, d := range tr.Documentos {
		if !strings.HasPrefix(d.ClaveObjeto, "docentes/123/"+tr.ID+"/") {
			t.Fatalf("clave fuera del directorio del docente: %s", d.ClaveObjeto)
		}
	}
}
