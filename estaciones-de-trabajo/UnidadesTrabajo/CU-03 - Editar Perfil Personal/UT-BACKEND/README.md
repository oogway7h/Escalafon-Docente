# Unidades de trabajo — CU-03-BE

**Caso de uso:** Lectura y edición del perfil propio y cambio de contraseña.
**Proyecto:** `backend`
**Fuente CU:** [`CU-03 - Editar Perfil Personal`](../../../CasosDeUso/CU-03%20-%20Editar%20Perfil%20Personal.md)
**Recurso HTTP / contrato:** `GET/PUT /api/perfil`, `PUT /api/perfil/contrasena`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/perfil | 200 perfil | 401 | `PERFIL_GESTIONAR` |
| PUT | /api/perfil `{telefono,…}` | 200 | 422 `CAMPO_NO_EDITABLE` | `PERFIL_GESTIONAR` |
| PUT | /api/perfil/contrasena `{actual, nueva}` | 204 | 422 `CONTRASENA_ACTUAL_INCORRECTA`, `POLITICA_CONTRASENA` | `PERFIL_GESTIONAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU03-BE-01](./UT-CU03-BE-01.md) | Casos de uso del perfil | — |
| [UT-CU03-BE-02](./UT-CU03-BE-02.md) | Handler y rutas de perfil | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Ver y editar contacto | 01, 02 |
| A1, A2 | Contraseña | 01 |
| A3 | Campos no editables | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Edición de otros usuarios (CU-04).
