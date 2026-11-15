# Unidades de trabajo — CU-26-FE

**Caso de uso:** Formulario de reportes con descarga sin salir de la página.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-26 - Generar y Exportar Reportes`](../../../CasosDeUso/CU-26%20-%20Generar%20y%20Exportar%20Reportes.md)
**Recurso HTTP / contrato:** `/panel/reportes`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/reportes | Tipo, filtros del tipo y formato; descarga por `/panel/reportes/descargar` | `REPORTE_GENERAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU26-FE-01](./UT-CU26-FE-01.md) | Formulario y descarga de reportes | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal y A1-A3 | Reportes | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- —
