# UT-CU17-BE-01 — Casos de uso ListarDocumentosDeBaja y RestituirDocumento

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU17-BE-01` |
| **Nombre** | Casos de uso ListarDocumentosDeBaja y RestituirDocumento |
| **Caso de uso** | CU-17-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Deshacer una baja equivocada con trazabilidad.

## Objetivo

Restitución con UPDATE condicional y auditoría.

## Alcance

- Tests con dobles.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/restitucion.go`
- `backend/internal/application/tramites/restitucion_test.go`
- `backend/internal/infrastructure/persistence/repositorio_documentos.go`

## Fuera de alcance (no tocar)

- Búsqueda.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Restitución con UPDATE condicional y auditoría.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
