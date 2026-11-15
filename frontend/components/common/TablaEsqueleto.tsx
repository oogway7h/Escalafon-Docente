import Tarjeta from "./Tarjeta";

// Silueta de la tabla que está por llegar. Reproduce la geometría de Tabla (mismo padding de
// celda, mismo borde entre filas, mismo min-w) y, sobre todo, su altura de fila: las catorce
// tablas del panel tienen al menos una columna de dos líneas (nombre + correo, categoría +
// facultad…), así que la fila mide ~64 px y no ~38. Con la silueta corta, la página crecía de
// golpe al llegar los datos.
//
// Los anchos varían por columna: una rejilla de barras idénticas se lee como un patrón, no como
// contenido. El ciclo arranca ancho porque la primera columna suele ser el nombre.
const ANCHOS = ["70%", "44%", "32%", "54%", "38%", "48%"];

export default function TablaEsqueleto({ columnas = 4, filas = 8 }: { columnas?: number; filas?: number }) {
  return (
    <Tarjeta relleno={false}>
      {/* aria-hidden: el estado de carga ya lo anuncia el role="status" del contenedor. */}
      <div className="overflow-x-auto" aria-hidden>
        <table className="w-full min-w-[640px] border-collapse">
          <thead>
            <tr>
              {Array.from({ length: columnas }, (_, c) => (
                <th key={c} className="border-b border-line px-4 py-2.5">
                  {/* min-h: cuadra el encabezado con los 36.9 px del real (py-2.5 + 16 + borde). */}
                  <span className="flex min-h-[16px] items-center">
                    <span className="esqueleto block h-2.5 w-14 rounded-full" />
                  </span>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {Array.from({ length: filas }, (_, f) => (
              <tr key={f}>
                {Array.from({ length: columnas }, (_, c) => (
                  <td key={c} className="border-b border-line px-4 py-3">
                    <span
                      className="flex min-h-[40px] flex-col justify-center gap-2"
                      // El desfase por fila evita que todas las barras parpadeen al unísono,
                      // que es lo que hace que un esqueleto parezca una pantalla rota.
                      style={{ animationDelay: `${f * 90}ms` }}
                    >
                      <span
                        className="esqueleto block h-3.5 rounded-full"
                        style={{ width: ANCHOS[(f + c) % ANCHOS.length], animationDelay: `${f * 90}ms` }}
                      />
                      {/* Segunda línea solo en la primera columna, como en las tablas reales. */}
                      {c === 0 && (
                        <span
                          className="esqueleto block h-2.5 w-[46%] rounded-full"
                          style={{ animationDelay: `${f * 90}ms` }}
                        />
                      )}
                    </span>
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Tarjeta>
  );
}
