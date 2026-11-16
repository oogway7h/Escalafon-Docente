import Image from "next/image";
import Link from "next/link";
import BotonTema from "@/components/common/BotonTema";
import EscalatorMark from "@/components/common/EscalatorMark";
import ModuloIcon from "@/components/common/ModuloIcon";
import CategoriasRequisitos from "@/components/publico/CategoriasRequisitos";
import { listarCategoriasPublicas } from "@/lib/api-tramites";
import { ANTIGUEDAD_REGULARIZACION, PUNTOS_AUTORIDAD } from "@/lib/requisitos";
import logoDed from "@/public/logo-ded.png";

const modulos = [
  {
    tipo: "seguridad" as const,
    nombre: "Seguridad",
    titulo: "Roles, permisos y sesiones",
    descripcion:
      "Acceso por rol (Director, Analista, Personal de Archivo) con bloqueo por intentos fallidos y recuperación segura de contraseña.",
  },
  {
    tipo: "docente" as const,
    nombre: "Dominio",
    titulo: "Categorías y docentes",
    descripcion:
      "Registro y categorización de la planta docente conforme al escalafón vigente de la universidad.",
  },
  {
    tipo: "tramites" as const,
    nombre: "Trámites",
    titulo: "Radicación, asignación y revisión",
    descripcion:
      "Cada trámite queda asignado a un responsable, revisado y trazado de principio a fin.",
  },
  {
    tipo: "documental" as const,
    nombre: "Documental",
    titulo: "Gestión, restitución y búsqueda",
    descripcion:
      "Documentos digitalizados con restitución de originales y búsqueda inmediata por legajo o docente.",
  },
  {
    tipo: "custodia" as const,
    nombre: "Custodia",
    titulo: "Legajos, carpetas y bitácora",
    descripcion:
      "Control físico de legajos y carpetas con bitácora de movimientos y solicitudes.",
  },
  {
    tipo: "indicadores" as const,
    nombre: "Indicadores",
    titulo: "Tableros y reportes",
    descripcion:
      "Indicadores de gestión en tiempo real y reportes en PDF o Excel en segundos.",
  },
];

