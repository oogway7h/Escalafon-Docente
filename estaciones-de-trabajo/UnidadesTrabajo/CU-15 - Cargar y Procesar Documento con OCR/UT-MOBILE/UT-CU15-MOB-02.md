# UT-CU15-MOB-02 — Subir documento desde el detalle

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-MOB-02` |
| **Nombre** | Subir documento desde el detalle |
| **Caso de uso** | CU-15-MOB |
| **Proyecto** | `ded_Mobile` |
| **Depende de** | UT-CU15-MOB-01 |
| **Estado** | Hecho |

## Antes de implementar

[ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md)

## Solicitud de negocio

Fotografiar el papel que trae el docente y subirlo en el momento.

## Objetivo

Botón «Subir documento» en el detalle (si el usuario tiene `DOCUMENTO_GESTIONAR` y el trámite no está archivado) con hoja de origen y requisito.

## Alcance

- Cámara / Galería / PDF; requisito opcional.
- Valida tipo y tamaño; muestra progreso y luego el estado OCR del documento en la lista.
- Recarga `tramiteDetailProvider` al terminar.

Archivos esperados (crear/tocar):

- `ded_Mobile/lib/features/tramites/widgets/subir_documento_sheet.dart`
- `ded_Mobile/lib/features/tramites/screens/tramite_detalle_screen.dart`

## Fuera de alcance (no tocar)

- Cualquier archivo de `ded_Mobile` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Botón «Subir documento» en el detalle (si el usuario tiene `DOCUMENTO_GESTIONAR` y el trámite no está archivado) con hoja de origen y requisito.

## Validación

1. En el emulador Android con la API local (`make dev`, `make seed`): subir una foto a un trámite asignado la muestra como «Pendiente/Procesando» y luego «Procesado».
2. `flutter analyze` sin issues y `flutter test` en verde en `ded_Mobile/`.
