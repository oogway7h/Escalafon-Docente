import { NextResponse } from "next/server";
import { interpretarReporte } from "@/lib/api-indicadores";

// Reenvía la frase del formulario a la API; el navegador nunca habla con ella.
export async function POST(request: Request) {
  const cuerpo = await request.json().catch(() => null);
  const frase = typeof cuerpo?.frase === "string" ? cuerpo.frase : "";
  const categorias = Array.isArray(cuerpo?.categorias) ? cuerpo.categorias : [];
  const respuesta = await interpretarReporte(frase, categorias);
  const datos = await respuesta.json().catch(() => null);
  return NextResponse.json(datos ?? { codigo: "ERROR_INTERPRETE" }, { status: respuesta.status, headers: { "cache-control": "private, no-store" } });
}
