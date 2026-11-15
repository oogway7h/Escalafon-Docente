# UT-CU09-FE-02 — Reasignación en el detalle

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU09-FE-02` |
| **Nombre** | Reasignación en el detalle |
| **Caso de uso** | CU-09-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU09-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Cambiar de analista dejando constancia.

## Objetivo

Formulario de reasignación con motivo en `FormulariosTramite`.

## Alcance

- Motivo obligatorio.

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/[id]/FormulariosTramite.tsx`
- `frontend/app/panel/tramites/[id]/page.tsx`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Formulario de reasignación con motivo en `FormulariosTramite`.

## Validación

1. Con `pnpm dev` y la API levantada, reasignar sin motivo muestra el error y con motivo aparece en el historial.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
