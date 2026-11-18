"use client";

import { useActionState, useRef, useState } from "react";
import Alerta from "@/components/common/Alerta";
import Avatar from "@/components/common/Avatar";
import BotonEnviar from "@/components/common/BotonEnviar";
import type { EstadoFormulario } from "@/lib/formularios";
import { quitarFoto, subirFoto } from "./actions";

const MAX = 2 * 1024 * 1024;
const TIPOS = ["image/jpeg", "image/png"];
const inicial: EstadoFormulario = {};

// Foto de perfil: elegir o arrastrar una imagen, ver cómo queda y guardarla. Es un <label> sobre un
// <input type="file"> nativo, así que el teclado y el clic funcionan sin código extra.
export default function FotoPerfil({
  nombres,
  apellidos,
  rol,
  correo,
  version,
}: {
  nombres: string;
  apellidos: string;
  rol: string;
  correo: string;
  version: number | null;
}) {
  const entrada = useRef<HTMLInputElement>(null);
  const [previa, setPrevia] = useState<string | undefined>();
  const [problema, setProblema] = useState("");
  const [arrastrando, setArrastrando] = useState(false);
  const [estadoSubida, accionSubir] = useActionState(subirFoto, inicial);
  const [estadoQuitar, accionQuitar] = useActionState(quitarFoto, inicial);

  function elegir(archivo: File | undefined) {
    if (previa) URL.revokeObjectURL(previa);
    setPrevia(undefined);
    if (!archivo) return setProblema("");
    if (!TIPOS.includes(archivo.type)) return rechazar("La foto debe ser JPG o PNG.");
    if (archivo.size > MAX) return rechazar("La foto no puede superar los 2 MB.");
    setProblema("");
    setPrevia(URL.createObjectURL(archivo));
  }

  function rechazar(texto: string) {
    if (entrada.current) entrada.current.value = "";
    setProblema(texto);
  }

  function soltar(e: React.DragEvent) {
    e.preventDefault();
    setArrastrando(false);
    const archivo = e.dataTransfer.files[0];
    if (!archivo || !entrada.current) return;
    const dt = new DataTransfer(); // un FileList no se edita: se arma con DataTransfer
    dt.items.add(archivo);
    entrada.current.files = dt.files;
    elegir(archivo);
  }

  const resultado = estadoSubida.mensaje ?? estadoQuitar.mensaje;
  const error = problema || estadoSubida.error || estadoQuitar.error;

  return (
    <div className="flex flex-col items-center gap-5">
      <Avatar nombres={nombres} apellidos={apellidos} version={version} tamano={120} vistaPrevia={previa} />

      <div className="text-center">
        <p className="text-lg font-semibold text-ink">
          {nombres} {apellidos}
        </p>
        <p className="text-sm text-ink-soft">{correo}</p>
        <p className="mt-2 inline-block rounded-full bg-accent-soft px-3 py-1 text-xs font-semibold text-accent-ink">{rol}</p>
      </div>

      <form action={accionSubir} className="flex w-full flex-col gap-3">
        <label
          htmlFor="foto"
          onDragOver={(e) => {
            e.preventDefault();
            setArrastrando(true);
          }}
          onDragLeave={() => setArrastrando(false)}
          onDrop={soltar}
          className={`flex cursor-pointer flex-col items-center gap-1 rounded-lg border-[1.5px] border-dashed px-4 py-4 text-center text-sm transition-colors focus-within:border-accent ${
            arrastrando ? "border-accent bg-accent-soft" : "border-line hover:border-accent"
          }`}
        >
          <span className="text-ink-soft">Arrastra una imagen aquí o elige un archivo</span>
          <span className="text-xs text-ink-faint">JPG o PNG, hasta 2 MB</span>
          <input
            ref={entrada}
            id="foto"
            name="foto"
            type="file"
            accept="image/jpeg,image/png"
            className="sr-only"
            onChange={(e) => elegir(e.target.files?.[0])}
          />
        </label>

        {error && <Alerta tipo="error">{error}</Alerta>}
        {resultado && !previa && <Alerta tipo="exito">{resultado}</Alerta>}

        <BotonEnviar tamano="sm" textoEnviando="Guardando…" deshabilitado={!previa}>
          Guardar foto
        </BotonEnviar>
      </form>

      {version !== null && !previa && (
        <form action={accionQuitar}>
          <BotonEnviar tamano="sm" variante="secundario" textoEnviando="Quitando…">
            Quitar foto
          </BotonEnviar>
        </form>
      )}
    </div>
  );
}
