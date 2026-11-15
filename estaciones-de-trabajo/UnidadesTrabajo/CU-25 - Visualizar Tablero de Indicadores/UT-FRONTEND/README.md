# Unidades de trabajo — CU-25-FE

**Caso de uso:** Tablero con periodo, barras y alertas, con actualización automática.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-25 - Visualizar Tablero de Indicadores`](../../../CasosDeUso/CU-25%20-%20Visualizar%20Tablero%20de%20Indicadores.md)
**Recurso HTTP / contrato:** `/panel/tablero`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/tablero | Periodo con atajos, 4 indicadores, barras por estado y por analista, alertas de mora | `TABLERO_VER` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU25-FE-01](./UT-CU25-FE-01.md) | Página del tablero | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Tablero | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Reportes.
