---
globs: ../ded_Mobile/**/*
---

# Navegación — Mobile (go_router)

Un único `GoRouter` en `lib/core/routing/app_router.dart`, expuesto como `routerProvider`. El chequeo
de sesión vive **solo** en el `redirect` raíz, nunca pantalla por pantalla:

| Estado de sesión | Ruta pedida | Resultado |
|---|---|---|
| cargando (leyendo el almacenamiento) | cualquiera | `/cargando` |
| sin sesión | `/inicio/...` | `/` (bienvenida) |
| con sesión | `/`, `/login` | `/inicio` |
| cualquiera | `/portal/...` | sin redirect (público) |

Rutas: `/` bienvenida (Soy docente / Soy personal) · `/login` · `/portal/seguimiento` ·
`/portal/solicitud` · `/inicio` (accesos por permiso) · `/inicio/tramites`, `/inicio/tramites/:id` ·
`/inicio/carpetas`, `/inicio/carpetas/nueva?tramiteId=&tramiteCodigo=&docenteCodigo=` · `/inicio/legajos` ·
`/inicio/custodia` · `/inicio/asignar` · `/inicio/tablero`.

- El router depende de `sesionProvider` (con `refreshListenable`), no lee `flutter_secure_storage`.
- El árbol completo se arma en `core/routing/`; las features solo exponen sus pantallas.
- Una ruta nueva entra con la UT del CU que la necesita, no antes.
