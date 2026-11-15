# Unidades de trabajo — CU-12-MOB

**Caso de uso:** El docente consulta su trámite desde el teléfono, sin cuenta.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-12 - Consultar Estado de Trámite (Portal Público)`](../../../CasosDeUso/CU-12%20-%20Consultar%20Estado%20de%20Tr%C3%A1mite%20%28Portal%20P%C3%BAblico%29.md)
**Recurso HTTP / contrato:** `GET /api/publico/tramites/:codigo`, `POST /api/publico/consulta-docente`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Mismo criterio que la web: más de 5 caracteres = código de seguimiento; hasta 5 = código de docente y se pide la cédula (que va en el cuerpo).
2. Las etiquetas y descripciones de estado se copian de `frontend/lib/tramites.ts` a `lib/shared/etiquetas.dart` (las usarán también los trámites internos).

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /publico/tramites/:codigo | 200 `{codigo,tipo,categoria,estado,radicadoEn,pendienteValidacion,resultado,etapas:[{estado,fecha}],observaciones:[{documento,motivo,plazo,resuelta}]}` | 404, 429 | público |
| POST | /publico/consulta-docente `{codigo, cedula}` | 200 `[{codigo,tipo,categoria,estado,radicadoEn,pendienteValidacion}]` | 422, 429 | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU12-MOB-01](./UT-CU12-MOB-01.md) | Modelos y `PortalService` (consulta) | — |
| [UT-CU12-MOB-02](./UT-CU12-MOB-02.md) | Pantalla de seguimiento | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01, CA04 | Estado, etapas y observaciones | 01, 02 |
| CA02 | Código inexistente | 02 |
| CA03 | Cédula incorrecta → lista vacía | 01, 02 |

## Follow-ups (no se implementan en estas UT)

- Leer el código desde un QR impreso en el comprobante de ventanilla.

## Qué NO incluyen estas UT

- Radicar una solicitud (CU-13).
- Ningún dato personal: la API no los entrega.
