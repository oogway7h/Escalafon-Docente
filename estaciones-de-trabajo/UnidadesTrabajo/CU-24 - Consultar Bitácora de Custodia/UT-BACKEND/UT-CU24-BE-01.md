# UT-CU24-BE-01 — Caso de uso ListarMovimientos y ruta

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU24-BE-01` |
| **Nombre** | Caso de uso ListarMovimientos y ruta |
| **Caso de uso** | CU-24-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Auditar la cadena de custodia.

## Objetivo

Consulta con filtros y ruta `/custodia/bitacora`.

## Alcance

- `MovimientoDTO` con `vencido`.

Archivos esperados (crear/tocar):

- `backend/internal/application/legajos/legajos.go`
- `backend/internal/interfaces/http/handlers/legajos.go`
- `backend/internal/interfaces/http/dto/legajos.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Consulta con filtros y ruta `/custodia/bitacora`.

## Validación

1. `go test ./internal/application/legajos -count=1` en verde.
2. `go test ./internal/interfaces/http -count=1` en verde.
3. `make arch` pasa (domain y application sin dependencias externas).
