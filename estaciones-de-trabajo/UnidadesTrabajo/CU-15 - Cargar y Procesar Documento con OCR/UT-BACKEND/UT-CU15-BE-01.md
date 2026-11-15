# UT-CU15-BE-01 — Migraciones de OCR y páginas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-BE-01` |
| **Nombre** | Migraciones de OCR y páginas |
| **Caso de uso** | CU-15-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El texto debe guardarse por página para saltar a la página exacta.

## Objetivo

Columnas OCR en `tramite_documento` y tabla `documento_pagina` con índice `es_unaccent`.

## Alcance

- Índices manuales (Atlas querrá borrarlos).

Archivos esperados (crear/tocar):

- `backend/migrations/20261020100000_documento_ocr.sql`
- `backend/migrations/20261021100000_documento_paginas.sql`
- `backend/migrations/20261021110000_busqueda_sin_tildes.sql`

## Fuera de alcance (no tocar)

- Casos de uso.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Columnas OCR en `tramite_documento` y tabla `documento_pagina` con índice `es_unaccent`.

## Validación

1. `make migrate-status` sin pendientes.
2. `make arch` pasa (domain y application sin dependencias externas).
