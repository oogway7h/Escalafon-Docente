# UT-CU25-MOB-01 — Modelo y `TableroService`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU25-MOB-01` |
| **Nombre** | Modelo y `TableroService` |
| **Caso de uso** | CU-25-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Indicadores tipados.

## Objetivo

`Tablero` con resumen, por estado, analistas y alertas; `obtener({desde, hasta})`.

## Alcance

- Tests de parseo y query.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tablero/models/tablero.dart`
- `ded_Mobile/lib/features/tablero/services/tablero_service.dart`
- `ded_Mobile/test/features/tablero/tablero_service_test.dart`

## Fuera de alcance (no tocar)

- Pantalla.
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Tablero` con resumen, por estado, analistas y alertas; `obtener({desde, hasta})`.

## Validación

1. `flutter test test/features/tablero` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
