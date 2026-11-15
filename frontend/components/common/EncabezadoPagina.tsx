export default function EncabezadoPagina({
  titulo,
  descripcion,
  accion,
}: {
  titulo: React.ReactNode;
  descripcion?: React.ReactNode;
  accion?: React.ReactNode;
}) {
  return (
    <div className="flex items-end justify-between gap-4">
      <div>
        <h1 className="font-display text-xl font-semibold tracking-[-0.01em] text-ink">{titulo}</h1>
        {descripcion && <p className="text-sm text-ink-soft">{descripcion}</p>}
      </div>
      {accion}
    </div>
  );
}
