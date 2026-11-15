# UT-CU01-MOB-05 — Pantalla de inicio de sesión

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-MOB-05` |
| **Nombre** | Pantalla de inicio de sesión |
| **Caso de uso** | CU-01-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU01-MOB-03, UT-CU01-MOB-04 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [FORMULARIOS_MOBILE.md](../../../.claude/rules/mobile/FORMULARIOS_MOBILE.md)

## Solicitud de negocio

El personal entra con su correo institucional y entiende por qué falló.

## Objetivo

`/login` con validación y mensajes del backend.

## Alcance

- `Form` con correo (obligatorio, con @) y contraseña (obligatoria, ocultable).
- Llama a `sesionProvider.iniciar`; muestra `mensajeDeError(e)` (credenciales, cuenta bloqueada, sin conexión).
- No navega a mano: el redirect lleva a `/inicio` al cambiar la sesión.
- Enlace «¿Olvidaste tu contraseña?» explica que se recupera desde la web.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/auth/screens/login_screen.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`

## Fuera de alcance (no tocar)

- Recuperación (CU-02).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/login` con validación y mensajes del backend.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): `archivo@ded.dev`/`123123` entra; una contraseña errónea muestra «credenciales inválidas».
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
