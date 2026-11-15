"use server";

import { revalidatePath } from "next/cache";
import { enviarJson } from "@/lib/api";
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
