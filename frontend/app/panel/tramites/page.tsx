import Link from "next/link";
import { claseBoton } from "@/components/common/Boton";
import { claseControl } from "@/components/common/Campo";
import { BotonFiltrar, FiltrosCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Paginacion from "@/components/common/Paginacion";
import TablaTramites from "@/components/tramites/TablaTramites";
import { obtenerSesion } from "@/lib/api";
import { listarTramitesPagina } from "@/lib/api-tramites";
import { ETIQUETA_ESTADO, FASES } from "@/lib/tramites";
import { notFound } from "next/navigation";

const PERMISOS = ["TRAMITE_RADICAR", "TRAMITE_ASIGNAR", "TRAMITE_REVISAR", "TRAMITE_ARCHIVAR"]; // los mismos que admite GET /api/tramites

export default async function TramitesPage({
  searchParams,
}: {
  searchParams: Promise<{ estado?: string; buscar?: string; pagina?: string }>;
}) {
  const sesion = await obtenerSesion();
  if (!PERMISOS.some((p) => sesion.permisos.includes(p))) notFound();
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const lista = await listarTramitesPagina(
    { estado: sp.estado, buscar: sp.buscar },
    pagina,
  );
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Trámites"
          descripcion={
            <>
              {lista.total} trámites
              {sp.estado
                ? ` en estado ${ETIQUETA_ESTADO[sp.estado as keyof typeof ETIQUETA_ESTADO] ?? sp.estado}`
                : ""}
              .
            </>
          }
          accion={
            sesion.permisos.includes("TRAMITE_RADICAR") && (
              <Link
                href="/panel/tramites/radicar"
                className={claseBoton("primario", "sm")}
              >
                + Radicar trámite
              </Link>
            )
          }
        />

        <FiltrosCarga action="/panel/tramites" className="flex flex-wrap gap-2.5">
          <input
            type="search"
            name="buscar"
            defaultValue={sp.buscar}
            placeholder="Código, cédula o nombre del docente"
            aria-label="Buscar trámite"
            className={`${claseControl()} min-w-0 max-w-[360px] flex-1`}
          />
          <select name="estado" defaultValue={sp.estado ?? ""} aria-label="Estado" className={claseControl()}>
            <option value="">Todos los estados</option>
            {FASES.map((f) => (
              <option key={f} value={f}>
                {ETIQUETA_ESTADO[f]}
              </option>
            ))}
          </select>
          <BotonFiltrar>Buscar</BotonFiltrar>
        </FiltrosCarga>

        <ZonaCarga>
          <TablaTramites
            tramites={lista.datos}
            vacio="No se encontraron trámites."
          />
        </ZonaCarga>

        <Paginacion
          pagina={lista.pagina}
          totalPaginas={totalPaginas}
          searchParams={{ estado: sp.estado, buscar: sp.buscar }}
        />
      </div>
    </ListadoConCarga>
  );
}
