package persistence

import "time"

// RolModel es la representación GORM de la tabla rol.
type RolModel struct {
	ID          int16     `gorm:"primaryKey"`
	Nombre      string    `gorm:"type:varchar(50);not null;unique"`
	Descripcion string    `gorm:"type:varchar(200)"`
	EsSistema   bool      `gorm:"not null;default:false"`
	Activo      bool      `gorm:"not null;default:true"`
	CreadoEn    time.Time `gorm:"not null;default:now()"`
}

func (RolModel) TableName() string { return "rol" }

// PermisoModel es la representación GORM de la tabla permiso.
type PermisoModel struct {
	ID      int16  `gorm:"primaryKey"`
	Codigo  string `gorm:"type:varchar(60);not null;unique"`
	Nombre  string `gorm:"type:varchar(120);not null"`
	Modulo  string `gorm:"type:varchar(40);not null"`
	CasoUso string `gorm:"column:caso_uso;type:varchar(6);not null"`
}

func (PermisoModel) TableName() string { return "permiso" }

// RolPermisoModel es la representación GORM de la tabla rol_permiso (PK compuesta).
type RolPermisoModel struct {
	RolID     int16 `gorm:"column:rol_id;primaryKey"`
	PermisoID int16 `gorm:"column:permiso_id;primaryKey"`

	Rol     RolModel     `gorm:"foreignKey:RolID;references:ID;constraint:OnDelete:CASCADE"`
	Permiso PermisoModel `gorm:"foreignKey:PermisoID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (RolPermisoModel) TableName() string { return "rol_permiso" }

// UsuarioModel es la representación GORM de la tabla usuario.
type UsuarioModel struct {
	ID               string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RolID            int16  `gorm:"column:rol_id;not null;index:idx_usuario_rol"`
	Nombres          string `gorm:"type:varchar(80);not null"`
	Apellidos        string `gorm:"type:varchar(80);not null"`
	Correo           string `gorm:"type:varchar(120);not null"`
	Telefono         string `gorm:"type:varchar(20)"`
	HashContrasena   string `gorm:"column:hash_contrasena;type:varchar(100);not null"`
	IntentosFallidos int16  `gorm:"not null;default:0;check:chk_usuario_intentos,intentos_fallidos >= 0"`
	BloqueadoHasta   *time.Time
	VersionSesion    int32 `gorm:"not null;default:1"`
	Activo           bool  `gorm:"not null;default:true"`
	UltimoAcceso     *time.Time
	CreadoEn         time.Time `gorm:"not null;default:now()"`
	ActualizadoEn    time.Time `gorm:"not null;default:now()"`

	Rol RolModel `gorm:"foreignKey:RolID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (UsuarioModel) TableName() string { return "usuario" }

// TokenRecuperacionModel es la representación GORM de la tabla token_recuperacion.
type TokenRecuperacionModel struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UsuarioID string    `gorm:"column:usuario_id;type:uuid;not null"`
	Proposito string    `gorm:"type:varchar(12);not null;check:chk_token_proposito,proposito IN ('RECUPERACION','ACTIVACION')"`
	TokenHash string    `gorm:"column:token_hash;type:char(64);not null;unique"`
	ExpiraEn  time.Time `gorm:"not null"`
	UsadoEn   *time.Time
	CreadoEn  time.Time `gorm:"not null;default:now()"`

	Usuario UsuarioModel `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnDelete:CASCADE"`
}

func (TokenRecuperacionModel) TableName() string { return "token_recuperacion" }

