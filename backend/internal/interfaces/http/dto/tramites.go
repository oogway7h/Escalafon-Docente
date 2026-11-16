package dto

import "time"

// --- CU-07: docentes y categorías ---

type LegajoDTO struct {
	Gaveta      string `json:"gaveta"`
	Orden       int    `json:"orden"`
	Anexos      string `json:"anexos"`
	Observacion string `json:"observacion"`
}

type FacultadDTO struct {
	ID     int16  `json:"id"`
	Nombre string `json:"nombre"`
}

type DocenteDTO struct {
	ID          string      `json:"id"`
	Codigo      string      `json:"codigo"`
	Cedula      string      `json:"cedula"`
	Nombres     string      `json:"nombres"`
	Apellidos   string      `json:"apellidos"`
	Correo      string      `json:"correo"`
	Telefono    string      `json:"telefono"`
	CategoriaID *int16      `json:"categoriaId"`
	Categoria   string      `json:"categoria"`
	FacultadID  *int16      `json:"facultadId"`
	Facultad    string      `json:"facultad"`
	Legajos     []LegajoDTO `json:"legajos"`
	Activo      bool        `json:"activo"`
}

type ListaDocentesDTO struct {
	Datos  []DocenteDTO `json:"datos"`
	Total  int64        `json:"total"`
	Pagina int          `json:"pagina"`
	Tamano int          `json:"tamano"`
}

type GuardarDocenteEntrada struct {
	Codigo      string `json:"codigo"` // obligatorio al crear
	Cedula      string `json:"cedula"` // opcional; inmutable una vez asignada
	Nombres     string `json:"nombres" validate:"required"`
	Apellidos   string `json:"apellidos" validate:"required"`
	Correo      string `json:"correo" validate:"omitempty,email"`
	Telefono    string `json:"telefono"`
	CategoriaID *int16 `json:"categoriaId"`
	FacultadID  *int16 `json:"facultadId"`
}

type RequisitoDTO struct {
	ID       int32  `json:"id"`
	Nombre   string `json:"nombre"`
	Presente *bool  `json:"presente,omitempty"`
}

type CategoriaDTO struct {
	ID              int16          `json:"id"`
	Nombre          string         `json:"nombre"`
	AniosAntiguedad int16          `json:"aniosAntiguedad"`
	PuntajeDesde    *int16         `json:"puntajeDesde"`
	PuntajeHasta    *int16         `json:"puntajeHasta"`
	Activa          bool           `json:"activa"`
	Requisitos      []RequisitoDTO `json:"requisitos"`
	Docentes        int64          `json:"docentes"`
	Tramites        int64          `json:"tramites"`
}

type GuardarCategoriaEntrada struct {
	Nombre          string   `json:"nombre" validate:"required"`
	AniosAntiguedad int16    `json:"aniosAntiguedad" validate:"min=0"`
	PuntajeDesde    *int16   `json:"puntajeDesde" validate:"omitempty,min=0"`
	PuntajeHasta    *int16   `json:"puntajeHasta" validate:"omitempty,min=0"`
	Requisitos      []string `json:"requisitos"`
}

// --- CU-08 al CU-11: trámites ---

type ResumenTramiteDTO struct {
	ID                  string    `json:"id"`
	Codigo              string    `json:"codigo"`
	DocenteNombre       string    `json:"docenteNombre"`
	DocenteCodigo       string    `json:"docenteCodigo"`
	Categoria           string    `json:"categoria"`
	Tipo                string    `json:"tipo"`
	Estado              string    `json:"estado"`
	Origen              string    `json:"origen"`
	PendienteValidacion bool      `json:"pendienteValidacion"`
	AnalistaNombre      string    `json:"analistaNombre"`
	CreadoEn            time.Time `json:"creadoEn"`
}

type HistorialDTO struct {
	Estado         string    `json:"estado"`
	EstadoAnterior string    `json:"estadoAnterior"`
	Usuario        string    `json:"usuario"`
	Detalle        string    `json:"detalle"`
	Fecha          time.Time `json:"fecha"`
}

type ObservacionDTO struct {
	ID         string     `json:"id"`
	Documento  string     `json:"documento"`
	Motivo     string     `json:"motivo"`
	Plazo      string     `json:"plazo"`
	Resuelta   bool       `json:"resuelta"`
	ResueltaEn *time.Time `json:"resueltaEn"`
}

// --- Expediente digital del docente ---

type SubidaExpedienteEntrada struct {
	Nombre        string `json:"nombre" validate:"required,max=200"`
	TipoMime      string `json:"tipoMime" validate:"required"`
	Tamano        int64  `json:"tamano" validate:"required,gt=0"`
	TipoDocumento string `json:"tipoDocumento" validate:"required"`
}

