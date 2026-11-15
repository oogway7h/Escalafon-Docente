# UT-CU09-BE-01 — Casos de uso Asignar y ListarAnalistas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU09-BE-01` |
| **Nombre** | Casos de uso Asignar y ListarAnalistas |
| **Caso de uso** | CU-09-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Repartir el trabajo con información de carga.

## Objetivo

`Asignar` (asignar/reasignar) y `ListarAnalistas` con carga vigente.

## Alcance

- Motivo obligatorio al reasignar.
- Tests: asignación, reasignación sin motivo, remitido.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/tramites_test.go`
- `backend/internal/infrastructure/persistence/repositorio_tramite.go`

## Fuera de alcance (no tocar)

- Revisión.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Asignar` (asignar/reasignar) y `ListarAnalistas` con carga vigente.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
