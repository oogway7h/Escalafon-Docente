# UT-CU07-FE-03 — Categorías: listado, formulario y eliminación

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU07-FE-03` |
| **Nombre** | Categorías: listado, formulario y eliminación |
| **Caso de uso** | CU-07-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU07-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Editar requisitos cambia el portal.

## Objetivo

`/panel/categorias` con `CategoriaForm` y `EliminarCategoriaForm`.

## Alcance

- Requisitos como lista editable.
- Eliminar solo sin uso.

Archivos esperados (crear/tocar):

- `frontend/app/panel/categorias/page.tsx`
- `frontend/app/panel/categorias/CategoriaForm.tsx`
- `frontend/app/panel/categorias/EliminarCategoriaForm.tsx`
- `frontend/app/panel/categorias/actions.ts`

## Fuera de alcance (no tocar)

- Docentes.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/categorias` con `CategoriaForm` y `EliminarCategoriaForm`.

## Validación

1. Con `pnpm dev` y la API levantada, eliminar una categoría en uso muestra cuántos docentes y trámites la usan.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
