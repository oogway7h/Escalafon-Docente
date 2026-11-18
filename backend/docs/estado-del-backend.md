# Backend — estado de lo implementado

API REST en Go del Sistema de Escalafón Docente (UAGRM). Stack: Fiber v2, GORM + Postgres (Neon), Redis, JWT (HS256) + bcrypt, Atlas (migraciones), SMTP (Mailpit en desarrollo), S3 para adjuntos.

> Última actualización: 2026-10-25.

## 1. Casos de uso cubiertos

| CU | Nombre | Estado |
|---|---|---|
| CU-01 | Iniciar sesión | ✅ |
| CU-02 | Recuperar contraseña | ✅ |
| CU-03 | Editar perfil personal | ✅ |
| CU-04 | Administrar cuentas de usuario | ✅ |
| CU-05 | Gestionar roles y permisos | ✅ |
| CU-06 | Consultar registro de auditoría | ✅ |
| CU-07 | Gestionar docentes y categorías | ✅ |
| CU-08 | Radicar trámite | ✅ |
| CU-09 | Asignar o reasignar trámite | ✅ |
| CU-10 | Revisar expediente y registrar observaciones | ✅ |
| CU-11 | Emitir dictamen y remitir el expediente (el archivo lo hace Personal de Archivo) | ✅ |
| CU-12 | Consultar estado de trámite (público) | ✅ |
| CU-13 | Radicar solicitud desde el portal público | ✅ (adjuntos requieren S3) |
| CU-19 | Gestionar inventario de legajos | ✅ |
| CU-20 | Solicitar legajo físico | ✅ |
| CU-21 | Atender solicitud de legajo | ✅ |
| CU-22 | Despachar carpeta física (check-out) | ✅ (sin evidencia fotográfica) |
| CU-23 | Registrar devolución de carpeta (check-in) | ✅ (sin evidencia fotográfica: solo estado material y observación) |
| CU-24 | Consultar bitácora de custodia | ✅ |
| CU-14 | Notificar cambios de estado al docente | ✅ |
| CU-15 | Cargar y procesar documento con OCR | ✅ (subida directa a S3; OCR y compresión en `ia-service`). **Expediente digital del docente**: el Personal de Archivo (`DOCUMENTO_ESCANEAR`) sube los escaneos de la carpeta física por docente y tipo (Currículum, Regularización, Adjunto, Asociado, Titular A/B/C, Emérito) en `/api/docentes/:id/documentos` (`subida`, `:docId/confirmar`, `:docId/estado` = baja, lista, `:docId` = URL firmada); se guardan en `tramite_documento` con `tramite_id` nulo y `docente_id`/`tipo_documento` (`application/tramites/expediente.go`). Archivo, analistas y Director (`DOCUMENTO_BUSCAR`) buscan por contenido en todos los documentos del docente: `GET /api/docentes/:id/documentos/buscar?q=` devuelve documento + página + fragmento resaltado. El texto es buscable al terminar el OCR (la corrección es opcional). Los PDF admiten hasta 400 MB (desvío consciente del límite de 20 MB del CU, que rige para imágenes). El hash del `ia-service` se coteja con el registrado al cargar (adjuntos del portal) y `Reanudar` no duplica trabajos en curso. Tras `make seed` hay que reiniciar sesión para ver el permiso nuevo. |
| CU-16 | Dar de baja documento | ✅ (baja lógica; solo el analista asignado; el objeto se conserva en S3) |
| CU-25 | Visualizar tablero de indicadores | ✅ (`GET /api/tablero`; ver sección 5 bis) |
| CU-26 | Generar y exportar reportes | ✅ (`POST /api/reportes`: PDF con membrete o Excel con autofiltro; queda en auditoría) |
| CU-17 | Restituir documento dado de baja | ✅ Solo el Director (`DOCUMENTO_RESTITUIR`). `GET /api/documentos/bajas?buscar=&pagina=` lista las bajas de trámites y de expedientes (docente, trámite, motivo, quién y cuándo); `POST /api/documentos/:docId/restituir {motivo}` (mín. 5) reactiva el documento con un `UPDATE` condicional, borra los datos de la baja y deja `RESTITUIR_DOCUMENTO` en auditoría con ambos motivos. No toca trámites remitidos o archivados (solo lectura). Las búsquedas filtran por `activo`, así que vuelve a ellas sin reprocesar el OCR (`application/tramites/restitucion.go`, `persistence/repositorio_documentos.go`). |
| CU-18 | Buscar documentación digital | ✅ `GET /api/documentos/buscar?q=&campo=&tipo=&facultad=&orden=&pagina=` (`DOCUMENTO_BUSCAR`) busca en todos los documentos activos (adjuntos y escaneos): `campo` vacío = contenido o docente, `CONTENIDO` = texto OCR por página (índice `es_unaccent`; los números de resolución son palabras del texto), `DOCENTE` = nombre, código o cédula (con dígitos, exacto). `facultad` filtra por la facultad del docente y `orden=recientes` pone primero los documentos más nuevos (por defecto, relevancia). Agrupa por documento y devuelve las páginas halladas y el mejor fragmento marcado con ⟦ ⟧. Los dados de baja nunca salen (`application/tramites/busqueda.go`). Para el buscador del panel, que pide primero al docente, `GET /api/docentes/sugerencias?q=` devuelve hasta 8 docentes activos mientras se escribe: código o cédula **por prefijo** (a diferencia del listado, que es exacto) y nombre por fragmento sin tildes y en cualquier orden, sin legajos ni conteo para responder rápido. |

