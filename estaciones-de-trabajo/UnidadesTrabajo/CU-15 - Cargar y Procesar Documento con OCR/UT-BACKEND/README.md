# Unidades de trabajo — CU-15-BE

**Caso de uso:** Subida por URL firmada, confirmación, OCR asíncrono por página y corrección.
**Proyecto:** `backend`
**Fuente CU:** [`CU-15 - Cargar y Procesar Documento con OCR`](../../../CasosDeUso/CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR.md)
**Recurso HTTP / contrato:** `/api/tramites/:id/documentos/*`, `/api/docentes/:id/documentos/*`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Subida directa a S3 para no cargar la API con archivos de cientos de MB.
2. OCR en un servicio aparte (FastAPI + Tesseract + ocrmypdf) detrás del puerto `ProcesadorDocumentos`.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/tramites/:id/documentos/subida `{nombre, tipoMime, tamano, requisitoId?}` | 201 `{documento, url}` | 422 `TIPO_NO_PERMITIDO`, `ARCHIVO_GRANDE`; 409 `NO_ASIGNADO` | `DOCUMENTO_GESTIONAR` |
| POST | /api/tramites/:id/documentos/:docId/confirmar | 200 documento | 409 `SUBIDA_INCOMPLETA` | idem |
| GET/PATCH | /api/tramites/:id/documentos/:docId/texto `{paginas}` | 200 | 422 `PAGINAS_INVALIDAS` | lectura / idem |
| POST | /api/docentes/:id/documentos/subida `{…, tipoDocumento}` | 201 | 422 `TIPO_DOCUMENTO_INVALIDO` | `DOCUMENTO_ESCANEAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU15-BE-01](./UT-CU15-BE-01.md) | Migraciones de OCR y páginas | — |
| [UT-CU15-BE-02](./UT-CU15-BE-02.md) | Puerto ProcesadorDocumentos y cliente ia-service | 01 |
| [UT-CU15-BE-03](./UT-CU15-BE-03.md) | Casos de uso PrepararSubida, ConfirmarSubida y CorregirPaginas | 02 |
| [UT-CU15-BE-04](./UT-CU15-BE-04.md) | Expediente digital del docente y rutas | 03 |

**Cadena:** 01 → 02 → 03 → 04

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Pasos 1-4 | Subida y confirmación | 01, 03 |
| Pasos 5-6 | OCR | 02 |
| Paso 7 | Corrección | 03 |
| Expediente del docente | Escaneos por tipo | 04 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Baja de documentos (CU-16).
- Búsqueda (CU-18).
