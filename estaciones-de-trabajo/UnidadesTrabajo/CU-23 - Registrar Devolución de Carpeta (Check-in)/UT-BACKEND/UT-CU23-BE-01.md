# UT-CU23-BE-01 — Caso de uso Devolver

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU23-BE-01` |
| **Nombre** | Caso de uso Devolver |
| **Caso de uso** | CU-23-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Cerrar el préstamo con responsable.

## Objetivo

`Devolver` con validación de observación.

## Alcance

- Transacción con UPDATE condicional.

Archivos esperados (crear/tocar):

- `backend/internal/application/legajos/legajos.go`
- `backend/internal/application/legajos/legajos_test.go`
- `backend/migrations/20261021120000_quitar_evidencia_custodia.sql`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Devolver` con validación de observación.

## Validación

1. `go test ./internal/application/legajos -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
