// Package auth implementa los casos de uso de HU-01: Acceso y Gestión de Credenciales.
package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/usuario"
)

// UsuarioResumen es lo que el cliente necesita saber de sí mismo tras autenticarse.
type UsuarioResumen struct {
	ID        string
	Nombres   string
	Apellidos string
	Correo    string
	RolID     int16
}

// TokenSalida es la forma común de una respuesta que entrega un JWT de sesión nuevo.
type TokenSalida struct {
	Token    string
	ExpiraEn int
}

// hashFicticio se compara cuando el usuario no existe, para no filtrar por temporización qué
// correos están registrados (sección 6).
const hashFicticio = "$2a$12$CwTycUXWue0Thq9StjUM0uJ8Nn3vD2G8IjF.Zi9v9y0k7Xz1B1J0S"

// AutenticarUsuario valida credenciales, aplica bloqueo por intentos y emite el JWT de sesión.
type AutenticarUsuario struct {
	Usuarios       usuario.Repository
	Roles          rol.Repository
	Auditoria      auditoria.Repository
	Hasher         ports.PasswordHasher
	Tokens         ports.TokenService
	Reloj          ports.Clock
	MaxIntentos    int
	MinutosBloqueo int
	SesionHoras    int
}

type AutenticarEntrada struct {
	Correo      string
	Contrasena  string
	DireccionIP string
}

type AutenticarSalida struct {
	Token    string
	ExpiraEn int
	Usuario  UsuarioResumen
	Permisos []string
	// TerminosAceptados: false obliga a mostrar el acuerdo de usuario antes de entrar al panel.
	TerminosAceptados bool
}

func (c *AutenticarUsuario) Ejecutar(ctx context.Context, in AutenticarEntrada) (*AutenticarSalida, error) {
	ahora := c.Reloj.Ahora()
	correo := strings.ToLower(strings.TrimSpace(in.Correo))

	u, err := c.Usuarios.BuscarPorCorreo(ctx, correo)
	if err != nil && !errors.Is(err, usuario.ErrNoEncontrado) {
		return nil, err
	}
	if u == nil || !u.Activo {
		c.Hasher.Comparar(hashFicticio, in.Contrasena)
		c.auditar(ctx, nil, auditoria.AccionLoginFallido, in.DireccionIP, ahora)
		return nil, usuario.ErrCredencialesInvalidas
	}
	// La contraseña se comprueba siempre: el bloqueo solo se informa a quien la conoce. Con una contraseña
	// incorrecta la respuesta es la misma que para un correo inexistente (no delata qué cuentas existen).
	correcta := c.Hasher.Comparar(u.HashContrasena, in.Contrasena)
	if u.EstaBloqueado(ahora) {
		if correcta {
			return nil, &usuario.ErrCuentaBloqueada{Hasta: *u.BloqueadoHasta}
		}
		c.auditar(ctx, &u.ID, auditoria.AccionLoginFallido, in.DireccionIP, ahora)
		return nil, usuario.ErrCredencialesInvalidas
	}
	if !correcta {
		if u.BloqueadoHasta != nil { // bloqueo vencido: la cuenta empieza otra vez desde cero
			u.IntentosFallidos, u.BloqueadoHasta = 0, nil
		}
		u.IntentosFallidos++
		if int(u.IntentosFallidos) >= c.MaxIntentos {
			hasta := ahora.Add(time.Duration(c.MinutosBloqueo) * time.Minute)
			u.BloqueadoHasta = &hasta
		}
		if err := c.Usuarios.Actualizar(ctx, u); err != nil {
			return nil, err
		}
		c.auditar(ctx, &u.ID, auditoria.AccionLoginFallido, in.DireccionIP, ahora)
		return nil, usuario.ErrCredencialesInvalidas
	}

	permisos, err := c.Roles.PermisosPorRol(ctx, u.RolID)
	if err != nil {
		return nil, err
	}

	u.IntentosFallidos = 0
	u.BloqueadoHasta = nil
	u.UltimoAcceso = &ahora
	if err := c.Usuarios.Actualizar(ctx, u); err != nil {
		return nil, err
	}

	vigencia := time.Duration(c.SesionHoras) * time.Hour
	tok, err := c.Tokens.Emitir(ports.ClaimsSesion{
		UsuarioID: u.ID,
		RolID:     u.RolID,
		Permisos:  permisos,
		Version:   u.VersionSesion,
		JTI:       uuid.NewString(),
	}, vigencia)
	if err != nil {
		return nil, err
	}

	c.auditar(ctx, &u.ID, auditoria.AccionLogin, in.DireccionIP, ahora)

	return &AutenticarSalida{
		Token:    tok,
		ExpiraEn: int(vigencia.Seconds()),
		Usuario: UsuarioResumen{
			ID: u.ID, Nombres: u.Nombres, Apellidos: u.Apellidos, Correo: u.Correo, RolID: u.RolID,
		},
		Permisos:          permisos,
		TerminosAceptados: u.AceptoTerminos(),
	}, nil
}

