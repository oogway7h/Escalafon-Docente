# UT-CU21-BE-02 — Ruta de atención

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU21-BE-02` |
| **Nombre** | Ruta de atención |
| **Caso de uso** | CU-21-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU21-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer a Archivo.

## Objetivo

`POST /solicitudes-carpeta/:id/atender`.

## Alcance

- `AtenderSolicitudEntrada`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/legajos.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`POST /solicitudes-carpeta/:id/atender`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
