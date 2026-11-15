# UT-CU03-BE-02 — Handler y rutas de perfil

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU03-BE-02` |
| **Nombre** | Handler y rutas de perfil |
| **Caso de uso** | CU-03-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU03-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer el perfil al panel.

## Objetivo

Rutas `/perfil` con `PERFIL_GESTIONAR`.

## Alcance

- DTO de entrada/salida; mapeo de errores.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/perfil.go`
- `backend/internal/interfaces/http/dto/perfil.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Usuarios (CU-04).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/perfil` con `PERFIL_GESTIONAR`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
