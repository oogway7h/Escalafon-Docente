# UT-CU15-BE-02 — Puerto ProcesadorDocumentos y cliente ia-service

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-BE-02` |
| **Nombre** | Puerto ProcesadorDocumentos y cliente ia-service |
| **Caso de uso** | CU-15-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU15-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El OCR es pesado y no debe bloquear la API.

## Objetivo

Cliente HTTP a `ia-service` con token y reanudación.

## Alcance

- Descarga en streaming de S3; hash del original.

Archivos esperados (crear/tocar):

- `backend/internal/application/ports/ports.go`
- `backend/internal/infrastructure/iaservice/cliente.go`
- `backend/../ia-service/app`

## Fuera de alcance (no tocar)

- Casos de uso (UT 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Cliente HTTP a `ia-service` con token y reanudación.

## Validación

1. `docker compose up -d ia-service` y `GET /health` del servicio responde 200.
2. `make arch` pasa (domain y application sin dependencias externas).
