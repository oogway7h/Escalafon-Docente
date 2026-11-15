# Unidades de trabajo — CU-09-FE

**Caso de uso:** Bandeja de asignación y reasignación desde el detalle.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-09 - Asignar o Reasignar Trámite`](../../../CasosDeUso/CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite.md)
**Recurso HTTP / contrato:** `/panel/tramites/asignar`, `/panel/tramites/[id]`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/tramites/asignar | Trámites sin asignar y analistas con carga | `TRAMITE_ASIGNAR` |
| /panel/tramites/[id] | Reasignar con motivo | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU09-FE-01](./UT-CU09-FE-01.md) | Bandeja de asignación | — |
| [UT-CU09-FE-02](./UT-CU09-FE-02.md) | Reasignación en el detalle | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Asignar | 01 |
| A1 | Reasignar con motivo | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Revisión.
