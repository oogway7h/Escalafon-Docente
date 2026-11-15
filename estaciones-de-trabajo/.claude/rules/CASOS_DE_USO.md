---
globs: CasosDeUso/**/*
---

# Casos de Uso

Define **cómo se escribe un Caso de Uso (CU)**: formato, secciones obligatorias y granularidad. Aplica a
todo archivo de `CasosDeUso/`. Es la **especificación** del ciclo SDD: nada se implementa sin un CU
escrito y completo.

La numeración (`CU-01` … `CU-26`) viene de `archivos/docsmd/cu.md` (tabla de la cátedra): es la fuente
de verdad de qué CU existe y cómo se llama. No se renumera ni se reasigna.

## Dónde vive y cómo se nombra

`CasosDeUso/CU-{NN} - {Nombre corto}.md`. El nombre coincide con el de la tabla de `cu.md`.

**Desvíos conscientes respecto de `cu.md`** (cada CU afectado lo repite en sus Reglas de negocio):

- **CU-01:** el inicio de sesión es solo por correo institucional (no por código de funcionario).
- **CU-11:** el analista dictamina y remite; **archiva el Personal de Archivo** (flujo real del
  Departamento, con el tramo «En facultad» sin seguimiento).
- **CU-15:** el límite de 20 MB rige para imágenes; los PDF admiten hasta 400 MB (escaneos de carpetas completas).
- **CU-22 / CU-23:** la salida y la devolución no llevan evidencia fotográfica (se quitó a pedido del Departamento).

## Granularidad: un CU, un flujo

Si una funcionalidad agrupa flujos con disparadores y criterios de aceptación independientes, son CU
distintos. Una dependencia entre CU se documenta en **precondiciones** con link, no fusionándolos.

## Plantilla

```markdown
# Caso de uso

| Id    | Nombre   | Paquete   | Prioridad       | Estado                      |
| ----- | -------- | --------- | --------------- | --------------------------- |
| CU-NN | {Nombre} | {Paquete} | Alta/Media/Baja | Pendiente/En progreso/Hecho |

**Plataformas:** Web (Hecho) · Mobile (Pendiente | Hecho | —)
**Actor principal:** …
**Actores secundarios:** … o "—"
**Objetivo:** una frase.

**Descripción:** párrafo con el camino feliz y los errores principales.

## Precondiciones
## Flujo principal                     (numerado)
## Flujos alternativos y excepciones   (A1 (paso N): condición → respuesta, con el código de error de la API)
## Reglas de negocio
## Postcondiciones
## Criterios de aceptación             (CA01, CA02… "Dado que …, el sistema debe …")
## Desarrollador
```

Paquetes: **P1 Seguridad** (CU-01…06), **P2 Trámites** (CU-07…14), **P3 Gestión documental**
(CU-15…18), **P4 Archivo y custodia** (CU-19…24), **P5 Indicadores** (CU-25, CU-26).

## Reglas de contenido

- Actor, Objetivo y Descripción obligatorios; flujos alternativos referenciados al paso que bifurcan.
- Reglas de negocio explícitas, nunca implícitas en la prosa.
- Criterios de aceptación numerados, binarios y verificables; un número eliminado no se reusa.
- **Estado** se actualiza a mano. Un CU puede estar Hecho en web y Pendiente en mobile: el campo
  Plataformas y las carpetas `UT-*` lo aclaran.

## Qué NO va en un CU

Nombres de clases, endpoints exactos, esquemas de tablas o listas de tareas técnicas: eso va en las
[Unidades de Trabajo](UNIDADES_DE_TRABAJO.md). Los **códigos de error** sí, porque el actor los ve
traducidos a un mensaje.
