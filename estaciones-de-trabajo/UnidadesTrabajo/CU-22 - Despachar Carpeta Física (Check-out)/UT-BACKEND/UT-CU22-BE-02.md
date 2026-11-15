# UT-CU22-BE-02 — Casos de uso Despachar y Entregar

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU22-BE-02` |
| **Nombre** | Casos de uso Despachar y Entregar |
| **Caso de uso** | CU-22-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU22-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Registrar la salida con todos sus datos.

## Objetivo

`Despachar` y `Entregar` (aprobar + despachar).

## Alcance

- Validar antes de aprobar.
- Tests de concurrencia en `repositorio_legajo_test.go`.

Archivos esperados (crear/tocar):

- `backend/internal/application/legajos/legajos.go`
- `backend/internal/application/legajos/legajos_test.go`
- `backend/internal/infrastructure/persistence/repositorio_legajo.go`
- `backend/internal/infrastructure/persistence/repositorio_legajo_test.go`

## Fuera de alcance (no tocar)

- Devolución.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Despachar` y `Entregar` (aprobar + despachar).

## Validación

1. `go test ./internal/application/legajos -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
