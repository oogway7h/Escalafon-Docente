# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-02 | Recuperar Contraseña | P1 Seguridad | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Usuario interno
**Actores secundarios:** Servidor de correo
**Objetivo:** Restablecer la contraseña olvidada con un enlace de un solo uso enviado al correo institucional.

**Descripción:**
El usuario pide recuperar su contraseña escribiendo su correo. El sistema responde siempre lo mismo (exista o no la cuenta) y, si existe, envía un enlace de un solo uso que vence a los treinta minutos. Al abrirlo, el usuario escribe una contraseña nueva que cumpla la política; el sistema la guarda, invalida el enlace y cierra todas las sesiones previas. El mismo mecanismo activa una cuenta nueva (primer acceso).

## Precondiciones

- El usuario tiene una cuenta registrada ([CU-04](CU-04%20-%20Administrar%20Cuentas%20de%20Usuario.md)).

## Flujo principal

1. El usuario abre «¿Olvidaste tu contraseña?» y escribe su correo.
2. El sistema genera un token aleatorio, guarda solo su hash con vencimiento de 30 minutos y envía el enlace por correo.
3. El sistema muestra un mensaje neutro de envío.
4. El usuario abre el enlace y escribe la contraseña nueva dos veces.
5. El sistema valida token y política, guarda el nuevo hash, marca el token como usado y sube la versión de sesión del usuario.
6. El sistema confirma y lleva al login.

## Flujos alternativos y excepciones

- **A1 (paso 2):** el correo no existe → mismo mensaje neutro, no se envía nada.
- **A2 (paso 5):** token vencido, usado o inexistente → 400 `TOKEN_INVALIDO`.
- **A3 (paso 5):** contraseña que no cumple la política → 422 `POLITICA_CONTRASENA`.
- **A4 (activación):** cuenta ya activada → 409 `CUENTA_YA_ACTIVADA`.

## Reglas de negocio

- El enlace es de un solo uso y vence a los `MINUTOS_TOKEN_RECUPERACION` (30) minutos; el de activación, a las `HORAS_TOKEN_ACTIVACION` horas.
- Solo se guarda el hash del token, nunca el token.
- Restablecer invalida todas las sesiones previas del usuario.
- Política: mínimo 8 caracteres con letras y números.

## Postcondiciones

- La contraseña cambió, el enlace ya no sirve y las sesiones anteriores están cerradas.
- El cambio queda en auditoría (sin la contraseña ni el token).

## Criterios de aceptación

- **CA01:** Dado que el correo existe, el sistema debe enviar un enlace de un solo uso.
- **CA02:** Dado que el correo no existe, el sistema debe responder el mismo mensaje sin enviar correo.
- **CA03:** Dado que el enlace tiene más de 30 minutos, el sistema debe rechazarlo con `TOKEN_INVALIDO`.
- **CA04:** Dado que el enlace ya se usó, el sistema debe rechazarlo.
- **CA05:** Dado que la contraseña se restableció, el sistema debe rechazar los tokens de sesión emitidos antes.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-02%20-%20Recuperar%20Contrase%C3%B1a/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-02%20-%20Recuperar%20Contrase%C3%B1a/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
