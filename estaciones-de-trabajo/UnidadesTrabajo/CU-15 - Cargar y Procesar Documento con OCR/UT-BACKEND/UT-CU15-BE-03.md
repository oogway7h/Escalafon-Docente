# UT-CU15-BE-03 — Casos de uso PrepararSubida, ConfirmarSubida y CorregirPaginas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-BE-03` |
| **Nombre** | Casos de uso PrepararSubida, ConfirmarSubida y CorregirPaginas |
| **Caso de uso** | CU-15-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU15-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Validar antes de subir y procesar al confirmar.

## Objetivo

Flujo SUBIENDO → PROCESADO y corrección por página.

## Alcance

- Límites por tipo; analista asignado.
- Tests con dobles de almacenamiento y procesador.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/documentos.go`
- `backend/internal/application/tramites/documentos_test.go`

## Fuera de alcance (no tocar)

- Expediente (UT 04).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Flujo SUBIENDO → PROCESADO y corrección por página.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
