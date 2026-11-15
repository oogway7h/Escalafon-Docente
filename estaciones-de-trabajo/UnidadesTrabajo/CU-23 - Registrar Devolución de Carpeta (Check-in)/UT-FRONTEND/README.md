# Unidades de trabajo — CU-23-FE

**Caso de uso:** Registrar devolución en las pestañas Prestadas y Vencidas.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-23 - Registrar Devolución de Carpeta (Check-in)`](../../../CasosDeUso/CU-23%20-%20Registrar%20Devoluci%C3%B3n%20de%20Carpeta%20%28Check-in%29.md)
**Recurso HTTP / contrato:** `/panel/custodia`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/custodia?tab=prestadas|vencidas | Registrar devolución | `CUSTODIA_REGISTRAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU23-FE-01](./UT-CU23-FE-01.md) | Formulario de devolución | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Devolución | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- —
