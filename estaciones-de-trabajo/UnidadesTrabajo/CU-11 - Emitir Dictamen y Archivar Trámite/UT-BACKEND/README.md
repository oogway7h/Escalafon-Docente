# Unidades de trabajo — CU-11-BE

**Caso de uso:** Dictamen, periodo, remisiones, regreso y archivo.
**Proyecto:** `backend`
**Fuente CU:** [`CU-11 - Emitir Dictamen y Archivar Trámite`](../../../CasosDeUso/CU-11%20-%20Emitir%20Dictamen%20y%20Archivar%20Tr%C3%A1mite.md)
**Recurso HTTP / contrato:** `/api/tramites/:id/{periodo,dictamen,remitir-facultad,regreso,remitir-rectorado,archivar}`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. El flujo se amplió al real del Departamento (En facultad, Revisión final, Remitido) en la migración `flujo_escalafon`.
2. La reasignación sigue permitida hasta Remitido.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| PUT | /periodo `{anioDesde, anioHasta}` | 200 | 422 `PERIODO_INVALIDO` | `TRAMITE_REVISAR` |
| POST | /dictamen `{resultado, texto}` | 200 | 409 `OBSERVACIONES_PENDIENTES`, `PERIODO_REQUERIDO` | idem |
| POST | /remitir-facultad `{nota}` | 200 | 409 `TRANSICION_INVALIDA` | idem |
| POST | /regreso `{resolucionNumero, resolucionFecha}` | 200 | 422 `RESOLUCION_REQUERIDA` | `TRAMITE_REVISAR` o `TRAMITE_RADICAR` |
| POST | /remitir-rectorado `{nota}` | 200 | 409 | `TRAMITE_REVISAR` |
| POST | /archivar `{ubicacion}` | 200 | 409 `CARPETAS_ABIERTAS` | `TRAMITE_ARCHIVAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU11-BE-01](./UT-CU11-BE-01.md) | Caso de uso Dictaminar | — |
| [UT-CU11-BE-02](./UT-CU11-BE-02.md) | Remisiones, regreso y archivo | 01 |
| [UT-CU11-BE-03](./UT-CU11-BE-03.md) | Rutas del flujo | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Pasos 1-2 | Periodo y dictamen | 01 |
| Pasos 3-5 | Remisiones y regreso | 02 |
| Paso 6 | Archivo | 02 |
| Rutas | HTTP | 03 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Notificaciones al docente (CU-14).
