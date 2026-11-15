# Unidades de trabajo — CU-25-MOB

**Caso de uso:** El Director ve los indicadores y las alertas de mora en el teléfono.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-25 - Visualizar Tablero de Indicadores`](../../../CasosDeUso/CU-25%20-%20Visualizar%20Tablero%20de%20Indicadores.md)
**Recurso HTTP / contrato:** `GET /api/tablero?desde=&hasta=`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Barras dibujadas con `FractionallySizedBox` (una serie, color `grafico`, valor escrito al lado): no se agrega librería de gráficos.
2. Periodo con atajos (Este mes, Últimos 3 meses, Este año, Todo).

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /tablero?desde=&hasta= | 200 `{generadoEn, periodo, diasMaximosPrestamo, tramites:{radicados,enCurso,resueltos,promedioDias}, porEstado, tiempos, analistas, prestamosVigentes, alertas}` | 422 | `TABLERO_VER` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU25-MOB-01](./UT-CU25-MOB-01.md) | Modelo y `TableroService` | — |
| [UT-CU25-MOB-02](./UT-CU25-MOB-02.md) | Pantalla del tablero | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Periodo | 01, 02 |
| CA02 | Alertas | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Reportes (CU-26, web).
