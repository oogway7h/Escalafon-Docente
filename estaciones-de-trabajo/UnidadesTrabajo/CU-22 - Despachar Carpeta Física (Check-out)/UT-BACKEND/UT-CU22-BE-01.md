# UT-CU22-BE-01 — TxManager y movimiento de custodia

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU22-BE-01` |
| **Nombre** | TxManager y movimiento de custodia |
| **Caso de uso** | CU-22-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

La salida y el cambio de estado deben ser atómicos.

## Objetivo

Puerto `TxManager` y tabla `movimiento_custodia` de solo anexión.

## Alcance

- Trigger que rechaza UPDATE/DELETE.

Archivos esperados (crear/tocar):

- `backend/internal/application/ports/ports.go`
- `backend/internal/infrastructure/persistence/tx.go`
- `backend/migrations/20261017120000_legajos_custodia.sql`

## Fuera de alcance (no tocar)

- Casos de uso.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Puerto `TxManager` y tabla `movimiento_custodia` de solo anexión.

## Validación

1. `go test ./internal/infrastructure/persistence -run Legajo -count=1` (con `DATABASE_URL_TEST`).
2. `make arch` pasa (domain y application sin dependencias externas).
