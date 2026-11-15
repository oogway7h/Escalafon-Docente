# UT-CU02-FE-02 — Restablecer y activar con token

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU02-FE-02` |
| **Nombre** | Restablecer y activar con token |
| **Caso de uso** | CU-02-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU02-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

El mismo formulario sirve para el primer acceso y para el olvido.

## Objetivo

`ContrasenaConTokenForm` compartido por `/restablecer` y `/activar`.

## Alcance

- Lee `token` de la URL, pide la contraseña dos veces.
- Mensajes para `TOKEN_INVALIDO`, `POLITICA_CONTRASENA`, `CUENTA_YA_ACTIVADA`.

Archivos esperados (crear/tocar):

- `frontend/app/(auth)/ContrasenaConTokenForm.tsx`
- `frontend/app/(auth)/restablecer/page.tsx`
- `frontend/app/(auth)/restablecer/actions.ts`
- `frontend/app/(auth)/activar/actions.ts`

## Fuera de alcance (no tocar)

- Login (CU-01).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`ContrasenaConTokenForm` compartido por `/restablecer` y `/activar`.

## Validación

1. Con `pnpm dev` y la API levantada, el enlace de Mailpit permite cambiar la contraseña y reusarlo muestra «enlace inválido».
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