func (c *AutenticarUsuario) auditar(ctx context.Context, usuarioID *string, accion, ip string, cuando time.Time) {
	entidadID := ""
	if usuarioID != nil {
		entidadID = *usuarioID
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: usuarioID, Accion: accion, Entidad: auditoria.EntidadUsuario,
		EntidadID: entidadID, DireccionIP: ip, Fecha: cuando,
	})
}

// CerrarSesion revoca el jti actual para que no pueda reutilizarse antes de que expire.
type CerrarSesion struct {
	Cache ports.SessionCache
}

func (c *CerrarSesion) Ejecutar(ctx context.Context, jti string, vigenciaRestante time.Duration) error {
	if vigenciaRestante <= 0 {
		return nil
	}
	return c.Cache.RevocarJTI(ctx, jti, vigenciaRestante)
}

// ObtenerSesion resuelve el usuario y su rol a partir del claim `sub`. Los permisos son los del propio token
// (los que la API hace cumplir): así el menú coincide con lo que el usuario puede hacer y se ahorra una consulta
// en cada navegación del panel.
type ObtenerSesion struct {
	Usuarios usuario.Repository
	Roles    rol.Repository
}

type ObtenerSesionSalida struct {
	Usuario           UsuarioResumen
	Rol               string
	Permisos          []string
	TerminosAceptados bool
}

func (c *ObtenerSesion) Ejecutar(ctx context.Context, usuarioID string, permisos []string) (*ObtenerSesionSalida, error) {
	u, err := c.Usuarios.BuscarPorID(ctx, usuarioID)
	if err != nil {
		return nil, err
	}
	r, err := c.Roles.BuscarPorID(ctx, u.RolID)
	if err != nil {
		return nil, err
	}
	return &ObtenerSesionSalida{
		Usuario:           UsuarioResumen{ID: u.ID, Nombres: u.Nombres, Apellidos: u.Apellidos, Correo: u.Correo, RolID: u.RolID},
		Rol:               r.Nombre,
		Permisos:          permisos,
		TerminosAceptados: u.AceptoTerminos(),
	}, nil
}

// AceptarTerminos registra que el usuario aceptó la versión vigente del acuerdo de usuario.
type AceptarTerminos struct {
	Usuarios  usuario.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *AceptarTerminos) Ejecutar(ctx context.Context, usuarioID, ip string) error {
	u, err := c.Usuarios.BuscarPorID(ctx, usuarioID)
	if err != nil {
		return err
	}
	if u.AceptoTerminos() {
		return nil
	}
	ahora := c.Reloj.Ahora()
	u.TerminosVersion, u.TerminosAceptadosEn = usuario.VersionTerminos, &ahora
	if err := c.Usuarios.Actualizar(ctx, u); err != nil {
		return err
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &u.ID, Accion: auditoria.AccionAceptarTerminos, Entidad: auditoria.EntidadUsuario,
		EntidadID: u.ID, Detalle: map[string]any{"version": usuario.VersionTerminos}, DireccionIP: ip, Fecha: ahora,
	})
	return nil
}

// SolicitarRecuperacion crea un token de un solo uso y envía el enlace por correo.
// Siempre debe interpretarse como éxito (202) exista o no el correo, para no filtrar información.
type SolicitarRecuperacion struct {
	Usuarios        usuario.Repository
	Tokens          token.Repository
	Generador       ports.GeneradorToken
	Mailer          ports.Mailer
	Reloj           ports.Clock
	MinutosVigencia int
	URLFrontend     string
	// Asincrono envía el correo en segundo plano: así la respuesta tarda lo mismo exista o no la cuenta.
	Asincrono bool
}

