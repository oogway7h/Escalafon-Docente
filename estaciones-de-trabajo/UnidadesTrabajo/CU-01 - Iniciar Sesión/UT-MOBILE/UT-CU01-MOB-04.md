# UT-CU01-MOB-04 — Router con redirect global, bienvenida y `main.dart`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-MOB-04` |
| **Nombre** | Router con redirect global, bienvenida y `main.dart` |
| **Caso de uso** | CU-01-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU01-MOB-03 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [NAVEGACION_MOBILE.md](../../../.claude/rules/mobile/NAVEGACION_MOBILE.md)

## Solicitud de negocio

La decisión de a dónde puede entrar el usuario debe tomarse en un solo lugar.

## Objetivo

`GoRouter` con `/cargando`, `/`, `/login`, `/inicio` y redirect por sesión.

## Alcance

- `core/routing/app_router.dart`: `routerProvider` con `refreshListenable` atado a `sesionProvider`; reglas de la tabla de NAVEGACION_MOBILE.md; `/portal/*` sin redirect.
- `features/inicio/screens/bienvenida_screen.dart`: logo, «Soy docente» (portal, se habilita en CU-12) y «Soy personal del DED» (login).
- `main.dart`: `ProviderScope` + `MaterialApp.router` con el tema claro/oscuro.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/core/routing/app_router.dart`
- `ded_Mobile/lib/features/inicio/screens/bienvenida_screen.dart`
- `ded_Mobile/lib/main.dart`

## Fuera de alcance (no tocar)

- Pantalla de login (UT 05).
- Rutas de los perfiles.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`GoRouter` con `/cargando`, `/`, `/login`, `/inicio` y redirect por sesión.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): sin sesión la app abre la bienvenida.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
