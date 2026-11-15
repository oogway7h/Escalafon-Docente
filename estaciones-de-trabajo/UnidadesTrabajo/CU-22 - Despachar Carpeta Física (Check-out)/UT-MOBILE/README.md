# Unidades de trabajo — CU-22-MOB

**Caso de uso:** Archivo registra la salida de la carpeta con fecha límite y estado material.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-22 - Despachar Carpeta Física (Check-out)`](../../../CasosDeUso/CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29.md)
**Recurso HTTP / contrato:** `POST /api/solicitudes-carpeta/:id/despachar`, `/entregar`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Desde Por atender se ofrece «Aprobar y entregar» (`/entregar`); desde Aprobadas, «Entregar» (`/despachar`). Igual que la web.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /solicitudes-carpeta/:id/despachar `{fechaLimite AAAA-MM-DD, estadoMaterial, observacion}` | 200 | 422 `FECHA_LIMITE_INVALIDA`, `ESTADO_MATERIAL_INVALIDO`; 409 `LEGAJO_NO_DISPONIBLE` | `CUSTODIA_REGISTRAR` |
| POST | /solicitudes-carpeta/:id/entregar (mismo cuerpo) | 200 | idem | `LEGAJO_GESTIONAR` + `CUSTODIA_REGISTRAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU22-MOB-01](./UT-CU22-MOB-01.md) | `despachar` y `entregar` en `CarpetasService` | — |
| [UT-CU22-MOB-02](./UT-CU22-MOB-02.md) | Hoja de entrega y pestaña Aprobadas | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Despacho | 01, 02 |
| CA02 | Fecha pasada | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Devolución (CU-23).