// AuditoriaModel es la representación GORM de la tabla auditoria (solo anexión).
type AuditoriaModel struct {
	ID          int64     `gorm:"primaryKey;type:bigserial"`
	UsuarioID   *string   `gorm:"column:usuario_id;type:uuid;index:idx_auditoria_usuario_fecha,priority:1"`
	Accion      string    `gorm:"type:varchar(40);not null"`
	Entidad     string    `gorm:"type:varchar(40);not null;index:idx_auditoria_entidad,priority:1"`
	EntidadID   string    `gorm:"column:entidad_id;type:varchar(40);index:idx_auditoria_entidad,priority:2"`
	Detalle     *string   `gorm:"type:jsonb"`
	DireccionIP *string   `gorm:"column:direccion_ip;type:inet"`
	Fecha       time.Time `gorm:"not null;default:now();index:idx_auditoria_usuario_fecha,priority:2,sort:desc;index:idx_auditoria_fecha,sort:desc"`

	Usuario *UsuarioModel `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (AuditoriaModel) TableName() string { return "auditoria" }

// CategoriaModel es la representación GORM de la tabla categoria.
type CategoriaModel struct {
	ID              int16     `gorm:"primaryKey"`
	Nombre          string    `gorm:"type:varchar(60);not null;unique"`
	AniosAntiguedad int16     `gorm:"column:anios_antiguedad;not null;default:0;check:chk_categoria_anios,anios_antiguedad >= 0"`
	Activa          bool      `gorm:"not null;default:true"`
	CreadoEn        time.Time `gorm:"not null;default:now()"`
}

func (CategoriaModel) TableName() string { return "categoria" }

// CategoriaRequisitoModel es un documento probatorio exigido por una categoría.
type CategoriaRequisitoModel struct {
	ID          int32  `gorm:"primaryKey"`
	CategoriaID int16  `gorm:"column:categoria_id;not null;uniqueIndex:ux_categoria_requisito,priority:1"`
	Nombre      string `gorm:"type:varchar(120);not null;uniqueIndex:ux_categoria_requisito,priority:2"`

	Categoria CategoriaModel `gorm:"foreignKey:CategoriaID;references:ID;constraint:OnDelete:CASCADE"`
}

func (CategoriaRequisitoModel) TableName() string { return "categoria_requisito" }

// DocenteModel es la representación GORM de la tabla docente (baja lógica, cédula inmutable).
type DocenteModel struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Codigo        *string   `gorm:"type:varchar(20);unique"`
	Cedula        *string   `gorm:"type:varchar(20);unique"`
	Nombres       string    `gorm:"type:varchar(80);not null"`
	Apellidos     string    `gorm:"type:varchar(80);not null"`
	Correo        string    `gorm:"type:varchar(120)"`
	Telefono      string    `gorm:"type:varchar(20)"`
	CategoriaID   *int16    `gorm:"column:categoria_id"`
	FacultadID    *int16    `gorm:"column:facultad_id"`
	Activo        bool      `gorm:"not null;default:true"`
	CreadoEn      time.Time `gorm:"not null;default:now()"`
	ActualizadoEn time.Time `gorm:"not null;default:now()"`

	Categoria *CategoriaModel `gorm:"foreignKey:CategoriaID;references:ID;constraint:OnDelete:RESTRICT"`
	Facultad  *FacultadModel  `gorm:"foreignKey:FacultadID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (DocenteModel) TableName() string { return "docente" }

// TramiteModel es la representación GORM de la tabla tramite.
type TramiteModel struct {
	ID                  string     `gorm:"type:uuid;primaryKey"`
	Codigo              string     `gorm:"type:varchar(16);not null;unique"`
	DocenteID           string     `gorm:"column:docente_id;type:uuid;not null;index:idx_tramite_docente"`
	CategoriaID         *int16     `gorm:"column:categoria_id;check:chk_tramite_categoria,tipo = 'REGULARIZACION' OR categoria_id IS NOT NULL"` // nulo en regularizaciones
	Tipo                string     `gorm:"type:varchar(14);not null;check:chk_tramite_tipo,tipo IN ('CATEGORIZACION','REGULARIZACION')"`
	Estado              string     `gorm:"type:varchar(20);not null;index:idx_tramite_estado;check:chk_tramite_estado,estado IN ('INGRESADO','ASIGNADO','EN_REVISION','OBSERVADO','DICTAMINADO','EN_FACULTAD','REVISION_FINAL','REMITIDO','ARCHIVADO')"`
	Origen              string     `gorm:"type:varchar(6);not null;check:chk_tramite_origen,origen IN ('STAFF','PORTAL')"`
	PendienteValidacion bool       `gorm:"column:pendiente_validacion;not null;default:false"`
	AnalistaID          *string    `gorm:"column:analista_id;type:uuid;index:idx_tramite_analista"`
	RadicadoPor         *string    `gorm:"column:radicado_por;type:uuid"`
	Resultado           *string    `gorm:"type:varchar(12);check:chk_tramite_resultado,resultado IN ('FAVORABLE','DESFAVORABLE')"`
	Dictamen            *string    `gorm:"type:text"`
	DictaminadoEn       *time.Time `gorm:"column:dictaminado_en"`
	ArchivadoEn         *time.Time `gorm:"column:archivado_en"`
	UbicacionLegajo     *string    `gorm:"column:ubicacion_legajo;type:varchar(160)"`
	AnioDesde           *int16     `gorm:"column:anio_desde"`
	AnioHasta           *int16     `gorm:"column:anio_hasta"`
	ResolucionNumero    *string    `gorm:"column:resolucion_numero;type:varchar(40)"`
	ResolucionFecha     *time.Time `gorm:"column:resolucion_fecha;type:date"`
	RemisionNota        *string    `gorm:"column:remision_nota;type:varchar(120)"`
	CreadoEn            time.Time  `gorm:"not null;default:now()"`
	ActualizadoEn       time.Time  `gorm:"not null;default:now()"`

	Docente   DocenteModel   `gorm:"foreignKey:DocenteID;references:ID;constraint:OnDelete:RESTRICT"`
	Categoria *CategoriaModel `gorm:"foreignKey:CategoriaID;references:ID;constraint:OnDelete:RESTRICT"`
	Analista  *UsuarioModel   `gorm:"foreignKey:AnalistaID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (TramiteModel) TableName() string { return "tramite" }

// TramiteHistorialModel es la bitácora de estados del trámite (solo anexión).
type TramiteHistorialModel struct {
	ID             int64     `gorm:"primaryKey;type:bigserial"`
	TramiteID      string    `gorm:"column:tramite_id;type:uuid;not null;index:idx_historial_tramite"`
	Estado         string    `gorm:"type:varchar(20);not null"`
	EstadoAnterior *string   `gorm:"column:estado_anterior;type:varchar(20)"`
	UsuarioID      *string   `gorm:"column:usuario_id;type:uuid"`
	Detalle        *string   `gorm:"type:text"`
	Fecha          time.Time `gorm:"not null;default:now()"`

	Tramite TramiteModel  `gorm:"foreignKey:TramiteID;references:ID;constraint:OnDelete:CASCADE"`
	Usuario *UsuarioModel `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnDelete:RESTRICT"`
}

func (TramiteHistorialModel) TableName() string { return "tramite_historial" }

// TramiteDocumentoModel es un adjunto del trámite almacenado en S3. La baja es lógica (CU-16):
// activo=false lo retira del expediente y de las búsquedas sin borrar el objeto.
type TramiteDocumentoModel struct {
	ID             string     `gorm:"type:uuid;primaryKey"`
	TramiteID      *string    `gorm:"column:tramite_id;type:uuid;index:idx_tramite_documento_tramite"` // nulo en los escaneos del expediente del docente
	DocenteID      *string    `gorm:"column:docente_id;type:uuid;index:idx_tramite_documento_docente"`
	TipoDocumento  *string    `gorm:"column:tipo_documento;type:varchar(20)"`
	RequisitoID    *int32     `gorm:"column:requisito_id"`
	NombreOriginal string     `gorm:"column:nombre_original;type:varchar(200);not null"`
	TipoMime       string     `gorm:"column:tipo_mime;type:varchar(80);not null"`
	Tamano         int64      `gorm:"not null"`
	HashSHA256     *string    `gorm:"column:hash_sha256;type:char(64)"` // nulo hasta que el ia-service lo calcula (CU-15)
	Bucket         string     `gorm:"type:varchar(100);not null"`
	ClaveObjeto    string     `gorm:"column:clave_objeto;type:varchar(300);not null"`
	Activo         bool       `gorm:"not null;default:true"`
	MotivoBaja     *string    `gorm:"column:motivo_baja;type:varchar(255)"`
	BajaEn         *time.Time `gorm:"column:baja_en"`
	BajaPor        *string    `gorm:"column:baja_por;type:uuid"`
	CreadoEn       time.Time  `gorm:"not null;default:now()"`

	// CU-15: OCR y compresión.
	EstadoOCR        string  `gorm:"column:estado_ocr;type:varchar(12);not null;default:NO_APLICA"`
	TextoOCR         *string `gorm:"column:texto_ocr;type:text"`
	TextoValidado    bool    `gorm:"column:texto_validado;not null;default:false"`
	ErrorOCR         *string `gorm:"column:error_ocr;type:varchar(300)"`
	ClaveComprimida  *string `gorm:"column:clave_comprimida;type:varchar(300)"`
	TamanoComprimido *int64  `gorm:"column:tamano_comprimido"`

	Tramite   *TramiteModel            `gorm:"foreignKey:TramiteID;references:ID;constraint:OnDelete:CASCADE"`
	Docente   *DocenteModel            `gorm:"foreignKey:DocenteID;references:ID;constraint:OnDelete:CASCADE"`
	Requisito *CategoriaRequisitoModel `gorm:"foreignKey:RequisitoID;references:ID;constraint:OnDelete:SET NULL"`
	Baja      *UsuarioModel            `gorm:"foreignKey:BajaPor;references:ID;constraint:OnDelete:RESTRICT"`
}

func (TramiteDocumentoModel) TableName() string { return "tramite_documento" }

// DocumentoPaginaModel es el texto OCR de una página de un documento (CU-15); su índice de texto
// completo es manual (migración documento_paginas) y lo usa el buscador del visor.
type DocumentoPaginaModel struct {
	DocumentoID string `gorm:"column:documento_id;type:uuid;primaryKey"`
	Pagina      int    `gorm:"primaryKey;autoIncrement:false"`
	Texto       string `gorm:"type:text;not null;default:''"`

	Documento TramiteDocumentoModel `gorm:"foreignKey:DocumentoID;references:ID;constraint:OnDelete:CASCADE"`
}

func (DocumentoPaginaModel) TableName() string { return "documento_pagina" }

// ObservacionModel es una observación técnica registrada durante la revisión.
type ObservacionModel struct {
	ID         string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TramiteID  string     `gorm:"column:tramite_id;type:uuid;not null;index:idx_observacion_tramite"`
	Documento  string     `gorm:"type:varchar(120);not null"`
	Motivo     string     `gorm:"type:text;not null"`
	Plazo      time.Time  `gorm:"type:date;not null"`
	Resuelta   bool       `gorm:"not null;default:false"`
	ResueltaEn *time.Time `gorm:"column:resuelta_en"`
	CreadaPor  string     `gorm:"column:creada_por;type:uuid;not null"`
	CreadoEn   time.Time  `gorm:"not null;default:now()"`

	Tramite TramiteModel `gorm:"foreignKey:TramiteID;references:ID;constraint:OnDelete:CASCADE"`
	Autor   UsuarioModel `gorm:"foreignKey:CreadaPor;references:ID;constraint:OnDelete:RESTRICT"`
}

func (ObservacionModel) TableName() string { return "observacion" }

// FacultadModel es la unidad académica (facultad, carrera o decanato) del docente.
type FacultadModel struct {
	ID     int16  `gorm:"primaryKey"`
	Nombre string `gorm:"type:varchar(160);not null;unique"`
}

func (FacultadModel) TableName() string { return "facultad" }

// LegajoModel es una carpeta física del docente: gaveta, posición, anexos y observación.
type LegajoModel struct {
	ID          int32   `gorm:"primaryKey"`
	Codigo      string  `gorm:"type:varchar(20);not null;uniqueIndex:ux_legajo_codigo"` // lo asigna un trigger si va vacío
	Estado      string  `gorm:"type:varchar(10);not null;default:EN_ARCHIVO"`
	DocenteID   string  `gorm:"column:docente_id;type:uuid;not null;index:idx_legajo_docente"`
	Gaveta      *string `gorm:"type:varchar(10);index:idx_legajo_gaveta"`
	Orden       *int32
	Anexos      *string `gorm:"type:varchar(80)"`
	Observacion *string `gorm:"type:varchar(200)"`

	Docente DocenteModel `gorm:"foreignKey:DocenteID;references:ID;constraint:OnDelete:CASCADE"`
}

func (LegajoModel) TableName() string { return "legajo" }

// SolicitudCarpetaModel es el pedido de una carpeta física (CU-20, CU-21).
type SolicitudCarpetaModel struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LegajoID      int32      `gorm:"column:legajo_id;not null;index:idx_solicitud_legajo"`
	SolicitanteID string     `gorm:"column:solicitante_id;type:uuid;not null"`
	TramiteID     *string    `gorm:"column:tramite_id;type:uuid;index:idx_solicitud_tramite"`
	Motivo        string     `gorm:"type:text;not null"`
	Estado        string     `gorm:"type:varchar(10);not null;default:PENDIENTE;index:idx_solicitud_estado"`
	Respuesta     *string    `gorm:"type:text"`
	AtendidoPor   *string    `gorm:"column:atendido_por;type:uuid"`
	AtendidoEn    *time.Time `gorm:"column:atendido_en"`
	CreadoEn      time.Time  `gorm:"not null;default:now()"`

	Legajo      LegajoModel   `gorm:"foreignKey:LegajoID;references:ID;constraint:OnDelete:RESTRICT"`
	Solicitante UsuarioModel  `gorm:"foreignKey:SolicitanteID;references:ID;constraint:OnDelete:RESTRICT"`
	Tramite     *TramiteModel `gorm:"foreignKey:TramiteID;references:ID;constraint:OnDelete:SET NULL"`
}

