// Borde de 1 px; el foco lo marca el anillo global de :focus-visible (globals.css) más el borde en el acento.
const base = "rounded-lg border border-line text-ink transition-colors duration-150 placeholder:text-ink-faint focus:border-accent";

// Clases para <input>/<select>. grande = formularios de autenticación.
export function claseControl({ grande = false, bloqueado = false } = {}) {
  const tamano = grande ? "px-3.5 py-3 text-base" : "px-3.5 py-2.5 text-sm";
  return `${base} ${tamano} ${bloqueado ? "bg-bg text-ink-soft" : "bg-surface"}`;
}

// Etiqueta + control + ayuda + error de validación.
export default function Campo({
  id,
  etiqueta,
  opcional,
  ayuda,
  error,
  className = "",
  children,
}: {
  id?: string;
  etiqueta: string;
  opcional?: boolean;
  ayuda?: string;
  error?: string;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div className={`flex flex-col gap-1.5 ${className}`}>
      <label htmlFor={id} className="text-sm font-semibold text-ink">
        {etiqueta} {opcional && <span className="font-normal text-ink-faint">(opcional)</span>}
      </label>
      {children}
      {ayuda && <span className="text-xs text-ink-faint">{ayuda}</span>}
      {error && <p className="text-sm text-accent-ink">{error}</p>}
    </div>
  );
}
