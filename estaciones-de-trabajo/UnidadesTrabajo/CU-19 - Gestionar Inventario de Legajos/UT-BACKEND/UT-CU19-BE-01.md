# UT-CU19-BE-01 — Dominio legajo, migración y trigger de código

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU19-BE-01` |
| **Nombre** | Dominio legajo, migración y trigger de código |
| **Caso de uso** | CU-19-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Cada carpeta debe ubicarse sin buscar a mano.

## Objetivo

`domain/legajo`, tabla con código y estado, trigger `fn_legajo_codigo`.

## Alcance

- Estados EN_ARCHIVO / PRESTADO.

Archivos esperados (crear/tocar):

- `backend/internal/domain/legajo/legajo.go`
- `backend/migrations/20261017120000_legajos_custodia.sql`

## Fuera de alcance (no tocar)

- Casos de uso.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`domain/legajo`, tabla con código y estado, trigger `fn_legajo_codigo`.

## Validación

1. Insertar un legajo sin código en la base asigna `LEG-…`.
2. `make arch` pasa (domain y application sin dependencias externas).
