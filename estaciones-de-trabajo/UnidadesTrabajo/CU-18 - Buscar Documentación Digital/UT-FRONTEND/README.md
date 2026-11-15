# Unidades de trabajo — CU-18-FE

**Caso de uso:** Buscador global y dentro del expediente, con visor en la página hallada.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-18 - Buscar Documentación Digital`](../../../CasosDeUso/CU-18%20-%20Buscar%20Documentaci%C3%B3n%20Digital.md)
**Recurso HTTP / contrato:** `/panel/documentos/buscar`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/documentos/buscar | Búsqueda global; `?docente=` dentro del expediente; visor con `?doc=&pagina=&q=` | `DOCUMENTO_BUSCAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU18-FE-01](./UT-CU18-FE-01.md) | Buscador global | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Buscar y abrir | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- —
