import Image from "next/image";
import Link from "next/link";
import BotonTema from "@/components/common/BotonTema";
import EscalatorMark from "@/components/common/EscalatorMark";
import logoDed from "@/public/logo-ded.png";

// Dos columnas: panel institucional a la izquierda, formulario a la derecha
// (iniciar sesión, recuperar contraseña).
export default function AuthDividido({
  heroTitulo,
  heroTexto,
  heroDetalle,
  titulo,
  descripcion,
  pie,
  children,
}: {
  heroTitulo: string;
  heroTexto: string;
  heroDetalle?: React.ReactNode;
  titulo: string;
  descripcion: string;
  pie: React.ReactNode;
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-full flex-col md:flex-row">
      <div
        className="sobre-oscuro relative flex flex-col justify-between gap-10 overflow-hidden px-8 py-12 md:w-[560px] md:shrink-0 md:px-14 md:py-14"
        style={{
          background:
            "radial-gradient(720px 480px at 88% -8%, rgba(200,16,46,0.20), transparent 60%), linear-gradient(165deg, var(--color-sidebar) 0%, var(--color-sidebar-2) 100%)",
        }}
      >
        <div className="flex items-center gap-3.5">
          <Link
            href="/"
            aria-label="Ir al inicio"
            className="flex items-center rounded-[14px] bg-white/96 px-3.5 py-2.5 transition-transform duration-100 hover:scale-[1.02] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
          >
            <Image src={logoDed} alt="Departamento Escalafón Docente" className="h-11 w-auto" priority />
          </Link>
        </div>

        <div className="flex flex-col gap-7">
          <h1 className="max-w-[400px] font-display text-[40px] font-extrabold leading-[1.14] tracking-[-0.01em] text-white text-balance">
            {heroTitulo}
          </h1>
          <p className="max-w-[380px] font-body text-[15.5px] leading-relaxed text-white/66">{heroTexto}</p>

          <div className="mt-2 flex items-center gap-5">
            <EscalatorMark className="h-[92px] w-[108px] shrink-0" />
            {heroDetalle}
          </div>
        </div>

        <p className="text-[12.5px] text-white/60">
          Departamento de Escalafón Docente — Universidad Autónoma Gabriel René Moreno
        </p>
      </div>

      <div className="relative flex flex-1 items-center justify-center px-8 py-14 md:px-14">
        <BotonTema className="absolute right-6 top-6" />
        <div className="flex w-full max-w-[380px] flex-col gap-7">
          <div className="flex flex-col gap-2">
            <h2 className="font-display text-[28px] font-extrabold tracking-[-0.01em] text-ink">{titulo}</h2>
            <p className="text-[15px] text-ink-soft">{descripcion}</p>
          </div>

          {children}

          <div className="border-t border-line pt-[18px]">
            <p className="text-[13px] text-ink-soft">{pie}</p>
          </div>
        </div>
      </div>
    </div>
  );
}
