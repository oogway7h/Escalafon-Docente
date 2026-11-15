package tramites

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/tramite"
)

// --- CU-15 Cargar y procesar documento con OCR -------------------------------------------
//
// El archivo nunca pasa por la API: PrepararSubida emite una URL firmada, el navegador sube a S3,
// ConfirmarSubida verifica el objeto y encola el procesamiento en el ia-service (Tesseract +
// compresión), cuyo resultado queda en el documento.

const (
	maxPDF    = 400 << 20 // los escaneos del archivo llegan a cientos de MB
	maxImagen = 20 << 20
)

func limitePorTipo(mime string) int64 {
	if mime == "application/pdf" {
		return maxPDF
	}
	return maxImagen
}

// carpetaDocente es el prefijo de las claves de objeto de un docente: su código universitario
// (o su id si aún no tiene).
func carpetaDocente(codigo, docenteID string) string {
	if codigo = strings.TrimSpace(codigo); codigo != "" {
		return "docentes/" + caracteresNoSeguros.ReplaceAllString(codigo, "_")
	}
	return "docentes/" + docenteID
}

// claveDocumento ubica el adjunto en el directorio del docente: docentes/<código>/<trámite>/<doc>-<nombre>.
func claveDocumento(carpeta, tramiteID, docID, nombre string) string {
	nombre = caracteresNoSeguros.ReplaceAllString(nombre, "_")
	if len(nombre) > 100 {
		nombre = nombre[len(nombre)-100:]
	}
	return fmt.Sprintf("%s/%s/%s-%s", carpeta, tramiteID, docID, nombre)
}

type PrepararSubida struct {
	Tramites tramite.Repository
	Almacen  ports.ObjectStorage
	Reloj    ports.Clock
}

type SubidaEntrada struct {
	ActorID, TramiteID, Nombre, TipoMime string
	Tamano                               int64
	RequisitoID                          *int32 // documento del catálogo que cubre; nil = documento adicional
}

type SubidaPreparada struct {
	Documento tramite.Documento
	URL       string
}

// Ejecutar valida el archivo declarado, registra el documento en estado SUBIENDO y devuelve la URL de subida.
func (c *PrepararSubida) Ejecutar(ctx context.Context, in SubidaEntrada) (*SubidaPreparada, error) {
	t, err := cargarDelAnalista(ctx, c.Tramites, in.ActorID, in.TramiteID)
	if err != nil {
		return nil, err
	}
	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" || in.Tamano <= 0 {
		return nil, tramite.ErrArchivoInvalido
	}
	if !tiposPermitidos[in.TipoMime] {
		return nil, tramite.ErrTipoNoPermitido
	}
	if in.Tamano > limitePorTipo(in.TipoMime) {
		return nil, tramite.ErrArchivoGrande
	}
	if in.RequisitoID != nil {
		ok := false
		for _, r := range t.Requisitos {
			ok = ok || r.ID == *in.RequisitoID
		}
		if !ok {
			return nil, tramite.ErrArchivoInvalido
		}
	}
	if r := []rune(nombre); len(r) > 200 {
		nombre = string(r[len(r)-200:])
	}

	id := uuid.NewString()
	clave := claveDocumento(carpetaDocente(t.DocenteCodigo, t.DocenteID), t.ID, id, nombre)
	url, err := c.Almacen.URLSubida(ctx, clave, in.TipoMime, time.Hour)
	if err != nil {
		return nil, err
	}
	d := tramite.Documento{
		ID: id, RequisitoID: in.RequisitoID, NombreOriginal: nombre, TipoMime: in.TipoMime, Tamano: in.Tamano,
		Bucket: c.Almacen.Bucket(), ClaveObjeto: clave, Activo: true, EstadoOCR: tramite.OCRSubiendo,
	}
	if err := c.Tramites.CrearDocumento(ctx, t.ID, &d); err != nil {
		return nil, err
	}
	return &SubidaPreparada{Documento: d, URL: url}, nil
}

