# Unidades de trabajo — CU-16-BE

**Caso de uso:** Baja lógica de adjuntos y escaneos.
**Proyecto:** `backend`
**Fuente CU:** [`CU-16 - Dar de Baja Documento`](../../../CasosDeUso/CU-16%20-%20Dar%20de%20Baja%20Documento.md)
**Recurso HTTP / contrato:** `PATCH /api/tramites/:id/documentos/:docId/estado`, `PATCH /api/docentes/:id/documentos/:docId/estado`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| PATCH | /api/tramites/:id/documentos/:docId/estado `{motivo}` | 200 | 422 `MOTIVO_BAJA_REQUERIDO`; 409 `DOCUMENTO_DE_BAJA`, `TRAMITE_ARCHIVADO`, `NO_ASIGNADO` | `DOCUMENTO_GESTIONAR` |
| PATCH | /api/docentes/:id/documentos/:docId/estado `{motivo}` | 200 | idem | `DOCUMENTO_ESCANEAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU16-BE-01](./UT-CU16-BE-01.md) | Migración y caso de uso DarDeBajaDocumento | — |
| [UT-CU16-BE-02](./UT-CU16-BE-02.md) | Rutas de baja | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Baja | 01, 02 |
| A1-A4 | Reglas | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Restitución (CU-17).
