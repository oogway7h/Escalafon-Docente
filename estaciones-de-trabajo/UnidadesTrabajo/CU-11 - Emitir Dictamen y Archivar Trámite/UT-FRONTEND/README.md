# Unidades de trabajo — CU-11-FE

**Caso de uso:** Acciones del flujo en el detalle y bandeja «Por archivar».
**Proyecto:** `frontend`
**Fuente CU:** [`CU-11 - Emitir Dictamen y Archivar Trámite`](../../../CasosDeUso/CU-11%20-%20Emitir%20Dictamen%20y%20Archivar%20Tr%C3%A1mite.md)
**Recurso HTTP / contrato:** `/panel/tramites/[id]`, `/panel/tramites/archivar`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/tramites/[id] | Años a evaluar, dictamen, remitir, regreso, remitir al Rectorado | `TRAMITE_REVISAR` |
| /panel/tramites/archivar | Remitidos por archivar | `TRAMITE_ARCHIVAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU11-FE-01](./UT-CU11-FE-01.md) | Formularios del flujo en el detalle | — |
| [UT-CU11-FE-02](./UT-CU11-FE-02.md) | Bandeja Por archivar | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Pasos 1-5 | Formularios del flujo | 01 |
| Paso 6 | Archivar | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Notificaciones.
