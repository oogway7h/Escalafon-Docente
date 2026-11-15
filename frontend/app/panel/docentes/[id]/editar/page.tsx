import { exigirPermiso } from "@/lib/api";
import { listarCategorias, listarFacultades, obtenerDocente } from "@/lib/api-tramites";
import DocenteForm from "../../DocenteForm";
import { guardarDocenteAction } from "../../actions";

export default async function EditarDocentePage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  await exigirPermiso("DOCENTE_GESTIONAR");
  const [docente, categorias, facultades] = await Promise.all([obtenerDocente(id), listarCategorias(), listarFacultades()]);
  return (
    <DocenteForm
      accion={guardarDocenteAction.bind(null, id)}
      categorias={categorias}
      facultades={facultades}
      docente={docente}
      titulo={`Editar docente — ${docente.nombres} ${docente.apellidos}`}
    />
  );
}
