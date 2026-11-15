# Estaciones de trabajo — SDD del Sistema de Escalafón Docente

Especificación del sistema DED (UAGRM) con **Spec-Driven Development (SDD)**, siguiendo la misma
metodología de *Casos de Uso → Unidades de Trabajo* del proyecto DriveHub. Cubre los tres clientes:

| Proyecto | Stack | Cómo se documentó |
|---|---|---|
| `../backend` | Go · Fiber · GORM/Postgres · Redis | **Retroactivo**: los 26 CU ya estaban hechos; las UT describen y citan el código existente |
| `../frontend` | Next.js 16 · React 19 | **Retroactivo**, igual que el backend |
| `../ded_Mobile` | Flutter (Android + iOS) | **Spec-first**: las UT-MOBILE se escribieron antes del código y el código las cita |

## El ciclo SDD

```
 cu.md (cátedra)          CasosDeUso/CU-NN.md             UnidadesTrabajo/CU-NN/UT-*/UT-…md        repo
 ─────────────── ──────▶ ────────────────────── ──────▶ ───────────────────────────────── ──────▶ ─────────
 1 línea por CU   SPEC   actor, flujos, reglas,   PLAN   alcance exacto, archivos, fuera de   CODE   commit que
                         errores y CA verificables       alcance y validación por UT                cita la UT
                                                                                                       │
                                     VALIDACIÓN ◀───────────────────────────────────────────────────┘
                                     (go test / flutter test / pasos manuales de la UT)
```

1. **Specify** — un CU por flujo con criterios de aceptación binarios ([regla](.claude/rules/CASOS_DE_USO.md)).
2. **Plan** — el CU se parte en UT por capa, cada una con 1-4 archivos, alcance, fuera de alcance y
   validación ([regla](.claude/rules/UNIDADES_DE_TRABAJO.md)).
3. **Implement** — se le pide a la IA *una* UT; las reglas de `.claude/rules/` son el contexto fijo.
4. **Validate** — se corre la validación de la UT; si pasa, commit `feat: UT-CUNN-MOB-NN …`.

La IA no decide el alcance: lo decide la UT. Eso es lo que evita que invente endpoints, mezcle capas o
construya trabajo de otro CU.

## Estructura

```
estaciones-de-trabajo/
├── README.md                 este archivo
├── BITACORA_IA.md            evidencia de uso de IA: pedido, resultado, revisión humana, tiempos
├── .claude/rules/            reglas que la IA debe respetar (adaptadas de DriveHub)
│   ├── CASOS_DE_USO.md · UNIDADES_DE_TRABAJO.md · CONVENCIONES_GIT.md
│   ├── ARQUITECTURA_BACKEND.md · ARQUITECTURA_FRONTEND.md · ARQUITECTURA_MOBILE.md
│   └── mobile/               consumo de API, estado, navegación, formularios, UI, testing
├── CasosDeUso/               26 CU
└── UnidadesTrabajo/          una carpeta por CU con UT-BACKEND, UT-FRONTEND y, si aplica, UT-MOBILE
```

## Índice de casos de uso

