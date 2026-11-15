"use server";

import { redirect } from "next/navigation";
import { revalidatePath } from "next/cache";
import { enviarJson } from "@/lib/api";
import { leerErrorApi } from "@/lib/formularios";

// Pantalla unificada de solicitudes y custodia (CU-21, CU-22 y CU-23). Tras cada acción se vuelve a la
// pestaña de origen (`tab`) y, si falló, se conserva la solicitud elegida (`sel`) para que el error se lea en su panel.
const RUTA = "/panel/custodia";

function destino(formData: FormData, extra: Record<string, string>): string {
  const qs = new URLSearchParams();
  const tab = String(formData.get("tab") ?? "");
  if (tab) qs.set("tab", tab);
  for (const [k, v] of Object.entries(extra)) if (v) qs.set(k, v);
  const s = qs.toString();
  return s ? `${RUTA}?${s}` : RUTA;
}

async function terminar(respuesta: Response, formData: FormData, exito: string): Promise<never> {
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    const error = codigo === "VALIDACION" ? "ESTADO_MATERIAL_INVALIDO" : (codigo ?? "ERROR");
    redirect(destino(formData, { error, sel: String(formData.get("sel") ?? "") }));
  }
  revalidatePath(RUTA, "layout");
  revalidatePath("/panel/legajos", "layout");
  redirect(destino(formData, { exito }));
}

// Aprueba (si sigue pendiente) y entrega la carpeta en un solo paso: fecha límite, estado material y observación.
export async function entregarAction(formData: FormData) {
  return terminar(
    await enviarJson(`/api/solicitudes-carpeta/${formData.get("id")}/entregar`, "POST", {
      fechaLimite: formData.get("fechaLimite"),
      estadoMaterial: formData.get("estadoMaterial"),
      observacion: formData.get("observacion") ?? "",
    }),
    formData,
    "ENTREGADA",
  );
}

// Solo aprueba (la carpeta queda lista para entregar después) o rechaza con motivo.
export async function atenderAction(formData: FormData) {
  const aprobar = formData.get("aprobar") === "true";
  return terminar(
    await enviarJson(`/api/solicitudes-carpeta/${formData.get("id")}/atender`, "POST", { aprobar, nota: formData.get("nota") ?? "" }),
    formData,
    aprobar ? "APROBADA" : "RECHAZADA",
  );
}

// Recibe la carpeta de vuelta y registra su estado material.
export async function devolverAction(formData: FormData) {
  return terminar(
    await enviarJson(`/api/legajos/${formData.get("legajoId")}/devolver`, "POST", {
      estadoMaterial: formData.get("estadoMaterial"),
      observacion: formData.get("observacion") ?? "",
    }),
    formData,
    "DEVUELTA",
  );
}
