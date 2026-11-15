# UT-CU19-BE-02 — Casos de uso ListarLegajos, ObtenerLegajo y GuardarLegajo

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU19-BE-02` |
| **Nombre** | Casos de uso ListarLegajos, ObtenerLegajo y GuardarLegajo |
| **Caso de uso** | CU-19-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU19-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Mantener el inventario.

## Objetivo

Listado paginado, detalle y guardado.

## Alcance

- Búsqueda exacta por código.

Archivos esperados (crear/tocar):

- `backend/internal/application/legajos/legajos.go`
- `backend/internal/application/legajos/legajos_test.go`
- `backend/internal/infrastructure/persistence/repositorio_legajo.go`

## Fuera de alcance (no tocar)

- Custodia.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Listado paginado, detalle y guardado.

## Validación

1. `go test ./internal/application/legajos -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
