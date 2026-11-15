# UT-CU07-BE-03 — Casos de uso de categorías

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU07-BE-03` |
| **Nombre** | Casos de uso de categorías |
| **Caso de uso** | CU-07-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU07-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El catálogo define qué documentos pide cada trámite.

## Objetivo

`Listar`, `Obtener`, `Guardar`, `CambiarEstado`, `Eliminar` con conteos de uso.

## Alcance

- Clave normalizada; requisitos ordenados.

Archivos esperados (crear/tocar):

- `backend/internal/application/categorias/categorias.go`
- `backend/internal/application/categorias/categorias_test.go`
- `backend/internal/infrastructure/persistence/claves.go`

## Fuera de alcance (no tocar)

- Docentes.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Listar`, `Obtener`, `Guardar`, `CambiarEstado`, `Eliminar` con conteos de uso.

## Validación

1. `go test ./internal/application/categorias -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
