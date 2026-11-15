# UT-CU09-MOB-02 — Pantalla «Por asignar» y reasignación en el detalle

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU09-MOB-02` |
| **Nombre** | Pantalla «Por asignar» y reasignación en el detalle |
| **Caso de uso** | CU-09-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU09-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [FORMULARIOS_MOBILE.md](../../../.claude/rules/mobile/FORMULARIOS_MOBILE.md)

## Solicitud de negocio

Repartir viendo la carga.

## Objetivo

`/inicio/asignar` con hoja de analistas ordenados por carga, y «Reasignar» en el detalle.

## Alcance

- Hoja inferior: analistas con su carga; elegir asigna.
- Reasignar (trámite con analista) exige motivo (validator, igual que la web).
- Errores del backend con `mensajeDeError`.
- `accesosPara` agrega «Por asignar» con `TRAMITE_ASIGNAR`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tramites/screens/asignar_screen.dart`
- `ded_Mobile/lib/features/tramites/widgets/elegir_analista_sheet.dart`
- `ded_Mobile/lib/features/tramites/screens/tramite_detalle_screen.dart`
- `ded_Mobile/lib/features/inicio/accesos.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/inicio/asignar` con hoja de analistas ordenados por carga, y «Reasignar» en el detalle.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): `admin@ded.dev` asigna un trámite ingresado y desaparece de «Por asignar».
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
