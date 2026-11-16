// Información normativa del Departamento de Escalafón Docente que no vive en el catálogo de categorías:
// los requisitos de la regularización y los puntos por función de autoridad. Fuentes: «Requisitos para
// regularización docente» y art. 12 del Reglamento del Escalafón Docente. Se usa en el portal público.

export const REQUISITOS_REGULARIZACION = [
  "Carta «Solicitud de Regularización» (3 ejemplares) al Jefe del Dpto. de Escalafón Docente, con la facultad",
  "Carta a la DAGA solicitando la carga horaria",
  "Formulario de Consanguinidad (original, RR.HH. – UAGRM)",
  "Certificado de Trabajo original de RR.HH., actualizado a la fecha de presentación",
  "Declaración jurada de Bienes y Rentas (copia del original, Contraloría General del Estado)",
  "Resolución Facultativa que te declara ganador de la materia o acredita la adjudicación de la cátedra",
  "Examen de competencia de la materia, llenado por el tribunal evaluador",
  "Acta de conformidad del tribunal evaluador",
  "Las 3 publicaciones de la convocatoria (externa) o la resolución facultativa si fue interna o interna ampliada",
  "Dos fotos 4x4 (cualquier fondo) y currículum personal",
  "Fotocopias legalizadas: título en provisión nacional, título académico, diplomado y maestría",
  "Fotocopia simple de la cédula de identidad",
  "Certificados de cursos y otros méritos",
];

// Requisito de antigüedad de la regularización (docente interino que pasa a ordinario o de planta).
export const ANTIGUEDAD_REGULARIZACION = "2 años o 4 semestres sin cortes en la materia que postuló y obtuvo por concurso de méritos o examen de competencia";

// Art. 12: el docente en función de autoridad (tiempo completo) suma puntos de categoría por cada año de gestión.
export const PUNTOS_AUTORIDAD = [
  { cargo: "Rector y Vicerrector", puntos: 50 },
  { cargo: "Decanos y Directores", puntos: 40 },
  { cargo: "Jefes de Carrera y Subdecanos (vicedecanos)", puntos: 30 },
];

// Art. 18: «30 a 89 puntos»; la última categoría no tiene tope («450 o más»).
export function textoPuntaje(desde: number | null, hasta: number | null): string | null {
  if (desde === null) return null;
  return hasta === null ? `${desde}+` : `${desde}–${hasta}`;
}
