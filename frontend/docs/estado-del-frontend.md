# Frontend — estado de lo implementado

Aplicación web del Sistema de Escalafón Docente (UAGRM). Stack: Next.js 16 (App Router), React 19, Tailwind CSS 4, zod, pnpm. Next.js 16 trae cambios incompatibles con versiones anteriores: leer `node_modules/next/dist/docs/` antes de tocar su API.

> Última actualización: 2026-10-25.

## 1. Casos de uso con pantalla

| CU | Pantalla | Estado |
|---|---|---|
| CU-01 | `/login` | ✅ |
| CU-02 | `/recuperar`, `/restablecer`, `/activar` | ✅ |
| CU-03 | `/panel/perfil` | ✅ |
| CU-04 | `/panel/usuarios` (listado, nuevo, editar, cambiar rol, baja/alta, reenviar activación) | ✅ |
| CU-05 | `/panel/roles` (listado, nuevo, editar, eliminar) | ✅ |
| CU-06 | `/panel/auditoria` (solo lectura: filtros por usuario, tipo de acción y rango de fechas, combinables; paginada) | ✅ |
| CU-07 | `/panel/docentes`, `/panel/categorias` (incluye eliminar categoría con confirmación en `/panel/categorias/[id]/eliminar`; solo si ningún docente ni trámite la usa) | ✅ |
| CU-08 | `/panel/tramites/radicar` | ✅ |
| CU-09 | `/panel/tramites/asignar` y asignación/reasignación en el detalle | ✅ |
| CU-10 | `/panel/tramites/revisar` (Mi bandeja) y revisión/observaciones en el detalle | ✅ |
| CU-11 | Dictamen, años a evaluar, remisión a la facultad, regreso y remisión al Rectorado en `/panel/tramites/[id]`; archivo en `/panel/tramites/archivar` (Personal de Archivo) | ✅ |
| CU-12 | `/seguimiento` (público): un solo campo de código. Más de 5 caracteres = código de seguimiento (el campo de cédula se desactiva); hasta 5 = código de docente y **exige la cédula** (acción POST: la cédula no viaja en la URL). Con un trámite lo muestra directo; con varios, una lista para elegir. Un error genérico no revela si falló el código o la cédula. Diseño: buscador en una franja oscura, la tarjeta de estado (código, estado, `FasesTramite`, tipo/categoría/radicado y resultado) y debajo las observaciones y el historial de etapas (la más reciente primero) en dos columnas, con una nota de privacidad | ✅ |
| CU-13 | `/solicitud` (público) | ✅ (los adjuntos dependen de S3 en el backend) |
| CU-19 | `/panel/legajos` (inventario con búsqueda exacta por código, paginado; `/nuevo`, `/[id]/editar`) | ✅ |
| CU-20 | `/panel/carpetas` (mis solicitudes, con el trámite) y `/panel/carpetas/nueva` (elegir el trámite y solicitar una carpeta de su docente) | ✅ |
| CU-21, CU-22, CU-23 | **`/panel/custodia` «Solicitudes y custodia»**: una sola pantalla para Archivo con pestañas Por atender / Prestadas / Vencidas (con conteos y buscador), lista a la izquierda y panel de acción a la derecha. En «Por atender», **«Aprobar y entregar carpeta»** hace en un paso la aprobación y la salida (fecha límite con atajos +3/+7/+15 días, estado material, observación); también «Solo aprobar» y «Rechazar con motivo». En «Prestadas»/«Vencidas», «Registrar devolución» (la observación es obligatoria si no vuelve en buen estado). `/panel/legajos/solicitudes` redirige aquí | ✅ (la devolución no lleva foto de evidencia) |
| CU-24 | `/panel/custodia/bitacora` (filtros por texto, tipo y vencidas; paginada) | ✅ |
| CU-14 | `/panel/notificaciones` (avisos enviados al docente, filtro por estado del envío; solo Director) | ✅ |
| CU-16 | Baja de adjuntos del portal en `/panel/tramites/[id]` (`BajaDocumentoForm`) y de escaneos en `/panel/documentos/[id]` (`BajaEscaneoForm`) | ✅ |
| CU-15 | Archivo sube escaneos por docente y tipo en `/panel/documentos` → `/panel/documentos/[id]` (OCR con sondeo y reintento); Archivo, analista y Director buscan por contenido en `/panel/documentos/buscar` (elegir docente → palabras → página exacta en el visor). Los adjuntos del portal se ven en `/panel/tramites/[id]` | ✅ |
| CU-25 | `/panel/tablero` (indicadores del periodo, carga por analista y alertas de mora; se refresca solo cada 5 minutos) | ✅ |
| CU-26 | `/panel/reportes` (tipo, filtros y formato PDF/Excel; descarga por `/panel/reportes/descargar`) | ✅ |
| CU-17 | `/panel/documentos/bajas` (Director, `DOCUMENTO_RESTITUIR`): bandeja paginada de documentos dados de baja; al elegir uno (`?doc=`) se ve su detalle y se restituye con motivo (`restituirDocumentoAction`, que redirige a `?restituido=1` para mostrar el aviso) | ✅ |
| CU-18 | `/panel/documentos/buscar` en dos pasos: 1) docente por nombre o código con sugerencias mientras se escribe; 2) contenido de sus documentos con sugerencias de página y fragmento, resultados completos y visor en la página exacta | ✅ |

