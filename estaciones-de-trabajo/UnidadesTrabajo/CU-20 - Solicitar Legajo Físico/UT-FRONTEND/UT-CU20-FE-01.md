# UT-CU20-FE-01 — Mis solicitudes y nueva solicitud

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU20-FE-01` |
| **Nombre** | Mis solicitudes y nueva solicitud |
| **Caso de uso** | CU-20-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Pedir y seguir carpetas.

## Objetivo

`/panel/carpetas` y `SolicitarCarpetaModal`.

## Alcance

- Etiquetas de `lib/legajos.ts`.

Archivos esperados (crear/tocar):

- `frontend/app/panel/carpetas/page.tsx`
- `frontend/app/panel/carpetas/nueva/page.tsx`
- `frontend/app/panel/carpetas/SolicitarCarpetaModal.tsx`
- `frontend/app/panel/carpetas/actions.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/carpetas` y `SolicitarCarpetaModal`.

## Validación

1. Con `pnpm dev` y la API levantada, solicitar una carpeta la muestra Pendiente.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
