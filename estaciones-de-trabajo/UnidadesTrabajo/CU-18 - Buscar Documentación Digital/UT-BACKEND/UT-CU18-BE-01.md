# UT-CU18-BE-01 — Caso de uso BuscarDocumentos

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU18-BE-01` |
| **Nombre** | Caso de uso BuscarDocumentos |
| **Caso de uso** | CU-18-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Encontrar antecedentes sin abrir decenas de PDF.

## Objetivo

Búsqueda por página con relevancia y fragmento resaltado.

## Alcance

- Agrupa por documento; filtra `activo`.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/busqueda.go`
- `backend/internal/infrastructure/persistence/repositorio_documentos.go`

## Fuera de alcance (no tocar)

- Rutas.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Búsqueda por página con relevancia y fragmento resaltado.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