type ConfirmarSubida struct {
	Tramites  tramite.Repository
	Almacen   ports.ObjectStorage
	OCR       *ProcesoOCR
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

// Ejecutar comprueba que el objeto llegó al almacén y lanza el procesamiento. También sirve para
// reintentar un documento cuyo procesamiento falló (estado ERROR).
func (c *ConfirmarSubida) Ejecutar(ctx context.Context, actorID, tramiteID, docID string) (*tramite.Documento, error) {
	t, err := cargarDelAnalista(ctx, c.Tramites, actorID, tramiteID)
	if err != nil {
		return nil, err
	}
	return c.confirmar(ctx, t.ID, "", docID, func(d *tramite.Documento, tam int64) {
		registrar(ctx, c.Auditoria, &actorID, auditoria.AccionCargarDocumento, t.ID, c.Reloj.Ahora(),
			map[string]any{"documentoId": d.ID, "documento": d.NombreOriginal, "tamano": tam})
	})
}

// confirmar es el núcleo compartido por los adjuntos de un trámite (tramiteID) y los escaneos del
// expediente de un docente (tramiteID "" y docenteID). auditar se llama solo en la primera confirmación.
func (c *ConfirmarSubida) confirmar(ctx context.Context, tramiteID, docenteID, docID string, auditar func(d *tramite.Documento, tam int64)) (*tramite.Documento, error) {
	d, err := c.Tramites.BuscarDocumento(ctx, tramiteID, docID)
	if err != nil {
		if errors.Is(err, tramite.ErrNoEncontrado) {
			return nil, tramite.ErrDocumentoNoEncontrado
		}
		return nil, err
	}
	if d.DocenteID != docenteID {
		return nil, tramite.ErrDocumentoNoEncontrado
	}
	if !d.Activo {
		return nil, tramite.ErrDocumentoDeBaja
	}
	if d.EstadoOCR != tramite.OCRSubiendo && d.EstadoOCR != tramite.OCRError {
		return nil, tramite.ErrDocumentoNoSubido
	}
	tam, err := c.Almacen.Tamano(ctx, d.ClaveObjeto)
	if errors.Is(err, tramite.ErrAlmacenNoDisponible) {
		return nil, err
	}
	if err != nil || tam <= 0 {
		return nil, tramite.ErrSubidaIncompleta
	}
	if tam > limitePorTipo(d.TipoMime) {
		return nil, tramite.ErrArchivoGrande
	}
	primera := d.EstadoOCR == tramite.OCRSubiendo
	d.Tamano, d.EstadoOCR, d.ErrorOCR = tam, tramite.OCRPendiente, ""
	if err := c.Tramites.GuardarDocumento(ctx, tramiteID, d); err != nil {
		return nil, err
	}
	if primera {
		auditar(d, tam)
	}
	c.OCR.Lanzar(tramiteID, *d)
	return d, nil
}

// ProcesoOCR lanza y sigue el trabajo del ia-service de cada documento (subidas del analista y adjuntos
// del portal) y guarda su resultado: estado, texto por página, tamaños y hash.
type ProcesoOCR struct {
	Tramites   tramite.Repository
	Procesador ports.ProcesadorDocumentos
	Asincrono  bool          // true en producción: el OCR corre en segundo plano
	Intervalo  time.Duration // pausa entre consultas al ia-service (por defecto 5 s)
	Espera     time.Duration // tiempo máximo por documento (por defecto 6 h)
	enCurso    sync.Map      // ids de documento con un procesamiento activo en esta instancia
}

// Lanzar procesa el documento (que ya debe estar guardado como PENDIENTE). Ignora el aviso si ese
// documento ya se está procesando.
func (c *ProcesoOCR) Lanzar(tramiteID string, d tramite.Documento) {
	if _, ya := c.enCurso.LoadOrStore(d.ID, true); ya {
		return
	}
	if c.Asincrono {
		go c.procesar(context.Background(), tramiteID, d)
		return
	}
	c.procesar(context.Background(), tramiteID, d)
}

// Reanudar vuelve a lanzar los documentos que quedaron PENDIENTES o PROCESANDO cuando la API se
// detuvo (los trabajos del ia-service viven en memoria). Se llama una vez al arrancar.
// ponytail: asume una sola instancia de la API; con varias habría que reclamar cada fila.
func (c *ProcesoOCR) Reanudar(ctx context.Context) error {
	docs, err := c.Tramites.ListarDocumentosEnProceso(ctx)
	if err != nil {
		return err
	}
	for _, e := range docs {
		c.Lanzar(e.TramiteID, e.Documento)
	}
	return nil
}

func (c *ProcesoOCR) procesar(ctx context.Context, tramiteID string, d tramite.Documento) {
	defer c.enCurso.Delete(d.ID)
	guardar := func() {
		if err := c.Tramites.GuardarDocumento(ctx, tramiteID, &d); err != nil {
			slog.Error("no se pudo guardar el estado del documento", "documento", d.ID, "error", err)
		}
	}
	fallar := func(motivo string) {
		slog.Warn("OCR del documento falló", "documento", d.ID, "motivo", motivo)
		if r := []rune(motivo); len(r) > 290 {
			motivo = string(r[:290])
		}
		d.EstadoOCR, d.ErrorOCR = tramite.OCRError, motivo
		guardar()
	}

	d.EstadoOCR = tramite.OCRProcesando
	guardar()
	solicitud := ports.SolicitudProceso{Bucket: d.Bucket, Clave: d.ClaveObjeto, TipoMime: d.TipoMime}
	if d.TipoMime == "application/pdf" {
		solicitud.ClaveSalida = path.Dir(d.ClaveObjeto) + "/" + d.ID + "-comprimido.pdf"
	}
	id, err := c.Procesador.Iniciar(ctx, solicitud)
	if err != nil {
		fallar(err.Error())
		return
	}

	intervalo, espera := c.Intervalo, c.Espera
	if intervalo <= 0 {
		intervalo = 5 * time.Second
	}
	if espera <= 0 {
		espera = 6 * time.Hour
	}
	limite := time.Now().Add(espera)
	fallos := 0
	for {
		r, err := c.Procesador.Estado(ctx, id)
		switch {
		case err != nil:
			if fallos++; fallos >= 5 { // el servicio pudo reiniciarse y perder el trabajo
				fallar("no se pudo consultar el servicio de OCR: " + err.Error())
				return
			}
		case r.Estado == "LISTO":
			if d.HashSHA256 != "" && r.HashSHA256 != "" && !strings.EqualFold(d.HashSHA256, r.HashSHA256) {
				fallar("el hash del archivo no coincide con el registrado al cargarlo")
				return
			}
			paginas := r.Paginas
			if len(paginas) == 0 && strings.TrimSpace(r.Texto) != "" { // servicio anterior: solo texto plano
				paginas = []string{r.Texto}
			}
			if err := c.Tramites.GuardarPaginas(ctx, tramiteID, d.ID, paginas); err != nil {
				fallar("no se pudo guardar el texto extraído: " + err.Error())
				return
			}
			// El texto es buscable de inmediato; corregirlo después es opcional (TextoValidado lo registra).
			d.EstadoOCR, d.ErrorOCR, d.TextoOCR, d.TextoValidado = tramite.OCRProcesado, "", strings.Join(paginas, "\n\n"), false
			if r.HashSHA256 != "" {
				d.HashSHA256 = r.HashSHA256
			}
			if r.TamanoOriginal > 0 {
				d.Tamano = r.TamanoOriginal
			}
			d.ClaveComprimida, d.TamanoComprimido = r.ClaveComprimida, 0
			if r.ClaveComprimida != "" {
				d.TamanoComprimido = r.TamanoFinal
			}
			guardar()
			return
		case r.Estado == "ERROR":
			fallar(r.Error)
			return
		default:
			fallos = 0
		}
		if time.Now().After(limite) {
			fallar("tiempo de espera agotado")
			return
		}
		time.Sleep(intervalo)
	}
}

// CorregirPaginas guarda el texto extraído tras la revisión del analista y lo marca como validado.
// Conserva el número de páginas: el índice y el visor las identifican por su posición.
type CorregirPaginas struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Tx        ports.TxManager
	Reloj     ports.Clock
}

