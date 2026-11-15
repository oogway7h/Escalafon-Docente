import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import TablaTramites from "@/components/tramites/TablaTramites";
import { exigirPermiso } from "@/lib/api";
import { listarTramites } from "@/lib/api-tramites";

export default async function MiBandejaPage() {
  await exigirPermiso("TRAMITE_REVISAR");
  const todos = await listarTramites({ mios: true });
  const activos = todos.filter((t) => t.estado !== "ARCHIVADO" && t.estado !== "REMITIDO");
  const archivados = todos.length - activos.length;

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina
        titulo="Mi bandeja"
        descripcion={<>{activos.length} trámites asignados a ti{archivados > 0 ? ` · ${archivados} remitidos o archivados` : ""}.</>}
      />
      <TablaTramites tramites={activos} vacio="No tienes trámites asignados por ahora." />
    </div>
  );
}
