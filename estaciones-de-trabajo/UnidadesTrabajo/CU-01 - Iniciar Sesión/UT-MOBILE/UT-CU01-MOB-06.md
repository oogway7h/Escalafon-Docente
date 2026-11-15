# UT-CU01-MOB-06 — Inicio con accesos por permiso y cerrar sesión

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-MOB-06` |
| **Nombre** | Inicio con accesos por permiso y cerrar sesión |
| **Caso de uso** | CU-01-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU01-MOB-03, UT-CU01-MOB-04 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [WIDGETS_UI_MOBILE.md](../../../.claude/rules/mobile/WIDGETS_UI_MOBILE.md)

## Solicitud de negocio

Cada perfil ve solo lo que puede hacer, igual que en la web.

## Objetivo

`/inicio` con saludo, tarjetas de acceso calculadas por permisos y botón de salir.

## Alcance

- `features/inicio/accesos.dart`: `List<Acceso> accesosPara(List<String> permisos)` (titulo, descripcion, icono, ruta, permiso). Empieza **vacía**: cada CU-MOB agrega su acceso.
- `inicio_screen.dart`: saludo con el nombre, cuadrícula de accesos, estado vacío «Tus herramientas están en la web» si no hay ninguno, menú con «Cerrar sesión».
- Test de `accesosPara` (lista vacía hoy; cada CU agrega su caso).

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/inicio/accesos.dart`
- `ded_Mobile/lib/features/inicio/screens/inicio_screen.dart`
- `ded_Mobile/test/features/inicio/accesos_test.dart`

## Fuera de alcance (no tocar)

- Los accesos de cada perfil (los agrega su CU).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/inicio` con saludo, tarjetas de acceso calculadas por permisos y botón de salir.

## Validación

1. `flutter test test/features/inicio` en verde.
2. En el emulador Android con la API local (`make dev`, `make seed`): cerrar sesión vuelve a la bienvenida y reabrir la app no entra.
3. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
