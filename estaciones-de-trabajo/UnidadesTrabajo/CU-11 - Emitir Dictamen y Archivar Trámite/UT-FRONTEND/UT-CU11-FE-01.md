# UT-CU11-FE-01 — Formularios del flujo en el detalle

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU11-FE-01` |
| **Nombre** | Formularios del flujo en el detalle |
| **Caso de uso** | CU-11-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Cada estado ofrece solo su acción.

## Objetivo

Dictamen, periodo, remisiones y regreso en `FormulariosTramite`.

## Alcance

- Visibles según estado, permiso y si es el asignado.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/[id]/FormulariosTramite.tsx`
- `frontend/app/panel/tramites/[id]/page.tsx`
- `frontend/app/panel/tramites/actions.ts`

## Fuera de alcance (no tocar)

- Archivar (UT 02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Dictamen, periodo, remisiones y regreso en `FormulariosTramite`.

## Validación

1. Con `pnpm dev` y la API levantada, un trámite En revisión sin periodo muestra «Años a evaluar» antes del dictamen.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
