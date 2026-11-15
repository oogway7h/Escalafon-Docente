package tramites_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/tramite"
)

// escaneado sube un documento al expediente del docente d1 y lo deja procesado con las páginas dadas.
func (e *entorno) escaneado(t *testing.T, tipo tramite.TipoDocumento, paginas ...string) string {
	t.Helper()
	ctx := context.Background()
	r, err := (&tramites.PrepararSubidaDocente{Docentes: e.docentes, Tramites: e.repo, Almacen: e.almacen}).Ejecutar(ctx,
		tramites.SubidaDocenteEntrada{DocenteID: "d1", Nombre: "escaneo.pdf", TipoMime: "application/pdf", Tipo: tipo, Tamano: 10})
	if err != nil {
		t.Fatalf("preparar: %v", err)
	}
	e.almacen.Guardados[r.Documento.ClaveObjeto] = make([]byte, 10)
	p := &testdobles.ProcesadorFalso{Resultado: ports.ResultadoProceso{Estado: "LISTO", Paginas: paginas}}
	if _, err := e.confirmador(p).EjecutarDocente(ctx, archivo, "d1", r.Documento.ID); err != nil {
		t.Fatalf("confirmar: %v", err)
	}
	return r.Documento.ID
}

func TestPrepararSubidaDocente_ValidaDocenteTipoYArchivo(t *testing.T) {
	e := nuevoEntorno()
	c := &tramites.PrepararSubidaDocente{Docentes: e.docentes, Tramites: e.repo, Almacen: e.almacen}
	entrada := tramites.SubidaDocenteEntrada{DocenteID: "d1", Nombre: "cv.pdf", TipoMime: "application/pdf", Tipo: tramite.TipoCurriculum, Tamano: 10}

	mala := entrada
	mala.DocenteID = "no-existe"
	if _, err := c.Ejecutar(context.Background(), mala); !errors.Is(err, docente.ErrNoEncontrado) {
		t.Fatalf("docente inexistente: %v", err)
	}
	mala = entrada
	mala.Tipo = "OTRO"
	if _, err := c.Ejecutar(context.Background(), mala); !errors.Is(err, tramite.ErrTipoDocumentoInvalido) {
		t.Fatalf("tipo inválido: %v", err)
	}
	mala = entrada
	mala.TipoMime = "text/html"
	if _, err := c.Ejecutar(context.Background(), mala); !errors.Is(err, tramite.ErrTipoNoPermitido) {
		t.Fatalf("tipo de archivo: %v", err)
	}
	r, err := c.Ejecutar(context.Background(), entrada)
	if err != nil {
		t.Fatalf("preparar: %v", err)
	}
	if !strings.HasPrefix(r.Documento.ClaveObjeto, "docentes/123/expediente/") || r.Documento.EstadoOCR != tramite.OCRSubiendo {
		t.Fatalf("va a la carpeta del docente y queda SUBIENDO: %+v", r.Documento)
	}
}

func TestBuscarEnExpediente_EncuentraLaPaginaEntreVariosDocumentos(t *testing.T) {
	e := nuevoEntorno()
	e.escaneado(t, tramite.TipoCurriculum, "Datos personales", "Título de Doctorado en Educación")
	e.escaneado(t, tramite.TipoRegularizacion, "Resolución rectoral 45/2019 de categorización")
	buscar := &tramites.BuscarEnExpediente{Docentes: e.docentes, Tramites: e.repo}

	lista, err := buscar.Ejecutar(context.Background(), "d1", "doctorado")
	if err != nil || len(lista) != 1 || lista[0].Pagina != 2 || lista[0].TipoDocumento != tramite.TipoCurriculum {
		t.Fatalf("esperaba la página 2 del currículum: %+v %v", lista, err)
	}
	if lista, _ := buscar.Ejecutar(context.Background(), "d1", "resolución"); len(lista) != 1 || lista[0].TipoDocumento != tramite.TipoRegularizacion {
		t.Fatalf("esperaba la regularización: %+v", lista)
	}
	if _, err := buscar.Ejecutar(context.Background(), "d1", "a"); !errors.Is(err, tramite.ErrConsultaInvalida) {
		t.Fatalf("consulta corta: %v", err)
	}
}

