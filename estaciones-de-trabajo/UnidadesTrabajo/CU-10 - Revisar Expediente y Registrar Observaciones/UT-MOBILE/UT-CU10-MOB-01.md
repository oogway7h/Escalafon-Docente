# UT-CU10-MOB-01 — Modelos y `TramitesService` (lectura)

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU10-MOB-01` |
| **Nombre** | Modelos y `TramitesService` (lectura) |
| **Caso de uso** | CU-10-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [CONSUMO_API_MOBILE.md](../../../.claude/rules/mobile/CONSUMO_API_MOBILE.md) · [ESTADO_GLOBAL_MOBILE.md](../../../.claude/rules/mobile/ESTADO_GLOBAL_MOBILE.md)

## Solicitud de negocio

El analista necesita su bandeja tipada.

## Objetivo

`ResumenTramite`, `TramiteDetalle` y `listar` / `obtener`.

## Alcance

- `TramiteDetalle` con `historial`, `observaciones`, `documentos` (`estadoOcr`, `activo`) y `requisitos`.
- `TramitesService(Dio)`: `listar({bool mios = false, bool sinAsignar = false})` (arreglo completo) y `obtener(id)`.
- Providers `misTramitesListProvider` y `tramiteDetailProvider` (family).
- Tests de parseo del detalle y de la query `mios=true`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tramites/models/tramite.dart`
- `ded_Mobile/lib/features/tramites/services/tramites_service.dart`
- `ded_Mobile/lib/features/tramites/providers/tramites_provider.dart`
- `ded_Mobile/test/features/tramites/tramites_service_test.dart`

## Fuera de alcance (no tocar)

- Asignación (CU-09).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`ResumenTramite`, `TramiteDetalle` y `listar` / `obtener`.

## Validación

1. `flutter test test/features/tramites` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
