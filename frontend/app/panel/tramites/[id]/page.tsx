import Link from "next/link";
import { notFound } from "next/navigation";
import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import BotonEnviar from "@/components/common/BotonEnviar";
import Insignia from "@/components/common/Insignia";
import Tarjeta from "@/components/common/Tarjeta";
import EstadoInsignia from "@/components/tramites/EstadoInsignia";
import FasesTramite from "@/components/tramites/FasesTramite";
import { obtenerSesion } from "@/lib/api";
import { listarAnalistas, obtenerTramite } from "@/lib/api-tramites";
import { mensajeError } from "@/lib/formularios";
import { fechaManana, formatearFecha, formatearTamano, tipoYCategoria } from "@/lib/tramites";
import { iniciarRevisionAction, remitirFacultadAction, resolverObservacionAction, validarAction } from "../actions";
import { EnlaceSinScroll } from "@/components/common/EnlaceCarga";
import { AutoRefrescar, ReintentarForm, TextoOcr } from "./CargaDocumentos";
import { ArchivarForm, AsignarForm, BajaDocumentoForm, DictamenForm, ObservacionForm, PeriodoForm, RegresoForm, RemisionForm } from "./FormulariosTramite";
import VisorDocumento from "./VisorDocumento";

const ESTADO_OCR: Record<string, { texto: string; tono: "acento" | "exito" | "neutro" }> = {
  SUBIENDO: { texto: "Subida incompleta", tono: "neutro" },
  PENDIENTE: { texto: "En cola", tono: "acento" },
  PROCESANDO: { texto: "Procesando OCR", tono: "acento" },
  PROCESADO: { texto: "Texto extraído", tono: "exito" },
  ERROR: { texto: "Error de OCR", tono: "acento" },
};

const MENSAJES_EXITO: Record<string, string> = {
  RADICADO: "Trámite radicado. Entrega el código de seguimiento al docente.",
  VALIDADO: "Solicitud validada: ya se puede asignar.",
  ASIGNADO: "Trámite asignado.",
  REVISION: "Revisión iniciada.",
  OBSERVADO: "Observación registrada: el trámite pasó a Observado.",
  RESUELTA: "Observación marcada como resuelta.",
  DICTAMINADO: "Dictamen emitido.",
  PERIODO: "Años a evaluar registrados.",
  FACULTAD: "Expediente remitido a la facultad. Desde ahora no se puede rastrear hasta que regrese.",
  REGRESO: "Regreso registrado: el expediente volvió con la resolución y está en revisión final.",
  RECTORADO: "Expediente remitido al Rectorado. Falta que Archivo archive la carpeta.",
  ARCHIVADO: "Trámite archivado: ahora es de solo lectura.",
  DOCUMENTO_BAJA: "Documento dado de baja: ya no cubre el requisito ni aparece en las búsquedas.",
};

function Dato({ etiqueta, children }: { etiqueta: string; children: React.ReactNode }) {
  return (
    <div className="flex flex-col gap-0.5">
      <dt className="text-xs font-semibold uppercase tracking-wide text-ink-faint">{etiqueta}</dt>
      <dd className="text-sm text-ink">{children}</dd>
    </div>
  );
}

