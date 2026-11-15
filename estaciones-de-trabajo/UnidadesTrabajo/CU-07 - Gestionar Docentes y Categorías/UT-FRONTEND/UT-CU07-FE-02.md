# UT-CU07-FE-02 — Docentes: listado y formulario

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU07-FE-02` |
| **Nombre** | Docentes: listado y formulario |
| **Caso de uso** | CU-07-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU07-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Mantener el padrón desde el panel.

## Objetivo

`/panel/docentes` con `DocenteForm`.

## Alcance

- Código primera columna; cédula opcional.

Archivos esperados (crear/tocar):

- `frontend/app/panel/docentes/page.tsx`
- `frontend/app/panel/docentes/DocenteForm.tsx`
- `frontend/app/panel/docentes/actions.ts`
- `frontend/lib/api-tramites.ts`

## Fuera de alcance (no tocar)

- Categorías.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/docentes` con `DocenteForm`.

## Validación

1. Con `pnpm dev` y la API levantada, un código repetido muestra «código ya registrado».
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
