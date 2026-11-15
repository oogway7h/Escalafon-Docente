# UT-CU07-BE-01 — Dominio y persistencia de docente, facultad y categoría

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU07-BE-01` |
| **Nombre** | Dominio y persistencia de docente, facultad y categoría |
| **Caso de uso** | CU-07-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El padrón es la base de trámites y legajos.

## Objetivo

Entidades y repositorios con búsqueda exacta/fragmento.

## Alcance

- `domain/docente`, `domain/categoria`.
- `tieneDigitos` decide búsqueda exacta.
- Migración `padron`.

Archivos esperados (crear/tocar):

- `backend/internal/domain/docente/docente.go`
- `backend/internal/domain/categoria/categoria.go`
- `backend/internal/infrastructure/persistence/repositorio_docente.go`
- `backend/internal/infrastructure/persistence/repositorio_categoria.go`

## Fuera de alcance (no tocar)

- Casos de uso (UT 02, 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Entidades y repositorios con búsqueda exacta/fragmento.

## Validación

1. `go test ./internal/infrastructure/persistence -count=1` (con `DATABASE_URL_TEST`).
2. `make arch` pasa (domain y application sin dependencias externas).
