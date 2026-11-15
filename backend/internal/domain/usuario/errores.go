package usuario

import "time"

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrNoEncontrado               erroDominio = "usuario no encontrado"
	ErrCredencialesInvalidas      erroDominio = "correo o contraseña incorrectos"
	ErrCorreoDuplicado            erroDominio = "el correo ya está registrado"
	ErrRolInvalido                erroDominio = "el rol no existe o está inactivo"
	ErrAutocambioRol              erroDominio = "no puedes cambiar tu propio rol"
	ErrAutodesactivacion          erroDominio = "no puedes desactivar tu propia cuenta"
	ErrCuentaYaActivada           erroDominio = "la cuenta ya fue activada"
	ErrPoliticaContrasena         erroDominio = "la contraseña no cumple la política mínima"
	ErrContrasenaActualIncorrecta erroDominio = "la contraseña actual es incorrecta"
	ErrCampoNoEditable            erroDominio = "ese campo no se puede editar desde aquí"
)

// ErrCuentaBloqueada se devuelve cuando el bloqueo por intentos fallidos sigue vigente.
type ErrCuentaBloqueada struct {
	Hasta time.Time
}

func (e *ErrCuentaBloqueada) Error() string { return "cuenta bloqueada temporalmente" }
