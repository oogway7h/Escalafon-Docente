package dto

import "time"

type AutorAuditoriaDTO struct {
	ID        string `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Correo    string `json:"correo"`
}

// RegistroAuditoriaDTO es un evento de la bitácora (CU-06). Usuario es null si la acción no tuvo
// usuario (login con un correo inexistente o solicitud del portal público).
type RegistroAuditoriaDTO struct {
	ID          int64              `json:"id"`
	Fecha       time.Time          `json:"fecha"`
	Usuario     *AutorAuditoriaDTO `json:"usuario"`
	Accion      string             `json:"accion"`
	Entidad     string             `json:"entidad"`
	EntidadID   string             `json:"entidadId"`
	Detalle     map[string]any     `json:"detalle"`
	DireccionIP string             `json:"direccionIp"`
}

type ListaAuditoriaDTO struct {
	Datos  []RegistroAuditoriaDTO `json:"datos"`
	Total  int64                  `json:"total"`
	Pagina int                    `json:"pagina"`
	Tamano int                    `json:"tamano"`
}

type FiltrosAuditoriaDTO struct {
	Acciones []string            `json:"acciones"`
	Usuarios []AutorAuditoriaDTO `json:"usuarios"`
}
