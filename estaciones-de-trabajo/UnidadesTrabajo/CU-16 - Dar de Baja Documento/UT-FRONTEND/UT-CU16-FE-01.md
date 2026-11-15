# UT-CU16-FE-01 — Formularios de baja

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU16-FE-01` |
| **Nombre** | Formularios de baja |
| **Caso de uso** | CU-16-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Retirar un documento con su motivo.

## Objetivo

`BajaDocumentoForm` y `BajaEscaneoForm`.

## Alcance

- Documento tachado con motivo y fecha.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/[id]/FormulariosTramite.tsx`
- `frontend/app/panel/documentos/[id]/page.tsx`
- `frontend/app/panel/tramites/actions.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`BajaDocumentoForm` y `BajaEscaneoForm`.

## Validación

1. Con `pnpm dev` y la API levantada, dar de baja un adjunto lo muestra tachado y el requisito como faltante.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
