# Unidades de trabajo — CU-23-BE

**Caso de uso:** Devolución con estado material y bitácora.
**Proyecto:** `backend`
**Fuente CU:** [`CU-23 - Registrar Devolución de Carpeta (Check-in)`](../../../CasosDeUso/CU-23%20-%20Registrar%20Devoluci%C3%B3n%20de%20Carpeta%20%28Check-in%29.md)
**Recurso HTTP / contrato:** `POST /api/legajos/:id/devolver`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/legajos/:id/devolver `{estadoMaterial, observacion}` | 200 | 422 `OBSERVACION_REQUERIDA`, `ESTADO_MATERIAL_INVALIDO`; 409 `LEGAJO_NO_PRESTADO` | `CUSTODIA_REGISTRAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU23-BE-01](./UT-CU23-BE-01.md) | Caso de uso Devolver | — |
| [UT-CU23-BE-02](./UT-CU23-BE-02.md) | Ruta de devolución | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal y A1-A3 | Devolución | 01, 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Foto de evidencia (eliminada en `quitar_evidencia_custodia`).
