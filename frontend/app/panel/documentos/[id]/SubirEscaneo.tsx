"use client";

import { useActionState, useState } from "react";
import { useRouter } from "next/navigation";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import ZonaArchivo from "@/components/common/ZonaArchivo";
import type { EstadoFormulario } from "@/lib/formularios";
import { LIMITE_MB, TIPOS_DOCUMENTO } from "@/lib/tramites";
import { subirAlAlmacen } from "../../tramites/[id]/CargaDocumentos";
import { confirmarSubidaExpedienteAction, darDeBajaExpedienteAction, prepararSubidaExpedienteAction } from "../actions";

const MB = 1024 * 1024;
const TIPOS_ARCHIVO = ["application/pdf", "image/jpeg", "image/png"];
const limiteMB = (tipo: string) => (tipo === "application/pdf" ? LIMITE_MB.pdf : LIMITE_MB.imagen);

// Archivo escanea la carpeta física y sube cada documento con su tipo; el navegador lo manda directo a S3
// y el servidor extrae el texto (OCR) en segundo plano.
export function SubirEscaneoForm({ docenteId }: { docenteId: string }) {
  const router = useRouter();
  const [intento, setIntento] = useState(0); // al terminar una subida, la zona de archivo se monta de nuevo vacía
  const [progreso, setProgreso] = useState<number | null>(null);
  const [error, setError] = useState("");

  async function cargar(formData: FormData) {
    setError("");
    const archivo = formData.get("archivo");
    const tipo = String(formData.get("tipoDocumento") ?? "");
    if (!tipo) return setError("Elige el tipo de documento.");
    if (!(archivo instanceof File) || archivo.size === 0) return setError("Selecciona el archivo escaneado.");
    if (!TIPOS_ARCHIVO.includes(archivo.type)) return setError("Solo se admiten archivos PDF, JPG o PNG.");
    if (archivo.size > limiteMB(archivo.type) * MB) return setError(`El archivo supera los ${limiteMB(archivo.type)} MB.`);

    setProgreso(0);
    try {
      const preparada = await prepararSubidaExpedienteAction(docenteId, {
        nombre: archivo.name,
        tipoMime: archivo.type,
        tamano: archivo.size,
        tipoDocumento: tipo,
      });
      if (preparada.error || !preparada.url || !preparada.documentoId) return setError(preparada.error ?? "No se pudo preparar la subida.");

      await subirAlAlmacen(preparada.url, archivo, setProgreso);
      const confirmada = await confirmarSubidaExpedienteAction(docenteId, preparada.documentoId);
      if (confirmada.error) return setError(confirmada.error);
      setIntento((n) => n + 1);
      router.refresh();
    } catch (e) {
      setError(e instanceof Error ? e.message : "No se pudo cargar el documento.");
    } finally {
      setProgreso(null);
    }
  }

  const cargando = progreso !== null;
  return (
    <Tarjeta titulo="Subir escaneo" descripcion="Se extrae el texto (OCR) para poder buscar cualquier página del expediente.">
      <form action={cargar} className="flex flex-col gap-3.5">
        <select name="tipoDocumento" defaultValue="" aria-label="Tipo de documento" className={`${claseControl()} sm:max-w-xs`} disabled={cargando}>
          <option value="" disabled>
            Tipo de documento
          </option>
          {TIPOS_DOCUMENTO.map((t) => (
            <option key={t.valor} value={t.valor}>
              {t.etiqueta}
            </option>
          ))}
        </select>
        <ZonaArchivo
          key={intento}
          name="archivo"
          maxMB={limiteMB}
          disabled={cargando}
          etiqueta={`PDF hasta ${LIMITE_MB.pdf} MB, JPG o PNG hasta ${LIMITE_MB.imagen} MB`}
        />
        <BotonEnviar tamano="sm" textoEnviando="Subiendo…">
          Subir
        </BotonEnviar>
        {cargando && (
          <div className="flex flex-col gap-1" role="status">
            <progress value={progreso} max={100} className="h-2 w-full" />
            <span className="text-xs text-ink-soft">{progreso < 100 ? `Subiendo… ${progreso}%` : "Verificando el archivo…"}</span>
          </div>
        )}
        {error && <Alerta tipo="error">{error}</Alerta>}
      </form>
    </Tarjeta>
  );
}

// Reintenta un documento cuyo procesamiento falló.
export function ReintentarEscaneo({ docenteId, docId }: { docenteId: string; docId: string }) {
  const router = useRouter();
  const [error, setError] = useState("");

  async function reintentar() {
    const r = await confirmarSubidaExpedienteAction(docenteId, docId);
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

const estadoInicial: EstadoFormulario = {};

// Baja lógica de un escaneo subido por error: el archivo se conserva y el motivo queda en la bitácora.
export function BajaEscaneoForm({ docenteId, docId }: { docenteId: string; docId: string }) {
  const [estado, formAction] = useActionState(darDeBajaExpedienteAction.bind(null, docenteId, docId), estadoInicial);
  return (
    <details className="text-sm">
      <summary className="cursor-pointer font-semibold text-ink-soft hover:text-accent-ink">Dar de baja</summary>
      <form action={formAction} className="mt-2 flex flex-col gap-2">
        <textarea name="motivo" rows={2} placeholder="Motivo de la baja" aria-label="Motivo de la baja" className={claseControl()} />
        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <BotonEnviar variante="secundario" tamano="sm" textoEnviando="Dando de baja…">
          Confirmar baja
        </BotonEnviar>
      </form>
    </details>
  );
}
