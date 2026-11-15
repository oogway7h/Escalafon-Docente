// Superficie con borde. nivel "h1" = la tarjeta es la página (formularios);
// "h2" = sección dentro de una página.
export default function Tarjeta({
  titulo,
  descripcion,
  nivel = "h2",
  relleno = true,
  className = "",
  children,
}: {
  titulo?: React.ReactNode;
  descripcion?: React.ReactNode;
  nivel?: "h1" | "h2";
  relleno?: boolean;
  className?: string;
  children: React.ReactNode;
}) {
  const Titulo = nivel;
  return (
    <section className={`rounded-xl border border-line bg-surface ${relleno ? "p-5" : "overflow-hidden"} ${nivel === "h1" ? "mx-auto" : ""} ${className}`}>
      {titulo && (
        <div className="mb-4">
          <Titulo className={nivel === "h1" ? "text-xl font-semibold text-ink" : "text-base font-semibold text-ink"}>{titulo}</Titulo>
          {descripcion && <p className="mt-1 text-sm text-ink-soft">{descripcion}</p>}
        </div>
      )}
      {children}
    </section>
  );
}
