import { LinkCarga } from "./CargaListado";

const claseEnlace = "rounded-lg border border-line px-3 py-1.5";

// Paginación por URL (?pagina=N): conserva el resto de los searchParams (filtros) al cambiar de página.
export default function Paginacion({
  pagina,
  totalPaginas,
  searchParams = {},
}: {
  pagina: number;
  totalPaginas: number;
  searchParams?: Record<string, string | undefined>;
}) {
  if (totalPaginas <= 1) return null;

  const enlaceA = (destino: number) => {
    const qs = new URLSearchParams();
    for (const [k, v] of Object.entries({ ...searchParams, pagina: String(destino) })) if (v) qs.set(k, v);
    return `?${qs}`;
  };

  const enlace = (destino: number, texto: string) =>
    destino < 1 || destino > totalPaginas ? (
      <span className={`${claseEnlace} opacity-40`} aria-disabled="true">
        {texto}
      </span>
    ) : (
      <LinkCarga href={enlaceA(destino)} className={`${claseEnlace} hover:border-accent hover:text-accent-ink`}>
        {texto}
      </LinkCarga>
    );

  return (
    <nav aria-label="Paginación" className="flex items-center justify-center gap-3 text-sm text-ink-soft">
      {enlace(pagina - 1, "Anterior")}
      <span>
        Página {pagina} de {totalPaginas}
      </span>
      {enlace(pagina + 1, "Siguiente")}
    </nav>
  );
}
