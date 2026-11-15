# Unidades de trabajo — CU-15-MOB

**Caso de uso:** El analista fotografía o elige un documento y lo sube al trámite; la API lo procesa por OCR.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-15 - Cargar y Procesar Documento con OCR`](../../../CasosDeUso/CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR.md)
**Recurso HTTP / contrato:** `POST /api/tramites/:id/documentos/subida`, PUT a la URL firmada, `POST …/:docId/confirmar`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. La subida a S3 usa un `Dio` **sin** baseUrl ni Bearer (S3 rechaza el `Authorization`).
2. Límites validados en el teléfono antes de pedir la URL: JPG/PNG ≤ 20 MB, PDF ≤ 400 MB.
3. El estado OCR se ve al recargar el detalle; no hay sondeo automático en el móvil.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /tramites/:id/documentos/subida `{nombre, tipoMime, tamano, requisitoId?}` | 201 `{documento, url}` | 422 `TIPO_NO_PERMITIDO`, `ARCHIVO_GRANDE`; 409 `NO_ASIGNADO` | `DOCUMENTO_GESTIONAR` |
| PUT | `url` firmada (S3), cabecera `Content-Type` = `tipoMime` | 200 | 403 si cambia el tipo | — |
| POST | /tramites/:id/documentos/:docId/confirmar | 200 documento (`estadoOcr`) | 409 `SUBIDA_INCOMPLETA` | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU15-MOB-01](./UT-CU15-MOB-01.md) | `DocumentosService`: preparar, subir y confirmar | — |
| [UT-CU15-MOB-02](./UT-CU15-MOB-02.md) | Subir documento desde el detalle | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Foto procesada | 01, 02 |
| CA02 | Archivo grande | 02 |
| CA03 | Subida incompleta | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Corregir el texto OCR (web).
- Escaneos del expediente del docente (web).
