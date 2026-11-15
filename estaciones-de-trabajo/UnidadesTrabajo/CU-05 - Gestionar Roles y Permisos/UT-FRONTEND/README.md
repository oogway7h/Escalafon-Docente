# Unidades de trabajo — CU-05-FE

**Caso de uso:** Listado, formulario con permisos agrupados y eliminación confirmada.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-05 - Gestionar Roles y Permisos`](../../../CasosDeUso/CU-05%20-%20Gestionar%20Roles%20y%20Permisos.md)
**Recurso HTTP / contrato:** `/panel/roles`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/roles | Listado | `ROL_GESTIONAR` |
| /panel/roles/nuevo | Alta con permisos | idem |
| /panel/roles/[id]/eliminar | Confirmación | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU05-FE-01](./UT-CU05-FE-01.md) | Listado de roles | — |
| [UT-CU05-FE-02](./UT-CU05-FE-02.md) | Formulario de rol y eliminación | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | CRUD | 01, 02 |
| A1-A3 | Errores | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Usuarios (CU-04).
