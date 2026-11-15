# UT-CU20-BE-02 — Rutas de solicitud

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU20-BE-02` |
| **Nombre** | Rutas de solicitud |
| **Caso de uso** | CU-20-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU20-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer al analista.

## Objetivo

Rutas `POST /legajos/:id/solicitudes` y `GET /solicitudes-carpeta`.

## Alcance

- `SolicitarCarpetaEntrada`, `SolicitudCarpetaDTO`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/legajos.go`
- `backend/internal/interfaces/http/dto/legajos.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `POST /legajos/:id/solicitudes` y `GET /solicitudes-carpeta`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
