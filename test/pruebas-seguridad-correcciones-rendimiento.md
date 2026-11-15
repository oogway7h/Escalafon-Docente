# Pruebas de seguridad, corrección de errores y rendimiento

Registro de cómo se revisó el sistema, qué se corrigió y cómo se comprobó cada corrección. Fecha: 14 de noviembre de 2026.

- **Revisión, plan y correcciones:** Claude Fable 5.2 (el plan completo está en `impruves.md`, en la raíz del proyecto).
- **Ejecución de las pruebas de servicios, `go test`, `govulncheck` y k6:** Claude Sonnet 5.5.
- **Entorno:** API en Go local (puerto 8081, porque el 8080 estaba ocupado por otro proceso), base de desarrollo en Neon, Redis y Mailpit en Docker, frontend de producción (`pnpm build` + `pnpm start`).

> Alcance honesto: la revisión de seguridad fue **estática** (lectura del código) más verificaciones dirigidas. No se hizo una prueba de penetración ni se usó un escáner dinámico (ZAP, Burp). Lo que sigue describe lo que realmente se hizo.

---

## 1. Cómo se hicieron las pruebas de seguridad

### 1.1 Revisión del código (análisis estático)

Se leyeron, sin ejecutar, el backend (`cmd/api`, autenticación, middleware, casos de uso de trámites y legajos, repositorios, correo, almacenamiento S3), el `ia-service` (FastAPI) y el frontend (`lib/api.ts`, sesión, `proxy.ts`, acciones del portal, entrega de archivos). Se buscó:

| Qué se buscó | Cómo |
|---|---|
| Inyección SQL | Se revisaron todas las consultas `Raw`, `Where`, `Order` y los `LIKE`. Todas van parametrizadas; la única concatenación (`columna+" >= ?"`) usa nombres fijos del código. |
| Autenticación y sesión | Se siguió el flujo completo: login, JWT (HS256, método e emisor validados), revocación por `jti`, versión de sesión, bloqueo de cuenta. |
| Autorización | Se leyó `router.go` ruta por ruta para comprobar que cada una exige un permiso en la API, no solo en el menú del frontend. |
| Enumeración de usuarios | Se compararon las respuestas de login, recuperación de contraseña y portal público ante usuario existente y no existente. |
| Archivos subidos | Se revisó el tipo, el tamaño y las claves de objeto en S3, y qué valida el `ia-service`. |
| Cabeceras, secretos y contenedores | Se revisaron `next.config.ts`, `.gitignore`, `docker-compose.yml` y los `Dockerfile`. |

De ahí salieron 13 hallazgos (S1–S13), clasificados por severidad en `impruves.md` §1.

### 1.2 Verificaciones dirigidas (dinámicas)

| Verificación | Herramienta | Qué demostró |
|---|---|---|
| Cabeceras de seguridad | `curl -I` contra el frontend y la API | CSP, `X-Frame-Options: DENY`, `nosniff`, `Referrer-Policy` y HSTS presentes en Next; `nosniff` y `X-Frame-Options` en la API. |
| La CSP no rompe nada | Navegador (visor de PDF, login, portal), consola sin errores de CSP | El visor abre PDF y el worker de pdf.js carga. |
| Sin token → 401; rol sin permiso → 403 | Colección Newman (filas 3 y 4) | La API rechaza en el servidor, no solo en la interfaz. |
| Mensaje genérico en credenciales incorrectas | Newman (fila 2) | Código `CREDENCIALES_INVALIDAS`, sin revelar si el correo existe ni si está bloqueado, sin token. |
| El portal no expone datos personales | Newman (fila 6) | La respuesta pública no contiene `cedula` ni `correo`. |
| Formato de archivo no permitido | Newman (fila 8) | Se rechaza con `TIPO_NO_PERMITIDO`. |
| Operación no permitida sobre la bitácora | Newman (fila 11) | No existe ninguna ruta que borre (responde 404). |
| Falla cerrada sin Redis | Prueba unitaria nueva | Con la lista de revocación inaccesible, la API responde 503 y no acepta el token. |
| Bloqueo de cuenta sin filtrar existencia | Prueba unitaria nueva | Con la cuenta bloqueada y una clave incorrecta se responde igual que con un correo inexistente. |
| Límite de tasa | k6 (50 usuarios, 1 minuto) | El portal público corta a 30 peticiones por minuto por IP (2.874 de 2.934 respuestas fueron 429). |
| Dependencias vulnerables | `govulncheck ./...` | 23 vulnerabilidades alcanzables (ver §4). |

