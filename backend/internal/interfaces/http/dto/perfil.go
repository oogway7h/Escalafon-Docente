package dto

type PerfilSalida struct {
	ID        string `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Correo    string `json:"correo"`
	Telefono  string `json:"telefono"`
	RolID     int16  `json:"rolId"`
}

// camposProhibidosContacto son las claves que, si llegan en el cuerpo de PUT /perfil, se
// rechazan con 403 CAMPO_NO_EDITABLE (sección 7.2): el teléfono es el único campo editable ahí.
var camposProhibidosContacto = []string{"rol", "rolId", "permisos", "correo"}

// CamposProhibidosPresentes revisa el cuerpo crudo (ya decodificado a un mapa) en busca de
// claves que este endpoint no permite editar.
func CamposProhibidosPresentes(cuerpo map[string]any) bool {
	for _, campo := range camposProhibidosContacto {
		if _, ok := cuerpo[campo]; ok {
			return true
		}
	}
	return false
}

type ActualizarContactoEntrada struct {
	Telefono string `json:"telefono" validate:"required"`
}

type CambiarContrasenaEntrada struct {
	Actual string `json:"actual" validate:"required"`
	Nueva  string `json:"nueva" validate:"required,min=8"`
}

type CambiarContrasenaSalida struct {
	Token    string `json:"token"`
	ExpiraEn int    `json:"expiraEn"`
}
