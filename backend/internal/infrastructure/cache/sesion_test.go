package cache

import (
	"context"
	"os"
	"testing"
	"time"
)

func clienteDePrueba(t *testing.T) *SessionCacheRedis {
	t.Helper()
	url := os.Getenv("REDIS_URL")
	if url == "" {
		t.Skip("REDIS_URL no está definida; omitiendo pruebas de integración")
	}
	c, err := NewSessionCacheRedis(url)
	if err != nil {
		t.Fatalf("no se pudo crear el cliente: %v", err)
	}
	if err := c.Cliente.Ping(context.Background()).Err(); err != nil {
		t.Skipf("Redis no está disponible: %v", err)
	}
	return c
}

func TestSessionCacheRedis_GuardarObtenerInvalidar(t *testing.T) {
	c := clienteDePrueba(t)
	ctx := context.Background()
	usuarioID := "usuario-prueba-sesion"
	defer c.Cliente.Del(ctx, claveSesion(usuarioID))

	if err := c.GuardarSesion(ctx, usuarioID, 5, true); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	version, activo, encontrado, err := c.ObtenerSesion(ctx, usuarioID)
	if err != nil || !encontrado || version != 5 || !activo {
		t.Fatalf("esperaba version=5 activo=true encontrado=true, obtuve %d %v %v (err=%v)", version, activo, encontrado, err)
	}

	if err := c.InvalidarSesion(ctx, usuarioID); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	_, _, encontrado, _ = c.ObtenerSesion(ctx, usuarioID)
	if encontrado {
		t.Fatalf("la sesión debería haberse invalidado")
	}
}

func TestSessionCacheRedis_RevocarJTI(t *testing.T) {
	c := clienteDePrueba(t)
	ctx := context.Background()
	jti := "jti-prueba"
	defer c.Cliente.Del(ctx, claveRevocado(jti))

	revocado, _ := c.EstaRevocado(ctx, jti)
	if revocado {
		t.Fatalf("no debería estar revocado todavía")
	}
	if err := c.RevocarJTI(ctx, jti, time.Minute); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	revocado, _ = c.EstaRevocado(ctx, jti)
	if !revocado {
		t.Fatalf("debería estar revocado tras RevocarJTI")
	}
}
