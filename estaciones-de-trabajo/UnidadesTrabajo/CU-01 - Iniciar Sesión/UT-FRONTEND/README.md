# Unidades de trabajo — CU-01-FE

**Caso de uso:** Pantalla de login, cookie httpOnly con el token, guardia de rutas y menú por permisos.
**Proyecto:** `frontend`
**Fuente CU:** [`CU-01 - Iniciar Sesión`](../../../CasosDeUso/CU-01%20-%20Iniciar%20Sesi%C3%B3n.md)
**Recurso HTTP / contrato:** `/login`, `app/panel/layout.tsx`, `proxy.ts`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Ruta | Qué hace | Permiso |
|---|---|---|
| /login | Formulario correo + contraseña; errores de la API traducidos | público |
| /panel | Portada con tarjetas por permiso | sesión |
| proxy.ts | Sin cookie en `/panel/*` → `/login` | — |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU01-FE-01](./UT-CU01-FE-01.md) | Cliente de API con sesión y guardia de rutas | — |
| [UT-CU01-FE-02](./UT-CU01-FE-02.md) | Pantalla y acción de login | 01 |
| [UT-CU01-FE-03](./UT-CU01-FE-03.md) | Layout del panel, menú por permisos y cierre de sesión | 01, 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Login y menú | 01, 02, 03 |
| A1-A4 | Mensajes de error | 02 |
| A6 | Cerrar sesión | 03 |
| CA06 | Con sesión, /login lleva al panel | 01 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Pantallas de recuperación y activación (CU-02).
