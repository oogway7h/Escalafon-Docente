# UT-CU26-BE-01 — Caso de uso GenerarReporte

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU26-BE-01` |
| **Nombre** | Caso de uso GenerarReporte |
| **Caso de uso** | CU-26-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Reportes oficiales sin armar planillas a mano.

## Objetivo

`GenerarReporte` con `indicador.Reporte` independiente del formato.

## Alcance

- Auditoría solo si hay filas.

Archivos esperados (crear/tocar):

- `backend/internal/application/indicadores/reportes.go`
- `backend/internal/application/indicadores/indicadores_test.go`

## Fuera de alcance (no tocar)

- Composición.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`GenerarReporte` con `indicador.Reporte` independiente del formato.

## Validación

1. `go test ./internal/application/indicadores -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
