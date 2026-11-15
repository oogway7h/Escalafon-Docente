# Unidades de trabajo — CU-20-BE

**Caso de uso:** Creación y listado de solicitudes de carpeta.
**Proyecto:** `backend`
**Fuente CU:** [`CU-20 - Solicitar Legajo Físico`](../../../CasosDeUso/CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico.md)
**Recurso HTTP / contrato:** `POST /api/legajos/:id/solicitudes`, `GET /api/solicitudes-carpeta`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/legajos/:id/solicitudes `{tramiteId, motivo}` | 201 solicitud | 409 `LEGAJO_NO_DISPONIBLE`, `SOLICITUD_ABIERTA`; 422 `TRAMITE_INVALIDO` | `CARPETA_SOLICITAR` |
| GET | /api/solicitudes-carpeta?estado=&pagina= | 200 `{datos:[{id,legajoId,legajoCodigo,docenteNombre,tramiteId,tramiteCodigo,solicitante,motivo,estado,respuesta,creadoEn,atendidoEn,fechaLimite}],…}` | — | lectura (el analista ve las suyas) |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU20-BE-01](./UT-CU20-BE-01.md) | Caso de uso Solicitar y ListarSolicitudes | — |
| [UT-CU20-BE-02](./UT-CU20-BE-02.md) | Rutas de solicitud | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Solicitar y seguir | 01, 02 |
| A1-A4 | Reglas | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Atención (CU-21).
