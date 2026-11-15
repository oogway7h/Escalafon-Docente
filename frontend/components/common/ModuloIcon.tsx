const paths: Record<string, React.ReactNode> = {
  seguridad: (
    <path d="M9 2.2 14.5 4v4.6c0 3.6-2.4 6.2-5.5 7.2-3.1-1-5.5-3.6-5.5-7.2V4L9 2.2Z" />
  ),
  docente: (
    <>
      <circle cx="9" cy="6.4" r="2.6" />
      <path d="M3.8 15.4c0-2.6 2.3-4.4 5.2-4.4s5.2 1.8 5.2 4.4" />
    </>
  ),
  tramites: (
    <>
      <path d="M4.5 2.8h6.2L13.5 6v9.2h-9V2.8Z" />
      <path d="M6.8 8.8h4.4M6.8 11.6h4.4" />
    </>
  ),
  documental: (
    <>
      <path d="M2.6 5.6h4.2l1.4 1.6h6.2v7.4H2.6V5.6Z" />
    </>
  ),
  custodia: (
    <>
      <rect x="3" y="4.2" width="12" height="10.6" rx="1.4" />
      <path d="M6.2 4.2V3a1.4 1.4 0 0 1 1.4-1.4h3.8A1.4 1.4 0 0 1 12.8 3v1.2" />
    </>
  ),
  indicadores: (
    <>
      <path d="M3 15.2V3" />
      <path d="M3 15.2h12" />
      <path d="M5.6 12.4V9.6M9 12.4V6.4M12.4 12.4V8" />
    </>
  ),
};

export default function ModuloIcon({ tipo, className }: { tipo: keyof typeof paths; className?: string }) {
  return (
    <svg
      viewBox="0 0 18 18"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.5"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      {paths[tipo]}
    </svg>
  );
}
