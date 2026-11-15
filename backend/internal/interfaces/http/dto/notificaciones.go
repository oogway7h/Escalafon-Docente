package dto

import "time"

type NotificacionDTO struct {
	ID            string     `json:"id"`
	TramiteCodigo string     `json:"tramiteCodigo"`
	EstadoTramite string     `json:"estadoTramite"`
	Destino       string     `json:"destino"`
	Estado        string     `json:"estado"`
	Intentos      int        `json:"intentos"`
	Error         string     `json:"error"`
	CreadoEn      time.Time  `json:"creadoEn"`
	EnviadoEn     *time.Time `json:"enviadoEn"`
}

type ListaNotificacionesDTO struct {
	Datos  []NotificacionDTO `json:"datos"`
	Total  int64             `json:"total"`
	Pagina int               `json:"pagina"`
	Tamano int               `json:"tamano"`
}
