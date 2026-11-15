import type { Metadata } from "next";
import Alerta from "@/components/common/Alerta";
import { BuscadorCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Paginacion from "@/components/common/Paginacion";
import Tabla from "@/components/common/Tabla";
import Tarjeta from "@/components/common/Tarjeta";
import { exigirPermiso } from "@/lib/api";
import { listarDocumentosDeBaja, type DocumentoDeBaja } from "@/lib/api-tramites";
import { ETIQUETA_ESTADO, ETIQUETA_TIPO_DOCUMENTO, formatearFecha } from "@/lib/tramites";
import RestituirForm from "./RestituirForm";

export const metadata: Metadata = { title: "Documentos dados de baja — Escalafón Docente" };

const RUTA = "/panel/documentos/bajas";

// CU-17: bandeja del Director con los documentos retirados por baja lógica (CU-16). Al elegir uno
// (?doc=) se ve su detalle y se restituye con un motivo; vuelve a las búsquedas sin reprocesar el OCR.
export default async function DocumentosDeBajaPage({ searchParams }: { searchParams: Promise<{ buscar?: string; pagina?: string; doc?: string; restituido?: string }> }) {
  await exigirPermiso("DOCUMENTO_RESTITUIR");
  const sp = await searchParams;
  const buscar = (sp.buscar ?? "").trim();
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const lista = await listarDocumentosDeBaja(buscar, pagina);
  const elegido = lista.datos.find((d) => d.id === sp.doc) ?? null;

  const url = (doc?: string) => {
    const qs = new URLSearchParams();
    if (buscar) qs.set("buscar", buscar);
    if (pagina > 1) qs.set("pagina", String(pagina));
    if (doc) qs.set("doc", doc);
    return qs.size ? `${RUTA}?${qs}` : RUTA;
  };
  const origen = (d: DocumentoDeBaja) => (d.tramiteCodigo ? `Trámite ${d.tramiteCodigo}` : `Expediente · ${ETIQUETA_TIPO_DOCUMENTO[d.tipoDocumento] ?? "Escaneo"}`);
  const cerrado = (d: DocumentoDeBaja) => d.tramiteEstado === "REMITIDO" || d.tramiteEstado === "ARCHIVADO";

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina
          titulo="Documentos dados de baja"
          descripcion="No aparecen en las búsquedas. Al restituirlos vuelven a estar activos y se pueden buscar de nuevo."
        />
        {sp.restituido && <Alerta tipo="exito">Documento restituido: ya vuelve a aparecer en las búsquedas.</Alerta>}
        <div className="flex flex-wrap items-start gap-6">
          <div className="flex min-w-0 flex-[999_1_560px] flex-col gap-4">
            <BuscadorCarga action={RUTA} defaultValue={buscar} placeholder="Docente, código, cédula, trámite o nombre del archivo" />
            <ZonaCarga>
              <Tabla
                filas={lista.datos}
                clave={(d) => d.id}
                href={(d) => url(d.id)}
                vacio={buscar ? "Ningún documento de baja coincide con la búsqueda." : "No hay documentos dados de baja."}
                columnas={[
                  {
                    titulo: "Documento",
                    celda: (d) => (
                      <div className={d.id === elegido?.id ? "border-l-[3px] border-accent pl-2.5" : ""}>
                        <div className="break-words font-semibold text-ink">{d.nombreOriginal}</div>
                        <div className="text-xs text-ink-faint">{origen(d)}</div>
                      </div>
                    ),
                  },
                  { titulo: "Docente", celda: (d) => d.docenteNombre },
                  { titulo: "Motivo de baja", celda: (d) => <span className="text-ink-soft">{d.motivoBaja}</span> },
                  {
                    titulo: "Baja",
                    className: "whitespace-nowrap",
                    celda: (d) => (
                      <span className="text-sm text-ink-soft">
                        {d.bajaPorNombre || "—"}
                        <br />
                        {d.bajaEn ? formatearFecha(d.bajaEn) : ""}
                      </span>
                    ),
                  },
                ]}
              />
            </ZonaCarga>
            <Paginacion pagina={lista.pagina} totalPaginas={Math.ceil(lista.total / lista.tamano)} searchParams={{ buscar }} />
          </div>

          <aside aria-labelledby="restituir" className="max-w-full flex-[1_1_340px]">
            <Tarjeta>
              <h2 id="restituir" className="mb-4 font-display text-xl font-semibold text-ink">
                Restituir documento
              </h2>
              {!elegido ? (
                <p className="text-sm text-ink-soft">Elige un documento de la lista para ver su detalle y restituirlo.</p>
              ) : (
                <div className="flex flex-col gap-5">
                  <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
                    <dt className="text-ink-faint">Documento</dt>
                    <dd className="break-words font-semibold text-ink">{elegido.nombreOriginal}</dd>
                    <dt className="text-ink-faint">Origen</dt>
                    <dd className="text-ink">
                      {origen(elegido)}
                      {elegido.tramiteEstado && ` · ${ETIQUETA_ESTADO[elegido.tramiteEstado]}`}
                    </dd>
                    <dt className="text-ink-faint">Docente</dt>
                    <dd className="text-ink">
                      {elegido.docenteNombre}
                      {elegido.docenteCodigo && <span className="font-mono text-ink-soft"> · {elegido.docenteCodigo}</span>}
                    </dd>
                    <dt className="text-ink-faint">Dado de baja</dt>
                    <dd className="text-ink">
                      {elegido.bajaEn ? formatearFecha(elegido.bajaEn, true) : "—"}
                      {elegido.bajaPorNombre && ` por ${elegido.bajaPorNombre}`}
                    </dd>
                    <dt className="text-ink-faint">Motivo</dt>
                    <dd className="text-ink">{elegido.motivoBaja}</dd>
                    {elegido.hashSha256 && (
                      <>
                        <dt className="text-ink-faint">Integridad</dt>
                        <dd className="break-all font-mono text-xs text-ink-soft">SHA-256 {elegido.hashSha256.slice(0, 16)}…</dd>
                      </>
                    )}
                  </dl>
                  {cerrado(elegido) ? (
                    <p className="rounded-lg bg-bg px-3.5 py-3 text-sm text-ink-soft">El trámite está cerrado y es de solo lectura: este documento ya no se puede restituir.</p>
                  ) : (
                    <RestituirForm key={elegido.id} docId={elegido.id} cancelar={url()} />
                  )}
                </div>
              )}
            </Tarjeta>
          </aside>
        </div>
      </div>
    </ListadoConCarga>
  );
}