func (SolicitudCarpetaModel) TableName() string { return "solicitud_carpeta" }

// MovimientoCustodiaModel es la bitácora de salidas y entradas de carpetas (solo anexión).
type MovimientoCustodiaModel struct {
	ID             int64      `gorm:"primaryKey;type:bigserial"`
	LegajoID       int32      `gorm:"column:legajo_id;not null;index:idx_movimiento_legajo"`
	SolicitudID    string     `gorm:"column:solicitud_id;type:uuid;not null"`
	Tipo           string     `gorm:"type:varchar(7);not null"`
	CustodioID     string     `gorm:"column:custodio_id;type:uuid;not null"`
	FechaLimite    *time.Time `gorm:"column:fecha_limite;type:date"`
	EstadoMaterial string     `gorm:"column:estado_material;type:varchar(11);not null"`
	Observacion    *string    `gorm:"type:text"`
	RegistradoPor  string     `gorm:"column:registrado_por;type:uuid;not null"`
	Fecha          time.Time  `gorm:"not null;default:now();index:idx_movimiento_fecha,sort:desc"`

	Legajo    LegajoModel           `gorm:"foreignKey:LegajoID;references:ID;constraint:OnDelete:RESTRICT"`
	Solicitud SolicitudCarpetaModel `gorm:"foreignKey:SolicitudID;references:ID;constraint:OnDelete:RESTRICT"`
	Custodio  UsuarioModel          `gorm:"foreignKey:CustodioID;references:ID;constraint:OnDelete:RESTRICT"`
	Registro  UsuarioModel          `gorm:"foreignKey:RegistradoPor;references:ID;constraint:OnDelete:RESTRICT"`
}

