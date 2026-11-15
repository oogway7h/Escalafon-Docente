# UT-CU14-BE-03 — Ruta de consulta de envíos

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU14-BE-03` |
| **Nombre** | Ruta de consulta de envíos |
| **Caso de uso** | CU-14-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU14-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El Director verifica que los avisos salen.

## Objetivo

`GET /notificaciones`.

## Alcance

- Paginada.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/notificaciones.go`
- `backend/internal/interfaces/http/dto/notificaciones.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`GET /notificaciones`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
