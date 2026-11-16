"use client";

import { useActionState, useState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonEnviar from "@/components/common/BotonEnviar";
import Campo, { claseControl } from "@/components/common/Campo";
import ZonaArchivo from "@/components/common/ZonaArchivo";
import type { Categoria } from "@/lib/api-tramites";
import type { EstadoFormulario } from "@/lib/formularios";
import { REQUISITOS_REGULARIZACION, textoPuntaje } from "@/lib/requisitos";
import { LIMITE_MB, MAX_ADJUNTOS_LIBRES } from "@/lib/tramites";
import { radicarPublicoAction } from "./actions";

const estadoInicial: EstadoFormulario = {};

export default function SolicitudForm({ categorias }: { categorias: Categoria[] }) {
  const [estado, formAction] = useActionState(radicarPublicoAction, estadoInicial);
  const [tipo, setTipo] = useState("CATEGORIZACION");
  const [categoriaId, setCategoriaId] = useState("");
  const errores = estado.camposInvalidos;
  const esCategorizacion = tipo === "CATEGORIZACION";
  const categoria = categorias.find((c) => String(c.id) === categoriaId);

  return (
    <form action={formAction} className="material flex w-full max-w-xl flex-col gap-5 rounded-[24px] p-6 sm:p-9">
      <div className="flex flex-col gap-5 sm:flex-row">
        <Campo id="codigo" etiqueta="Código de docente" ayuda="El que te otorgó la universidad." error={errores?.codigo} className="flex-1">
          <input id="codigo" name="codigo" inputMode="numeric" autoComplete="off" required className={claseControl({ grande: true })} />
        </Campo>
        <Campo id="cedula" etiqueta="Cédula de identidad" ayuda="Debe coincidir con la del padrón docente." error={errores?.cedula} className="flex-1">
          <input id="cedula" name="cedula" autoComplete="off" required className={claseControl({ grande: true })} />
        </Campo>
      </div>

      <Campo id="tipo" etiqueta="Tipo de trámite">
        <select id="tipo" name="tipo" value={tipo} onChange={(e) => setTipo(e.target.value)} className={claseControl({ grande: true })}>
          <option value="CATEGORIZACION">Categorización</option>
          <option value="REGULARIZACION">Regularización</option>
        </select>
      </Campo>

      {/* La categoría solo aplica a la categorización: una regularización no postula a ninguna. */}
      {esCategorizacion && (
        <Campo id="categoriaId" etiqueta="Categoría a la que postulas" error={errores?.categoriaId}>
          <select
            id="categoriaId"
            name="categoriaId"
            value={categoriaId}
            onChange={(e) => setCategoriaId(e.target.value)}
            required
            className={claseControl({ grande: true })}
          >
            <option value="" disabled>
              Selecciona
            </option>
            {categorias.map((c) => (
              <option key={c.id} value={c.id}>
                {c.nombre}
              </option>
            ))}
          </select>
        </Campo>
      )}

      {esCategorizacion && categoria && (
        <fieldset key={categoria.id} className="animate-materializar flex flex-col gap-4 rounded-[16px] border border-line p-5">
          <legend className="px-2 text-[13px] font-bold text-ink">Documentos para {categoria.nombre}</legend>
          <p className="-mt-1 text-xs text-ink-faint">
            {textoPuntaje(categoria.puntajeDesde, categoria.puntajeHasta) && `Puntaje de categoría: ${textoPuntaje(categoria.puntajeDesde, categoria.puntajeHasta)} puntos. `}
            Todos son obligatorios.
          </p>
          {categoria.requisitos.map((r) => (
            <Campo key={r.id} id={`doc_${r.id}`} etiqueta={r.nombre} error={errores?.[`doc_${r.id}`]}>
              <ZonaArchivo id={`doc_${r.id}`} name={`doc_${r.id}`} maxMB={LIMITE_MB.portal} required />
            </Campo>
          ))}
        </fieldset>
      )}

      {!esCategorizacion && (
        <div className="animate-materializar flex flex-col gap-2 rounded-[16px] border border-line p-5">
          <p className="text-[13px] font-bold text-ink">Documentos que exige la regularización</p>
          <ul className="flex list-disc flex-col gap-1 pl-5 text-xs leading-snug text-ink-soft">
            {REQUISITOS_REGULARIZACION.map((r) => (
              <li key={r}>{r}</li>
            ))}
          </ul>
        </div>
      )}

      {!esCategorizacion && (
        <Campo
          id="adjunto"
          etiqueta="Documentos de respaldo"
          ayuda={`Hasta ${MAX_ADJUNTOS_LIBRES} archivos: si tienes más documentos, júntalos en un mismo PDF.`}
          error={errores?.adjunto}
        >
          <ZonaArchivo id="adjunto" name="adjunto" maxMB={LIMITE_MB.portal} maxArchivos={MAX_ADJUNTOS_LIBRES} required />
        </Campo>
      )}

      {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

      <BotonEnviar tamano="lg" textoEnviando="Enviando solicitud…">
        Radicar solicitud
      </BotonEnviar>
    </form>
  );
}