func (MovimientoCustodiaModel) TableName() string { return "movimiento_custodia" }

// NotificacionEnvioModel es un aviso de cambio de estado enviado (o por enviar) al docente (CU-14).
type NotificacionEnvioModel struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TramiteID     string     `gorm:"column:tramite_id;type:uuid;not null;index:idx_notificacion_tramite"`
	TramiteCodigo string     `gorm:"column:tramite_codigo;type:varchar(30);not null"`
	EstadoTramite string     `gorm:"column:estado_tramite;type:varchar(20);not null"`
	Destino       string     `gorm:"type:varchar(120);not null;default:''"`
	Nombres       string     `gorm:"type:varchar(80);not null;default:''"`
	Estado        string     `gorm:"type:varchar(9);not null;default:PENDIENTE;index:idx_notificacion_estado,priority:1"`
	Intentos      int16      `gorm:"not null;default:0"`
	Error         *string    `gorm:"type:text"`
	CreadoEn      time.Time  `gorm:"not null;default:now();index:idx_notificacion_estado,priority:2"`
	EnviadoEn     *time.Time `gorm:"column:enviado_en"`

	Tramite TramiteModel `gorm:"foreignKey:TramiteID;references:ID;constraint:OnDelete:CASCADE"`
}

func (NotificacionEnvioModel) TableName() string { return "notificacion_envio" }