El mailer envía correos de activación, recuperación y cambio de estado del trámite (CU-14). Constancias no es un caso de uso: su permiso se quitó del seed (en una BD ya sembrada la fila `CONSTANCIA_EMITIR` queda huérfana y sin efecto).

## 2. Arquitectura

```
interfaces/http  →  application  →  domain
infrastructure   →  application  →  domain
```

| Capa | Contenido |
|---|---|
| `internal/domain` | Entidades, reglas y puertos de repositorio: `usuario`, `rol`, `token`, `auditoria`, `docente`, `categoria`, `tramite`, `legajo`, `notificacion`, `indicador` (CU-25/26: regla de mora, agrupación de tiempos y contenido de los reportes). No importa GORM, Fiber ni Redis. |
| `internal/application` | Un struct por caso de uso: `auth`, `perfil`, `usuarios`, `roles`, `auditorias`, `docentes`, `categorias`, `tramites`, `legajos`, `notificaciones`, `indicadores`. Puertos externos en `ports` (`TokenService`, `PasswordHasher`, `Mailer`, `SessionCache`, `TxManager`, `Clock`, `ObjectStorage`, `GeneradorToken`, `ProcesadorDocumentos`, `ExportadorReportes`). Dobles en `testdobles`. |
| `internal/infrastructure` | `persistence` (repositorios GORM y modelos), `cache` (sesión en Redis), `security` (JWT, bcrypt, tokens), `mail` (SMTP), `storage` (S3), `iaservice` (OCR), `reportes` (PDF con `go-pdf/fpdf` y Excel con `excelize`), `config`. |
| `internal/interfaces/http` | Router, middlewares, handlers, DTO y `apierror` (traduce errores de dominio a HTTP). |

Puntos de entrada: `cmd/api` (único composition root), `cmd/seed`, `cmd/importdocentes`, `cmd/atlasloader`.
`make arch` verifica que `domain` y `application` no importen capas externas.

**Sesión:** `Autenticar` valida el JWT localmente y consulta Redis (`version_sesion`, revocación por `jti`); si Redis no responde, cae a Postgres. **Sin Redis la API no arranca** (el limitador de tasa lo necesita). Tres intentos fallidos bloquean la cuenta temporalmente.

## 3. Modelo de datos

| Tabla | Notas |
|---|---|
| `rol`, `permiso`, `rol_permiso`, `usuario`, `token_recuperacion` | Seguridad. Los roles de sistema no se eliminan ni renombran. |
| `auditoria` | Solo anexión (trigger). El dominio separa el puerto de escritura (`auditoria.Repository`, solo `Registrar`) del de lectura (`auditoria.Consulta`, CU-06). |
| `categoria`, `categoria_requisito` | Catálogo normativo y documentos exigidos. Los cargos (Decano, Rector…) existen como categorías inactivas. |
| `facultad` | Facultad, carrera o decanato. |
| `docente` | `codigo` único (identificador principal), `cedula` opcional e inmutable una vez asignada, nombres y apellidos separados, categoría, facultad, baja lógica. |
| `legajo` | Carpeta física: `codigo` único (`LEG-NNNNN`, trigger `fn_legajo_codigo` si va vacío), gaveta, orden, anexos, observación y `estado`. Varios por docente. |
| `solicitud_carpeta` | PENDIENTE → ATENDIDA → DESPACHADA → DEVUELTA, o RECHAZADA. Índice único parcial `ux_solicitud_abierta`: una solicitud abierta por legajo. |
| `movimiento_custodia` | Salidas/entradas con custodio, fecha límite y estado material. Solo anexión (trigger). |
| `tramite`, `tramite_historial`, `tramite_documento`, `observacion` | Trámites; el historial es de solo anexión. En `tramite_documento`, `activo=false` es la baja lógica de CU-16 (con `motivo_baja`, `baja_en` y `baja_por`): el documento deja de cubrir el requisito, pero el objeto sigue en el almacén. |

### Migraciones aplicadas (`migrations/`)

