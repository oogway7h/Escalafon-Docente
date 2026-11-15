# UT-CU25-MOB-02 — Pantalla del tablero

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU25-MOB-02` |
| **Nombre** | Pantalla del tablero |
| **Caso de uso** | CU-25-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU25-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [WIDGETS_UI_MOBILE.md](../../../.claude/rules/mobile/WIDGETS_UI_MOBILE.md)

## Solicitud de negocio

Decidir con datos desde cualquier lugar.

## Objetivo

`/inicio/tablero` con 4 indicadores, barras por estado y por analista y alertas de mora.

## Alcance

- Aviso arriba si hay carpetas en mora.
- `accesosPara` agrega «Tablero» con `TABLERO_VER`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tablero/screens/tablero_screen.dart`
- `ded_Mobile/lib/features/inicio/accesos.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/inicio/tablero` con 4 indicadores, barras por estado y por analista y alertas de mora.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): `admin@ded.dev` ve los indicadores; cambiar a «Este año» recalcula.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
