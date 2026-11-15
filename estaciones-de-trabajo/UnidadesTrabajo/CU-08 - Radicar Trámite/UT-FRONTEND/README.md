# Unidades de trabajo — CU-08-FE

**Caso de uso:** Formulario de radicación y listado general.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-08 - Radicar Trámite`](../../../CasosDeUso/CU-08%20-%20Radicar%20Tr%C3%A1mite.md)
**Recurso HTTP / contrato:** `/panel/tramites/radicar`, `/panel/tramites`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/tramites/radicar | Buscar docente, elegir tipo y categoría | `TRAMITE_RADICAR` |
| /panel/tramites | Listado paginado con filtros | lectura |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU08-FE-01](./UT-CU08-FE-01.md) | Formulario de radicación | — |
| [UT-CU08-FE-02](./UT-CU08-FE-02.md) | Listado de trámites y etiquetas de estado | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Radicar y ver | 01, 02 |
| A3 | Validar solicitud del portal | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Asignar (CU-09).