export default async function TramitePage({
  params,
  searchParams,
}: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ error?: string; exito?: string; doc?: string; q?: string; pagina?: string }>;
}) {
  const { id } = await params;
  const sp = await searchParams;
  const sesion = await obtenerSesion();
  const puede = (p: string) => sesion.permisos.includes(p);
  if (!["TRAMITE_RADICAR", "TRAMITE_ASIGNAR", "TRAMITE_REVISAR", "TRAMITE_ARCHIVAR"].some(puede)) notFound();

  const t = await obtenerTramite(id);
  const esMio = t.analistaId === sesion.usuario.id;
  const archivado = t.estado === "ARCHIVADO";
  const cerrado = archivado || t.estado === "REMITIDO"; // remitido o archivado: solo lectura para el analista
  const puedeAsignar =
    puede("TRAMITE_ASIGNAR") &&
    !t.pendienteValidacion &&
    ["INGRESADO", "ASIGNADO", "EN_REVISION", "OBSERVADO", "DICTAMINADO", "EN_FACULTAD", "REVISION_FINAL"].includes(t.estado);
  const analistas = puedeAsignar ? await listarAnalistas() : [];
  const pendientes = t.observaciones.filter((o) => !o.resuelta);
  const etapas = t.historial.filter((h) => h.estado !== h.estadoAnterior).map((h) => ({ estado: h.estado, fecha: h.fecha }));
  const revisando = esMio && (t.estado === "EN_REVISION" || t.estado === "OBSERVADO");
  // Visor: ?doc=<id> abre el documento debajo, en lugar del historial.
  const abierto = t.documentos.find((d) => d.id === sp.doc) ?? null;
  const verDocumento = (docId: string) => `/panel/tramites/${t.id}?doc=${docId}`;
  // CU-15: documentos cargados por el analista (con OCR); los del portal no pasan por el OCR.
  const digitales = t.documentos.filter((d) => d.activo && d.estadoOcr !== "NO_APLICA");
  const enProceso = digitales.some((d) => d.estadoOcr === "PENDIENTE" || d.estadoOcr === "PROCESANDO");
  const puedeCorregir = puede("DOCUMENTO_GESTIONAR") && esMio && !cerrado;
  const hoy = new Date().toISOString().slice(0, 10);
  const hayResumen = t.anioDesde > 0 || t.resolucionNumero || t.remisionNota;

  return (
    <div className="mx-auto flex w-full max-w-[1200px] flex-col gap-6">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <div>
          <Link href="/panel/tramites" className="text-sm font-semibold text-ink-soft hover:text-accent-ink">
            ← Trámites
          </Link>
          <h1 className="mt-1 flex flex-wrap items-center gap-3 font-display text-3xl font-semibold leading-tight tracking-[-0.02em] text-ink">
            <span className="font-mono">{t.codigo}</span>
            <EstadoInsignia estado={t.estado} />
            {t.pendienteValidacion && <Insignia tono="acento">Por validar</Insignia>}
          </h1>
          <p className="text-sm text-ink-soft">
            {tipoYCategoria(t)} · radicado el {formatearFecha(t.creadoEn, true)}
          </p>
        </div>
      </div>

      {sp.error && <Alerta tipo="error">{mensajeError(sp.error, "No se pudo completar la acción.")}</Alerta>}
      {sp.exito && MENSAJES_EXITO[sp.exito] && <Alerta tipo="exito">{MENSAJES_EXITO[sp.exito]}</Alerta>}

      {/* Fila superior: resumen (izquierda, con su altura natural) y acciones (derecha; si son más cortas que el resumen, la última tarjeta se estira para igualarlo). El resto del expediente ocupa todo el ancho. */}
      <div className="grid grid-cols-1 gap-5 lg:grid-cols-[minmax(0,1fr)_380px] lg:[&:has(>:last-child:empty)]:grid-cols-1">
        <div className="flex min-w-0 flex-col gap-5 lg:col-start-1 lg:row-start-1 lg:self-start">
          <Tarjeta>
            <dl className="mb-7 grid grid-cols-2 gap-5 sm:grid-cols-4">
              <Dato etiqueta="Docente">{t.docenteNombre}</Dato>
              <Dato etiqueta="Código">{t.docenteCodigo}</Dato>
              <Dato etiqueta="Analista">{t.analistaNombre || <span className="text-ink-faint">Sin asignar</span>}</Dato>
              <Dato etiqueta="Origen">{t.origen === "PORTAL" ? "Portal público" : "Ventanilla"}</Dato>
            </dl>
            <FasesTramite estado={t.estado} etapas={etapas} />
          </Tarjeta>

          {hayResumen && (
            <Tarjeta titulo="Evaluación y resolución">
              <dl className="grid grid-cols-1 gap-4 sm:grid-cols-3">
                {t.anioDesde > 0 && (
                  <Dato etiqueta="Años evaluados">
                    {t.anioDesde} a {t.anioHasta}
                  </Dato>
                )}
                {t.resolucionNumero && (
                  <Dato etiqueta="Resolución de la facultad">
                    {t.resolucionNumero}
                    {t.resolucionFecha && <span className="block text-xs text-ink-faint">{formatearFecha(t.resolucionFecha)}</span>}
                  </Dato>
                )}
                {t.remisionNota && <Dato etiqueta="Remisión al Rectorado">{t.remisionNota}</Dato>}
              </dl>
            </Tarjeta>
          )}
          {/* Años a evaluar bajo el resumen: rellena la columna izquierda en lugar de alargar la de acciones. */}
          {revisando && <PeriodoForm key={`${t.anioDesde}-${t.anioHasta}`} id={t.id} desde={t.anioDesde} hasta={t.anioHasta} />}
        </div>

        <div className="flex min-w-0 flex-col gap-5 lg:col-span-full lg:row-start-2">
          <Tarjeta titulo="Documentación exigida" descripcion={t.categoria ? `Catálogo normativo de la categoría ${t.categoria}.` : "Una regularización no postula a ninguna categoría."}>
            {t.requisitos.length === 0 ? (
              <p className="text-sm text-ink-faint">{t.categoria ? "La categoría no exige documentos." : "No hay documentos exigidos por catálogo."}</p>
            ) : (
              <ul className="flex flex-col divide-y divide-line">
                {t.requisitos.map((r) => {
                  const doc = t.documentos.find((d) => d.requisitoId === r.id && d.activo && d.estadoOcr !== "SUBIENDO");
                  const bajas = t.documentos.filter((d) => d.requisitoId === r.id && !d.activo);
                  return (
                    <li key={r.id} className="flex flex-col gap-1.5 py-2.5 text-sm">
                      <div className="flex items-center justify-between gap-3">
                        <span className="text-ink">{r.nombre}</span>
                        {doc ? (
                          <EnlaceSinScroll href={verDocumento(doc.id)} className="text-sm font-semibold text-accent-ink hover:underline">
                            {doc.nombreOriginal}
                          </EnlaceSinScroll>
                        ) : (
                          <span className="text-xs text-ink-faint">{t.origen === "PORTAL" ? "No adjuntado" : "Se verifica en el legajo físico"}</span>
                        )}
                      </div>
                      {/* CU-16: el documento dado de baja se conserva, con su motivo, para que el Director pueda restituirlo. */}
                      {bajas.map((b) => (
                        <div key={b.id} className="flex flex-wrap items-center gap-2 text-xs text-ink-faint">
                          <Insignia tono="neutro">De baja</Insignia>
                          <span className="line-through">{b.nombreOriginal}</span>
                          <span>
                            {b.motivoBaja}
                            {b.bajaEn && ` · ${formatearFecha(b.bajaEn)}`}
                          </span>
                        </div>
                      ))}
                    </li>
                  );
                })}
              </ul>
            )}
          </Tarjeta>

          {puede("DOCUMENTO_BUSCAR") && (
            <Tarjeta
              titulo="Expediente digital del docente"
              descripcion="Escaneos de su carpeta física: busca por contenido (resolución, título, provisión…) y abre la página exacta."
            >
              <Link href={`/panel/documentos/buscar?docente=${t.docenteId}`} className={claseBoton("secundario", "sm")}>
                Buscar en el expediente
              </Link>
            </Tarjeta>
          )}

          {digitales.length > 0 && (
            <Tarjeta titulo="Documentos digitales" descripcion="Cargados por el analista. El texto extraído por OCR se puede revisar y corregir.">
              <AutoRefrescar activo={enProceso} />
              <ul className="flex flex-col divide-y divide-line">
                {digitales.map((d) => {
                  const estado = ESTADO_OCR[d.estadoOcr];
                  return (
                    <li key={d.id} className="flex flex-col gap-2 py-3 text-sm">
                      <div className="flex flex-wrap items-center justify-between gap-2">
                        <span className="min-w-0 break-words text-ink">
                          {t.requisitos.find((r) => r.id === d.requisitoId)?.nombre ?? "Documento adicional"} — {d.nombreOriginal}
                        </span>
                        {estado && <Insignia tono={estado.tono}>{estado.texto}</Insignia>}
                      </div>
                      <p className="text-xs text-ink-faint">
                        {formatearTamano(d.tamano)}
                        {d.tamanoComprimido > 0 && ` → ${formatearTamano(d.tamanoComprimido)} comprimido`}
                        {d.textoValidado && " · texto validado"}
                      </p>
                      {d.estadoOcr === "SUBIENDO" && (
                        <p className="text-xs text-ink-soft">El archivo no llegó al almacén. Dalo de baja y vuelve a cargarlo.</p>
                      )}
                      {d.estadoOcr === "ERROR" && (
                        <>
                          <p className="text-xs text-ink-soft">{d.errorOcr || "El procesamiento falló."}</p>
                          {puedeCorregir && <ReintentarForm id={t.id} docId={d.id} />}
                        </>
                      )}
                      {d.estadoOcr !== "SUBIENDO" && (
                        <EnlaceSinScroll href={verDocumento(d.id)} className="w-fit text-sm font-semibold text-accent-ink hover:underline">
                          Ver documento
                        </EnlaceSinScroll>
                      )}
                      {d.estadoOcr === "PROCESADO" && <TextoOcr id={t.id} docId={d.id} editable={puedeCorregir} />}
                    </li>
                  );
                })}
              </ul>
            </Tarjeta>
          )}

          {puedeCorregir && digitales.length > 0 && ( // baja de adjuntos del portal cargados por error (CU-16)
            <BajaDocumentoForm
              id={t.id}
              dictaminado={!!t.dictamen}
              documentos={digitales.map((d) => ({ id: d.id, etiqueta: d.nombreOriginal }))}
            />
          )}

          {t.observaciones.length > 0 && (
            <Tarjeta titulo="Observaciones" descripcion={`${pendientes.length} pendientes de ${t.observaciones.length}.`}>
              <ul className="flex flex-col gap-3">
                {t.observaciones.map((o) => (
                  <li key={o.id} className="flex items-start justify-between gap-4 rounded-xl border border-line p-4">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2 text-sm font-semibold text-ink">
                        {o.documento}
                        <Insignia tono={o.resuelta ? "exito" : "acento"}>{o.resuelta ? "Resuelta" : "Pendiente"}</Insignia>
                      </div>
                      <p className="mt-1 text-sm text-ink-soft">{o.motivo}</p>
                      <p className="mt-1 text-xs text-ink-faint">Plazo: {formatearFecha(o.plazo)}</p>
                    </div>
                    {revisando && !o.resuelta && (
                      <form action={resolverObservacionAction}>
                        <input type="hidden" name="id" value={t.id} />
                        <input type="hidden" name="observacionId" value={o.id} />
                        <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Guardando…">
                          Marcar resuelta
                        </BotonEnviar>
                      </form>
                    )}
                  </li>
                ))}
              </ul>
            </Tarjeta>
          )}

          {t.dictamen && (
            <Tarjeta titulo="Dictamen" descripcion={t.dictaminadoEn ? `Emitido el ${formatearFecha(t.dictaminadoEn, true)}` : undefined}>
              <Insignia tono={t.resultado === "FAVORABLE" ? "exito" : "acento"}>{t.resultado === "FAVORABLE" ? "Favorable" : "Desfavorable"}</Insignia>
              <p className="mt-3 whitespace-pre-line text-sm leading-relaxed text-ink">{t.dictamen}</p>
              {archivado && (
                <p className="mt-4 border-t border-line pt-3 text-sm text-ink-soft">
                  Legajo archivado en <strong className="text-ink">{t.ubicacionLegajo}</strong>
                  {t.archivadoEn && ` el ${formatearFecha(t.archivadoEn)}`}.
                </p>
              )}
            </Tarjeta>
          )}

          {abierto && (
            <VisorDocumento
              base={`/api/tramites/${t.id}`}
              archivo={`/panel/tramites/${t.id}/archivo/${abierto.id}?v=${abierto.tamanoComprimido}`}
              subtitulo={`${t.docenteNombre} · ${t.codigo}`}
              rutaBuscar={`/panel/tramites/${t.id}/documentos/${abierto.id}/buscar`}
              documento={abierto}
              etiqueta={t.requisitos.find((r) => r.id === abierto.requisitoId)?.nombre ?? "Documento adicional"}
              cerrar={`/panel/tramites/${t.id}`}
              consulta={sp.q ?? ""}
              pagina={Number(sp.pagina) || 1}
            />
          )}
        </div>

        <div className="flex min-w-0 flex-col gap-5 empty:hidden lg:col-start-2 lg:row-start-1 [&>:last-child]:flex-1">
          {archivado && <Alerta tipo="aviso">Trámite archivado: es de solo lectura.</Alerta>}
          {t.estado === "REMITIDO" && !puede("TRAMITE_ARCHIVAR") && (
            <Alerta tipo="aviso">Remitido al Rectorado: es de solo lectura. Archivo archivará la carpeta.</Alerta>
          )}
          {t.estado === "EN_FACULTAD" && (
            <Alerta tipo="aviso">
              Fuera de Escalafón: el expediente está en la facultad (Tribunal, Decano y Consejo Facultativo) y no se puede rastrear. Registra su regreso cuando vuelva con la resolución.
            </Alerta>
          )}

          {t.pendienteValidacion && puede("TRAMITE_RADICAR") && (
            <Tarjeta titulo="Validar solicitud" descripcion="Llegó por el portal público. Confirma en ventanilla que los documentos coinciden con la persona.">
              <form action={validarAction}>
                <input type="hidden" name="id" value={t.id} />
                <BotonEnviar tamano="sm" textoEnviando="Validando…">
                  Validar solicitud
                </BotonEnviar>
              </form>
            </Tarjeta>
          )}

          {puedeAsignar && <AsignarForm id={t.id} analistas={analistas} analistaActualId={t.analistaId} />}

          {esMio && t.estado === "ASIGNADO" && t.carpetasLigadas === 0 && (
            <Tarjeta titulo="Solicita la carpeta" descripcion="La revisión parte de la carpeta física del docente: pídela al archivo antes de iniciar.">
              <Link href={`/panel/carpetas/nueva?tramite=${t.id}`} className={claseBoton("primario", "sm")}>
                Solicitar carpeta
              </Link>
            </Tarjeta>
          )}
          {esMio && t.estado === "ASIGNADO" && t.carpetasLigadas > 0 && (
            <Tarjeta titulo="Iniciar revisión" descripcion="Pasa el trámite a En revisión para contrastar la documentación.">
              <form action={iniciarRevisionAction}>
                <input type="hidden" name="id" value={t.id} />
                <BotonEnviar tamano="sm" textoEnviando="Iniciando…">
                  Iniciar revisión
                </BotonEnviar>
              </form>
            </Tarjeta>
          )}

          {revisando && <ObservacionForm key={t.observaciones.length} id={t.id} documentos={t.requisitos.map((r) => r.nombre)} plazoMinimo={fechaManana()} />}
          {esMio && t.estado === "EN_REVISION" &&
            (t.anioDesde > 0 ? <DictamenForm id={t.id} /> : <Alerta tipo="aviso">Registra los años a evaluar para poder emitir el dictamen.</Alerta>)}
          {esMio && t.estado === "DICTAMINADO" && (
            <Tarjeta
              titulo="Remitir a la facultad"
              descripcion="El expediente sale de Escalafón hacia el Tribunal, el Decano y el Consejo Facultativo. En esa etapa no se puede rastrear: registra el regreso cuando vuelva con la resolución."
            >
              <form action={remitirFacultadAction}>
                <input type="hidden" name="id" value={t.id} />
                <BotonEnviar tamano="sm" textoEnviando="Remitiendo…">
                  Remitir a la facultad
                </BotonEnviar>
              </form>
            </Tarjeta>
          )}
          {t.estado === "EN_FACULTAD" && (esMio || puede("TRAMITE_RADICAR")) && <RegresoForm id={t.id} hoy={hoy} />}
          {esMio && t.estado === "REVISION_FINAL" && <RemisionForm id={t.id} />}
          {t.estado === "REMITIDO" && puede("TRAMITE_ARCHIVAR") &&
            (t.carpetasAbiertas > 0 ? (
              <Alerta tipo="aviso">
                Hay {t.carpetasAbiertas} {t.carpetasAbiertas === 1 ? "carpeta" : "carpetas"} de este trámite sin devolver al archivo: registra la devolución en Custodia antes de archivar.
              </Alerta>
            ) : (
              <ArchivarForm id={t.id} />
            ))}

          {!archivado && !esMio && t.analistaId && puede("TRAMITE_REVISAR") && (
            <p className="text-sm text-ink-faint">Este trámite está asignado a {t.analistaNombre}: solo esa persona puede revisarlo.</p>
          )}
          {puede("TRAMITE_ASIGNAR") && t.pendienteValidacion && (
            <p className="text-sm text-ink-faint">No se puede asignar hasta que la solicitud se valide en ventanilla.</p>
          )}
        </div>
      </div>
    </div>
  );
}
