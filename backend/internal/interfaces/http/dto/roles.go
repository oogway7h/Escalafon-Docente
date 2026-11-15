package dto

type PermisoDTO struct {
	ID      int16  `json:"id"`
	Codigo  string `json:"codigo"`
	Nombre  string `json:"nombre"`
	Modulo  string `json:"modulo"`
	CasoUso string `json:"casoUso"`
}

type RolConConteoDTO struct {
	ID            int16  `json:"id"`
	Nombre        string `json:"nombre"`
	Descripcion   string `json:"descripcion"`
	EsSistema     bool   `json:"esSistema"`
	Activo        bool   `json:"activo"`
	UsuariosCount int64  `json:"usuariosCount"`
}

type RolConPermisosDTO struct {
	ID          int16        `json:"id"`
	Nombre      string       `json:"nombre"`
	Descripcion string       `json:"descripcion"`
	EsSistema   bool         `json:"esSistema"`
	Activo      bool         `json:"activo"`
	Permisos    []PermisoDTO `json:"permisos"`
}

type GuardarRolEntrada struct {
	Nombre      string  `json:"nombre" validate:"required"`
	Descripcion string  `json:"descripcion"`
	PermisoIDs  []int16 `json:"permisoIds"`
}
