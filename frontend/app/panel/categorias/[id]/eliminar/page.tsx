import Link from "next/link";
import { claseBoton } from "@/components/common/Boton";
import Tarjeta from "@/components/common/Tarjeta";
import { exigirPermiso } from "@/lib/api";
import { obtenerCategoria } from "@/lib/api-tramites";
import { eliminarCategoriaAction } from "../../actions";
import EliminarCategoriaForm from "../../EliminarCategoriaForm";

export default async function EliminarCategoriaPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("CATEGORIA_GESTIONAR");
  const categoria = await obtenerCategoria(id);

  const enUso = categoria.docentes + categoria.tramites > 0;

  return (
    <Tarjeta
      nivel="h1"
      titulo={enUso ? `No se puede eliminar "${categoria.nombre}"` : `Eliminar "${categoria.nombre}"`}
      descripcion={
        enUso &&
        `La usan ${categoria.docentes} ${categoria.docentes === 1 ? "docente" : "docentes"} y ${categoria.tramites} ${categoria.tramites === 1 ? "trámite" : "trámites"}. Dala de baja para que no se ofrezca más: se conserva en su historial.`
      }
      className="max-w-md"
    >
      {enUso ? (
        <Link href="/panel/categorias" className={claseBoton("secundario")}>
          Entendido
        </Link>
      ) : (
        <EliminarCategoriaForm accion={eliminarCategoriaAction.bind(null, categoria.id)} nombre={categoria.nombre} />
      )}
    </Tarjeta>
  );
}
