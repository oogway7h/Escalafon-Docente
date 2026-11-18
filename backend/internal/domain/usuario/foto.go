package usuario

import (
	"context"
	"time"
)

// Foto es la imagen de perfil de un usuario: el objeto vive en el almacén y aquí solo se guarda su referencia.
type Foto struct {
	UsuarioID     string
	Clave         string
	TipoMime      string
	ActualizadoEn time.Time
}

// FotoRepository guarda la referencia a la foto (una por usuario).
type FotoRepository interface {
	// Obtener devuelve la foto del usuario o ErrFotoNoEncontrada.
	Obtener(ctx context.Context, usuarioID string) (*Foto, error)
	// Guardar crea o reemplaza la foto del usuario.
	Guardar(ctx context.Context, f *Foto) error
	Eliminar(ctx context.Context, usuarioID string) error
}

const (
	// FotoMaxBytes es el tamaño máximo de la foto de perfil.
	FotoMaxBytes = 2 << 20
	// FotoMaxLado es el lado máximo en píxeles; evita imágenes enormes que consumen memoria al mostrarse.
	FotoMaxLado = 4096
)
