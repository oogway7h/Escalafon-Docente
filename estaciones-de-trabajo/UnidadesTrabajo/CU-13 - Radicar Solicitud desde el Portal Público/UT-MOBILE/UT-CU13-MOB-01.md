# UT-CU13-MOB-01 — Categorías públicas y radicación en `PortalService`

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU13-MOB-01` |
| **Nombre** | Categorías públicas y radicación en `PortalService` |
| **Caso de uso** | CU-13-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | — |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) · [CONSUMO_API_MOBILE.md](../../../.claude/rules/mobile/CONSUMO_API_MOBILE.md)

## Solicitud de negocio

La app debe pedir exactamente los documentos que define el catálogo.

## Objetivo

`PortalService` gana `categorias()` y `radicar(...)` multipart.

## Alcance

- `CategoriaPublica` y `RequisitoPublico`.
- `radicar({identificador, tipo, categoriaId, Map<int, ArchivoAdjunto>})` → `FormData` con `doc_<id>`; devuelve el código.
- `ArchivoAdjunto(nombre, bytes, mime)` y `elegirArchivo`/`elegirOrigen` (cámara, galería, PDF) en `lib/shared/archivo_adjunto.dart`: los usa también CU-15.
- Tests: nombres de campo exactos del multipart y parseo del código.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/portal/models/portal.dart`
- `ded_Mobile/lib/features/portal/services/portal_service.dart`
- `ded_Mobile/lib/shared/archivo_adjunto.dart`
- `ded_Mobile/test/features/portal/portal_service_test.dart`

## Fuera de alcance (no tocar)

- Pantalla (UT 02).
- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`PortalService` gana `categorias()` y `radicar(...)` multipart.

## Validación

1. `flutter test test/features/portal` en verde.
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
