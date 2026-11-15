# UT-CU12-FE-01 — Página de seguimiento

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU12-FE-01` |
| **Nombre** | Página de seguimiento |
| **Caso de uso** | CU-12-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Consulta sin cuenta, también sin JavaScript.

## Objetivo

`/seguimiento` con `BuscadorSeguimiento` y `FasesTramite`.

## Alcance

- Código > 5 = seguimiento; ≤ 5 = docente + cédula (POST).
- Nota de privacidad.

Archivos esperados (crear/tocar):

- `frontend/app/(publico)/seguimiento/page.tsx`
- `frontend/app/(publico)/seguimiento/BuscadorSeguimiento.tsx`
- `frontend/app/(publico)/seguimiento/actions.ts`
- `frontend/components/layout/PublicoShell.tsx`

## Fuera de alcance (no tocar)

- Solicitud.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/seguimiento` con `BuscadorSeguimiento` y `FasesTramite`.

## Validación

1. Con `pnpm dev` y la API levantada, un código ED-… real muestra su estado y uno falso un mensaje genérico.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
