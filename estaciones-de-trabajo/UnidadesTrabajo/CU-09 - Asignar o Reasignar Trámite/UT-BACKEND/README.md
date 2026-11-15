# Unidades de trabajo — CU-09-BE

**Caso de uso:** Asignación y reasignación con motivo y carga de analistas.
**Proyecto:** `backend`
**Fuente CU:** [`CU-09 - Asignar o Reasignar Trámite`](../../../CasosDeUso/CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite.md)
**Recurso HTTP / contrato:** `POST /api/tramites/:id/asignar`, `GET /api/analistas`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/analistas | 200 `[{id, nombre, carga}]` | — | `TRAMITE_ASIGNAR` |
| POST | /api/tramites/:id/asignar `{analistaId, motivo?}` | 200 detalle | 422 `MOTIVO_REQUERIDO`, `ANALISTA_INVALIDO`; 409 `PENDIENTE_VALIDACION`, `TRAMITE_REMITIDO` | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU09-BE-01](./UT-CU09-BE-01.md) | Casos de uso Asignar y ListarAnalistas | — |
| [UT-CU09-BE-02](./UT-CU09-BE-02.md) | Rutas de asignación | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Asignar | 01, 02 |
| A1-A4 | Reglas | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Revisión (CU-10).