## 2. Estructura

| Carpeta | Contenido |
|---|---|
| `app/page.tsx` | Landing pública con enlaces a login, seguimiento y nueva solicitud. |
| `app/(auth)/` | Login, recuperar, restablecer y activar. |
| `app/(publico)/` | Portal sin sesión: `seguimiento` y `solicitud`, con `PublicoShell`. |
| `app/panel/` | Área autenticada. `layout.tsx` carga la sesión y la pasa a `PanelShell`. |
| `components/common/` | UI reutilizable: `Boton` (`claseBoton`), `BotonEnviar`, `Campo` (`claseControl`), `Tarjeta`, `Tabla`, `Insignia`, `Alerta`, `EncabezadoPagina`, `AccionesFormulario`, `Paginacion`, `Loader`, `BotonTema`. |
| `components/layout/` | `PanelShell`, `Sidebar`, `Header`, `AuthDividido`, `AuthTarjeta`, `PublicoShell`, `IconosPanel`. |
| `components/tramites/` | `FasesTramite` (recorrido del trámite), `EstadoInsignia`, `TablaTramites`. |
| `lib/` | `api.ts` (cliente de la API, sesión, usuarios y roles), `api-tramites.ts` (docentes, categorías, trámites y portal), `api-auditoria.ts` (consulta de auditoría), `auditoria.ts` (etiquetas de acciones, fecha y hora de Bolivia, resumen del detalle), `esquemas.ts` (zod), `formularios.ts` (estado de formulario y mensajes de error de la API), `navegacion.ts` (menú por permisos; cada ítem lleva además una `descripcion` que usa la portada del panel), `tramites.ts` (etiquetas y fechas), `sesion.ts`, `tema.ts`. |
| `proxy.ts` | Redirige a `/login` si falta la cookie de sesión en `/panel/*`. La validación real la hace la API en cada petición. |

## 3. Cómo funciona

