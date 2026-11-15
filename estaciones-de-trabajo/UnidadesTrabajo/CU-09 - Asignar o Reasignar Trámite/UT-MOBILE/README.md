# Unidades de trabajo — CU-09-MOB

**Caso de uso:** El Director asigna los trámites ingresados viendo la carga de cada analista, y reasigna con motivo.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-09 - Asignar o Reasignar Trámite`](../../../CasosDeUso/CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite.md)
**Recurso HTTP / contrato:** `GET /api/tramites?sinAsignar=true`, `GET /api/analistas`, `POST /api/tramites/:id/asignar`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /tramites?sinAsignar=true | 200 lista | — | `TRAMITE_ASIGNAR` |
| GET | /analistas | 200 `[{id,nombre,carga}]` | — | idem |
| POST | /tramites/:id/asignar `{analistaId, motivo?}` | 200 detalle | 422 `MOTIVO_REQUERIDO`, `ANALISTA_INVALIDO`; 409 `PENDIENTE_VALIDACION` | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU09-MOB-01](./UT-CU09-MOB-01.md) | Asignación en `TramitesService` | — |
| [UT-CU09-MOB-02](./UT-CU09-MOB-02.md) | Pantalla «Por asignar» y reasignación en el detalle | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Asignar | 01, 02 |
| CA02, CA03 | Reasignar con motivo | 02 |
| CA04 | Carga visible | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Validar solicitudes del portal (ventanilla, web).
