# UT-CU25-BE-02 — Consulta del tablero

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU25-BE-02` |
| **Nombre** | Consulta del tablero |
| **Caso de uso** | CU-25-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU25-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Consolidar en una consulta.

## Objetivo

`ConsultarTablero` y `repositorio_indicador.go`.

## Alcance

- Tests contra Postgres.

Archivos esperados (crear/tocar):

- `backend/internal/application/indicadores/tablero.go`
- `backend/internal/infrastructure/persistence/repositorio_indicador.go`
- `backend/internal/infrastructure/persistence/repositorio_indicador_test.go`

## Fuera de alcance (no tocar)

- Reportes.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`ConsultarTablero` y `repositorio_indicador.go`.

## Validación

1. `go test ./internal/application/indicadores -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
