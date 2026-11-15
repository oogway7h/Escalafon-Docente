# UT-CU03-FE-01 — Perfil con contacto y contraseña

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU03-FE-01` |
| **Nombre** | Perfil con contacto y contraseña |
| **Caso de uso** | CU-03-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Autogestión de datos.

## Objetivo

`/panel/perfil` con `ContactoForm` y `ContrasenaForm`.

## Alcance

- Acciones de servidor con zod.
- El correo y el rol se muestran como solo lectura.

Archivos esperados (crear/tocar):

- `frontend/app/panel/perfil/page.tsx`
- `frontend/app/panel/perfil/ContactoForm.tsx`
- `frontend/app/panel/perfil/ContrasenaForm.tsx`
- `frontend/app/panel/perfil/actions.ts`

## Fuera de alcance (no tocar)

- Usuarios y roles.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/perfil` con `ContactoForm` y `ContrasenaForm`.

## Validación

1. Con `pnpm dev` y la API levantada, cambiar el teléfono persiste y una contraseña actual errónea muestra el error.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