| Archivo | Contenido |
|---|---|
| `…_inicial` | Seguridad y auditoría. |
| `…_ajustes_manuales` | Índices únicos y triggers escritos a mano (correo único sin distinguir mayúsculas, auditoría de solo anexión, `actualizado_en`). |
| `…_tramites` | Docentes, categorías, trámites, historial, documentos y observaciones. |
| `…_padron` | `codigo` y `facultad` en docente, cédula opcional, tablas `facultad` y `legajo`. |
| `…_legajos_custodia` | `idx_auditoria_fecha`, código y estado en `legajo`, tablas `solicitud_carpeta` y `movimiento_custodia` (escrita a mano; validada contra Postgres 16). |
| `…_notificacion_envio` | Tabla `notificacion_envio` (CU-14: aviso, estado PENDIENTE/ENVIADO/FALLIDO, intentos y error). |
| `…_custodia_evidencia` | `evidencia_clave` en `movimiento_custodia` (foto de la devolución). |
| `…_quitar_evidencia_custodia` | Elimina `evidencia_clave`: la devolución ya no admite foto. |
| `…_documento_baja` | `activo`, `motivo_baja`, `baja_en` y `baja_por` en `tramite_documento` (CU-16; escrita a mano y validada contra Postgres 16). |
| `…_documento_paginas` | Texto OCR por página: tabla `documento_pagina (documento_id, pagina, texto)` con índice GIN de texto completo (manual: Atlas querrá borrarlo). `texto_ocr` queda como la concatenación de las páginas; lo ya procesado se parte por el salto de página. Servirá también a CU-18. |
| `…_busqueda_sin_tildes` | Extensión `unaccent` y configuración de texto `es_unaccent` (español sin tildes); el índice por página pasa a usarla. |
| `…_documento_ocr` | CU-15: `estado_ocr`, `texto_ocr`, `texto_validado`, `error_ocr`, `clave_comprimida`, `tamano_comprimido` en `tramite_documento`; `hash_sha256` pasa a opcional (se calcula tras la subida); índice GIN de texto completo (`ix_tramite_documento_texto_fts`, manual: Atlas querrá borrarlo). |
| `…_expediente_digital_docente` | `tramite_id` pasa a opcional en `tramite_documento`; añade `docente_id` (FK a `docente`) y `tipo_documento`, con CHECK de un solo dueño y de tipos válidos; elimina `ix_tramite_documento_texto_fts` (ya no se consulta). Escrita a mano; aplicar con `make migrate-apply`. |

**Convención:** nombre descriptivo en snake_case (`make migrate-diff name=agregar_tabla_notificaciones`; el comando falla sin `name=`), nunca editar una migración aplicada y revisar el SQL antes de aplicar. Atlas quiere borrar los índices manuales `idx_token_recuperacion_pendiente` y `ux_usuario_correo` porque no están en los modelos: quitar esos `DROP INDEX` del SQL generado. Después, `make migrate-hash` y `make migrate-apply`.

## 4. Trámites

**Estados** (`domain/tramite`): `Ingresado → Asignado → En revisión ⇄ Observado → Dictaminado → En facultad → Revisión final → Remitido → Archivado`. **En facultad** es el tramo sin seguimiento (Facultad, Tribunal, Decano y Consejo Facultativo): solo se registran la salida y el regreso. «Observado» es un desvío de «En revisión».

- Código de seguimiento `ED-AAAA-XXXXXX` (alfabeto sin 0/O/1/I), generado con `crypto/rand`.
- La reasignación no cambia de estado y exige un motivo, que queda en el historial.
- Asignado → En revisión exige **una solicitud de carpeta ligada al trámite** (no rechazada); En revisión → Dictaminado exige haber registrado los **años a evaluar** (`PUT /tramites/:id/periodo`) y no tener observaciones pendientes.
- El analista asignado revisa, observa, dictamina, **remite a la facultad** y **remite al Rectorado**; el **regreso** con la resolución (número y fecha) lo registra el analista asignado o quien tenga `TRAMITE_RADICAR` (ventanilla) y el trámite vuelve **al mismo analista** (Revisión final).
- **Archiva el Personal de Archivo** (`TRAMITE_ARCHIVAR`) un trámite Remitido y solo si todas sus carpetas ya volvieron al archivo. Remitido y Archivado son de solo lectura para el analista.
- La carga vigente cuenta Asignado, En revisión, Observado, Dictaminado y Revisión final (no En facultad). Se puede reasignar (con motivo) en cualquier estado salvo Remitido y Archivado.
- El dictamen es un informe interno: el docente ve el **Resultado** recién desde Remitido. CU-14 avisa en Observado, En facultad, Remitido y Archivado.
- Las solicitudes de carpeta llevan `tramiteId` (obligatorio; el trámite debe ser del solicitante, no estar cerrado y la carpeta ser del mismo docente).
- No se dictamina con observaciones pendientes; al resolver la última, el trámite vuelve a *En revisión*.
- Migración `…_flujo_escalafon` (estados ampliados a `varchar(20)`, años a evaluar, resolución, remisión y `solicitud_carpeta.tramite_id`).
- Una solicitud del portal queda pendiente de validación y no se puede asignar hasta que ventanilla la valide.
- Adjuntos del portal: PDF/JPG/PNG de hasta 8 MB, validados por contenido, con hash SHA-256, en S3. Sin `S3_BUCKET` la API arranca pero los adjuntos responden 503.
- Claves de objeto: `docentes/<código del docente>/<id del trámite>/<id del documento>-<nombre>` (portal y analista). La versión comprimida va junto al original como `<id>-comprimido.pdf`.
- Documentos del analista (CU-15): el archivo **no pasa por la API**. `POST /tramites/:id/documentos/subida` valida (PDF ≤ 400 MB, JPG/PNG ≤ 20 MB, analista asignado, trámite no archivado) y devuelve una URL firmada de S3; el navegador hace `PUT` directo; `POST …/:docId/confirmar` comprueba el objeto (HeadObject) y lanza el procesamiento en segundo plano contra `ia-service` (`IA_SERVICE_URL`), que descarga de S3 en streaming, calcula el SHA-256 del original, extrae el texto **por página** (Tesseract; en PDF, `pdftotext` sobre el PDF ya procesado) y comprime el PDF con ocrmypdf. Las páginas se guardan en `documento_pagina`. Estados: `SUBIENDO → PENDIENTE → PROCESANDO → PROCESADO | ERROR` (`NO_APLICA` en los adjuntos del portal). Al arrancar, la API reanuda los PENDIENTE/PROCESANDO (los trabajos del ia-service viven en memoria; supone una sola instancia de la API).

