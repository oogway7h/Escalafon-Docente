import { NextResponse } from "next/server";
import { api } from "@/lib/api";
import { leerErrorApi, mensajeError } from "@/lib/formularios";

// Páginas de un adjunto del trámite que contienen lo buscado (índice OCR). Lo llama el visor desde el navegador:
// así buscar dentro de un documento ya abierto no vuelve a renderizar toda la página del trámite.
export async function GET(request: Request, { params }: { params: Promise<{ id: string; docId: string }> }) {
  const { id, docId } = await params;
  const q = new URL(request.url).searchParams.get("q") ?? "";
  const respuesta = await api(`/api/tramites/${id}/documentos/${docId}/buscar?q=${encodeURIComponent(q)}`);
  if (!respuesta.ok) {
    const { codigo } = await leerErrorApi(respuesta);
    return NextResponse.json({ error: mensajeError(codigo, "No se pudo buscar en el documento.") }, { status: respuesta.status });
  }
  return NextResponse.json(await respuesta.json());
}
