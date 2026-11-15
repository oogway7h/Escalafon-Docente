"use server";

import { enviarJson } from "@/lib/api";
import { restablecerEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, type EstadoFormulario } from "@/lib/formularios";

export async function restablecerContrasena(
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = restablecerEsquema.safeParse({
    token: formData.get("token"),
    nueva: formData.get("nueva"),
    confirmar: formData.get("confirmar"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  let respuesta: Response;
  try {
    const { token, nueva } = analisis.data;
    respuesta = await enviarJson("/api/auth/restablecer", "POST", { token, nueva });
  } catch {
    return { error: "No se pudo conectar con el servidor. Intenta nuevamente." };
  }

  if (!respuesta.ok) return errorApi(respuesta, "No se pudo restablecer la contraseña.");

  return { mensaje: "Tu contraseña se actualizó. Ya puedes iniciar sesión con la nueva contraseña." };
}
