// Package apierror traduce errores de dominio a códigos HTTP (sección 7): es la única capa que
// conoce ambos lenguajes, así los casos de uso no dependen de HTTP y los handlers no adivinan
// códigos de estado.
package apierror

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/categoria"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/domain/usuario"
)

type ErrorAPI struct {
	Status   int
	Codigo   string
	Mensaje  string
	Detalles any
}

func (e *ErrorAPI) Error() string { return e.Mensaje }

func Nuevo(status int, codigo, mensaje string) *ErrorAPI {
	return &ErrorAPI{Status: status, Codigo: codigo, Mensaje: mensaje}
}

// Validacion construye el 422 uniforme para errores de validación de entrada.
func Validacion(campos map[string]string) *ErrorAPI {
	return &ErrorAPI{Status: 422, Codigo: "VALIDACION", Mensaje: "La solicitud tiene datos inválidos", Detalles: campos}
}

// Mapear traduce un error de dominio al código HTTP y código de negocio de la sección 7.
// Un error no reconocido se trata como 500 ERROR_INTERNO (sin filtrar su mensaje interno).
func Mapear(err error) *ErrorAPI {
	if err == nil {
		return nil
	}
	var apiErr *ErrorAPI
	if errors.As(err, &apiErr) {
		return apiErr
	}

	// Errores propios de Fiber (ruta inexistente, método no soportado, body malformado):
	// conservan el código HTTP que trae Fiber en vez de caer a 500.
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		codigo := "ERROR_INTERNO"
		if fiberErr.Code == fiber.StatusNotFound {
			codigo = "NO_ENCONTRADO"
		}
		return &ErrorAPI{Status: fiberErr.Code, Codigo: codigo, Mensaje: fiberErr.Message}
	}

	var bloqueo *usuario.ErrCuentaBloqueada
	if errors.As(err, &bloqueo) {
		return &ErrorAPI{Status: 423, Codigo: "CUENTA_BLOQUEADA", Mensaje: "La cuenta está bloqueada temporalmente",
			Detalles: map[string]any{"bloqueadoHasta": bloqueo.Hasta}}
	}
	var conUsuarios *rol.ErrConUsuarios
	if errors.As(err, &conUsuarios) {
		return &ErrorAPI{Status: 409, Codigo: "ROL_CON_USUARIOS", Mensaje: "El rol tiene usuarios asignados",
			Detalles: map[string]any{"cantidad": conUsuarios.Cantidad}}
	}

	var enUso *categoria.ErrEnUso
	if errors.As(err, &enUso) {
		return &ErrorAPI{Status: 409, Codigo: "CATEGORIA_EN_USO", Mensaje: "La categoría está en uso: dala de baja en vez de eliminarla",
			Detalles: map[string]any{"docentes": enUso.Docentes, "tramites": enUso.Tramites}}
	}
	var faltantes *tramite.ErrDocumentosFaltantes
	if errors.As(err, &faltantes) {
		return &ErrorAPI{Status: 422, Codigo: "DOCUMENTOS_FALTANTES", Mensaje: "Faltan documentos exigidos por la categoría",
			Detalles: map[string]any{"faltan": faltantes.Faltan}}
	}

	switch {
	case errors.Is(err, usuario.ErrCredencialesInvalidas):
		return Nuevo(401, "CREDENCIALES_INVALIDAS", "Correo o contraseña incorrectos")
	case errors.Is(err, usuario.ErrNoEncontrado):
		return Nuevo(404, "NO_ENCONTRADO", "El usuario no existe")
	case errors.Is(err, usuario.ErrCorreoDuplicado):
		return Nuevo(409, "CORREO_DUPLICADO", "El correo ya está registrado")
	case errors.Is(err, usuario.ErrRolInvalido):
		return Nuevo(422, "ROL_INVALIDO", "El rol no existe o está inactivo")
	case errors.Is(err, usuario.ErrAutocambioRol):
		return Nuevo(409, "AUTOCAMBIO_ROL", "No puedes cambiar tu propio rol")
	case errors.Is(err, usuario.ErrAutodesactivacion):
		return Nuevo(409, "AUTODESACTIVACION", "No puedes desactivar tu propia cuenta")
	case errors.Is(err, usuario.ErrCuentaYaActivada):
		return Nuevo(409, "CUENTA_YA_ACTIVADA", "La cuenta ya fue activada")
	case errors.Is(err, usuario.ErrPoliticaContrasena):
		return Nuevo(422, "POLITICA_CONTRASENA", "La contraseña no cumple la política mínima")
	case errors.Is(err, usuario.ErrContrasenaActualIncorrecta):
		return Nuevo(422, "CONTRASENA_ACTUAL_INCORRECTA", "La contraseña actual es incorrecta")
	case errors.Is(err, usuario.ErrCampoNoEditable):
		return Nuevo(403, "CAMPO_NO_EDITABLE", "Ese campo no se puede editar desde aquí")

	case errors.Is(err, auditoria.ErrAccionInvalida):
		return Nuevo(422, "ACCION_INVALIDA", "El tipo de acción no existe")
	case errors.Is(err, auditoria.ErrRangoFechas), errors.Is(err, indicador.ErrRangoFechas):
		return Nuevo(422, "RANGO_FECHAS_INVALIDO", "La fecha inicial no puede ser posterior a la final")

	case errors.Is(err, indicador.ErrTipoReporte):
		return Nuevo(422, "TIPO_REPORTE_INVALIDO", "El tipo de reporte no existe")
	case errors.Is(err, indicador.ErrFormato):
		return Nuevo(422, "FORMATO_INVALIDO", "El formato debe ser PDF o XLSX")
	case errors.Is(err, indicador.ErrFiltroInvalido):
		return Nuevo(422, "FILTRO_INVALIDO", "Uno de los filtros no tiene un valor válido")
	case errors.Is(err, indicador.ErrInterpreteNoDisponible):
		return Nuevo(503, "INTERPRETE_NO_DISPONIBLE", "El asistente de reportes no está disponible")
	case errors.Is(err, indicador.ErrNoEntendido):
		return Nuevo(422, "PETICION_NO_ENTENDIDA", "No se pudo entender qué reporte se pide")
	case errors.Is(err, indicador.ErrReporteVacio):
		return Nuevo(422, "REPORTE_VACIO", "No hay información para los filtros elegidos")

	case errors.Is(err, token.ErrInvalido):
		return Nuevo(400, "TOKEN_INVALIDO", "El enlace es inválido o expiró")

	case errors.Is(err, rol.ErrNoEncontrado):
		return Nuevo(404, "NO_ENCONTRADO", "El rol no existe")
	case errors.Is(err, rol.ErrDuplicado):
		return Nuevo(409, "ROL_DUPLICADO", "Ya existe un rol con ese nombre")
	case errors.Is(err, rol.ErrPermisoInvalido):
		return Nuevo(422, "PERMISO_INVALIDO", "Uno o más permisos no existen")
	case errors.Is(err, rol.ErrSistemaProtegido):
		return Nuevo(409, "ROL_SISTEMA_PROTEGIDO", "El rol de sistema no se puede modificar de esa forma")

	case errors.Is(err, docente.ErrNoEncontrado):
		return Nuevo(404, "NO_ENCONTRADO", "El docente no existe")
	case errors.Is(err, docente.ErrCedulaDuplicada):
		return Nuevo(409, "CEDULA_DUPLICADA", "Ya existe un docente con esa cédula")
	case errors.Is(err, docente.ErrCodigoDuplicado):
		return Nuevo(409, "CODIGO_DUPLICADO", "Ya existe un docente con ese código")
	case errors.Is(err, docente.ErrFacultadInvalida):
		return Nuevo(422, "FACULTAD_INVALIDA", "La facultad no existe")
	case errors.Is(err, docente.ErrCategoriaInvalida), errors.Is(err, tramite.ErrCategoriaInvalida):
		return Nuevo(422, "CATEGORIA_INVALIDA", "La categoría no existe o está inactiva")

	case errors.Is(err, categoria.ErrNoEncontrado):
		return Nuevo(404, "NO_ENCONTRADO", "La categoría no existe")
	case errors.Is(err, categoria.ErrDuplicada):
		return Nuevo(409, "CATEGORIA_DUPLICADA", "Ya existe una categoría con ese nombre")

	case errors.Is(err, tramite.ErrNoEncontrado):
		return Nuevo(404, "NO_ENCONTRADO", "El trámite no existe")
	case errors.Is(err, tramite.ErrObservacionNoEncontrada):
		return Nuevo(404, "NO_ENCONTRADO", "La observación no existe")
	case errors.Is(err, tramite.ErrDocenteNoEncontrado):
		return Nuevo(404, "DOCENTE_NO_ENCONTRADO", "No existe ese docente en el padrón")
	case errors.Is(err, tramite.ErrDocenteNoVerificado):
		return Nuevo(422, "DOCENTE_NO_VERIFICADO", "El código y la cédula no corresponden a un docente activo del padrón")
	case errors.Is(err, tramite.ErrAdjuntosLibres):
		return Nuevo(422, "ADJUNTOS_REQUERIDOS", "Adjunta entre 1 y 5 documentos")
	case errors.Is(err, tramite.ErrDemasiadasPendientes):
		return Nuevo(409, "SOLICITUDES_PENDIENTES", "Ya hay solicitudes pendientes de validación en ventanilla")
	case errors.Is(err, tramite.ErrDocenteInactivo):
		return Nuevo(409, "DOCENTE_INACTIVO", "El docente está dado de baja")
	case errors.Is(err, tramite.ErrTransicionInvalida):
		return Nuevo(409, "TRANSICION_INVALIDA", "El trámite no puede pasar a ese estado desde el actual")
	case errors.Is(err, tramite.ErrArchivado):
		return Nuevo(409, "TRAMITE_ARCHIVADO", "El trámite está archivado y es de solo lectura")
	case errors.Is(err, tramite.ErrRemitido):
		return Nuevo(409, "TRAMITE_REMITIDO", "El trámite fue remitido al Rectorado y es de solo lectura")
	case errors.Is(err, tramite.ErrCarpetaNoSolicitada):
		return Nuevo(409, "CARPETA_NO_SOLICITADA", "Antes de revisar hay que solicitar la carpeta física del trámite")
	case errors.Is(err, tramite.ErrCarpetasAbiertas):
		return Nuevo(409, "CARPETAS_ABIERTAS", "Hay carpetas del trámite sin devolver al archivo")
	case errors.Is(err, tramite.ErrPeriodoRequerido):
		return Nuevo(409, "PERIODO_REQUERIDO", "Registra los años a evaluar antes de dictaminar")
	case errors.Is(err, tramite.ErrPeriodoInvalido):
		return Nuevo(422, "PERIODO_INVALIDO", "Los años a evaluar no son válidos")
	case errors.Is(err, tramite.ErrResolucionRequerida):
		return Nuevo(422, "RESOLUCION_REQUERIDA", "Indica el número y la fecha (no futura) de la resolución")
	case errors.Is(err, tramite.ErrNotaRequerida):
		return Nuevo(422, "NOTA_REQUERIDA", "Indica la nota de remisión al Rectorado")
	case errors.Is(err, legajo.ErrTramiteRequerido):
		return Nuevo(422, "TRAMITE_REQUERIDO", "Indica el trámite para el que se solicita la carpeta")
	case errors.Is(err, legajo.ErrTramiteInvalido):
		return Nuevo(422, "TRAMITE_INVALIDO", "El trámite no existe, no está a tu cargo, es de otro docente o ya está cerrado")
	case errors.Is(err, tramite.ErrMotivoRequerido):
		return Nuevo(422, "MOTIVO_REQUERIDO", "La reasignación exige un motivo")
	case errors.Is(err, tramite.ErrNoAsignado):
		return Nuevo(409, "NO_ASIGNADO", "El trámite no está asignado a tu usuario")
	case errors.Is(err, tramite.ErrAnalistaInvalido):
		return Nuevo(422, "ANALISTA_INVALIDO", "El usuario no es un analista activo")
	case errors.Is(err, tramite.ErrObservacionesPendientes):
		return Nuevo(409, "OBSERVACIONES_PENDIENTES", "Hay observaciones pendientes por resolver")
	case errors.Is(err, tramite.ErrPendienteValidacion):
		return Nuevo(409, "PENDIENTE_VALIDACION", "La solicitud del portal aún no fue validada en ventanilla")
	case errors.Is(err, tramite.ErrYaValidado):
		return Nuevo(409, "YA_VALIDADO", "La solicitud ya fue validada")
	case errors.Is(err, tramite.ErrArchivoInvalido):
		return Nuevo(422, "ARCHIVO_INVALIDO", "Los adjuntos deben ser PDF, JPG o PNG de hasta 8 MB")
	case errors.Is(err, tramite.ErrDocumentoNoEncontrado):
		return Nuevo(404, "DOCUMENTO_NO_ENCONTRADO", "El documento no existe")
	case errors.Is(err, tramite.ErrProcesadorNoDisponible):
		return Nuevo(503, "OCR_NO_DISPONIBLE", "El servicio de OCR no está disponible")
	case errors.Is(err, tramite.ErrTipoNoPermitido):
		return Nuevo(422, "TIPO_NO_PERMITIDO", "Solo se admiten archivos PDF, JPG o PNG")
	case errors.Is(err, tramite.ErrArchivoGrande):
		return Nuevo(422, "ARCHIVO_GRANDE", "El archivo supera el tamaño máximo (PDF 400 MB, imágenes 20 MB)")
	case errors.Is(err, tramite.ErrSinOCR):
		return Nuevo(409, "SIN_OCR", "El documento aún no tiene texto extraído para buscar")
	case errors.Is(err, tramite.ErrTipoDocumentoInvalido):
		return Nuevo(422, "TIPO_DOCUMENTO_INVALIDO", "El tipo de documento no es válido")
	case errors.Is(err, tramite.ErrConsultaInvalida):
		return Nuevo(422, "CONSULTA_INVALIDA", "La búsqueda debe tener entre 2 y 100 caracteres")
	case errors.Is(err, tramite.ErrPaginasInvalidas):
		return Nuevo(422, "PAGINAS_INVALIDAS", "El texto corregido debe conservar las páginas del documento")
	case errors.Is(err, tramite.ErrSubidaIncompleta):
		return Nuevo(409, "SUBIDA_INCOMPLETA", "El archivo todavía no llegó al almacén")
	case errors.Is(err, tramite.ErrDocumentoNoSubido):
		return Nuevo(409, "DOCUMENTO_ESTADO_INVALIDO", "El documento no está en un estado que permita esta operación")
	case errors.Is(err, tramite.ErrDocumentoDeBaja):
		return Nuevo(409, "DOCUMENTO_DE_BAJA", "El documento ya está dado de baja")
	case errors.Is(err, tramite.ErrMotivoBajaRequerido):
		return Nuevo(422, "MOTIVO_BAJA_REQUERIDO", "Indica el motivo de la baja del documento")
	case errors.Is(err, tramite.ErrDocumentoActivo):
		return Nuevo(409, "DOCUMENTO_ACTIVO", "El documento no está dado de baja")
	case errors.Is(err, tramite.ErrMotivoRestitucion):
		return Nuevo(422, "MOTIVO_RESTITUCION_REQUERIDO", "Indica el motivo de la restitución")
	case errors.Is(err, tramite.ErrPlazoInvalido):
		return Nuevo(422, "PLAZO_INVALIDO", "El plazo no puede estar en el pasado")
	case errors.Is(err, tramite.ErrAlmacenNoDisponible):
		return Nuevo(503, "ALMACEN_NO_DISPONIBLE", "El almacenamiento de documentos no está disponible")

	case errors.Is(err, legajo.ErrNoEncontrado):
		return Nuevo(404, "NO_ENCONTRADO", "El legajo no existe")
	case errors.Is(err, legajo.ErrSolicitudNoExiste):
		return Nuevo(404, "NO_ENCONTRADO", "La solicitud no existe")
	case errors.Is(err, legajo.ErrCodigoDuplicado):
		return Nuevo(409, "LEGAJO_CODIGO_DUPLICADO", "Ya existe un legajo con ese código")
	case errors.Is(err, legajo.ErrDocenteInvalido):
		return Nuevo(422, "DOCENTE_INVALIDO", "El docente no existe")
	case errors.Is(err, legajo.ErrNoDisponible):
		return Nuevo(409, "LEGAJO_NO_DISPONIBLE", "El legajo no está disponible en archivo")
	case errors.Is(err, legajo.ErrNoPrestado):
		return Nuevo(409, "LEGAJO_NO_PRESTADO", "El legajo no está prestado")
	case errors.Is(err, legajo.ErrSolicitudAbierta):
		return Nuevo(409, "SOLICITUD_ABIERTA", "Ya hay una solicitud en curso para ese legajo")
	case errors.Is(err, legajo.ErrEstadoSolicitud):
		return Nuevo(409, "SOLICITUD_ESTADO_INVALIDO", "La solicitud no está en el estado requerido")
	case errors.Is(err, legajo.ErrMotivoRequerido):
		return Nuevo(422, "MOTIVO_REQUERIDO", "El motivo es obligatorio")
	case errors.Is(err, legajo.ErrObservacionRequerida):
		return Nuevo(422, "OBSERVACION_REQUERIDA", "Describe la novedad: la carpeta no está en buen estado")
	case errors.Is(err, legajo.ErrFechaLimiteInvalid):
		return Nuevo(422, "FECHA_LIMITE_INVALIDA", "La fecha límite debe ser posterior a hoy")
	case errors.Is(err, legajo.ErrEstadoMaterial):
		return Nuevo(422, "ESTADO_MATERIAL_INVALIDO", "El estado material debe ser BUENO, REGULAR o DETERIORADO")

	default:
		return &ErrorAPI{Status: 500, Codigo: "ERROR_INTERNO", Mensaje: "Ocurrió un error inesperado"}
	}
}
