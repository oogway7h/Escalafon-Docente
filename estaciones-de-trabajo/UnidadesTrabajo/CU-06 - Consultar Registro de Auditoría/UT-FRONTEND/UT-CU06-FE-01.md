# UT-CU06-FE-01 — Página de auditoría

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU06-FE-01` |
| **Nombre** | Página de auditoría |
| **Caso de uso** | CU-06-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Consultar la bitácora sin poder alterarla.

## Objetivo

`/panel/auditoria` con `FiltrosCarga` y `Paginacion`.

## Alcance

- Etiquetas de acciones en `lib/auditoria.ts`.
- Tabla sin acciones.

Archivos esperados (crear/tocar):

- `frontend/app/panel/auditoria/page.tsx`
- `frontend/lib/api-auditoria.ts`
- `frontend/lib/auditoria.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/auditoria` con `FiltrosCarga` y `Paginacion`.

## Validación

1. Con `pnpm dev` y la API levantada, filtrar por usuario y fechas conserva los filtros al paginar.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
