# Unidades de trabajo — CU-01-BE

**Caso de uso:** Autenticación por correo y contraseña, bloqueo por intentos, emisión y revocación del JWT.
**Proyecto:** `backend`
**Fuente CU:** [`CU-01 - Iniciar Sesión`](../../../CasosDeUso/CU-01%20-%20Iniciar%20Sesi%C3%B3n.md)
**Recurso HTTP / contrato:** `POST /api/auth/login`, `POST /api/auth/logout`, `GET /api/auth/sesion`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. JWT HS256 firmado localmente + `version_sesion` en Redis: revocar todas las sesiones de un usuario es subir su versión (lo usa CU-02).
2. El limitador de tasa de `/auth/*` necesita Redis: sin Redis la API no arranca a propósito.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/auth/login `{correo, contrasena}` | 200 `{token, expiraEn, usuario, permisos}` | 401 `CREDENCIALES_INVALIDAS`, 423 `CUENTA_BLOQUEADA`, 422 `VALIDACION`, 429 | público |
| POST | /api/auth/logout | 204 | 401 | sesión |
| GET | /api/auth/sesion | 200 `{usuario, rol, permisos}` | 401 | sesión |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU01-BE-01](./UT-CU01-BE-01.md) | Entidad Usuario y reglas de bloqueo | — |
| [UT-CU01-BE-02](./UT-CU01-BE-02.md) | Puertos de seguridad y sus implementaciones | 01 |
| [UT-CU01-BE-03](./UT-CU01-BE-03.md) | Caso de uso AutenticarUsuario | 01, 02 |
| [UT-CU01-BE-04](./UT-CU01-BE-04.md) | CerrarSesion y ObtenerSesion + middleware Autenticar | 03 |
| [UT-CU01-BE-05](./UT-CU01-BE-05.md) | Handler y rutas de autenticación | 03, 04 |

**Cadena:** 01 → 02 → 03 → 04 → 05

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Validar credenciales y emitir sesión | 02, 03, 05 |
| A1, A2, A3 | Credenciales inválidas, bloqueo, cuenta inactiva | 01, 03 |
| A5 | Limitador de tasa | 05 |
| A6 | Cerrar sesión y revocación | 04, 05 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Recuperación y activación de contraseña (CU-02).
- Alta de cuentas y asignación de roles (CU-04, CU-05).
