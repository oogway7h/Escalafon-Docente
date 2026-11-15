package rol

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrNoEncontrado     erroDominio = "rol no encontrado"
	ErrDuplicado        erroDominio = "ya existe un rol con ese nombre"
	ErrPermisoInvalido  erroDominio = "uno o más permisos no existen"
	ErrSistemaProtegido erroDominio = "el rol de sistema no se puede modificar de esa forma"
)

// ErrConUsuarios se devuelve al intentar eliminar un rol que aún tiene usuarios asignados.
type ErrConUsuarios struct {
	Cantidad int64
}

func (e *ErrConUsuarios) Error() string { return "el rol tiene usuarios asignados" }
