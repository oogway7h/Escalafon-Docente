# Unidades de trabajo — CU-25-BE

**Caso de uso:** Cálculo del tablero sobre los datos existentes.
**Proyecto:** `backend`
**Fuente CU:** [`CU-25 - Visualizar Tablero de Indicadores`](../../../CasosDeUso/CU-25%20-%20Visualizar%20Tablero%20de%20Indicadores.md)
**Recurso HTTP / contrato:** `GET /api/tablero`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/tablero?desde=&hasta= | 200 `{generadoEn, periodo, diasMaximosPrestamo, tramites:{radicados,enCurso,resueltos,promedioDias}, porEstado:[{estado,cantidad}], tiempos, analistas:[{id,nombre,activos,enFacultad,resueltos}], prestamosVigentes, alertas:[{legajoCodigo,docenteNombre,custodio,tramiteCodigo,salida,fechaLimite,diasPrestado,diasMora,motivo}]}` | 422 `RANGO_FECHAS_INVALIDO`, `VALIDACION` | `TABLERO_VER` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU25-BE-01](./UT-CU25-BE-01.md) | Dominio indicador: mora y tiempos | — |
| [UT-CU25-BE-02](./UT-CU25-BE-02.md) | Consulta del tablero | 01 |
| [UT-CU25-BE-03](./UT-CU25-BE-03.md) | Ruta del tablero | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Indicadores | 01-03 |
| Alertas | Mora | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Reportes (CU-26).
