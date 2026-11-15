# UT-CU14-FE-01 — Página de notificaciones

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU14-FE-01` |
| **Nombre** | Página de notificaciones |
| **Caso de uso** | CU-14-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Ver qué avisos salieron y cuáles fallaron.

## Objetivo

`/panel/notificaciones`.

## Alcance

- Tabla paginada con filtro.

Archivos esperados (crear/tocar):

- `frontend/app/panel/notificaciones/page.tsx`
- `frontend/lib/api-tramites.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/notificaciones`.

## Validación

1. Con `pnpm dev` y la API levantada, observar un trámite hace aparecer un aviso ENVIADO y el correo en Mailpit.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
