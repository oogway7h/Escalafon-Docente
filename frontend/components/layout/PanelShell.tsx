"use client";

import { useState, useSyncExternalStore } from "react";
import type { Sesion } from "@/lib/api";
import Header from "./Header";
import Sidebar from "./Sidebar";

const CLAVE_COLAPSADO = "escalafon:sidebar-colapsado";

function suscribirseSinCambios() {
  return () => {};
}

function leerColapsadoGuardado() {
  try {
    return window.localStorage.getItem(CLAVE_COLAPSADO) === "1";
  } catch {
    return false;
  }
}

export default function PanelShell({ sesion, fotoVersion, children }: { sesion: Sesion; fotoVersion: number | null; children: React.ReactNode }) {
  // useSyncExternalStore evita el desajuste de hidratación: en el servidor usa el snapshot fijo
  // (false), en el cliente lee localStorage sin pasar por un efecto.
  const colapsadoGuardado = useSyncExternalStore(suscribirseSinCambios, leerColapsadoGuardado, () => false);
  const [colapsadoManual, setColapsadoManual] = useState<boolean | null>(null);
  const [abiertoMovil, setAbiertoMovil] = useState(false);

  const colapsado = colapsadoManual ?? colapsadoGuardado;

  const alternarColapso = () => {
    const siguiente = !colapsado;
    try {
      window.localStorage.setItem(CLAVE_COLAPSADO, siguiente ? "1" : "0");
    } catch {}
    setColapsadoManual(siguiente);
  };

  return (
    <div className="flex min-h-full bg-bg">
      <Sidebar
        sesion={sesion}
        colapsado={colapsado}
        abiertoMovil={abiertoMovil}
        onCerrarMovil={() => setAbiertoMovil(false)}
      />
      <div className="flex min-h-screen min-w-0 flex-1 flex-col">
        <Header
          sesion={sesion}
          fotoVersion={fotoVersion}
          colapsado={colapsado}
          onAlternarColapso={alternarColapso}
          onAbrirMovil={() => setAbiertoMovil(true)}
        />
        <main className="flex-1 px-6 py-8">{children}</main>
      </div>
    </div>
  );
}
