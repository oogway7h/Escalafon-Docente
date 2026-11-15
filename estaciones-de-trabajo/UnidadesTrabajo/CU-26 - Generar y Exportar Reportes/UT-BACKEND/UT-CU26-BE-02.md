# UT-CU26-BE-02 — Exportadores PDF y Excel

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU26-BE-02` |
| **Nombre** | Exportadores PDF y Excel |
| **Caso de uso** | CU-26-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU26-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Membrete institucional y planilla filtrable.

## Objetivo

`ports.ExportadorReportes` con `go-pdf/fpdf` y `excelize`.

## Alcance

- Tests de composición real.

Archivos esperados (crear/tocar):

- `backend/internal/infrastructure/reportes/pdf.go`
- `backend/internal/infrastructure/reportes/excel.go`
- `backend/internal/infrastructure/reportes/reportes_test.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`ports.ExportadorReportes` con `go-pdf/fpdf` y `excelize`.

## Validación

1. `go test ./internal/infrastructure/reportes -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
