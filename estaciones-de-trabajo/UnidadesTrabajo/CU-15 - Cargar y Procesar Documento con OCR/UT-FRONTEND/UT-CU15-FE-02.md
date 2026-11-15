# UT-CU15-FE-02 — Visor y corrección del texto OCR

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-FE-02` |
| **Nombre** | Visor y corrección del texto OCR |
| **Caso de uso** | CU-15-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU15-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Leer el documento y corregir lo que el OCR leyó mal.

## Objetivo

Visor pdf.js propio y edición por página.

## Alcance

- Route handler que pasa rangos de S3.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/[id]/VisorDocumento.tsx`
- `frontend/components/documentos/VisorPdf.tsx`
- `frontend/app/panel/tramites/[id]/archivo/[docId]/route.ts`
- `frontend/app/panel/tramites/[id]/documentos/[docId]/texto/route.ts`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Visor pdf.js propio y edición por página.

## Validación

1. Con `pnpm dev` y la API levantada, corregir la página 1 la marca como validada.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
