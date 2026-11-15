# Unidades de trabajo — CU-08-BE

**Caso de uso:** Radicación interna, código de seguimiento y validación de solicitudes del portal.
**Proyecto:** `backend`
**Fuente CU:** [`CU-08 - Radicar Trámite`](../../../CasosDeUso/CU-08%20-%20Radicar%20Tr%C3%A1mite.md)
**Recurso HTTP / contrato:** `POST /api/tramites`, `POST /api/tramites/:id/validar`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/tramites `{docenteId, tipo, categoriaId}` | 201 resumen | 409 `DOCENTE_INACTIVO`, 422 `CATEGORIA_INVALIDA` | `TRAMITE_RADICAR` |
| POST | /api/tramites/:id/validar | 200 | 409 `YA_VALIDADO` | idem |
| GET | /api/tramites?estado=&buscar=&mios=&sinAsignar=&pagina= | 200 lista | — | lectura de trámites |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU08-BE-01](./UT-CU08-BE-01.md) | Agregado Trámite y máquina de estados | — |
| [UT-CU08-BE-02](./UT-CU08-BE-02.md) | Casos de uso Radicar, ValidarSolicitud y Listar | 01 |
| [UT-CU08-BE-03](./UT-CU08-BE-03.md) | Handler, DTO y rutas de trámites | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Radicar | 01-03 |
| A1-A4 | Validaciones | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Asignación (CU-09).
- Solicitud del portal (CU-13).
