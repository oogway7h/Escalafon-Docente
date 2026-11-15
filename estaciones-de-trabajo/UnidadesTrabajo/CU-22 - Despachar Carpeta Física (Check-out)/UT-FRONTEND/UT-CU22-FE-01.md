# UT-CU22-FE-01 — Formulario de entrega

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU22-FE-01` |
| **Nombre** | Formulario de entrega |
| **Caso de uso** | CU-22-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Despachar en un paso.

## Objetivo

Panel de acción con atajos de fecha.

## Alcance

- `enDias(n)` de `lib/legajos.ts`.

Archivos esperados (crear/tocar):

- `frontend/app/panel/custodia/AccionesCarpeta.tsx`
- `frontend/app/panel/custodia/actions.ts`
- `frontend/lib/legajos.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Panel de acción con atajos de fecha.

## Validación

1. Con `pnpm dev` y la API levantada, entregar una carpeta la mueve a la pestaña Prestadas.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
