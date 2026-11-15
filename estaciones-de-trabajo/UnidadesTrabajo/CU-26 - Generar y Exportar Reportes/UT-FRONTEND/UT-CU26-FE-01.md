# UT-CU26-FE-01 — Formulario y descarga de reportes

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU26-FE-01` |
| **Nombre** | Formulario y descarga de reportes |
| **Caso de uso** | CU-26-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Generar desde el panel.

## Objetivo

`ReporteForm` y route handler de descarga.

## Alcance

- Sin JavaScript redirige con `?error=`.

Archivos esperados (crear/tocar):

- `frontend/app/panel/reportes/page.tsx`
- `frontend/app/panel/reportes/ReporteForm.tsx`
- `frontend/app/panel/reportes/descargar/route.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`ReporteForm` y route handler de descarga.

## Validación

1. Con `pnpm dev` y la API levantada, generar un PDF de trámites descarga el archivo con membrete.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
