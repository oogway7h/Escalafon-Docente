# UT-CU25-BE-03 — Ruta del tablero

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU25-BE-03` |
| **Nombre** | Ruta del tablero |
| **Caso de uso** | CU-25-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU25-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer al Director.

## Objetivo

`GET /tablero`.

## Alcance

- `TableroDTO`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/indicadores.go`
- `backend/internal/interfaces/http/dto/indicadores.go`
- `backend/internal/interfaces/http/router_indicadores_test.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`GET /tablero`.

## Validación

1. `go test ./internal/interfaces/http -run Indicadores -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
