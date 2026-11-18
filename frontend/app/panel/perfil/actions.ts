"use server";

import { revalidatePath } from "next/cache";
import { api, enviarJson } from "@/lib/api";
import { cambiarContrasenaEsquema, contactoEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, type EstadoFormulario } from "@/lib/formularios";
import { escribirTokenSesion } from "@/lib/sesion";

export async function actualizarContacto(
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = contactoEsquema.safeParse({ telefono: formData.get("telefono") });
  if (!analisis.success) return camposDeZod(analisis.error);

  const respuesta = await enviarJson("/api/perfil", "PUT", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo actualizar el teléfono.");

  revalidatePath("/panel/perfil");
  return { mensaje: "Teléfono actualizado." };
}

export async function cambiarContrasenaPerfil(
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = cambiarContrasenaEsquema.safeParse({
    actual: formData.get("actual"),
    nueva: formData.get("nueva"),
    confirmar: formData.get("confirmar"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  const { actual, nueva } = analisis.data;
  const respuesta = await enviarJson("/api/perfil/contrasena", "PUT", { actual, nueva });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo cambiar la contraseña.");

  // La API emite un JWT nuevo para la sesión actual (sección 8): reemplazamos la cookie.
  const { token, expiraEn } = await respuesta.json();
  await escribirTokenSesion(token, expiraEn);

  revalidatePath("/panel/perfil");
  return { mensaje: "Contraseña actualizada." };
}

const FOTO_MAX = 2 * 1024 * 1024;

// El navegador ya valida tipo y tamaño; aquí se repite porque el cliente no es de fiar. La API valida
// además el contenido real de la imagen.
export async function subirFoto(_estadoPrevio: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const foto = formData.get("foto");
  if (!(foto instanceof File) || foto.size === 0) return { error: "Elige una imagen JPG o PNG." };
  if (!["image/jpeg", "image/png"].includes(foto.type)) return { error: "La foto debe ser JPG o PNG." };
  if (foto.size > FOTO_MAX) return { error: "La foto no puede superar los 2 MB." };

  const cuerpo = new FormData();
  cuerpo.set("foto", foto);
  const respuesta = await api("/api/perfil/foto", { method: "PUT", body: cuerpo });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo guardar la foto.");

  revalidatePath("/panel", "layout"); // el avatar del encabezado también cambia
  return { mensaje: "Foto actualizada." };
}

export async function quitarFoto(_estadoPrevio: EstadoFormulario): Promise<EstadoFormulario> {
  const respuesta = await api("/api/perfil/foto", { method: "DELETE" });
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo quitar la foto.");

  revalidatePath("/panel", "layout");
  return { mensaje: "Foto eliminada." };
}