---

## 2. Cómo se hizo la corrección de errores

### 2.1 Método

1. Se **trazó el flujo completo** de cada problema antes de editar (por ejemplo, de dónde sale la IP que ve el limitador, o qué consultas arma el detalle de un trámite).
2. Se corrigió **en el lugar común** y no en cada llamador (por ejemplo, un solo `docenteVerificado` para el portal y un solo `plantilla()` para escapar los correos).
3. Cada cambio de comportamiento llevó **una prueba** que falla si la lógica se rompe.
4. Se verificó en tres niveles: pruebas del backend, compilación y lint del frontend, y prueba manual en el navegador.

### 2.2 Seguridad corregida

| Hallazgo | Corrección | Dónde |
|---|---|---|
| S2: el limitador era global | Next reenvía la IP real del navegador en `X-Forwarded-For` | `frontend/lib/api.ts` (`cabecerasCliente`), login y `apiPublica` |
| S3: radicar a nombre de otro docente desde el portal | Se exige código **y** cédula; error único para «no existe», «inactivo» y «cédula distinta»; 5 solicitudes por hora por IP; máximo 2 pendientes por docente; el OCR se lanza al validar, no al radicar | `application/tramites/tramites.go`, `main.go`, `router.go` |
| S4: sesión cerrada válida si cae Redis | `Autenticar` responde 503 en lugar de continuar | `middleware/autenticar.go`, `cache/sesion.go` |
| S5: sin cabeceras de seguridad | CSP, `X-Frame-Options`, `nosniff`, `Referrer-Policy`, `Permissions-Policy`, HSTS en producción; `helmet` en Fiber | `frontend/next.config.ts`, `cmd/api/main.go` |
| S6: el bloqueo delataba qué cuentas existen | La contraseña se comprueba siempre; el bloqueo solo se informa a quien la conoce; el contador se reinicia al vencer el bloqueo | `application/auth/auth.go` |
| S7: la recuperación delataba el correo por tiempo de respuesta | Envío en segundo plano y tope de 3 solicitudes por hora por cuenta | `application/auth/auth.go`, `token.Repository.ContarDesde` |
| S8: contenedores | Puertos publicados solo en `127.0.0.1`; `ia-service` sin privilegios de root | `docker-compose.yml`, `ia-service/Dockerfile` |
| S9: tipo real de la imagen | El `ia-service` comprueba que sea JPG o PNG de verdad | `ia-service/app/procesar.py` |
| S10: nombres sin escapar en correos | `html.EscapeString` en `plantilla()` | `infrastructure/mail/mailer.go` |
| S11: caché del visor | `max-age` baja de 24 h a 10 min | `frontend/lib/archivo.ts` |
| S12: comodines de `LIKE` | `escaparLike` en todas las búsquedas | `infrastructure/persistence/repositorio_*.go` |

### 2.3 Errores de funcionamiento corregidos

| Error | Corrección |
|---|---|
| Una regularización pedía categoría | `tramite.categoria_id` pasa a ser opcional (migración `20261114100000`, con `CHECK` que la exige solo en categorizaciones); el campo desaparece de ventanilla y del portal; los listados y reportes usan `LEFT JOIN`. En el portal, la regularización pide de 1 a 5 adjuntos libres. |
| Los campos de archivo no eran arrastrar y soltar | Nuevo componente `ZonaArchivo.tsx`: un `<label>` sobre un `<input type="file">` nativo (clic y teclado funcionan), más arrastrar y soltar, con validación de tipo y tamaño. |
| No se podía solicitar una carpeta sin trámite abierto | `tramiteId` pasa a ser opcional; si se envía, se mantiene la comprobación de que sea del solicitante, esté abierto y sea del mismo docente. |
| Mensajes de tamaño contradictorios (8, 20 y 400 MB) | `LIMITE_MB` en un solo módulo del frontend. |
| Lista de trámites con permisos distintos a los de la API | Se alineó `TRAMITE_ARCHIVAR`. |

### 2.4 Problemas que aparecieron durante la corrección

