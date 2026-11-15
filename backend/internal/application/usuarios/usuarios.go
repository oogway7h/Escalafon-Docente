// Package usuarios implementa los casos de uso de HU-02: Administración de Usuarios.
package usuarios

import (
	"context"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/usuario"
)

type Listar struct{ Usuarios usuario.Repository }

func (c *Listar) Ejecutar(ctx context.Context, filtro usuario.ListarFiltro) ([]usuario.Usuario, int64, error) {
	return c.Usuarios.Listar(ctx, filtro)
}

type Obtener struct{ Usuarios usuario.Repository }

func (c *Obtener) Ejecutar(ctx context.Context, id string) (*usuario.Usuario, error) {
	return c.Usuarios.BuscarPorID(ctx, id)
}

// Crear implementa CrearUsuario (sección 8): rol existente y activo, contraseña aleatoria
// inutilizable, token de ACTIVACION, auditoría en la misma transacción y correo después.
type Crear struct {
	Usuarios             usuario.Repository
	Roles                rol.Repository
	Tokens               token.Repository
	Auditoria            auditoria.Repository
	Hasher               ports.PasswordHasher
	Generador            ports.GeneradorToken
	Mailer               ports.Mailer
	Tx                   ports.TxManager
	Reloj                ports.Clock
	HorasTokenActivacion int
	URLFrontend          string
}

type CrearEntrada struct {
	ActorID   string
	Nombres   string
	Apellidos string
	Correo    string
	Telefono  string
	RolID     int16
}

func (c *Crear) Ejecutar(ctx context.Context, in CrearEntrada) (*usuario.Usuario, error) {
	r, err := c.Roles.BuscarPorID(ctx, in.RolID)
	if err != nil {
		return nil, usuario.ErrRolInvalido
	}
	if !r.Activo {
		return nil, usuario.ErrRolInvalido
	}

	claveAleatoria, _, err := c.Generador.Generar()
	if err != nil {
		return nil, err
	}
	hashClave, err := c.Hasher.Hash(claveAleatoria)
	if err != nil {
		return nil, err
	}

	u := &usuario.Usuario{
		RolID: in.RolID, Nombres: in.Nombres, Apellidos: in.Apellidos, Correo: strings.ToLower(in.Correo),
		Telefono: in.Telefono, HashContrasena: hashClave, VersionSesion: 1, Activo: true,
	}

	ahora := c.Reloj.Ahora()
	var tokenPlano string
	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		if err := c.Usuarios.Crear(ctx, u); err != nil {
			return err
		}
		plano, hash, err := c.Generador.Generar()
		if err != nil {
			return err
		}
		tokenPlano = plano
		t := &token.TokenRecuperacion{
			UsuarioID: u.ID, Proposito: token.PropositoActivacion, TokenHash: hash,
			ExpiraEn: ahora.Add(time.Duration(c.HorasTokenActivacion) * time.Hour), CreadoEn: ahora,
		}
		if err := c.Tokens.Crear(ctx, t); err != nil {
			return err
		}
		return c.Auditoria.Registrar(ctx, &auditoria.Evento{
			UsuarioID: &in.ActorID, Accion: auditoria.AccionCrear, Entidad: auditoria.EntidadUsuario, EntidadID: u.ID, Fecha: ahora,
		})
	})
	if err != nil {
		return nil, err
	}

	enlace := c.URLFrontend + "/activar?token=" + tokenPlano
	_ = c.Mailer.EnviarActivacion(u.Correo, u.Nombres, enlace)
	return u, nil
}

// Editar implementa EditarUsuario: datos personales y correo, auditando valores anteriores y nuevos.
type Editar struct {
	Usuarios  usuario.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

type EditarEntrada struct {
	ActorID   string
	Nombres   string
	Apellidos string
	Correo    string
	Telefono  string
}

func (c *Editar) Ejecutar(ctx context.Context, id string, in EditarEntrada) (*usuario.Usuario, error) {
	u, err := c.Usuarios.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	anterior := map[string]any{"nombres": u.Nombres, "apellidos": u.Apellidos, "correo": u.Correo, "telefono": u.Telefono}

	u.Nombres, u.Apellidos, u.Telefono = in.Nombres, in.Apellidos, in.Telefono
	u.Correo = strings.ToLower(in.Correo)
	if err := c.Usuarios.Actualizar(ctx, u); err != nil {
		return nil, err
	}

	nuevo := map[string]any{"nombres": u.Nombres, "apellidos": u.Apellidos, "correo": u.Correo, "telefono": u.Telefono}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &in.ActorID, Accion: auditoria.AccionEditar, Entidad: auditoria.EntidadUsuario, EntidadID: u.ID,
		Detalle: map[string]any{"anterior": anterior, "nuevo": nuevo}, Fecha: c.Reloj.Ahora(),
	})
	return u, nil
}

