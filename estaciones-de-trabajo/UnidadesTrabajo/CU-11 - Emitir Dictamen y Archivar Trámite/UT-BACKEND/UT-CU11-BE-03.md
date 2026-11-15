# UT-CU11-BE-03 — Rutas del flujo

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU11-BE-03` |
| **Nombre** | Rutas del flujo |
| **Caso de uso** | CU-11-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU11-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer cada transición con su permiso.

## Objetivo

Rutas de periodo, dictamen, remisiones, regreso y archivar.

## Alcance

- Regreso con `RequiereAlguno`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/tramites.go`
- `backend/internal/interfaces/http/dto/tramites.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas de periodo, dictamen, remisiones, regreso y archivar.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
