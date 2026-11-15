"use server";

import { enviarJson } from "@/lib/api";
import { recuperarEsquema } from "@/lib/esquemas";
import { camposDeZod, type EstadoFormulario } from "@/lib/formularios";

export async function solicitarRecuperacion(
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = recuperarEsquema.safeParse({ correo: formData.get("correo") });
  if (!analisis.success) return camposDeZod(analisis.error);

  try {
    await enviarJson("/api/auth/recuperacion", "POST", analisis.data);
  } catch {
    return { error: "No se pudo conectar con el servidor. Intenta nuevamente." };
  }

  // La API responde 202 exista o no el correo (sección 8): nunca revelamos cuál es el caso.
  return { mensaje: "Si el correo está registrado, te enviamos un enlace para restablecer tu contraseña." };
}
