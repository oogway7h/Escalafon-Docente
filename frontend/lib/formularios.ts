import type { ZodError } from "zod";
import { LIMITE_MB } from "./tramites";

export type EstadoFormulario = {
  error?: string;
  mensaje?: string;
  camposInvalidos?: Record<string, string>;
};

export function camposDeZod(error: ZodError): EstadoFormulario {
  const campos: Record<string, string> = {};
  for (const issue of error.issues) {
    const campo = issue.path[0];
    if (typeof campo === "string" && !campos[campo]) campos[campo] = issue.message;
  }
  return { camposInvalidos: campos };
}

// Traduce los códigos de negocio de la API (sección 7 de Inicializacion.md) a mensajes en español.
export const MENSAJES_API: Record<string, string> = {
  CREDENCIALES_INVALIDAS: "Correo o contraseña incorrectos.",
  CUENTA_BLOQUEADA: "Tu cuenta está bloqueada temporalmente por intentos fallidos.",
  TOKEN_INVALIDO: "El enlace es inválido o ya expiró.",
  POLITICA_CONTRASENA: "La contraseña no cumple la política mínima (8 caracteres, una mayúscula y un dígito).",
  CORREO_DUPLICADO: "Ese correo ya está registrado.",
  ROL_INVALIDO: "El rol seleccionado no existe o está inactivo.",
  AUTOCAMBIO_ROL: "No puedes cambiar tu propio rol.",
  AUTODESACTIVACION: "No puedes desactivar tu propia cuenta.",
  CUENTA_YA_ACTIVADA: "Esa cuenta ya fue activada.",
  CONTRASENA_ACTUAL_INCORRECTA: "La contraseña actual es incorrecta.",
  CAMPO_NO_EDITABLE: "Ese campo no se puede editar desde aquí.",
  ROL_DUPLICADO: "Ya existe un rol con ese nombre.",
  PERMISO_INVALIDO: "Uno o más permisos seleccionados no existen.",
  ROL_SISTEMA_PROTEGIDO: "Ese rol de sistema no se puede modificar de esa forma.",
  ROL_CON_USUARIOS: "El rol tiene usuarios asignados: reasígnalos antes de eliminarlo.",
  NO_ENCONTRADO: "No se encontró el recurso.",
  CEDULA_DUPLICADA: "Ya existe un docente con esa cédula.",
  CODIGO_DUPLICADO: "Ya existe un docente con ese código.",
  FACULTAD_INVALIDA: "La facultad seleccionada no existe.",
  CATEGORIA_INVALIDA: "La categoría no existe o está inactiva.",
  CATEGORIA_DUPLICADA: "Ya existe una categoría con ese nombre.",
  CATEGORIA_EN_USO: "La categoría está en uso por docentes o trámites: dala de baja en vez de eliminarla.",
  DOCENTE_NO_ENCONTRADO: "No encontramos un docente con ese código o cédula en el padrón.",
  DOCENTE_INACTIVO: "El docente está dado de baja.",
  DOCENTE_NO_VERIFICADO: "El código y la cédula no corresponden a un docente activo del padrón. Revisa ambos datos o acércate a ventanilla.",
  ADJUNTOS_REQUERIDOS: "Adjunta entre 1 y 5 documentos de respaldo.",
  SOLICITUDES_PENDIENTES: "Ya tienes solicitudes pendientes de validación en ventanilla. Espera a que se validen antes de radicar otra.",
  SESION_NO_VERIFICABLE: "No se pudo verificar tu sesión. Intenta de nuevo en unos segundos.",
  TRANSICION_INVALIDA: "El trámite no puede pasar a ese estado desde el actual.",
  TRAMITE_REMITIDO: "El trámite fue remitido al Rectorado y es de solo lectura.",
  CARPETA_NO_SOLICITADA: "Antes de revisar hay que solicitar la carpeta física del trámite.",
  CARPETAS_ABIERTAS: "Hay carpetas del trámite sin devolver al archivo: regístralas antes de archivar.",
  PERIODO_REQUERIDO: "Registra los años a evaluar antes de dictaminar.",
  PERIODO_INVALIDO: "Los años a evaluar no son válidos (el hasta no puede ser anterior al desde ni futuro).",
  RESOLUCION_REQUERIDA: "Indica el número de la resolución y una fecha que no sea futura.",
  NOTA_REQUERIDA: "Indica la nota de remisión al Rectorado.",
  TRAMITE_INVALIDO: "El trámite no existe, no está a tu cargo, es de otro docente o ya está cerrado.",
  TRAMITE_ARCHIVADO: "El trámite está archivado y es de solo lectura.",
  MOTIVO_REQUERIDO: "La reasignación exige un motivo.",
  NO_ASIGNADO: "Este trámite no está asignado a tu usuario.",
  ANALISTA_INVALIDO: "El usuario seleccionado no es un analista activo.",
  OBSERVACIONES_PENDIENTES: "Hay observaciones pendientes: resuélvelas antes de dictaminar.",
  PENDIENTE_VALIDACION: "La solicitud del portal debe validarse en ventanilla antes de asignarla.",
  YA_VALIDADO: "La solicitud ya fue validada.",
  ARCHIVO_INVALIDO: `Los adjuntos deben ser PDF, JPG o PNG de hasta ${LIMITE_MB.portal} MB.`,
  DOCUMENTO_NO_ENCONTRADO: "Ese documento ya no está en el expediente: actualiza la página.",
  DOCUMENTO_DE_BAJA: "Ese documento ya está dado de baja: actualiza la página.",
  MOTIVO_BAJA_REQUERIDO: "Indica el motivo de la baja del documento.",
  DOCUMENTO_ACTIVO: "Ese documento ya fue restituido: actualiza la página.",
  MOTIVO_RESTITUCION_REQUERIDO: "Indica el motivo de la restitución.",
  PLAZO_INVALIDO: "El plazo no puede estar en el pasado.",
  TIPO_NO_PERMITIDO: "Solo se admiten archivos PDF, JPG o PNG.",
  ARCHIVO_GRANDE: `El archivo supera el tamaño máximo (PDF ${LIMITE_MB.pdf} MB, imágenes ${LIMITE_MB.imagen} MB).`,
  SUBIDA_INCOMPLETA: "El archivo no llegó al almacén. Vuelve a cargarlo.",
  DOCUMENTO_ESTADO_INVALIDO: "El documento no está en un estado que permita esta operación: actualiza la página.",
  SIN_OCR: "El documento aún no tiene texto extraído.",
  OCR_NO_DISPONIBLE:"El servicio de OCR no está disponible por ahora.",
  ALMACEN_NO_DISPONIBLE: "El almacenamiento de documentos no está disponible por ahora. Intenta más tarde.",
  LIMITE_TASA: "Demasiadas solicitudes. Espera un momento e intenta de nuevo.",
  ACCION_INVALIDA: "El tipo de acción seleccionado no existe.",
  LEGAJO_CODIGO_DUPLICADO: "Ya existe un legajo con ese código.",
  DOCENTE_INVALIDO: "El docente no existe.",
  LEGAJO_NO_DISPONIBLE: "La carpeta no está disponible en archivo.",
  LEGAJO_NO_PRESTADO: "La carpeta no figura como prestada.",
  SOLICITUD_ABIERTA: "Ya hay una solicitud en curso para esa carpeta.",
  SOLICITUD_ESTADO_INVALIDO: "La solicitud ya cambió de estado: actualiza la lista.",
  FECHA_LIMITE_INVALIDA: "La fecha límite debe ser posterior a hoy.",
  OBSERVACION_REQUERIDA: "Describe la novedad: la carpeta no está en buen estado.",
  ESTADO_MATERIAL_INVALIDO: "Selecciona el estado material de la carpeta.",
  RANGO_FECHAS_INVALIDO: "La fecha inicial no puede ser posterior a la fecha final.",
  REPORTE_VACIO: "No hay información para los filtros elegidos: amplía el periodo o quita algún filtro.",
  TIPO_REPORTE_INVALIDO: "Elige un tipo de reporte.",
  FORMATO_INVALIDO: "Elige el formato: PDF o Excel.",
  FILTRO_INVALIDO: "Uno de los filtros no tiene un valor válido.",
};

export function mensajeError(codigo: string | undefined, fallback: string): string {
  if (codigo && MENSAJES_API[codigo]) return MENSAJES_API[codigo];
  return fallback;
}

export async function leerErrorApi(respuesta: Response): Promise<{ codigo?: string; mensaje?: string }> {
  const cuerpo = await respuesta.json().catch(() => null);
  return { codigo: cuerpo?.codigo, mensaje: cuerpo?.mensaje };
}

// Estado de formulario con el error de la API ya traducido.
export async function errorApi(respuesta: Response, fallback: string): Promise<EstadoFormulario> {
  const { codigo } = await leerErrorApi(respuesta);
  return { error: mensajeError(codigo, fallback) };
}
