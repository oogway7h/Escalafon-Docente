# UT-CU07-BE-02 — Casos de uso de docentes

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU07-BE-02` |
| **Nombre** | Casos de uso de docentes |
| **Caso de uso** | CU-07-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU07-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Mantener el padrón sin duplicados.

## Objetivo

`Listar`, `ListarPagina`, `Obtener`, `Guardar`, `CambiarEstado`, `ListarFacultades`.

## Alcance

- Unicidad de código y cédula, inmutabilidad de cédula, auditoría.

Archivos esperados (crear/tocar):

- `backend/internal/application/docentes/docentes.go`
- `backend/internal/application/docentes/docentes_test.go`

## Fuera de alcance (no tocar)

- Categorías (UT 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Listar`, `ListarPagina`, `Obtener`, `Guardar`, `CambiarEstado`, `ListarFacultades`.

## Validación

1. `go test ./internal/application/docentes -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
