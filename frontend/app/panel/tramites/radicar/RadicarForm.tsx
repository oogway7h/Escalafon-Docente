"use client";

import { useActionState, useRef, useState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import type { Categoria, Docente } from "@/lib/api-tramites";
import type { EstadoFormulario } from "@/lib/formularios";
import { buscarDocentePorCodigoAction, radicarAction } from "../actions";

const estadoInicial: EstadoFormulario = {};
const RETARDO_BUSQUEDA_MS = 400;

export default function RadicarForm({ categorias }: { categorias: Categoria[] }) {
  const [estado, formAction] = useActionState(radicarAction, estadoInicial);
  const errores = estado.camposInvalidos;

  const [tipo, setTipo] = useState("CATEGORIZACION");
  const [codigo, setCodigo] = useState("");
  const [docente, setDocente] = useState<Docente | null>(null);
  const [buscando, setBuscando] = useState(false);
  const [errorBusqueda, setErrorBusqueda] = useState("");
  const temporizador = useRef<ReturnType<typeof setTimeout> | null>(null);
  // Id de la última búsqueda disparada: descarta respuestas de búsquedas anteriores
  // que lleguen tarde (p. ej. "594" resolviendo después de "5948").
  const idBusqueda = useRef(0);

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
      if (!encontrado) {
        setErrorBusqueda("No se encontró ningún docente con ese código.");
      } else if (!encontrado.activo) {
        setErrorBusqueda(`${encontrado.nombres} ${encontrado.apellidos} está dado de baja y no se le puede radicar un trámite.`);
      } else {
        setDocente(encontrado);
      }
    }, RETARDO_BUSQUEDA_MS);
  }

  return (
    <Tarjeta
      nivel="h1"
      titulo="Iniciar trámite"
      descripcion="Registra la solicitud en ventanilla. Se genera un código de seguimiento y el trámite queda en estado Ingresado."
      className="max-w-xl"
    >
      <form action={formAction} className="flex flex-col gap-4">
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
            <div className="mt-1 flex flex-col gap-1 rounded-lg border border-line bg-bg px-3.5 py-3">
              <p className="text-sm font-semibold text-ink">
                {docente.apellidos}, {docente.nombres}
              </p>
              <dl className="grid grid-cols-2 gap-x-3 gap-y-0.5 text-xs text-ink-soft">
                <div className="flex gap-1">
                  <dt className="text-ink-faint">Código:</dt>
                  <dd>{docente.codigo || "—"}</dd>
                </div>
                <div className="flex gap-1">
                  <dt className="text-ink-faint">Cédula:</dt>
                  <dd>{docente.cedula || "—"}</dd>
                </div>
                <div className="col-span-2 flex gap-1">
                  <dt className="text-ink-faint">Categoría:</dt>
                  <dd>{docente.categoria || "Sin categoría"}</dd>
                </div>
                <div className="col-span-2 flex gap-1">
                  <dt className="shrink-0 text-ink-faint">Facultad:</dt>
                  <dd className="truncate" title={docente.facultad}>
                    {docente.facultad || "Sin facultad"}
                  </dd>
                </div>
              </dl>
            </div>
          )}
        </Campo>

        <div className="flex gap-4">
          <Campo id="tipo" etiqueta="Tipo de trámite" error={errores?.tipo} className="flex-1">
            <select id="tipo" name="tipo" value={tipo} onChange={(e) => setTipo(e.target.value)} className={claseControl()}>
              <option value="CATEGORIZACION">Categorización</option>
              <option value="REGULARIZACION">Regularización</option>
            </select>
          </Campo>
          {/* Una regularización no postula a ninguna categoría: el campo solo existe en la categorización. */}
          {tipo === "CATEGORIZACION" && (
            <Campo id="categoriaId" etiqueta="Categoría postulada" error={errores?.categoriaId} className="flex-1">
              <select id="categoriaId" name="categoriaId" defaultValue="" className={claseControl()}>
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
        </div>

        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

        <AccionesFormulario cancelarHref="/panel/tramites" textoEnviar="Radicar trámite" textoEnviando="Radicando…" />
      </form>
    </Tarjeta>
  );
}
