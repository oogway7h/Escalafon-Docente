# UT-CU25-FE-01 — Página del tablero

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU25-FE-01` |
| **Nombre** | Página del tablero |
| **Caso de uso** | CU-25-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Decidir con datos.

## Objetivo

`/panel/tablero` con `ActualizacionAutomatica`.

## Alcance

- Barras de una serie con `--color-grafico`; rojo solo para mora.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tablero/page.tsx`
- `frontend/app/panel/tablero/ActualizacionAutomatica.tsx`
- `frontend/lib/api-indicadores.ts`
- `frontend/lib/indicadores.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/tablero` con `ActualizacionAutomatica`.

## Validación

1. Con `pnpm dev` y la API levantada, con una carpeta vencida aparece el aviso que lleva a las alertas.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
