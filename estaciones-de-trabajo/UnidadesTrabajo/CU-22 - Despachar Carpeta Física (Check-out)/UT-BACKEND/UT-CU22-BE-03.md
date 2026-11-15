# UT-CU22-BE-03 — Rutas de despacho

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU22-BE-03` |
| **Nombre** | Rutas de despacho |
| **Caso de uso** | CU-22-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU22-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer a Archivo.

## Objetivo

Rutas `/despachar` y `/entregar`.

## Alcance

- `DespacharEntrada`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/legajos.go`
- `backend/internal/interfaces/http/dto/legajos.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/despachar` y `/entregar`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