- **Datos:** los Server Components llaman a la API desde el servidor con `api()` (`lib/api.ts`), que agrega el token de la cookie `httpOnly`. Un 401 redirige a `/login` (la cookie solo se borra si la llamada viene de una acción de servidor: Next no permite modificar cookies durante el render de una página); un 403 responde 404. El navegador nunca habla con la API.
- **Formularios:** `actions.ts` por segmento con `useActionState` y validación zod; los errores de la API se traducen con `MENSAJES_API` (`lib/formularios.ts`). Las acciones de un solo botón (baja, validar, iniciar revisión, resolver) devuelven el error por la URL (`?error=CODIGO`).
- **Portada del panel** (`/panel`): tarjetas clicables con lo que el rol puede hacer, una por pantalla, con la descripción de `NAVEGACION`. Se agrupa por sección solo cuando las secciones tienen cuerpo (de media, dos pantallas o más): el Director ve 21 pantallas en 6 secciones y se agrupan, mientras que Analista (5), Personal de Archivo (9) y Ventanilla (3) las ven en una rejilla plana, porque un encabezado sobre una única tarjeta es ruido. Los ítems con `disponible: false` salen apagados, con borde discontinuo y «Pronto», y no son enlaces. «Mi bandeja» y «Solicitudes y custodia» llevan además una insignia con los pendientes (expedientes asignados sin remitir y solicitudes de carpeta sin atender); con cero no se pinta nada, y el número lleva `title` y `aria-label` porque una cifra suelta no dice de qué es. Cada insignia va en su propio `<Suspense>`: `/panel` no tiene `loading.tsx` propio, así que esperar la consulta dejaría a la vista la pantalla anterior. Medido: el armazón con las tarjetas sale a 728 ms y la insignia llega ~280 ms después. Si la consulta falla, la tarjeta se pinta sin insignia.
- **Permisos:** `exigirPermiso("X")` en cada página; el menú (`lib/navegacion.ts`) muestra solo lo permitido y marca como "próximamente" lo que aún no existe (`disponible: false`).
- **Portal público:** `apiPublica()` no usa sesión ni redirecciones. `/seguimiento` consulta por código con un formulario GET (funciona sin JavaScript). `/solicitud` identifica al docente por **código o cédula**, pide un archivo por cada documento que exige la categoría y, al terminar, redirige a `/seguimiento?codigo=…&nueva=1`.
- **Trámite en el panel** (`/panel/tramites/[id]`): datos, línea de fases, documentación exigida (con descarga firmada de los adjuntos), observaciones, dictamen, historial y las acciones que correspondan según el permiso, el estado y si el usuario es el analista asignado.
- **Baja de documentos (CU-16):** con `DOCUMENTO_GESTIONAR` y el expediente sin archivar aparece la tarjeta «Dar de baja un documento» (elegir el adjunto vigente y escribir el motivo, mínimo 5 caracteres). Si el trámite ya tiene dictamen emitido, la tarjeta avisa que hay que justificar la baja. El documento dado de baja deja de cubrir su requisito y se muestra tachado en la lista de documentación exigida, con su motivo y fecha, porque solo el Director podrá restituirlo (CU-17).
- **Docentes:** listado paginado en el servidor (10 por página, `?pagina=N`) con búsqueda por código, cédula o nombre; el **código es la primera columna**, seguido de nombre, categoría y facultad, gavetas del legajo y estado. Formulario con código obligatorio y cédula opcional.
- **Carga del listado de docentes:** `components/common/CargaListado.tsx` (`ListadoConCarga`, `ZonaCarga`, `BuscadorCarga`, `LinkCarga`). Al buscar o cambiar de página se navega con `useTransition`: `ZonaCarga` sustituye la tabla por su silueta (mismas columnas y filas, medidas del DOM) hasta que llega el resultado. Donde no envuelve una tabla, conserva el contenido difuminado con el `Loader` compacto encima, centrado con `sticky`. `Paginacion` usa `LinkCarga`, y los botones de acción (dar de baja, reactivar) usan `FormularioCarga`, que mantiene la carga hasta que la acción de servidor termina y llega la vista actualizada; los iconos de editar usan `LinkCarga`. Docentes y Categorías los usan completos (búsqueda, paginación, editar, baja/alta); Usuarios y Roles los usan en el botón "Nuevo" y sobre la tabla. Cualquier listado puede reutilizarlos envolviéndose en `ListadoConCarga`. Sin JavaScript sigue funcionando como enlaces y formulario GET.
- **Auditoría (CU-06):** `/panel/auditoria` exige `AUDITORIA_CONSULTAR`. Los filtros viajan por URL con los mismos nombres que la API (`usuarioId`, `accion`, `desde`, `hasta`, `pagina`) mediante `FiltrosCarga` (`CargaListado.tsx`), que omite los campos vacíos. Las opciones de los filtros salen de `GET /api/auditoria/filtros`; la validación (fechas, rango, acción) la hace la API y la pantalla muestra su mensaje. Paginada en el servidor como Docentes (10 por página, `?pagina=N`, con `Paginacion`, que conserva los filtros). Columnas: fecha y hora, usuario, acción, detalle y dirección IP (la entidad afectada no se muestra, aunque la API la devuelve). Tabla sin acciones: la bitácora es de solo lectura. En pantallas angostas la tabla se desplaza horizontalmente.
- **Legajos y custodia:** `lib/api-legajos.ts` (API, `server-only`) y `lib/legajos.ts` (etiquetas, sin `server-only`). Las acciones por fila (atender, solicitar, entregar, recibir) son formularios con server action que devuelven el error por la URL; el botón pulsado manda `aprobar=true|false` por su `name/value`. La barra lateral resalta solo el ítem más específico (`/panel/custodia/bitacora` no enciende también `/panel/custodia`).
- **Trámites:** el listado general pagina en el servidor (10 por página, `?pagina=N`, conserva `estado` y `buscar`); las bandejas de asignar y revisar siguen usando el arreglo completo.
- **Buscadores de usuarios y trámites:** usan `next/form`, que navega en el cliente y conserva el layout. Un `<form method="get">` plano recarga todo el documento y hace aparecer el loader de pantalla completa.

