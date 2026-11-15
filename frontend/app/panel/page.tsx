import Link from "next/link";
import { Suspense } from "react";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import ModuloIcon from "@/components/common/ModuloIcon";
import { obtenerSesion } from "@/lib/api";
import { listarSolicitudes } from "@/lib/api-legajos";
import { listarTramites } from "@/lib/api-tramites";
import { NAVEGACION, type SeccionNav } from "@/lib/navegacion";

type Modulo = SeccionNav["items"][number] & { seccion: SeccionNav["id"] };

const plural = (n: number, uno: string, varios: string) => `${n} ${n === 1 ? uno : varios}`;

/** Pendientes que esperan una acción de esta persona. `etiqueta` es lo que se lee al pasar el
 *  cursor y lo que anuncia el lector de pantalla: un número suelto no dice de qué es. */
function Insignia({ n, etiqueta }: { n: number; etiqueta: string }) {
  // Un cero no es una notificación: sin pendientes la tarjeta va limpia.
  if (n <= 0) return null;
  return (
    <span
      title={etiqueta}
      aria-label={etiqueta}
      className="rounded-full bg-accent px-2 py-0.5 text-xs font-semibold tabular-nums text-white transition-transform duration-100 group-hover:scale-110"
    >
      {n > 99 ? "99+" : n}
    </span>
  );
}

// Cada insignia se pide por su cuenta y va envuelta en <Suspense>: así la portada se pinta de
// inmediato y el número llega después, en vez de retrasar toda la página una ida y vuelta a la
// API. Importa porque /panel no tiene loading.tsx propio (se lo lleva un ancestro y volvería a
// cubrir a todos sus hijos), de modo que esperar aquí deja la pantalla anterior a la vista.
// Si la consulta falla no hay insignia: la portada tiene que funcionar con la API de legajos caída.
async function InsigniaBandeja() {
  const n = await listarTramites({ mios: true })
    .then((ts) => ts.filter((t) => t.estado !== "ARCHIVADO" && t.estado !== "REMITIDO").length)
    .catch(() => 0);
  return <Insignia n={n} etiqueta={`${plural(n, "expediente asignado", "expedientes asignados")} sin remitir`} />;
}

async function InsigniaCustodia() {
  // tamaño 1: solo interesa el total que devuelve la paginación, no las filas.
  const n = await listarSolicitudes({ estado: "PENDIENTE" }, 1, 1)
    .then((r) => r.total)
    .catch(() => 0);
  return <Insignia n={n} etiqueta={`${plural(n, "solicitud", "solicitudes")} de carpeta sin atender`} />;
}

// Una tarjeta por pantalla. Es un enlace entero, no un título enlazado: el objetivo táctil es
// toda la tarjeta.
function TarjetaModulo({ modulo, insignia }: { modulo: Modulo; insignia?: React.ReactNode }) {
  if (!modulo.disponible) {
    // Mismo trato que en la barra lateral: el permiso existe en el seed pero la pantalla
    // todavía no, así que se ve apagada y no se puede pulsar.
    return (
      <div className="flex flex-1 items-start gap-3.5 rounded-xl border border-dashed border-line p-5">
        <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-bg text-ink-faint">
          <ModuloIcon tipo={modulo.seccion} className="h-[18px] w-[18px]" />
        </span>
        <span className="flex min-w-0 flex-col gap-1">
          <span className="flex items-center gap-2">
            <span className="font-display text-base font-semibold text-ink-soft">{modulo.etiqueta}</span>
            <span className="rounded-full bg-bg px-1.5 py-0.5 text-xs font-semibold uppercase tracking-wide text-ink-faint">
              Pronto
            </span>
          </span>
          <span className="text-sm leading-relaxed text-ink-faint">{modulo.descripcion}</span>
        </span>
      </div>
    );
  }

  return (
    <Link
      href={modulo.href}
      className="group flex flex-1 items-start gap-3.5 rounded-xl border border-line bg-surface p-5 transition-[border-color,transform] duration-100 hover:border-accent active:scale-[0.995]"
    >
      <span className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent-soft text-accent-ink">
        <ModuloIcon tipo={modulo.seccion} className="h-[18px] w-[18px]" />
      </span>
      <span className="flex min-w-0 flex-1 flex-col gap-1">
        <span className="flex items-center gap-2">
          <span className="font-display text-base font-semibold text-ink">{modulo.etiqueta}</span>
          {insignia}
        </span>
        <span className="text-sm leading-relaxed text-ink-soft">{modulo.descripcion}</span>
      </span>
    </Link>
  );
}

