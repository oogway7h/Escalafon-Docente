# UT-CU16-BE-01 — Migración y caso de uso DarDeBajaDocumento

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU16-BE-01` |
| **Nombre** | Migración y caso de uso DarDeBajaDocumento |
| **Caso de uso** | CU-16-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Corregir errores de carga sin perder evidencia.

## Objetivo

Columnas de baja y `DarDeBajaDocumento` / `DarDeBajaExpediente`.

## Alcance

- Auditoría con motivo.

Archivos esperados (crear/tocar):

- `backend/migrations/20261019120000_documento_baja.sql`
- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/expediente.go`
- `backend/internal/application/tramites/tramites_test.go`

## Fuera de alcance (no tocar)

- Restitución.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Columnas de baja y `DarDeBajaDocumento` / `DarDeBajaExpediente`.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
