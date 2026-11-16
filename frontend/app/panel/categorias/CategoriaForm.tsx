"use client";

import { useActionState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import type { Categoria } from "@/lib/api-tramites";
import type { EstadoFormulario } from "@/lib/formularios";

const estadoInicial: EstadoFormulario = {};

export default function CategoriaForm({
  accion,
  categoria,
  titulo,
}: {
  accion: (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;
  categoria?: Categoria;
  titulo: string;
}) {
  const [estado, formAction] = useActionState(accion, estadoInicial);
  const errores = estado.camposInvalidos;

  return (
    <Tarjeta nivel="h1" titulo={titulo} className="max-w-xl">
      <form action={formAction} className="flex flex-col gap-4">
        <div className="flex gap-4">
          <Campo id="nombre" etiqueta="Nombre" error={errores?.nombre} className="flex-[2]">
            <input id="nombre" name="nombre" defaultValue={categoria?.nombre} placeholder="Titular A" className={claseControl()} />
          </Campo>
          <Campo id="aniosAntiguedad" etiqueta="Años de antigüedad" error={errores?.aniosAntiguedad} className="flex-1">
            <input
              id="aniosAntiguedad"
              name="aniosAntiguedad"
              type="number"
              min={0}
              defaultValue={categoria?.aniosAntiguedad ?? 0}
              className={claseControl()}
            />
          </Campo>
        </div>

        <div className="flex gap-4">
          <Campo id="puntajeDesde" etiqueta="Puntaje desde" ayuda="Puntos de categoría (art. 18)." error={errores?.puntajeDesde} className="flex-1">
            <input id="puntajeDesde" name="puntajeDesde" type="number" min={0} defaultValue={categoria?.puntajeDesde ?? ""} className={claseControl()} />
          </Campo>
          <Campo id="puntajeHasta" etiqueta="Puntaje hasta" ayuda="Vacío = «o más». Sin puntaje, no sale en las tarjetas del portal." error={errores?.puntajeHasta} className="flex-1">
            <input id="puntajeHasta" name="puntajeHasta" type="number" min={0} defaultValue={categoria?.puntajeHasta ?? ""} className={claseControl()} />
          </Campo>
        </div>

        <Campo
          id="requisitos"
          etiqueta="Documentos probatorios exigidos"
          ayuda="Uno por línea. El portal público pedirá un archivo por cada uno."
          error={errores?.requisitos}
        >
          <textarea
            id="requisitos"
            name="requisitos"
            rows={6}
            defaultValue={categoria?.requisitos.map((r) => r.nombre).join("\n")}
            className={claseControl()}
          />
        </Campo>

        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

        <AccionesFormulario cancelarHref="/panel/categorias" textoEnviar="Guardar categoría" />
      </form>
    </Tarjeta>
  );
}
