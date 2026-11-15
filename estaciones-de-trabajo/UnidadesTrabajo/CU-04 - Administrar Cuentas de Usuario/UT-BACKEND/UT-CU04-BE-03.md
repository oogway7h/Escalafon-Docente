# UT-CU04-BE-03 — Handler y rutas de usuarios

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU04-BE-03` |
| **Nombre** | Handler y rutas de usuarios |
| **Caso de uso** | CU-04-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU04-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer la gestión al panel.

## Objetivo

Rutas `/usuarios` con `USUARIO_GESTIONAR`.

## Alcance

- DTOs y mapeo de errores.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/usuarios.go`
- `backend/internal/interfaces/http/dto/usuarios.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/usuarios` con `USUARIO_GESTIONAR`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
