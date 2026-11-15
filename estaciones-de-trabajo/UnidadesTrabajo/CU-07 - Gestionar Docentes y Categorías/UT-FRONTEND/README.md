# Unidades de trabajo — CU-07-FE

**Caso de uso:** Listados paginados con carga sin recarga, formularios y eliminación de categoría.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-07 - Gestionar Docentes y Categorías`](../../../CasosDeUso/CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas.md)
**Recurso HTTP / contrato:** `/panel/docentes`, `/panel/categorias`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/docentes | Listado paginado y búsqueda | `DOCENTE_GESTIONAR` |
| /panel/docentes/nuevo, /[id]/editar | Formulario | idem |
| /panel/categorias | Listado con conteos | `CATEGORIA_GESTIONAR` |
| /panel/categorias/[id]/eliminar | Confirmación | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU07-FE-01](./UT-CU07-FE-01.md) | Carga de listados sin recargar | — |
| [UT-CU07-FE-02](./UT-CU07-FE-02.md) | Docentes: listado y formulario | 01 |
| [UT-CU07-FE-03](./UT-CU07-FE-03.md) | Categorías: listado, formulario y eliminación | 01 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Listados y formularios | 01-03 |
| A1-A5 | Errores | 02, 03 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Legajos.
