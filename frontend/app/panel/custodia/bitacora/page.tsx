import { claseControl } from "@/components/common/Campo";
import { BotonFiltrar, FiltrosCarga, ListadoConCarga, ZonaCarga } from "@/components/common/CargaListado";
import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import Insignia from "@/components/common/Insignia";
import Paginacion from "@/components/common/Paginacion";
import Tabla from "@/components/common/Tabla";
import { exigirPermiso } from "@/lib/api";
import { listarBitacora } from "@/lib/api-legajos";
import { ESTADOS_MATERIAL, formatearFecha } from "@/lib/legajos";

// CU-24: bitácora de custodia (solo lectura).
export default async function BitacoraPage({
  searchParams,
}: {
  searchParams: Promise<{ buscar?: string; tipo?: string; vencidos?: string; pagina?: string }>;
}) {
  await exigirPermiso("BITACORA_CONSULTAR");
  const sp = await searchParams;
  const pagina = Math.max(1, Number.parseInt(sp.pagina ?? "1", 10) || 1);
  const lista = await listarBitacora({ buscar: sp.buscar, tipo: sp.tipo, vencidos: sp.vencidos }, pagina);
  const totalPaginas = Math.max(1, Math.ceil(lista.total / lista.tamano));
  const material = Object.fromEntries(ESTADOS_MATERIAL.map((e) => [e.valor, e.etiqueta]));

  return (
    <ListadoConCarga>
      <div className="flex flex-col gap-6">
        <EncabezadoPagina titulo="Bitácora de custodia" descripcion={<>{lista.total} movimientos registrados.</>} />

        <FiltrosCarga action="/panel/custodia/bitacora" className="flex flex-wrap items-center gap-2.5">
          <input
            type="search"
            name="buscar"
            defaultValue={sp.buscar}
            placeholder="Código de legajo o de docente, cédula o nombre"
            className={`${claseControl()} min-w-64 max-w-[360px] flex-1`}
          />
          <select name="tipo" defaultValue={sp.tipo ?? ""} aria-label="Tipo de movimiento" className={claseControl()}>
            <option value="">Entradas y salidas</option>
            <option value="SALIDA">Salidas</option>
            <option value="ENTRADA">Entradas</option>
          </select>
          <label className="flex items-center gap-2 text-sm text-ink-soft">
            <input type="checkbox" name="vencidos" value="true" defaultChecked={sp.vencidos === "true"} />
            Solo vencidas
          </label>
          <BotonFiltrar>Filtrar</BotonFiltrar>
        </FiltrosCarga>

        <ZonaCarga>
          <Tabla
            filas={lista.datos}
            clave={(m) => m.id}
            vacio="No hay movimientos con esos filtros."
            columnas={[
              { titulo: "Fecha", celda: (m) => formatearFecha(m.fecha) },
              { titulo: "Carpeta", className: "w-28", celda: (m) => <span className="font-mono text-sm font-semibold">{m.legajoCodigo}</span> },
              {
                titulo: "Docente",
                celda: (m) => (
                  <>
                    <div className="font-semibold">{m.docenteNombre}</div>
                    {m.observacion && <div className="text-xs text-ink-faint">{m.observacion}</div>}
                  </>
                ),
              },
              {
                titulo: "Movimiento",
                celda: (m) => (
                  <div className="flex items-center gap-2">
                    <Insignia tono={m.tipo === "SALIDA" ? "acento" : "exito"}>{m.tipo === "SALIDA" ? "Salida" : "Entrada"}</Insignia>
                    {m.vencido && <Insignia tono="neutro">Vencida</Insignia>}
                  </div>
                ),
              },
              {
                titulo: "Custodio",
                celda: (m) => (
                  <>
                    <div>{m.custodio}</div>
                    <div className="text-xs text-ink-faint">Registró {m.registradoPor}</div>
                  </>
                ),
              },
              { titulo: "Límite", celda: (m) => formatearFecha(m.fechaLimite) },
              { titulo: "Estado material", celda: (m) => material[m.estadoMaterial] ?? m.estadoMaterial },
            ]}
          />
        </ZonaCarga>

        <Paginacion pagina={lista.pagina} totalPaginas={totalPaginas} searchParams={{ buscar: sp.buscar, tipo: sp.tipo, vencidos: sp.vencidos }} />
      </div>
    </ListadoConCarga>
  );
}
