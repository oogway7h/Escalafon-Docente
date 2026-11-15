package docente

import (
	"context"
	"time"
)

type erroDominio string

func (e erroDominio) Error() string { return string(e) }

const (
	ErrNoEncontrado      erroDominio = "docente no encontrado"
	ErrCedulaDuplicada   erroDominio = "ya existe un docente con esa cédula"
	ErrCodigoDuplicado   erroDominio = "ya existe un docente con ese código"
	ErrFacultadInvalida  erroDominio = "la facultad no existe"
	ErrCategoriaInvalida erroDominio = "la categoría no existe o está inactiva"
)

// Legajo es una carpeta física del docente (un docente puede tener varias).
type Legajo struct {
	Gaveta      string
	Orden       int
	Anexos      string
	Observacion string
}

// Facultad es la unidad académica (facultad, carrera o decanato) a la que pertenece el docente.
type Facultad struct {
	ID     int16
	Nombre string
}

type Docente struct {
	ID          string
	Codigo      string // código que la universidad otorga al docente; identificador principal
	Cedula      string // opcional; inmutable una vez asignada
	Nombres     string
	Apellidos   string
	Correo      string
	Telefono    string
	CategoriaID *int16
	Categoria   string
	FacultadID  *int16
	Facultad    string
	Legajos     []Legajo
	Activo      bool
	CreadoEn    time.Time
}

type Repository interface {
	// Listar filtra por texto (cédula, nombres o apellidos); vacío = todos.
	Listar(ctx context.Context, busqueda string, soloActivos bool) ([]Docente, error)
	// ListarPagina es Listar con paginación; devuelve también el total de coincidencias.
	ListarPagina(ctx context.Context, busqueda string, soloActivos bool, pagina, tamano int) ([]Docente, int64, error)
	// Sugerir es la búsqueda mientras se escribe: código o cédula por prefijo, nombre por fragmento sin
	// tildes en cualquier orden; solo activos, sin legajos ni total (debe responder rápido).
	Sugerir(ctx context.Context, texto string, limite int) ([]Docente, error)
	BuscarPorID(ctx context.Context, id string) (*Docente, error)
	// BuscarPorIdentificador busca por código o por cédula.
	BuscarPorIdentificador(ctx context.Context, valor string) (*Docente, error)
	Crear(ctx context.Context, d *Docente) error
	// Actualizar no toca la cédula.
	Actualizar(ctx context.Context, d *Docente) error
	CambiarActivo(ctx context.Context, id string, activo bool) error
	ListarFacultades(ctx context.Context) ([]Facultad, error)
}
