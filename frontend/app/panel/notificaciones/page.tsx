import { claseControl } from "@/components/common/Campo";
import { BotonFiltrar, FiltrosCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Insignia from "@/components/common/Insignia";
import Paginacion from "@/components/common/Paginacion";
import Tabla from "@/components/common/Tabla";
import { exigirPermiso } from "@/lib/api";
import { listarNotificaciones } from "@/lib/api-legajos";
import { ETIQUETA_ESTADO, formatearFecha } from "@/lib/tramites";

const ESTADO_ENVIO = {
  ENVIADO: { etiqueta: "Enviado", tono: "exito" },
  PENDIENTE: { etiqueta: "Pendiente", tono: "neutro" },
  FALLIDO: { etiqueta: "Fallido", tono: "acento" },
} as const;

// CU-14: registro de los avisos por correo al docente (el Director revisa los que fallaron).
export default async function NotificacionesPage({ searchParams }: { searchParams: Promise<{ estado?: string; pagina?: string }> }) {
  await exigirPermiso("NOTIFICACION_CONSULTAR");
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const lista = await listarNotificaciones({ estado: sp.estado }, pagina);
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina titulo="Notificaciones a docentes" descripcion={<>{lista.total} avisos registrados. Los fallidos se reintentan hasta 3 veces.</>} />

        <FiltrosCarga action="/panel/notificaciones" className="flex flex-wrap items-center gap-2.5">
          <select name="estado" defaultValue={sp.estado ?? ""} aria-label="Estado del envío" className={claseControl()}>
            <option value="">Todos los envíos</option>
            <option value="ENVIADO">Enviados</option>
            <option value="PENDIENTE">Pendientes</option>
            <option value="FALLIDO">Fallidos</option>
          </select>
          <BotonFiltrar>Filtrar</BotonFiltrar>
        </FiltrosCarga>

        <ZonaCarga>
          <Tabla
            filas={lista.datos}
            clave={(n) => n.id}
            vacio="No hay notificaciones con ese filtro."
            columnas={[
              { titulo: "Fecha", celda: (n) => formatearFecha(n.creadoEn, true) },
              { titulo: "Trámite", celda: (n) => <span className="font-mono text-sm font-semibold">{n.tramiteCodigo}</span> },
              { titulo: "Nuevo estado", celda: (n) => ETIQUETA_ESTADO[n.estadoTramite] ?? n.estadoTramite },
              { titulo: "Destino", celda: (n) => n.destino || "—" },
              {
                titulo: "Envío",
                celda: (n) => (
                  <>
                    <Insignia tono={ESTADO_ENVIO[n.estado].tono}>{ESTADO_ENVIO[n.estado].etiqueta}</Insignia>
                    {n.error && <div className="mt-1 text-xs text-ink-faint">{n.error}</div>}
                  </>
                ),
              },
              { titulo: "Intentos", className: "w-20", celda: (n) => n.intentos },
            ]}
          />
        </ZonaCarga>

        <Paginacion pagina={lista.pagina} totalPaginas={totalPaginas} searchParams={{ estado: sp.estado }} />
      </div>
    </ListadoConCarga>
  );
}
