package rol

import (
	"context"
	"time"
)

// Permisos protegidos: el rol Director de sistema nunca puede perderlos (sección 4.4).
const (
	PermisoPerfilGestionar    = "PERFIL_GESTIONAR"
	PermisoUsuarioGestionar   = "USUARIO_GESTIONAR"
	PermisoRolGestionar       = "ROL_GESTIONAR"
	PermisoAuditoriaConsultar = "AUDITORIA_CONSULTAR"

	PermisoDocenteGestionar   = "DOCENTE_GESTIONAR"
	PermisoCategoriaGestionar = "CATEGORIA_GESTIONAR"
	PermisoTramiteRadicar     = "TRAMITE_RADICAR"
	PermisoTramiteAsignar     = "TRAMITE_ASIGNAR"
	PermisoTramiteRevisar     = "TRAMITE_REVISAR"
	PermisoTramiteArchivar    = "TRAMITE_ARCHIVAR"
	PermisoDocumentoGestionar = "DOCUMENTO_GESTIONAR"
	PermisoDocumentoRestituir = "DOCUMENTO_RESTITUIR"
	PermisoDocumentoBuscar    = "DOCUMENTO_BUSCAR"
	PermisoDocumentoEscanear  = "DOCUMENTO_ESCANEAR" // subir los escaneos del expediente digital del docente
	PermisoLegajoGestionar    = "LEGAJO_GESTIONAR"
	PermisoCarpetaSolicitar   = "CARPETA_SOLICITAR"
	PermisoCustodiaRegistrar  = "CUSTODIA_REGISTRAR"
	PermisoBitacoraConsultar  = "BITACORA_CONSULTAR"
	PermisoNotificacionVer    = "NOTIFICACION_CONSULTAR"
	PermisoTableroVer         = "TABLERO_VER"     // CU-25
	PermisoReporteGenerar     = "REPORTE_GENERAR" // CU-26
)

var PermisosProtegidosDirector = []string{PermisoUsuarioGestionar, PermisoRolGestionar}

// NombreDirector es el nombre del rol de sistema con protección adicional.
const NombreDirector = "Director"

type Permiso struct {
	ID      int16
	Codigo  string
	Nombre  string
	Modulo  string
	CasoUso string
}

type Rol struct {
	ID          int16
	Nombre      string
	Descripcion string
	EsSistema   bool
	Activo      bool
	CreadoEn    time.Time
	Permisos    []Permiso
}

// TienePermiso indica si el rol incluye el código de permiso dado.
func (r *Rol) TienePermiso(codigo string) bool {
	for _, p := range r.Permisos {
		if p.Codigo == codigo {
			return true
		}
	}
	return false
}

type ConConteo struct {
	Rol           Rol
	UsuariosCount int64
}

// Repository es el puerto de persistencia de roles y su catálogo de permisos.
type Repository interface {
	Listar(ctx context.Context) ([]ConConteo, error)
	BuscarPorID(ctx context.Context, id int16) (*Rol, error)
	BuscarPorNombre(ctx context.Context, nombre string) (*Rol, error)
	// Guardar crea el rol si ID es 0, o lo actualiza junto con sus permisos.
	Guardar(ctx context.Context, r *Rol, permisoIDs []int16) error
	Eliminar(ctx context.Context, id int16) error
	ContarUsuarios(ctx context.Context, id int16) (int64, error)
	IDsUsuariosDelRol(ctx context.Context, id int16) ([]string, error)
	ListarPermisos(ctx context.Context) ([]Permiso, error)
	PermisosExisten(ctx context.Context, ids []int16) (bool, error)
	PermisosPorRol(ctx context.Context, rolID int16) ([]string, error)
}
