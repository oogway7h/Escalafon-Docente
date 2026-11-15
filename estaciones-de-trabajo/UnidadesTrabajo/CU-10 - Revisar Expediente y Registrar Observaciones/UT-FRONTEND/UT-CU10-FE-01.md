# UT-CU10-FE-01 — Mi bandeja

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU10-FE-01` |
| **Nombre** | Mi bandeja |
| **Caso de uso** | CU-10-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Cada analista ve lo suyo.

## Objetivo

`/panel/tramites/revisar`.

## Alcance

- Insignia de pendientes en la portada.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/revisar/page.tsx`
- `frontend/app/panel/page.tsx`

## Fuera de alcance (no tocar)

- Detalle.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/tramites/revisar`.

## Validación

1. Con `pnpm dev` y la API levantada, `analista@ded.dev` ve solo sus trámites.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
