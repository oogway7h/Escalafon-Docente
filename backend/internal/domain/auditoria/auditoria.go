package auditoria

import (
	"context"
	"time"
)

// Acciones auditadas por los casos de uso de HU-01 y HU-02.
const (
	AccionLogin              = "LOGIN"
	AccionLoginFallido       = "LOGIN_FALLIDO"
	AccionRestablecer        = "RESTABLECER_CONTRASENA"
	AccionActivar            = "ACTIVAR_CUENTA"
	AccionCambiarContrasena  = "CAMBIAR_CONTRASENA"
	AccionActualizarContacto = "ACTUALIZAR_CONTACTO"
	AccionCrear              = "CREAR"
	AccionEditar             = "EDITAR"
	AccionCambioRol          = "CAMBIO_ROL"
	AccionBaja               = "BAJA"
	AccionReactivar          = "REACTIVAR"
	AccionReenviarActivacion = "REENVIAR_ACTIVACION"
	AccionGuardarRol         = "GUARDAR_ROL"
	AccionEliminarRol        = "ELIMINAR_ROL"

	AccionEliminarCategoria   = "ELIMINAR_CATEGORIA"
	AccionRadicar             = "RADICAR_TRAMITE"
	AccionRadicarPortal       = "RADICAR_TRAMITE_PORTAL"
	AccionValidarSolicitud    = "VALIDAR_SOLICITUD"
	AccionAsignarTramite      = "ASIGNAR_TRAMITE"
	AccionReasignarTramite    = "REASIGNAR_TRAMITE"
	AccionIniciarRevision     = "INICIAR_REVISION"
	AccionObservar            = "OBSERVAR_TRAMITE"
	AccionResolverObservacion = "RESOLVER_OBSERVACION"
	AccionDictaminar          = "DICTAMINAR_TRAMITE"
	AccionArchivar            = "ARCHIVAR_TRAMITE"
	AccionBajaDocumento       = "BAJA_DOCUMENTO"
	AccionRestituirDocumento  = "RESTITUIR_DOCUMENTO"
	AccionRegistrarPeriodo    = "REGISTRAR_PERIODO"
	AccionRemitirFacultad     = "REMITIR_FACULTAD"
	AccionRegresoFacultad     = "REGRESO_FACULTAD"
	AccionRemitirRectorado    = "REMITIR_RECTORADO"
	AccionCargarDocumento     = "CARGAR_DOCUMENTO"
	AccionCorregirTextoOCR    = "CORREGIR_TEXTO_OCR"
	AccionSolicitarCarpeta    = "SOLICITAR_CARPETA"
	AccionAtenderSolicitud    = "ATENDER_SOLICITUD"
	AccionRechazarSolicitud   = "RECHAZAR_SOLICITUD"
	AccionDespacharCarpeta    = "DESPACHAR_CARPETA"
	AccionDevolverCarpeta     = "DEVOLVER_CARPETA"
	AccionGenerarReporte      = "GENERAR_REPORTE"
	AccionAceptarTerminos     = "ACEPTAR_TERMINOS"

	EntidadUsuario   = "usuario"
	EntidadRol       = "rol"
	EntidadDocente   = "docente"
	EntidadCategoria = "categoria"
	EntidadTramite   = "tramite"
	EntidadLegajo    = "legajo"
	EntidadSolicitud = "solicitud_carpeta"
	EntidadReporte   = "reporte"
)

// Acciones es el catálogo de tipos de acción que se pueden consultar (CU-06).
var Acciones = []string{
	AccionLogin, AccionLoginFallido, AccionRestablecer, AccionActivar, AccionCambiarContrasena,
	AccionActualizarContacto, AccionCrear, AccionEditar, AccionCambioRol, AccionBaja, AccionReactivar,
	AccionReenviarActivacion, AccionGuardarRol, AccionEliminarRol,
	AccionEliminarCategoria, AccionRadicar, AccionRadicarPortal, AccionValidarSolicitud, AccionAsignarTramite,
	AccionReasignarTramite, AccionIniciarRevision, AccionObservar, AccionResolverObservacion, AccionDictaminar,
	AccionArchivar, AccionBajaDocumento, AccionRestituirDocumento, AccionSolicitarCarpeta, AccionAtenderSolicitud, AccionRechazarSolicitud,
	AccionDespacharCarpeta, AccionDevolverCarpeta, AccionCargarDocumento, AccionCorregirTextoOCR,
	AccionRegistrarPeriodo, AccionRemitirFacultad, AccionRegresoFacultad, AccionRemitirRectorado,
	AccionGenerarReporte, AccionAceptarTerminos,
}

// AccionValida indica si la acción pertenece al catálogo.
func AccionValida(accion string) bool {
	for _, a := range Acciones {
		if a == accion {
			return true
		}
	}
	return false
}

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrAccionInvalida erroDominio = "el tipo de acción no existe"
	ErrRangoFechas    erroDominio = "la fecha inicial es posterior a la final"
)

// Evento es un registro de solo anexión: la tabla auditoria rechaza UPDATE/DELETE/TRUNCATE.
type Evento struct {
	ID          int64
	UsuarioID   *string
	Accion      string
	Entidad     string
	EntidadID   string
	Detalle     map[string]any
	DireccionIP string
	Fecha       time.Time
}

// Repository es el puerto de persistencia de auditoría (solo anexión).
type Repository interface {
	Registrar(ctx context.Context, e *Evento) error
}

// Autor es el usuario que realizó la acción (datos actuales de su cuenta).
type Autor struct {
	ID        string
	Nombres   string
	Apellidos string
	Correo    string
}

// Registro es un evento con su autor; Autor es nil si la acción no tuvo usuario (login con un
// correo inexistente o solicitud del portal público).
type Registro struct {
	Evento
	Autor *Autor
}

// Filtro acota la consulta de la bitácora. Los campos vacíos no filtran. Hasta es exclusivo.
type Filtro struct {
	UsuarioID string
	Accion    string
	Desde     *time.Time
	Hasta     *time.Time
	Pagina    int
	Tamano    int
}

// Consulta es el puerto de lectura de la bitácora (CU-06): no expone escritura.
type Consulta interface {
	Listar(ctx context.Context, f Filtro) ([]Registro, int64, error)
	// ListarAutores devuelve los usuarios que figuran en la bitácora, para el filtro por usuario.
	ListarAutores(ctx context.Context) ([]Autor, error)
}