## 5. Endpoints

Todos bajo `/api`. Los de sesión usan el JWT; las lecturas compartidas usan `RequiereAlguno`.

| Grupo | Rutas | Permiso |
|---|---|---|
| Autenticación (con límite de 10 peticiones/min por IP) | `POST /auth/login`, `/auth/recuperacion`, `/auth/restablecer`, `/auth/activar`; con sesión: `POST /auth/logout`, `GET /auth/sesion` | públicas / sesión |
| Acuerdo de usuario | `POST /auth/terminos` (sesión): registra que el usuario aceptó la versión vigente (`usuario.VersionTerminos`; columnas `terminos_version` y `terminos_aceptados_en`) y audita `ACEPTAR_TERMINOS`. `POST /auth/login` y `GET /auth/sesion` devuelven `terminosAceptados`; el frontend muestra el acuerdo antes del panel mientras sea `false`. Subir la versión lo vuelve a pedir a todos | sesión |
| Perfil | `GET/PUT /perfil`, `PUT /perfil/contrasena` | `PERFIL_GESTIONAR` |
| Usuarios | `GET/POST /usuarios`, `GET/PUT /usuarios/:id`, `PATCH /usuarios/:id/rol`, `PATCH /usuarios/:id/estado`, `POST /usuarios/:id/reenviar-activacion` | `USUARIO_GESTIONAR` |
| Roles | `GET /permisos`, `GET/POST /roles`, `GET/PUT/DELETE /roles/:id` | `ROL_GESTIONAR` |
| Auditoría (CU-06, solo lectura) | `GET /auditoria?usuarioId=&accion=&desde=AAAA-MM-DD&hasta=AAAA-MM-DD&pagina=1&tamano=20` devuelve `{datos,total,pagina,tamano}` (del más reciente al más antiguo; filtros combinables; fechas en hora de Bolivia, ambos extremos inclusivos); `GET /auditoria/filtros` devuelve `{acciones, usuarios}` (catálogo de acciones y usuarios que figuran en la bitácora). Errores 422: `VALIDACION`, `ACCION_INVALIDA`, `RANGO_FECHAS_INVALIDO` | `AUDITORIA_CONSULTAR` (el Director lo tiene por seed) |
| Docentes | `GET/POST /docentes` (con `?pagina=N&tamano=10&buscar=` devuelve `{datos,total,pagina,tamano}`; sin `pagina`, el arreglo completo), `GET/PUT /docentes/:id`, `PATCH /docentes/:id/estado`, `GET /facultades` | lectura: `DOCENTE_GESTIONAR`, `TRAMITE_RADICAR` o `LEGAJO_GESTIONAR`; escritura: `DOCENTE_GESTIONAR` |
| Categorías | `GET/POST /categorias`, `GET/PUT /categorias/:id`, `PATCH /categorias/:id/estado`, `DELETE /categorias/:id` (borra de la BD; responde 409 `CATEGORIA_EN_USO` con los conteos si algún docente o trámite la usa) | lectura: `CATEGORIA_GESTIONAR`, `DOCENTE_GESTIONAR` o `TRAMITE_RADICAR`; escritura: `CATEGORIA_GESTIONAR` |
| Legajos (CU-19) | `GET /legajos?buscar=&estado=&pagina=`, `GET /legajos/:id`; `POST /legajos`, `PUT /legajos/:id` (el estado solo lo cambian los movimientos) | lectura: `LEGAJO_GESTIONAR`, `CARPETA_SOLICITAR` o `CUSTODIA_REGISTRAR`; escritura: `LEGAJO_GESTIONAR` |
| Solicitudes (CU-20, CU-21) | `POST /legajos/:id/solicitudes` `{motivo}` (409 `LEGAJO_NO_DISPONIBLE`, `SOLICITUD_ABIERTA`); `GET /solicitudes-carpeta?estado=` (quien solo puede solicitar ve únicamente las suyas); `POST /solicitudes-carpeta/:id/atender` `{aprobar,nota}` (la nota es obligatoria al rechazar) | `CARPETA_SOLICITAR` / lectura como arriba / `LEGAJO_GESTIONAR` |
| Custodia (CU-22 al CU-24) | `POST /solicitudes-carpeta/:id/despachar` `{fechaLimite AAAA-MM-DD, estadoMaterial, observacion}`; `POST /solicitudes-carpeta/:id/entregar` (mismo cuerpo; exige `LEGAJO_GESTIONAR` y `CUSTODIA_REGISTRAR`: aprueba la solicitud pendiente y la despacha en un paso, validando fecha, estado material y disponibilidad antes de aprobar; una ya aprobada se entrega directamente); `POST /legajos/:id/devolver` `{estadoMaterial, observacion}` (422 `OBSERVACION_REQUERIDA` si no vuelve en buen estado y no hay observación) ; `GET /custodia/bitacora?buscar=&tipo=&vencidos=true&pagina=` (cada fila indica si la salida está vencida) | `CUSTODIA_REGISTRAR`; bitácora: `BITACORA_CONSULTAR` |
| Notificaciones (CU-14, solo lectura) | `GET /notificaciones?estado=&pagina=` devuelve `{datos,total,pagina,tamano}` (avisos de Observado, En facultad, Remitido y Archivado enviados al docente) | `NOTIFICACION_CONSULTAR` (Director) |
| Tablero (CU-25, solo lectura) | `GET /tablero?desde=AAAA-MM-DD&hasta=AAAA-MM-DD` → `{generadoEn, periodo, diasMaximosPrestamo, tramites:{radicados,enCurso,resueltos,promedioDias}, porEstado, tiempos, analistas, prestamosVigentes, alertas}`. 422 `VALIDACION`, `RANGO_FECHAS_INVALIDO` | `TABLERO_VER` (Director) |
| Reportes (CU-26) | `POST /reportes` `{tipo: TRAMITES\|INVENTARIO\|TIEMPOS, formato: PDF\|XLSX, desde, hasta, estado, tipoTramite, categoriaId, estadoLegajo}` → el archivo como adjunto (`Content-Disposition`, `X-Reporte-Filas`). 422 `VALIDACION`, `RANGO_FECHAS_INVALIDO`, `FILTRO_INVALIDO`, `REPORTE_VACIO` | `REPORTE_GENERAR` (Director y, por seed, Personal de Archivo) |
| Reportes en lenguaje natural (CU-26) | `POST /reportes/interpretar` `{frase (≤300), categorias:[{id,nombre}]}` → `{tipo, formato, desde, hasta, estado, tipoTramite, categoriaId, estadoLegajo, explicacion}`: propuesta de filtros que el usuario revisa; no genera ni audita. Solo existe con `GROQ_API_KEY` (`GROQ_MODEL`, por defecto `openai/gpt-oss-120b`); `infrastructure/llm` usa llamada a función forzada y solo envía la frase, la fecha y los nombres de categorías (nunca datos de docentes). `InterpretarReporte` descarta lo que no sea un valor real. 422 `PETICION_NO_ENTENDIDA`, 503 `INTERPRETE_NO_DISPONIBLE` | `REPORTE_GENERAR` |
| Trámites (lectura) | `GET /tramites` (con `?pagina=N&tamano=10` devuelve `{datos,total,pagina,tamano}`; sin `pagina`, el arreglo completo), `GET /tramites/:id`, `GET /tramites/:id/documentos/:docId` | `TRAMITE_RADICAR`, `TRAMITE_ASIGNAR` o `TRAMITE_REVISAR` |
| Trámites (radicar y validar) | `POST /tramites`, `POST /tramites/:id/validar` | `TRAMITE_RADICAR` |
| Trámites (asignar) | `POST /tramites/:id/asignar`, `GET /analistas` | `TRAMITE_ASIGNAR` |
| Trámites (revisar) | `POST /tramites/:id/{revision,observaciones,dictamen,remitir-facultad,remitir-rectorado}`, `PUT /tramites/:id/periodo` `{anioDesde,anioHasta}`, `POST /tramites/:id/observaciones/:oid/resolver` | `TRAMITE_REVISAR` + ser el analista asignado |
| Trámites (regreso) | `POST /tramites/:id/regreso` `{resolucionNumero, resolucionFecha AAAA-MM-DD}` | `TRAMITE_REVISAR` (analista asignado) o `TRAMITE_RADICAR` (ventanilla) |
| Trámites (archivar) | `POST /tramites/:id/archivar` `{ubicacion}` (409 `CARPETAS_ABIERTAS`) | `TRAMITE_ARCHIVAR` (Personal de Archivo, Director) |
| Documentos (CU-15, CU-16) | `POST /tramites/:id/documentos/subida` `{nombre, tipoMime, tamano, requisitoId?}` → 201 `{documento, url}` (`DOCUMENTO_GESTIONAR`; 422 `TIPO_NO_PERMITIDO`/`ARCHIVO_GRANDE`, 403/409 `NO_ASIGNADO`); `POST …/documentos/:docId/confirmar` (también reintenta un documento en `ERROR`; 409 `SUBIDA_INCOMPLETA`); `GET …/documentos/:docId/texto` → `{paginas:[…]}`; `GET …/documentos/:docId/buscar?q=` → `{coincidencias:[{pagina,fragmento,rango}]}` (páginas que contienen todas las palabras, sin tildes; el fragmento marca lo hallado entre ⟦ y ⟧); `PATCH …/documentos/:docId/texto` `{paginas:[…]}` (debe traer exactamente tantas páginas como el documento; lo marca validado y reindexa); `PATCH /tramites/:id/documentos/:docId/estado` `{motivo}` (baja lógica; solo el analista asignado; 409 `DOCUMENTO_DE_BAJA`, 409 `TRAMITE_ARCHIVADO`). `GET …/documentos/:docId` entrega la versión comprimida si existe; un documento de baja solo lo descarga quien tiene `DOCUMENTO_RESTITUIR` (Director). |
| Portal público (30 peticiones/min por IP) | `GET /publico/tramites/:codigo`, `POST /publico/consulta-docente` `{codigo, cedula}` (código de docente **y** cédula del docente → lista mínima: código de seguimiento, tipo, categoría, estado y fecha; código desconocido, cédula que no coincide o docente sin cédula registrada responden lista vacía, para no revelar qué docentes existen; la cédula va en el cuerpo, no en la URL), `GET /publico/categorias`, `POST /publico/solicitudes` (multipart: `identificador`, `tipo`, `categoriaId` y un archivo `doc_<id>` por requisito) | sin sesión |

