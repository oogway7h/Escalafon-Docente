# Unidades de trabajo — CU-21-FE

**Caso de uso:** Pestaña Por atender de la pantalla de custodia.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-21 - Atender Solicitud de Legajo`](../../../CasosDeUso/CU-21%20-%20Atender%20Solicitud%20de%20Legajo.md)
**Recurso HTTP / contrato:** `/panel/custodia`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/custodia?tab=atender | Lista y panel de acción: aprobar, rechazar con motivo | `LEGAJO_GESTIONAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU21-FE-01](./UT-CU21-FE-01.md) | Pestaña Por atender | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Atender | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Bitácora.
