# UT-CU25-BE-01 — Dominio indicador: mora y tiempos

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU25-BE-01` |
| **Nombre** | Dominio indicador: mora y tiempos |
| **Caso de uso** | CU-25-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Reglas de mora y agrupación probadas sin base de datos.

## Objetivo

`indicador.EvaluarMora` y agrupación de tiempos.

## Alcance

- Días calendario de Bolivia.

Archivos esperados (crear/tocar):

- `backend/internal/domain/indicador/indicador.go`
- `backend/internal/domain/indicador/indicador_test.go`

## Fuera de alcance (no tocar)

- Consultas SQL.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`indicador.EvaluarMora` y agrupación de tiempos.

## Validación

1. `go test ./internal/domain/indicador -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
