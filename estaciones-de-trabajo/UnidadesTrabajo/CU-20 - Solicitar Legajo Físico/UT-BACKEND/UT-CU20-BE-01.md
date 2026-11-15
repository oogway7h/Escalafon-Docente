# UT-CU20-BE-01 — Caso de uso Solicitar y ListarSolicitudes

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU20-BE-01` |
| **Nombre** | Caso de uso Solicitar y ListarSolicitudes |
| **Caso de uso** | CU-20-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El analista pide la carpeta sin ir al archivo.

## Objetivo

`Solicitar` con sus validaciones y listado filtrado por solicitante.

## Alcance

- Índice `ux_solicitud_abierta` manual.

Archivos esperados (crear/tocar):

- `backend/internal/application/legajos/legajos.go`
- `backend/internal/application/legajos/legajos_test.go`
- `backend/internal/infrastructure/persistence/repositorio_legajo.go`

## Fuera de alcance (no tocar)

- Atender.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Solicitar` con sus validaciones y listado filtrado por solicitante.

## Validación

1. `go test ./internal/application/legajos -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
