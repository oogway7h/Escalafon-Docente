import { FilaClicable } from "./EnlaceCarga";
import Tarjeta from "./Tarjeta";

export type Columna<T> = {
  titulo?: string;
  celda: (fila: T) => React.ReactNode;
  className?: string; // ancho/alineación de la columna, p. ej. "w-24 text-right"
};

// Tabla dentro de una tarjeta. Las columnas describen cómo pintar cada celda a partir de la fila.
export default function Tabla<T>({
  columnas,
  filas,
  clave,
  vacio = "No hay registros.",
  pie,
  className = "",
  href,
}: {
  columnas: Columna<T>[];
  filas: T[];
  clave: (fila: T) => React.Key;
  vacio?: string;
  pie?: React.ReactNode;
  className?: string;
  href?: (fila: T) => string; // si se indica, toda la fila navega a esa ruta al hacer clic
}) {
  return (
    <Tarjeta relleno={false} className={className}>
      {/* El desborde se queda dentro de la tarjeta: a 375px una tabla de 6 columnas no debe
          empujar scroll horizontal a toda la pagina. */}
      <div className="overflow-x-auto">
        <table className="w-full min-w-[640px] border-collapse">
        <thead>
          <tr className="text-left text-xs font-semibold uppercase tracking-wide text-ink-faint">
            {columnas.map((columna, i) => (
              <th key={i} className={`border-b border-line px-4 py-2.5 ${columna.className ?? ""}`}>
                {columna.titulo}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {filas.map((fila) => {
            const celdas = columnas.map((columna, i) => (
              <td key={i} className={`border-b border-line px-4 py-3 ${columna.className ?? ""}`}>
                {columna.celda(fila)}
              </td>
            ));
            return href ? (
              <FilaClicable key={clave(fila)} href={href(fila)}>
                {celdas}
              </FilaClicable>
            ) : (
              <tr key={clave(fila)} className="text-sm text-ink">
                {celdas}
              </tr>
            );
          })}
          {filas.length === 0 && (
            <tr>
              <td colSpan={columnas.length} className="px-4 py-8 text-center text-sm text-ink-faint">
                {vacio}
              </td>
            </tr>
          )}
        </tbody>
        </table>
      </div>
      {pie && <p className="px-4 py-3 text-xs text-ink-faint">{pie}</p>}
    </Tarjeta>
  );
}
