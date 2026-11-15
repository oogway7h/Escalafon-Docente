# UT-CU20-MOB-02 — Mis solicitudes y solicitar desde el trámite

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU20-MOB-02` |
| **Nombre** | Mis solicitudes y solicitar desde el trámite |
| **Caso de uso** | CU-20-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU20-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [FORMULARIOS_MOBILE.md](../../../.claude/rules/mobile/FORMULARIOS_MOBILE.md)

## Solicitud de negocio

Pedir y seguir carpetas desde el celular.

## Objetivo

`/inicio/carpetas` (mis solicitudes) y `/inicio/carpetas/nueva` (carpetas del docente + motivo).

## Alcance

- El detalle del trámite muestra «Solicitar carpeta» si el usuario tiene `CARPETA_SOLICITAR` (navega por ruta).
- Motivo obligatorio.
- `accesosPara` agrega «Mis carpetas» con `CARPETA_SOLICITAR`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/carpetas/screens/mis_solicitudes_screen.dart`
- `ded_Mobile/lib/features/carpetas/screens/solicitar_carpeta_screen.dart`
- `ded_Mobile/lib/features/tramites/screens/tramite_detalle_screen.dart`
- `ded_Mobile/lib/features/inicio/accesos.dart`
- `ded_Mobile/lib/core/routing/app_router.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/inicio/carpetas` (mis solicitudes) y `/inicio/carpetas/nueva` (carpetas del docente + motivo).

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): solicitar una carpeta desde un trámite la muestra Pendiente en «Mis carpetas».
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
