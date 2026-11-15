# UT-CU06-BE-03 — Handler y rutas de auditoría

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU06-BE-03` |
| **Nombre** | Handler y rutas de auditoría |
| **Caso de uso** | CU-06-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU06-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer la consulta.

## Objetivo

Rutas `/auditoria` y `/auditoria/filtros`.

## Alcance

- Sin rutas de escritura.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/auditoria.go`
- `backend/internal/interfaces/http/dto/auditoria.go`
- `backend/internal/interfaces/http/router_auditoria_test.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/auditoria` y `/auditoria/filtros`.

## Validación

1. `go test ./internal/interfaces/http -run Auditoria -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
