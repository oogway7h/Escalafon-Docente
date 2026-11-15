# UT-CU19-FE-01 — Inventario y formulario de legajo

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU19-FE-01` |
| **Nombre** | Inventario y formulario de legajo |
| **Caso de uso** | CU-19-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Archivo mantiene el inventario.

## Objetivo

`/panel/legajos` con `LegajoForm`.

## Alcance

- Sin campo de estado.

Archivos esperados (crear/tocar):

- `frontend/app/panel/legajos/page.tsx`
- `frontend/app/panel/legajos/LegajoForm.tsx`
- `frontend/app/panel/legajos/actions.ts`
- `frontend/lib/api-legajos.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/legajos` con `LegajoForm`.

## Validación

1. Con `pnpm dev` y la API levantada, crear un legajo sin código lo muestra con `LEG-…`.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
