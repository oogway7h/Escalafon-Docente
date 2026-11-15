# Unidades de trabajo — CU-12-BE

**Caso de uso:** Consulta pública por código o por código de docente + cédula, con limitador.
**Proyecto:** `backend`
**Fuente CU:** [`CU-12 - Consultar Estado de Trámite (Portal Público)`](../../../CasosDeUso/CU-12%20-%20Consultar%20Estado%20de%20Tr%C3%A1mite%20%28Portal%20P%C3%BAblico%29.md)
**Recurso HTTP / contrato:** `GET /api/publico/tramites/:codigo`, `POST /api/publico/consulta-docente`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/publico/tramites/:codigo | 200 `{codigo,tipo,categoria,estado,radicadoEn,pendienteValidacion,resultado,etapas,observaciones}` | 404, 429 | público |
| POST | /api/publico/consulta-docente `{codigo, cedula}` | 200 `[{codigo,tipo,categoria,estado,radicadoEn,pendienteValidacion}]` | 422, 429 | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU12-BE-01](./UT-CU12-BE-01.md) | Casos de uso ConsultarPublico y ConsultarPublicoPorDocente | — |
| [UT-CU12-BE-02](./UT-CU12-BE-02.md) | Rutas públicas con limitador | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Consulta por código | 01, 02 |
| Paso 4, A2 | Consulta por docente | 01 |
| A3 | Limitador | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Radicar desde el portal (CU-13).