function Rejilla({ modulos, insignias }: { modulos: Modulo[]; insignias: Record<string, React.ReactNode> }) {
  return (
    <ul className="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
      {modulos.map((modulo) => (
        <li key={modulo.href} className="flex">
          <TarjetaModulo modulo={modulo} insignia={insignias[modulo.href]} />
        </li>
      ))}
    </ul>
  );
}

// Portada del panel: lo que este rol puede hacer, con una línea que explica para qué sirve cada
// pantalla (la barra lateral solo tiene la etiqueta). Sale de NAVEGACION y de los permisos de la
// sesión, así que un rol nuevo, o un permiso quitado, se refleja aquí sin tocar esta pantalla.
export default async function PanelInicio() {
  const sesion = await obtenerSesion();
  const permisos = new Set(sesion.permisos);
  // Solo las dos pantallas que son una cola de trabajo llevan contador.
  const insignias: Record<string, React.ReactNode> = {};
  if (permisos.has("TRAMITE_REVISAR")) {
    insignias["/panel/tramites/revisar"] = (
      <Suspense fallback={null}>
        <InsigniaBandeja />
      </Suspense>
    );
  }
  if (permisos.has("CUSTODIA_REGISTRAR")) {
    insignias["/panel/custodia"] = (
      <Suspense fallback={null}>
        <InsigniaCustodia />
      </Suspense>
    );
  }
  const secciones = NAVEGACION.map((seccion) => ({
    ...seccion,
    items: seccion.items
      .filter((item) => permisos.has(item.permiso))
      .map((item): Modulo => ({ ...item, seccion: seccion.id })),
  })).filter((seccion) => seccion.items.length > 0);

  const modulos = secciones.flatMap((seccion) => seccion.items);

  // Agrupar por sección solo cuando las secciones tienen cuerpo: de media, dos pantallas o más.
  // Un encabezado sobre una única tarjeta a un tercio de ancho es ruido, y los roles acotados
  // caen justo ahí (el Analista ve cinco pantallas repartidas en cinco secciones, y el Personal
  // de Archivo nueve en cinco). El Director ve veintiuna en seis y ahí la agrupación sí orienta.
  const agrupar = modulos.length >= secciones.length * 2;

  return (
    <div className="flex flex-col gap-8">
      <EncabezadoPagina
        titulo={`Hola, ${sesion.usuario.nombres}`}
        descripcion={
          <>
            Sesión iniciada como <span className="font-semibold text-ink">{sesion.rol}</span>.{" "}
            {modulos.length === 0
              ? "Tu rol todavía no tiene módulos asignados; habla con el Director."
              : `Esto es lo que tu rol puede hacer (${modulos.length} ${modulos.length === 1 ? "pantalla" : "pantallas"}).`}
          </>
        }
      />

      {agrupar ? (
        secciones.map((seccion) => (
          <section key={seccion.id} className="flex flex-col gap-3.5">
            <h2 className="text-xs font-semibold uppercase tracking-[0.1em] text-ink-faint">{seccion.etiqueta}</h2>
            <Rejilla modulos={seccion.items} insignias={insignias} />
          </section>
        ))
      ) : (
        <Rejilla modulos={modulos} insignias={insignias} />
      )}
    </div>
  );
}
