# UT-CU02-FE-01 — Pantalla de recuperación

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU02-FE-01` |
| **Nombre** | Pantalla de recuperación |
| **Caso de uso** | CU-02-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

El usuario sin acceso necesita pedir ayuda sin depender del Director.

## Objetivo

`/recuperar` con acción de servidor y mensaje neutro.

## Alcance

- `recuperarAction` con zod.
- Mensaje igual exista o no el correo.

Archivos esperados (crear/tocar):

- `frontend/app/(auth)/recuperar/page.tsx`
- `frontend/app/(auth)/recuperar/RecuperarForm.tsx`
- `frontend/app/(auth)/recuperar/actions.ts`

## Fuera de alcance (no tocar)

- Restablecer (UT 02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/recuperar` con acción de servidor y mensaje neutro.

## Validación

1. Con `pnpm dev` y la API levantada, pedir recuperación llega un correo a Mailpit (http://localhost:8025).
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
