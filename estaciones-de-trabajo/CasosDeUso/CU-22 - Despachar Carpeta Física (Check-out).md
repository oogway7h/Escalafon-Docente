# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-22 | Despachar Carpeta Física (Check-out) | P4 Archivo y custodia | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Personal de Archivo
**Actores secundarios:** —
**Objetivo:** Registrar la salida de una carpeta con custodio, fecha límite y estado material.

**Descripción:**
Sobre una solicitud aprobada (o pendiente, con «Aprobar y entregar»), el Personal de Archivo registra la salida: fecha límite de devolución, estado material y observación. El sistema marca el legajo como Prestado, la solicitud como Prestada y crea un asiento inmutable en la bitácora con fecha y hora, custodio receptor (el solicitante) y encargado de despacho.

## Precondiciones

- Sesión con `CUSTODIA_REGISTRAR` (y `LEGAJO_GESTIONAR` para aprobar y entregar en un paso).
- Solicitud Aprobada o Pendiente ([CU-21](CU-21%20-%20Atender%20Solicitud%20de%20Legajo.md)); legajo En archivo.

## Flujo principal

1. Archivo elige la solicitud.
2. Escribe la fecha límite (atajos +3/+7/+15 días), el estado material y una observación.
3. El sistema valida y, en una transacción, marca el legajo Prestado con un UPDATE condicional, la solicitud Prestada y crea el movimiento de salida.
4. El sistema audita.

## Flujos alternativos y excepciones

- **A1:** fecha límite pasada o inválida → 422 `FECHA_LIMITE_INVALIDA`.
- **A2:** estado material desconocido → 422 `ESTADO_MATERIAL_INVALIDO`.
- **A3:** legajo ya prestado → 409 `LEGAJO_NO_DISPONIBLE`.
- **A4:** solicitud en otro estado → 409 `SOLICITUD_ESTADO_INVALIDO`.

## Reglas de negocio

- Desvío de `cu.md`: sin evidencia fotográfica.
- `movimiento_custodia` es de solo anexión (trigger).
- Estados materiales: Buen estado, Hojas sueltas, Deteriorado.

## Postcondiciones

- El legajo está Prestado y la salida está en la bitácora.

## Criterios de aceptación

- **CA01:** Dada una solicitud aprobada y fecha futura, el legajo debe quedar Prestado y la salida en la bitácora.
- **CA02:** Dada una fecha pasada, el sistema debe rechazar el despacho.
- **CA03:** Dos despachos simultáneos del mismo legajo: solo uno debe tener éxito.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
