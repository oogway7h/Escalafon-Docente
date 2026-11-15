import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import TablaTramites from "@/components/tramites/TablaTramites";
import { exigirPermiso } from "@/lib/api";
import { listarTramites } from "@/lib/api-tramites";

// Bandeja del Personal de Archivo: trámites ya remitidos al Rectorado, a la espera de archivar su carpeta.
export default async function PorArchivarPage() {
  await exigirPermiso("TRAMITE_ARCHIVAR");
  const remitidos = await listarTramites({ estado: "REMITIDO" });

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina
        titulo="Por archivar"
        descripcion={<>{remitidos.length} trámites remitidos al Rectorado esperan el archivo de su carpeta. Abre uno para registrar su ubicación.</>}
      />
      <TablaTramites tramites={remitidos} vacio="No hay trámites por archivar." />
    </div>
  );
}
