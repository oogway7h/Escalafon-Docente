# UT-CU10-BE-02 — Rutas de revisión

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU10-BE-02` |
| **Nombre** | Rutas de revisión |
| **Caso de uso** | CU-10-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU10-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer las acciones del analista.

## Objetivo

Rutas `/revision`, `/observaciones`, `/resolver` con `TRAMITE_REVISAR`.

## Alcance

- `ObservacionEntrada`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/tramites.go`
- `backend/internal/interfaces/http/dto/tramites.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/revision`, `/observaciones`, `/resolver` con `TRAMITE_REVISAR`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
