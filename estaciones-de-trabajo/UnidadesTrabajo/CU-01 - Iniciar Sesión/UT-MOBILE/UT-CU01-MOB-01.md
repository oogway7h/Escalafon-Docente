# UT-CU01-MOB-01 — Proyecto Flutter y núcleo (`core/`)

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-MOB-01` |
| **Nombre** | Proyecto Flutter y núcleo (`core/`) |
| **Caso de uso** | CU-01-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [CONSUMO_API_MOBILE.md](../../../.claude/rules/mobile/CONSUMO_API_MOBILE.md) · [WIDGETS_UI_MOBILE.md](../../../.claude/rules/mobile/WIDGETS_UI_MOBILE.md) · [TESTING_MOBILE.md](../../../.claude/rules/mobile/TESTING_MOBILE.md)

## Solicitud de negocio

Toda pantalla necesita hablar con la API con el mismo manejo de errores y verse con la identidad del DED.

## Objetivo

Proyecto creado para Android/iOS con `ApiClient`, `ApiException`, configuración y tema.

## Alcance

- `flutter create --org bo.edu.uagrm --platforms android,ios --project-name ded_mobile .`
- Dependencias: `flutter_riverpod`, `go_router`, `dio`, `flutter_secure_storage` (más adelante CU-13/15 agregan las de archivos).
- `core/config.dart`: `const apiUrl = String.fromEnvironment('API_URL', defaultValue: 'http://10.0.2.2:8080/api')`.
- `core/network/api_client.dart`: `crearDio({token, onNoAutorizado})` con interceptor Bearer; `ApiException(status, codigo, mensaje)`; `mensajeDeError(Object)`.
- `core/theme/app_theme.dart`: claro y oscuro con la paleta de `globals.css`.
- Android: permiso `INTERNET` y `usesCleartextTraffic` (API local por http).
- `test/helpers/dio_falso.dart` y `test/core/api_client_test.dart` (error JSON → ApiException; sin respuesta → `SIN_CONEXION`; 401 llama a `onNoAutorizado`).

Archivos esperados (crear/tocar):

- `ded_Mobile/pubspec.yaml`
- `ded_Mobile/lib/core/config.dart`
- `ded_Mobile/lib/core/network/api_client.dart`
- `ded_Mobile/lib/core/theme/app_theme.dart`
- `ded_Mobile/android/app/src/main/AndroidManifest.xml`
- `ded_Mobile/test/helpers/dio_falso.dart`
- `ded_Mobile/test/core/api_client_test.dart`

## Fuera de alcance (no tocar)

- Router y pantallas (UT 04-06).
- Cualquier service de feature.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Proyecto creado para Android/iOS con `ApiClient`, `ApiException`, configuración y tema.

## Validación

1. `flutter test test/core` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
