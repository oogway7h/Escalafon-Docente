package dto

import "time"

// --- CU-25: tablero de indicadores ---

type PeriodoDTO struct {
	Desde *string `json:"desde"` // AAAA-MM-DD, inclusive; null = sin límite
	Hasta *string `json:"hasta"`
}

type ConteoEstadoDTO struct {
	Estado   string `json:"estado"`
	Cantidad int64  `json:"cantidad"`
}

type ResumenTramitesDTO struct {
	Radicados    int64   `json:"radicados"`
	EnCurso      int64   `json:"enCurso"`
	Resueltos    int     `json:"resueltos"`
	PromedioDias float64 `json:"promedioDias"`
}

type TiempoResolucionDTO struct {
	Tipo         string  `json:"tipo"`
	Categoria    string  `json:"categoria"`
	Cantidad     int     `json:"cantidad"`
	PromedioDias float64 `json:"promedioDias"`
	MinimoDias   float64 `json:"minimoDias"`
	MaximoDias   float64 `json:"maximoDias"`
}

type CargaAnalistaDTO struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Activos    int64  `json:"activos"`
	EnFacultad int64  `json:"enFacultad"`
	Resueltos  int64  `json:"resueltos"`
}

type AlertaMoraDTO struct {
	LegajoCodigo  string    `json:"legajoCodigo"`
	DocenteNombre string    `json:"docenteNombre"`
	Custodio      string    `json:"custodio"`
	TramiteCodigo string    `json:"tramiteCodigo"`
	Salida        time.Time `json:"salida"`
	FechaLimite   *string   `json:"fechaLimite"`
	DiasPrestado  int       `json:"diasPrestado"`
	DiasMora      int       `json:"diasMora"`
	Motivo        string    `json:"motivo"`
}

type TableroDTO struct {
	GeneradoEn          time.Time             `json:"generadoEn"`
	Periodo             PeriodoDTO            `json:"periodo"`
	DiasMaximosPrestamo int                   `json:"diasMaximosPrestamo"`
	Tramites            ResumenTramitesDTO    `json:"tramites"`
	PorEstado           []ConteoEstadoDTO     `json:"porEstado"`
	Tiempos             []TiempoResolucionDTO `json:"tiempos"`
	Analistas           []CargaAnalistaDTO    `json:"analistas"`
	PrestamosVigentes   int                   `json:"prestamosVigentes"`
	Alertas             []AlertaMoraDTO       `json:"alertas"`
}

// --- CU-26: reportes ---

// InterpretarReporteEntrada: la frase del usuario y las categorías vigentes que el formulario ya conoce.
type InterpretarReporteEntrada struct {
	Frase      string `json:"frase" validate:"required,max=300"`
	Categorias []struct {
		ID     int16  `json:"id"`
		Nombre string `json:"nombre"`
	} `json:"categorias" validate:"max=60"`
}

type GenerarReporteEntrada struct {
	Tipo         string `json:"tipo" validate:"required,oneof=TRAMITES INVENTARIO TIEMPOS"`
	Formato      string `json:"formato" validate:"required,oneof=PDF XLSX"`
	Desde        string `json:"desde"` // AAAA-MM-DD, inclusive
	Hasta        string `json:"hasta"` // AAAA-MM-DD, inclusive
	Estado       string `json:"estado"`
	TipoTramite  string `json:"tipoTramite"`
	CategoriaID  int16  `json:"categoriaId" validate:"min=0"`
	EstadoLegajo string `json:"estadoLegajo"`
}
