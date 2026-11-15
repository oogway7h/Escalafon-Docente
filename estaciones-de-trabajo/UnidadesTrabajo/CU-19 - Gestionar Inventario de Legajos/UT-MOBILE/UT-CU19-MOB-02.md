# UT-CU19-MOB-02 — Pantalla de legajos

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU19-MOB-02` |
| **Nombre** | Pantalla de legajos |
| **Caso de uso** | CU-19-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU19-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Encontrar la carpeta y su ubicación al instante.

## Objetivo

`/inicio/legajos` con buscador, filtro En archivo/Prestado y lista con gaveta y orden.

## Alcance

- `accesosPara` agrega «Legajos» con `LEGAJO_GESTIONAR`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/legajos/screens/legajos_screen.dart`
- `ded_Mobile/lib/features/inicio/accesos.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/inicio/legajos` con buscador, filtro En archivo/Prestado y lista con gaveta y orden.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): buscar un `LEG-…` muestra solo esa carpeta con su gaveta.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
