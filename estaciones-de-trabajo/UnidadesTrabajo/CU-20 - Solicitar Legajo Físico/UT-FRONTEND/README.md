# Unidades de trabajo — CU-20-FE

**Caso de uso:** Mis solicitudes y nueva solicitud ligada a un trámite.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-20 - Solicitar Legajo Físico`](../../../CasosDeUso/CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico.md)
**Recurso HTTP / contrato:** `/panel/carpetas`, `/panel/carpetas/nueva`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/carpetas | Mis solicitudes con estado | `CARPETA_SOLICITAR` |
| /panel/carpetas/nueva | Elegir trámite y carpeta | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU20-FE-01](./UT-CU20-FE-01.md) | Mis solicitudes y nueva solicitud | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Solicitar | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Custodia.