## 4. Diseño

- Tokens de color y modo claro/oscuro en `app/globals.css` (el script de `lib/tema.ts` aplica la clase antes de hidratar, sin parpadeo).
- Materiales translúcidos (`.material`: fondo semitransparente con `backdrop-filter`), animación de aparición (`animate-materializar`), línea de progreso animada (`progreso-fase`) y tipografía con tracking negativo en titulares.
- Respeta `prefers-reduced-motion`, `prefers-reduced-transparency` y `prefers-contrast`.
- **Esqueletos de listado**: las rutas del panel que muestran una tabla tienen su propio
  `loading.tsx` con `ListadoEsqueleto` (título real + silueta de recuento, filtros y tabla) en vez
  del `Loader` genérico. `TablaEsqueleto` reproduce la geometría de `Tabla` —fila de 65 px por la
  columna de dos líneas, encabezado de 37 px— para que al llegar los datos no salte la
  maquetación; medido contra las tablas reales (auditoría 64.6 px, categorías 66.8 px). La barra
  late con `.esqueleto` (`globals.css`), con desfase por fila y sin animación bajo
  `prefers-reduced-motion`.
  También lo usa `ZonaCarga` al buscar o paginar: mide del DOM las columnas y filas de la tabla
  que el usuario acaba de ver y dibuja esa silueta exacta, así que la página no cambia de altura
  (medido: 685 px de silueta frente a 683 px de tabla). Donde `ZonaCarga` no envuelve una tabla
  (tablero, custodia, buscador de documentos) se conserva el contenido anterior difuminado con el
  `Loader` encima, porque ahí no hay forma que replicar.
- **Una sola frontera de carga por ruta.** Cada ruta del panel tiene su propio `loading.tsx`: las
  de listado con el esqueleto, las demás con el `Loader`. Ya **no** existe `app/panel/loading.tsx`.
  Al ser ancestro de todas, creaba un segundo `<Suspense>` por encima del de cada listado y en el
  streaming React emitía primero el fallback de fuera: se veía el `Loader` y justo después el
  esqueleto. La única ruta sin frontera propia es `/panel` (la portada), porque un `loading.tsx`
  ahí volvería a cubrir a todos sus hijos; al no tenerla, Next mantiene visible la página anterior
  mientras carga, que para una sola tarjeta es preferible.
- El `Loader` tiene fondo translúcido con desenfoque. Los `loading.tsx` que no son listados
  re-exportan ese mismo `Loader` (requisito de convención de Next); no hay uno en la raíz a
  propósito, para que al recargar el panel no aparezca un loader de pantalla completa que oculte
  el layout; el login y la solicitud pública lo muestran a pantalla completa mientras se envía el
  formulario. El `Loader mini` que aparece dentro del botón «Buscar» al paginar es el estado
  ocupado del propio botón y es intencionado.
- El botón de tema va arriba a la derecha en las pantallas con `AuthDividido`.
- **Contraste**: todos los pares texto/fondo del sistema cumplen WCAG 2.1 AA (4.5:1 normal, 3:1 grande)
  en los dos modos. `node scripts/contraste.mjs` lo comprueba leyendo los tokens reales de
  `globals.css` y falla si alguien baja un valor; al añadir una combinación nueva (`bg-X text-Y`)
  hay que añadir el par a esa lista.
- `--color-accent` es el **relleno** de marca (lleva texto blanco encima) y `--color-accent-ink` es el
  acento **como texto**. En modo oscuro no pueden ser el mismo valor: el relleno necesita ser oscuro
  para que el blanco encima se lea, y el texto necesita ser claro. Para texto, borde de énfasis o
  icono, usar `text-accent-ink`; para fondos, `bg-accent`.
