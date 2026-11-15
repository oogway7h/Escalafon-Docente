# UT-CU22-MOB-01 — `despachar` y `entregar` en `CarpetasService`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU22-MOB-01` |
| **Nombre** | `despachar` y `entregar` en `CarpetasService` |
| **Caso de uso** | CU-22-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Registrar la salida.

## Objetivo

Dos métodos con el mismo cuerpo.

## Alcance

- Test de rutas y cuerpo.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/carpetas/services/carpetas_service.dart`
- `ded_Mobile/test/features/carpetas/carpetas_service_test.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Dos métodos con el mismo cuerpo.

## Validación

1. `flutter test test/features/carpetas` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
