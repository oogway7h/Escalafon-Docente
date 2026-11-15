package dto

import "time"

// --- CU-19 al CU-24: legajos, solicitudes y custodia ---

type CarpetaDTO struct {
	ID            int32  `json:"id"`
	Codigo        string `json:"codigo"`
	DocenteID     string `json:"docenteId"`
	DocenteCodigo string `json:"docenteCodigo"`
	DocenteNombre string `json:"docenteNombre"`
	Gaveta        string `json:"gaveta"`
	Orden         int    `json:"orden"`
	Anexos        string `json:"anexos"`
	Observacion   string `json:"observacion"`
	Estado        string `json:"estado"`
}

type ListaLegajosDTO struct {
	Datos  []CarpetaDTO `json:"datos"`
	Total  int64        `json:"total"`
	Pagina int          `json:"pagina"`
	Tamano int          `json:"tamano"`
}

type GuardarLegajoEntrada struct {
	DocenteID   string `json:"docenteId"` // obligatorio al crear
	Codigo      string `json:"codigo" validate:"max=20"`
	Gaveta      string `json:"gaveta" validate:"max=10"`
	Orden       int    `json:"orden" validate:"min=0"`
	Anexos      string `json:"anexos" validate:"max=80"`
	Observacion string `json:"observacion" validate:"max=200"`
}

type SolicitarCarpetaEntrada struct {
	TramiteID string `json:"tramiteId" validate:"omitempty,uuid"` // opcional: trámite al que se liga la solicitud
	Motivo    string `json:"motivo" validate:"required"`
}

type SolicitudCarpetaDTO struct {
	ID            string     `json:"id"`
	LegajoID      int32      `json:"legajoId"`
	LegajoCodigo  string     `json:"legajoCodigo"`
	DocenteNombre string     `json:"docenteNombre"`
	TramiteID     string     `json:"tramiteId"`
	TramiteCodigo string     `json:"tramiteCodigo"`
	Solicitante   string     `json:"solicitante"`
	Motivo        string     `json:"motivo"`
	Estado        string     `json:"estado"`
	Respuesta     string     `json:"respuesta"`
	CreadoEn      time.Time  `json:"creadoEn"`
	AtendidoEn    *time.Time `json:"atendidoEn"`
	FechaLimite   *string    `json:"fechaLimite"`
}

type ListaSolicitudesDTO struct {
	Datos  []SolicitudCarpetaDTO `json:"datos"`
	Total  int64                 `json:"total"`
	Pagina int                   `json:"pagina"`
	Tamano int                   `json:"tamano"`
}

type AtenderSolicitudEntrada struct {
	Aprobar bool   `json:"aprobar"`
	Nota    string `json:"nota"` // obligatoria al rechazar
}

type DespacharEntrada struct {
	FechaLimite    string `json:"fechaLimite" validate:"required"` // AAAA-MM-DD
	EstadoMaterial string `json:"estadoMaterial" validate:"required"`
	Observacion    string `json:"observacion"`
}

type DevolverEntrada struct {
	EstadoMaterial string `json:"estadoMaterial" validate:"required"`
	Observacion    string `json:"observacion"`
}

type MovimientoDTO struct {
	ID             int64     `json:"id"`
	LegajoCodigo   string    `json:"legajoCodigo"`
	DocenteNombre  string    `json:"docenteNombre"`
	Tipo           string    `json:"tipo"`
	Custodio       string    `json:"custodio"`
	FechaLimite    *string   `json:"fechaLimite"`
	EstadoMaterial string    `json:"estadoMaterial"`
	Observacion    string    `json:"observacion"`
	RegistradoPor  string    `json:"registradoPor"`
	Fecha          time.Time `json:"fecha"`
	Vencido        bool      `json:"vencido"`
}

type ListaMovimientosDTO struct {
	Datos  []MovimientoDTO `json:"datos"`
	Total  int64           `json:"total"`
	Pagina int             `json:"pagina"`
	Tamano int             `json:"tamano"`
}
