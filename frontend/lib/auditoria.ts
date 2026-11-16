// Sin "server-only": etiquetas y formato de la bitácora de auditoría (CU-06).

// Códigos del catálogo del backend (domain/auditoria). Un código nuevo sin etiqueta se muestra tal cual.
export const ETIQUETA_ACCION: Record<string, string> = {
  LOGIN: "Inicio de sesión",
  LOGIN_FALLIDO: "Inicio de sesión fallido",
  RESTABLECER_CONTRASENA: "Restablecer contraseña",
  ACTIVAR_CUENTA: "Activar cuenta",
  CAMBIAR_CONTRASENA: "Cambiar contraseña",
  ACTUALIZAR_CONTACTO: "Actualizar contacto",
  CREAR: "Crear",
  EDITAR: "Editar",
  CAMBIO_ROL: "Cambiar rol",
  BAJA: "Dar de baja",
  REACTIVAR: "Reactivar",
  REENVIAR_ACTIVACION: "Reenviar activación",
  GUARDAR_ROL: "Guardar rol",
  ELIMINAR_ROL: "Eliminar rol",
  ELIMINAR_CATEGORIA: "Eliminar categoría",
  RADICAR_TRAMITE: "Radicar trámite",
  RADICAR_TRAMITE_PORTAL: "Radicar desde el portal",
  VALIDAR_SOLICITUD: "Validar solicitud",
  ASIGNAR_TRAMITE: "Asignar trámite",
  REASIGNAR_TRAMITE: "Reasignar trámite",
  INICIAR_REVISION: "Iniciar revisión",
  OBSERVAR_TRAMITE: "Observar trámite",
  RESOLVER_OBSERVACION: "Resolver observación",
  DICTAMINAR_TRAMITE: "Dictaminar trámite",
  ARCHIVAR_TRAMITE: "Archivar trámite",
  BAJA_DOCUMENTO: "Dar de baja documento",
  RESTITUIR_DOCUMENTO: "Restituir documento",
  SOLICITAR_CARPETA: "Solicitar carpeta",
  ATENDER_SOLICITUD: "Aprobar solicitud de carpeta",
  RECHAZAR_SOLICITUD: "Rechazar solicitud de carpeta",
  DESPACHAR_CARPETA: "Despachar carpeta",
  DEVOLVER_CARPETA: "Recibir devolución de carpeta",
  GENERAR_REPORTE: "Generar reporte",
  ACEPTAR_TERMINOS: "Aceptar acuerdo de usuario",
};

export function etiquetaAccion(codigo: string): string {
  return ETIQUETA_ACCION[codigo] ?? codigo;
}

// Acciones que conviene resaltar al revisar la bitácora.
export function tonoAccion(codigo: string): "acento" | "neutro" {
  return codigo === "LOGIN_FALLIDO" || codigo.startsWith("BAJA") || codigo.startsWith("ELIMINAR_") ? "acento" : "neutro";
}

// La API filtra por días de Bolivia: las fechas se muestran en esa misma zona horaria.
export function formatearFechaHora(iso: string): string {
  return new Date(iso).toLocaleString("es-BO", {
    timeZone: "America/La_Paz",
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  });
}

// rolAnterior -> "rol anterior"
function etiquetaCampo(clave: string): string {
  return clave.replace(/([a-z])([A-Z])/g, "$1 $2").toLowerCase();
}

function texto(valor: unknown): string {
  if (valor === null || valor === undefined || valor === "") return "—";
  if (Array.isArray(valor)) return valor.length ? valor.map(texto).join(", ") : "ninguno";
  if (typeof valor === "object") return JSON.stringify(valor);
  return String(valor);
}

// Resume el detalle en líneas legibles. Si trae {anterior, nuevo}, muestra solo los campos que cambiaron.
export function resumirDetalle(detalle: Record<string, unknown> | null): string[] {
  if (!detalle) return [];
  const { anterior, nuevo } = detalle as { anterior?: unknown; nuevo?: unknown };
  if (anterior && nuevo && typeof anterior === "object" && typeof nuevo === "object") {
    const a = anterior as Record<string, unknown>;
    const n = nuevo as Record<string, unknown>;
    return Object.keys({ ...a, ...n })
      .filter((k) => texto(a[k]) !== texto(n[k]))
      .map((k) => `${etiquetaCampo(k)}: ${texto(a[k])} → ${texto(n[k])}`);
  }
  return Object.entries(detalle).map(([k, v]) => `${etiquetaCampo(k)}: ${texto(v)}`);
}
