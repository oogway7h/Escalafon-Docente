package dto

type LoginEntrada struct {
	Correo     string `json:"correo" validate:"required,email"`
	Contrasena string `json:"contrasena" validate:"required"`
}

type LoginSalida struct {
	Token             string         `json:"token"`
	ExpiraEn          int            `json:"expiraEn"`
	Usuario           UsuarioResumen `json:"usuario"`
	Permisos          []string       `json:"permisos"`
	TerminosAceptados bool           `json:"terminosAceptados"`
}

type UsuarioResumen struct {
	ID        string `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Correo    string `json:"correo"`
	RolID     int16  `json:"rolId"`
}

type RecuperacionEntrada struct {
	Correo string `json:"correo" validate:"required,email"`
}

type RestablecerEntrada struct {
	Token string `json:"token" validate:"required"`
	Nueva string `json:"nueva" validate:"required,min=8"`
}

type SesionSalida struct {
	Usuario           UsuarioResumen `json:"usuario"`
	Rol               string         `json:"rol"`
	Permisos          []string       `json:"permisos"`
	TerminosAceptados bool           `json:"terminosAceptados"`
}
