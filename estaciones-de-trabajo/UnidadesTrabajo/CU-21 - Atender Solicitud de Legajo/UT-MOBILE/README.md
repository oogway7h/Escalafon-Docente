# Unidades de trabajo — CU-21-MOB

**Caso de uso:** Archivo aprueba o rechaza (con motivo) las solicitudes pendientes.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-21 - Atender Solicitud de Legajo`](../../../CasosDeUso/CU-21%20-%20Atender%20Solicitud%20de%20Legajo.md)
**Recurso HTTP / contrato:** `GET /api/solicitudes-carpeta?estado=PENDIENTE`, `POST /api/solicitudes-carpeta/:id/atender`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Una pantalla «Custodia» con pestañas (Por atender, Aprobadas, Prestadas) como `/panel/custodia` de la web; esta UT crea la pantalla y su primera pestaña.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /solicitudes-carpeta/:id/atender `{aprobar, nota}` | 200 | 422 `NOTA_REQUERIDA`; 409 `SOLICITUD_ESTADO_INVALIDO` | `LEGAJO_GESTIONAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU21-MOB-01](./UT-CU21-MOB-01.md) | `atender` en `CarpetasService` | — |
| [UT-CU21-MOB-02](./UT-CU21-MOB-02.md) | Pantalla Custodia: pestaña Por atender | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Aprobar | 01, 02 |
| CA02 | Rechazo sin nota | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Despachar (CU-22) y devolver (CU-23), que agregan acciones y pestañas.
