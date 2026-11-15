# UT-CU01-BE-05 — Handler y rutas de autenticación

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-BE-05` |
| **Nombre** | Handler y rutas de autenticación |
| **Caso de uso** | CU-01-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU01-BE-03, UT-CU01-BE-04 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer el login al frontend con errores estables y límite de tasa.

## Objetivo

`POST /auth/login`, `/auth/logout`, `GET /auth/sesion` registrados con su limitador.

## Alcance

- `AuthHandler.Login/Logout/Sesion` traducen DTO ↔ caso de uso.
- `apierror.Mapear` para los errores de dominio.
- Limitador de 10/min por IP en `/auth/*` (Redis).
- Tests de router: login correcto, credenciales inválidas, validación.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/auth.go`
- `backend/internal/interfaces/http/dto/auth.go`
- `backend/internal/interfaces/http/router.go`
- `backend/internal/interfaces/http/router_test.go`

## Fuera de alcance (no tocar)

- Rutas de recuperación (CU-02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`POST /auth/login`, `/auth/logout`, `GET /auth/sesion` registrados con su limitador.

## Validación

1. `go test ./internal/interfaces/http -run Login -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
