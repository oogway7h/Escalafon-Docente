export function IconoMenuMovil({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" className={className} aria-hidden="true">
      <path d="M2.5 5h13M2.5 9h13M2.5 13h13" />
    </svg>
  );
}

export function IconoColapsar({ className, colapsado }: { className?: string; colapsado?: boolean }) {
  return (
    <svg
      viewBox="0 0 18 18"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.6"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      style={{ transform: colapsado ? "rotate(180deg)" : undefined, transition: "transform 150ms" }}
      aria-hidden="true"
    >
      <path d="M11.5 3.5 6 9l5.5 5.5" />
      <path d="M14 3.5v11" opacity="0.4" />
    </svg>
  );
}

export function IconoMas({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 18 18" fill="currentColor" className={className} aria-hidden="true">
      <circle cx="4" cy="9" r="1.3" />
      <circle cx="9" cy="9" r="1.3" />
      <circle cx="14" cy="9" r="1.3" />
    </svg>
  );
}

export function IconoEditar({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      <path d="M11.5 2.8 15.2 6.5 6.3 15.4H2.6v-3.7Z" />
    </svg>
  );
}

export function IconoEliminar({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      <path d="M3.5 5.2h11M7.3 5.2V3.6h3.4v1.6M6 5.2l.6 9.2a1 1 0 0 0 1 .95h2.8a1 1 0 0 0 1-.95l.6-9.2" />
    </svg>
  );
}

export function IconoSalir({ className }: { className?: string }) {
  return (
    <svg viewBox="0 0 18 18" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      <path d="M7 15.5H4a1.5 1.5 0 0 1-1.5-1.5V4A1.5 1.5 0 0 1 4 2.5h3" />
      <path d="M12 12.5 15.5 9 12 5.5" />
      <path d="M15.5 9h-9" />
    </svg>
  );
}
