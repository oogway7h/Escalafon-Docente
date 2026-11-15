# UT-CU05-BE-01 — Dominio de rol, permisos y seed

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU05-BE-01` |
| **Nombre** | Dominio de rol, permisos y seed |
| **Caso de uso** | CU-05-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Los permisos son la base de toda la autorización.

## Objetivo

`domain/rol` con constantes `Permiso*`, roles de sistema y seed idempotente.

## Alcance

- `PermisosProtegidosDirector`.
- `cmd/seed` carga catálogo, roles y usuarios de prueba.

Archivos esperados (crear/tocar):

- `backend/internal/domain/rol/rol.go`
- `backend/internal/domain/rol/errores.go`
- `backend/cmd/seed/main.go`

## Fuera de alcance (no tocar)

- Casos de uso (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`domain/rol` con constantes `Permiso*`, roles de sistema y seed idempotente.

## Validación

1. `make seed` dos veces no duplica filas.
2. `make arch` pasa (domain y application sin dependencias externas).
