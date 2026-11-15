# UT-CU01-FE-02 — Pantalla y acción de login

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-FE-02` |
| **Nombre** | Pantalla y acción de login |
| **Caso de uso** | CU-01-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU01-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

El usuario necesita entender por qué no pudo entrar sin revelar si el correo existe.

## Objetivo

`/login` con `useActionState` + zod y mensajes de `MENSAJES_API`.

## Alcance

- `loginAction` valida con zod, llama a la API y guarda la cookie.
- Mensajes para `CREDENCIALES_INVALIDAS`, `CUENTA_BLOQUEADA`, 429.

Archivos esperados (crear/tocar):

- `frontend/app/(auth)/login/page.tsx`
- `frontend/app/(auth)/login/LoginForm.tsx`
- `frontend/app/(auth)/login/actions.ts`
- `frontend/lib/formularios.ts`

## Fuera de alcance (no tocar)

- Recuperar contraseña (CU-02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/login` con `useActionState` + zod y mensajes de `MENSAJES_API`.

## Validación

1. Con `pnpm dev` y la API levantada, un login con `admin@ded.dev`/`123123` entra y uno con contraseña errónea muestra el mensaje genérico.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
