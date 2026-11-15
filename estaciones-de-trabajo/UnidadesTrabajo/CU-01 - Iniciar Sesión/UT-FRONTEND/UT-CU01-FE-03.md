# UT-CU01-FE-03 — Layout del panel, menú por permisos y cierre de sesión

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-FE-03` |
| **Nombre** | Layout del panel, menú por permisos y cierre de sesión |
| **Caso de uso** | CU-01-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU01-FE-01, UT-CU01-FE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Cada rol debe ver solo sus módulos y poder salir.

## Objetivo

`PanelShell` con `Sidebar` armado desde `lib/navegacion.ts` y botón de cerrar sesión.

## Alcance

- `app/panel/layout.tsx` carga la sesión (`GET /auth/sesion`).
- `NAVEGACION` filtra por permisos.
- `cerrarSesionAction` llama a `/auth/logout` y borra la cookie.

Archivos esperados (crear/tocar):

- `frontend/app/panel/layout.tsx`
- `frontend/components/layout/PanelShell.tsx`
- `frontend/components/layout/Sidebar.tsx`
- `frontend/lib/navegacion.ts`

## Fuera de alcance (no tocar)

- Contenido de cada módulo.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`PanelShell` con `Sidebar` armado desde `lib/navegacion.ts` y botón de cerrar sesión.

## Validación

1. Con `pnpm dev` y la API levantada, `archivo@ded.dev` ve solo legajos y custodia; cerrar sesión vuelve a `/login`.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
