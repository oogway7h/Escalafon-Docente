# UT-CU09-MOB-01 — Asignación en `TramitesService`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU09-MOB-01` |
| **Nombre** | Asignación en `TramitesService` |
| **Caso de uso** | CU-09-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

El Director reparte desde el teléfono.

## Objetivo

`analistas()` y `asignar(id, analistaId, {motivo})`.

## Alcance

- `Analista(id, nombre, carga)`.
- Provider `porAsignarListProvider` (`sinAsignar=true`).
- Tests: cuerpo con y sin motivo.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tramites/models/tramite.dart`
- `ded_Mobile/lib/features/tramites/services/tramites_service.dart`
- `ded_Mobile/lib/features/tramites/providers/tramites_provider.dart`
- `ded_Mobile/test/features/tramites/tramites_service_test.dart`

## Fuera de alcance (no tocar)

- Pantallas (UT 02).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`analistas()` y `asignar(id, analistaId, {motivo})`.

## Validación

1. `flutter test test/features/tramites` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
