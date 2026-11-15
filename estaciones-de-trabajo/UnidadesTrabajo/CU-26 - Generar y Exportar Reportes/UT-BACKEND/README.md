# Unidades de trabajo — CU-26-BE

**Caso de uso:** Generación de reportes PDF/Excel con auditoría.
**Proyecto:** `backend`
**Fuente CU:** [`CU-26 - Generar y Exportar Reportes`](../../../CasosDeUso/CU-26%20-%20Generar%20y%20Exportar%20Reportes.md)
**Recurso HTTP / contrato:** `POST /api/reportes`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/reportes `{tipo, formato, desde, hasta, estado, tipoTramite, categoriaId, estadoLegajo}` | 200 archivo (`Content-Disposition`, `X-Reporte-Filas`) | 422 `REPORTE_VACIO`, `RANGO_FECHAS_INVALIDO`, `FILTRO_INVALIDO` | `REPORTE_GENERAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU26-BE-01](./UT-CU26-BE-01.md) | Caso de uso GenerarReporte | — |
| [UT-CU26-BE-02](./UT-CU26-BE-02.md) | Exportadores PDF y Excel | 01 |
| [UT-CU26-BE-03](./UT-CU26-BE-03.md) | Ruta de reportes | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Generar | 01-03 |
| A1-A3 | Validación | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Tablero (CU-25).
