# UT-CU09-FE-01 — Bandeja de asignación

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU09-FE-01` |
| **Nombre** | Bandeja de asignación |
| **Caso de uso** | CU-09-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

El Director reparte desde una pantalla.

## Objetivo

`/panel/tramites/asignar`.

## Alcance

- Selector de analista con su carga.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/asignar/page.tsx`
- `frontend/app/panel/tramites/actions.ts`

## Fuera de alcance (no tocar)

- Reasignación (UT 02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/tramites/asignar`.

## Validación

1. Con `pnpm dev` y la API levantada, asignar un trámite lo saca de la bandeja.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
