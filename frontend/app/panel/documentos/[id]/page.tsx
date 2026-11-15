import Link from "next/link";
import Insignia from "@/components/common/Insignia";
import Tarjeta from "@/components/common/Tarjeta";
import { EnlaceSinScroll } from "@/components/common/EnlaceCarga";
import { exigirPermiso } from "@/lib/api";
import { listarExpediente, obtenerDocente, type DocumentoTramite } from "@/lib/api-tramites";
import { formatearFecha, formatearTamano, TIPOS_DOCUMENTO } from "@/lib/tramites";
import { AutoRefrescar } from "../../tramites/[id]/CargaDocumentos";
import VisorDocumento from "../../tramites/[id]/VisorDocumento";
import { BajaEscaneoForm, ReintentarEscaneo, SubirEscaneoForm } from "./SubirEscaneo";

const ESTADO_OCR: Record<string, { texto: string; tono: "acento" | "exito" | "neutro" }> = {
  SUBIENDO: { texto: "Subida incompleta", tono: "neutro" },
  PENDIENTE: { texto: "En cola", tono: "acento" },
  PROCESANDO: { texto: "Procesando OCR", tono: "acento" },
  PROCESADO: { texto: "Listo para buscar", tono: "exito" },
  ERROR: { texto: "Error de OCR", tono: "acento" },
};

// Expediente digital de un docente: Archivo sube los escaneos de su carpeta física, uno por documento y tipo.
export default async function ExpedienteDocentePage({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ doc?: string }>;
}) {
  const { id } = await params;
  const sp = await searchParams;
  await exigirPermiso("DOCUMENTO_ESCANEAR");
  const [docente, docs] = await Promise.all([obtenerDocente(id), listarExpediente(id)]);
  const enProceso = docs.some((d) => d.activo && (d.estadoOcr === "PENDIENTE" || d.estadoOcr === "PROCESANDO"));
  const abierto = docs.find((d) => d.id === sp.doc) ?? null;
  const porTipo = (tipo: string) => docs.filter((d) => d.tipoDocumento === tipo);
  const verDocumento = (d: DocumentoTramite) => `/panel/documentos/${id}?doc=${d.id}`;

  return (
    <div className="mx-auto flex w-full max-w-[1200px] flex-col gap-6">
      <div>
        <Link href="/panel/documentos" className="text-sm font-semibold text-ink-soft hover:text-accent-ink">
          ← Docentes
        </Link>
        <h1 className="mt-1 font-display text-3xl font-semibold leading-tight tracking-[-0.02em] text-ink">
          {docente.apellidos}, {docente.nombres}
        </h1>
        <p className="text-sm text-ink-soft">
          Código <span className="font-mono">{docente.codigo || "—"}</span>
          {docente.categoria && ` · ${docente.categoria}`}
          {docente.legajos.length > 0 && ` · Gaveta ${docente.legajos.map((l) => l.gaveta || "s/g").join(" · ")}`}
        </p>
      </div>

      <SubirEscaneoForm docenteId={id} />

      <Tarjeta titulo="Documentos escaneados" descripcion="El texto de cada página se extrae con OCR; los analistas lo buscan por contenido.">
        <AutoRefrescar activo={enProceso} />
        {docs.length === 0 ? (
          <p className="text-sm text-ink-faint">Este docente aún no tiene escaneos. Sube el primero arriba.</p>
        ) : (
          <div className="flex flex-col gap-6">
            {TIPOS_DOCUMENTO.map(({ valor, etiqueta }) => {
              const lista = porTipo(valor);
              if (lista.length === 0) return null;
              return (
                <section key={valor}>
                  <h3 className="mb-1 text-xs font-semibold uppercase tracking-wide text-ink-faint">{etiqueta}</h3>
                  <ul className="flex flex-col divide-y divide-line">
                    {lista.map((d) => {
                      const estado = ESTADO_OCR[d.estadoOcr];
                      return (
                        <li key={d.id} className="flex flex-col gap-2 py-3 text-sm">
                          <div className="flex flex-wrap items-center justify-between gap-2">
                            <span className={`min-w-0 break-words text-ink ${d.activo ? "" : "line-through"}`}>{d.nombreOriginal}</span>
                            {!d.activo ? <Insignia tono="neutro">De baja</Insignia> : estado && <Insignia tono={estado.tono}>{estado.texto}</Insignia>}
                          </div>
                          <p className="text-xs text-ink-faint">
                            {formatearTamano(d.tamano)}
                            {d.tamanoComprimido > 0 && ` → ${formatearTamano(d.tamanoComprimido)} comprimido`}
                            {!d.activo && d.motivoBaja && ` · ${d.motivoBaja}`}
                            {!d.activo && d.bajaEn && ` · ${formatearFecha(d.bajaEn)}`}
                          </p>
                          {d.activo && d.estadoOcr === "SUBIENDO" && (
                            <p className="text-xs text-ink-soft">El archivo no llegó al almacén. Dalo de baja y vuelve a subirlo.</p>
                          )}
                          {d.activo && d.estadoOcr === "ERROR" && (
                            <>
                              <p className="text-xs text-ink-soft">{d.errorOcr || "El procesamiento falló."}</p>
                              <ReintentarEscaneo docenteId={id} docId={d.id} />
                            </>
                          )}
                          <div className="flex flex-wrap items-start gap-x-5 gap-y-2">
                            {d.estadoOcr !== "SUBIENDO" && (
                              <EnlaceSinScroll href={verDocumento(d)} className="w-fit text-sm font-semibold text-accent-ink hover:underline">
                                Ver documento
                              </EnlaceSinScroll>
                            )}
                            {d.activo && <BajaEscaneoForm docenteId={id} docId={d.id} />}
                          </div>
                        </li>
                      );
                    })}
                  </ul>
                </section>
              );
            })}
          </div>
        )}
      </Tarjeta>

      {abierto && (
        <VisorDocumento
          base={`/api/docentes/${id}`}
          archivo={`/panel/documentos/${id}/archivo/${abierto.id}?v=${abierto.tamanoComprimido}`}
          subtitulo={`${docente.apellidos}, ${docente.nombres}`}
          documento={abierto}
          etiqueta={TIPOS_DOCUMENTO.find((t) => t.valor === abierto.tipoDocumento)?.etiqueta ?? "Documento"}
          cerrar={`/panel/documentos/${id}`}
        />
      )}
    </div>
  );
}
