import Link from "next/link";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import { claseControl } from "@/components/common/Campo";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Tabla from "@/components/common/Tabla";
import EstadoInsignia from "@/components/tramites/EstadoInsignia";
import { exigirPermiso } from "@/lib/api";
import { listarCarpetas } from "@/lib/api-legajos";
import { listarTramites, obtenerTramite } from "@/lib/api-tramites";
import { mensajeError } from "@/lib/formularios";
import { tipoYCategoria } from "@/lib/tramites";
import { solicitarCarpetaAction } from "../actions";

// La carpeta se solicita para un trámite concreto: primero se elige el trámite (de los que el analista
// tiene a cargo) y luego una carpeta de ese docente que esté disponible en el archivo.
export default async function NuevaSolicitudPage({
  searchParams,
}: {
  searchParams: Promise<{ tramite?: string; error?: string }>;
}) {
  await exigirPermiso("CARPETA_SOLICITAR");
  const sp = await searchParams;
  const error = sp.error && <Alerta tipo="error">{mensajeError(sp.error, "No se pudo enviar la solicitud.")}</Alerta>;

  if (!sp.tramite) {
    const mios = (await listarTramites({ mios: true })).filter((t) => t.estado !== "ARCHIVADO" && t.estado !== "REMITIDO");
    return (
      <div className="flex flex-col gap-6">
        <EncabezadoPagina titulo="Solicitar carpeta" descripcion="Elige el trámite para el que necesitas la carpeta física del docente." />
        {error}
        <Tabla
          filas={mios}
          clave={(t) => t.id}
          href={(t) => `/panel/carpetas/nueva?tramite=${t.id}`}
          vacio="No tienes trámites asignados a los que solicitar una carpeta."
          columnas={[
            {
              titulo: "Trámite",
              celda: (t) => (
                <>
                  <div className="font-mono text-sm font-semibold">{t.codigo}</div>
                  <div className="text-xs text-ink-faint">
                    {tipoYCategoria(t)}
                  </div>
                </>
              ),
            },
            {
              titulo: "Docente",
              celda: (t) => (
                <>
                  <div className="font-semibold">{t.docenteNombre}</div>
                  <div className="text-xs text-ink-faint">Cód. {t.docenteCodigo}</div>
                </>
              ),
            },
            { titulo: "Estado", celda: (t) => <EstadoInsignia estado={t.estado} /> },
          ]}
        />
      </div>
    );
  }

  const t = await obtenerTramite(sp.tramite);
  const lista = t.docenteCodigo ? await listarCarpetas({ buscar: t.docenteCodigo, estado: "EN_ARCHIVO" }, 1, 20) : null;

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina
        titulo="Solicitar carpeta"
        descripcion={
          <>
            Trámite <span className="font-mono font-semibold">{t.codigo}</span> · {t.docenteNombre}. Pide al personal de archivo la carpeta de este docente.
          </>
        }
        accion={
          <Link href="/panel/carpetas/nueva" className="text-sm font-semibold text-ink-soft hover:text-accent-ink">
            Cambiar de trámite
          </Link>
        }
      />
      {error}
      <Tabla
        filas={lista?.datos ?? []}
        clave={(l) => l.id}
        vacio="El docente no tiene carpetas disponibles en el archivo (pueden estar prestadas o sin inventariar)."
        columnas={[
          { titulo: "Carpeta", className: "w-28", celda: (l) => <span className="font-mono text-sm font-semibold">{l.codigo}</span> },
          { titulo: "Ubicación", celda: (l) => <span className="text-sm text-ink-soft">Gaveta {l.gaveta || "s/g"}</span> },
          {
            titulo: "Solicitar",
            className: "w-96",
            celda: (l) => (
              <form action={solicitarCarpetaAction} className="flex items-center gap-2">
                <input type="hidden" name="legajoId" value={l.id} />
                <input type="hidden" name="tramiteId" value={t.id} />
                <input
                  name="motivo"
                  required
                  defaultValue="Revisión del expediente"
                  placeholder="Motivo"
                  aria-label="Motivo"
                  className={`${claseControl()} min-w-0 flex-1 !py-1.5 !text-sm`}
                />
                <BotonEnviar tamano="sm" textoEnviando="Solicitando…">
                  Solicitar
                </BotonEnviar>
              </form>
            ),
          },
        ]}
      />
    </div>
  );
}
