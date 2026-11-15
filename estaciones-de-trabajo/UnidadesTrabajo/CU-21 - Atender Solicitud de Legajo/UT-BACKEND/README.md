# Unidades de trabajo — CU-21-BE

**Caso de uso:** Atención de solicitudes.
**Proyecto:** `backend`
**Fuente CU:** [`CU-21 - Atender Solicitud de Legajo`](../../../CasosDeUso/CU-21%20-%20Atender%20Solicitud%20de%20Legajo.md)
**Recurso HTTP / contrato:** `POST /api/solicitudes-carpeta/:id/atender`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/solicitudes-carpeta/:id/atender `{aprobar, nota}` | 200 solicitud | 422 `NOTA_REQUERIDA`; 409 `SOLICITUD_ESTADO_INVALIDO` | `LEGAJO_GESTIONAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU21-BE-01](./UT-CU21-BE-01.md) | Caso de uso Atender | — |
| [UT-CU21-BE-02](./UT-CU21-BE-02.md) | Ruta de atención | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal y A1-A2 | Atender | 01, 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Despacho (CU-22).
