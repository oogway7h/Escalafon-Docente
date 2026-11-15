# Unidades de trabajo — CU-06-BE

**Caso de uso:** Consulta paginada y filtrable de la bitácora de solo anexión.
**Proyecto:** `backend`
**Fuente CU:** [`CU-06 - Consultar Registro de Auditoría`](../../../CasosDeUso/CU-06%20-%20Consultar%20Registro%20de%20Auditor%C3%ADa.md)
**Recurso HTTP / contrato:** `GET /api/auditoria`, `GET /api/auditoria/filtros`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. El dominio separa el puerto de escritura (`auditoria.Repository.Registrar`) del de lectura (`auditoria.Consulta`) para que ningún caso de uso pueda modificar la bitácora.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/auditoria?usuarioId=&accion=&desde=&hasta=&pagina=&tamano= | 200 `{datos,total,pagina,tamano}` | 422 `ACCION_INVALIDA`, `RANGO_FECHAS_INVALIDO` | `AUDITORIA_CONSULTAR` |
| GET | /api/auditoria/filtros | 200 `{acciones, usuarios}` | — | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU06-BE-01](./UT-CU06-BE-01.md) | Tabla de auditoría de solo anexión | — |
| [UT-CU06-BE-02](./UT-CU06-BE-02.md) | Caso de uso Consultar y ListarOpciones | 01 |
| [UT-CU06-BE-03](./UT-CU06-BE-03.md) | Handler y rutas de auditoría | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Consulta paginada | 02, 03 |
| A1, A2 | Validación de filtros | 02 |
| Regla de solo anexión | Trigger | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Bitácora de custodia (CU-24).