### 5 bis. Tablero (CU-25) y reportes (CU-26)

- **Solo lectura** sobre trámites, historial, legajos y movimientos de custodia: no hay tablas nuevas ni migración. Repositorio `persistence/repositorio_indicador.go` (puerto `indicador.Consulta`).
- **Resolución de un trámite** = su primer paso a Remitido en `tramite_historial` (o a Archivado en datos anteriores al flujo actual). El tiempo de resolución va de la radicación (`tramite.creado_en`) a ese momento.
- **Periodo del tablero:** acota los trámites *radicados* (distribución por estado, «en curso») y los *resueltos* (tiempo promedio, resueltos por analista). La carga vigente por analista y las alertas de mora son siempre las de hoy. Sin periodo, todo el historial. Las fechas son días de Bolivia y el «hasta» es inclusive.
- **Carga por analista:** los mismos usuarios que la bandeja de asignación (activos con `TRAMITE_REVISAR`, por eso aparece también el Director). «A cargo» usa los estados de `tramite.Estado.Activo()`; En facultad se muestra aparte.
- **Mora** (`indicador.EvaluarMora`): una carpeta prestada está en mora si pasó su fecha límite o lleva más de `DIAS_MAXIMOS_PRESTAMO` días fuera (opcional, 15 por defecto); los días de mora son el mayor de ambos. Se cuentan días calendario de Bolivia. Se calcula en cada consulta: al registrar la devolución la alerta desaparece sola (4a). El frontend vuelve a pedir el tablero cada 5 minutos (el CU pide al menos cada 15).
- **Reportes:** TRAMITES (filtro por radicación, estado, tipo y categoría), INVENTARIO (estado de la carpeta; incluye préstamo vigente y días de mora) y TIEMPOS (filtro por resolución, tipo y categoría; resumen con promedios por tipo y categoría). El caso de uso arma un `indicador.Reporte` independiente del formato y el puerto `ports.ExportadorReportes` lo compone: **PDF** A4 apaisado con el logo del DED y membrete «Universidad Autónoma Gabriel René Moreno — Departamento de Escalafón Docente» en cada página, fecha de emisión, responsable, filtros, resumen, tabla con cabecera repetida y «Página X de N»; **Excel** con hoja «Datos» (tabla con autofiltro, cabecera fija, fechas y números con su tipo) y hoja «Resumen». Un conjunto vacío responde 422 `REPORTE_VACIO` y no se audita; cada emisión registra `GENERAR_REPORTE` (tipo, formato, filas y filtros).