- `app/globals.css` tematiza también las superficies que pinta el navegador: anillo de foco único
  (`:focus-visible`, blanco sobre las superficies navy vía `.sobre-oscuro` y `aside`), `::selection`,
  `caret-color`, `accent-color` (casillas y radios), barra de desplazamiento y `tabular-nums` global
  para que las columnas de cifras no bailen al paginar.
- `Tabla` envuelve la `<table>` en un contenedor `overflow-x-auto` con `min-w-[640px]`: en pantallas
  estrechas la tabla se desplaza dentro de su tarjeta en vez de empujar scroll horizontal a la página.

## 5. Comandos y configuración

| Comando | Qué hace |
|---|---|
| `pnpm dev` | Servidor de desarrollo (puerto 3000). |
| `pnpm build` / `pnpm start` | Compilación y servidor de producción. |
| `pnpm lint` | ESLint (incluye reglas de pureza de React 19). |

- Variable: `API_URL_INTERNA` (URL de la API; solo la usa el servidor de Next).
- `next.config.ts` sube el límite de cuerpo de las server actions a 40 MB por los adjuntos del portal.
- Si ya corre un `next dev` en el proyecto, Next 16 no permite abrir otro: para probar en paralelo usar `next build` + `next start -p <puerto>`.
- No hay tests automatizados en el frontend.

## 6. Pendiente

- Paginar las bandejas de asignar y revisar (usan el arreglo completo).
- Pruebas automatizadas de los flujos principales (login, radicar, portal).

## Expediente digital y búsqueda por contenido (CU-15)

El Personal de Archivo escanea la carpeta física y sube cada documento al **docente** (no a un trámite): `/panel/documentos` (permiso `DOCUMENTO_ESCANEAR`) lista el padrón con `SelectorDocente` y `/panel/documentos/[id]` muestra «Subir escaneo» (tipo: Currículum, Regularización, Adjunto, Asociado, Titular A/B/C, Emérito + archivo; PDF hasta 400 MB, JPG/PNG hasta 20 MB) y los documentos agrupados por tipo con su estado de OCR (se refresca solo cada 6 s), reintento y baja con motivo. `SubirEscaneo.tsx` pide la URL firmada con una server action (`prepararSubidaExpedienteAction`), sube **directo a S3 con XHR** (barra de progreso; hace falta CORS con `PUT` en el bucket) y confirma con `confirmarSubidaExpedienteAction`, que lanza el OCR.

Archivo, analistas y Director (`DOCUMENTO_BUSCAR`) usan `/panel/documentos/buscar` (CU-18), en dos pasos marcados en pantalla. **Paso 1:** se busca al docente por nombre o código; `BuscadorSugerencias` (`components/common`) muestra debajo, como la barra de un navegador, hasta 8 docentes mientras se escribe (`buscar/sugerencias`, que llama a `GET /api/docentes/sugerencias`: código por prefijo, nombre sin tildes); Enter sin elegir lista todos los que coinciden (`SelectorDocente`). **Paso 2:** dentro del expediente, la misma barra sugiere las páginas halladas con su fragmento (`buscar/coincidencias`) y Enter lista todas las coincidencias (tipo de documento, página, fragmento resaltado); al pulsar una se abre el visor en esa página (`?doc=&pagina=&q=`). La barra navega con la transición de `ListadoConCarga` (`useCargaListado`), así el botón Buscar muestra el Loader mientras llega el resultado; ↑/↓, Enter y Esc funcionan como en un navegador y las respuestas se recuerdan en la sesión de la página. El texto es buscable apenas termina el OCR. En `/panel/tramites/[id]` hay un acceso «Buscar en el expediente» y siguen viéndose los adjuntos del portal (con su OCR y la corrección opcional por página).

## Listados y visor de documentos

