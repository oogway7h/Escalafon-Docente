# Unidades de trabajo — CU-01-MOB

**Caso de uso:** Base de la app (cliente HTTP, tema, router) e inicio de sesión del personal con sesión guardada y accesos por permiso.
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-01 - Iniciar Sesión`](../../../CasosDeUso/CU-01%20-%20Iniciar%20Sesi%C3%B3n.md)
**Recurso HTTP / contrato:** `POST /api/auth/login`, `GET /api/auth/sesion`, `POST /api/auth/logout`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Esta es la **primera UT-MOBILE**: crea el proyecto Flutter y `core/` (UT 01). Ningún otro CU lo vuelve a crear.
2. Token solo en `flutter_secure_storage`. Al abrir la app se valida con `GET /auth/sesion`: si responde 401 (vencido o revocado) se borra y se va a la bienvenida. No hay refresh token en la API.
3. `sesionProvider` (AsyncNotifier<Sesion?>) es la única fuente de verdad; el router lo escucha con `refreshListenable`.
4. Accesos del inicio calculados por permisos con una función pura `accesosPara(permisos)`: un acceso aparece recién cuando su CU-MOB se implementa (cada CU edita esa lista).

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /auth/login `{correo, contrasena}` | 200 `{token, expiraEn (segundos), usuario:{id,nombres,apellidos,correo,rolId}, permisos:[…]}` | 401 `CREDENCIALES_INVALIDAS`, 423 `CUENTA_BLOQUEADA`, 422, 429 | público |
| GET | /auth/sesion | 200 `{usuario, rol, permisos}` | 401 | Bearer |
| POST | /auth/logout | 204 | 401 | Bearer |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU01-MOB-01](./UT-CU01-MOB-01.md) | Proyecto Flutter y núcleo (`core/`) | — |
| [UT-CU01-MOB-02](./UT-CU01-MOB-02.md) | Modelos y `AuthService` | 01 |
| [UT-CU01-MOB-03](./UT-CU01-MOB-03.md) | Sesión guardada y `sesionProvider` | 02 |
| [UT-CU01-MOB-04](./UT-CU01-MOB-04.md) | Router con redirect global, bienvenida y `main.dart` | 03 |
| [UT-CU01-MOB-05](./UT-CU01-MOB-05.md) | Pantalla de inicio de sesión | 03, 04 |
| [UT-CU01-MOB-06](./UT-CU01-MOB-06.md) | Inicio con accesos por permiso y cerrar sesión | 03, 04 |

**Cadena:** 01 → 02 → 03 → 04 → 05, 06

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA01 | Login con permisos y accesos | 02, 03, 05, 06 |
| CA02, CA03, CA04 | Mensajes de error del backend | 01, 05 |
| CA05 | Cerrar sesión revoca | 02, 06 |
| CA06 | Con sesión, `/login` lleva al inicio | 04 |

## Follow-ups (no se implementan en estas UT)

- Biometría para reabrir la sesión.
- Notificaciones push del cambio de estado (CU-14).

## Qué NO incluyen estas UT

- Recuperar contraseña (CU-02): el login enlaza a la web.
- Pantallas de cada perfil: las crea su propio CU.
