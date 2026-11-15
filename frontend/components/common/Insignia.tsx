const tonos = {
  acento: "bg-accent-soft text-accent-ink",
  exito: "bg-success/10 text-success",
  neutro: "bg-bg text-ink-faint",
};

// Pastilla pequeña para rol, estado, etc.
export default function Insignia({
  tono = "acento",
  className = "",
  children,
}: {
  tono?: keyof typeof tonos;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <span className={`whitespace-nowrap rounded-full px-2.5 py-1 text-xs font-semibold ${tonos[tono]} ${className}`}>{children}</span>
  );
}
