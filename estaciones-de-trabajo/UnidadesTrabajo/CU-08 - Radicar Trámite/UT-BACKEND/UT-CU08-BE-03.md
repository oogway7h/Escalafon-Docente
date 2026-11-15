# UT-CU08-BE-03 — Handler, DTO y rutas de trámites

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU08-BE-03` |
| **Nombre** | Handler, DTO y rutas de trámites |
| **Caso de uso** | CU-08-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU08-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer radicación y lectura.

## Objetivo

Rutas de lectura (`RequiereAlguno`) y de radicación.

## Alcance

- `ResumenTramiteDTO`, `TramiteDetalleDTO`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/tramites.go`
- `backend/internal/interfaces/http/dto/tramites.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Rutas del flujo (CU-10/11).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas de lectura (`RequiereAlguno`) y de radicación.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
