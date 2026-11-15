# Unidades de trabajo — CU-17-BE

**Caso de uso:** Bandeja de bajas y restitución.
**Proyecto:** `backend`
**Fuente CU:** [`CU-17 - Restituir Documento Dado de Baja`](../../../CasosDeUso/CU-17%20-%20Restituir%20Documento%20Dado%20de%20Baja.md)
**Recurso HTTP / contrato:** `GET /api/documentos/bajas`, `POST /api/documentos/:docId/restituir`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/documentos/bajas?buscar=&pagina= | 200 `{datos,total,pagina,tamano}` | — | `DOCUMENTO_RESTITUIR` |
| POST | /api/documentos/:docId/restituir `{motivo}` | 200 | 422 `MOTIVO_RESTITUCION_REQUERIDO`; 409 `DOCUMENTO_ACTIVO` | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU17-BE-01](./UT-CU17-BE-01.md) | Casos de uso ListarDocumentosDeBaja y RestituirDocumento | — |
| [UT-CU17-BE-02](./UT-CU17-BE-02.md) | Rutas de restitución | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Bandeja y restitución | 01, 02 |
| A1-A3 | Reglas | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Baja (CU-16).
