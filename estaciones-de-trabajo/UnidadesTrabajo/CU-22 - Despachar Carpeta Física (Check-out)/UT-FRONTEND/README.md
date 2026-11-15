# Unidades de trabajo — CU-22-FE

**Caso de uso:** «Aprobar y entregar carpeta» en la pantalla de custodia.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-22 - Despachar Carpeta Física (Check-out)`](../../../CasosDeUso/CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29.md)
**Recurso HTTP / contrato:** `/panel/custodia`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/custodia | Fecha límite con atajos, estado material, observación | `CUSTODIA_REGISTRAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU22-FE-01](./UT-CU22-FE-01.md) | Formulario de entrega | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Despacho | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- —
