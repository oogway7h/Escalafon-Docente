# UT-CU10-MOB-02 — Bandeja y detalle del trámite

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU10-MOB-02` |
| **Nombre** | Bandeja y detalle del trámite |
| **Caso de uso** | CU-10-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU10-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Ver lo asignado desde cualquier lugar.

## Objetivo

`/inicio/tramites` y `/inicio/tramites/:id` con acceso «Mi bandeja» para `TRAMITE_REVISAR`.

## Alcance

- Bandeja con estado y docente; tirar para recargar; estado vacío.
- Detalle: datos, fases, requisitos (presente o faltante), observaciones e historial.
- `accesosPara` agrega «Mi bandeja» con `TRAMITE_REVISAR` (+ caso en el test).

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tramites/screens/bandeja_screen.dart`
- `ded_Mobile/lib/features/tramites/screens/tramite_detalle_screen.dart`
- `ded_Mobile/lib/shared/widgets/linea_fases.dart (movido desde portal)`
- `ded_Mobile/lib/shared/widgets/vista_async.dart`
- `ded_Mobile/lib/features/inicio/accesos.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`

## Fuera de alcance (no tocar)

- Acciones de revisión.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/inicio/tramites` y `/inicio/tramites/:id` con acceso «Mi bandeja» para `TRAMITE_REVISAR`.

## Validación

1. `flutter test test/features/inicio` en verde.
2. En el emulador Android con la API local (`make dev`, `make seed`): `analista@ded.dev` ve «Mi bandeja» con sus trámites.
3. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
