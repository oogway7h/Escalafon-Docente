import Image from "next/image";
import Link from "next/link";
import BotonTema from "@/components/common/BotonTema";
import logoDed from "@/public/logo-ded.png";

// Marco del portal público: barra translúcida (el contenido se desliza por debajo) y pie sobrio.
export default function PublicoShell({ children }: { children: React.ReactNode }) {
  const enlace = "rounded-lg px-3 py-1.5 text-[14px] font-semibold text-ink-soft transition-colors hover:text-ink";
  return (
    <div className="flex min-h-full flex-col bg-bg">
      <header className="material sticky top-0 z-20 border-x-0 border-t-0">
        <div className="mx-auto flex max-w-5xl items-center justify-between gap-4 px-6 py-3">
          <Link href="/" className="flex items-center gap-3">
            <Image src={logoDed} alt="Departamento Escalafón Docente" className="h-8 w-auto" priority />
            <span className="hidden font-display text-[15px] font-extrabold tracking-[-0.01em] text-ink sm:inline">Escalafón</span>
          </Link>
          <nav className="flex items-center gap-1">
            <Link href="/seguimiento" className={enlace}>
              Seguimiento
            </Link>
            <Link href="/solicitud" className={enlace}>
              Nueva solicitud
            </Link>
            <BotonTema className="ml-2" />
          </nav>
        </div>
      </header>

      <main className="mx-auto w-full max-w-5xl flex-1 px-6 py-14 sm:py-20">{children}</main>

      <footer className="border-t border-line px-6 py-6 text-center text-[12.5px] text-ink-faint">
        Departamento de Escalafón Docente — Universidad Autónoma Gabriel René Moreno ·{" "}
        <Link href="/login" className="font-semibold hover:text-ink">
          Acceso del personal
        </Link>
      </footer>
    </div>
  );
}
