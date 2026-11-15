# UT-CU23-MOB-02 — Pestaña Prestadas y hoja de devolución

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU23-MOB-02` |
| **Nombre** | Pestaña Prestadas y hoja de devolución |
| **Caso de uso** | CU-23-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU23-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [FORMULARIOS_MOBILE.md](../../../.claude/rules/mobile/FORMULARIOS_MOBILE.md)

## Solicitud de negocio

Recibir la carpeta.

## Objetivo

Pestaña Prestadas (vencidas resaltadas) y hoja con estado material y observación.

## Alcance

- Observación obligatoria si no es Buen estado (validator).

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/carpetas/widgets/devolucion_sheet.dart`
- `ded_Mobile/lib/features/carpetas/screens/custodia_screen.dart`
- `ded_Mobile/test/features/carpetas/devolucion_test.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Pestaña Prestadas (vencidas resaltadas) y hoja con estado material y observación.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): devolver una carpeta deteriorada sin observación muestra el error.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
