# Unidades de trabajo — CU-06-FE

**Caso de uso:** Página de auditoría con filtros por URL y paginación.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-06 - Consultar Registro de Auditoría`](../../../CasosDeUso/CU-06%20-%20Consultar%20Registro%20de%20Auditor%C3%ADa.md)
**Recurso HTTP / contrato:** `/panel/auditoria`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/auditoria | Filtros (usuario, acción, fechas), tabla de solo lectura, paginación | `AUDITORIA_CONSULTAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU06-FE-01](./UT-CU06-FE-01.md) | Página de auditoría | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal y A1-A3 | Consulta con filtros | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Bitácora de custodia (CU-24).
