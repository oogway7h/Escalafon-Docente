# UT-CU18-FE-01 — Buscador global

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU18-FE-01` |
| **Nombre** | Buscador global |
| **Caso de uso** | CU-18-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Buscar desde una sola pantalla.

## Objetivo

`BusquedaGlobal` con filtros y botones «Pág. N».

## Alcance

- Abre el visor en la página.

Archivos esperados (crear/tocar):

- `frontend/app/panel/documentos/buscar/page.tsx`
- `frontend/app/panel/documentos/buscar/BusquedaGlobal.tsx`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`BusquedaGlobal` con filtros y botones «Pág. N».

## Validación

1. Con `pnpm dev` y la API levantada, buscar «resolucion» muestra fragmentos con «resolución» resaltado.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
