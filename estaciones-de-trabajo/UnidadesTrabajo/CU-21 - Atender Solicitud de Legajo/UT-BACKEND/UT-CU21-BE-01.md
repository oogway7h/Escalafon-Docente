# UT-CU21-BE-01 — Caso de uso Atender

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU21-BE-01` |
| **Nombre** | Caso de uso Atender |
| **Caso de uso** | CU-21-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Archivo responde con trazabilidad.

## Objetivo

`Atender` con aprobación o rechazo.

## Alcance

- Tests: rechazo sin nota, estado inválido.

Archivos esperados (crear/tocar):

- `backend/internal/application/legajos/legajos.go`
- `backend/internal/application/legajos/legajos_test.go`

## Fuera de alcance (no tocar)

- Despachar.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Atender` con aprobación o rechazo.

## Validación

1. `go test ./internal/application/legajos -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
