package tramite

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrNoEncontrado            erroDominio = "trámite no encontrado"
	ErrCodigoDuplicado         erroDominio = "el código de seguimiento ya existe"
	ErrTransicionInvalida      erroDominio = "el trámite no puede pasar a ese estado desde el actual"
	ErrArchivado               erroDominio = "el trámite está archivado y es de solo lectura"
	ErrRemitido                erroDominio = "el trámite fue remitido al Rectorado y es de solo lectura"
	ErrCarpetaNoSolicitada     erroDominio = "antes de revisar hay que solicitar la carpeta física del trámite"
	ErrCarpetasAbiertas        erroDominio = "hay carpetas del trámite sin devolver al archivo"
	ErrPeriodoRequerido        erroDominio = "registra los años a evaluar antes de dictaminar"
	ErrPeriodoInvalido         erroDominio = "los años a evaluar no son válidos"
	ErrResolucionRequerida     erroDominio = "indica el número y la fecha de la resolución"
	ErrNotaRequerida           erroDominio = "indica la nota de remisión al Rectorado"
	ErrMotivoRequerido         erroDominio = "la reasignación exige un motivo"
	ErrNoAsignado              erroDominio = "el trámite no está asignado a este usuario"
	ErrAnalistaInvalido        erroDominio = "el usuario no es un analista activo"
	ErrObservacionesPendientes erroDominio = "hay observaciones pendientes"
	ErrObservacionNoEncontrada erroDominio = "observación no encontrada"
	ErrPendienteValidacion     erroDominio = "la solicitud del portal aún no fue validada en ventanilla"
	ErrYaValidado              erroDominio = "la solicitud ya fue validada"
	ErrDocenteNoEncontrado     erroDominio = "no existe ese docente en el padrón"
	ErrDocenteNoVerificado     erroDominio = "el código y la cédula no corresponden a un docente activo del padrón"
	ErrAdjuntosLibres          erroDominio = "una regularización exige entre 1 y 5 documentos adjuntos"
	ErrDemasiadasPendientes    erroDominio = "el docente ya tiene solicitudes pendientes de validación en ventanilla"
	ErrDocenteInactivo         erroDominio = "el docente está dado de baja"
	ErrCategoriaInvalida       erroDominio = "la categoría no existe o está inactiva"
	ErrArchivoInvalido         erroDominio = "el archivo adjunto no es válido"
	ErrDocumentoNoEncontrado   erroDominio = "documento no encontrado"
	ErrDocumentoDeBaja         erroDominio = "el documento ya está dado de baja"
	ErrMotivoBajaRequerido     erroDominio = "la baja del documento exige un motivo"
	ErrAlmacenNoDisponible     erroDominio = "el almacenamiento de documentos no está disponible"
	ErrProcesadorNoDisponible  erroDominio = "el servicio de OCR no está disponible"
	ErrTipoNoPermitido         erroDominio = "solo se admiten archivos PDF, JPG o PNG"
	ErrArchivoGrande           erroDominio = "el archivo supera el tamaño máximo permitido"
	ErrSubidaIncompleta        erroDominio = "el archivo todavía no se subió al almacén"
	ErrDocumentoNoSubido       erroDominio = "el documento no está en un estado que permita esta operación"
	ErrSinOCR                  erroDominio = "el documento aún no tiene texto extraído"
	ErrTipoDocumentoInvalido   erroDominio = "el tipo de documento no es válido"
	ErrConsultaInvalida       erroDominio = "la búsqueda debe tener entre 2 y 100 caracteres"
	ErrPaginasInvalidas        erroDominio = "el texto corregido debe conservar las páginas del documento"
	ErrPlazoInvalido           erroDominio = "el plazo no puede estar en el pasado"
	ErrDocumentoActivo         erroDominio = "el documento no está dado de baja"
	ErrMotivoRestitucion       erroDominio = "la restitución del documento exige un motivo"
)

