# Unidades de trabajo — CU-13-FE

**Caso de uso:** Formulario público con requisitos por categoría.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-13 - Radicar Solicitud desde el Portal Público`](../../../CasosDeUso/CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico.md)
**Recurso HTTP / contrato:** `/solicitud`, portada `app/page.tsx`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /solicitud | Identificación, tipo, categoría y un archivo por requisito | público |
| / | Requisitos por categoría (tarjetas) | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU13-FE-01](./UT-CU13-FE-01.md) | Formulario de solicitud | — |
| [UT-CU13-FE-02](./UT-CU13-FE-02.md) | Requisitos por categoría en la portada | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal, A1-A5 | Formulario | 01 |
| Paso 3 | Requisitos visibles antes de empezar | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Seguimiento (CU-12).
