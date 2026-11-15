# UT-CU22-MOB-02 — Hoja de entrega y pestaña Aprobadas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU22-MOB-02` |
| **Nombre** | Hoja de entrega y pestaña Aprobadas |
| **Caso de uso** | CU-22-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU22-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [FORMULARIOS_MOBILE.md](../../../.claude/rules/mobile/FORMULARIOS_MOBILE.md)

## Solicitud de negocio

Entregar la carpeta en el mostrador.

## Objetivo

Hoja con fecha límite (atajos +3/+7/+15, selector de fecha), estado material y observación.

## Alcance

- Pestaña Aprobadas en Custodia.
- Botón «Aprobar y entregar» en Por atender.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/carpetas/widgets/entrega_sheet.dart`
- `ded_Mobile/lib/features/carpetas/screens/custodia_screen.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Hoja con fecha límite (atajos +3/+7/+15, selector de fecha), estado material y observación.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): entregar una carpeta la pasa a Prestadas y el legajo a Prestado.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
