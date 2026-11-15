# UT-CU07-FE-01 — Carga de listados sin recargar

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU07-FE-01` |
| **Nombre** | Carga de listados sin recargar |
| **Caso de uso** | CU-07-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Buscar en 1 640 docentes no debe dar saltos de pantalla.

## Objetivo

`CargaListado.tsx` (`ListadoConCarga`, `ZonaCarga`, `BuscadorCarga`, `LinkCarga`) y `Paginacion`.

## Alcance

- Silueta medida del DOM al buscar/paginar.

Archivos esperados (crear/tocar):

- `frontend/components/common/CargaListado.tsx`
- `frontend/components/common/Paginacion.tsx`
- `frontend/components/common/ListadoEsqueleto.tsx`

## Fuera de alcance (no tocar)

- Formularios.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`CargaListado.tsx` (`ListadoConCarga`, `ZonaCarga`, `BuscadorCarga`, `LinkCarga`) y `Paginacion`.

## Validación

1. Con `pnpm dev` y la API levantada, paginar docentes no cambia la altura de la página.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