const maxRecuperacionesPorHora = 3

func (c *SolicitarRecuperacion) Ejecutar(ctx context.Context, correo string) error {
	u, err := c.Usuarios.BuscarPorCorreo(ctx, strings.ToLower(strings.TrimSpace(correo)))
	if err != nil && !errors.Is(err, usuario.ErrNoEncontrado) {
		return err
	}
	if u == nil || !u.Activo {
		return nil
	}

	ahora := c.Reloj.Ahora()
	// Tope por cuenta: evita inundar el buzón de un usuario con correos de recuperación.
	recientes, err := c.Tokens.ContarDesde(ctx, u.ID, token.PropositoRecuperacion, ahora.Add(-time.Hour))
	if err != nil {
		return err
	}
	if recientes >= maxRecuperacionesPorHora {
		return nil
	}

	plano, hash, err := c.Generador.Generar()
	if err != nil {
		return err
	}
	t := &token.TokenRecuperacion{
		UsuarioID: u.ID,
		Proposito: token.PropositoRecuperacion,
		TokenHash: hash,
		ExpiraEn:  ahora.Add(time.Duration(c.MinutosVigencia) * time.Minute),
		CreadoEn:  ahora,
	}
	if err := c.Tokens.Crear(ctx, t); err != nil {
		return err
	}

	enlace := c.URLFrontend + "/restablecer?token=" + plano
	if c.Asincrono {
		go func() {
			if err := c.Mailer.EnviarRecuperacion(u.Correo, u.Nombres, enlace); err != nil {
				slog.Warn("no se pudo enviar el correo de recuperación", "error", err)
			}
		}()
		return nil
	}
	return c.Mailer.EnviarRecuperacion(u.Correo, u.Nombres, enlace)
}

// RestablecerContrasena consume un token de recuperación o activación y fija una nueva contraseña.
type RestablecerContrasena struct {
	Usuarios          usuario.Repository
	Tokens            token.Repository
	Auditoria         auditoria.Repository
	Hasher            ports.PasswordHasher
	Generador         ports.GeneradorToken
	Cache             ports.SessionCache
	Tx                ports.TxManager
	Reloj             ports.Clock
	PropositoEsperado token.Proposito
}

func (c *RestablecerContrasena) Ejecutar(ctx context.Context, tokenPlano, nuevaContrasena string) error {
	if !usuario.PoliticaContrasenaValida(nuevaContrasena) {
		return usuario.ErrPoliticaContrasena
	}

	hash := c.Generador.Hash(tokenPlano)
	t, err := c.Tokens.BuscarPorHash(ctx, hash)
	if err != nil {
		return token.ErrInvalido
	}
	ahora := c.Reloj.Ahora()
	if t.Proposito != c.PropositoEsperado || !t.EsValido(ahora) {
		return token.ErrInvalido
	}

	u, err := c.Usuarios.BuscarPorID(ctx, t.UsuarioID)
	if err != nil {
		return err
	}
	hashNueva, err := c.Hasher.Hash(nuevaContrasena)
	if err != nil {
		return err
	}

	accion := auditoria.AccionRestablecer
	if c.PropositoEsperado == token.PropositoActivacion {
		accion = auditoria.AccionActivar
	}

	err = c.Tx.WithinTx(ctx, func(ctx context.Context) error {
		u.HashContrasena = hashNueva
		u.VersionSesion++
		u.IntentosFallidos = 0
		u.BloqueadoHasta = nil
		if err := c.Usuarios.Actualizar(ctx, u); err != nil {
			return err
		}
		if err := c.Tokens.MarcarUsado(ctx, t.ID, ahora); err != nil {
			return err
		}
		return c.Auditoria.Registrar(ctx, &auditoria.Evento{
			UsuarioID: &u.ID, Accion: accion, Entidad: auditoria.EntidadUsuario, EntidadID: u.ID, Fecha: ahora,
		})
	})
	if err != nil {
		return err
	}

	return c.Cache.InvalidarSesion(ctx, u.ID)
}
