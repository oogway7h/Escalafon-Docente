package dto

// ErrorRespuesta es el formato uniforme de error de la API (sección 7).
type ErrorRespuesta struct {
	Codigo   string `json:"codigo"`
	Mensaje  string `json:"mensaje"`
	Detalles any    `json:"detalles,omitempty"`
}
