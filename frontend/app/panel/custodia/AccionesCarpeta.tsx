"use client";

import { useState } from "react";
import BotonEnviar from "@/components/common/BotonEnviar";
import { claseControl } from "@/components/common/Campo";
import { ESTADOS_MATERIAL, enDias } from "@/lib/legajos";
import { atenderAction, devolverAction, entregarAction } from "./actions";

// Selector de estado material con botones segmentados (radios accesibles con aspecto de pastillas).
function EstadoMaterial({ titulo, alCambiar }: { titulo: string; alCambiar?: (valor: string) => void }) {
  return (
    <fieldset className="flex flex-col gap-2">
      <legend className="mb-2 text-sm font-semibold text-ink">{titulo}</legend>
      <div className="grid grid-cols-3 gap-2">
        {ESTADOS_MATERIAL.map((e) => (
          <label key={e.valor} className="cursor-pointer">
            <input type="radio" name="estadoMaterial" value={e.valor} defaultChecked={e.valor === "BUENO"} onChange={() => alCambiar?.(e.valor)} className="peer sr-only" />
            <span className="block rounded-xl border border-line px-1.5 py-2.5 text-center text-sm font-semibold text-ink-soft transition-colors peer-checked:border-ink peer-checked:bg-ink peer-checked:text-surface peer-focus-visible:ring-2 peer-focus-visible:ring-accent">
              {e.etiqueta}
            </span>
          </label>
        ))}
      </div>
    </fieldset>
  );
}

// Aprobar y entregar en un paso (o solo entregar, si la solicitud ya estaba aprobada).
export function EntregaForm({ id, tab, aprobada }: { id: string; tab: string; aprobada: boolean }) {
  const [fecha, setFecha] = useState(enDias(7));
  const chip = (dias: number) => (
    <button
      key={dias}
      type="button"
      onClick={() => setFecha(enDias(dias))}
      className={`rounded-full border px-3 py-1.5 text-sm font-semibold transition-colors ${fecha === enDias(dias) ? "border-ink bg-ink text-surface" : "border-line text-ink-soft hover:border-accent hover:text-accent-ink"}`}
    >
      +{dias} días
    </button>
  );

  return (
    <div className="flex flex-col gap-4">
      <form action={entregarAction} className="flex flex-col gap-4">
        <input type="hidden" name="id" value={id} />
        <input type="hidden" name="sel" value={id} />
        <input type="hidden" name="tab" value={tab} />
        <div className="flex flex-col gap-2">
          <label htmlFor={`limite-${id}`} className="text-sm font-semibold text-ink">
            Devolver antes de
          </label>
          <div className="flex flex-wrap items-center gap-2">
            <input id={`limite-${id}`} type="date" name="fechaLimite" required min={enDias(1)} value={fecha} onChange={(e) => setFecha(e.target.value)} className={`${claseControl()} !w-auto`} />
            {[3, 7, 15].map(chip)}
          </div>
        </div>
        <EstadoMaterial titulo="Estado de la carpeta al salir" />
        <div className="flex flex-col gap-2">
          <label htmlFor={`obs-${id}`} className="text-sm font-semibold text-ink">
            Observación <span className="font-normal text-ink-faint">(opcional)</span>
          </label>
          <textarea id={`obs-${id}`} name="observacion" rows={2} placeholder="Ej.: falta la hoja 12" className={claseControl()} />
        </div>
        <BotonEnviar clase="w-full !py-3.5 text-base" textoEnviando={aprobada ? "Entregando…" : "Aprobando y entregando…"}>
          {aprobada ? "Entregar carpeta" : "Aprobar y entregar carpeta"}
        </BotonEnviar>
        <p className="-mt-1 text-center text-xs leading-relaxed text-ink-faint">
          {aprobada ? "Registra la salida a nombre del solicitante." : "Un solo paso: aprueba la solicitud y registra la salida a nombre del solicitante."} Queda en la bitácora.
        </p>
      </form>

      {!aprobada && (
        <div className="flex flex-col gap-3 border-t border-line pt-3.5 text-sm font-semibold">
          <form action={atenderAction}>
            <input type="hidden" name="id" value={id} />
            <input type="hidden" name="sel" value={id} />
            <input type="hidden" name="tab" value={tab} />
            <input type="hidden" name="aprobar" value="true" />
            <button type="submit" className="text-ink-soft hover:text-accent-ink hover:underline">
              Solo aprobar, entregar después
            </button>
          </form>
          <details>
            <summary className="w-fit cursor-pointer text-accent-dark hover:underline">Rechazar con motivo…</summary>
            <form action={atenderAction} className="mt-2.5 flex flex-col gap-2">
              <input type="hidden" name="id" value={id} />
              <input type="hidden" name="sel" value={id} />
              <input type="hidden" name="tab" value={tab} />
              <input type="hidden" name="aprobar" value="false" />
              <textarea name="nota" required rows={2} placeholder="Motivo del rechazo (lo verá el solicitante)" aria-label="Motivo del rechazo" className={claseControl()} />
              <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Rechazando…">
                Rechazar solicitud
              </BotonEnviar>
            </form>
          </details>
        </div>
      )}
    </div>
  );
}

// Recibir la carpeta: estado material y observación (obligatoria si no vuelve en buen estado).
export function DevolucionForm({ legajoId, id, tab }: { legajoId: number; id: string; tab: string }) {
  const [material, setMaterial] = useState("BUENO");
  return (
    <form action={devolverAction} className="flex flex-col gap-4">
      <input type="hidden" name="legajoId" value={legajoId} />
      <input type="hidden" name="sel" value={id} />
      <input type="hidden" name="tab" value={tab} />
      <EstadoMaterial titulo="Estado de la carpeta al volver" alCambiar={setMaterial} />
      <div className="flex flex-col gap-2">
        <label htmlFor={`obs-${id}`} className="text-sm font-semibold text-ink">
          Observación <span className="font-normal text-ink-faint">{material === "BUENO" ? "(opcional)" : "(obligatoria: la carpeta no está en buen estado)"}</span>
        </label>
        <textarea
          id={`obs-${id}`}
          name="observacion"
          rows={2}
          required={material !== "BUENO"}
          placeholder="Ej.: faltan dos hojas del tomo 2"
          className={claseControl()}
        />
      </div>
      <BotonEnviar clase="w-full !py-3.5 text-base" textoEnviando="Registrando…">
        Registrar devolución
      </BotonEnviar>
      <p className="-mt-1 text-center text-xs leading-relaxed text-ink-faint">La carpeta vuelve a «En archivo» y el ingreso queda en la bitácora.</p>
    </form>
  );
}
