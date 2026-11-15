# UT-CU20-MOB-01 — Modelo `SolicitudCarpeta` y `CarpetasService` (solicitar y listar)

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU20-MOB-01` |
| **Nombre** | Modelo `SolicitudCarpeta` y `CarpetasService` (solicitar y listar) |
| **Caso de uso** | CU-20-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Pedir la carpeta sin ir al archivo.

## Objetivo

`solicitar(legajoId, tramiteId, motivo)` y `listar({estado, pagina})`.

## Alcance

- Etiquetas de solicitud y estado material en `lib/shared/etiquetas.dart`.
- Tests de cuerpo y parseo.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/carpetas/models/solicitud_carpeta.dart`
- `ded_Mobile/lib/features/carpetas/services/carpetas_service.dart`
- `ded_Mobile/lib/shared/etiquetas.dart`
- `ded_Mobile/test/features/carpetas/carpetas_service_test.dart`

## Fuera de alcance (no tocar)

- Atender, despachar, devolver.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`solicitar(legajoId, tramiteId, motivo)` y `listar({estado, pagina})`.

## Validación

1. `flutter test test/features/carpetas` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
