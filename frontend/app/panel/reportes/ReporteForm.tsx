"use client";

import { useState } from "react";
import Alerta from "@/components/common/Alerta";
import BotonCarga from "@/components/common/BotonCarga";
import Campo, { claseControl } from "@/components/common/Campo";
import Tarjeta from "@/components/common/Tarjeta";
import { mensajeReporte, TIPOS_REPORTE, type TipoReporte } from "@/lib/indicadores";
import { ETIQUETA_ESTADO, ETIQUETA_TIPO, FASES } from "@/lib/tramites";

type Mensaje = { tipo: "error" | "exito"; texto: string };

const FORMATOS = [
  { valor: "PDF", etiqueta: "PDF", ayuda: "Con membrete, fecha y responsable" },
  { valor: "XLSX", etiqueta: "Excel", ayuda: "Tabla con filtros por columna" },
] as const;

const opcion = (activa: boolean) =>
  `flex cursor-pointer flex-col gap-1 rounded-xl border-2 px-4 py-3 transition-colors has-[:focus-visible]:shadow-[0_0_0_4px_var(--color-accent-soft)] ${
    activa ? "border-accent bg-accent-soft/40" : "border-line bg-surface hover:border-accent"
  }`;

// Nombre del archivo que propone la API en Content-Disposition.
function nombreArchivo(cabecera: string | null, porDefecto: string): string {
  return /filename="([^"]+)"/.exec(cabecera ?? "")?.[1] ?? porDefecto;
}

