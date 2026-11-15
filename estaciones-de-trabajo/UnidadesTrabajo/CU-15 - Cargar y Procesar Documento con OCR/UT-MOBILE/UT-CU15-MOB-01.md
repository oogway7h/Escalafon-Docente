# UT-CU15-MOB-01 — `DocumentosService`: preparar, subir y confirmar

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-MOB-01` |
| **Nombre** | `DocumentosService`: preparar, subir y confirmar |
| **Caso de uso** | CU-15-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [CONSUMO_API_MOBILE.md](../../../.claude/rules/mobile/CONSUMO_API_MOBILE.md)

## Solicitud de negocio

El archivo debe ir directo al almacenamiento como en la web.

## Objetivo

Servicio con el flujo de tres pasos.

## Alcance

- `subir(tramiteId, ArchivoAdjunto, {requisitoId})` → `POST /subida`, `PUT` con un Dio limpio, `POST /confirmar`; devuelve el `Documento`.
- Tests: cuerpo de `/subida`, el PUT va a la URL firmada sin `Authorization` y con `Content-Type`, luego `/confirmar`.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tramites/services/documentos_service.dart`
- `ded_Mobile/test/features/tramites/documentos_service_test.dart`

## Fuera de alcance (no tocar)

- Pantalla (UT 02).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Servicio con el flujo de tres pasos.

## Validación

1. `flutter test test/features/tramites` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
