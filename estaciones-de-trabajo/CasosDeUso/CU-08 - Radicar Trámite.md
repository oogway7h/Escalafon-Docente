# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-08 | Radicar Trámite | P2 Trámites | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Personal de Ventanilla
**Actores secundarios:** —
**Objetivo:** Registrar un trámite de categorización o regularización con su código único de seguimiento.

**Descripción:**
El personal de ventanilla busca al docente, elige el tipo (categorización o regularización) y la categoría, y radica. El sistema genera un código único de seguimiento con formato `ED-AAAA-XXXXXX`, crea el trámite en estado Ingresado y registra la primera etapa del historial. También valida las solicitudes que llegan del portal público, que no pueden asignarse hasta ser validadas.

## Precondiciones

- Sesión con `TRAMITE_RADICAR`.
- El docente existe y está activo ([CU-07](CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas.md)).

## Flujo principal

1. Ventanilla abre «Radicar trámite» y busca al docente.
2. Elige tipo y categoría.
3. El sistema verifica docente activo y categoría vigente.
4. El sistema genera el código de seguimiento, crea el trámite en Ingresado y registra el historial y la auditoría.
5. El sistema muestra el código para entregarlo al docente.

## Flujos alternativos y excepciones

- **A1:** docente inactivo → 409 `DOCENTE_INACTIVO`.
- **A2:** categoría inexistente o dada de baja → 422 `CATEGORIA_INVALIDA`.
- **A3 (validar):** solicitud del portal ya validada → 409 `YA_VALIDADO`.
- **A4:** asignar una solicitud del portal sin validar → 409 `PENDIENTE_VALIDACION`.

## Reglas de negocio

- Código `ED-AAAA-XXXXXX` con alfabeto sin 0/O/1/I, generado con `crypto/rand`.
- El historial del trámite es de solo anexión.
- Tipo: CATEGORIZACION o REGULARIZACION.

## Postcondiciones

- El trámite existe en Ingresado con su código y la primera etapa.

## Criterios de aceptación

- **CA01:** Dado un docente activo y una categoría vigente, el sistema debe crear el trámite en Ingresado con un código único.
- **CA02:** Dado un docente inactivo, el sistema debe rechazar la radicación.
- **CA03:** Dada una solicitud del portal sin validar, el sistema debe impedir asignarla.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-08%20-%20Radicar%20Tr%C3%A1mite/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-08%20-%20Radicar%20Tr%C3%A1mite/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
