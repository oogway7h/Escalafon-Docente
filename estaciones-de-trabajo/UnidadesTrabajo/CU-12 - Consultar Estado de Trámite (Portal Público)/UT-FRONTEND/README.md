# Unidades de trabajo — CU-12-FE

**Caso de uso:** Página pública de seguimiento.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-12 - Consultar Estado de Trámite (Portal Público)`](../../../CasosDeUso/CU-12%20-%20Consultar%20Estado%20de%20Tr%C3%A1mite%20%28Portal%20P%C3%BAblico%29.md)
**Recurso HTTP / contrato:** `/seguimiento`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /seguimiento | Un campo de código; con ≤ 5 caracteres pide la cédula; tarjeta de estado, fases, observaciones e historial | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU12-FE-01](./UT-CU12-FE-01.md) | Página de seguimiento | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal, A1-A4 | Seguimiento | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Solicitud (CU-13).
