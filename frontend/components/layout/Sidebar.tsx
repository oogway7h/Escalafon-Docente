"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import ModuloIcon from "@/components/common/ModuloIcon";
import type { Sesion } from "@/lib/api";
import { NAVEGACION } from "@/lib/navegacion";
import logoDed from "@/public/logo-ded.png";

function ContenidoSidebar({
  sesion,
  colapsado,
  pathname,
  onNavegar,
}: {
  sesion: Sesion;
  colapsado: boolean;
  pathname: string;
  onNavegar?: () => void;
}) {
  const permisos = new Set(sesion.permisos);
  const secciones = NAVEGACION.map((seccion) => ({
    ...seccion,
    items: seccion.items.filter((item) => permisos.has(item.permiso)),
  })).filter((seccion) => seccion.items.length > 0);

  return (
    <>
      {/* El logo ya lleva el nombre del Departamento; colapsado se muestra entero, escalado al ancho. */}
      <div className={`flex items-center px-4 py-5 ${colapsado ? "justify-center px-0" : ""}`}>
        <Image
          src={logoDed}
          alt="Departamento Escalafón Docente"
          priority
          className={colapsado ? "h-auto w-10 rounded-[5px]" : "h-12 w-auto rounded-lg"}
        />
      </div>

      <nav className="scrollbar-hide flex-1 overflow-y-auto px-2.5 pb-4">
        {secciones.map((seccion) => (
          <div key={seccion.id} className="mb-5">
            {!colapsado && (
              <p className="px-2.5 pb-1.5 text-xs text-white/60">
                {seccion.etiqueta}
              </p>
            )}
            <ul className="flex flex-col gap-0.5">
              {seccion.items.map((item) => {
                const coincide = (href: string) => pathname === href || pathname.startsWith(`${href}/`);
                // Gana el ítem más específico (/panel/custodia/bitacora sobre /panel/custodia).
                const activo = coincide(item.href) && !NAVEGACION.some((s) => s.items.some((o) => o.href.length > item.href.length && coincide(o.href)));
                return (
                  <li key={item.href}>
                    <Link
                      href={item.href}
                      onClick={onNavegar}
                      title={colapsado ? item.etiqueta : undefined}
                      className={`flex items-center gap-2.5 rounded-lg px-2.5 py-2 text-sm font-semibold transition-colors duration-100 ${colapsado ? "justify-center" : ""} ${
                        activo ? "bg-white/12 text-white" : "text-white/64 hover:bg-white/8 hover:text-white"
                      }`}
                    >
                      <ModuloIcon tipo={seccion.id} className="h-[17px] w-[17px] shrink-0" />
                      {!colapsado && item.etiqueta}
                    </Link>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>

      <div className={`border-t border-white/10 px-4 py-4 ${colapsado ? "px-0 text-center" : ""}`}>
        {!colapsado ? (
          <>
            <p className="truncate text-sm font-semibold text-white">
              {sesion.usuario.nombres} {sesion.usuario.apellidos}
            </p>
            <p className="text-xs text-white/60">{sesion.rol}</p>
          </>
        ) : (
          <div className="mx-auto flex h-7 w-7 items-center justify-center rounded-full bg-white/10 text-xs font-semibold text-white">
            {sesion.usuario.nombres[0]}
          </div>
        )}
      </div>
    </>
  );
}

export default function Sidebar({
  sesion,
  colapsado,
  abiertoMovil,
  onCerrarMovil,
}: {
  sesion: Sesion;
  colapsado: boolean;
  abiertoMovil: boolean;
  onCerrarMovil: () => void;
}) {
  const pathname = usePathname();
  const fondo = "bg-sidebar";

  return (
    <>
      <aside
        className={`sticky top-0 hidden h-screen shrink-0 flex-col ${fondo} transition-[width] duration-150 md:flex ${
          colapsado ? "w-[68px]" : "w-[248px]"
        }`}
      >
        <ContenidoSidebar sesion={sesion} colapsado={colapsado} pathname={pathname} />
      </aside>

      {abiertoMovil && (
        <div className="fixed inset-0 z-40 md:hidden">
          <div className="absolute inset-0 bg-black/40" onClick={onCerrarMovil} />
          <aside className={`absolute inset-y-0 left-0 flex w-[248px] flex-col ${fondo}`}>
            <ContenidoSidebar sesion={sesion} colapsado={false} pathname={pathname} onNavegar={onCerrarMovil} />
          </aside>
        </div>
      )}
    </>
  );
}
