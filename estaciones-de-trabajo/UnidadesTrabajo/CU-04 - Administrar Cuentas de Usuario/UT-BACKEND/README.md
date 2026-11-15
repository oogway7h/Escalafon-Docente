# Unidades de trabajo — CU-04-BE

**Caso de uso:** CRUD de usuarios con baja lógica, cambio de rol y reenvío de activación.
**Proyecto:** `backend`
**Fuente CU:** [`CU-04 - Administrar Cuentas de Usuario`](../../../CasosDeUso/CU-04%20-%20Administrar%20Cuentas%20de%20Usuario.md)
**Recurso HTTP / contrato:** `/api/usuarios`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/usuarios | 200 lista | — | `USUARIO_GESTIONAR` |
| POST | /api/usuarios | 201 | 409 `CORREO_DUPLICADO`, 422 `ROL_INVALIDO` | idem |
| PUT | /api/usuarios/:id | 200 | 409 | idem |
| PATCH | /api/usuarios/:id/rol | 200 | 409 `AUTOCAMBIO_ROL` | idem |
| PATCH | /api/usuarios/:id/estado | 200 | 409 `AUTODESACTIVACION` | idem |
| POST | /api/usuarios/:id/reenviar-activacion | 202 | 409 `CUENTA_YA_ACTIVADA` | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU04-BE-01](./UT-CU04-BE-01.md) | Repositorio de usuarios y unicidad de correo | — |
| [UT-CU04-BE-02](./UT-CU04-BE-02.md) | Casos de uso de usuarios | 01 |
| [UT-CU04-BE-03](./UT-CU04-BE-03.md) | Handler y rutas de usuarios | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Alta, edición, baja | 01, 02, 03 |
| A1-A4 | Reglas de unicidad y autoprotección | 02 |
| A5 | Reenvío | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Definición de roles y permisos (CU-05).
