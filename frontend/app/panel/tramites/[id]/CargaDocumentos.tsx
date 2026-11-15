"use client";

import { useActionState, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import { claseControl } from "@/components/common/Campo";
import type { EstadoFormulario } from "@/lib/formularios";
import { confirmarSubidaAction, corregirTextoAction } from "../actions";

// Sube con XHR (y no con fetch) para informar el avance: los PDF escaneados llegan a cientos de MB.
export function subirAlAlmacen(url: string, archivo: File, alAvanzar: (porcentaje: number) => void): Promise<void> {
  return new Promise((resolver, rechazar) => {
    const xhr = new XMLHttpRequest();
    xhr.open("PUT", url);
    xhr.setRequestHeader("Content-Type", archivo.type);
    xhr.upload.onprogress = (e) => e.lengthComputable && alAvanzar(Math.round((e.loaded / e.total) * 100));
    xhr.onload = () => (xhr.status >= 200 && xhr.status < 300 ? resolver() : rechazar(new Error(`El almacén respondió ${xhr.status}.`)));
    xhr.onerror = () => rechazar(new Error("No se pudo subir el archivo al almacén. Revisa tu conexión."));
    xhr.send(archivo);
  });
}

// Mientras algún documento se procesa, vuelve a pedir la página cada pocos segundos.
export function AutoRefrescar({ activo }: { activo: boolean }) {
  const router = useRouter();
  useEffect(() => {
    if (!activo) return;
    // Con la pestaña oculta no hay nadie mirando el avance: se salta el refresco.
    const t = setInterval(() => document.visibilityState === "visible" && router.refresh(), 10_000);
    return () => clearInterval(t);
  }, [activo, router]);
  return null;
}

const estadoInicial: EstadoFormulario = {};

// Texto extraído por OCR: se carga al abrirlo (puede pesar varios MB) y se guarda ya corregido.
export function TextoOcr({ id, docId, editable }: { id: string; docId: string; editable: boolean }) {
  const [paginas, setPaginas] = useState<string[] | null>(null);
  const [fallo, setFallo] = useState(false);
  const [estado, formAction] = useActionState(corregirTextoAction.bind(null, id, docId), estadoInicial);

  async function abrir(e: React.SyntheticEvent<HTMLDetailsElement>) {
    if (!e.currentTarget.open || paginas !== null) return;
    const r = await fetch(`/panel/tramites/${id}/documentos/${docId}/texto`).catch(() => null);
    if (!r?.ok) return setFallo(true);
    setPaginas((await r.json()).paginas);
  }

  return (
    <details onToggle={abrir} className="text-sm">
      <summary className="cursor-pointer font-semibold text-accent-ink">Ver texto extraído</summary>
      <form action={formAction} className="mt-2 flex flex-col gap-2">
        {fallo && <Alerta tipo="error">No se pudo cargar el texto.</Alerta>}
        {paginas === null && !fallo && <p className="text-ink-faint">Cargando…</p>}
        {paginas !== null && (
          <>
            <div className="flex max-h-[60vh] flex-col gap-3 overflow-y-auto pr-1">
              {paginas.map((p, i) => (
                <label key={i} className="flex flex-col gap-1 font-semibold text-ink-soft">
                  Página {i + 1}
                  <textarea name="pagina" defaultValue={p} rows={8} readOnly={!editable} className={`${claseControl()} font-mono text-xs font-normal`} />
                </label>
              ))}
            </div>
            {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
            {estado.mensaje && <Alerta tipo="exito">{estado.mensaje}</Alerta>}
            {editable && (
              <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Guardando…">
                Guardar texto corregido
              </BotonEnviar>
            )}
          </>
        )}
      </form>
    </details>
  );
}

// Reintenta un documento cuyo procesamiento falló.
export function ReintentarForm({ id, docId }: { id: string; docId: string }) {
  const router = useRouter();
  const [error, setError] = useState("");

  async function reintentar() {
    const r = await confirmarSubidaAction(id, docId);
    if (r.error) setError(r.error);
    else router.refresh();
  }

  return (
    <form action={reintentar} className="flex items-center gap-2">
      <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Reintentando…">
        Reintentar OCR
      </BotonEnviar>
      {error && <span className="text-xs text-ink-faint">{error}</span>}
    </form>
  );
}
