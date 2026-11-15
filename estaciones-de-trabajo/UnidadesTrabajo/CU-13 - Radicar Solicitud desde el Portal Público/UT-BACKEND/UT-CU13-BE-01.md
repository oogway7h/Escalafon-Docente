# UT-CU13-BE-01 — Puerto ObjectStorage y adaptador S3

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU13-BE-01` |
| **Nombre** | Puerto ObjectStorage y adaptador S3 |
| **Caso de uso** | CU-13-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Los adjuntos no deben vivir en la base de datos.

## Objetivo

`ports.ObjectStorage` con implementación S3 y URLs firmadas.

## Alcance

- Sin `S3_BUCKET` responde `ALMACEN_NO_DISPONIBLE`.

Archivos esperados (crear/tocar):

- `backend/internal/application/ports/ports.go`
- `backend/internal/infrastructure/storage/s3.go`

## Fuera de alcance (no tocar)

- Caso de uso (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`ports.ObjectStorage` con implementación S3 y URLs firmadas.

## Validación

1. La API arranca sin `S3_BUCKET` y la subida responde 503.
2. `make arch` pasa (domain y application sin dependencias externas).
