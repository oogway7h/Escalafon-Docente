"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { z } from "zod";
import { enviarJson } from "@/lib/api";
import { camposDeZod, errorApi, type EstadoFormulario } from "@/lib/formularios";

const legajoEsquema = z.object({
  docenteId: z.string().optional(),
  codigo: z.string().max(20, "Máximo 20 caracteres").optional(),
  gaveta: z.string().max(10, "Máximo 10 caracteres").optional(),
  orden: z.coerce.number().int("Debe ser un entero").min(0, "No puede ser negativo"),
  anexos: z.string().max(80, "Máximo 80 caracteres").optional(),
  observacion: z.string().max(200, "Máximo 200 caracteres").optional(),
});

// CU-19: crea (id = null) o edita una carpeta. El docente solo se elige al crear.
export async function guardarLegajoAction(id: string | null, _previo: EstadoFormulario, formData: FormData): Promise<EstadoFormulario> {
  const analisis = legajoEsquema.safeParse({
    docenteId: formData.get("docenteId") || undefined,
    codigo: formData.get("codigo") || undefined,
    gaveta: formData.get("gaveta") || undefined,
    orden: formData.get("orden") || 0,
    anexos: formData.get("anexos") || undefined,
    observacion: formData.get("observacion") || undefined,
  });
  if (!analisis.success) return camposDeZod(analisis.error);
  if (!id && !analisis.data.docenteId) return { camposInvalidos: { docenteId: "Busca y selecciona un docente" } };

  const respuesta = await enviarJson(id ? `/api/legajos/${id}` : "/api/legajos", id ? "PUT" : "POST", analisis.data);
  if (!respuesta.ok) return errorApi(respuesta, "No se pudo guardar el legajo.");

  revalidatePath("/panel/legajos");
  redirect("/panel/legajos");
}
