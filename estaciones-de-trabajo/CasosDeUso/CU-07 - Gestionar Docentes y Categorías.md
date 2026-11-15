# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-07 | Gestionar Docentes y Categorías | P2 Trámites | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director / Analista con `DOCENTE_GESTIONAR` o `CATEGORIA_GESTIONAR`
**Actores secundarios:** —
**Objetivo:** Mantener el padrón docente y el catálogo normativo de categorías con sus requisitos.

**Descripción:**
El usuario administra el padrón: busca docentes por código, cédula o nombre, registra uno nuevo (código universitario obligatorio y único, cédula opcional pero inmutable una vez asignada), lo edita o lo da de baja lógica. También administra las categorías del escalafón: nombre, años de antigüedad exigidos y documentos probatorios. Una categoría se puede dar de baja (se oculta) o eliminar solo si ningún docente ni trámite la usa.

## Precondiciones

- Sesión con `DOCENTE_GESTIONAR` o `CATEGORIA_GESTIONAR`.
- El padrón inicial se cargó con el importador (`make import-docentes`).

## Flujo principal

1. El usuario abre Docentes y busca.
2. Registra o edita un docente con código, nombres, apellidos, categoría, facultad y contacto.
3. El sistema valida unicidad del código y de la cédula y guarda con auditoría.
4. En Categorías, define nombre, antigüedad y requisitos y guarda.
5. Da de baja o reactiva docentes y categorías según corresponda.

## Flujos alternativos y excepciones

- **A1:** código repetido → 409 `CODIGO_DUPLICADO`.
- **A2:** cédula repetida → 409 `CEDULA_DUPLICADA`.
- **A3:** cambiar una cédula ya asignada → 422 `CAMPO_NO_EDITABLE`.
- **A4:** eliminar una categoría en uso → 409 `CATEGORIA_EN_USO` con los conteos.
- **A5:** categoría con nombre repetido (sin tildes ni prefijo «Profesor») → 409 `CATEGORIA_DUPLICADA`.

## Reglas de negocio

- Desvío de `cu.md`: el identificador principal es el **código universitario**; la cédula es opcional e inmutable una vez asignada.
- Bajas lógicas; eliminar una categoría solo sin docentes ni trámites asociados.
- Búsqueda con dígitos = código o cédula exactos; sin dígitos = fragmento del nombre.
- Los cargos (Decano, Rector…) existen como categorías inactivas.

## Postcondiciones

- Padrón y catálogo actualizados y auditados; los requisitos cambian en el portal público al instante.

## Criterios de aceptación

- **CA01:** Dado un código ya existente, el sistema debe rechazar el alta.
- **CA02:** Dada una cédula asignada, el sistema debe impedir cambiarla.
- **CA03:** Dada una categoría con trámites, el sistema debe impedir eliminarla y ofrecer darla de baja.
- **CA04:** Dado «123» en el buscador, el sistema no debe traer el código «1234».

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
