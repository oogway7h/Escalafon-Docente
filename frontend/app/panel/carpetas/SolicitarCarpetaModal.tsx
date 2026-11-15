"use client";

import { useActionState, useRef, useState } from "react";
import Alerta from "@/components/common/Alerta";
import { claseBoton } from "@/components/common/Boton";
import BotonEnviar from "@/components/common/BotonEnviar";
import { claseControl } from "@/components/common/Campo";
import Insignia from "@/components/common/Insignia";
import type { EstadoFormulario } from "@/lib/formularios";
import { datosParaSolicitudAction, solicitarCarpetaModalAction, type DatosSolicitud } from "./actions";

const estadoInicial: EstadoFormulario = {};

// CU-20: el botón abre un modal donde se busca al docente por su código; al encontrarlo se muestra su tarjeta
// con los datos y, debajo, el formulario para pedir su carpeta. Ligar la solicitud a un trámite es opcional: se
// puede pedir cualquier carpeta aunque no haya un trámite abierto de ese docente.
export default function SolicitarCarpetaModal() {
  const dialogo = useRef<HTMLDialogElement>(null);
  const [codigo, setCodigo] = useState("");
  const [datos, setDatos] = useState<DatosSolicitud | null>(null);
  const [fallo, setFallo] = useState("");
  const [estado, solicitar] = useActionState(solicitarCarpetaModalAction, estadoInicial);

  const abrir = () => dialogo.current?.showModal();
  const cerrar = () => dialogo.current?.close();
  const reiniciar = () => {
    setCodigo("");
    setDatos(null);
    setFallo("");
  };

  async function buscar(formData: FormData) {
    const valor = String(formData.get("codigo") ?? "").trim();
    setFallo("");
    if (!valor) return setFallo("Escribe el código del docente.");
    try {
      setDatos(await datosParaSolicitudAction(valor));
    } catch {
      setDatos(null);
      setFallo("No se pudo buscar al docente. Intenta de nuevo.");
    }
  }

  const d = datos?.docente;
  const iniciales = d ? `${d.nombres.charAt(0)}${d.apellidos.charAt(0)}`.toUpperCase() : "";
  const sinCarpeta = !!datos?.docente && datos.carpetas.length === 0;

  return (
    <>
      <button type="button" onClick={abrir} className={claseBoton("primario", "sm")}>
        + Solicitar carpeta
      </button>

      <dialog
        ref={dialogo}
        onClose={reiniciar}
        onClick={(e) => e.target === e.currentTarget && cerrar()} // clic en el fondo oscuro
        aria-labelledby="titulo-solicitar-carpeta"
        className="m-auto w-[min(560px,calc(100vw-32px))] rounded-xl border border-line bg-surface p-0 text-ink backdrop:bg-black/50"
      >
        <div className="flex flex-col gap-5 p-6">
          <header className="flex items-start justify-between gap-4">
            <div>
              <h2 id="titulo-solicitar-carpeta" className="font-display text-xl font-semibold tracking-[-0.01em]">
                Solicitar carpeta
              </h2>
              <p className="mt-1 text-sm text-ink-soft">Busca al docente por su código para pedir su carpeta física al archivo.</p>
            </div>
            <button type="button" onClick={cerrar} aria-label="Cerrar" className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg border border-line text-ink-soft hover:border-accent hover:text-accent-ink">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
                <path d="M18 6 6 18" />
                <path d="m6 6 12 12" />
              </svg>
            </button>
          </header>

          <form action={buscar} className="flex gap-2.5">
            <input
              name="codigo"
              value={codigo}
              onChange={(e) => setCodigo(e.target.value)}
              placeholder="Código del docente (ej. 5948)"
              aria-label="Código del docente"
              autoComplete="off"
              className={`${claseControl()} min-w-0 flex-1`}
            />
            <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Buscando…">
              Buscar
            </BotonEnviar>
          </form>
          {fallo && <Alerta tipo="error">{fallo}</Alerta>}

          {datos && !d && <Alerta tipo="aviso">No existe un docente con el código «{codigo.trim()}». Revisa que esté completo.</Alerta>}

          {d && (
            <div className="flex items-center gap-4 rounded-xl border border-line bg-bg p-4">
              <span className="flex h-12 w-12 shrink-0 items-center justify-center rounded-full bg-accent-soft font-display text-base font-semibold text-accent-dark">{iniciales}</span>
              <div className="min-w-0 flex-1">
                <div className="flex flex-wrap items-center gap-2">
                  <span className="font-display text-base font-semibold">
                    {d.apellidos}, {d.nombres}
                  </span>
                  <Insignia tono={d.activo ? "exito" : "neutro"}>{d.activo ? "Activo" : "De baja"}</Insignia>
                </div>
                <dl className="mt-1.5 grid grid-cols-[auto_1fr] gap-x-3 gap-y-0.5 text-sm">
                  <dt className="text-ink-faint">Código</dt>
                  <dd className="font-mono font-semibold">{d.codigo}</dd>
                  {d.cedula && (
                    <>
                      <dt className="text-ink-faint">Cédula</dt>
                      <dd>{d.cedula}</dd>
                    </>
                  )}
                  <dt className="text-ink-faint">Categoría</dt>
                  <dd>{d.categoria || "Sin categoría"}</dd>
                  <dt className="text-ink-faint">Facultad</dt>
                  <dd className="break-words">{d.facultad || "Sin facultad"}</dd>
                </dl>
              </div>
            </div>
          )}

          {sinCarpeta && <Alerta tipo="aviso">El docente no tiene carpetas disponibles en el archivo (pueden estar prestadas o sin inventariar).</Alerta>}

          {datos && d && !sinCarpeta && (
            <form action={solicitar} className="flex flex-col gap-3.5">
              {/* Con trámites a su cargo se propone el primero (ligarlo habilita iniciar la revisión); sin ellos, la solicitud va suelta. */}
              {datos.tramites.length > 0 && (
                <select name="tramiteId" defaultValue={datos.tramites[0].id} aria-label="Trámite" className={claseControl()}>
                  {datos.tramites.map((t) => (
                    <option key={t.id} value={t.id}>
                      Trámite {t.descripcion}
                    </option>
                  ))}
                  <option value="">Sin trámite</option>
                </select>
              )}
              {datos.carpetas.length === 1 ? (
                <>
                  <input type="hidden" name="legajoId" value={datos.carpetas[0].id} />
                  <p className="text-sm text-ink-soft">
                    Carpeta: <span className="font-mono font-semibold text-ink">{datos.carpetas[0].codigo}</span> · Gaveta {datos.carpetas[0].gaveta || "s/g"}
                  </p>
                </>
              ) : (
                <select name="legajoId" defaultValue={datos.carpetas[0].id} aria-label="Carpeta" className={claseControl()}>
                  {datos.carpetas.map((c) => (
                    <option key={c.id} value={c.id}>
                      Carpeta {c.codigo} · Gaveta {c.gaveta || "s/g"}
                    </option>
                  ))}
                </select>
              )}
              <input name="motivo" required defaultValue="Revisión del expediente" placeholder="Motivo de la solicitud" aria-label="Motivo" className={claseControl()} />
              {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
              <div className="flex justify-end gap-2.5">
                <button type="button" onClick={cerrar} className={claseBoton("secundario", "sm")}>
                  Cancelar
                </button>
                <BotonEnviar tamano="sm" textoEnviando="Solicitando…">
                  Solicitar carpeta
                </BotonEnviar>
              </div>
            </form>
          )}
        </div>
      </dialog>
    </>
  );
}
