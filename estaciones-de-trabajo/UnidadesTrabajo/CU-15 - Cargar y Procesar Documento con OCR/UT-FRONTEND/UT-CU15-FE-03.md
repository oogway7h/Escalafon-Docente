# UT-CU15-FE-03 — Expediente digital del docente

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-FE-03` |
| **Nombre** | Expediente digital del docente |
| **Caso de uso** | CU-15-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU15-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Digitalizar la carpeta física por tipo.

## Objetivo

`/panel/documentos` y `/panel/documentos/[id]` con `SubirEscaneo`.

## Alcance

- Sondeo del OCR cada 6 s.

Archivos esperados (crear/tocar):

- `frontend/app/panel/documentos/page.tsx`
- `frontend/app/panel/documentos/[id]/page.tsx`
- `frontend/app/panel/documentos/[id]/SubirEscaneo.tsx`
- `frontend/app/panel/documentos/actions.ts`

## Fuera de alcance (no tocar)

- Búsqueda.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/documentos` y `/panel/documentos/[id]` con `SubirEscaneo`.

## Validación

1. Con `pnpm dev` y la API levantada, `archivo@ded.dev` sube un escaneo y ve su estado OCR.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