// ErrDocumentosFaltantes lista los requisitos de la categoría que no se adjuntaron.
type ErrDocumentosFaltantes struct{ Faltan []string }

func (e *ErrDocumentosFaltantes) Error() string {
	return "faltan documentos exigidos por la categoría"
}

type Estado string

const (
	Ingresado   Estado = "INGRESADO"
	Asignado    Estado = "ASIGNADO"
	EnRevision  Estado = "EN_REVISION"
	Observado   Estado = "OBSERVADO"
	Dictaminado Estado = "DICTAMINADO"
	// EnFacultad: la carpeta salió del Departamento (Facultad, Tribunal, Decano, Consejo Facultativo).
	// En esta etapa el sistema no puede rastrearla: solo se registran la salida y el regreso.
	EnFacultad    Estado = "EN_FACULTAD"
	RevisionFinal Estado = "REVISION_FINAL" // el mismo analista revisa la resolución recibida
	Remitido      Estado = "REMITIDO"       // remitido al Rectorado; falta que Archivo archive la carpeta
	Archivado     Estado = "ARCHIVADO"
)

type Tipo string

const (
	Categorizacion Tipo = "CATEGORIZACION"
	Regularizacion Tipo = "REGULARIZACION"
)

func (t Tipo) Valido() bool { return t == Categorizacion || t == Regularizacion }

type Origen string

const (
	OrigenStaff  Origen = "STAFF"
	OrigenPortal Origen = "PORTAL"
)

type Resultado string

const (
	Favorable    Resultado = "FAVORABLE"
	Desfavorable Resultado = "DESFAVORABLE"
)

func (r Resultado) Valido() bool { return r == Favorable || r == Desfavorable }

// transiciones es la máquina de estados del trámite. La reasignación no cambia de estado,
// por eso no aparece aquí.
var transiciones = map[Estado][]Estado{
	Ingresado:   {Asignado},
	Asignado:    {EnRevision},
	EnRevision:  {Observado, Dictaminado},
	Observado:     {EnRevision},
	Dictaminado:   {EnFacultad},
	EnFacultad:    {RevisionFinal},
	RevisionFinal: {Remitido},
	Remitido:      {Archivado},
}

// Validar comprueba que desde → hasta es una transición permitida.
func Validar(desde, hasta Estado) error {
	if desde == Archivado {
		return ErrArchivado
	}
	for _, e := range transiciones[desde] {
		if e == hasta {
			return nil
		}
	}
	return ErrTransicionInvalida
}

// Activo indica si el trámite cuenta para la carga vigente de un analista. En facultad no cuenta:
// el expediente está fuera del Departamento y el analista queda libre hasta su regreso.
func (e Estado) Activo() bool {
	return e == Asignado || e == EnRevision || e == Observado || e == Dictaminado || e == RevisionFinal
}

// Cerrado indica que el analista ya no puede modificar el trámite (remitido o archivado).
func (e Estado) Cerrado() bool { return e == Remitido || e == Archivado }

// ErrCerrado es el error de solo lectura que corresponde al estado cerrado.
func (e Estado) ErrCerrado() error {
	if e == Remitido {
		return ErrRemitido
	}
	return ErrArchivado
}

type Tramite struct {
	ID                  string
	Codigo              string
	DocenteID           string
	DocenteCodigo       string
	DocenteNombre       string
	CategoriaID         int16 // 0 = sin categoría (las regularizaciones no postulan a ninguna)
	Categoria           string
	Tipo                Tipo
	Estado              Estado
	Origen              Origen
	PendienteValidacion bool
	AnalistaID          *string
	AnalistaNombre      string
	RadicadoPor         *string
	Resultado           Resultado
	Dictamen            string
	DictaminadoEn       *time.Time
	ArchivadoEn         *time.Time
	UbicacionLegajo     string
	CreadoEn            time.Time

	// Años evaluados (se verifican antes de dictaminar), resolución recibida de la facultad y nota de
	// remisión al Rectorado.
	AnioDesde, AnioHasta int16
	ResolucionNumero     string
	ResolucionFecha      *time.Time
	RemisionNota         string
	// Solicitudes de carpeta ligadas: no rechazadas, y las que aún no volvieron al archivo.
	CarpetasLigadas, CarpetasAbiertas int64

	Historial     []Historial
	Observaciones []Observacion
	Documentos    []Documento
	Requisitos    []Requisito
}

