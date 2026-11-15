# Unidades de trabajo — CU-03-FE

**Caso de uso:** Página de perfil con dos formularios.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-03 - Editar Perfil Personal`](../../../CasosDeUso/CU-03%20-%20Editar%20Perfil%20Personal.md)
**Recurso HTTP / contrato:** `/panel/perfil`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/perfil | Datos, contacto editable y cambio de contraseña | `PERFIL_GESTIONAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU03-FE-01](./UT-CU03-FE-01.md) | Perfil con contacto y contraseña | — |

**Cadena:** 01

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal y A1-A3 | Formularios con errores de la API | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Administración de usuarios (CU-04).
