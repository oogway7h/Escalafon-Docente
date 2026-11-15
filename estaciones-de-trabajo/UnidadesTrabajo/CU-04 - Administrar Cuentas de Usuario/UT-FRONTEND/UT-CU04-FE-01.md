# UT-CU04-FE-01 — Listado de usuarios

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU04-FE-01` |
| **Nombre** | Listado de usuarios |
| **Caso de uso** | CU-04-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Ver al personal y su estado.

## Objetivo

`/panel/usuarios` con búsqueda `next/form` y acciones de baja/alta/reenvío.

## Alcance

- Tabla con estado y rol.
- `FormularioCarga` para las acciones de una fila.

Archivos esperados (crear/tocar):

- `frontend/app/panel/usuarios/page.tsx`
- `frontend/app/panel/usuarios/actions.ts`

## Fuera de alcance (no tocar)

- Formularios (UT 02).
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/panel/usuarios` con búsqueda `next/form` y acciones de baja/alta/reenvío.

## Validación

1. Con `pnpm dev` y la API levantada, el listado filtra por nombre sin recargar el documento.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