type CoincidenciaExpedienteDTO struct {
	DocumentoID    string  `json:"documentoId"`
	TipoDocumento  string  `json:"tipoDocumento"`
	NombreOriginal string  `json:"nombreOriginal"`
	Pagina         int     `json:"pagina"`
	Fragmento      string  `json:"fragmento"` // las palabras halladas van entre ⟦ y ⟧
	Rango          float64 `json:"rango"`
}

type BusquedaExpedienteDTO struct {
	Coincidencias []CoincidenciaExpedienteDTO `json:"coincidencias"`
}

type DocumentoDTO struct {
	ID             string     `json:"id"`
	TipoDocumento  string     `json:"tipoDocumento"` // solo en los escaneos del expediente del docente
	NombreOriginal string     `json:"nombreOriginal"`
	TipoMime       string     `json:"tipoMime"`
	Tamano         int64      `json:"tamano"`
	RequisitoID    *int32     `json:"requisitoId"`
	Activo         bool       `json:"activo"`
	MotivoBaja     string     `json:"motivoBaja"`
	BajaEn         *time.Time `json:"bajaEn"`

	// CU-15: procesamiento OCR y compresión.
	EstadoOcr        string `json:"estadoOcr"`
	ErrorOcr         string `json:"errorOcr"`
	TextoValidado    bool   `json:"textoValidado"`
	HashSha256       string `json:"hashSha256"`
	TamanoComprimido int64  `json:"tamanoComprimido"`
}

// --- CU-15: carga de documentos con OCR ---

type SubidaEntrada struct {
	Nombre      string `json:"nombre" validate:"required,max=200"`
	TipoMime    string `json:"tipoMime" validate:"required"`
	Tamano      int64  `json:"tamano" validate:"required,gt=0"`
	RequisitoID *int32 `json:"requisitoId"`
}

type SubidaPreparadaDTO struct {
	Documento DocumentoDTO `json:"documento"`
	URL       string       `json:"url"`
}

// TextoDocumentoDTO es el texto OCR de un documento, una entrada por página.
type TextoDocumentoDTO struct {
	Paginas []string `json:"paginas" validate:"required"`
}

type CoincidenciaDTO struct {
	Pagina    int     `json:"pagina"`
	Fragmento string  `json:"fragmento"` // las palabras halladas van entre ⟦ y ⟧
	Rango     float64 `json:"rango"`
}

type BusquedaDocumentoDTO struct {
	Coincidencias []CoincidenciaDTO `json:"coincidencias"`
}

type ListaTramitesDTO struct {
	Datos  []ResumenTramiteDTO `json:"datos"`
	Total  int64               `json:"total"`
	Pagina int                 `json:"pagina"`
	Tamano int                 `json:"tamano"`
}

type TramiteDetalleDTO struct {
	ResumenTramiteDTO
	DocenteID       string           `json:"docenteId"`
	CategoriaID     int16            `json:"categoriaId"`
	AnalistaID      *string          `json:"analistaId"`
	Resultado       string           `json:"resultado"`
	Dictamen        string           `json:"dictamen"`
	DictaminadoEn   *time.Time       `json:"dictaminadoEn"`
	ArchivadoEn     *time.Time       `json:"archivadoEn"`
	UbicacionLegajo string           `json:"ubicacionLegajo"`
	AnioDesde        int16           `json:"anioDesde"`
	AnioHasta        int16           `json:"anioHasta"`
	ResolucionNumero string          `json:"resolucionNumero"`
	ResolucionFecha  *string         `json:"resolucionFecha"`
	RemisionNota     string          `json:"remisionNota"`
	CarpetasLigadas  int64           `json:"carpetasLigadas"`
	CarpetasAbiertas int64           `json:"carpetasAbiertas"`
	Historial       []HistorialDTO   `json:"historial"`
	Observaciones   []ObservacionDTO `json:"observaciones"`
	Documentos      []DocumentoDTO   `json:"documentos"`
	Requisitos      []RequisitoDTO   `json:"requisitos"`
}

type AnalistaDTO struct {
	ID     string `json:"id"`
	Nombre string `json:"nombre"`
	Carga  int64  `json:"carga"`
}

type RadicarEntrada struct {
	DocenteID   string `json:"docenteId" validate:"required,uuid"`
	Tipo        string `json:"tipo" validate:"required,oneof=CATEGORIZACION REGULARIZACION"`
	CategoriaID int16  `json:"categoriaId"` // obligatoria solo en categorización (lo valida el caso de uso)
}

type AsignarEntrada struct {
	AnalistaID string `json:"analistaId" validate:"required,uuid"`
	Motivo     string `json:"motivo"` // obligatorio al reasignar
}

