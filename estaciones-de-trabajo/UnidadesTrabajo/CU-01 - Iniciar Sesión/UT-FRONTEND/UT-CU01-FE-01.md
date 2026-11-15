# UT-CU01-FE-01 — Cliente de API con sesión y guardia de rutas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-FE-01` |
| **Nombre** | Cliente de API con sesión y guardia de rutas |
| **Caso de uso** | CU-01-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

El navegador nunca debe ver el token ni hablar con la API.

## Objetivo

`lib/api.ts` agrega el token desde la cookie httpOnly y `proxy.ts` protege `/panel/*`.

## Alcance

- `api()` en servidor con `Authorization: Bearer`.
- 401 → redirige a `/login`; 403 → 404.
- `lib/sesion.ts` lee/escribe la cookie.
- `proxy.ts` redirige sin cookie.

Archivos esperados (crear/tocar):

- `frontend/lib/api.ts`
- `frontend/lib/sesion.ts`
- `frontend/proxy.ts`

## Fuera de alcance (no tocar)

- Formularios (UT 02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`lib/api.ts` agrega el token desde la cookie httpOnly y `proxy.ts` protege `/panel/*`.

## Validación

1. Con `pnpm dev` y la API levantada, abrir `/panel` sin cookie redirige a `/login`.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