- Las tablas de trámites tienen la **fila completa clicable** (`Tabla` con `href`, `FilaClicable` en `components/common/EnlaceCarga.tsx`) y el botón de «Revisar solicitud» / «Asignar» muestra el Loader mientras navega (`EnlaceBoton`).
- **Visor del sistema** (`app/panel/tramites/[id]/VisorDocumento.tsx` + `components/documentos/VisorPdf.tsx`): lector propio con **pdf.js** (`pdfjs-dist@4.10.38`, worker en `public/pdf.worker.min.mjs`: copiarlo de `node_modules/pdfjs-dist/build` si se actualiza la dependencia) en lugar del visor del navegador. Tiene cabecera con volver/Descargar, barra flotante (página anterior/siguiente y salto a una página, zoom, navegación entre coincidencias), miniaturas de la página actual, las vecinas y las que tienen coincidencias, resaltado de la búsqueda sobre la capa de texto del OCR y un panel con coincidencias y detalles. El archivo llega por un route handler propio (`…/archivo/[docId]`, `lib/archivo.ts`) que pasa los rangos de S3, así el bucket no necesita CORS de lectura. Se usa en `/panel/tramites/[id]?doc=`, `/panel/documentos/[id]?doc=` y `/panel/documentos/buscar`. Las páginas van **una debajo de otra y se recorren con el scroll** (la página actual es la que ocupa el tercio superior de la vista); solo se dibujan las que están a menos de una pantalla, el resto son recuadros del tamaño correcto. **Rendimiento:** pdf.js pide tramos de 1 MB; `lib/archivo.ts` recuerda la URL firmada por sesión y documento 10 min (los tramos no pasan por la API), descarga el objeto completo a una caché en disco (`$TMPDIR/escalafon-documentos`, límite 2 GB, se borran los más antiguos) y sirve los rangos desde ahí; la ruta lleva `?v=<tamaño comprimido>` y responde `private, max-age=86400`, así el navegador también lo cachea.
- **Buscador** (`/panel/documentos/buscar`): ver «Expediente digital y búsqueda por contenido». Los chips por tipo de documento (`?tipo=`) filtran los resultados del paso 2 y, al abrir uno (`?doc=&pagina=`), el visor reemplaza a la lista.

## Flujo del trámite y línea de tiempo

- `FasesTramite` dibuja el recorrido real: Ingresado · Asignado · En revisión · Dictaminado · **tramo punteado «Sin seguimiento» (Fuera de Escalafón: Facultad · Tribunal · Decano · Consejo, con fecha de salida, de regreso y días fuera)** · Revisión final · Remitido · Archivado. Observado se marca sobre En revisión. Se usa en el detalle y en `/seguimiento`, donde el docente lee que en esa etapa no se puede rastrear su avance.
- Detalle del trámite según estado y rol: sin carpeta solicitada se ofrece «Solicitar carpeta» en vez de iniciar la revisión; en revisión aparece «Años a evaluar» (obligatorio para dictaminar); con dictamen, «Remitir a la facultad»; en la facultad, «Registrar regreso» (analista asignado o ventanilla: resolución y fecha); en revisión final, «Remitir al Rectorado»; remitido, «Archivar carpeta» solo para `TRAMITE_ARCHIVAR`, y bloqueado mientras haya carpetas sin devolver.
- Personal de Archivo: «Por archivar» (`/panel/tramites/archivar`) lista los trámites remitidos.

## Requisitos por categoría en el portal

- La pantalla de inicio (`app/page.tsx`) tiene, antes de «Alcance del sistema», la sección «Requisitos por categoría»: una tarjeta grande de dos caras por categoría vigente (`CategoriasRequisitos.tsx`). Con ratón gira al pasar el cursor; al hacer clic o tocar queda **bloqueada** girada (en móvil el toque es el único gesto) y otro clic la devuelve. Frente: nombre, antigüedad mínima y n.º de documentos; reverso: los documentos exigidos y «Iniciar mi solicitud». Respeta `prefers-reduced-motion`. Lee `GET /publico/categorias`; si la API no responde, la sección se omite.
- En `/solicitud`, al elegir la categoría se listan sus documentos (con la antigüedad mínima) como campos de carga.
- Ambos salen del catálogo que el Director edita en `/panel/categorias` (CU-07): cambiar allí un requisito lo cambia en el portal.


- **Fade al hacer scroll (inicio):** los bloques de `app/page.tsx` llevan las clases `fade-scroll`, `fade-izq`, `fade-der` y `fade-salida` (definidas en `app/globals.css`, animación CSS ligada al scroll con `animation-timeline: view()`): aparecen al entrar por abajo y se desvanecen al salir por arriba. Sin soporte del navegador o con «reducir movimiento» se ven sin animación. Las tarjetas de categorías no se tocan: solo se envuelven.

