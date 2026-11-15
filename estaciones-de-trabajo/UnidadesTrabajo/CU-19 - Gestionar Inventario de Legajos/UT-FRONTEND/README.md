# Unidades de trabajo — CU-19-FE

**Caso de uso:** Inventario paginado con formulario.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-19 - Gestionar Inventario de Legajos`](../../../CasosDeUso/CU-19%20-%20Gestionar%20Inventario%20de%20Legajos.md)
**Recurso HTTP / contrato:** `/panel/legajos`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/legajos | Inventario con búsqueda y estado | `LEGAJO_GESTIONAR` |
| /panel/legajos/nuevo, /[id]/editar | Formulario de ubicación | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU19-FE-01](./UT-CU19-FE-01.md) | Inventario y formulario de legajo | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Inventario | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Custodia.
