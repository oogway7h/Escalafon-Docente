package perfil

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg" // registra los decodificadores que usa DecodeConfig
	_ "image/png"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/usuario"
)

// tipoFoto valida el contenido real (no el tipo que declara el cliente) y devuelve su tipo MIME.
func tipoFoto(datos []byte) (string, error) {
	if len(datos) == 0 {
		return "", usuario.ErrFotoInvalida
	}
	if len(datos) > usuario.FotoMaxBytes {
		return "", usuario.ErrFotoGrande
	}
	cfg, formato, err := image.DecodeConfig(bytes.NewReader(datos))
	if err != nil || cfg.Width < 1 || cfg.Height < 1 || cfg.Width > usuario.FotoMaxLado || cfg.Height > usuario.FotoMaxLado {
		return "", usuario.ErrFotoInvalida
	}
	switch formato {
	case "jpeg":
		return "image/jpeg", nil
	case "png":
		return "image/png", nil
	}
	return "", usuario.ErrFotoInvalida
}

func claveFoto(usuarioID string) string { return "usuarios/" + usuarioID + "/foto" }

// SubirFoto guarda (o reemplaza) la foto de perfil del propio usuario.
type SubirFoto struct {
	Fotos     usuario.FotoRepository
	Almacen   ports.ObjectStorage
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *SubirFoto) Ejecutar(ctx context.Context, usuarioID string, datos []byte) (*usuario.Foto, error) {
	tipo, err := tipoFoto(datos)
	if err != nil {
		return nil, err
	}
	clave := claveFoto(usuarioID)
	if err := c.Almacen.Guardar(ctx, clave, tipo, datos); err != nil {
		return nil, err
	}
	f := &usuario.Foto{UsuarioID: usuarioID, Clave: clave, TipoMime: tipo, ActualizadoEn: c.Reloj.Ahora()}
	if err := c.Fotos.Guardar(ctx, f); err != nil {
		return nil, err
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &usuarioID, Accion: auditoria.AccionActualizarContacto, Entidad: auditoria.EntidadUsuario,
		EntidadID: usuarioID, Detalle: map[string]any{"foto": "subida"}, Fecha: f.ActualizadoEn,
	})
	return f, nil
}

// QuitarFoto elimina la foto de perfil. Quitar una foto que no existe no es un error.
type QuitarFoto struct {
	Fotos     usuario.FotoRepository
	Almacen   ports.ObjectStorage
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *QuitarFoto) Ejecutar(ctx context.Context, usuarioID string) error {
	f, err := c.Fotos.Obtener(ctx, usuarioID)
	if err == usuario.ErrFotoNoEncontrada {
		return nil
	}
	if err != nil {
		return err
	}
	if err := c.Fotos.Eliminar(ctx, usuarioID); err != nil {
		return err
	}
	_ = c.Almacen.Eliminar(ctx, f.Clave) // si falla (permisos del bucket) queda un objeto huérfano sin referencia
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &usuarioID, Accion: auditoria.AccionActualizarContacto, Entidad: auditoria.EntidadUsuario,
		EntidadID: usuarioID, Detalle: map[string]any{"foto": "quitada"}, Fecha: c.Reloj.Ahora(),
	})
	return nil
}

// FotoEnlace es un enlace temporal de lectura a la foto.
type FotoEnlace struct {
	URL       string
	TipoMime  string
	Version   int64 // cambia con cada foto nueva: sirve para invalidar cachés
	VigenciaS int
}

// ObtenerFoto devuelve un enlace firmado y de corta vida a la foto del usuario.
type ObtenerFoto struct {
	Fotos   usuario.FotoRepository
	Almacen ports.ObjectStorage
}

const vigenciaEnlaceFoto = 10 * time.Minute

func (c *ObtenerFoto) Ejecutar(ctx context.Context, usuarioID string) (*FotoEnlace, error) {
	f, err := c.Fotos.Obtener(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	url, err := c.Almacen.URLFirmada(ctx, f.Clave, vigenciaEnlaceFoto)
	if err != nil {
		return nil, err
	}
	return &FotoEnlace{URL: url, TipoMime: f.TipoMime, Version: f.ActualizadoEn.Unix(), VigenciaS: int(vigenciaEnlaceFoto.Seconds())}, nil
}
