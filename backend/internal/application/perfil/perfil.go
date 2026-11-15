// Package perfil implementa los casos de uso de perfil propio de HU-01.
package perfil

import (
	"context"
	"time"

	"github.com/google/uuid"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/usuario"
)

// ObtenerPerfil devuelve los datos del usuario autenticado.
type ObtenerPerfil struct {
	Usuarios usuario.Repository
}

func (c *ObtenerPerfil) Ejecutar(ctx context.Context, usuarioID string) (*usuario.Usuario, error) {
	return c.Usuarios.BuscarPorID(ctx, usuarioID)
}

// ActualizarContacto permite editar solo el teléfono; rol, permisos y correo se rechazan en el DTO HTTP.
type ActualizarContacto struct {
	Usuarios  usuario.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *ActualizarContacto) Ejecutar(ctx context.Context, usuarioID, telefono string) (*usuario.Usuario, error) {
	u, err := c.Usuarios.BuscarPorID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	anterior := u.Telefono
	u.Telefono = telefono
	if err := c.Usuarios.Actualizar(ctx, u); err != nil {
		return nil, err
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &u.ID, Accion: auditoria.AccionActualizarContacto, Entidad: auditoria.EntidadUsuario,
		EntidadID: u.ID, Detalle: map[string]any{"telefonoAnterior": anterior, "telefonoNuevo": telefono},
		Fecha: c.Reloj.Ahora(),
	})
	return u, nil
}

// CambiarContrasena exige la contraseña actual, incrementa version_sesion y emite un JWT nuevo
// solo para la sesión que hizo el cambio; las demás quedan invalidadas por la caché.
type CambiarContrasena struct {
	Usuarios    usuario.Repository
	Roles       rol.Repository
	Auditoria   auditoria.Repository
	Hasher      ports.PasswordHasher
	Tokens      ports.TokenService
	Cache       ports.SessionCache
	Tx          ports.TxManager
	Reloj       ports.Clock
	SesionHoras int
}

type TokenSalida struct {
	Token    string
	ExpiraEn int
}

func (c *CambiarContrasena) Ejecutar(ctx context.Context, usuarioID, actual, nueva string) (*TokenSalida, error) {
	u, err := c.Usuarios.BuscarPorID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	if !c.Hasher.Comparar(u.HashContrasena, actual) {
		return nil, usuario.ErrContrasenaActualIncorrecta
	}
	if !usuario.PoliticaContrasenaValida(nueva) {
		return nil, usuario.ErrPoliticaContrasena
	}
	hashNueva, err := c.Hasher.Hash(nueva)
	if err != nil {
		return nil, err
	}

	ahora := c.Reloj.Ahora()
	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		u.HashContrasena = hashNueva
		u.VersionSesion++
		if err := c.Usuarios.Actualizar(ctx, u); err != nil {
			return err
		}
		return c.Auditoria.Registrar(ctx, &auditoria.Evento{
			UsuarioID: &u.ID, Accion: auditoria.AccionCambiarContrasena, Entidad: auditoria.EntidadUsuario,
			EntidadID: u.ID, Fecha: ahora,
		})
	})
	if err != nil {
		return nil, err
	}
	if err := c.Cache.InvalidarSesion(ctx, u.ID); err != nil {
		return nil, err
	}

	permisos, err := c.Roles.PermisosPorRol(ctx, u.RolID)
	if err != nil {
		return nil, err
	}
	vigencia := time.Duration(c.SesionHoras) * time.Hour
	tok, err := c.Tokens.Emitir(ports.ClaimsSesion{
		UsuarioID: u.ID, RolID: u.RolID, Permisos: permisos, Version: u.VersionSesion, JTI: uuid.NewString(),
	}, vigencia)
	if err != nil {
		return nil, err
	}
	return &TokenSalida{Token: tok, ExpiraEn: int(vigencia.Seconds())}, nil
}
