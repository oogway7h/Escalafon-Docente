# Unidades de trabajo — CU-13-MOB

**Caso de uso:** El docente radica una solicitud fotografiando o eligiendo los documentos que exige la categoría.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-13 - Radicar Solicitud desde el Portal Público`](../../../CasosDeUso/CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico.md)
**Recurso HTTP / contrato:** `GET /api/publico/categorias`, `POST /api/publico/solicitudes` (multipart)
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Cada requisito se cubre con **cámara**, **galería** (`image_picker`, JPG comprimido a calidad 85 para no pasar 8 MB) o **archivo** PDF (`file_picker`).
2. Antes de enviar se valida en el teléfono que haya un archivo por requisito y que ninguno pase 8 MB (misma regla que la API).

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /publico/categorias | 200 `[{id,nombre,aniosAntiguedad,requisitos:[{id,nombre}]}]` | — | público |
| POST | /publico/solicitudes multipart: `identificador`, `tipo` (CATEGORIZACION/REGULARIZACION), `categoriaId`, `doc_<requisitoId>` | 201 `{codigo}` | 422 `DOCUMENTOS_FALTANTES`, `ARCHIVO_INVALIDO`, `ARCHIVO_GRANDE`, `DOCENTE_INVALIDO`; 503 | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU13-MOB-01](./UT-CU13-MOB-01.md) | Categorías públicas y radicación en `PortalService` | — |
| [UT-CU13-MOB-02](./UT-CU13-MOB-02.md) | Pantalla de nueva solicitud | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Solicitud completa devuelve código | 01, 02 |
| CA02, CA03 | Faltantes y tamaño | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- OCR de los adjuntos del portal (la API los guarda como `NO_APLICA`).
- Seguimiento (CU-12), salvo navegar a él al terminar.
