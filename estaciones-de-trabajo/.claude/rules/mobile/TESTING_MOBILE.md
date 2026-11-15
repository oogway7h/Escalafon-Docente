---
globs: ../ded_Mobile/**/*
---

# Testing — Mobile

Caso a caso, no cobertura exhaustiva (igual criterio que DriveHub):

- **Sí:** parseo de DTOs no triviales, services (cuerpo y ruta exactos con un adaptador de Dio falso),
  conversión de errores a `ApiException`, lógica de providers (sesión, menú por permisos), funciones
  puras con reglas de negocio.
- **No:** widgets triviales o pantallas que solo arman un formulario sobre un service ya probado.

El adaptador falso vive en `test/helpers/dio_falso.dart` (`HttpClientAdapter` que registra la petición
y responde lo programado). `flutter test` nunca abre la cámara ni habla con la red.

Un test que existe se mantiene en verde: un test roto que nadie arregla desinforma.
