# UT-CU08-FE-01 — Formulario de radicación

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU08-FE-01` |
| **Nombre** | Formulario de radicación |
| **Caso de uso** | CU-08-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Ventanilla atiende al docente en el mostrador.

## Objetivo

`RadicarForm` con buscador de docente.

## Alcance

- Muestra el código al terminar.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/radicar/page.tsx`
- `frontend/app/panel/tramites/radicar/RadicarForm.tsx`
- `frontend/app/panel/tramites/actions.ts`

## Fuera de alcance (no tocar)

- Asignación.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`RadicarForm` con buscador de docente.

## Validación

1. Con `pnpm dev` y la API levantada, radicar con `ventanilla@ded.dev` muestra un código ED-….
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