## Tablero de indicadores y reportes (CU-25, CU-26)

- **Tablero** (`/panel/tablero`, permiso `TABLERO_VER`): `lib/api-indicadores.ts` (`obtenerTablero`, mismo patrón de resultado `{ok}` que la auditoría) y `lib/indicadores.ts` (etiquetas, formato de días, atajos de periodo; sin `server-only`). Un solo filtro de periodo arriba (`desde`/`hasta` por URL con `FiltrosCarga`, más los atajos «Este mes», «Últimos 3 meses», «Este año» y «Todo el historial»). Con carpetas en mora abre con un aviso que lleva a la tabla de alertas. Muestra cuatro indicadores (radicados, en curso, tiempo promedio de resolución, carpetas en mora), trámites por estado y carga por analista lado a lado, la tabla de **carpetas retenidas fuera de plazo** (legajo, docente y trámite, custodio, salida y límite, días de mora y motivo; «Registrar devolución» lleva a `/panel/custodia?tab=vencidas` si el usuario tiene `CUSTODIA_REGISTRAR`) y el tiempo de resolución por tipo y categoría. `ActualizacionAutomatica.tsx` llama a `router.refresh()` cada 5 minutos.
- **Barras:** una sola serie por gráfico, color `--color-grafico` (token nuevo en `globals.css`, `#2c5aa0` claro / `#6a96d8` oscuro, validado contra la superficie de cada modo), extremo redondeado de 4px y valor siempre escrito al lado (la tabla es la vista accesible; el `title` repite el dato al pasar el cursor). El rojo queda reservado para la mora.
- **Reportes** (`/panel/reportes`, permiso `REPORTE_GENERAR`): `ReporteForm.tsx` elige el tipo (Trámites, Inventario de legajos, Tiempos de atención), muestra solo los filtros de ese tipo y el formato (PDF o Excel). Con JavaScript envía el formulario por `fetch` a `/panel/reportes/descargar`, descarga el archivo sin salir de la página y explica los errores en el mismo formulario (`REPORTE_VACIO` = 3a del CU); sin JavaScript el formulario se envía igual y el route handler redirige con `?error=`. Las categorías salen del catálogo público (`listarCategoriasPublicas`) para que el Personal de Archivo, que no tiene permisos de catálogo, también pueda filtrar.
- La auditoría muestra la emisión como «Generar reporte» (`GENERAR_REPORTE`).
- Si `next build` falla con `TS1128` en `.next/dev/types/routes.d.ts`, ese archivo generado quedó truncado por un `next dev` anterior: se regenera al volver a levantar `next dev` (o borrando `.next/dev/types`).

## Cambios del plan de mejoras (noviembre 2026)

Detalle en `impruves.md` (raíz del proyecto).

