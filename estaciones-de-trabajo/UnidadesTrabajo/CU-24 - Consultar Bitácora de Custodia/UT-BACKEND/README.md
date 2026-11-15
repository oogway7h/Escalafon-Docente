# Unidades de trabajo — CU-24-BE

**Caso de uso:** Consulta paginada de movimientos.
**Proyecto:** `backend`
**Fuente CU:** [`CU-24 - Consultar Bitácora de Custodia`](../../../CasosDeUso/CU-24%20-%20Consultar%20Bit%C3%A1cora%20de%20Custodia.md)
**Recurso HTTP / contrato:** `GET /api/custodia/bitacora`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/custodia/bitacora?buscar=&tipo=&vencidos=&pagina= | 200 `{datos:[{id,legajoCodigo,docenteNombre,tipo,custodio,fechaLimite,estadoMaterial,observacion,registradoPor,fecha,vencido}],…}` | 422 | `BITACORA_CONSULTAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU24-BE-01](./UT-CU24-BE-01.md) | Caso de uso ListarMovimientos y ruta | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Consulta | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Registro de movimientos (CU-22, CU-23).
