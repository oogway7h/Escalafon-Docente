# UT-CU11-FE-02 — Bandeja Por archivar

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU11-FE-02` |
| **Nombre** | Bandeja Por archivar |
| **Caso de uso** | CU-11-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU11-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Archivo cierra lo remitido.

## Objetivo

`/panel/tramites/archivar`.

## Alcance

- Bloquea si hay carpetas sin devolver.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/archivar/page.tsx`
- `frontend/app/panel/tramites/actions.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/tramites/archivar`.

## Validación

1. Con `pnpm dev` y la API levantada, `archivo@ded.dev` archiva un remitido y pasa a Archivado.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
