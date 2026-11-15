# Unidades de trabajo — CU-24-FE

**Caso de uso:** Página de bitácora.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-24 - Consultar Bitácora de Custodia`](../../../CasosDeUso/CU-24%20-%20Consultar%20Bit%C3%A1cora%20de%20Custodia.md)
**Recurso HTTP / contrato:** `/panel/custodia/bitacora`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/custodia/bitacora | Filtros por texto, tipo y vencidas; paginada | `BITACORA_CONSULTAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU24-FE-01](./UT-CU24-FE-01.md) | Página de bitácora | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Consulta | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- —
