# UT-CU21-FE-01 — Pestaña Por atender

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU21-FE-01` |
| **Nombre** | Pestaña Por atender |
| **Caso de uso** | CU-21-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Una sola pantalla para Archivo.

## Objetivo

Pestañas con conteo y panel de acción.

## Alcance

- `aprobar=true|false` por name/value del botón.

Archivos esperados (crear/tocar):

- `frontend/app/panel/custodia/page.tsx`
- `frontend/app/panel/custodia/AccionesCarpeta.tsx`
- `frontend/app/panel/custodia/actions.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Pestañas con conteo y panel de acción.

## Validación

1. Con `pnpm dev` y la API levantada, rechazar con motivo saca la solicitud de Por atender.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