type ObservacionEntrada struct {
	Documento string `json:"documento" validate:"required"`
	Motivo    string `json:"motivo" validate:"required,min=5"`
	Plazo     string `json:"plazo" validate:"required"` // AAAA-MM-DD
}

type DictamenEntrada struct {
	Resultado string `json:"resultado" validate:"required,oneof=FAVORABLE DESFAVORABLE"`
	Texto     string `json:"texto" validate:"required,min=10"`
}

// --- Etapas posteriores al dictamen ---

type PeriodoEntrada struct {
	AnioDesde int16 `json:"anioDesde" validate:"required"`
	AnioHasta int16 `json:"anioHasta" validate:"required"`
}

type RegresoEntrada struct {
	ResolucionNumero string `json:"resolucionNumero" validate:"required"`
	ResolucionFecha  string `json:"resolucionFecha" validate:"required"` // AAAA-MM-DD
}

type RemisionEntrada struct {
	Nota string `json:"nota" validate:"required"`
}

type ArchivarEntrada struct {
	Ubicacion string `json:"ubicacion" validate:"required"`
}

// --- CU-16: baja lógica de un documento ---

type BajaDocumentoEntrada struct {
	Motivo string `json:"motivo" validate:"required,min=5"`
}

// --- CU-12 / CU-13: portal público ---

type EtapaPublicaDTO struct {
	Estado string    `json:"estado"`
	Fecha  time.Time `json:"fecha"`
}

type ObservacionPublicaDTO struct {
	Documento string `json:"documento"`
	Motivo    string `json:"motivo"`
	Plazo     string `json:"plazo"`
	Resuelta  bool   `json:"resuelta"`
}

type VistaPublicaDTO struct {
	Codigo              string                  `json:"codigo"`
	Tipo                string                  `json:"tipo"`
	Categoria           string                  `json:"categoria"`
	Estado              string                  `json:"estado"`
	RadicadoEn          time.Time               `json:"radicadoEn"`
	PendienteValidacion bool                    `json:"pendienteValidacion"`
	Resultado           string                  `json:"resultado"`
	Etapas              []EtapaPublicaDTO       `json:"etapas"`
	Observaciones       []ObservacionPublicaDTO `json:"observaciones"`
}

// ConsultaDocenteEntrada: el docente que perdió su código de seguimiento se identifica con su código y su cédula.
type ConsultaDocenteEntrada struct {
	Codigo string `json:"codigo" validate:"required,max=20"`
	Cedula string `json:"cedula" validate:"required,max=20"`
}

// ResumenPublicoDTO es un trámite en la lista que ve el docente al buscar por su código de docente.
type ResumenPublicoDTO struct {
	Codigo              string    `json:"codigo"`
	Tipo                string    `json:"tipo"`
	Categoria           string    `json:"categoria"`
	Estado              string    `json:"estado"`
	RadicadoEn          time.Time `json:"radicadoEn"`
	PendienteValidacion bool      `json:"pendienteValidacion"`
}

type SolicitudRadicadaDTO struct {
	Codigo string `json:"codigo"`
}

// --- CU-17 y CU-18: bajas, restitución y búsqueda global ---

type RestituirEntrada struct {
	Motivo string `json:"motivo" validate:"required,min=5"`
}

type DocumentoDeBajaDTO struct {
	DocumentoDTO
	TramiteID     string `json:"tramiteId"`
	TramiteCodigo string `json:"tramiteCodigo"`
	TramiteEstado string `json:"tramiteEstado"`
	DocenteID     string `json:"docenteId"`
	DocenteNombre string `json:"docenteNombre"`
	DocenteCodigo string `json:"docenteCodigo"`
	BajaPorNombre string `json:"bajaPorNombre"`
}

type ListaBajasDTO struct {
	Datos  []DocumentoDeBajaDTO `json:"datos"`
	Total  int64                `json:"total"`
	Pagina int                  `json:"pagina"`
	Tamano int                  `json:"tamano"`
}

type ResultadoBusquedaDTO struct {
	DocumentoID    string `json:"documentoId"`
	TramiteID      string `json:"tramiteId"`
	TramiteCodigo  string `json:"tramiteCodigo"`
	TipoDocumento  string `json:"tipoDocumento"`
	NombreOriginal string `json:"nombreOriginal"`
	DocenteID      string `json:"docenteId"`
	DocenteNombre  string `json:"docenteNombre"`
	DocenteCodigo  string `json:"docenteCodigo"`
	Paginas        []int  `json:"paginas"`
	Fragmento      string `json:"fragmento"` // palabras halladas entre ⟦ ⟧
}

type ListaResultadosDTO struct {
	Datos  []ResultadoBusquedaDTO `json:"datos"`
	Total  int64                  `json:"total"`
	Pagina int                    `json:"pagina"`
	Tamano int                    `json:"tamano"`
}
