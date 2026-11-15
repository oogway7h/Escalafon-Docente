# Unidades de trabajo — CU-22-BE

**Caso de uso:** Despacho con transacción y bitácora inmutable.
**Proyecto:** `backend`
**Fuente CU:** [`CU-22 - Despachar Carpeta Física (Check-out)`](../../../CasosDeUso/CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29.md)
**Recurso HTTP / contrato:** `POST /api/solicitudes-carpeta/:id/despachar`, `/entregar`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. El cambio EN_ARCHIVO → PRESTADO es un `UPDATE … WHERE estado='EN_ARCHIVO'` dentro de `TxManager`: la base resuelve la concurrencia.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/solicitudes-carpeta/:id/despachar `{fechaLimite, estadoMaterial, observacion}` | 200 | 422 `FECHA_LIMITE_INVALIDA`, `ESTADO_MATERIAL_INVALIDO`; 409 `LEGAJO_NO_DISPONIBLE` | `CUSTODIA_REGISTRAR` |
| POST | /api/solicitudes-carpeta/:id/entregar (mismo cuerpo) | 200 | idem | `LEGAJO_GESTIONAR` + `CUSTODIA_REGISTRAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU22-BE-01](./UT-CU22-BE-01.md) | TxManager y movimiento de custodia | — |
| [UT-CU22-BE-02](./UT-CU22-BE-02.md) | Casos de uso Despachar y Entregar | 01 |
| [UT-CU22-BE-03](./UT-CU22-BE-03.md) | Rutas de despacho | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Despacho | 01-03 |
| A1-A4 | Reglas | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Devolución (CU-23).
