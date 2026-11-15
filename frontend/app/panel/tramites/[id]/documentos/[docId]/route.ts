import { NextResponse } from "next/server";
import { api } from "@/lib/api";

// Redirige a un enlace firmado y temporal del adjunto (válido unos minutos).
export async function GET(
  request: Request,
  { params }: { params: Promise<{ id: string; docId: string }> },
) {
  const { id, docId } = await params;
  const respuesta = await api(`/api/tramites/${id}/documentos/${docId}`);
  if (!respuesta.ok) {
    return NextResponse.redirect(new URL(`/panel/tramites/${id}?error=ALMACEN_NO_DISPONIBLE`, request.url));
  }
  const { url } = await respuesta.json();
  return NextResponse.redirect(url);
}
