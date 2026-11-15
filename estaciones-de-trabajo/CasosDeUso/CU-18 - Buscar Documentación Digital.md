# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-18 | Buscar Documentación Digital | P3 Gestión documental | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director, Analista, Personal de Archivo (con `DOCUMENTO_BUSCAR`)
**Actores secundarios:** —
**Objetivo:** Encontrar documentos por contenido, número de resolución, nombre o cédula del docente en menos de un segundo.

**Descripción:**
El usuario escribe palabras y opcionalmente elige buscar en todo, solo en el contenido o solo por docente, y un tipo de documento. El sistema busca con índices de texto completo sin tildes en los documentos activos, agrupa por documento, ordena por relevancia y muestra el docente, el trámite, las páginas halladas y el mejor fragmento con las coincidencias resaltadas. Al elegir una página, el visor se abre en ella.

## Precondiciones

- Sesión con `DOCUMENTO_BUSCAR`.
- Hay documentos procesados ([CU-15](CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR.md)).

## Flujo principal

1. El usuario abre el buscador.
2. Escribe la consulta y elige filtros.
3. El sistema busca en los documentos activos y devuelve los resultados paginados.
4. El usuario elige un resultado y una página; el visor se abre en esa página con la búsqueda resaltada.

## Flujos alternativos y excepciones

- **A1:** consulta vacía → 422 `CONSULTA_INVALIDA`.
- **A2:** campo o tipo desconocido → 422 `FILTRO_INVALIDO` / `TIPO_DOCUMENTO_INVALIDO`.
- **A3:** sin resultados → lista vacía con mensaje.

## Reglas de negocio

- Los documentos de baja nunca aparecen.
- Con dígitos en la búsqueda por docente, código y cédula son exactos.
- Índice `es_unaccent`: «resolucion» encuentra «resolución».

## Postcondiciones

- La consulta no modifica nada.

## Criterios de aceptación

- **CA01:** Dada una palabra del texto OCR, el sistema debe devolver el documento y la página.
- **CA02:** Dado un documento de baja que contiene la palabra, no debe aparecer.
- **CA03:** Dada una búsqueda sin tildes, deben aparecer coincidencias con tildes.
- **CA04:** La respuesta típica debe tardar menos de un segundo.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-18%20-%20Buscar%20Documentaci%C3%B3n%20Digital/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-18%20-%20Buscar%20Documentaci%C3%B3n%20Digital/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
