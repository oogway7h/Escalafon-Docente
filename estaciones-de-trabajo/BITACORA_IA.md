# Bitácora de uso de IA

Evidencia de cómo se usó la IA (Claude Code, modelo Claude Opus) para acelerar el proceso con SDD.
Cada entrada registra **qué se pidió**, **qué produjo la IA**, **cómo se verificó** y **qué decidió la
persona**. Las horas "manual" son una **estimación** de lo que habría tomado hacerlo a mano; no una
medición.

## Reglas del uso de IA en este proyecto

- La IA trabaja **siempre** con las reglas de `.claude/rules/` cargadas como contexto fijo.
- La persona decide el alcance (CU y UT); la IA no agrega funcionalidad fuera de una UT.
- Nada se da por terminado sin la **Validación** de la UT (tests, `flutter analyze`, prueba manual).
- La IA no inventa contratos: los endpoints, DTOs y códigos de error se leyeron del código real
  (`backend/internal/interfaces/http/router.go`, `dto/*.go`, `apierror.go`).

## Entradas

| # | Fecha | Fase SDD | Pedido a la IA | Resultado | Verificación | Decisión humana | Manual (est.) |
|---|---|---|---|---|---|---|---|
| 1 | 2026-11-08 | Análisis | Estudiar la metodología de `examen2/DriveHub` (CU → UT por capa) y el estado del proyecto DED | Resumen de la metodología y de lo que ya existe (26 CU hechos en web) | Lectura de `.claude/rules/` de DriveHub y de `docs/estado-del-*.md` | Elegir carpeta `estaciones-de-trabajo`, UT condensadas para la web, 4 perfiles para la app, iOS por CI | 4 h |
| 2 | 2026-11-08 | Reglas | Adaptar las reglas de DriveHub al stack DED (Go, Next, Flutter) | 13 archivos en `.claude/rules/` | Revisión de que cada regla cite la API real (sin sobre, `{codigo,mensaje}`) | Se descartan microservicios, gateway y diseño Cockpit de DriveHub | 6 h |
| 3 | 2026-11-08 | Specify (retroactivo) | Escribir los 26 CU con la plantilla, a partir de `cu.md` y del código | 26 archivos en `CasosDeUso/` con flujos, códigos de error reales y CA | Script de links relativos; códigos de error contrastados con `apierror.go` | Se documentan los desvíos conscientes respecto de `cu.md` (CU-01, 11, 15, 22, 23) | 16 h |
| 4 | 2026-11-08 | Plan (retroactivo) | Partir cada CU en UT de backend y frontend que citen los archivos existentes | 163 archivos (52 READMEs + 111 UT) generados desde datos con un script | Script que comprueba que **todos** los archivos citados existen en `backend/` y `frontend/` (0 faltantes) | UT condensadas (2-5 por capa) en vez de las 10-17 de DriveHub | 30 h |
| 5 | 2026-11-08 | Plan (spec-first) | Escribir las UT-MOBILE de 12 CU **antes** del código | 12 READMEs + 28 UT-MOBILE con alcance, archivos, fuera de alcance y validación | Revisión contra los DTOs reales (`dto/tramites.go`, `dto/legajos.go`, `dto/indicadores.go`) | CU-10 y CU-19 solo consulta en móvil; CU-15 dentro de la feature `tramites` | 10 h |

| 6 | 2026-11-08 | Implement | «Implementa UT-CU01-MOB-01..06» | Proyecto Flutter, `ApiClient`, tema, sesión cifrada, router con redirect, login e inicio por permisos | `flutter analyze` 0 issues; 16 tests (cliente HTTP, AuthService, sesión, redirect) | Sin refresh token (la API no lo tiene): una sesión vencida o revocada vuelve a la bienvenida | 14 h |
| 7 | 2026-11-08 | Implement | «Implementa UT-CU12-MOB-01..02» | Seguimiento público por código o código + cédula | Tests de la cédula en el cuerpo y no en la URL | Mismo umbral de 5 caracteres que la web | 4 h |
| 8 | 2026-11-08 | Implement | «Implementa UT-CU13-MOB-01..02» | Solicitud pública con cámara, galería o PDF por requisito | Test de los campos exactos del multipart (`doc_<id>`) | `ArchivoAdjunto` va a `shared/` porque CU-15 también lo usa (se actualizó la UT) | 6 h |
| 9 | 2026-11-08 | Implement | «Implementa UT-CU10-MOB-01..02» | Bandeja y detalle del trámite | Tests de parseo del detalle y de `mios=true` | `LineaFases` pasa a `shared/widgets` al tener un segundo uso (se actualizó la UT de CU-12) | 5 h |
| 10 | 2026-11-08 | Implement | «Implementa UT-CU09-MOB-01..02» | Por asignar, hoja de analistas por carga, reasignar con motivo | Tests del cuerpo con y sin motivo | — | 4 h |
| 11 | 2026-11-08 | Implement | «Implementa UT-CU15-MOB-01..02» | Subida directa a S3 con URL firmada y confirmación | Test: el PUT va a la URL firmada **sin** `Authorization` y con `Content-Type`; si la API rechaza no se toca S3 | Lectura completa en memoria marcada como deuda (`ponytail:`) para PDF muy grandes | 6 h |
| 12 | 2026-11-08 | Implement | «Implementa UT-CU19-MOB-01..02» | Búsqueda de legajos paginada | Test de query y parseo | — | 3 h |
| 13 | 2026-11-08 | Implement | «Implementa UT-CU20 a CU23-MOB» | Mis carpetas, solicitar desde el trámite, pantalla Custodia con pestañas Por atender / Aprobadas / Prestadas, hojas de entrega y devolución | Tests de cada endpoint, de «vencida» y de la regla de observación obligatoria | Una pantalla con pestañas como `/panel/custodia` de la web | 12 h |
| 14 | 2026-11-08 | Implement | «Implementa UT-CU25-MOB-01..02» | Tablero con periodo, indicadores, barras y alertas de mora | Test de parseo y del periodo AAAA-MM-DD | Barras propias: sin librería de gráficos | 5 h |
| 15 | 2026-11-08 | CI | Compilar iOS sin Mac | Workflow de GitHub Actions: Android (analyze, test, APK) e iOS (`--no-codesign`) | Push a `oogway7h/DED-Mobile` | Se descarta Hackintosh por la licencia de macOS | 3 h |

**Total:** 45 tests en verde y `flutter analyze` sin issues en el último commit.

**Pendiente de verificar por la persona:** los pasos de Validación que dicen «En el emulador Android con la
API local» (login real, subir una foto, entregar y recibir una carpeta) no se ejecutaron todavía: en la
sesión de desarrollo Docker no estaba corriendo (la API necesita Redis). El resultado del workflow de iOS
se revisa en la pestaña *Actions* del repositorio.
