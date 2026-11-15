# Unidades de trabajo — CU-18-BE

**Caso de uso:** Búsqueda global y dentro del expediente con texto completo.
**Proyecto:** `backend`
**Fuente CU:** [`CU-18 - Buscar Documentación Digital`](../../../CasosDeUso/CU-18%20-%20Buscar%20Documentaci%C3%B3n%20Digital.md)
**Recurso HTTP / contrato:** `GET /api/documentos/buscar`, `GET /api/docentes/:id/documentos/buscar`, `GET /api/tramites/:id/documentos/:docId/buscar`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/documentos/buscar?q=&campo=&tipo=&pagina= | 200 `{datos,total,pagina,tamano}` con `paginas` y `fragmento` (⟦ ⟧) | 422 `CONSULTA_INVALIDA`, `FILTRO_INVALIDO` | `DOCUMENTO_BUSCAR` |
| GET | /api/docentes/:id/documentos/buscar?q= | 200 `{coincidencias}` | 422 | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU18-BE-01](./UT-CU18-BE-01.md) | Caso de uso BuscarDocumentos | — |
| [UT-CU18-BE-02](./UT-CU18-BE-02.md) | Rutas de búsqueda | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Búsqueda | 01, 02 |
| A1-A3 | Validación | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Indexación (CU-15).
