# UT-CU08-FE-02 — Listado de trámites y etiquetas de estado

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU08-FE-02` |
| **Nombre** | Listado de trámites y etiquetas de estado |
| **Caso de uso** | CU-08-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU08-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Ver todos los trámites con su estado.

## Objetivo

`/panel/tramites` con `TablaTramites`, `EstadoInsignia` y `lib/tramites.ts`.

## Alcance

- Fila clicable; filtros por estado y búsqueda.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/page.tsx`
- `frontend/components/tramites/TablaTramites.tsx`
- `frontend/components/tramites/EstadoInsignia.tsx`
- `frontend/lib/tramites.ts`

## Fuera de alcance (no tocar)

- Detalle (CU-10).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/tramites` con `TablaTramites`, `EstadoInsignia` y `lib/tramites.ts`.

## Validación

1. Con `pnpm dev` y la API levantada, filtrar por estado Ingresado conserva el filtro al paginar.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
