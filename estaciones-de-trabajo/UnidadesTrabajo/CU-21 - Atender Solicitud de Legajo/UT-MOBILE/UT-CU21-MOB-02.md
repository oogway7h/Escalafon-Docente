# UT-CU21-MOB-02 — Pantalla Custodia: pestaña Por atender

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU21-MOB-02` |
| **Nombre** | Pantalla Custodia: pestaña Por atender |
| **Caso de uso** | CU-21-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU21-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Archivo atiende desde el depósito, con el teléfono.

## Objetivo

`/inicio/custodia` con la pestaña Por atender (aprobar / rechazar con motivo).

## Alcance

- Nota obligatoria al rechazar.
- `accesosPara` agrega «Custodia» con `LEGAJO_GESTIONAR` o `CUSTODIA_REGISTRAR`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/carpetas/screens/custodia_screen.dart`
- `ded_Mobile/lib/features/inicio/accesos.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`

## Fuera de alcance (no tocar)

- Entrega y devolución.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/inicio/custodia` con la pestaña Por atender (aprobar / rechazar con motivo).

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): `archivo@ded.dev` rechaza una solicitud con motivo y el analista la ve Rechazada.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