func TestExpediente_BajaOcultaDeLaBusquedaYNoSeCruzaConOtroDocente(t *testing.T) {
	e := nuevoEntorno()
	e.docentes.Sembrar(docente.Docente{ID: "d2", Codigo: "456", Nombres: "Luis", Apellidos: "Rojas", Activo: true})
	doc := e.escaneado(t, tramite.TipoAdjunto, "Provisión nacional")
	baja := &tramites.DarDeBajaExpediente{Tramites: e.repo, Auditoria: e.aud, Reloj: e.reloj}
	buscar := &tramites.BuscarEnExpediente{Docentes: e.docentes, Tramites: e.repo}

	if err := baja.Ejecutar(context.Background(), archivo, "d1", doc, " "); !errors.Is(err, tramite.ErrMotivoBajaRequerido) {
		t.Fatalf("sin motivo: %v", err)
	}
	if err := baja.Ejecutar(context.Background(), archivo, "d2", doc, "error"); !errors.Is(err, tramite.ErrDocumentoNoEncontrado) {
		t.Fatalf("el documento es de otro docente: %v", err)
	}
	if err := baja.Ejecutar(context.Background(), archivo, "d1", doc, "Escaneo equivocado"); err != nil {
		t.Fatalf("baja: %v", err)
	}
	if lista, _ := buscar.Ejecutar(context.Background(), "d1", "provisión"); len(lista) != 0 {
		t.Fatalf("un documento de baja no se busca: %+v", lista)
	}
	url := &tramites.UrlExpediente{Tramites: e.repo, Almacen: e.almacen}
	if _, err := url.Ejecutar(context.Background(), "d1", doc, false); !errors.Is(err, tramite.ErrDocumentoDeBaja) {
		t.Fatalf("sin permiso de restituir no se descarga una baja: %v", err)
	}
	if _, err := url.Ejecutar(context.Background(), "d1", doc, true); err != nil {
		t.Fatalf("quien restituye sí: %v", err)
	}
}

func TestConsultarPublicoPorDocente_ExigeCodigoYCedulaYNoRevelaSiExiste(t *testing.T) {
	e := nuevoEntorno()
	e.docentes.Sembrar(docente.Docente{ID: "d1", Codigo: "123", Cedula: "1111111", Nombres: "María", Apellidos: "Vargas", Activo: true})
	e.docentes.Sembrar(docente.Docente{ID: "d2", Codigo: "456", Cedula: "9876543", Nombres: "Luis", Apellidos: "Rojas", Activo: true})
	e.docentes.Sembrar(docente.Docente{ID: "d3", Codigo: "789", Nombres: "Sin", Apellidos: "Cédula", Activo: true})
	e.radicar(t) // trámite del docente d1 (código 123)
	caso := &tramites.ConsultarPublicoPorDocente{Docentes: e.docentes, Tramites: e.repo}
	vacio := func(codigo, cedula, motivo string) {
		t.Helper()
		if lista, err := caso.Ejecutar(context.Background(), codigo, cedula); err != nil || len(lista) != 0 {
			t.Fatalf("%s debe responder como «sin trámites»: %+v %v", motivo, lista, err)
		}
	}

	lista, err := caso.Ejecutar(context.Background(), " 123 ", " 1111111 ")
	if err != nil || len(lista) != 1 || lista[0].Codigo == "" || lista[0].Estado != tramite.Ingresado {
		t.Fatalf("código y cédula correctos listan el trámite: %+v %v", lista, err)
	}
	vacio("123", "", "sin cédula")
	vacio("123", "9876543", "cédula de otro docente")
	vacio("123", "0000000", "cédula equivocada")
	vacio("99999", "1111111", "código desconocido")
	vacio("1111111", "1111111", "la cédula usada como código")
	vacio("789", "123", "docente sin cédula registrada")
	vacio("456", "9876543", "docente sin trámites")
}
