// Indicador de carga: un anillo que gira (animación en globals.css, .animate-girar). Toma el color del texto.
// Ocupa el espacio del contenedor donde se renderiza (en el panel, solo el área de contenido).
// compacto: para superponerlo a una tabla ya difuminada (sin altura de página).
// mini: versión en línea (16 px) para dentro de un botón.
export default function Loader({ compacto = false, mini = false }: { compacto?: boolean; mini?: boolean }) {
  if (mini) {
    return (
      <span
        role="status"
        aria-label="Cargando"
        className="animate-girar inline-block h-4 w-4 shrink-0 rounded-full border-2 border-current border-t-transparent"
      />
    );
  }

  return (
    <div className={`flex w-full items-center justify-center text-ink-soft ${compacto ? "py-6" : "min-h-[60vh]"}`} role="status" aria-label="Cargando">
      <span className="animate-girar h-8 w-8 rounded-full border-2 border-current border-t-transparent" />
    </div>
  );
}
