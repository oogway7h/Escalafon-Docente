const estilos = {
  error: "rounded-lg bg-accent-soft px-3.5 py-2.5 text-accent-ink",
  exito: "rounded-lg bg-success/10 px-3.5 py-2.5 text-success",
  aviso: "rounded-lg bg-accent-soft px-4 py-3.5 text-accent-dark dark:text-accent-ink",
};

export default function Alerta({
  tipo,
  className = "",
  children,
}: {
  tipo: keyof typeof estilos;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div role={tipo === "error" ? "alert" : "status"} className={`text-sm ${estilos[tipo]} ${className}`}>
      {children}
    </div>
  );
}
