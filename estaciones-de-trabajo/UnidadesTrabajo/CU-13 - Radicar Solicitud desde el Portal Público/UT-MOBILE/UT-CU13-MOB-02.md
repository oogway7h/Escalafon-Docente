# UT-CU13-MOB-02 — Pantalla de nueva solicitud

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU13-MOB-02` |
| **Nombre** | Pantalla de nueva solicitud |
| **Caso de uso** | CU-13-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU13-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [FORMULARIOS_MOBILE.md](../../../.claude/rules/mobile/FORMULARIOS_MOBILE.md)

## Solicitud de negocio

Radicar desde la casa con el celular.

## Objetivo

`/portal/solicitud` con identificación, tipo, categoría y un adjunto por requisito.

## Alcance

- Dependencias `image_picker` y `file_picker`; permisos de cámara (Android) y `NSCameraUsageDescription`/`NSPhotoLibraryUsageDescription` (iOS).
- Por requisito: fila con estado y botones Cámara / Galería / PDF.
- Valida faltantes y 8 MB antes de enviar; muestra errores del backend.
- Al terminar, muestra el código y lleva al seguimiento con ese código.

Archivos esperados (crear/tocar):

- `ded_Mobile/pubspec.yaml`
- `ded_Mobile/lib/features/portal/screens/solicitud_screen.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`
- `ded_Mobile/lib/features/inicio/screens/bienvenida_screen.dart`
- `ded_Mobile/ios/Runner/Info.plist`

## Fuera de alcance (no tocar)

- Seguimiento.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/portal/solicitud` con identificación, tipo, categoría y un adjunto por requisito.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): adjuntar una foto por requisito y enviar devuelve un código ED-… que se ve en el seguimiento.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
