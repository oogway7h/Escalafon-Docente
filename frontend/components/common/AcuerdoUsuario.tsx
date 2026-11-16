"use client";

import { useState } from "react";
import BotonEnviar from "@/components/common/BotonEnviar";
import { cerrarSesion } from "@/app/panel/actions";
import { aceptarTerminos } from "@/app/panel/terminos";
import { DESARROLLADOR, FECHA_TERMINOS, SECCIONES_TERMINOS, VERSION_TERMINOS } from "@/lib/terminos";

// Pantalla completa que se muestra al iniciar sesión mientras el usuario no haya aceptado la versión vigente del
// acuerdo de usuario. Sin aceptar no se ve el panel; rechazar cierra la sesión.
export default function AcuerdoUsuario({ nombre }: { nombre: string }) {
  const [acepto, setAcepto] = useState(false);

  return (
    <main className="flex min-h-screen items-center justify-center bg-bg px-4 py-8">
      <div className="flex max-h-[calc(100vh-4rem)] w-full max-w-3xl flex-col gap-5 rounded-2xl border border-line bg-surface p-6 shadow-sm sm:p-8">
        <header className="flex flex-col gap-1.5">
          <p className="text-xs font-bold uppercase tracking-[0.12em] text-ink-faint">
            Versión {VERSION_TERMINOS} · {FECHA_TERMINOS}
          </p>
          <h1 className="font-display text-2xl font-extrabold tracking-[-0.01em] text-ink">Acuerdo de usuario</h1>
          <p className="text-sm text-ink-soft">
            Hola, {nombre}. Antes de entrar por primera vez lee y acepta este acuerdo con {DESARROLLADOR}.
          </p>
        </header>

        <div tabIndex={0} aria-label="Texto del acuerdo de usuario" className="min-h-0 flex-1 overflow-y-auto rounded-xl border border-line bg-bg p-5">
          <div className="flex flex-col gap-5">
            {SECCIONES_TERMINOS.map((s) => (
              <section key={s.titulo} className="flex flex-col gap-2">
                <h2 className="text-[15px] font-bold text-ink">{s.titulo}</h2>
                {s.parrafos.map((p, i) => (
                  <p key={i} className="text-sm leading-relaxed text-ink-soft">
                    {p}
                  </p>
                ))}
              </section>
            ))}
          </div>
        </div>

        <form action={aceptarTerminos} className="flex flex-col gap-4">
          <label className="flex cursor-pointer items-start gap-3 text-sm text-ink">
            <input
              type="checkbox"
              name="acepto"
              required
              checked={acepto}
              onChange={(e) => setAcepto(e.target.checked)}
              className="mt-0.5 h-5 w-5 shrink-0 accent-[var(--color-accent)]"
            />
            He leído el acuerdo de usuario, lo entiendo y me comprometo a cumplirlo.
          </label>
          <div className="flex flex-wrap items-center gap-3">
            <BotonEnviar variante="primario" tamano="md" textoEnviando="Guardando…">
              Aceptar y continuar
            </BotonEnviar>
            <button
              formAction={cerrarSesion}
              formNoValidate
              className="min-h-11 rounded-lg px-4 text-sm font-semibold text-ink-soft hover:text-accent-ink"
            >
              No acepto, cerrar sesión
            </button>
          </div>
        </form>
      </div>
    </main>
  );
}
