"use client";

import { useActionState, useRef, useState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import type { Docente } from "@/lib/api-tramites";
import type { Carpeta } from "@/lib/api-legajos";
import type { EstadoFormulario } from "@/lib/formularios";
import { buscarDocentePorCodigoAction } from "../tramites/actions";

const estadoInicial: EstadoFormulario = {};
const RETARDO_BUSQUEDA_MS = 400;

export default function LegajoForm({
  accion,
  carpeta,
  titulo,
}: {
  accion: (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;
  carpeta?: Carpeta;
  titulo: string;
}) {
  const [estado, formAction] = useActionState(accion, estadoInicial);
  const errores = estado.camposInvalidos;

  const [codigo, setCodigo] = useState("");
  const [docente, setDocente] = useState<Docente | null>(null);
  const [buscando, setBuscando] = useState(false);
  const [errorBusqueda, setErrorBusqueda] = useState("");
  const temporizador = useRef<ReturnType<typeof setTimeout> | null>(null);
  const idBusqueda = useRef(0); // descarta respuestas de búsquedas anteriores que lleguen tarde

  function alCambiarCodigo(valor: string) {
    setCodigo(valor);
    setDocente(null);
    setErrorBusqueda("");
    if (temporizador.current) clearTimeout(temporizador.current);
    const idActual = ++idBusqueda.current;
    const buscado = valor.trim();
    if (!buscado) {
      setBuscando(false);
      return;
    }
    setBuscando(true);
    temporizador.current = setTimeout(async () => {
      const encontrado = await buscarDocentePorCodigoAction(buscado);
      if (idActual !== idBusqueda.current) return;
      setBuscando(false);
      if (encontrado) setDocente(encontrado);
      else setErrorBusqueda("No se encontró ningún docente con ese código.");
    }, RETARDO_BUSQUEDA_MS);
  }

  return (
    <Tarjeta nivel="h1" titulo={titulo} className="max-w-xl">
      <form action={formAction} className="flex flex-col gap-4">
        {carpeta ? (
          <Campo etiqueta="Docente">
            <div className={claseControl({ bloqueado: true })}>
              {carpeta.docenteNombre} · código {carpeta.docenteCodigo || "—"}
            </div>
          </Campo>
        ) : (
          <Campo id="codigoDocente" etiqueta="Código del docente" error={errorBusqueda || errores?.docenteId}>
            <input
              id="codigoDocente"
              autoComplete="off"
              placeholder="Ej. 4821"
              value={codigo}
              onChange={(e) => alCambiarCodigo(e.target.value)}
              className={claseControl()}
            />
            <input type="hidden" name="docenteId" value={docente?.id ?? ""} />
            {buscando && <p className="mt-1 text-xs text-ink-faint">Buscando…</p>}
            {docente && (
              <p className="mt-1 text-sm font-semibold text-ink">
                {docente.apellidos}, {docente.nombres}
              </p>
            )}
          </Campo>
        )}

        <div className="flex gap-4">
          <Campo id="codigo" etiqueta="Código del legajo" opcional ayuda="Si lo dejas vacío se asigna uno (LEG-00001)." error={errores?.codigo} className="flex-1">
            <input id="codigo" name="codigo" defaultValue={carpeta?.codigo} className={claseControl()} />
          </Campo>
          <Campo id="gaveta" etiqueta="Gaveta" opcional error={errores?.gaveta} className="w-32">
            <input id="gaveta" name="gaveta" defaultValue={carpeta?.gaveta} className={claseControl()} />
          </Campo>
          <Campo id="orden" etiqueta="Orden" error={errores?.orden} className="w-24">
            <input id="orden" name="orden" type="number" min={0} defaultValue={carpeta?.orden ?? 0} className={claseControl()} />
          </Campo>
        </div>

        <Campo id="anexos" etiqueta="Anexos" opcional error={errores?.anexos}>
          <input id="anexos" name="anexos" defaultValue={carpeta?.anexos} className={claseControl()} />
        </Campo>
        <Campo id="observacion" etiqueta="Observación" opcional error={errores?.observacion}>
          <input id="observacion" name="observacion" defaultValue={carpeta?.observacion} className={claseControl()} />
        </Campo>

        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}
        <AccionesFormulario cancelarHref="/panel/legajos" textoEnviar="Guardar legajo" />
      </form>
    </Tarjeta>
  );
}