- **`components/common/ZonaArchivo.tsx`:** campo de archivo con arrastrar y soltar. Es un `<label>` sobre un `<input type="file">` nativo (clic y teclado abren el selector, el formulario lo envía como cualquier campo). Props: `name`, `accept`, `maxMB` (número o función por tipo), `maxArchivos`, `required`, `disabled`. Lo usan `/solicitud` y `SubirEscaneo.tsx`.
- **Regularización sin categoría:** en `RadicarForm.tsx` y `SolicitudForm.tsx` el campo de categoría solo existe con «Categorización»; `radicarEsquema` la exige solo en ese caso. El portal pide código + cédula y, en regularización, de 1 a 5 adjuntos libres. `tipoYCategoria()` (`lib/tramites.ts`) arma «Tipo · Categoría» y omite la categoría si no hay.
- **Límites de tamaño:** `LIMITE_MB` y `MAX_ADJUNTOS_LIBRES` en `lib/tramites.ts` son la única fuente para textos y validaciones.
- **Carpetas:** `SolicitarCarpetaModal.tsx` permite pedir la carpeta de cualquier docente; si el analista tiene trámites de ese docente se ofrecen en un `<select>` con la opción «Sin trámite».
- **Seguridad:** `cabecerasCliente()` (`lib/api.ts`) reenvía `X-Forwarded-For` a la API en `api()`, `apiPublica()` y el login (sin eso el limitador y la auditoría veían una sola IP). `next.config.ts` envía CSP, `X-Frame-Options`, `nosniff`, `Referrer-Policy`, `Permissions-Policy` y HSTS en producción; `CSP_CONNECT_EXTRA` agrega un origen para subir documentos (MinIO local). El visor entrega documentos con `max-age=600`.
- **Rendimiento:** `listarCategoriasPublicas()` se guarda 5 minutos en la caché de datos (etiqueta `categorias-publicas`, invalidada desde `app/panel/categorias/actions.ts`); llama a `connection()` para que las páginas sigan siendo dinámicas (la API no existe durante el build). `AutoRefrescar` no refresca con la pestaña oculta.
- **Diseño minimalista (panel):** componentes base (`Campo`, `Tarjeta`, `Loader`, `Boton`) con bordes de 1 px, sin sombras decorativas y con una escala reducida de tamaños y radios en el panel. `Loader` es un anillo que gira (`.animate-girar`). El filtro de estado de `/panel/tramites` es un `<select>` dentro del buscador. **La parte pública conserva su aspecto original** (portada, seguimiento, solicitud e inicio de sesión): tipografía Manrope para los títulos (`font-display`) con Source Sans 3 para el texto, tarjetas de vidrio translúcido (`.material`, con desenfoque) y la barra superior translúcida, degradados en la portada y en el panel de acceso, los fade al hacer scroll (`.fade-*`) y la aparición de los resultados del seguimiento (`.animate-materializar`).
- **Visor y búsqueda dentro del documento:** `components/documentos/VisorConBusqueda.tsx` (cliente) envuelve `VisorPdf` y el panel lateral. Buscar y saltar entre coincidencias ocurre en el navegador: la búsqueda pide solo las páginas a `/panel/tramites/[id]/documentos/[docId]/buscar` y el salto mueve el visor (la URL se refleja con `history.replaceState`), sin volver a renderizar la página en el servidor. En `VisorPdf`: el salto a una página lejana es instantáneo (solo se anima a la vecina), el resaltado se calcula aparte del dibujo, el riel muestra solo las miniaturas vecinas (con retardo, para no ocupar el único trabajador de pdf.js) y una barra de progreso muestra la descarga del documento (`loadingTask.onProgress`).

## Acuerdo de usuario y requisitos oficiales (noviembre 2026)

- **Acuerdo de usuario:** `app/panel/layout.tsx` muestra `components/common/AcuerdoUsuario.tsx` (texto en `lib/terminos.ts`) en lugar del panel mientras `sesion.terminosAceptados` sea `false`; aceptar llama a `POST /api/auth/terminos` (`app/panel/terminos.ts`) y rechazar cierra la sesión. Al cambiar el texto hay que subir `VERSION_TERMINOS` aquí y `VersionTerminos` en el backend. El nombre del desarrollador y la fecha se cambian en `lib/terminos.ts` (`DESARROLLADOR`, `FECHA_TERMINOS`).
- **Tarjetas del portal:** muestran el puntaje de categoría (art. 18 del Reglamento; `puntajeDesde`/`puntajeHasta` de la categoría, vacíos = sin puntaje y no sale la tarjeta) en vez de los años de antigüedad, más una tarjeta de Regularización y la tabla de puntos por función de autoridad (`lib/requisitos.ts`). Los documentos de la categorización son los mismos para todas las categorías y los carga `make seed`; el Director los puede ajustar en `/panel/categorias`.


## Foto de perfil

`/panel/perfil` rediseñada (foto, datos de contacto y contraseña). `app/panel/perfil/FotoPerfil.tsx` (arrastrar y soltar, vista previa con `blob:`, validación JPG/PNG ≤ 2 MB) usa las acciones `subirFoto`/`quitarFoto`. `GET /panel/perfil/foto` (route handler) entrega la imagen desde el mismo origen leyendo el enlace firmado de la API, porque la CSP solo permite `img-src 'self' data: blob:`; no pasa por `lib/archivo.ts` porque su caché en disco se indexa por la ruta del objeto y la foto reemplaza siempre la misma clave. `components/common/Avatar.tsx` muestra la foto o las iniciales y se usa en el encabezado del panel; `obtenerVersionFoto()` (`lib/api.ts`) **nunca lanza ni redirige**: si la API o la migración fallan, el panel abre con iniciales.

