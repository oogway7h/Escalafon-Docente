# UT-CU02-BE-03 — Handlers y rutas públicas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU02-BE-03` |
| **Nombre** | Handlers y rutas públicas |
| **Caso de uso** | CU-02-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU02-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El enlace del correo debe poder usarse sin sesión.

## Objetivo

Rutas `/auth/recuperacion`, `/restablecer`, `/activar` con el limitador de `/auth`.

## Alcance

- DTOs con validación.
- Tests de router: token inexistente, política débil, activación doble.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/auth.go`
- `backend/internal/interfaces/http/dto/auth.go`
- `backend/internal/interfaces/http/router_test.go`

## Fuera de alcance (no tocar)

- Login (CU-01).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/auth/recuperacion`, `/restablecer`, `/activar` con el limitador de `/auth`.

## Validación

1. `go test ./internal/interfaces/http -run Restablecer -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