// Requisito es un documento exigido por la categoría del trámite y si ya se adjuntó.
type Requisito struct {
	ID       int32
	Nombre   string
	Presente bool
}

type Historial struct {
	ID             int64
	Estado         Estado
	EstadoAnterior Estado
	UsuarioID      *string
	UsuarioNombre  string
	Detalle        string
	Fecha          time.Time
}

type Observacion struct {
	ID         string
	Documento  string
	Motivo     string
	Plazo      time.Time
	Resuelta   bool
	ResueltaEn *time.Time
	CreadoEn   time.Time
}

// Documento es un adjunto del expediente. La baja es lógica (CU-16): el objeto nunca se borra del
// almacén y conserva el motivo y la fecha para que el Director pueda restituirlo.
type Documento struct {
	ID             string
	DocenteID      string // solo en los escaneos del expediente digital del docente (sin trámite)
	TipoDocumento  TipoDocumento
	RequisitoID    *int32
	NombreOriginal string
	TipoMime       string
	Tamano         int64
	HashSHA256     string
	Bucket         string
	ClaveObjeto    string
	Activo         bool
	MotivoBaja     string
	BajaEn         *time.Time
	BajaPor        *string

	// CU-15: carga con OCR. El original se conserva; ClaveComprimida apunta a la versión reducida con texto buscable.
	EstadoOCR        EstadoOCR
	TextoOCR         string
	TextoValidado    bool
	ErrorOCR         string
	ClaveComprimida  string
	TamanoComprimido int64
}

// TipoDocumento clasifica un escaneo del expediente digital del docente (lo sube el Personal de Archivo).
type TipoDocumento string

const (
	TipoCurriculum     TipoDocumento = "CURRICULUM"
	TipoRegularizacion TipoDocumento = "REGULARIZACION"
	TipoAdjunto        TipoDocumento = "ADJUNTO"
	TipoAsociado       TipoDocumento = "ASOCIADO"
	TipoTitularA       TipoDocumento = "TITULAR_A"
	TipoTitularB       TipoDocumento = "TITULAR_B"
	TipoTitularC       TipoDocumento = "TITULAR_C"
	TipoEmerito        TipoDocumento = "EMERITO"
)

var TiposDocumento = []TipoDocumento{
	TipoCurriculum, TipoRegularizacion, TipoAdjunto, TipoAsociado, TipoTitularA, TipoTitularB, TipoTitularC, TipoEmerito,
}

func (t TipoDocumento) Valido() bool {
	for _, v := range TiposDocumento {
		if v == t {
			return true
		}
	}
	return false
}

// CoincidenciaDocente es una página de un documento del expediente donde aparece lo buscado.
type CoincidenciaDocente struct {
	DocumentoID    string
	TipoDocumento  TipoDocumento
	NombreOriginal string
	Coincidencia
}

// DocumentoDeBaja es un documento retirado por baja lógica, con lo necesario para decidir si se
// restituye (CU-17). TramiteID va vacío en los escaneos del expediente del docente.
type DocumentoDeBaja struct {
	Documento
	TramiteID     string
	TramiteCodigo string
	TramiteEstado Estado
	DocenteNombre string
	DocenteCodigo string
	BajaPorNombre string
}

// CampoBusqueda acota la búsqueda global de documentos (CU-18).
type CampoBusqueda string

const (
	BuscarTodo      CampoBusqueda = ""          // contenido o docente
	BuscarContenido CampoBusqueda = "CONTENIDO" // texto OCR (incluye números de resolución)
	BuscarDocente   CampoBusqueda = "DOCENTE"   // nombre, código o cédula del docente
)

