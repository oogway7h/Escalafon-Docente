# UT-CU23-FE-01 — Formulario de devolución

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU23-FE-01` |
| **Nombre** | Formulario de devolución |
| **Caso de uso** | CU-23-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Recibir la carpeta.

## Objetivo

Panel de devolución con estado material.

## Alcance

- Observación obligatoria si no es Buen estado.

Archivos esperados (crear/tocar):

- `frontend/app/panel/custodia/AccionesCarpeta.tsx`
- `frontend/app/panel/custodia/actions.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Panel de devolución con estado material.

## Validación

1. Con `pnpm dev` y la API levantada, devolver una carpeta deteriorada sin observación muestra el error.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