// CU-26: elegir el tipo de reporte, los filtros y el formato. Con JavaScript la descarga ocurre sin salir
// de la página (y un conjunto vacío se explica aquí mismo); sin JavaScript el formulario se envía igual.
export default function ReporteForm({
  categorias,
  tipoInicial,
  errorInicial,
}: {
  categorias: { id: number; nombre: string }[];
  tipoInicial: TipoReporte;
  errorInicial?: string;
}) {
  const [tipo, setTipo] = useState<TipoReporte>(tipoInicial);
  const [formato, setFormato] = useState<string>("PDF");
  const [cargando, setCargando] = useState(false);
  const [frase, setFrase] = useState("");
  const [interpretando, setInterpretando] = useState(false);
  const [nota, setNota] = useState<Mensaje | null>(null);
  // Valores propuestos por el asistente; `version` remonta los filtros para que los tomen como valor inicial.
  const [valores, setValores] = useState<Record<string, string>>({});
  const [version, setVersion] = useState(0);
  const [mensaje, setMensaje] = useState<Mensaje | null>(errorInicial ? { tipo: "error", texto: errorInicial } : null);

  async function generar(formulario: HTMLFormElement) {
    setCargando(true);
    setMensaje(null);
    try {
      const respuesta = await fetch(formulario.action, {
        method: "POST",
        body: new FormData(formulario),
        headers: { "x-descarga": "fetch" },
      });
      if (respuesta.redirected) {
        window.location.assign(respuesta.url); // la sesión expiró: el servidor mandó al login
        return;
      }
      if (!respuesta.ok) {
        const cuerpo = await respuesta.json().catch(() => null);
        setMensaje({ tipo: "error", texto: mensajeReporte(cuerpo?.codigo) });
        return;
      }
      const archivo = await respuesta.blob();
      const enlace = document.createElement("a");
      enlace.href = URL.createObjectURL(archivo);
      enlace.download = nombreArchivo(respuesta.headers.get("content-disposition"), `reporte.${formato === "PDF" ? "pdf" : "xlsx"}`);
      enlace.click();
      URL.revokeObjectURL(enlace.href);
      const filas = Number(respuesta.headers.get("x-reporte-filas") ?? 0);
      setMensaje({ tipo: "exito", texto: `Reporte generado con ${filas.toLocaleString("es-BO")} ${filas === 1 ? "registro" : "registros"}. La emisión quedó registrada en la auditoría.` });
    } catch {
      setMensaje({ tipo: "error", texto: "No se pudo conectar con el servidor. Revisa tu conexión e intenta de nuevo." });
    } finally {
      setCargando(false);
    }
  }

  // Pide a la API los filtros que se entienden de la frase y los carga en el formulario; no genera nada.
  async function interpretar() {
    setInterpretando(true);
    setNota(null);
    try {
      const respuesta = await fetch("/panel/reportes/interpretar", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ frase, categorias }),
      });
      const cuerpo = await respuesta.json().catch(() => null);
      if (!respuesta.ok) {
        const texto =
          respuesta.status === 503 || respuesta.status === 404
            ? "El asistente no está disponible por ahora. Usa el formulario de abajo."
            : cuerpo?.codigo === "PETICION_NO_ENTENDIDA"
              ? "No entendí qué reporte pides. Prueba con algo como «trámites observados de octubre en Excel»."
              : "No se pudo interpretar la frase. Intenta de nuevo.";
        setNota({ tipo: "error", texto });
        return;
      }
      setTipo(cuerpo.tipo);
      setFormato(cuerpo.formato);
      setValores({
        desde: cuerpo.desde ?? "",
        hasta: cuerpo.hasta ?? "",
        estado: cuerpo.estado ?? "",
        tipoTramite: cuerpo.tipoTramite ?? "",
        categoriaId: cuerpo.categoriaId ? String(cuerpo.categoriaId) : "",
        estadoLegajo: cuerpo.estadoLegajo ?? "",
      });
      setVersion((v) => v + 1);
      setNota({ tipo: "exito", texto: `${cuerpo.explicacion || "Listo."} Revisa los filtros y pulsa «Generar y descargar».` });
    } catch {
      setNota({ tipo: "error", texto: "No se pudo conectar con el servidor. Revisa tu conexión e intenta de nuevo." });
    } finally {
      setInterpretando(false);
    }
  }

  const conTramites = tipo === "TRAMITES" || tipo === "TIEMPOS";

  return (
    <form
      action="/panel/reportes/descargar"
      method="post"
      onSubmit={(e) => {
        e.preventDefault();
        void generar(e.currentTarget);
      }}
      className="flex flex-col gap-6"
    >
      <Tarjeta titulo="Descríbelo con tus palabras" descripcion="Opcional: escribe qué necesitas y el formulario se completa solo. Siempre puedes corregirlo antes de generar.">
        <div className="flex flex-col gap-3">
          <div className="flex flex-wrap gap-3">
            <input
              type="text"
              value={frase}
              maxLength={300}
              onChange={(e) => setFrase(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter") {
                  e.preventDefault(); // Enter interpreta; no envía el formulario
                  if (frase.trim() && !interpretando) void interpretar();
                }
              }}
              placeholder="Ej.: trámites observados de octubre en Excel"
              aria-label="Describe el reporte que necesitas"
              className={`${claseControl()} min-w-0 flex-[1_1_320px]`}
            />
            <BotonCarga type="button" cargando={interpretando} disabled={interpretando || !frase.trim()} variante="secundario" tamano="md" onClick={() => void interpretar()}>
              Completar filtros
            </BotonCarga>
          </div>
          {nota && <Alerta tipo={nota.tipo}>{nota.texto}</Alerta>}
        </div>
      </Tarjeta>

      <Tarjeta titulo="1. Tipo de reporte">
        <div className="grid gap-3 md:grid-cols-3">
          {TIPOS_REPORTE.map((t) => (
            <label key={t.valor} className={opcion(tipo === t.valor)}>
              <span className="flex items-center gap-2">
                <input type="radio" name="tipo" value={t.valor} checked={tipo === t.valor} onChange={() => setTipo(t.valor)} className="accent-[var(--color-accent)]" />
                <span className="text-base font-semibold text-ink">{t.etiqueta}</span>
              </span>
              <span className="text-sm leading-snug text-ink-soft">{t.descripcion}</span>
            </label>
          ))}
        </div>
      </Tarjeta>

      <Tarjeta
        titulo="2. Filtros"
        descripcion={
          tipo === "TRAMITES"
            ? "Las fechas se aplican a la radicación del trámite. Deja un campo vacío para no filtrar por él."
            : tipo === "TIEMPOS"
              ? "Las fechas se aplican a la resolución (remisión al Rectorado). Deja un campo vacío para no filtrar por él."
              : "El inventario siempre muestra la situación de hoy."
        }
      >
        <div key={version} className="grid items-end gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {conTramites && (
            <>
              <Campo id="desde" etiqueta={tipo === "TRAMITES" ? "Radicados desde" : "Resueltos desde"}>
                <input id="desde" type="date" name="desde" defaultValue={valores.desde} className={claseControl()} />
              </Campo>
              <Campo id="hasta" etiqueta={tipo === "TRAMITES" ? "Radicados hasta" : "Resueltos hasta"}>
                <input id="hasta" type="date" name="hasta" defaultValue={valores.hasta} className={claseControl()} />
              </Campo>
              {tipo === "TRAMITES" && (
                <Campo id="estado" etiqueta="Estado">
                  <select id="estado" name="estado" defaultValue={valores.estado ?? ""} className={claseControl()}>
                    <option value="">Todos los estados</option>
                    {FASES.map((e) => (
                      <option key={e} value={e}>
                        {ETIQUETA_ESTADO[e]}
                      </option>
                    ))}
                  </select>
                </Campo>
              )}
              <Campo id="tipoTramite" etiqueta="Tipo de trámite">
                <select id="tipoTramite" name="tipoTramite" defaultValue={valores.tipoTramite ?? ""} className={claseControl()}>
                  <option value="">Categorización y regularización</option>
                  {Object.entries(ETIQUETA_TIPO).map(([valor, etiqueta]) => (
                    <option key={valor} value={valor}>
                      {etiqueta}
                    </option>
                  ))}
                </select>
              </Campo>
              <Campo id="categoriaId" etiqueta="Categoría">
                <select id="categoriaId" name="categoriaId" defaultValue={valores.categoriaId ?? ""} disabled={categorias.length === 0} className={claseControl({ bloqueado: categorias.length === 0 })}>
                  <option value="">Todas las categorías</option>
                  {categorias.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.nombre}
                    </option>
                  ))}
                </select>
              </Campo>
            </>
          )}
          {tipo === "INVENTARIO" && (
            <Campo id="estadoLegajo" etiqueta="Estado de la carpeta">
              <select id="estadoLegajo" name="estadoLegajo" defaultValue={valores.estadoLegajo ?? ""} className={claseControl()}>
                <option value="">En archivo y prestadas</option>
                <option value="EN_ARCHIVO">Solo en archivo</option>
                <option value="PRESTADO">Solo prestadas</option>
              </select>
            </Campo>
          )}
        </div>
      </Tarjeta>

      <Tarjeta titulo="3. Formato">
        <div className="flex flex-col gap-5">
          <div className="grid gap-3 sm:max-w-[560px] sm:grid-cols-2">
            {FORMATOS.map((f) => (
              <label key={f.valor} className={opcion(formato === f.valor)}>
                <span className="flex items-center gap-2">
                  <input type="radio" name="formato" value={f.valor} checked={formato === f.valor} onChange={() => setFormato(f.valor)} className="accent-[var(--color-accent)]" />
                  <span className="text-base font-semibold text-ink">{f.etiqueta}</span>
                </span>
                <span className="text-sm text-ink-soft">{f.ayuda}</span>
              </label>
            ))}
          </div>

          {mensaje && <Alerta tipo={mensaje.tipo}>{mensaje.texto}</Alerta>}

          <BotonCarga type="submit" cargando={cargando} disabled={cargando} variante="primario" tamano="md">
            Generar y descargar
          </BotonCarga>
        </div>
      </Tarjeta>
    </form>
  );
}
