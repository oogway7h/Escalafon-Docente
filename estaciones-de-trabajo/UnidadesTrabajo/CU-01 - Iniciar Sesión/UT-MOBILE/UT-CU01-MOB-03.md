# UT-CU01-MOB-03 — Sesión guardada y `sesionProvider`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-MOB-03` |
| **Nombre** | Sesión guardada y `sesionProvider` |
| **Caso de uso** | CU-01-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU01-MOB-02 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [ESTADO_GLOBAL_MOBILE.md](../../../.claude/rules/mobile/ESTADO_GLOBAL_MOBILE.md)

## Solicitud de negocio

El personal no debe escribir la contraseña cada vez que abre la app, y una sesión revocada no debe seguir valiendo.

## Objetivo

`sesionProvider` restaura, inicia y cierra la sesión.

## Alcance

- `SesionStorage` sobre `flutter_secure_storage` (guardar, leer, borrar el JSON de la sesión).
- `SesionNotifier extends AsyncNotifier<Sesion?>`: `build()` lee el almacenamiento, descarta si `expiraEn` pasó y confirma con `sesion()` (401 → null); `iniciar(correo, clave)`; `cerrar()` (logout best-effort + borrar).
- `dioProvider` y `authServiceProvider`: el Dio lee el token de la sesión y en 401 llama a `cerrarLocal()`.
- Tests con un almacenamiento en memoria y el Dio falso: restaurar válida, vencida, revocada; iniciar guarda; cerrar borra.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/auth/services/sesion_storage.dart`
- `ded_Mobile/lib/features/auth/providers/sesion_provider.dart`
- `ded_Mobile/test/features/auth/sesion_provider_test.dart`

## Fuera de alcance (no tocar)

- Router (UT 04).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`sesionProvider` restaura, inicia y cierra la sesión.

## Validación

1. `flutter test test/features/auth` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
