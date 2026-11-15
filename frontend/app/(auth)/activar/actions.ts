"use server";

import { enviarJson } from "@/lib/api";
import { activarEsquema } from "@/lib/esquemas";
import { camposDeZod, errorApi, type EstadoFormulario } from "@/lib/formularios";

export async function activarCuenta(
  _estadoPrevio: EstadoFormulario,
  formData: FormData,
): Promise<EstadoFormulario> {
  const analisis = activarEsquema.safeParse({
    token: formData.get("token"),
    nueva: formData.get("nueva"),
    confirmar: formData.get("confirmar"),
  });
  if (!analisis.success) return camposDeZod(analisis.error);

  let respuesta: Response;
  try {
    const { token, nueva } = analisis.data;
    respuesta = await enviarJson("/api/auth/activar", "POST", { token, nueva });
  } catch {
    return { error: "No se pudo conectar con el servidor. Intenta nuevamente." };
  }

  if (!respuesta.ok) return errorApi(respuesta, "No se pudo activar la cuenta.");

  return { mensaje: "Tu cuenta quedó activada. Ya puedes iniciar sesión." };
}
