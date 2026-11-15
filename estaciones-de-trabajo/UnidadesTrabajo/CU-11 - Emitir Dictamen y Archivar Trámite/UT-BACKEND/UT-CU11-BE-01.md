# UT-CU11-BE-01 — Caso de uso Dictaminar

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU11-BE-01` |
| **Nombre** | Caso de uso Dictaminar |
| **Caso de uso** | CU-11-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El dictamen cierra la revisión técnica.

## Objetivo

`Dictaminar` y `RegistrarPeriodo` con sus precondiciones.

## Alcance

- Sin observaciones pendientes; periodo registrado.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/flujo.go`
- `backend/internal/application/tramites/tramites_test.go`

## Fuera de alcance (no tocar)

- Remisiones (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Dictaminar` y `RegistrarPeriodo` con sus precondiciones.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
