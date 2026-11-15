const variantes = {
  primario: "bg-accent text-white transition-[transform,background] duration-100 hover:bg-accent-dark active:scale-[0.975]",
  secundario: "border border-line text-ink transition-colors duration-100 hover:border-accent hover:text-accent-ink",
};

const tamanos = {
  sm: "w-fit rounded-lg px-4 py-2.5 text-sm",
  md: "w-fit rounded-lg px-6 py-2.5 text-sm",
  lg: "w-full rounded-lg py-3.5 text-base",
};

export type VarianteBoton = keyof typeof variantes;
export type TamanoBoton = keyof typeof tamanos;

// Clases de botón. Sin "use client" a propósito: se usa desde Server y Client Components
// con <button>, <Link> y <a>.
export function claseBoton(variante: VarianteBoton = "primario", tamano: TamanoBoton = "md") {
  return `inline-block text-center font-semibold disabled:cursor-wait disabled:opacity-70 ${variantes[variante]} ${tamanos[tamano]}`;
}
