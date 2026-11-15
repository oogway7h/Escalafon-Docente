# UT-CU14-BE-02 — Caso de uso Notificar con reintentos

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU14-BE-02` |
| **Nombre** | Caso de uso Notificar con reintentos |
| **Caso de uso** | CU-14-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU14-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El docente se entera sin llamar al Departamento.

## Objetivo

`Notificar.Transicion` y `ReintentarPendientes` (ticker en `main.go`).

## Alcance

- Solo en las cuatro transiciones.
- Máx. 3 intentos.

Archivos esperados (crear/tocar):

- `backend/internal/application/notificaciones/notificaciones.go`
- `backend/internal/application/notificaciones/notificaciones_test.go`
- `backend/cmd/api/main.go`

## Fuera de alcance (no tocar)

- Consulta (UT 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Notificar.Transicion` y `ReintentarPendientes` (ticker en `main.go`).

## Validación

1. `go test ./internal/application/notificaciones -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
