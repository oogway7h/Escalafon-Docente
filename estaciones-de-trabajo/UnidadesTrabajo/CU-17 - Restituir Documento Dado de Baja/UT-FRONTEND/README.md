# Unidades de trabajo — CU-17-FE

**Caso de uso:** Bandeja de bajas con detalle y restitución.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-17 - Restituir Documento Dado de Baja`](../../../CasosDeUso/CU-17%20-%20Restituir%20Documento%20Dado%20de%20Baja.md)
**Recurso HTTP / contrato:** `/panel/documentos/bajas`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/documentos/bajas | Bandeja paginada, detalle `?doc=`, restituir con motivo | `DOCUMENTO_RESTITUIR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU17-FE-01](./UT-CU17-FE-01.md) | Bandeja de bajas | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Restituir | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- —
