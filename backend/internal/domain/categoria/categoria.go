package categoria

import "context"

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrNoEncontrado erroDominio = "categoría no encontrada"
	ErrDuplicada    erroDominio = "ya existe una categoría con ese nombre"
)

// ErrEnUso se devuelve al eliminar una categoría que docentes o trámites todavía usan.
type ErrEnUso struct{ Docentes, Tramites int64 }

func (e *ErrEnUso) Error() string { return "la categoría está en uso" }

// Requisito es un documento probatorio que exige una categoría.
type Requisito struct {
	ID     int32
	Nombre string
}

// Categoria es una entrada del catálogo normativo (Adjunto, Asociado, Titular A/B/C, Emérito).
type Categoria struct {
	ID              int16
	Nombre          string
	AniosAntiguedad int16
	// Rango de puntos de categoría (Reglamento del Escalafón, art. 18); nil = la categoría no se define por puntaje.
	// PuntajeHasta nil con PuntajeDesde definido significa «o más» (Titular C: 450 o más).
	PuntajeDesde, PuntajeHasta *int16
	Activa          bool
	Requisitos      []Requisito
	Docentes        int64 // docentes que la tienen como categoría actual
	Tramites        int64 // trámites que la postulan
}

type Repository interface {
	Listar(ctx context.Context, soloActivas bool) ([]Categoria, error)
	BuscarPorID(ctx context.Context, id int16) (*Categoria, error)
	// Guardar crea la categoría si ID es 0 o la actualiza, reemplazando sus requisitos por nombres.
	Guardar(ctx context.Context, c *Categoria, requisitos []string) error
	CambiarActiva(ctx context.Context, id int16, activa bool) error
	// Eliminar borra la categoría y sus requisitos; falla con *ErrEnUso si algo la referencia.
	Eliminar(ctx context.Context, id int16) error
}
