// Layout centrado de una sola tarjeta (activar cuenta, restablecer contraseña).
export default function AuthTarjeta({
  titulo,
  descripcion,
  children,
}: {
  titulo: string;
  descripcion: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex min-h-full items-center justify-center bg-bg px-6 py-14">
      <div className="w-full max-w-[380px] rounded-xl border border-line bg-surface p-8 sm:p-9">
        <div className="mb-6 flex h-10 w-10 items-center justify-center rounded-lg bg-sidebar font-display text-sm font-semibold text-white">
          ED
        </div>
        <h1 className="font-display text-xl font-semibold tracking-[-0.01em] text-ink">{titulo}</h1>
        <p className="mb-6 mt-1.5 text-sm text-ink-soft">{descripcion}</p>
        {children}
      </div>
    </div>
  );
}
