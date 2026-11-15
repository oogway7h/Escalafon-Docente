import { NextResponse } from "next/server";
import { api } from "@/lib/api";

// Texto extraído por OCR de un documento (se pide aparte porque puede pesar varios MB).
export async function GET(_request: Request, { params }: { params: Promise<{ id: string; docId: string }> }) {
  const { id, docId } = await params;
  const respuesta = await api(`/api/tramites/${id}/documentos/${docId}/texto`);
  if (!respuesta.ok) return NextResponse.json({ error: true }, { status: respuesta.status });
  return NextResponse.json(await respuesta.json());
}
