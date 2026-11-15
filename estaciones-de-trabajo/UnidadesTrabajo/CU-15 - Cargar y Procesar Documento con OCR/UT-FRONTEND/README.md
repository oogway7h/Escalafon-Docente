# Unidades de trabajo — CU-15-FE

**Caso de uso:** Carga con progreso directo a S3, estado OCR con sondeo, visor y corrección.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-15 - Cargar y Procesar Documento con OCR`](../../../CasosDeUso/CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR.md)
**Recurso HTTP / contrato:** `/panel/tramites/[id]`, `/panel/documentos/[id]`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/tramites/[id] | Cargar documentos del trámite, ver OCR y corregir | `DOCUMENTO_GESTIONAR` |
| /panel/documentos/[id] | Escaneos del docente por tipo | `DOCUMENTO_ESCANEAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU15-FE-01](./UT-CU15-FE-01.md) | Carga de documentos del trámite | — |
| [UT-CU15-FE-02](./UT-CU15-FE-02.md) | Visor y corrección del texto OCR | 01 |
| [UT-CU15-FE-03](./UT-CU15-FE-03.md) | Expediente digital del docente | 01 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Pasos 1-4 | Subida | 01 |
| Pasos 5-7 | Estado y corrección | 02 |
| Expediente | Escaneos | 03 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Búsqueda global (CU-18).
