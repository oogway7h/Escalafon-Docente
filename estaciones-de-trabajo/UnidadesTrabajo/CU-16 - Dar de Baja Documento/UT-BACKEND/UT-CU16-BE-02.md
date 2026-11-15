# UT-CU16-BE-02 — Rutas de baja

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU16-BE-02` |
| **Nombre** | Rutas de baja |
| **Caso de uso** | CU-16-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU16-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer la baja.

## Objetivo

Rutas PATCH de estado.

## Alcance

- `BajaDocumentoEntrada`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/tramites.go`
- `backend/internal/interfaces/http/handlers/expediente.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas PATCH de estado.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
