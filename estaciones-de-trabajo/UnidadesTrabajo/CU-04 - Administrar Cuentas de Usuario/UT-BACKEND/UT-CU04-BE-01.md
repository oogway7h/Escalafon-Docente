# UT-CU04-BE-01 — Repositorio de usuarios y unicidad de correo

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU04-BE-01` |
| **Nombre** | Repositorio de usuarios y unicidad de correo |
| **Caso de uso** | CU-04-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Dos cuentas con el mismo correo romperían el login y la recuperación.

## Objetivo

Repositorio GORM con búsqueda y el índice único sin mayúsculas.

## Alcance

- Modelo `usuario`; `ux_usuario_correo` en `ajustes_manuales.sql`.

Archivos esperados (crear/tocar):

- `backend/internal/infrastructure/persistence/repositorio_usuario.go`
- `backend/internal/infrastructure/persistence/modelos.go`
- `backend/migrations/20260922003354_ajustes_manuales.sql`

## Fuera de alcance (no tocar)

- Casos de uso (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Repositorio GORM con búsqueda y el índice único sin mayúsculas.

## Validación

1. `go test ./internal/infrastructure/persistence -run Usuario -count=1` (con `DATABASE_URL_TEST`).
2. `make arch` pasa (domain y application sin dependencias externas).
