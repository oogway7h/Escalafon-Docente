// Acuerdo de usuario del Sistema de Escalafón Docente. Subir VERSION_TERMINOS aquí exige también subir
// VersionTerminos en backend/internal/domain/usuario/usuario.go: todos vuelven a aceptarlo al iniciar sesión.
export const VERSION_TERMINOS = 1;

export const DESARROLLADOR = "el equipo desarrollador del Sistema de Escalafón Docente";
export const FECHA_TERMINOS = "15 de noviembre de 2026";

export type SeccionTerminos = { titulo: string; parrafos: string[] };

export const SECCIONES_TERMINOS: SeccionTerminos[] = [
  {
    titulo: "1. Quiénes participan y para qué sirve este acuerdo",
    parrafos: [
      `Este acuerdo se celebra entre ${DESARROLLADOR} (en adelante, «el Desarrollador»), que diseñó, construyó y mantiene el software, y la persona que inicia sesión con una cuenta del personal (en adelante, «el Usuario»). El sistema es utilizado por el Departamento de Escalafón Docente de la Universidad Autónoma Gabriel René Moreno (UAGRM), que es la titular de la información que se gestiona en él.`,
      "Al marcar «Acepto» declaras que leíste este documento, que lo entendiste y que te comprometes a cumplirlo. Si no estás de acuerdo, no podrás usar el sistema y tu sesión se cerrará.",
    ],
  },
  {
    titulo: "2. Qué es el sistema",
    parrafos: [
      "El Sistema de Gestión de Flujos de Trabajo y Rastreo Documental del Departamento de Escalafón Docente permite radicar y dar seguimiento a los trámites de categorización y regularización docente, revisar expedientes y registrar observaciones, emitir dictámenes, cargar documentos con reconocimiento de texto (OCR), controlar el inventario y la custodia de los legajos físicos, consultar indicadores de gestión y generar reportes.",
      "El sistema es una herramienta de apoyo: no reemplaza el procedimiento ni la normativa universitaria (Estatuto Orgánico, Reglamento General del Profesor Universitario, Reglamento del Escalafón Docente y resoluciones vigentes). La documentación física sigue presentándose en dos ejemplares, acompañada de sus originales para la autenticación.",
    ],
  },
  {
    titulo: "3. Tu cuenta y tus credenciales",
    parrafos: [
      "La cuenta es personal e intransferible. Eres responsable de todo lo que se haga con ella. No compartas tu contraseña, no la anotes en lugares visibles y cámbiala si sospechas que otra persona la conoce.",
      "Tras tres intentos fallidos de inicio de sesión la cuenta se bloquea temporalmente. La recuperación de contraseña se hace con un enlace de un solo uso enviado a tu correo institucional, que vence a los treinta minutos.",
      "Cierra la sesión al terminar, sobre todo en equipos compartidos. Puedes actualizar tus datos de contacto y tu contraseña desde tu perfil; tu rol y tus permisos solo los modifica el Director.",
    ],
  },
  {
    titulo: "4. Uso permitido",
    parrafos: [
      "Puedes usar el sistema únicamente para las funciones que tu rol te autoriza y para los fines del Departamento de Escalafón Docente. Cada rol (Director, Analista, Personal de Ventanilla y Personal de Archivo) tiene permisos distintos; si ves una opción que no te corresponde, avisa al Director.",
      "Está prohibido: acceder con la cuenta de otra persona; intentar saltarte los controles de acceso; consultar expedientes sin una razón de trabajo; modificar, ocultar o falsear información de un trámite; cargar documentos que no correspondan al trámite o cuya autenticidad desconozcas; y usar el sistema para fines personales, comerciales o ajenos a la Universidad.",
    ],
  },
  {
    titulo: "5. Confidencialidad y datos personales",
    parrafos: [
      "El sistema contiene datos personales de docentes (nombres, cédula, código universitario, títulos, trayectoria y documentos de respaldo). Te comprometes a tratarlos con estricta reserva, a usarlos solo para tramitar el escalafón y a no copiarlos, enviarlos ni divulgarlos fuera de los canales institucionales, incluso después de dejar de usar el sistema.",
      "No tomes capturas ni descargues documentos para guardarlos en dispositivos personales o servicios externos. Las notificaciones por correo al docente solo incluyen su código de seguimiento y el enlace al portal, sin información sensible; no la amplíes por otros medios.",
      "El tratamiento de estos datos se rige por la normativa boliviana vigente y por los reglamentos de la UAGRM.",
    ],
  },
  {
    titulo: "6. Documentos y reconocimiento de texto (OCR)",
    parrafos: [
      "Al cargar un documento (PDF, JPG o PNG) el sistema calcula una huella (hash) para garantizar su integridad y lo procesa de forma automática para extraer el texto y comprimirlo. El texto extraído puede contener errores: antes de indexarlo debes revisarlo y corregirlo. Eres responsable de que lo que queda registrado corresponda al documento original.",
      "Dar de baja un documento cargado por error lo retira de las búsquedas, pero no lo borra: permanece en la auditoría y solo el Director puede restituirlo.",
    ],
  },
  {
    titulo: "7. Custodia de carpetas físicas",
    parrafos: [
      "Quien despacha o recibe una carpeta física registra en el sistema la fecha, el custodio y el estado material del legajo. Estos asientos forman una bitácora que no puede modificarse ni borrarse, y la responsabilidad por el estado de la carpeta recae en el último custodio registrado. Los plazos de devolución vencidos generan alertas visibles para el Director.",
    ],
  },
  {
    titulo: "8. Auditoría y monitoreo",
    parrafos: [
      "Aceptas que el sistema registre tus acciones: inicios de sesión, cambios en trámites, documentos, carpetas, cuentas, roles y reportes emitidos, junto con la fecha, la hora y la dirección IP desde la que las hiciste. Estos registros son de solo lectura, los consulta el Director y pueden usarse para investigar usos indebidos o respaldar decisiones administrativas.",
    ],
  },
  {
    titulo: "9. Asistente de reportes con inteligencia artificial",
    parrafos: [
      "En la sección de reportes puedes describir con tus palabras lo que necesitas. Esa frase se envía a un proveedor externo de modelos de lenguaje únicamente para convertirla en filtros del reporte; junto con ella viajan la fecha de hoy y los nombres de las categorías, nunca datos de docentes ni de trámites.",
      "Por eso no escribas en esa frase nombres, cédulas ni otros datos personales. El asistente solo propone filtros: tú los revisas y decides si generas el reporte. Cada reporte emitido queda registrado en la auditoría.",
    ],
  },
  {
    titulo: "10. Propiedad intelectual",
    parrafos: [
      "El software, su código, su diseño y su documentación pertenecen al Desarrollador, que concede a la UAGRM y a sus usuarios autorizados el derecho de usarlo para la gestión del escalafón docente. No puedes copiarlo, modificarlo, redistribuirlo ni intentar extraer su código. La información y los documentos cargados en el sistema pertenecen a la UAGRM.",
    ],
  },
  {
    titulo: "11. Disponibilidad y responsabilidad",
    parrafos: [
      "El Desarrollador procura que el sistema esté disponible y se mantenga seguro, pero no garantiza un funcionamiento ininterrumpido: puede haber mantenimientos, fallas del servicio de internet, de correo o de los proveedores externos. Las notificaciones por correo fallidas se reintentan y se registran, sin revertir la operación realizada.",
      "El Desarrollador no responde por decisiones tomadas con información cargada o corregida de forma errónea por los usuarios, ni por el uso del sistema contrario a este acuerdo. Las decisiones sobre categorías, puntajes y dictámenes corresponden a las autoridades y al personal competente del Departamento.",
    ],
  },
  {
    titulo: "12. Suspensión de la cuenta",
    parrafos: [
      "El incumplimiento de este acuerdo puede llevar a la suspensión o baja de tu cuenta por parte del Director, sin perjuicio de las responsabilidades administrativas o legales que correspondan. Dar de baja una cuenta conserva el historial de actividad.",
    ],
  },
  {
    titulo: "13. Cambios en este acuerdo",
    parrafos: [
      "El Desarrollador puede actualizar este documento cuando cambien las funciones del sistema o la normativa. Cada versión nueva se te volverá a presentar al iniciar sesión y deberás aceptarla para seguir usando el sistema. Registramos la versión que aceptaste y cuándo lo hiciste.",
    ],
  },
  {
    titulo: "14. Ley aplicable",
    parrafos: [
      "Este acuerdo se rige por las leyes del Estado Plurinacional de Bolivia. Cualquier controversia se resolverá primero por la vía institucional y, si no es posible, ante los tribunales competentes de la ciudad de Santa Cruz de la Sierra.",
    ],
  },
];
