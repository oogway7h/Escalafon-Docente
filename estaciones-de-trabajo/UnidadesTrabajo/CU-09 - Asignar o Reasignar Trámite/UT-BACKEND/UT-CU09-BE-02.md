# UT-CU09-BE-02 — Rutas de asignación

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU09-BE-02` |
| **Nombre** | Rutas de asignación |
| **Caso de uso** | CU-09-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU09-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer la bandeja al Director.

## Objetivo

Rutas `/analistas` y `/tramites/:id/asignar`.

## Alcance

- `AsignarEntrada`, `AnalistaDTO`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/tramites.go`
- `backend/internal/interfaces/http/dto/tramites.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/analistas` y `/tramites/:id/asignar`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