export default async function LandingPage() {
  // Categorías vigentes con sus documentos exigidos (las administra el Director). Si la API no responde, la sección se omite.
  const categorias = await listarCategoriasPublicas()
    .then((lista) => [...lista].sort((a, b) => (a.puntajeDesde ?? 0) - (b.puntajeDesde ?? 0) || a.nombre.localeCompare(b.nombre)))
    .catch(() => []);

  return (
    <div className="flex min-h-full flex-col bg-bg">
      <header className="sticky top-0 z-10 border-b border-line/80 bg-bg/75 backdrop-blur-md">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-3.5">
          <div className="flex items-center gap-3">
            <Image src={logoDed} alt="Departamento Escalafón Docente" className="h-8 w-auto" priority />
            <span className="font-display text-[15px] font-extrabold tracking-[-0.01em] text-ink">
              Escalafón
            </span>
          </div>
          <nav className="hidden items-center gap-7 text-sm font-semibold text-ink-soft md:flex">
            <Link href="/seguimiento" className="hover:text-ink">Consultar trámite</Link>
            <Link href="/solicitud" className="hover:text-ink">Nueva solicitud</Link>
            {categorias.length > 0 && <a href="#requisitos" className="hover:text-ink">Requisitos</a>}
            <a href="#modulos" className="hover:text-ink">Módulos</a>
            <a href="#acerca" className="hover:text-ink">Acerca del sistema</a>
          </nav>
          <div className="flex items-center gap-3">
            <BotonTema />
            <Link
              href="/login"
              className="rounded-[9px] bg-accent px-4 py-2 text-sm font-bold text-white transition-[transform,background] duration-100 hover:bg-accent-dark active:scale-[0.97]"
            >
              Iniciar sesión
            </Link>
          </div>
        </div>
      </header>

      <section
        className="sobre-oscuro relative overflow-hidden px-6 py-20"
        style={{
          background:
            "radial-gradient(900px 520px at 82% -10%, rgba(200,16,46,0.16), transparent 60%), linear-gradient(165deg, var(--color-sidebar) 0%, var(--color-sidebar-2) 100%)",
        }}
      >
        <div className="fade-salida mx-auto flex max-w-6xl flex-col items-center gap-10 text-center md:flex-row md:items-center md:gap-16 md:text-left">
          <div className="flex max-w-xl flex-col items-center gap-6 md:items-start">
            <h1 className="font-display text-[44px] font-extrabold leading-[1.08] tracking-[-0.01em] text-white text-balance">
              Un peldaño a la vez, sin perder el rastro de ningún trámite
            </h1>
            <p className="max-w-md text-[16px] leading-relaxed text-white/68">
              El sistema de gestión de flujos de trabajo y rastreo documental centraliza usuarios,
              trámites, documentos y legajos de la carrera docente en un solo lugar, con trazabilidad
              completa de cada paso.
            </p>
            <div className="flex flex-wrap items-center justify-center gap-3 md:justify-start">
              <Link
                href="/login"
                className="rounded-[10px] bg-accent px-6 py-3.5 text-[15px] font-bold text-white transition-[transform,background] duration-100 hover:bg-accent-dark active:scale-[0.975]"
              >
                Iniciar sesión
              </Link>
              <Link
                href="/seguimiento"
                className="rounded-[10px] border border-white/20 px-6 py-3.5 text-[15px] font-bold text-white/90 transition-colors duration-100 hover:bg-white/10"
              >
                Consultar mi trámite
              </Link>
            </div>
          </div>

          <div className="flex shrink-0 flex-col items-center gap-3">
            <EscalatorMark className="h-[190px] w-[220px]" />
            <p className="text-[12.5px] text-white/60">Categoría Titular · en revisión</p>
          </div>
        </div>
      </section>

      {categorias.length > 0 && (
        <section id="requisitos" className="px-6 pt-20">
          <div className="mx-auto flex max-w-6xl flex-col gap-10">
            <div className="fade-scroll flex flex-col gap-3">
              <h2 className="font-display text-[26px] font-extrabold tracking-[-0.01em] text-ink text-balance">
                ¿Qué documentos necesito para mi trámite?
              </h2>
              <p className="max-w-2xl text-[15px] leading-relaxed text-ink-soft">
                La categorización pide los mismos documentos para todas las categorías: lo que cambia es el puntaje que alcanzas. La regularización tiene su propia lista. Pasa el cursor por una tarjeta o tócala para ver los documentos; al tocarla queda fija hasta que la vuelvas a tocar. Toda la documentación se presenta en dos ejemplares, con sus originales para la autenticación.
              </p>
            </div>
            <div className="fade-scroll">
              <CategoriasRequisitos categorias={categorias} />
            </div>
            <div className="fade-scroll grid gap-6 rounded-[22px] border border-line bg-surface p-7 md:grid-cols-2">
              <div className="flex flex-col gap-2">
                <h3 className="font-display text-[18px] font-extrabold text-ink">Docentes en función de autoridad</h3>
                <p className="text-[14.5px] leading-relaxed text-ink-soft">
                  Se consideran de tiempo completo los docentes que ejercen una autoridad. Por cada año de gestión suman puntos de categoría extra, según el cargo:
                </p>
              </div>
              <ul className="flex flex-col divide-y divide-line">
                {PUNTOS_AUTORIDAD.map((a) => (
                  <li key={a.cargo} className="flex items-center justify-between gap-4 py-2.5 text-[14.5px]">
                    <span className="font-semibold text-ink">{a.cargo}</span>
                    <span className="shrink-0 font-mono font-medium text-accent-ink">{a.puntos} puntos por año</span>
                  </li>
                ))}
              </ul>
              <p className="text-[13px] text-ink-faint md:col-span-2">
                Regularización: es para el docente interino que acredita {ANTIGUEDAD_REGULARIZACION}. Consultas al Departamento de Escalafón Docente: teléfono 3393212.
              </p>
            </div>
          </div>
        </section>
      )}

      <section id="modulos" className="px-6 py-20">
        <div className="mx-auto max-w-6xl">
          <div className="fade-scroll mb-12 flex flex-col gap-3">
            <h2 className="font-display text-[30px] font-extrabold tracking-[-0.01em] text-ink text-balance">
              Seis módulos, un mismo expediente
            </h2>
          </div>
          <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {modulos.map((m, i) => (
              <div
                key={m.nombre}
                // En cada fila de tres, la 1.ª entra por la izquierda, la 2.ª por abajo y la 3.ª por la derecha.
                className={`${["fade-izq", "fade-scroll", "fade-der"][i % 3]} flex flex-col gap-4 rounded-2xl border border-line bg-surface p-6`}
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-[10px] bg-accent-soft text-accent-ink">
                  <ModuloIcon tipo={m.tipo} className="h-[18px] w-[18px]" />
                </div>
                <h3 className="font-display text-[17px] font-bold text-ink">{m.titulo}</h3>
                <p className="text-[14.5px] leading-relaxed text-ink-soft">{m.descripcion}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="acerca" className="px-6 pb-20">
        <div className="fade-scroll mx-auto flex max-w-6xl flex-col gap-4 rounded-2xl border border-line bg-surface p-10 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex flex-col gap-2">
            <h2 className="font-display text-[22px] font-extrabold tracking-[-0.01em] text-ink">
              ¿Tienes credenciales institucionales?
            </h2>
            <p className="max-w-md text-[14.5px] text-ink-soft">
              Ingresa con tu correo @uagrm.edu.bo para acceder al panel según los permisos de tu rol.
            </p>
          </div>
          <Link
            href="/login"
            className="w-fit rounded-[10px] bg-accent px-6 py-3.5 text-[15px] font-bold text-white transition-[transform,background] duration-100 hover:bg-accent-dark active:scale-[0.975]"
          >
            Iniciar sesión
          </Link>
        </div>
      </section>

      <footer className="mt-auto border-t border-line px-6 py-8">
        <div className="mx-auto flex max-w-6xl flex-col items-center justify-between gap-4 text-center sm:flex-row sm:text-left">
          <div className="flex items-center gap-2.5">
            <Image src={logoDed} alt="Departamento Escalafón Docente" className="h-6 w-auto" />
            <span className="text-[13px] text-ink-soft">
              Departamento de Escalafón Docente — Universidad Autónoma Gabriel René Moreno
            </span>
          </div>
          <span className="text-[13px] text-ink-faint">sistemas@uagrm.edu.bo</span>
        </div>
      </footer>
    </div>
  );
}
