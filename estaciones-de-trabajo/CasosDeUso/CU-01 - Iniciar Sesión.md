# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-01 | Iniciar Sesión | P1 Seguridad | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Usuario interno (Director, Analista, Personal de Archivo, Personal de Ventanilla)
**Actores secundarios:** —
**Objetivo:** Entrar al sistema con su correo institucional y contraseña y ver solo lo que su rol permite.

**Descripción:**
El usuario escribe su correo institucional y su contraseña. El sistema valida las credenciales, emite un token de sesión y devuelve sus datos y la lista de permisos de su rol; con ella arma el menú. Si la contraseña es incorrecta suma un intento fallido y al tercero bloquea la cuenta durante unos minutos. Una cuenta dada de baja o sin activar no puede entrar. El usuario puede cerrar la sesión, que queda revocada aunque el token no haya vencido.

## Precondiciones

- La cuenta fue creada por el Director ([CU-04](CU-04%20-%20Administrar%20Cuentas%20de%20Usuario.md)) y activada con el enlace de activación ([CU-02](CU-02%20-%20Recuperar%20Contrase%C3%B1a.md)).
- El usuario no tiene la cuenta bloqueada.

## Flujo principal

1. El usuario abre la pantalla de inicio de sesión.
2. Escribe su correo institucional y su contraseña y confirma.
3. El sistema verifica que la cuenta exista, esté activa y no esté bloqueada.
4. El sistema compara la contraseña con su hash y reinicia el contador de intentos fallidos.
5. El sistema emite un token de sesión con vencimiento y registra el inicio de sesión en auditoría con la IP.
6. El sistema devuelve los datos del usuario y sus permisos; la interfaz muestra el menú correspondiente.

## Flujos alternativos y excepciones

- **A1 (paso 3-4):** correo inexistente o contraseña incorrecta → 401 `CREDENCIALES_INVALIDAS` con un mensaje genérico (no revela cuál falló) y se audita el intento fallido.
- **A2 (paso 4):** tercer intento fallido consecutivo → 423 `CUENTA_BLOQUEADA`; la cuenta queda bloqueada `MINUTOS_BLOQUEO` minutos.
- **A3 (paso 3):** cuenta inactiva o sin activar → 401 `CREDENCIALES_INVALIDAS`.
- **A4 (paso 2):** correo con formato inválido o campos vacíos → 422 `VALIDACION`.
- **A5:** más de 10 peticiones por minuto desde la misma IP → 429 (limitador de tasa).
- **A6 (cerrar sesión):** el usuario cierra sesión → el token queda revocado y cualquier uso posterior responde 401.

## Reglas de negocio

- Desvío de `cu.md`: se inicia sesión solo por correo institucional, no por código de funcionario.
- El bloqueo se dispara al llegar a `MAX_INTENTOS_LOGIN` (3) intentos fallidos consecutivos.
- La sesión se valida en cada petición contra Redis (versión de sesión y revocación por `jti`); si Redis no responde, contra Postgres.
- El menú se arma por **permisos**, no por nombre de rol.
- Nunca se guardan contraseñas ni tokens en auditoría.

## Postcondiciones

- El usuario tiene una sesión vigente y ve el menú de su rol.
- El inicio de sesión (exitoso o fallido) queda en auditoría.

## Criterios de aceptación

- **CA01:** Dado que el usuario escribe credenciales válidas, el sistema debe emitir la sesión y mostrar solo los módulos de sus permisos.
- **CA02:** Dado que la contraseña es incorrecta, el sistema debe responder `CREDENCIALES_INVALIDAS` sin indicar si el correo existe.
- **CA03:** Dado que es el tercer intento fallido seguido, el sistema debe bloquear la cuenta y responder `CUENTA_BLOQUEADA`.
- **CA04:** Dado que la cuenta está dada de baja, el sistema debe rechazar el inicio de sesión.
- **CA05:** Dado que el usuario cerró sesión, el sistema debe rechazar con 401 cualquier petición con ese token.
- **CA06:** Dado que hay sesión vigente, abrir la pantalla de login debe llevar al panel.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-01%20-%20Iniciar%20Sesi%C3%B3n/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-01%20-%20Iniciar%20Sesi%C3%B3n/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-01%20-%20Iniciar%20Sesi%C3%B3n/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
