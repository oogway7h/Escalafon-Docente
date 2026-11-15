# Unidades de trabajo — CU-16-FE

**Caso de uso:** Formularios de baja en el detalle del trámite y del expediente.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-16 - Dar de Baja Documento`](../../../CasosDeUso/CU-16%20-%20Dar%20de%20Baja%20Documento.md)
**Recurso HTTP / contrato:** `/panel/tramites/[id]`, `/panel/documentos/[id]`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/tramites/[id] | Tarjeta «Dar de baja un documento» | `DOCUMENTO_GESTIONAR` |
| /panel/documentos/[id] | Baja de escaneo | `DOCUMENTO_ESCANEAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU16-FE-01](./UT-CU16-FE-01.md) | Formularios de baja | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal y A1-A4 | Formularios | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Restitución (CU-17).
