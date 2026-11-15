---
globs: ../ded_Mobile/**/*
---

# Arquitectura Mobile (Flutter)

Repo `ded_Mobile/` (https://github.com/oogway7h/DED-Mobile). Flutter para **Android e iOS**, Riverpod,
go_router, dio. Adaptada de la arquitectura mobile de DriveHub: **feature-based**, lo transversal en
`core/`.

## Para quién es la app

| Perfil | Entra | Qué hace |
|---|---|---|
| Docente | sin cuenta | consulta su trámite (CU-12) y radica una solicitud con fotos o PDF (CU-13) |
| Analista | login | bandeja de trámites (CU-10, lectura), sube documentos desde la cámara (CU-15), pide carpetas (CU-20) |
| Personal de Archivo | login | busca legajos (CU-19), atiende, despacha y recibe carpetas (CU-21/22/23) |
| Director | login | asigna trámites (CU-09) y ve el tablero (CU-25) |

Lo administrativo (usuarios, roles, auditoría, reportes, catálogos) queda en la web.

## Estructura de `lib/`

```
lib/
├── core/
│   ├── config.dart      # API_URL por --dart-define (por defecto http://10.0.2.2:8080/api)
│   ├── network/         # ApiClient (Dio + Bearer + ApiException)
│   ├── routing/         # GoRouter único con redirect global
│   └── theme/           # AppTheme claro/oscuro, espejo de frontend/app/globals.css
├── shared/widgets/      # widgets usados por 2+ features (se crea con el primero real)
├── features/{feature}/
│   ├── models/          # DTOs a mano (fromJson/toJson)
│   ├── services/        # llamadas a la API de ESTA feature
│   ├── providers/       # Riverpod
│   ├── screens/         # pantallas a las que apunta el router
│   └── widgets/         # solo si hace falta
└── main.dart
```

Features: `auth` (CU-01), `inicio` (bienvenida y accesos), `portal` (CU-12, CU-13), `tramites`
(CU-09, CU-10 y la carga de documentos de CU-15, que siempre ocurre dentro de un trámite), `legajos`
(CU-19), `carpetas` (CU-20 a CU-23), `tablero` (CU-25).

Etiquetas de estados (trámite, solicitud, estado material) que usan dos o más features viven en
`lib/shared/etiquetas.dart`, con los mismos textos que `frontend/lib/tramites.ts` y `frontend/lib/legajos.ts`.

## Regla de dependencia

`features/` depende de `core/` y `shared/`, nunca al revés. `core/network` no importa `features/auth`:
recibe el token por un callback. Una feature no importa otra salvo un modelo justificado: `carpetas`
usa el `Legajo` y el `LegajosService` de `legajos` para elegir la carpeta a solicitar. Entre pantallas
de features distintas se navega por **ruta** (`context.push('/inicio/carpetas/nueva?…')`), no importando
la pantalla.

## Menú por permisos

La app decide qué ve el usuario con los `permisos` que devuelve el login (no por nombre de rol), igual
que `frontend/lib/navegacion.ts`. El inicio muestra un acceso por permiso: `TRAMITE_REVISAR`,
`TRAMITE_ASIGNAR`, `CARPETA_SOLICITAR`, `DOCUMENTO_GESTIONAR`, `LEGAJO_GESTIONAR`,
`CUSTODIA_REGISTRAR`, `TABLERO_VER`.

## iOS sin Mac

No hay Mac en el equipo y virtualizar macOS en hardware no Apple viola su licencia (y choca con el
capítulo de Aspectos Legales). El build de iOS corre en **GitHub Actions (`macos-latest`)**:
`flutter build ios --release --no-codesign` en cada push. Firmar e instalar en un iPhone o TestFlight
requiere cuenta Apple Developer (follow-up). Alternativa equivalente: Codemagic (plan gratuito con macOS).

## Validación estándar de una UT mobile

`flutter analyze` sin issues + `flutter test` en verde (+ prueba manual en emulador cuando es pantalla).
