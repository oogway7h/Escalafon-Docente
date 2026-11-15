"use client";

import Link from "next/link";
import { useActionState, useState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import { claseControl } from "@/components/common/Campo";
import { MAX_CODIGO_DOCENTE, TITULO_ESTADO, formatearFecha, tipoYCategoria } from "@/lib/tramites";
import { consultarAction, type EstadoBusqueda } from "./actions";

const estadoInicial: EstadoBusqueda = {};

// Un solo campo para el código: más de 5 caracteres es un código de seguimiento (no pide cédula); hasta 5 es un código
// de docente y se exige la cédula. La cédula no viaja en la URL: el formulario es una acción POST.
export default function BuscadorSeguimiento({ codigoInicial, sobreOscuro = false }: { codigoInicial?: string; sobreOscuro?: boolean }) {
  const [codigo, setCodigo] = useState(codigoInicial ?? "");
  const [estado, accion] = useActionState(consultarAction, estadoInicial);
  const escrito = codigo.trim();
  const esSeguimiento = escrito.length > MAX_CODIGO_DOCENTE;
  const pideCedula = escrito.length > 0 && !esSeguimiento;

  return (
    <div className="flex w-full max-w-xl flex-col gap-4">
      <form action={accion} className="flex flex-col gap-3">
        <div className="flex gap-2.5">
          <input
            type="text"
            name="codigo"
            value={codigo}
            onChange={(e) => setCodigo(e.target.value)}
            required
            autoFocus={!codigoInicial}
            autoCapitalize="characters"
            autoComplete="off"
            spellCheck={false}
            placeholder="ED-2026-XXXXXX o tu código de docente"
            aria-label="Código de seguimiento o código de docente"
            className={`${claseControl({ grande: true })} flex-1 font-mono uppercase tracking-wider`}
          />
          <BotonEnviar textoEnviando="Consultando…">Consultar</BotonEnviar>
        </div>
        <div className="flex flex-col gap-1.5">
          <input
            type="text"
            name="cedula"
            inputMode="numeric"
            disabled={!pideCedula}
            required={pideCedula}
            autoComplete="off"
            placeholder={esSeguimiento ? "No necesitas cédula con el código de seguimiento" : "Cédula de identidad (solo con tu código de docente)"}
            aria-label="Cédula de identidad"
            className={`${claseControl({ grande: true })} transition-opacity ${pideCedula ? "" : "cursor-not-allowed opacity-50"}`}
          />
          <p className={`px-1 text-[12.5px] ${sobreOscuro ? "text-white/75" : "text-ink-faint"}`}>
            {esSeguimiento
              ? "Con el código de seguimiento no hace falta la cédula."
              : "Si perdiste tu código de seguimiento, escribe tu código de docente (hasta 5 caracteres) y tu cédula."}
          </p>
        </div>
      </form>

      {estado.error && <Alerta tipo="aviso" className="animate-materializar text-center">{estado.error}</Alerta>}

      {estado.lista && (
        <section className="animate-materializar material flex flex-col gap-4 rounded-[24px] p-6 shadow-[0_24px_60px_-24px_rgba(15,27,46,0.35)] sm:p-8">
          <div className="text-center">
            <h2 className="font-display text-[clamp(1.25rem,3vw,1.75rem)] font-extrabold tracking-[-0.02em] text-ink">Tienes {estado.lista.length} trámites</h2>
            <p className="mt-1 text-[14.5px] text-ink-soft">Elige el que quieres consultar.</p>
          </div>
          <ul className="flex flex-col gap-2.5">
            {estado.lista.map((t) => (
              <li key={t.codigo}>
                <Link href={`/seguimiento?codigo=${encodeURIComponent(t.codigo)}`} className="flex flex-wrap items-center gap-x-4 gap-y-1 rounded-[14px] border border-line px-4 py-3.5 hover:border-accent">
                  <span className="font-mono text-[14px] font-bold text-ink">{t.codigo}</span>
                  <span className="min-w-0 flex-1 text-[14px] text-ink-soft">
                    {tipoYCategoria(t)}
                  </span>
                  <span className="text-[13px] font-semibold text-ink">{t.estado === "INGRESADO" && t.pendienteValidacion ? "Pendiente de validación" : TITULO_ESTADO[t.estado]}</span>
                  <span className="text-[12.5px] text-ink-faint">{formatearFecha(t.radicadoEn)}</span>
                </Link>
              </li>
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}