const maxTextoOCR = 5 << 20

func (c *CorregirPaginas) Ejecutar(ctx context.Context, actorID, tramiteID, docID string, paginas []string) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, actorID, tramiteID)
	if err != nil {
		return err
	}
	d, err := c.Tramites.BuscarDocumento(ctx, t.ID, docID)
	if err != nil {
		if errors.Is(err, tramite.ErrNoEncontrado) {
			return tramite.ErrDocumentoNoEncontrado
		}
		return err
	}
	if !d.Activo {
		return tramite.ErrDocumentoDeBaja
	}
	if d.EstadoOCR != tramite.OCRProcesado {
		return tramite.ErrDocumentoNoSubido
	}
	actuales, err := c.Tramites.ListarPaginas(ctx, t.ID, d.ID)
	if err != nil {
		return err
	}
	if len(paginas) != len(actuales) {
		return tramite.ErrPaginasInvalidas
	}
	total := 0
	limpias := make([]string, len(paginas))
	for i, p := range paginas {
		limpias[i] = strings.TrimSpace(p)
		total += len(limpias[i])
	}
	if total > maxTextoOCR {
		return tramite.ErrArchivoGrande
	}
	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := c.Tramites.GuardarPaginas(ctx, t.ID, d.ID, limpias); err != nil {
			return err
		}
		d.TextoOCR, d.TextoValidado = strings.Join(limpias, "\n\n"), true
		return c.Tramites.GuardarDocumento(ctx, t.ID, d)
	})
	if err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionCorregirTextoOCR, t.ID, c.Reloj.Ahora(),
		map[string]any{"documentoId": d.ID, "documento": d.NombreOriginal})
	return nil
}

