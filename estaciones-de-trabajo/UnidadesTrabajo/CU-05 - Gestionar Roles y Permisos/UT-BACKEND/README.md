# Unidades de trabajo — CU-05-BE

**Caso de uso:** CRUD de roles con catálogo de permisos y protecciones.
**Proyecto:** `backend`
**Fuente CU:** [`CU-05 - Gestionar Roles y Permisos`](../../../CasosDeUso/CU-05%20-%20Gestionar%20Roles%20y%20Permisos.md)
**Recurso HTTP / contrato:** `/api/roles`, `/api/permisos`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/permisos | 200 catálogo | — | `ROL_GESTIONAR` |
| GET/POST | /api/roles | 200/201 | 409 `ROL_DUPLICADO` | idem |
| PUT | /api/roles/:id | 200 | 409 `ROL_SISTEMA_PROTEGIDO` | idem |
| DELETE | /api/roles/:id | 204 | 409 `ROL_CON_USUARIOS` | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU05-BE-01](./UT-CU05-BE-01.md) | Dominio de rol, permisos y seed | — |
| [UT-CU05-BE-02](./UT-CU05-BE-02.md) | Casos de uso de roles | 01 |
| [UT-CU05-BE-03](./UT-CU05-BE-03.md) | Handler y rutas de roles | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | CRUD | 02, 03 |
| A1-A4 | Protecciones | 01, 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Asignar rol a un usuario (CU-04).
