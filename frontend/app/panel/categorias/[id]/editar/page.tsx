import { exigirPermiso } from "@/lib/api";
import { obtenerCategoria } from "@/lib/api-tramites";
import CategoriaForm from "../../CategoriaForm";
import { guardarCategoriaAction } from "../../actions";

export default async function EditarCategoriaPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("CATEGORIA_GESTIONAR");
  const categoria = await obtenerCategoria(id);
  return (
    <CategoriaForm
      accion={guardarCategoriaAction.bind(null, id)}
      categoria={categoria}
      titulo={`Editar categoría — ${categoria.nombre}`}
    />
  );
}