La consulta pública devuelve solo código, tipo, categoría, estado, etapas y observaciones; nunca datos personales. El portal identifica al docente por **código o cédula**.
`api/openapi.yaml` todavía no documenta los endpoints de docentes, categorías, trámites y portal.

Cada categoría incluye cuántos `docentes` y `tramites` la usan; solo se puede eliminar si ambos son 0 (si no, darla de baja).

Errores: `{ codigo, mensaje, detalles? }` con códigos de negocio estables (por ejemplo `CODIGO_DUPLICADO`, `MOTIVO_REQUERIDO`, `OBSERVACIONES_PENDIENTES`, `DOCUMENTOS_FALTANTES`).

## 6. Datos iniciales y carga del padrón

- `make seed` (idempotente): roles de sistema (Director, Analista, Personal de Archivo, Personal de Ventanilla), catálogo de 21 permisos (Ventanilla: `PERFIL_GESTIONAR` y `TRAMITE_RADICAR`; el Analista ya no radica), usuarios de prueba (`admin@ded.dev`, `analista@ded.dev`, `archivo@ded.dev`, `ventanilla@ded.dev`, contraseña `123123`; solo desarrollo), categorías del escalafón con documentos de ejemplo y, si hay `SEED_DIRECTOR_*`, el primer Director con enlace de activación.
- `make import-docentes file=../archivos/docsAyuda/datos.csv` (idempotente; `import-docentes-dry` solo informa y escribe `reporte-importacion.txt`, que no se versiona).
  - Resultado cargado: 1 640 docentes, 111 facultades, 16 categorías (9 cargos inactivos) y 1 641 legajos.
  - Separa "APELLIDOS NOMBRES" con una heurística (2 apellidos y el resto nombres, uniendo *de/del/la…*); unos 750 casos ambiguos quedan en el reporte.
  - Los códigos repetidos entre personas distintas (2188 y 3004) se importan sin código para corregirlos desde el panel.

