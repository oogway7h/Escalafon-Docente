# UT-CU19-MOB-01 — Modelo `Legajo` y `LegajosService`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU19-MOB-01` |
| **Nombre** | Modelo `Legajo` y `LegajosService` |
| **Caso de uso** | CU-19-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Archivo y analistas necesitan ubicar carpetas.

## Objetivo

`listar({buscar, estado, pagina})` paginado.

## Alcance

- `Legajo` y `Pagina<T>` (datos, total, pagina, tamano) en `lib/shared/pagina.dart` (la usan legajos y carpetas).
- Tests de query y parseo.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/legajos/models/legajo.dart`
- `ded_Mobile/lib/features/legajos/services/legajos_service.dart`
- `ded_Mobile/lib/shared/pagina.dart`
- `ded_Mobile/test/features/legajos/legajos_service_test.dart`

## Fuera de alcance (no tocar)

- Pantalla.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`listar({buscar, estado, pagina})` paginado.

## Validación

1. `flutter test test/features/legajos` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
