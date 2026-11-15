# UT-CU10-FE-02 — Detalle del trámite y línea de fases

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU10-FE-02` |
| **Nombre** | Detalle del trámite y línea de fases |
| **Caso de uso** | CU-10-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU10-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Contrastar requisitos y ver el recorrido.

## Objetivo

`/panel/tramites/[id]` con `FasesTramite`.

## Alcance

- Requisitos con documento presente o no; historial.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/[id]/page.tsx`
- `frontend/components/tramites/FasesTramite.tsx`

## Fuera de alcance (no tocar)

- Formularios de acción (UT 03).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/tramites/[id]` con `FasesTramite`.

## Validación

1. Con `pnpm dev` y la API levantada, el detalle muestra la fase actual resaltada.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
