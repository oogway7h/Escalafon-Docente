# UT-CU24-FE-01 — Página de bitácora

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU24-FE-01` |
| **Nombre** | Página de bitácora |
| **Caso de uso** | CU-24-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Ver el historial de una carpeta.

## Objetivo

`/panel/custodia/bitacora`.

## Alcance

- El menú resalta solo el ítem más específico.

Archivos esperados (crear/tocar):

- `frontend/app/panel/custodia/bitacora/page.tsx`
- `frontend/lib/api-legajos.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/custodia/bitacora`.

## Validación

1. Con `pnpm dev` y la API levantada, filtrar por un `LEG-…` muestra su salida y su entrada.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
