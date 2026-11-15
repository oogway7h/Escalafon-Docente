# Unidades de trabajo — CU-10-FE

**Caso de uso:** Mi bandeja y detalle del trámite con revisión y observaciones.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-10 - Revisar Expediente y Registrar Observaciones`](../../../CasosDeUso/CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones.md)
**Recurso HTTP / contrato:** `/panel/tramites/revisar`, `/panel/tramites/[id]`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/tramites/revisar | Mi bandeja | `TRAMITE_REVISAR` |
| /panel/tramites/[id] | Detalle, fases, requisitos, observaciones | lectura |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU10-FE-01](./UT-CU10-FE-01.md) | Mi bandeja | — |
| [UT-CU10-FE-02](./UT-CU10-FE-02.md) | Detalle del trámite y línea de fases | 01 |
| [UT-CU10-FE-03](./UT-CU10-FE-03.md) | Acciones de revisión y observaciones | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Bandeja y detalle | 01, 02 |
| Paso 4-5 | Observaciones | 03 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Dictamen y remisiones (CU-11).