// documentoLegible carga el documento y comprueba que se pueda leer: uno dado de baja solo lo ve
// quien puede restituirlo.
func documentoLegible(ctx context.Context, repo tramite.Repository, tramiteID, docID string, verBajas bool) (*tramite.Documento, error) {
	d, err := repo.BuscarDocumento(ctx, tramiteID, docID)
	if err != nil {
		if errors.Is(err, tramite.ErrNoEncontrado) {
			return nil, tramite.ErrDocumentoNoEncontrado
		}
		return nil, err
	}
	if !d.Activo && !verBajas {
		return nil, tramite.ErrDocumentoDeBaja
	}
	return d, nil
}

// TextoDocumento devuelve el texto extraído por OCR, una entrada por página (puede pesar varios MB).
type TextoDocumento struct{ Tramites tramite.Repository }

func (c *TextoDocumento) Ejecutar(ctx context.Context, tramiteID, docID string, verBajas bool) ([]string, error) {
	if _, err := documentoLegible(ctx, c.Tramites, tramiteID, docID, verBajas); err != nil {
		return nil, err
	}
	return c.Tramites.ListarPaginas(ctx, tramiteID, docID)
}

// BuscarEnDocumento busca en el índice OCR del documento y devuelve las páginas donde aparece.
type BuscarEnDocumento struct{ Tramites tramite.Repository }

func (c *BuscarEnDocumento) Ejecutar(ctx context.Context, tramiteID, docID, consulta string, verBajas bool) ([]tramite.Coincidencia, error) {
	consulta = strings.TrimSpace(consulta)
	if n := len([]rune(consulta)); n < 2 || n > 100 {
		return nil, tramite.ErrConsultaInvalida
	}
	d, err := documentoLegible(ctx, c.Tramites, tramiteID, docID, verBajas)
	if err != nil {
		return nil, err
	}
	if d.EstadoOCR != tramite.OCRProcesado {
		return nil, tramite.ErrSinOCR
	}
	return c.Tramites.BuscarEnDocumento(ctx, tramiteID, docID, consulta)
}
