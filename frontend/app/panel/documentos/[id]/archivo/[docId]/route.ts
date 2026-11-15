import { entregarArchivo } from "@/lib/archivo";

// Archivo del expediente de un docente, para el visor del sistema.
export async function GET(request: Request, { params }: { params: Promise<{ id: string; docId: string }> }) {
  const { id, docId } = await params;
  return entregarArchivo(request, `/api/docentes/${id}/documentos/${docId}`);
}
