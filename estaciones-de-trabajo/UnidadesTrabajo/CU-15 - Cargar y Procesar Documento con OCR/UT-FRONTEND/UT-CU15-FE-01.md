# UT-CU15-FE-01 — Carga de documentos del trámite

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-FE-01` |
| **Nombre** | Carga de documentos del trámite |
| **Caso de uso** | CU-15-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Subir sin pasar el archivo por Next.

## Objetivo

`CargaDocumentos` pide URL firmada, sube con XHR (progreso) y confirma.

## Alcance

- Requiere CORS con PUT en el bucket.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/[id]/CargaDocumentos.tsx`
- `frontend/app/panel/tramites/actions.ts`

## Fuera de alcance (no tocar)

- Visor (UT 02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`CargaDocumentos` pide URL firmada, sube con XHR (progreso) y confirma.

## Validación

1. Con `pnpm dev` y la API levantada, subir un PDF muestra la barra de progreso y luego «Procesando».
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
