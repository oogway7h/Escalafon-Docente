# Unidades de trabajo — CU-14-FE

**Caso de uso:** Página de envíos para el Director.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-14 - Notificar Cambios de Estado al Docente`](../../../CasosDeUso/CU-14%20-%20Notificar%20Cambios%20de%20Estado%20al%20Docente.md)
**Recurso HTTP / contrato:** `/panel/notificaciones`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/notificaciones | Avisos con estado, intentos y error; filtro por estado | `NOTIFICACION_CONSULTAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU14-FE-01](./UT-CU14-FE-01.md) | Página de notificaciones | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Paso 5 | Consulta | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Configuración del correo.
