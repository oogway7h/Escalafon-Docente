# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-19 | Gestionar Inventario de Legajos | P4 Archivo y custodia | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Personal de Archivo
**Actores secundarios:** —
**Objetivo:** Registrar y mantener el inventario de carpetas físicas con su código y ubicación precisa.

**Descripción:**
El Personal de Archivo lista los legajos, busca por código exacto o por docente y filtra por estado (En archivo o Prestado). Registra un legajo nuevo para un docente con su ubicación (gaveta, orden, anexos, observación); si deja el código vacío, el sistema asigna `LEG-NNNNN`. Edita la ubicación de un legajo existente. El estado no se edita a mano: lo cambian los movimientos de custodia.

## Precondiciones

- Sesión con `LEGAJO_GESTIONAR` (lectura también con `CARPETA_SOLICITAR` o `CUSTODIA_REGISTRAR`).
- El docente existe ([CU-07](CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas.md)).

## Flujo principal

1. El usuario abre Legajos y busca.
2. Elige «Nuevo», selecciona el docente y escribe la ubicación.
3. El sistema valida y guarda; asigna el código si va vacío.
4. Para editar, cambia la ubicación y guarda.

## Flujos alternativos y excepciones

- **A1:** código repetido → 409 `LEGAJO_CODIGO_DUPLICADO`.
- **A2:** docente inexistente → 422 `DOCENTE_INVALIDO`.
- **A3:** campos que superan su largo → 422 `VALIDACION`.

## Reglas de negocio

- Código único `LEG-NNNNN`, asignado por trigger si va vacío.
- Un docente puede tener varias carpetas.
- El estado cambia solo por despacho (CU-22) y devolución (CU-23).

## Postcondiciones

- El legajo está en el inventario con su ubicación y estado En archivo.

## Criterios de aceptación

- **CA01:** Dado un legajo sin código, el sistema debe asignar el siguiente `LEG-NNNNN`.
- **CA02:** Dado un código existente, el sistema debe rechazarlo.
- **CA03:** Dada una búsqueda por código exacto, debe aparecer solo ese legajo.
- **CA04:** El formulario no debe permitir cambiar el estado.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-19%20-%20Gestionar%20Inventario%20de%20Legajos/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-19%20-%20Gestionar%20Inventario%20de%20Legajos/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-19%20-%20Gestionar%20Inventario%20de%20Legajos/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
