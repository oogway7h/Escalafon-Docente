package dto

type UsuarioListado struct {
	ID        string `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Correo    string `json:"correo"`
	Telefono  string `json:"telefono"`
	RolID     int16  `json:"rolId"`
	Activo    bool   `json:"activo"`
}

type ListarUsuariosSalida struct {
	Datos  []UsuarioListado `json:"datos"`
	Total  int64            `json:"total"`
	Pagina int              `json:"pagina"`
	Tamano int              `json:"tamano"`
}

type CrearUsuarioEntrada struct {
	Nombres   string `json:"nombres" validate:"required"`
	Apellidos string `json:"apellidos" validate:"required"`
	Correo    string `json:"correo" validate:"required,email"`
	Telefono  string `json:"telefono"`
	RolID     int16  `json:"rolId" validate:"required"`
}

type EditarUsuarioEntrada struct {
	Nombres   string `json:"nombres" validate:"required"`
	Apellidos string `json:"apellidos" validate:"required"`
	Correo    string `json:"correo" validate:"required,email"`
	Telefono  string `json:"telefono"`
}

type CambiarRolEntrada struct {
	RolID int16 `json:"rolId" validate:"required"`
}

type CambiarEstadoEntrada struct {
	Activo *bool `json:"activo" validate:"required"`
}
