# Unidades de trabajo — CU-10-MOB

**Caso de uso:** El analista consulta su bandeja y el detalle de cada trámite (solo lectura en el móvil).
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-10 - Revisar Expediente y Registrar Observaciones`](../../../CasosDeUso/CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones.md)
**Recurso HTTP / contrato:** `GET /api/tramites?mios=true`, `GET /api/tramites/:id`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. En el móvil la revisión es de **consulta**: observar, dictaminar y remitir siguen en la web, donde está el visor de documentos. Lo que el móvil agrega es ver la bandeja en cualquier lugar y, en otros CU, subir fotos y pedir carpetas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /tramites?mios=true | 200 `[{id,codigo,docenteNombre,docenteCodigo,categoria,tipo,estado,origen,pendienteValidacion,analistaNombre,creadoEn}]` | 401, 403 | `TRAMITE_REVISAR` (lectura) |
| GET | /tramites/:id | 200 detalle con `historial`, `observaciones`, `documentos`, `requisitos` | 404 | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU10-MOB-01](./UT-CU10-MOB-01.md) | Modelos y `TramitesService` (lectura) | — |
| [UT-CU10-MOB-02](./UT-CU10-MOB-02.md) | Bandeja y detalle del trámite | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Pasos 1-2 | Bandeja y detalle | 01, 02 |

## Follow-ups (no se implementan en estas UT)

- Registrar observaciones desde el móvil.

## Qué NO incluyen estas UT

- Iniciar revisión, observar o dictaminar (se hacen en la web).
- Subir documentos (CU-15) y pedir carpeta (CU-20), que editan esta pantalla.
