import Alerta from "@/components/common/Alerta";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Insignia from "@/components/common/Insignia";
import Paginacion from "@/components/common/Paginacion";
import Tabla from "@/components/common/Tabla";
import { exigirPermiso } from "@/lib/api";
import { listarSolicitudes } from "@/lib/api-legajos";
import { ETIQUETA_SOLICITUD, TONO_SOLICITUD, formatearFecha } from "@/lib/legajos";
import SolicitarCarpetaModal from "./SolicitarCarpetaModal";

// CU-20: las solicitudes del analista (el backend filtra por usuario) y acceso a una nueva.
export default async function CarpetasPage({ searchParams }: { searchParams: Promise<{ pagina?: string; exito?: string }> }) {
  await exigirPermiso("CARPETA_SOLICITAR");
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const lista = await listarSolicitudes({}, pagina);
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina
        titulo="Mis carpetas"
        descripcion={<>{lista.total} solicitudes de carpeta física.</>}
        accion={
          <SolicitarCarpetaModal />
        }
      />
      {sp.exito === "SOLICITADA" && <Alerta tipo="exito">Solicitud enviada al personal de archivo.</Alerta>}

      <Tabla
        filas={lista.datos}
        clave={(s) => s.id}
        vacio="Aún no solicitaste ninguna carpeta."
        columnas={[
          { titulo: "Carpeta", className: "w-28", celda: (s) => <span className="font-mono text-sm font-semibold">{s.legajoCodigo}</span> },
          {
            titulo: "Docente y motivo",
            celda: (s) => (
              <>
                <div className="font-semibold">{s.docenteNombre}</div>
                {s.tramiteCodigo && <div className="font-mono text-xs text-ink-faint">Trámite {s.tramiteCodigo}</div>}
                <div className="text-xs text-ink-soft">{s.motivo}</div>
                {s.estado === "RECHAZADA" && s.respuesta && <div className="text-xs text-accent-ink">Rechazada: {s.respuesta}</div>}
              </>
            ),
          },
          { titulo: "Solicitada", celda: (s) => formatearFecha(s.creadoEn) },
          { titulo: "Devolver antes de", celda: (s) => formatearFecha(s.fechaLimite) },
          { titulo: "Estado", celda: (s) => <Insignia tono={TONO_SOLICITUD[s.estado]}>{ETIQUETA_SOLICITUD[s.estado]}</Insignia> },
        ]}
      />

      <Paginacion pagina={lista.pagina} totalPaginas={totalPaginas} />
    </div>
  );
}
