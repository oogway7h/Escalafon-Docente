# UT-CU12-MOB-01 — Modelos y `PortalService` (consulta)

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU12-MOB-01` |
| **Nombre** | Modelos y `PortalService` (consulta) |
| **Caso de uso** | CU-12-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [CONSUMO_API_MOBILE.md](../../../.claude/rules/mobile/CONSUMO_API_MOBILE.md)

## Solicitud de negocio

El docente necesita su estado en objetos claros, sin datos personales.

## Objetivo

DTOs públicos y `consultar` / `consultarPorDocente`.

## Alcance

- `VistaPublica`, `EtapaPublica`, `ObservacionPublica`, `ResumenPublico` con `fromJson`.
- `PortalService(Dio)`: `consultar(codigo)` y `consultarPorDocente(codigo, cedula)`.
- `lib/shared/etiquetas.dart`: `etiquetaEstado`, `descripcionEstado`, `fasesTramite`.
- Tests: parseo completo, 404 como `ApiException`, cédula en el cuerpo y no en la URL.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/portal/models/portal.dart`
- `ded_Mobile/lib/features/portal/services/portal_service.dart`
- `ded_Mobile/lib/shared/etiquetas.dart`
- `ded_Mobile/test/features/portal/portal_service_test.dart`

## Fuera de alcance (no tocar)

- Solicitud (CU-13).
- Pantalla (UT 02).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

DTOs públicos y `consultar` / `consultarPorDocente`.

## Validación

1. `flutter test test/features/portal` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
