import { EnlaceBoton } from "@/components/common/EnlaceCarga";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Tarjeta from "@/components/common/Tarjeta";
import TablaTramites from "@/components/tramites/TablaTramites";
import { exigirPermiso } from "@/lib/api";
import { listarAnalistas, listarTramites } from "@/lib/api-tramites";

export default async function AsignarPage() {
  await exigirPermiso("TRAMITE_ASIGNAR");
  const [pendientes, analistas] = await Promise.all([listarTramites({ sinAsignar: true }), listarAnalistas()]);
  const maxCarga = Math.max(1, ...analistas.map((a) => a.carga));

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina
        titulo="Asignar trámites"
        descripcion={<>{pendientes.length} trámites esperan analista. Reparte según la carga vigente.</>}
      />

      <div className="flex flex-col gap-5 lg:flex-row lg:items-start">
        <div className="flex-[1.6]">
          <TablaTramites
            tramites={pendientes}
            vacio="No hay trámites sin asignar."
            accion={(t) => (
              <div className="flex justify-end">
                <EnlaceBoton href={`/panel/tramites/${t.id}`} variante="secundario" tamano="sm">
                  {t.pendienteValidacion ? "Revisar solicitud" : "Asignar"}
                </EnlaceBoton>
              </div>
            )}
          />
        </div>

        <Tarjeta titulo="Carga vigente" descripcion="Trámites activos por analista (asignados, en revisión u observados)." className="max-w-[360px] flex-1">
          <ul className="flex flex-col gap-4">
            {analistas.map((a) => (
              <li key={a.id} className="flex flex-col gap-1.5">
                <div className="flex justify-between text-sm">
                  <span className="font-semibold text-ink">{a.nombre}</span>
                  <span className="text-ink-faint">{a.carga}</span>
                </div>
                <div className="h-1.5 overflow-hidden rounded-full bg-bg" role="presentation">
                  <div className="h-full rounded-full bg-accent" style={{ width: `${(a.carga / maxCarga) * 100}%` }} />
                </div>
              </li>
            ))}
          </ul>
        </Tarjeta>
      </div>
    </div>
  );
}
