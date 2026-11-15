import { entregarArchivo } from "@/lib/archivo";

// Adjunto de un trámite, para el visor del sistema.
export async function GET(request: Request, { params }: { params: Promise<{ id: string; docId: string }> }) {
  const { id, docId } = await params;
  return entregarArchivo(request, `/api/tramites/${id}/documentos/${docId}`);
}
