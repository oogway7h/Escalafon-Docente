package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/usuario"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// repoContador envuelve un usuario.Repository y cuenta las llamadas a BuscarPorID, para
// demostrar que, con la caché vigente, validar una sesión no consulta la base de datos
// (sección 12, criterios técnicos).
type repoContador struct {
	usuario.Repository
	Llamadas int
}

func (r *repoContador) BuscarPorID(ctx context.Context, id string) (*usuario.Usuario, error) {
	r.Llamadas++
	return r.Repository.BuscarPorID(ctx, id)
}

func TestAutenticar_ConCacheVigenteNoConsultaLaBaseDeDatos(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	u := usuarios.Sembrar(usuario.Usuario{RolID: 1, VersionSesion: 1, Activo: true})
	contador := &repoContador{Repository: usuarios}

	tokens := testdobles.NewTokenServiceFalso()
	cache := testdobles.NewSessionCacheFalso()
	tok, err := tokens.Emitir(ports.ClaimsSesion{UsuarioID: u.ID, RolID: 1, Version: 1, JTI: "jti-1"}, 0)
	if err != nil {
		t.Fatalf("no se pudo emitir el token de prueba: %v", err)
	}

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ManejadorErrores})
	app.Get("/probar", middleware.Autenticar(tokens, cache, contador), func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	hacerPeticion := func() int {
		req := httptest.NewRequest(http.MethodGet, "/probar", nil)
		req.Header.Set("Authorization", "Bearer "+tok)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("no esperaba error: %v", err)
		}
		return resp.StatusCode
	}

	if status := hacerPeticion(); status != fiber.StatusOK {
		t.Fatalf("primera petición: esperaba 200, obtuve %d", status)
	}
	if contador.Llamadas != 1 {
		t.Fatalf("la primera petición (fallo de caché) debería consultar la BD una vez, llamadas=%d", contador.Llamadas)
	}

	if status := hacerPeticion(); status != fiber.StatusOK {
		t.Fatalf("segunda petición: esperaba 200, obtuve %d", status)
	}
	if contador.Llamadas != 1 {
		t.Fatalf("con la caché vigente, la segunda petición no debería consultar la BD; llamadas=%d", contador.Llamadas)
	}
}

func TestAutenticar_RechazaSinEncabezado(t *testing.T) {
	app := fiber.New(fiber.Config{ErrorHandler: middleware.ManejadorErrores})
	app.Get("/probar", middleware.Autenticar(testdobles.NewTokenServiceFalso(), testdobles.NewSessionCacheFalso(), testdobles.NewUsuarioRepoFalso()),
		func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	resp, _ := app.Test(httptest.NewRequest(http.MethodGet, "/probar", nil))
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("esperaba 401 sin encabezado, obtuve %d", resp.StatusCode)
	}
}

func TestAutenticar_RechazaVersionDesactualizada(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	u := usuarios.Sembrar(usuario.Usuario{RolID: 1, VersionSesion: 2, Activo: true}) // versión actual = 2
	tokens := testdobles.NewTokenServiceFalso()
	cache := testdobles.NewSessionCacheFalso()
	tok, _ := tokens.Emitir(ports.ClaimsSesion{UsuarioID: u.ID, RolID: 1, Version: 1, JTI: "jti-2"}, 0) // token con versión vieja

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ManejadorErrores})
	app.Get("/probar", middleware.Autenticar(tokens, cache, usuarios), func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/probar", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("esperaba 401 por versión de sesión desactualizada, obtuve %d", resp.StatusCode)
	}
}

// cacheCaida simula Redis sin servicio al consultar la lista de revocación.
type cacheCaida struct{ *testdobles.SessionCacheFalso }

func (cacheCaida) EstaRevocado(context.Context, string) (bool, error) {
	return false, errors.New("redis caído")
}

func TestAutenticar_FallaCerradoSiNoPuedeComprobarLaRevocacion(t *testing.T) {
	usuarios := testdobles.NewUsuarioRepoFalso()
	u := usuarios.Sembrar(usuario.Usuario{RolID: 1, VersionSesion: 1, Activo: true})
	tokens := testdobles.NewTokenServiceFalso()
	tok, _ := tokens.Emitir(ports.ClaimsSesion{UsuarioID: u.ID, RolID: 1, Version: 1, JTI: "jti-1"}, 0)

	app := fiber.New(fiber.Config{ErrorHandler: middleware.ManejadorErrores})
	app.Get("/probar", middleware.Autenticar(tokens, cacheCaida{testdobles.NewSessionCacheFalso()}, usuarios),
		func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/probar", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusServiceUnavailable {
		t.Fatalf("esperaba 503 con la lista de revocación inaccesible, obtuve %d", resp.StatusCode)
	}
}