func (c CampoBusqueda) Valido() bool {
	return c == BuscarTodo || c == BuscarContenido || c == BuscarDocente
}

type FiltroBusqueda struct {
	Consulta   string
	Campo      CampoBusqueda
	Tipo       TipoDocumento // vacío = cualquiera
	FacultadID int16         // 0 = cualquiera (facultad del docente)
	Recientes  bool          // true = los más recientes primero; si no, por relevancia
}

// ResultadoBusqueda es un documento activo hallado por la búsqueda global: las páginas con
// coincidencia (vacías si solo coincide el docente) y el mejor fragmento, marcado con ⟦ ⟧.
type ResultadoBusqueda struct {
	DocumentoID    string
	TramiteID      string
	TramiteCodigo  string
	TipoDocumento  TipoDocumento
	NombreOriginal string
	DocenteID      string
	DocenteNombre  string
	DocenteCodigo  string
	Paginas        []int
	Fragmento      string
	Rango          float64
}

// EstadoOCR es el avance del procesamiento de un documento cargado por el analista.
type EstadoOCR string

const (
	OCRNoAplica   EstadoOCR = "NO_APLICA" // adjuntos del portal público, sin OCR
	OCRSubiendo   EstadoOCR = "SUBIENDO"  // URL de subida emitida, el archivo aún no llegó al almacén
	OCRPendiente  EstadoOCR = "PENDIENTE"
	OCRProcesando EstadoOCR = "PROCESANDO"
	OCRProcesado  EstadoOCR = "PROCESADO"
	OCRError      EstadoOCR = "ERROR"
)

// ClaveDescarga es el objeto que se entrega al descargar: la versión comprimida si existe.
func (d *Documento) ClaveDescarga() string {
	if d.ClaveComprimida != "" {
		return d.ClaveComprimida
	}
	return d.ClaveObjeto
}

// Coincidencia es una página del documento donde aparece lo buscado. Fragmento trae las palabras
// halladas entre ⟦ y ⟧ (las formas reales del texto, no su raíz).
type Coincidencia struct {
	Pagina    int
	Fragmento string
	Rango     float64
}

// DocumentoEnProceso es un documento cuyo OCR quedó a medias, con el trámite al que pertenece.
type DocumentoEnProceso struct {
	TramiteID string // vacío en los escaneos del expediente del docente
	Documento Documento
}

type Resumen struct {
	ID                  string
	Codigo              string
	DocenteNombre       string
	DocenteCodigo       string
	Categoria           string
	Tipo                Tipo
	Estado              Estado
	Origen              Origen
	PendienteValidacion bool
	AnalistaNombre      string
	CreadoEn            time.Time
}

type Filtro struct {
	Estado     Estado
	AnalistaID string // vacío = cualquiera
	SinAsignar bool
	Busqueda   string // código de trámite o de docente, cédula o nombre
}

type Analista struct {
	ID     string
	Nombre string
	Carga  int64 // trámites activos asignados
}

