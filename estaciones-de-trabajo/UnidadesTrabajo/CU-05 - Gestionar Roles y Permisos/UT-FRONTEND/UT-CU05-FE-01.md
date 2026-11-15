# UT-CU05-FE-01 — Listado de roles

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU05-FE-01` |
| **Nombre** | Listado de roles |
| **Caso de uso** | CU-05-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Ver qué roles existen y cuántos usuarios tienen.

## Objetivo

`/panel/roles`.

## Alcance

- Tabla con conteo y acciones.

Archivos esperados (crear/tocar):

- `frontend/app/panel/roles/page.tsx`
- `frontend/app/panel/roles/actions.ts`

## Fuera de alcance (no tocar)

- Formulario (UT 02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/roles`.

## Validación

1. Con `pnpm dev` y la API levantada, el listado muestra los cuatro roles de sistema.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