- **`pnpm build` falló al prerenderizar `/solicitud`** después de cachear el catálogo de categorías: la API no existe durante el build. Se resolvió llamando a `connection()` para que las páginas sigan siendo dinámicas y solo se cachee el dato.
- **`atlas.sum` no coincidía** porque el repositorio tiene finales de línea CRLF y Atlas calcula sobre LF. Se calculó el hash sobre una copia en LF y se validó aplicando la migración primero a `postgres-test` y después a la base de desarrollo.
- **Las pruebas de integración fallaban** con «la relación no existe» hasta aplicar las migraciones a `postgres-test`.
- **Newman ejecutó 0 solicitudes** la primera vez («request url is empty»): el formato de URL como objeto no era válido. Como ninguna llegó a salir, no hubo efectos en la base.
- **Un proceso ajeno ocupaba `localhost:8080`** y respondía 400 en lugar de la API. Todas las pruebas se hicieron en el 8081.

---

## 3. Cómo se hizo la mejora de rendimiento

### 3.1 Mejoras generales

| Mejora | Qué se hizo |
|---|---|
| Pool de conexiones | `SetMaxOpenConns(20)`, `SetMaxIdleConns(5)` y tiempos de vida acotados (Neon cierra las conexiones inactivas). |
| Compresión | `compress` de Fiber (el texto OCR de un documento pesa varios MB). |
| Sesión | `GET /api/auth/sesion` toma los permisos del propio JWT y se ahorra una consulta por navegación del panel. |
| Catálogo público | Las categorías se guardan 5 minutos en la caché de datos de Next y se invalidan al editarlas. |
| Refresco de OCR | Se pausa con la pestaña oculta y baja de 6 s a 10 s. |
| Diseño | Una sola fuente variable, sin vidrio ni degradados ni animaciones de scroll. |

### 3.2 El visor de PDF y la búsqueda dentro de un documento

**Síntoma reportado:** buscar «maestria» dentro de un documento ya abierto tardaba en mostrar la página, y se pidió una barra de progreso de descarga.

**Diagnóstico** (leyendo `VisorPdf.tsx` y `VisorDocumento.tsx`; no se perfiló antes del cambio):

1. Cada clic en una coincidencia era un enlace que **recargaba toda la página en el servidor**.
2. El salto era un desplazamiento animado: a una página lejana hacía **dibujar y cancelar** todas las intermedias.
3. Se dibujaba una **miniatura por cada página con coincidencia** (hasta 24); pdf.js tiene un solo hilo de trabajo y la página buscada esperaba detrás de ellas.
4. Cambiar la búsqueda **redibujaba** las páginas cercanas solo para recalcular el resaltado.

**Cambios:**

- Nuevo `VisorConBusqueda.tsx`: buscar y saltar entre coincidencias ocurre en el navegador; la búsqueda pide solo las páginas a una ruta nueva (`.../documentos/[docId]/buscar`) y la URL se refleja con `history.replaceState`.
- El salto a una página lejana es instantáneo; solo se anima el paso a la vecina.
- El riel muestra solo las miniaturas vecinas, con retardo para no competir con la página principal.
- El resaltado se calcula en un efecto aparte, sin tocar el lienzo.
- **Barra de progreso de descarga** con porcentaje y MB (`loadingTask.onProgress`).

**Medición** (build de producción, PDF real de 78 páginas, búsqueda «maestria», 6 coincidencias):

| Salto | Página en vista | Dibujada | Resaltada |
|---|---|---|---|
| 18 → 68 | 29 ms | 517 ms | 740 ms |
| 68 → 18 | 29 ms | 281 ms | 88 ms |
| 18 → 77 | 27 ms | 107 ms | 634 ms |

Límites: no se midió la versión anterior (no hay cifra de «antes»), y la barra de progreso apenas se vio porque el archivo se sirve desde la misma máquina; falta probarla con un escaneo de cientos de MB.

---

## 4. Resultados de las pruebas posteriores

Detalle completo en la conversación; aquí el resumen.

### 4.1 Colección de servicios (Newman)

Archivos para repetirla: `test/servicios/coleccion.json` y `test/servicios/entorno.json`.

```
pnpm dlx newman run test/servicios/coleccion.json -e test/servicios/entorno.json \
  --reporters cli,json,junit --reporter-json-export resultados-newman.json --reporter-junit-export resultados-newman.xml
```

Requiere la API en `http://127.0.0.1:8081`, Redis activo y el seed cargado (`make seed`). Para la fila 10 usa el legajo 2397 (si ya no existe o está prestado, cambiar el número en la colección). La limpieza devuelve la carpeta al archivo, pero deja asientos en la bitácora de custodia, que es de solo anexión.