type Repository interface {
	// Crear inserta el trámite, sus documentos y la primera fila de historial.
	Crear(ctx context.Context, t *Tramite) error
	BuscarPorID(ctx context.Context, id string) (*Tramite, error)
	BuscarPorCodigo(ctx context.Context, codigo string) (*Tramite, error)
	Listar(ctx context.Context, f Filtro) ([]Resumen, error)
	// ListarPagina es Listar con paginación; devuelve también el total de coincidencias.
	ListarPagina(ctx context.Context, f Filtro, pagina, tamano int) ([]Resumen, int64, error)
	// Guardar persiste los campos mutables del trámite y anexa h al historial (si no es nil).
	Guardar(ctx context.Context, t *Tramite, h *Historial) error
	AgregarObservacion(ctx context.Context, tramiteID string, o *Observacion, creadaPor string) error
	BuscarObservacion(ctx context.Context, tramiteID, id string) (*Observacion, error)
	ResolverObservacion(ctx context.Context, id string, en time.Time) error
	ContarObservacionesPendientes(ctx context.Context, tramiteID string) (int64, error)
	BuscarDocumento(ctx context.Context, tramiteID, id string) (*Documento, error)
	// DarDeBajaDocumento marca el adjunto como inactivo (CU-16). Es condicional: si el documento
	// ya estaba de baja devuelve ErrDocumentoDeBaja.
	DarDeBajaDocumento(ctx context.Context, tramiteID, id, motivo, actorID string, en time.Time) error
	// CrearDocumento agrega un adjunto a un trámite existente (CU-15).
	CrearDocumento(ctx context.Context, tramiteID string, d *Documento) error
	// GuardarDocumento persiste los campos mutables del documento (hash, tamaños, OCR y compresión).
	GuardarDocumento(ctx context.Context, tramiteID string, d *Documento) error
	// ListarDocumentosEnProceso devuelve los documentos PENDIENTES o PROCESANDO (para reanudar tras un reinicio).
	ListarDocumentosEnProceso(ctx context.Context) ([]DocumentoEnProceso, error)
	// GuardarPaginas reemplaza el texto OCR del documento, una entrada por página (la 1.ª es la página 1),
	// y deja en texto_ocr su concatenación.
	GuardarPaginas(ctx context.Context, tramiteID, docID string, paginas []string) error
	ListarPaginas(ctx context.Context, tramiteID, docID string) ([]string, error)
	// BuscarEnDocumento devuelve, en orden de página, las páginas que contienen todas las palabras buscadas.
	BuscarEnDocumento(ctx context.Context, tramiteID, docID, consulta string) ([]Coincidencia, error)
	ListarAnalistas(ctx context.Context) ([]Analista, error)
	// ListarPorDocente devuelve todos los trámites del docente, el más reciente primero (portal público).
	ListarPorDocente(ctx context.Context, docenteID string) ([]Resumen, error)

	// Expediente digital del docente: sus escaneos viven en la misma tabla de documentos, sin trámite
	// (tramiteID "" en los métodos de arriba los alcanza).
	CrearDocumentoDocente(ctx context.Context, docenteID string, d *Documento) error
	ListarDocumentosDocente(ctx context.Context, docenteID string, incluirBajas bool) ([]Documento, error)
	// BuscarEnDocente busca en el texto OCR de todos los documentos activos del docente.
	BuscarEnDocente(ctx context.Context, docenteID, consulta string) ([]CoincidenciaDocente, error)

	// CU-17: documentos dados de baja (de trámites y de expedientes) y su restitución.
	ListarDocumentosDeBaja(ctx context.Context, busqueda string, pagina, tamano int) ([]DocumentoDeBaja, int64, error)
	// ObtenerDocumentoDeBaja trae el documento con su contexto, esté o no de baja (ErrDocumentoNoEncontrado si no existe).
	ObtenerDocumentoDeBaja(ctx context.Context, id string) (*DocumentoDeBaja, error)
	// RestituirDocumento lo reactiva y borra los datos de la baja. Es condicional: si ya estaba activo devuelve ErrDocumentoActivo.
	RestituirDocumento(ctx context.Context, id string) error

	// CU-18: búsqueda de texto completo en todos los documentos activos.
	BuscarDocumentos(ctx context.Context, f FiltroBusqueda, pagina, tamano int) ([]ResultadoBusqueda, int64, error)
}

const alfabetoCodigo = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // sin 0/O/1/I

// NuevoCodigo genera un código de seguimiento ED-AAAA-XXXXXX con crypto/rand.
func NuevoCodigo(anio int) (string, error) {
	suf := make([]byte, 6)
	for i := range suf {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alfabetoCodigo))))
		if err != nil {
			return "", err
		}
		suf[i] = alfabetoCodigo[n.Int64()]
	}
	return fmt.Sprintf("ED-%d-%s", anio, suf), nil
}
