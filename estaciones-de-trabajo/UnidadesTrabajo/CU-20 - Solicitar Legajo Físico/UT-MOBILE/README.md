# Unidades de trabajo — CU-20-MOB

**Caso de uso:** El analista pide la carpeta física del docente de su trámite y sigue el estado.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-20 - Solicitar Legajo Físico`](../../../CasosDeUso/CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico.md)
**Recurso HTTP / contrato:** `POST /api/legajos/:id/solicitudes`, `GET /api/solicitudes-carpeta`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Las carpetas del docente se obtienen con `GET /legajos?buscar=<docenteCodigo>` (búsqueda exacta por código).

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /legajos/:id/solicitudes `{tramiteId, motivo}` | 201 solicitud | 409 `LEGAJO_NO_DISPONIBLE`, `SOLICITUD_ABIERTA`; 422 `TRAMITE_INVALIDO` | `CARPETA_SOLICITAR` |
| GET | /solicitudes-carpeta?estado=&pagina=&tamano= | 200 `{datos:[{id,legajoId,legajoCodigo,docenteNombre,tramiteId,tramiteCodigo,solicitante,motivo,estado,respuesta,creadoEn,atendidoEn,fechaLimite}],…}` | — | lectura (el analista ve las suyas) |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU20-MOB-01](./UT-CU20-MOB-01.md) | Modelo `SolicitudCarpeta` y `CarpetasService` (solicitar y listar) | — |
| [UT-CU20-MOB-02](./UT-CU20-MOB-02.md) | Mis solicitudes y solicitar desde el trámite | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Solicitud Pendiente | 01, 02 |
| CA02 | Solicitud abierta | 02 |
| CA03 | Mis solicitudes | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Atender (CU-21).
