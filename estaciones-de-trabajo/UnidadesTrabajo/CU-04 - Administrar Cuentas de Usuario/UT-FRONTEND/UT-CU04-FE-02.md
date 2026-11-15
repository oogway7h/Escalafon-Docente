# UT-CU04-FE-02 — Formularios de usuario y de rol

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU04-FE-02` |
| **Nombre** | Formularios de usuario y de rol |
| **Caso de uso** | CU-04-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU04-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Alta y edición con validación consistente con la API.

## Objetivo

`UsuarioForm` y `CambiarRolForm`.

## Alcance

- Zod en `lib/esquemas.ts`; errores de `MENSAJES_API`.

Archivos esperados (crear/tocar):

- `frontend/app/panel/usuarios/UsuarioForm.tsx`
- `frontend/app/panel/usuarios/CambiarRolForm.tsx`
- `frontend/app/panel/usuarios/nuevo/page.tsx`
- `frontend/app/panel/usuarios/[id]/rol/page.tsx`

## Fuera de alcance (no tocar)

- Roles.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`UsuarioForm` y `CambiarRolForm`.

## Validación

1. Con `pnpm dev` y la API levantada, crear un usuario con correo repetido muestra «correo ya registrado».
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
