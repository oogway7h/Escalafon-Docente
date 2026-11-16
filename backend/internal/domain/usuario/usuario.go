package usuario

import (
	"context"
	"regexp"
	"time"
	"unicode"
)

// Usuario es el agregado central de seguridad: credenciales, estado de sesión y bloqueo.
type Usuario struct {
	ID               string
	RolID            int16
	Nombres          string
	Apellidos        string
	Correo           string
	Telefono         string
	HashContrasena   string
	IntentosFallidos int16
	BloqueadoHasta   *time.Time
	VersionSesion    int32
	Activo           bool
	UltimoAcceso     *time.Time
	CreadoEn         time.Time
	ActualizadoEn    time.Time
	// Acuerdo de usuario: versión aceptada (0 = nunca) y cuándo. Subir VersionTerminos lo vuelve a pedir a todos.
	TerminosVersion     int16
	TerminosAceptadosEn *time.Time
}

// VersionTerminos es la versión vigente del acuerdo de usuario (el texto vive en el frontend: lib/terminos.ts).
const VersionTerminos int16 = 1

// AceptoTerminos indica si el usuario aceptó la versión vigente del acuerdo.
func (u *Usuario) AceptoTerminos() bool { return u.TerminosVersion >= VersionTerminos }

// EstaBloqueado indica si, en el instante ahora, el usuario tiene un bloqueo vigente.
func (u *Usuario) EstaBloqueado(ahora time.Time) bool {
	return u.BloqueadoHasta != nil && u.BloqueadoHasta.After(ahora)
}

var correoRe = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

// CorreoValido valida el formato básico de un correo electrónico.
func CorreoValido(correo string) bool {
	return correoRe.MatchString(correo)
}

// PoliticaContrasenaValida exige mínimo 8 caracteres, una mayúscula y un dígito.
func PoliticaContrasenaValida(contrasena string) bool {
	if len(contrasena) < 8 {
		return false
	}
	var tieneMayuscula, tieneDigito bool
	for _, r := range contrasena {
		switch {
		case unicode.IsUpper(r):
			tieneMayuscula = true
		case unicode.IsDigit(r):
			tieneDigito = true
		}
	}
	return tieneMayuscula && tieneDigito
}

// ListarFiltro acota la búsqueda paginada de usuarios (HU-02).
type ListarFiltro struct {
	Buscar string
	RolID  *int16
	Activo *bool
	Pagina int
	Tamano int
}

// Repository es el puerto de persistencia del agregado Usuario.
type Repository interface {
	BuscarPorID(ctx context.Context, id string) (*Usuario, error)
	BuscarPorCorreo(ctx context.Context, correo string) (*Usuario, error)
	Listar(ctx context.Context, filtro ListarFiltro) ([]Usuario, int64, error)
	Crear(ctx context.Context, u *Usuario) error
	Actualizar(ctx context.Context, u *Usuario) error
	// IncrementarVersionSesionPorRol invalida de un golpe el JWT de todos los usuarios de un rol
	// (cambian sus permisos, sección 8 HU-02 "GuardarRol").
	IncrementarVersionSesionPorRol(ctx context.Context, rolID int16) error
}
