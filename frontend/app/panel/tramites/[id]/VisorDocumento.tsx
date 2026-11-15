import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import { LinkCarga } from "@/components/common/CargaListado";
import { AlVisor } from "@/components/common/EnlaceCarga";
import VisorConBusqueda, { type Coincidencia } from "@/components/documentos/VisorConBusqueda";
import { api } from "@/lib/api";
import { type DocumentoTramite } from "@/lib/api-tramites";
import { leerErrorApi, mensajeError } from "@/lib/formularios";
import { formatearTamano } from "@/lib/tramites";

export { Fragmento } from "@/components/documentos/VisorConBusqueda";

// Visor del documento abierto: lector propio (pdf.js) con barra de herramientas, miniaturas y resaltado de la
// búsqueda, más un panel con las coincidencias y los detalles. El archivo lo entrega el servidor de la
// aplicación (`archivo`), así que no depende del visor del navegador ni de permisos del bucket.
export default async function VisorDocumento({
  base,
  archivo,
  rutaBuscar,
  documento,
  etiqueta,
  subtitulo,
  cerrar,
  consulta = "",
  pagina = 1,
  coincidencias,
}: {
  base: string; // prefijo de la API del dueño del documento: /api/tramites/<id> o /api/docentes/<id>
  archivo: string; // ruta del servidor que entrega el archivo: /panel/.../archivo/<docId>
  rutaBuscar?: string; // ruta que busca dentro del documento (?q=); sin ella no hay buscador interno
  documento: DocumentoTramite;
  etiqueta: string;
  subtitulo?: string;
  cerrar: string;
  consulta?: string; // ?q= búsqueda dentro del texto OCR
  pagina?: number; // ?pagina= página a mostrar
  coincidencias?: Coincidencia[]; // las páginas halladas, si quien lo usa ya buscó
}) {
  // Enlace firmado solo para «Descargar» (el visor usa `archivo`).
  const respuesta = await api(`${base}/documentos/${documento.id}`);
  const urlDescarga: string | null = respuesta.ok ? (await respuesta.json()).url : null;
  const error = urlDescarga ? null : mensajeError((await leerErrorApi(respuesta)).codigo, "No se pudo abrir el documento.");
  const esImagen = documento.tipoMime.startsWith("image/") && !documento.tamanoComprimido;

  // Búsqueda inicial (al llegar con ?q= en la URL); las siguientes las hace el visor desde el navegador.
  let halladas: Coincidencia[] = coincidencias ?? [];
  let errorBusqueda: string | null = null;
  if (!coincidencias && rutaBuscar && documento.estadoOcr === "PROCESADO" && consulta) {
    const r = await api(`${base}/documentos/${documento.id}/buscar?q=${encodeURIComponent(consulta)}`);
    if (r.ok) halladas = (await r.json()).coincidencias;
    else errorBusqueda = mensajeError((await leerErrorApi(r)).codigo, "No se pudo buscar en el documento.");
  }

  return (
    <section className="flex flex-col gap-5 rounded-xl border border-line bg-surface p-5 sm:p-6">
      <AlVisor clave={documento.id} />
      <header className="flex flex-wrap items-center gap-x-5 gap-y-3">
        <LinkCarga
          href={cerrar}
          scroll={false}
          aria-label="Cerrar visor"
          className="flex h-10 w-10 items-center justify-center rounded-xl border border-line text-ink hover:border-accent hover:text-accent-ink"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
            <path d="M19 12H5" />
            <path d="m11 18-6-6 6-6" />
          </svg>
        </LinkCarga>
        <div className="min-w-0 flex-1 basis-64">
          <div className="flex flex-wrap items-center gap-x-2.5 gap-y-1">
            <span className="rounded-full bg-bg px-2.5 py-1 text-xs font-semibold text-ink-soft">{etiqueta}</span>
            {subtitulo && <span className="text-sm text-ink-soft">{subtitulo}</span>}
          </div>
          <h2 className="mt-1 break-words font-display text-xl font-semibold leading-tight tracking-[-0.01em] text-ink">{documento.nombreOriginal}</h2>
        </div>
        {urlDescarga && (
          <a href={urlDescarga} target="_blank" rel="noreferrer" className={claseBoton("secundario", "sm")}>
            Descargar
          </a>
        )}
      </header>

      {error && <Alerta tipo="error">{error}</Alerta>}

      {errorBusqueda && <Alerta tipo="error">{errorBusqueda}</Alerta>}

      {/* Buscar y saltar entre coincidencias ocurre en el navegador (VisorConBusqueda), sin volver al servidor. */}
      <VisorConBusqueda
        key={documento.id}
        src={archivo}
        esImagen={esImagen}
        nombre={documento.nombreOriginal}
        paginaInicial={pagina}
        consultaInicial={consulta}
        coincidenciasIniciales={halladas}
        rutaBuscar={rutaBuscar && documento.estadoOcr === "PROCESADO" ? rutaBuscar : undefined}
      >
        <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2.5 rounded-xl border border-line p-4 text-sm">
          <dt className="text-ink-faint">Tamaño</dt>
          <dd className="text-right font-semibold text-ink">
            {formatearTamano(documento.tamano)}
            {documento.tamanoComprimido > 0 && ` → ${formatearTamano(documento.tamanoComprimido)}`}
          </dd>
          <dt className="text-ink-faint">Texto (OCR)</dt>
          <dd className={`text-right font-semibold ${documento.estadoOcr === "PROCESADO" ? "text-success" : "text-ink"}`}>
            {documento.estadoOcr === "PROCESADO" ? "Listo para buscar" : documento.estadoOcr === "NO_APLICA" ? "No aplica" : "En proceso"}
          </dd>
        </dl>
      </VisorConBusqueda>
    </section>
  );
}
