# UT-CU05-FE-02 — Formulario de rol y eliminación

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU05-FE-02` |
| **Nombre** | Formulario de rol y eliminación |
| **Caso de uso** | CU-05-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU05-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Elegir permisos con claridad y no borrar por error.

## Objetivo

`RolForm` y `EliminarRolForm`.

## Alcance

- Permisos agrupados por módulo.
- Página de confirmación para eliminar.

Archivos esperados (crear/tocar):

- `frontend/app/panel/roles/RolForm.tsx`
- `frontend/app/panel/roles/EliminarRolForm.tsx`
- `frontend/app/panel/roles/[id]/editar/page.tsx`
- `frontend/app/panel/roles/[id]/eliminar/page.tsx`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`RolForm` y `EliminarRolForm`.

## Validación

1. Con `pnpm dev` y la API levantada, eliminar un rol con usuarios muestra «tiene usuarios asignados».
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
