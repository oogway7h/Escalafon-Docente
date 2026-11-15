# Unidades de trabajo — CU-13-BE

**Caso de uso:** Radicación pública multipart con adjuntos a S3 y catálogo público.
**Proyecto:** `backend`
**Fuente CU:** [`CU-13 - Radicar Solicitud desde el Portal Público`](../../../CasosDeUso/CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico.md)
**Recurso HTTP / contrato:** `POST /api/publico/solicitudes`, `GET /api/publico/categorias`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/publico/categorias | 200 `[{id,nombre,aniosAntiguedad,requisitos:[{id,nombre}]}]` | — | público |
| POST | /api/publico/solicitudes (multipart: `identificador`, `tipo`, `categoriaId`, `doc_<requisitoId>`) | 201 `{codigo}` | 422 `DOCUMENTOS_FALTANTES`, `ARCHIVO_INVALIDO`, `ARCHIVO_GRANDE`, `DOCENTE_INVALIDO`; 503 | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU13-BE-01](./UT-CU13-BE-01.md) | Puerto ObjectStorage y adaptador S3 | — |
| [UT-CU13-BE-02](./UT-CU13-BE-02.md) | Caso de uso RadicarPublico | 01 |
| [UT-CU13-BE-03](./UT-CU13-BE-03.md) | Rutas públicas de solicitud y categorías | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Radicar con adjuntos | 02, 03 |
| A3-A5 | Archivos | 01, 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- OCR de los adjuntos (no aplica: `NO_APLICA`).
