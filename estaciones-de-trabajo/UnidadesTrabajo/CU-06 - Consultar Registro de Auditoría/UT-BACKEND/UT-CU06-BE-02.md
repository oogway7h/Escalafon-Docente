# UT-CU06-BE-02 — Caso de uso Consultar y ListarOpciones

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU06-BE-02` |
| **Nombre** | Caso de uso Consultar y ListarOpciones |
| **Caso de uso** | CU-06-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU06-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El Director necesita encontrar rápido quién hizo qué.

## Objetivo

Consulta con filtros combinables y catálogo de filtros.

## Alcance

- Valida acción y rango; fechas de Bolivia.
- Tests con dobles.

Archivos esperados (crear/tocar):

- `backend/internal/application/auditorias/auditorias.go`
- `backend/internal/application/auditorias/auditorias_test.go`

## Fuera de alcance (no tocar)

- HTTP (UT 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Consulta con filtros combinables y catálogo de filtros.

## Validación

1. `go test ./internal/application/auditorias -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
