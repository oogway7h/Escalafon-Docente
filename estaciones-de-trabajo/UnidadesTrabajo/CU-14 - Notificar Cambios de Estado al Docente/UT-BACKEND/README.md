# Unidades de trabajo — CU-14-BE

**Caso de uso:** Avisos por correo con reintentos y consulta de envíos.
**Proyecto:** `backend`
**Fuente CU:** [`CU-14 - Notificar Cambios de Estado al Docente`](../../../CasosDeUso/CU-14%20-%20Notificar%20Cambios%20de%20Estado%20al%20Docente.md)
**Recurso HTTP / contrato:** `GET /api/notificaciones`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. `tramites.Notificador` es una interfaz opcional inyectada en los casos de uso del flujo: sin ella el flujo funciona igual.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/notificaciones?estado=&pagina= | 200 `{datos,total,pagina,tamano}` | — | `NOTIFICACION_CONSULTAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU14-BE-01](./UT-CU14-BE-01.md) | Dominio y tabla de envíos | — |
| [UT-CU14-BE-02](./UT-CU14-BE-02.md) | Caso de uso Notificar con reintentos | 01 |
| [UT-CU14-BE-03](./UT-CU14-BE-03.md) | Ruta de consulta de envíos | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Aviso y envío | 01, 02 |
| A1-A3 | Reintentos | 02 |
| Paso 5 | Consulta | 03 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Notificaciones push en el móvil (follow-up de la app).
