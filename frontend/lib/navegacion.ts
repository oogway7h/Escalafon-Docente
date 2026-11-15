export type ItemNav = {
  href: string;
  etiqueta: string;
  /** Qué se hace en esa pantalla. La usa la portada del panel; la barra lateral la ignora. */
  descripcion: string;
  permiso: string;
  disponible: boolean;
};

export type SeccionNav = {
  id: "seguridad" | "docente" | "tramites" | "documental" | "custodia" | "indicadores";
  etiqueta: string;
  items: ItemNav[];
};

// disponible:false = el permiso ya existe en el catálogo de seed, pero la pantalla es de una
// historia de usuario futura (ver Inicializacion.md, fuera de alcance de HU-01/HU-02).
export const NAVEGACION: SeccionNav[] = [
  {
    id: "seguridad",
    etiqueta: "Seguridad",
    items: [
      { href: "/panel/perfil", etiqueta: "Mi perfil", descripcion: "Actualiza tus datos de contacto y tu contraseña.", permiso: "PERFIL_GESTIONAR", disponible: true },
      { href: "/panel/usuarios", etiqueta: "Usuarios", descripcion: "Da de alta al personal, edita sus datos y gestiona las bajas.", permiso: "USUARIO_GESTIONAR", disponible: true },
      { href: "/panel/roles", etiqueta: "Roles y permisos", descripcion: "Define los roles del sistema y qué puede hacer cada uno.", permiso: "ROL_GESTIONAR", disponible: true },
      { href: "/panel/auditoria", etiqueta: "Auditoría", descripcion: "Consulta quién hizo qué y cuándo, con filtros por usuario y fecha.", permiso: "AUDITORIA_CONSULTAR", disponible: true },
    ],
  },
  {
    id: "docente",
    etiqueta: "Dominio",
    items: [
      { href: "/panel/docentes", etiqueta: "Docentes", descripcion: "Mantén el padrón docente: altas, datos y categoría de cada uno.", permiso: "DOCENTE_GESTIONAR", disponible: true },
      { href: "/panel/categorias", etiqueta: "Categorías", descripcion: "Administra las categorías del escalafón y los documentos que exige cada una.", permiso: "CATEGORIA_GESTIONAR", disponible: true },
    ],
  },
  {
    id: "tramites",
    etiqueta: "Trámites",
    items: [
      { href: "/panel/tramites/radicar", etiqueta: "Iniciar Trámite", descripcion: "Radica un trámite nuevo y genera su código de seguimiento.", permiso: "TRAMITE_RADICAR", disponible: true },
      { href: "/panel/tramites", etiqueta: "Seguimiento", descripcion: "Busca cualquier trámite y sigue su estado y su historial.", permiso: "TRAMITE_RADICAR", disponible: true },
      { href: "/panel/tramites/asignar", etiqueta: "Asignar", descripcion: "Reparte los trámites ingresados entre los analistas.", permiso: "TRAMITE_ASIGNAR", disponible: true },
      { href: "/panel/tramites/revisar", etiqueta: "Mi bandeja", descripcion: "Revisa los expedientes que tienes asignados y registra observaciones.", permiso: "TRAMITE_REVISAR", disponible: true },
      { href: "/panel/tramites/archivar", etiqueta: "Por archivar", descripcion: "Registra la ubicación del legajo de los trámites ya remitidos.", permiso: "TRAMITE_ARCHIVAR", disponible: true },
      { href: "/panel/notificaciones", etiqueta: "Notificaciones", descripcion: "Revisa los avisos enviados al docente y los envíos fallidos.", permiso: "NOTIFICACION_CONSULTAR", disponible: true },
    ],
  },
  {
    id: "documental",
    etiqueta: "Documental",
    items: [
      { href: "/panel/documentos", etiqueta: "Subir escaneos", descripcion: "Carga documentos al expediente y procésalos con OCR.", permiso: "DOCUMENTO_ESCANEAR", disponible: true },
      { href: "/panel/documentos/buscar", etiqueta: "Buscar", descripcion: "Busca por contenido, número de resolución o docente en todos los documentos.", permiso: "DOCUMENTO_BUSCAR", disponible: true },
      { href: "/panel/documentos/bajas", etiqueta: "Dados de baja", descripcion: "Restituye los documentos retirados por error para que vuelvan a las búsquedas.", permiso: "DOCUMENTO_RESTITUIR", disponible: true },
    ],
  },
  {
    id: "custodia",
    etiqueta: "Carpetas",
    items: [
      { href: "/panel/legajos", etiqueta: "Legajos", descripcion: "Mantén el inventario de legajos físicos y su ubicación exacta.", permiso: "LEGAJO_GESTIONAR", disponible: true },
      { href: "/panel/carpetas", etiqueta: "Solicitar carpeta", descripcion: "Pide al archivo la carpeta física de un docente.", permiso: "CARPETA_SOLICITAR", disponible: true },
      { href: "/panel/custodia", etiqueta: "Solicitudes y custodia", descripcion: "Atiende las solicitudes recibidas y registra salidas y devoluciones.", permiso: "CUSTODIA_REGISTRAR", disponible: true },
      { href: "/panel/custodia/bitacora", etiqueta: "Bitácora", descripcion: "Consulta el historial de movimientos de cada carpeta.", permiso: "BITACORA_CONSULTAR", disponible: true },
    ],
  },
  {
    id: "indicadores",
    etiqueta: "Indicadores",
    items: [
      { href: "/panel/tablero", etiqueta: "Tablero", descripcion: "Indicadores de gestión del periodo y alertas de carpetas en mora.", permiso: "TABLERO_VER", disponible: true },
      { href: "/panel/reportes", etiqueta: "Reportes", descripcion: "Exporta trámites, inventario y tiempos de atención en PDF o Excel.", permiso: "REPORTE_GENERAR", disponible: true },
    ],
  },
];
