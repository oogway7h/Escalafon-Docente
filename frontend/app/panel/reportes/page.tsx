import EncabezadoPagina from "@/components/common/EncabezadoPagina";
import { exigirPermiso } from "@/lib/api";
import { listarCategoriasPublicas } from "@/lib/api-tramites";
import { mensajeReporte, TIPOS_REPORTE, type TipoReporte } from "@/lib/indicadores";
import ReporteForm from "./ReporteForm";

// CU-26: generar y exportar reportes consolidados (PDF o Excel). La emisión queda en la auditoría.
export default async function ReportesPage({ searchParams }: { searchParams: Promise<{ tipo?: string; error?: string }> }) {
  await exigirPermiso("REPORTE_GENERAR");
  const sp = await searchParams;
  // Las categorías vigentes del portal: no exigen permisos de catálogo (el Personal de Archivo también genera reportes).
  const categorias = await listarCategoriasPublicas().catch(() => []);
  const tipoInicial: TipoReporte = TIPOS_REPORTE.find((t) => t.valor === sp.tipo)?.valor ?? "TRAMITES";

  return (
    <div className="flex flex-col gap-6">
      <EncabezadoPagina
        titulo="Reportes"
        descripcion="Reportes consolidados de trámites, inventario y tiempos de atención, en PDF para presentar o en Excel para analizar. Cada emisión queda registrada en la auditoría."
      />
      <ReporteForm
        categorias={categorias.map((c) => ({ id: c.id, nombre: c.nombre }))}
        tipoInicial={tipoInicial}
        errorInicial={sp.error ? mensajeReporte(sp.error) : undefined}
      />
    </div>
  );
}