| CU | Nombre | Web | Mobile |
|---|---|---|---|
| [CU-01](CasosDeUso/CU-01%20-%20Iniciar%20Sesi%C3%B3n.md) | Iniciar Sesión | ✅ | ✅ |
| [CU-02](CasosDeUso/CU-02%20-%20Recuperar%20Contrase%C3%B1a.md) | Recuperar Contraseña | ✅ | — |
| [CU-03](CasosDeUso/CU-03%20-%20Editar%20Perfil%20Personal.md) | Editar Perfil Personal | ✅ | — |
| [CU-04](CasosDeUso/CU-04%20-%20Administrar%20Cuentas%20de%20Usuario.md) | Administrar Cuentas de Usuario | ✅ | — |
| [CU-05](CasosDeUso/CU-05%20-%20Gestionar%20Roles%20y%20Permisos.md) | Gestionar Roles y Permisos | ✅ | — |
| [CU-06](CasosDeUso/CU-06%20-%20Consultar%20Registro%20de%20Auditor%C3%ADa.md) | Consultar Registro de Auditoría | ✅ | — |
| [CU-07](CasosDeUso/CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas.md) | Gestionar Docentes y Categorías | ✅ | — |
| [CU-08](CasosDeUso/CU-08%20-%20Radicar%20Tr%C3%A1mite.md) | Radicar Trámite | ✅ | — |
| [CU-09](CasosDeUso/CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite.md) | Asignar o Reasignar Trámite | ✅ | ✅ |
| [CU-10](CasosDeUso/CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones.md) | Revisar Expediente y Registrar Observaciones | ✅ | ✅ (consulta) |
| [CU-11](CasosDeUso/CU-11%20-%20Emitir%20Dictamen%20y%20Archivar%20Tr%C3%A1mite.md) | Emitir Dictamen y Archivar Trámite | ✅ | — |
| [CU-12](CasosDeUso/CU-12%20-%20Consultar%20Estado%20de%20Tr%C3%A1mite%20%28Portal%20P%C3%BAblico%29.md) | Consultar Estado de Trámite (Portal Público) | ✅ | ✅ |
| [CU-13](CasosDeUso/CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico.md) | Radicar Solicitud desde el Portal Público | ✅ | ✅ |
| [CU-14](CasosDeUso/CU-14%20-%20Notificar%20Cambios%20de%20Estado%20al%20Docente.md) | Notificar Cambios de Estado al Docente | ✅ | — |
| [CU-15](CasosDeUso/CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR.md) | Cargar y Procesar Documento con OCR | ✅ | ✅ |
| [CU-16](CasosDeUso/CU-16%20-%20Dar%20de%20Baja%20Documento.md) | Dar de Baja Documento | ✅ | — |
| [CU-17](CasosDeUso/CU-17%20-%20Restituir%20Documento%20Dado%20de%20Baja.md) | Restituir Documento Dado de Baja | ✅ | — |
| [CU-18](CasosDeUso/CU-18%20-%20Buscar%20Documentaci%C3%B3n%20Digital.md) | Buscar Documentación Digital | ✅ | — |
| [CU-19](CasosDeUso/CU-19%20-%20Gestionar%20Inventario%20de%20Legajos.md) | Gestionar Inventario de Legajos | ✅ | ✅ (consulta) |
| [CU-20](CasosDeUso/CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico.md) | Solicitar Legajo Físico | ✅ | ✅ |
| [CU-21](CasosDeUso/CU-21%20-%20Atender%20Solicitud%20de%20Legajo.md) | Atender Solicitud de Legajo | ✅ | ✅ |
| [CU-22](CasosDeUso/CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29.md) | Despachar Carpeta Física (Check-out) | ✅ | ✅ |
| [CU-23](CasosDeUso/CU-23%20-%20Registrar%20Devoluci%C3%B3n%20de%20Carpeta%20%28Check-in%29.md) | Registrar Devolución de Carpeta (Check-in) | ✅ | ✅ |
| [CU-24](CasosDeUso/CU-24%20-%20Consultar%20Bit%C3%A1cora%20de%20Custodia.md) | Consultar Bitácora de Custodia | ✅ | — |
| [CU-25](CasosDeUso/CU-25%20-%20Visualizar%20Tablero%20de%20Indicadores.md) | Visualizar Tablero de Indicadores | ✅ | ✅ |
| [CU-26](CasosDeUso/CU-26%20-%20Generar%20y%20Exportar%20Reportes.md) | Generar y Exportar Reportes | ✅ | — |

## La app móvil

Una app con dos entradas: **«Soy docente»** (sin cuenta: seguimiento y nueva solicitud con fotos de los
documentos) y **«Soy personal»** (login; el inicio muestra solo los accesos de sus permisos). Detalle
en [ARQUITECTURA_MOBILE.md](.claude/rules/ARQUITECTURA_MOBILE.md).

**iOS sin Mac:** no se usa Hackintosh (la licencia de macOS solo permite hardware Apple). El repo
`DED-Mobile` compila para iOS en GitHub Actions (`macos-latest`, `flutter build ios --no-codesign`) en
cada push; el artefacto es la evidencia. Para instalar en un iPhone hace falta firmar con una cuenta
Apple Developer (Codemagic o el mismo workflow con certificados).

## Cómo pedirle una UT a la IA

```
Lee estaciones-de-trabajo/.claude/rules/ y la UT UT-CU21-MOB-01 con su README.
Implementa solo esa UT. Respeta "Fuera de alcance". Corre su Validación y reporta.
```
