# UT-CU12-MOB-02 — Pantalla de seguimiento

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU12-MOB-02` |
| **Nombre** | Pantalla de seguimiento |
| **Caso de uso** | CU-12-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU12-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [WIDGETS_UI_MOBILE.md](../../../.claude/rules/mobile/WIDGETS_UI_MOBILE.md)

## Solicitud de negocio

Consultar desde el celular sin ir a ventanilla.

## Objetivo

`/portal/seguimiento` con buscador y detalle del trámite.

## Alcance

- Un campo de código; con ≤ 5 caracteres aparece el de cédula.
- Detalle: código, estado con su descripción, línea de fases (En facultad = «sin seguimiento»), observaciones y etapas, resultado si existe.
- Con varios trámites del docente, lista para elegir.
- Errores genéricos (no revela qué falló).
- Ruta en el router y botón «Soy docente» de la bienvenida.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/portal/screens/seguimiento_screen.dart`
- `ded_Mobile/lib/shared/widgets/linea_fases.dart (creado en portal/widgets y movido a shared en CU-10)`
- `ded_Mobile/lib/core/routing/app_router.dart`
- `ded_Mobile/lib/features/inicio/screens/bienvenida_screen.dart`

## Fuera de alcance (no tocar)

- Solicitud (CU-13).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/portal/seguimiento` con buscador y detalle del trámite.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): un código ED-… real muestra su fase; uno falso, «No encontramos un trámite con ese código».
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
