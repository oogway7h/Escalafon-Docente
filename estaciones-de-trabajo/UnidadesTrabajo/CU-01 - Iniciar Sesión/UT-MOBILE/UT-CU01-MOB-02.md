# UT-CU01-MOB-02 — Modelos y `AuthService`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-MOB-02` |
| **Nombre** | Modelos y `AuthService` |
| **Caso de uso** | CU-01-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU01-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [CONSUMO_API_MOBILE.md](../../../.claude/rules/mobile/CONSUMO_API_MOBILE.md)

## Solicitud de negocio

La app necesita enviar las credenciales y recibir la sesión tipada.

## Objetivo

DTOs de sesión y service con `login`, `sesion` y `logout`.

## Alcance

- `Sesion` (`token`, `expiraEn` como `DateTime`, `usuario`, `permisos`) y `UsuarioSesion` (`id`, `nombres`, `apellidos`, `correo`); `toJson/fromJson` para guardarla.
- `AuthService(Dio)`: `login(correo, contrasena)` → `POST /auth/login`; `sesion()` → `GET /auth/sesion` (devuelve permisos actualizados); `logout()`.
- Tests con el Dio falso: cuerpo exacto del login, parseo, 401 `CREDENCIALES_INVALIDAS` como `ApiException`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/auth/models/sesion.dart`
- `ded_Mobile/lib/features/auth/services/auth_service.dart`
- `ded_Mobile/test/features/auth/auth_service_test.dart`

## Fuera de alcance (no tocar)

- Guardar la sesión (UT 03).
- Pantallas.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

DTOs de sesión y service con `login`, `sesion` y `logout`.

## Validación

1. `flutter test test/features/auth` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
