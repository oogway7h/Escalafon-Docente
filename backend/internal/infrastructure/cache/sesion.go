// Package cache implementa la caché de sesión sobre Redis (sección 5), con caída a la base de
// datos cuando Redis no está disponible: la aplicación sigue funcionando, solo más lenta.
package cache

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const ttlSesion = 15 * time.Minute

type SessionCacheRedis struct {
	Cliente *redis.Client
}

func NewSessionCacheRedis(url string) (*SessionCacheRedis, error) {
	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &SessionCacheRedis{Cliente: redis.NewClient(opts)}, nil
}

func claveSesion(usuarioID string) string { return "sesion:" + usuarioID }
func claveRevocado(jti string) string     { return "jwt:rev:" + jti }

// ObtenerSesion lee "{version}|{activo}" de sesion:{usuarioId}. Si Redis falla, lo reporta como
// "no encontrado" (encontrado=false) para que el llamador caiga a la base de datos.
func (c *SessionCacheRedis) ObtenerSesion(ctx context.Context, usuarioID string) (int32, bool, bool, error) {
	val, err := c.Cliente.Get(ctx, claveSesion(usuarioID)).Result()
	if err != nil {
		if err != redis.Nil {
			slog.Warn("redis no disponible, se cae a la base de datos", "error", err)
		}
		return 0, false, false, nil
	}
	partes := strings.SplitN(val, "|", 2)
	if len(partes) != 2 {
		return 0, false, false, nil
	}
	version, err := strconv.Atoi(partes[0])
	if err != nil {
		return 0, false, false, nil
	}
	return int32(version), partes[1] == "1", true, nil
}

func (c *SessionCacheRedis) GuardarSesion(ctx context.Context, usuarioID string, version int32, activo bool) error {
	activoStr := "0"
	if activo {
		activoStr = "1"
	}
	if err := c.Cliente.Set(ctx, claveSesion(usuarioID), fmt.Sprintf("%d|%s", version, activoStr), ttlSesion).Err(); err != nil {
		slog.Warn("no se pudo escribir en la caché de sesión", "error", err)
	}
	return nil
}

// InvalidarSesion borra la clave en lugar de escribir un valor que podría quedar obsoleto por
// concurrencia (sección 5.3): la siguiente petición repuebla desde la base de datos.
func (c *SessionCacheRedis) InvalidarSesion(ctx context.Context, usuarioID string) error {
	if err := c.Cliente.Del(ctx, claveSesion(usuarioID)).Err(); err != nil {
		slog.Warn("no se pudo invalidar la caché de sesión", "error", err)
	}
	return nil
}

func (c *SessionCacheRedis) RevocarJTI(ctx context.Context, jti string, vigencia time.Duration) error {
	if err := c.Cliente.Set(ctx, claveRevocado(jti), "1", vigencia).Err(); err != nil {
		slog.Warn("no se pudo registrar la revocación del jti (logout best-effort)", "error", err)
	}
	return nil
}

// EstaRevocado devuelve el error si Redis no está disponible: quien llama debe fallar cerrado, porque
// sin la lista de revocación un token de una sesión ya cerrada volvería a ser válido hasta expirar.
func (c *SessionCacheRedis) EstaRevocado(ctx context.Context, jti string) (bool, error) {
	n, err := c.Cliente.Exists(ctx, claveRevocado(jti)).Result()
	if err != nil {
		slog.Warn("redis no disponible: no se puede comprobar la lista de revocación", "error", err)
		return false, err
	}
	return n > 0, nil
}
