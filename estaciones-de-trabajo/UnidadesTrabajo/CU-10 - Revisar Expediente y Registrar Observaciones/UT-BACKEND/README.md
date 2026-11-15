# Unidades de trabajo — CU-10-BE

**Caso de uso:** Inicio de revisión, observaciones y su resolución.
**Proyecto:** `backend`
**Fuente CU:** [`CU-10 - Revisar Expediente y Registrar Observaciones`](../../../CasosDeUso/CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones.md)
**Recurso HTTP / contrato:** `POST /api/tramites/:id/revision`, `/observaciones`, `/observaciones/:oid/resolver`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/tramites/:id/revision | 200 detalle | 409 `CARPETA_NO_SOLICITADA`, `NO_ASIGNADO` | `TRAMITE_REVISAR` |
| POST | /api/tramites/:id/observaciones `{documento, motivo, plazo}` | 201 | 422 `VALIDACION`, `PLAZO_INVALIDO` | idem |
| POST | /api/tramites/:id/observaciones/:oid/resolver | 200 | 409 | idem |
| GET | /api/tramites?mios=true | 200 mi bandeja | — | lectura |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU10-BE-01](./UT-CU10-BE-01.md) | Casos de uso IniciarRevision, RegistrarObservacion y ResolverObservacion | — |
| [UT-CU10-BE-02](./UT-CU10-BE-02.md) | Rutas de revisión | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Revisar y observar | 01, 02 |
| A1-A4 | Reglas | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Dictamen (CU-11).
- Solicitud de carpeta (CU-20).
