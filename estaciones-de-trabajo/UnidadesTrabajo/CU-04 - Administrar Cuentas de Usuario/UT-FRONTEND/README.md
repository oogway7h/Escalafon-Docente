# Unidades de trabajo — CU-04-FE

**Caso de uso:** Listado, alta, edición, cambio de rol, baja y reenvío.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-04 - Administrar Cuentas de Usuario`](../../../CasosDeUso/CU-04%20-%20Administrar%20Cuentas%20de%20Usuario.md)
**Recurso HTTP / contrato:** `/panel/usuarios`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /panel/usuarios | Listado con búsqueda y acciones | `USUARIO_GESTIONAR` |
| /panel/usuarios/nuevo | Alta | idem |
| /panel/usuarios/[id]/editar | Edición | idem |
| /panel/usuarios/[id]/rol | Cambio de rol | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU04-FE-01](./UT-CU04-FE-01.md) | Listado de usuarios | — |
| [UT-CU04-FE-02](./UT-CU04-FE-02.md) | Formularios de usuario y de rol | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Listado y formularios | 01, 02 |
| A1-A5 | Errores y reenvío | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Roles (CU-05).