// CambiarRol y CambiarEstado comparten dependencias e invalidan la caché del usuario afectado.
type CambiarRol struct {
	Usuarios  usuario.Repository
	Roles     rol.Repository
	Auditoria auditoria.Repository
	Cache     ports.SessionCache
	Tx        ports.TxManager
	Reloj     ports.Clock
}

func (c *CambiarRol) Ejecutar(ctx context.Context, actorID, objetivoID string, nuevoRolID int16) error {
	if actorID == objetivoID {
		return usuario.ErrAutocambioRol
	}
	r, err := c.Roles.BuscarPorID(ctx, nuevoRolID)
	if err != nil || !r.Activo {
		return usuario.ErrRolInvalido
	}
	u, err := c.Usuarios.BuscarPorID(ctx, objetivoID)
	if err != nil {
		return err
	}
	rolAnterior := u.RolID
	ahora := c.Reloj.Ahora()

	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		u.RolID = nuevoRolID
		u.VersionSesion++
		if err := c.Usuarios.Actualizar(ctx, u); err != nil {
			return err
		}
		return c.Auditoria.Registrar(ctx, &auditoria.Evento{
			UsuarioID: &actorID, Accion: auditoria.AccionCambioRol, Entidad: auditoria.EntidadUsuario, EntidadID: u.ID,
			Detalle: map[string]any{"rolAnterior": rolAnterior, "rolNuevo": nuevoRolID}, Fecha: ahora,
		})
	})
	if err != nil {
		return err
	}
	return c.Cache.InvalidarSesion(ctx, u.ID)
}

type CambiarEstado struct {
	Usuarios  usuario.Repository
	Auditoria auditoria.Repository
	Cache     ports.SessionCache
	Tx        ports.TxManager
	Reloj     ports.Clock
}

func (c *CambiarEstado) Ejecutar(ctx context.Context, actorID, objetivoID string, activo bool) error {
	if actorID == objetivoID && !activo {
		return usuario.ErrAutodesactivacion
	}
	u, err := c.Usuarios.BuscarPorID(ctx, objetivoID)
	if err != nil {
		return err
	}
	ahora := c.Reloj.Ahora()
	accion := auditoria.AccionBaja
	if activo {
		accion = auditoria.AccionReactivar
	}

	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		u.Activo = activo
		u.VersionSesion++
		if err := c.Usuarios.Actualizar(ctx, u); err != nil {
			return err
		}
		return c.Auditoria.Registrar(ctx, &auditoria.Evento{
			UsuarioID: &actorID, Accion: accion, Entidad: auditoria.EntidadUsuario, EntidadID: u.ID, Fecha: ahora,
		})
	})
	if err != nil {
		return err
	}
	return c.Cache.InvalidarSesion(ctx, u.ID)
}

// ReenviarActivacion solo procede si el usuario nunca activó su cuenta, e invalida los tokens
// de activación previos que quedaron pendientes.
type ReenviarActivacion struct {
	Usuarios             usuario.Repository
	Tokens               token.Repository
	Auditoria            auditoria.Repository
	Generador            ports.GeneradorToken
	Mailer               ports.Mailer
	Reloj                ports.Clock
	HorasTokenActivacion int
	URLFrontend          string
}

func (c *ReenviarActivacion) Ejecutar(ctx context.Context, actorID, id string) error {
	u, err := c.Usuarios.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	yaActivada, err := c.Tokens.ExisteTokenUsado(ctx, u.ID, token.PropositoActivacion)
	if err != nil {
		return err
	}
	if yaActivada {
		return usuario.ErrCuentaYaActivada
	}

	ahora := c.Reloj.Ahora()
	if err := c.Tokens.InvalidarPendientes(ctx, u.ID, token.PropositoActivacion, ahora); err != nil {
		return err
	}
	plano, hash, err := c.Generador.Generar()
	if err != nil {
		return err
	}
	t := &token.TokenRecuperacion{
		UsuarioID: u.ID, Proposito: token.PropositoActivacion, TokenHash: hash,
		ExpiraEn: ahora.Add(time.Duration(c.HorasTokenActivacion) * time.Hour), CreadoEn: ahora,
	}
	if err := c.Tokens.Crear(ctx, t); err != nil {
		return err
	}
	// Se audita la acción, nunca el token ni el enlace.
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actorID, Accion: auditoria.AccionReenviarActivacion, Entidad: auditoria.EntidadUsuario, EntidadID: u.ID, Fecha: ahora,
	})

	enlace := c.URLFrontend + "/activar?token=" + plano
	return c.Mailer.EnviarActivacion(u.Correo, u.Nombres, enlace)
}