| N.° | Solicitud | CU | Resultado | Tiempo (ms) |
|---|---|---|---|---|
| 1 | Login correcto | CU-01 | **Incorrecto** (solo el tiempo: pide < 1.000) | 1337 |
| 2 | Login con contraseña incorrecta | CU-01 | Correcto | 1219 |
| 3 | Radicar sin token | CU-08 | Correcto (401) | 2 |
| 4 | Radicar como Analista | CU-08 | Correcto (403) | 4 |
| 5 | Radicar como Ventanilla | CU-08 | Correcto (responde 201, no 200) | 3100 |
| 6 | Consulta pública por código | CU-12 | Correcto | 1080 |
| 7 | Código inexistente | CU-12 | Correcto (404) | 138 |
| 8 | Archivo de formato no permitido | CU-15 | **Incorrecto** (rechaza con 422, no 400) | 1210 |
| 9 | Buscar por contenido | CU-18 | Correcto | 142 |
| 10 | Despachar legajo ya prestado | CU-22 | Correcto (409, bitácora sin cambios) | 137 |
| 11 | Eliminar entrada de la bitácora | CU-24 | **Incorrecto** (404 en lugar de 403) | 2 |

### 4.2 `go test` y cobertura

`go test ./... -count=1 -coverprofile=cobertura.out` (con `DATABASE_URL_TEST` definida; sin ella las pruebas de repositorio se omiten en silencio): **174 correctas, 0 fallidas, 2 omitidas; cobertura total 36,5 %**. Más cubiertos: reportes 95 %, indicadores 92–95 %, seguridad 91 %. Menos cubiertos: categorías 24 %, rutas HTTP 23 %; `domain`, `handlers`, `dto`, `config`, `mail`, `storage` e `iaservice` no tienen pruebas propias.

### 4.3 `govulncheck`

23 vulnerabilidades alcanzables: 6 en la librería estándar (Go 1.27.1 → 1.27.2), 5 en `golang.org/x/net` (→ v0.60.0), 8 en `excelize` (→ v2.11.x), 2 en `grpc`, 1 en `brotli` y 1 en `go-jose` (esta última entra por la herramienta de migraciones, no por la API). **No están corregidas.**

### 4.4 Carga con k6 (50 usuarios, 1 minuto)

Se apuntó a `/api/publico/tramites/:codigo`, el endpoint real (el ejemplo original usaba `/consulta`, que no existe). El umbral `p(95) < 1000 ms` pasa (6,66 ms), pero solo porque 2.874 de 2.934 respuestas fueron 429 del límite de tasa. Las 60 que sí respondieron 200 tardaron 2,38 s de promedio y 2,80 s en p95. **La prueba no demuestra que el portal aguante 50 usuarios simultáneos.**

---

## 5. Pendiente y riesgos aceptados

| Pendiente | Motivo |
|---|---|
| Actualizar Go y las dependencias vulnerables (§4.3) | No se hizo; es el resultado más grave de las pruebas. |
| Rotar la llave de AWS (`ded-backend_accessKeys.csv`) | Decisión del responsable: no se rota. Riesgo aceptado. |
| Contraseña de Redis y límites del `ia-service` | Cambiarían el `.env` y podrían cortar el OCR de PDF grandes. |
| Borrar del bucket los objetos que exceden el tamaño | El almacén no tiene operación de borrado; mejor una regla de ciclo de vida en S3. |
| Latencia de login (1,3 s), radicar (3,1 s) y consulta pública (1,1 s) | Causa probable: viaje a Neon y consultas en serie (el detalle de un trámite hace unas 8). No se midió ni se paralelizó. |
| Criterios de las filas 8 y 11 | Decidir si el servicio cambia a 400 y 403, o si el criterio pasa a «4xx» y «sin ruta». |
| Prueba de carga real del portal | Repetir k6 desde varias IP o con el límite de tasa ajustado en un entorno de pruebas. |
| Modo oscuro y móvil del rediseño | Se interrumpió la captura de pantalla del navegador de pruebas. |
| Barra de progreso con un documento grande | Solo se probó con uno pequeño servido en local. |
| Flujo de GitHub Actions | No se creó: `backend/` es un repositorio propio y la carpeta de trabajo del ejemplo habría que ajustarla. |

## 6. Efectos en la base de desarrollo

Trámites de prueba creados: `ED-2026-MKHMXL` (regularización) y `ED-2026-2SN2WX` (categorización). La carpeta del legajo 2397 se solicitó, despachó y devolvió (queda en archivo, con asientos en la bitácora). La migración `20261114100000` está aplicada.
