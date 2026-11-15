import Link from "next/link";
import Alerta from "@/components/common/Alerta";
import BuscadorSugerencias from "@/components/common/BuscadorSugerencias";
import { LinkCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import { exigirPermiso } from "@/lib/api";
import { buscarEnExpediente, listarExpediente, obtenerDocente, type CoincidenciaExpediente } from "@/lib/api-tramites";
import { mensajeError } from "@/lib/formularios";
import { ETIQUETA_TIPO_DOCUMENTO } from "@/lib/tramites";
import VisorDocumento, { Fragmento } from "../../tramites/[id]/VisorDocumento";
import SelectorDocente from "../SelectorDocente";

const RUTA = "/panel/documentos/buscar";

// Buscador en dos pasos (CU-18): 1) el docente, por nombre o código, con sugerencias mientras se escribe;
// 2) el contenido de sus documentos, también con sugerencias (página y fragmento) y el visor en la página exacta.
export default async function BuscarDocumentosPage({
  searchParams,
}: {
  searchParams: Promise<{ docente?: string; q?: string; doc?: string; pagina?: string; tipo?: string; buscar?: string }>;
}) {
  await exigirPermiso("DOCUMENTO_BUSCAR");
  const sp = await searchParams;

  if (!sp.docente) {
    const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
    const buscar = (sp.buscar ?? "").trim();
    return (
      <ListadoConCarga>
        <div className="mx-auto flex w-full max-w-[1100px] flex-col gap-6">
          <div className="flex flex-col gap-1.5">
            <h1 className="font-display text-3xl font-semibold leading-tight tracking-[-0.02em] text-ink">Buscar documentos</h1>
            <Pasos actual={1} />
          </div>
          <div className="flex flex-col gap-2">
            <p className="text-base font-semibold text-ink">Paso 1 · Busca al docente por su nombre o su código de docente</p>
            <BuscadorSugerencias
              accion={RUTA}
              nombre="buscar"
              valorInicial={buscar}
              fuente={`${RUTA}/sugerencias`}
              placeholder="Nombre o código del docente, p. ej. «Gutiérrez» o «1234»"
              etiqueta="Nombre o código del docente"
              autoFocus
            />
            <p className="text-sm text-ink-faint">Elige al docente en la lista de sugerencias; después podrás buscar dentro de sus documentos.</p>
          </div>
          {buscar && (
            <SelectorDocente ruta={RUTA} buscar={buscar} pagina={pagina} etiquetaAccion="Buscar en sus documentos" destino={(id) => `${RUTA}?docente=${id}`} sinBuscador />
          )}
        </div>
      </ListadoConCarga>
    );
  }

  const q = (sp.q ?? "").trim();
  const [docente, docs, busqueda] = await Promise.all([
    obtenerDocente(sp.docente),
    listarExpediente(sp.docente),
    q ? buscarEnExpediente(sp.docente, q) : Promise.resolve(null),
  ]);
  const activos = docs.filter((d) => d.activo);
  const abierto = docs.find((d) => d.id === sp.doc) ?? null;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const pendientes = activos.filter((d) => d.estadoOcr !== "PROCESADO" && d.estadoOcr !== "NO_APLICA").length;

  const url = (p: { tipo?: string; doc?: string; pagina?: number }) => {
    const qs = new URLSearchParams({ docente: docente.id });
    if (q) qs.set("q", q);
    const tipo = "tipo" in p ? p.tipo : sp.tipo;
    if (tipo) qs.set("tipo", tipo);
    if (p.doc) qs.set("doc", p.doc);
    if (p.pagina) qs.set("pagina", String(p.pagina));
    return `${RUTA}?${qs}`;
  };

  // Chips por tipo, con los documentos que tiene el docente de cada uno.
  const tipos = [...new Set(activos.map((d) => d.tipoDocumento))].map((t) => ({ valor: t, etiqueta: ETIQUETA_TIPO_DOCUMENTO[t] ?? t, n: activos.filter((d) => d.tipoDocumento === t).length }));
  const filtradas = (busqueda?.coincidencias ?? []).filter((c) => !sp.tipo || c.tipoDocumento === sp.tipo);
  const porDocumento = new Map<string, CoincidenciaExpediente[]>();
  for (const c of filtradas) porDocumento.set(c.documentoId, [...(porDocumento.get(c.documentoId) ?? []), c]);
  const delAbierto = abierto ? (busqueda?.coincidencias ?? []).filter((c) => c.documentoId === abierto.id).sort((a, b) => a.pagina - b.pagina) : [];
  const verdocs = activos.filter((d) => !sp.tipo || d.tipoDocumento === sp.tipo);

  const chip = "rounded-full border px-3.5 py-1.5 text-sm font-semibold transition-colors";
  const chipActivo = `${chip} border-ink bg-ink text-surface`;
  const chipInactivo = `${chip} border-line bg-surface text-ink hover:border-accent hover:text-accent-ink`;

  const resultados = (
    <div className="flex min-w-0 flex-col gap-4">
      {busqueda?.error && <Alerta tipo="error">{mensajeError(busqueda.error, "No se pudo buscar en el expediente.")}</Alerta>}
      {busqueda && !busqueda.error && (
        <p className="text-sm text-ink-soft">
          <strong className="text-ink">{filtradas.length === 1 ? "1 página" : `${filtradas.length} páginas`}</strong> en {porDocumento.size === 1 ? "1 documento" : `${porDocumento.size} documentos`} para «{q}»
        </p>
      )}
      {busqueda && !busqueda.error && filtradas.length === 0 && (
        <p className="rounded-xl border border-line bg-surface p-5 text-sm text-ink-soft">Prueba con otras palabras. Solo se buscan documentos con el OCR terminado.</p>
      )}

      {[...porDocumento.entries()].map(([docId, lista]) => (
        <article key={docId} className="overflow-hidden rounded-xl border border-line bg-surface">
          <header className="flex flex-wrap items-center gap-2.5 border-b border-line px-5 py-3.5">
            <span className="rounded-full bg-bg px-2.5 py-1 text-xs font-semibold text-ink-soft">{ETIQUETA_TIPO_DOCUMENTO[lista[0].tipoDocumento] ?? lista[0].tipoDocumento}</span>
            <span className="min-w-0 flex-1 truncate text-sm font-semibold text-ink">{lista[0].nombreOriginal}</span>
            <span className="text-xs text-ink-faint">{lista.length === 1 ? "1 página" : `${lista.length} páginas`}</span>
          </header>
          <ul className="divide-y divide-line">
            {lista
              .slice()
              .sort((a, b) => a.pagina - b.pagina)
              .map((c) => {
                const actual = c.documentoId === sp.doc && c.pagina === pagina;
                return (
                  <li key={c.pagina}>
                    <LinkCarga href={url({ doc: c.documentoId, pagina: c.pagina })} scroll={false} className={`flex items-center gap-4 px-5 py-3.5 hover:bg-bg ${actual ? "bg-accent-soft" : ""}`}>
                      <span className={`min-w-14 shrink-0 rounded-lg px-2 py-1.5 text-center text-xs font-semibold ${actual ? "bg-accent text-white" : "bg-bg text-ink-soft"}`}>pág. {c.pagina}</span>
                      <span className="min-w-0 flex-1 text-sm leading-relaxed text-ink-soft">
                        <Fragmento texto={c.fragmento} />
                      </span>
                      <svg className="shrink-0 text-ink-faint" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                        <path d="M5 12h14" />
                        <path d="m13 6 6 6-6 6" />
                      </svg>
                    </LinkCarga>
                  </li>
                );
              })}
          </ul>
        </article>
      ))}

      {!q && (
        <section className="rounded-xl border border-line bg-surface p-5">
          <h2 className="text-base font-semibold text-ink">Documentos del expediente</h2>
          <p className="mb-3 text-sm text-ink-soft">Escribe arriba para buscar por contenido, o abre uno directamente.</p>
          {verdocs.length === 0 ? (
            <p className="text-sm text-ink-faint">Archivo aún no subió escaneos de este docente.</p>
          ) : (
            <ul className="divide-y divide-line">
              {verdocs.map((d) => {
                const contenido = (
                  <span className="flex min-w-0 items-center gap-2.5">
                    <span className="rounded-full bg-bg px-2.5 py-1 text-xs font-semibold text-ink-soft">{ETIQUETA_TIPO_DOCUMENTO[d.tipoDocumento] ?? d.tipoDocumento}</span>
                    <span className="min-w-0 break-words">{d.nombreOriginal}</span>
                  </span>
                );
                const fila = "-mx-3 flex flex-wrap items-center justify-between gap-2 rounded-xl px-3 py-2.5 text-sm";
                return (
                  <li key={d.id}>
                    {d.estadoOcr === "SUBIENDO" ? (
                      <div className={fila}>
                        {contenido}
                        <span className="text-xs text-ink-faint">Subida incompleta</span>
                      </div>
                    ) : (
                      // Toda la fila abre el documento, no solo «Abrir».
                      <LinkCarga href={url({ doc: d.id })} scroll={false} className={`${fila} group text-ink transition-colors hover:bg-bg`}>
                        {contenido}
                        <span className="text-sm font-semibold text-accent-ink group-hover:underline">Abrir</span>
                      </LinkCarga>
                    )}
                  </li>
                );
              })}
            </ul>
          )}
        </section>
      )}
    </div>
  );

  return (
    <ListadoConCarga>
    <div className="mx-auto flex w-full max-w-[1320px] flex-col gap-7">
      <div>
        <Link href={RUTA} className="text-sm font-semibold text-ink-soft hover:text-accent-ink">
          ← Cambiar de docente
        </Link>
        <Pasos actual={2} />
        <div className="mt-1 flex flex-wrap items-baseline gap-x-4 gap-y-1">
          <h1 className="font-display text-3xl font-semibold leading-tight tracking-[-0.02em] text-ink">
            {docente.apellidos}, {docente.nombres}
          </h1>
          <p className="text-sm text-ink-soft">
            Código <span className="font-mono text-ink">{docente.codigo || "—"}</span>
            {docente.categoria && ` · ${docente.categoria}`} · {activos.length} {activos.length === 1 ? "documento" : "documentos"}
          </p>
        </div>
      </div>

      <section className="flex flex-col gap-4">
        <p className="-mb-2 text-base font-semibold text-ink">Paso 2 · Busca en sus documentos</p>
        <BuscadorSugerencias
          accion={RUTA}
          ocultos={{ docente: docente.id, ...(sp.tipo ? { tipo: sp.tipo } : {}) }}
          valorInicial={q}
          fuente={`${RUTA}/coincidencias?docente=${docente.id}`}
          placeholder="Busca por contenido: resolución, título, provisión nacional…"
          etiqueta="Buscar en los documentos del docente"
          autoFocus={!q && !abierto}
        />

        {tipos.length > 1 && (
          <div className="flex flex-wrap items-center gap-2">
            <LinkCarga href={url({ tipo: "" })} scroll={false} className={!sp.tipo ? chipActivo : chipInactivo}>
              Todos · {activos.length}
            </LinkCarga>
            {tipos.map((t) => (
              <LinkCarga key={t.valor} href={url({ tipo: t.valor })} scroll={false} className={sp.tipo === t.valor ? chipActivo : chipInactivo}>
                {t.etiqueta} · {t.n}
              </LinkCarga>
            ))}
          </div>
        )}
        {pendientes > 0 && (
          <p className="text-xs text-ink-faint">
            {pendientes} {pendientes === 1 ? "documento aún se está procesando" : "documentos aún se están procesando"}: sus páginas se podrán buscar al terminar.
          </p>
        )}
      </section>

      <ZonaCarga>
        {abierto ? (
          <VisorDocumento
            base={`/api/docentes/${docente.id}`}
            archivo={`/panel/documentos/${docente.id}/archivo/${abierto.id}?v=${abierto.tamanoComprimido}`}
            documento={abierto}
            etiqueta={ETIQUETA_TIPO_DOCUMENTO[abierto.tipoDocumento] ?? "Documento"}
            subtitulo={`${docente.apellidos}, ${docente.nombres}`}
            cerrar={url({ doc: "" })}
            consulta={q}
            pagina={pagina}
            coincidencias={q ? delAbierto.map((c) => ({ pagina: c.pagina, fragmento: c.fragmento })) : undefined}
          />
        ) : (
          resultados
        )}
      </ZonaCarga>
    </div>
    </ListadoConCarga>
  );
}

// Indicador de en qué paso del buscador se está.
function Pasos({ actual }: { actual: 1 | 2 }) {
  const paso = (n: 1 | 2, texto: string) => (
    <li className={`flex items-center gap-2 ${n === actual ? "font-semibold text-ink" : "text-ink-faint"}`} aria-current={n === actual ? "step" : undefined}>
      <span className={`flex h-6 w-6 items-center justify-center rounded-full text-xs font-semibold ${n === actual ? "bg-accent text-white" : n < actual ? "bg-ink text-surface" : "border border-line"}`}>
        {n < actual ? "✓" : n}
      </span>
      {texto}
    </li>
  );
  return (
    <ol aria-label="Pasos de la búsqueda" className="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm">
      {paso(1, "Docente")}
      <li aria-hidden="true" className="h-px w-8 bg-line" />
      {paso(2, "Contenido de sus documentos")}
    </ol>
  );
}