## 7. Comandos

| Comando | Qué hace |
|---|---|
| `docker compose up -d redis mailpit` (desde la raíz) | Dependencias locales; obligatorio antes de `make dev`. |
| `make dev` | Levanta la API. |
| `make test` / `make arch` / `make lint` | Tests, regla de dependencias y `go vet`. En Windows con rutas con espacios, `make test` falla por `$(MAKE)`: ejecutar `go test ./...` y `make arch` por separado. |
| Un solo test | `go test ./internal/application/tramites -run NombreTest -count=1` |
| `make seed`, `make import-docentes`, `make import-docentes-dry` | Datos iniciales y padrón. |
| `make migrate-diff name=…`, `migrate-new`, `migrate-hash`, `migrate-apply`, `migrate-status` | Migraciones. |

Variables de entorno: `.env.example`. S3 usa los nombres estándar `AWS_REGION`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` más `S3_BUCKET` (los antiguos `S3_REGION`/`S3_ACCESS_KEY`/`S3_SECRET_KEY` siguen valiendo); el OCR usa `IA_SERVICE_URL` e `IA_SERVICE_TOKEN` (ver `ia-service/README.md`). El usuario IAM necesita `s3:PutObject`, `s3:GetObject` y `s3:ListBucket` sobre el bucket, y el bucket debe tener CORS que permita `PUT` desde el origen del frontend. Obligatorias: Obligatorias: `APP_URL_FRONTEND`, `DATABASE_URL`, `REDIS_URL`, `JWT_SECRETO` (≥ 32 bytes), `JWT_EMISOR`, `SESION_HORAS`, `MAX_INTENTOS_LOGIN`, `MINUTOS_BLOQUEO`, `MINUTOS_TOKEN_RECUPERACION`, `HORAS_TOKEN_ACTIVACION` y las `SMTP_*`. Opcionales: `S3_BUCKET`, `S3_REGION`, `S3_ACCESS_KEY`, `S3_SECRET_KEY`, `S3_ENDPOINT`, `DIAS_MAXIMOS_PRESTAMO` (umbral de mora del tablero, entero > 0, 15 por defecto), `SEED_DIRECTOR_*`.

## 8. Pruebas

Tests unitarios de los casos de uso con dobles en memoria (`auth`, `perfil`, `usuarios`, `roles`, `auditorias`, `docentes`, `tramites`, `legajos`), de seguridad, caché, middleware, router y repositorios (estos últimos usan `DATABASE_URL_TEST`, Postgres local en el puerto 5433), y del flujo real de custodia contra Postgres (`repositorio_legajo_test.go`), y del importador (`cmd/importdocentes`: separación de nombres, consolidación de códigos, lectura del CSV). CU-25/26: regla de mora y agrupación de tiempos (`domain/indicador`), casos de uso con dobles (`application/indicadores`), composición real del PDF y del Excel (`infrastructure/reportes`), rutas y permisos (`router_indicadores_test.go`) y consultas contra Postgres (`repositorio_indicador_test.go`, con `DATABASE_URL_TEST`).

## 9. Particularidades y pendientes

- GORM ignora `false` en columnas con `default:true` y devuelve el valor por defecto tras `Create`: para guardar `false` hay que hacer un `Update` aparte.
- Acciones auditadas: login (exitoso y fallido, con IP), restablecer/activar/cambiar contraseña, contacto del perfil, alta/edición/cambio de rol/baja/reactivación de usuarios, reenvío de activación, guardar/eliminar rol, docentes, categorías, cada paso de los trámites y la baja de documentos (`BAJA_DOCUMENTO`, con el documento y el motivo). Nunca se guardan contraseñas, tokens ni enlaces de activación.
- Búsquedas de docentes, trámites y legajos: con dígitos en el texto, el código y la cédula coinciden **exactos** (no trae códigos parecidos); sin dígitos, el nombre coincide por fragmento.
- La fecha límite de un préstamo se compara por fecha en UTC (Bolivia es UTC-4: de noche puede adelantarse un día).
- Notificaciones (CU-14): el envío corre en una goroutine tras responder la transición y un ticker reintenta los PENDIENTES cada minuto; si el proceso se cae a mitad de un envío, queda PENDIENTE y se reintenta al arrancar. Solo hay una instancia de API: con varias, un mismo aviso podría enviarse dos veces.
- `api/openapi.yaml` describe seguridad (CU-01 a CU-06), legajos/custodia, notificaciones, tablero y reportes; faltan docentes, categorías y trámites.
- Pendiente: describir docentes, categorías y trámites en `openapi.yaml` y consolidar las facultades con variantes de escritura.
- CU-16 tiene prueba de caso de uso con dobles (`tramites_test.go`) y de repositorio contra Postgres (`repositorio_tramite_test.go`).
- Las pruebas de repositorio solo corren con las variables del `.env` cargadas (`make test`); con `go test ./...` a secas, `DATABASE_URL_TEST` no está definida y se saltan en silencio. En bash, `source .env` no exporta `DATABASE_URL` porque la URL de Neon lleva un `&`.
- Un código numérico corto es un factor débil para identificar al docente en el portal; se mitiga con el límite de peticiones y la validación en ventanilla.

## 10. Cambios del plan de mejoras (noviembre 2026)

Detalle y justificación en `impruves.md` (raíz del proyecto).

- **Regularización sin categoría:** `tramite.categoria_id` es opcional (migración `20261114100000_tramite_categoria_opcional.sql`, con `chk_tramite_categoria`: obligatoria salvo en `REGULARIZACION`). En el dominio `Tramite.CategoriaID == 0` significa «sin categoría»; los listados y reportes usan `LEFT JOIN categoria`. `Radicar` ignora la categoría que llegue en una regularización.
- **Portal público (CU-13):** `POST /api/publico/solicitudes` recibe `codigo` + `cedula` (ya no `identificador`); ambos deben coincidir con el padrón (`docenteVerificado`, compartido con la consulta por docente) y cualquier fallo responde `DOCENTE_NO_VERIFICADO`. Un docente sin cédula registrada no puede radicar desde el portal. En regularización se envían de 1 a 5 archivos en el campo `adjunto` (sin requisito). Límites: 5 solicitudes por hora por IP (`LimitadorSolicitudes`) y 2 solicitudes pendientes de validación por docente (`SOLICITUDES_PENDIENTES`). El OCR de los adjuntos se lanza en `ValidarSolicitud`, no al radicar.
- **Carpetas (CU-20):** `tramiteId` es opcional en `POST /api/legajos/:id/solicitudes`: se puede pedir cualquier carpeta sin trámite abierto. Si se envía, se sigue exigiendo que sea del solicitante, esté abierto y sea del mismo docente (esa solicitud ligada habilita «Iniciar revisión»).
- **Sesión:** `Autenticar` falla cerrado: si Redis no deja consultar la lista de revocación responde 503 `SESION_NO_VERIFICABLE` (antes aceptaba el token). `GET /api/auth/sesion` devuelve los permisos del propio JWT (los que la API hace cumplir) y ya no los consulta en la BD.
- **Login y recuperación:** el bloqueo solo se informa (`CUENTA_BLOQUEADA`) a quien envía la contraseña correcta; con una incorrecta responde `CREDENCIALES_INVALIDAS`. Al vencer el bloqueo el contador vuelve a cero. La recuperación envía el correo en segundo plano y admite 3 solicitudes por hora por cuenta (`token.Repository.ContarDesde`).
- **Otros:** `helmet` y `compress` en Fiber; pool de conexiones acotado en `persistence.Conectar`; comodines de `LIKE` escapados (`escaparLike`); nombres escapados en todos los correos; `docker-compose.yml` publica los puertos solo en `127.0.0.1` y el `ia-service` corre sin root y verifica que las imágenes sean JPG/PNG reales.
- **Pendiente:** contraseña de Redis, borrar del bucket los objetos que exceden el tamaño al confirmar, límites de memoria del `ia-service` y rotar la llave de AWS (decisión del responsable: no se rota por ahora).

## Foto de perfil

`GET/PUT/DELETE /api/perfil/foto` (permiso `PERFIL_GESTIONAR`, solo la foto propia). `PUT` recibe multipart (campo `foto`, JPG o PNG, máx. 2 MB y 4096 px por lado); el tipo se valida por el **contenido** (`image.DecodeConfig`), no por la cabecera. El objeto vive en S3 en `usuarios/<id>/foto` (una clave fija que se sobrescribe) y la referencia en la tabla `usuario_foto` (migración `20261118120000_usuario_foto`, **tabla aparte a propósito**: si el código corre antes que la migración solo fallan estos endpoints, no el login). `GET` devuelve `{url, tipoMime, version, vigenciaSegundos}` con un enlace firmado de 10 min; `version` (unix de la última subida) sirve para invalidar cachés. `DELETE` borra la referencia y el objeto (si el IAM no permite `DeleteObject`, queda un objeto huérfano sin referencia). `ports.ObjectStorage` ganó `Eliminar`. Auditoría: reutiliza `ACTUALIZAR_CONTACTO` con `detalle.foto`.

