# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-05 | Gestionar Roles y Permisos | P1 Seguridad | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director
**Actores secundarios:** —
**Objetivo:** Definir roles, asociarles permisos granulares y protegerlos de eliminaciones peligrosas.

**Descripción:**
El Director lista los roles, crea uno nuevo con nombre y permisos del catálogo (cada permiso corresponde a un caso de uso), lo edita o lo elimina. Un rol con usuarios asignados no puede eliminarse, y los roles de sistema no se eliminan ni se renombran. El Director no puede quitarse los permisos de usuarios y roles.

## Precondiciones

- Sesión con `ROL_GESTIONAR`.
- El catálogo de permisos fue cargado por el seed.

## Flujo principal

1. El Director abre Roles.
2. Elige crear o editar; escribe el nombre y marca permisos.
3. El sistema valida nombre único y permisos existentes y guarda.
4. Para eliminar, confirma; el sistema verifica que no haya usuarios con ese rol.
5. Toda operación queda en auditoría.

## Flujos alternativos y excepciones

- **A1 (paso 4):** rol con usuarios → 409 `ROL_CON_USUARIOS`.
- **A2 (paso 3):** nombre repetido → 409 `ROL_DUPLICADO`.
- **A3:** rol de sistema → 409 `ROL_SISTEMA_PROTEGIDO`.
- **A4:** permiso inexistente → 422 `PERMISO_INVALIDO`.

## Reglas de negocio

- Roles de sistema: Director, Analista, Personal de Archivo, Personal de Ventanilla.
- El Director conserva siempre `USUARIO_GESTIONAR` y `ROL_GESTIONAR`.
- La asignación de un rol a un usuario se hace en CU-04.

## Postcondiciones

- El rol existe con sus permisos y la operación está auditada.

## Criterios de aceptación

- **CA01:** Dado un rol con usuarios, el sistema debe impedir eliminarlo.
- **CA02:** Dado un rol de sistema, el sistema debe impedir renombrarlo o eliminarlo.
- **CA03:** Dado un rol guardado con nuevos permisos, sus usuarios deben ver el menú nuevo al volver a iniciar sesión.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-05%20-%20Gestionar%20Roles%20y%20Permisos/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-05%20-%20Gestionar%20Roles%20y%20Permisos/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
