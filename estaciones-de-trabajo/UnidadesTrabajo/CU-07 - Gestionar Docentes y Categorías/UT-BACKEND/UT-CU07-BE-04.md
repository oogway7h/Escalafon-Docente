# UT-CU07-BE-04 — Handlers y rutas de catálogo

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU07-BE-04` |
| **Nombre** | Handlers y rutas de catálogo |
| **Caso de uso** | CU-07-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU07-BE-02, UT-CU07-BE-03 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer padrón y catálogo.

## Objetivo

Rutas `/docentes`, `/facultades`, `/categorias` con lecturas compartidas por `RequiereAlguno`.

## Alcance

- DTOs con requisitos y conteos.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/catalogo.go`
- `backend/internal/interfaces/http/dto/tramites.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Trámites.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/docentes`, `/facultades`, `/categorias` con lecturas compartidas por `RequiereAlguno`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
