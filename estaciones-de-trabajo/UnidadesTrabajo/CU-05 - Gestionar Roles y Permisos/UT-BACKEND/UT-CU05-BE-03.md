# UT-CU05-BE-03 — Handler y rutas de roles

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU05-BE-03` |
| **Nombre** | Handler y rutas de roles |
| **Caso de uso** | CU-05-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU05-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer la gestión al panel.

## Objetivo

Rutas `/roles` y `/permisos`.

## Alcance

- DTOs y errores.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/roles.go`
- `backend/internal/interfaces/http/dto/roles.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/roles` y `/permisos`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
