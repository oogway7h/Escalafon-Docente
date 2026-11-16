import { z } from "zod";

export const loginEsquema = z.object({
  correo: z.email("Ingresa un correo válido"),
  contrasena: z.string().min(1, "Ingresa tu contraseña"),
});

export type LoginDatos = z.infer<typeof loginEsquema>;

export const recuperarEsquema = z.object({
  correo: z.email("Ingresa un correo válido"),
});

// Misma política que valida el backend (sección 6): 8+ caracteres, una mayúscula, un dígito.
const contrasenaPolitica = z
  .string()
  .min(8, "Debe tener al menos 8 caracteres")
  .regex(/[A-Z]/, "Debe incluir al menos una mayúscula")
  .regex(/[0-9]/, "Debe incluir al menos un dígito");

export const restablecerEsquema = z
  .object({
    token: z.string().min(1, "Falta el token del enlace"),
    nueva: contrasenaPolitica,
    confirmar: z.string().min(1, "Repite la contraseña"),
  })
  .refine((datos) => datos.nueva === datos.confirmar, {
    message: "Las contraseñas no coinciden",
    path: ["confirmar"],
  });

export const activarEsquema = restablecerEsquema;

export const contactoEsquema = z.object({
  telefono: z.string().min(1, "Ingresa un teléfono"),
});

export const cambiarContrasenaEsquema = z
  .object({
    actual: z.string().min(1, "Ingresa tu contraseña actual"),
    nueva: contrasenaPolitica,
    confirmar: z.string().min(1, "Repite la contraseña"),
  })
  .refine((datos) => datos.nueva === datos.confirmar, {
    message: "Las contraseñas no coinciden",
    path: ["confirmar"],
  });

export const usuarioEsquema = z.object({
  nombres: z.string().min(1, "Ingresa el nombre"),
  apellidos: z.string().min(1, "Ingresa el apellido"),
  correo: z.email("Ingresa un correo válido"),
  telefono: z.string().optional(),
  rolId: z.coerce.number("Selecciona un rol").int().positive("Selecciona un rol"),
});

export const cambiarRolEsquema = z.object({
  rolId: z.coerce.number("Selecciona un rol").int().positive("Selecciona un rol"),
});

export const rolEsquema = z.object({
  nombre: z.string().min(1, "Ingresa un nombre"),
  descripcion: z.string().optional(),
  permisoIds: z.array(z.coerce.number()).default([]),
});

// --- CU-07 al CU-11 ---

export const docenteEsquema = z.object({
  codigo: z.string().min(1, "Ingresa el código del docente"),
  cedula: z.string().optional(),
  nombres: z.string().min(1, "Ingresa los nombres"),
  apellidos: z.string().min(1, "Ingresa los apellidos"),
  correo: z.union([z.email("Ingresa un correo válido"), z.literal("")]).optional(),
  telefono: z.string().optional(),
  categoriaId: z.coerce.number().int().positive().optional(),
  facultadId: z.coerce.number().int().positive().optional(),
});

export const categoriaEsquema = z.object({
  nombre: z.string().min(1, "Ingresa el nombre"),
  aniosAntiguedad: z.coerce.number().int().min(0, "Debe ser 0 o más"),
  // Vacío = sin puntaje (no aparece en las tarjetas del portal).
  puntajeDesde: z.preprocess((v) => (v === "" || v == null ? undefined : v), z.coerce.number().int().min(0, "Debe ser 0 o más").optional()),
  puntajeHasta: z.preprocess((v) => (v === "" || v == null ? undefined : v), z.coerce.number().int().min(0, "Debe ser 0 o más").optional()),
  requisitos: z.array(z.string()),
});

// La categoría solo se exige en la categorización: una regularización no postula a ninguna (viaja como 0).
export const radicarEsquema = z
  .object({
    docenteId: z.string().min(1, "Ingresa el código de un docente válido"),
    tipo: z.enum(["CATEGORIZACION", "REGULARIZACION"], "Selecciona el tipo de trámite"),
    categoriaId: z.coerce.number().int().min(0).catch(0),
  })
  .refine((d) => d.tipo === "REGULARIZACION" || d.categoriaId > 0, { path: ["categoriaId"], message: "Selecciona la categoría" })
  .transform((d) => (d.tipo === "REGULARIZACION" ? { ...d, categoriaId: 0 } : d));

export const asignarEsquema = z.object({
  analistaId: z.string().min(1, "Selecciona un analista"),
  motivo: z.string().optional(),
});

export const observacionEsquema = z.object({
  documento: z.string().min(1, "Indica el documento afectado"),
  motivo: z.string().min(5, "Describe el motivo (mínimo 5 caracteres)"),
  plazo: z.string().min(1, "Indica el plazo"),
});

export const dictamenEsquema = z.object({
  resultado: z.enum(["FAVORABLE", "DESFAVORABLE"], "Selecciona el resultado"),
  texto: z.string().min(10, "Fundamenta el dictamen (mínimo 10 caracteres)"),
});

export const archivarEsquema = z.object({
  ubicacion: z.string().min(1, "Indica dónde queda archivada la carpeta"),
});

// Etapas posteriores al dictamen.

export const periodoEsquema = z
  .object({
    anioDesde: z.coerce.number().int("Debe ser un año").min(1950, "Año no válido").max(new Date().getFullYear(), "No puede ser futuro"),
    anioHasta: z.coerce.number().int("Debe ser un año").min(1950, "Año no válido").max(new Date().getFullYear(), "No puede ser futuro"),
  })
  .refine((v) => v.anioHasta >= v.anioDesde, { path: ["anioHasta"], message: "No puede ser anterior al año inicial" });

export const regresoEsquema = z.object({
  resolucionNumero: z.string().trim().min(1, "Indica el número de la resolución"),
  resolucionFecha: z.string().min(1, "Indica la fecha de la resolución"),
});

export const remisionEsquema = z.object({
  nota: z.string().trim().min(3, "Indica la nota de remisión (número o referencia)"),
});

// --- CU-16 ---

export const bajaDocumentoEsquema = z.object({
  documentoId: z.string().min(1, "Selecciona el documento"),
  motivo: z.string().min(5, "Describe el motivo de la baja (mínimo 5 caracteres)"),
});
