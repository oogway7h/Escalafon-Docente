---
globs: ../ded_Mobile/**/*
---

# Estado global — Mobile (Riverpod)

Todo es Riverpod (`flutter_riverpod`, sin generación de código). Se distingue por **nombre y ubicación**
si un provider envuelve datos del servidor o estado de pantalla.

## Datos de la API

En `features/{feature}/providers/`, junto al service:

- Listado: `FutureProvider` (o `.autoDispose`) → `AsyncValue<List<T>>`.
- Detalle: `FutureProvider.family` por id.
- Mutaciones: método en el service llamado desde la pantalla; al terminar, `ref.invalidate(...)` del
  listado afectado (equivalente a `revalidatePath` de la web).

## Sesión

`sesionProvider` (`features/auth/providers/`) es el **único** que sabe si hay sesión. Lo leen el
router (redirect) y el `ApiClient` (token). Se guarda en `flutter_secure_storage`.

## Nombres

| Tipo | Convención | Ejemplo |
|---|---|---|
| Listado | `{recurso}ListProvider` | `tramitesListProvider` |
| Detalle | `{recurso}DetailProvider` (family) | `tramiteDetailProvider` |
| Service | `{feature}ServiceProvider` | `legajosServiceProvider` |
| Estado de UI | `{algo}UiProvider` | — (se usa `StatefulWidget` mientras sea local a una pantalla) |

Nunca un nombre genérico (`provider`, `state`).
