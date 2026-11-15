import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import { exigirPermiso } from "@/lib/api";
import SelectorDocente from "./SelectorDocente";

// Archivo elige al docente cuyo expediente digital va a completar con los escaneos de su carpeta física.
export default async function DocumentosPage({ searchParams }: { searchParams: Promise<{ buscar?: string; pagina?: string }> }) {
  await exigirPermiso("DOCUMENTO_ESCANEAR");
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina
        titulo="Expediente digital"
        descripcion="Elige al docente para subir los documentos escaneados de su carpeta (currículum, regularización, titular…)."
      />
      <SelectorDocente
        ruta="/panel/documentos"
        buscar={sp.buscar}
        pagina={pagina}
        etiquetaAccion="Abrir expediente"
        destino={(id) => `/panel/documentos/${id}`}
      />
    </div>
  );
}
