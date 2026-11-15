# UT-CU17-FE-01 — Bandeja de bajas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU17-FE-01` |
| **Nombre** | Bandeja de bajas |
| **Caso de uso** | CU-17-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

El Director revisa y corrige bajas.

## Objetivo

`/panel/documentos/bajas` con `RestituirForm`.

## Alcance

- `restituirDocumentoAction`.

Archivos esperados (crear/tocar):

- `frontend/app/panel/documentos/bajas/page.tsx`
- `frontend/app/panel/documentos/bajas/RestituirForm.tsx`
- `frontend/app/panel/documentos/actions.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/documentos/bajas` con `RestituirForm`.

## Validación

1. Con `pnpm dev` y la API levantada, restituir un documento lo saca de la bandeja.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
