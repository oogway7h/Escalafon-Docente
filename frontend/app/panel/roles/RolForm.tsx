"use client";

import { useActionState } from "react";
import AccionesFormulario from "@/components/common/AccionesFormulario";
import Alerta from "@/components/common/Alerta";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import type { Permiso } from "@/lib/api";
import type { EstadoFormulario } from "@/lib/formularios";

type AccionRol = (estado: EstadoFormulario, formData: FormData) => Promise<EstadoFormulario>;

const estadoInicial: EstadoFormulario = {};

export default function RolForm({
  accion,
  permisos,
  valoresIniciales,
  esSistema,
  titulo,
  textoBoton,
}: {
  accion: AccionRol;
  permisos: Permiso[];
  valoresIniciales?: { nombre: string; descripcion: string; permisoIds: number[] };
  esSistema?: boolean;
  titulo: string;
  textoBoton: string;
}) {
  const [estado, formAction] = useActionState(accion, estadoInicial);
  const marcados = new Set(valoresIniciales?.permisoIds ?? []);

  const porModulo = new Map<string, Permiso[]>();
  for (const permiso of permisos) {
    const lista = porModulo.get(permiso.modulo) ?? [];
    lista.push(permiso);
    porModulo.set(permiso.modulo, lista);
  }

  return (
    <Tarjeta nivel="h1" titulo={titulo} className="max-w-xl">
      <form action={formAction} className="flex flex-col gap-4">
        <Campo
          id="nombre"
          etiqueta="Nombre del rol"
          ayuda={esSistema ? "Bloqueado: es un rol de sistema y no se puede renombrar." : undefined}
          error={estado.camposInvalidos?.nombre}
        >
          <input
            id="nombre"
            name="nombre"
            defaultValue={valoresIniciales?.nombre}
            readOnly={esSistema}
            className={claseControl({ bloqueado: esSistema })}
          />
        </Campo>

        <Campo id="descripcion" etiqueta="Descripción" opcional>
          <input id="descripcion" name="descripcion" defaultValue={valoresIniciales?.descripcion} className={claseControl()} />
        </Campo>

        <Campo etiqueta="Permisos">
          <div className="flex max-h-[360px] flex-col gap-2.5 overflow-y-auto pr-1">
            {Array.from(porModulo.entries()).map(([modulo, items]) => (
              <div key={modulo} className="rounded-xl border border-line px-4 py-3">
                <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-ink-faint">{modulo}</p>
                <div className="flex flex-col gap-1.5">
                  {items.map((permiso) => (
                    <label key={permiso.id} className="flex items-center gap-2.5 text-sm text-ink">
                      <input
                        type="checkbox"
                        name="permisoIds"
                        value={permiso.id}
                        defaultChecked={marcados.has(permiso.id)}
                        className="h-[17px] w-[17px] rounded border border-line accent-accent"
                      />
                      {permiso.nombre}
                    </label>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </Campo>

        {estado.error && <Alerta tipo="error">{estado.error}</Alerta>}

        <AccionesFormulario cancelarHref="/panel/roles" textoEnviar={textoBoton} />
      </form>
    </Tarjeta>
  );
}
