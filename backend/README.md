# Escalafón UAGRM — Backend

API REST en Go para el Sistema de Gestión de Flujos de Trabajo y Rastreo Documental del Escalafón Docente (UAGRM). Este backend cubre **HU-01 (Acceso y Gestión de Credenciales)** y **HU-02 (Administración de Usuarios, Roles y Permisos)**, según el plan de implementación en [`../archivos/docsmd/Inicializacion.md`](../archivos/docsmd/Inicializacion.md).

## Stack

| Capa                    | Tecnología                                                  |
| ----------------------- | ------------------------------------------------------------ |
| API                     | Go + [Fiber v2](https://gofiber.io/)                         |
| ORM                     | GORM (`gorm.io/driver/postgres`, pgx)                         |
| Migraciones             | [Atlas](https://atlasgo.io/) (`ariga.io/atlas-provider-gorm`) |
| Base de datos           | PostgreSQL (Neon en desarrollo/producción)                    |
| Caché de sesión         | Redis                                                          |
| Autenticación           | JWT (HS256) + bcrypt (costo 12)                                |
| Correo (desarrollo)     | Mailpit (SMTP local)                                           |
| Documentación de la API | OpenAPI 3 servido con Swagger UI en `/api/docs`                |

## Arquitectura

El backend sigue **Clean Architecture** con una regla de dependencia verificada automáticamente (`make arch`, incluido en `make test`): las capas externas dependen de las internas, nunca al revés.

```
interfaces/http  ──►  application  ──►  domain
infrastructure   ──►  application  ──►  domain
```

- **`internal/domain`** — entidades (`usuario`, `rol`, `token`, `auditoria`), reglas de negocio y puertos (interfaces) de repositorio. No importa GORM, Fiber, Redis, JWT ni ninguna otra capa.
- **`internal/application`** — casos de uso (`auth`, `perfil`, `usuarios`, `roles`) y los puertos que necesitan (`ports`: `TokenService`, `PasswordHasher`, `Mailer`, `SessionCache`, `TxManager`, `Clock`). Solo depende de `domain`.
- **`internal/infrastructure`** — implementaciones concretas: repositorios GORM, caché de sesión en Redis, JWT/bcrypt, generador de tokens, mailer SMTP y configuración por variables de entorno.
- **`internal/interfaces/http`** — router, middlewares (`Autenticar`, `RequierePermiso`, manejo de errores, request-id) y handlers Fiber. Traducen HTTP a casos de uso; no contienen lógica de negocio.
- **`cmd/api/main.go`** — único punto donde se instancian las implementaciones concretas y se inyectan en los casos de uso (composition root).

Otros puntos de entrada:

- **`cmd/seed`** — carga los datos iniciales (idempotente).
- **`cmd/atlasloader`** — expone los modelos GORM a Atlas para generar las migraciones.

### Sesión sin consultar la BD en cada petición

El middleware `Autenticar` valida el JWT localmente (firma, `exp`, `iss`) y luego hace una sola ida y vuelta a Redis para comprobar `version_sesion` y revocación (`jti`). Si Redis no está disponible, la validación cae a Postgres (la app sigue funcionando, más lenta). Ver la sección 5 de `Inicializacion.md` para el detalle completo.

## Usuarios creados por el seed

El comando `make seed` (`cmd/seed`) crea los datos iniciales del módulo de seguridad de forma idempotente:

**Roles de sistema** (`es_sistema = true`, no se pueden eliminar ni renombrar):

| Rol                  | Permisos                                                                                                        |
| --------------------- | ----------------------------------------------------------------------------------------------------------------- |
| `Director`            | Todos los permisos del catálogo (no se le pueden retirar `USUARIO_GESTIONAR` ni `ROL_GESTIONAR`)                  |
| `Analista`            | `PERFIL_GESTIONAR`, `DOCENTE_GESTIONAR`, `TRAMITE_RADICAR`, `TRAMITE_REVISAR`, `DOCUMENTO_GESTIONAR`, `DOCUMENTO_BUSCAR`, `CARPETA_SOLICITAR` |
| `Personal de Archivo` | `PERFIL_GESTIONAR`, `LEGAJO_GESTIONAR`, `CUSTODIA_REGISTRAR`, `CONSTANCIA_EMITIR`, `REPORTE_GENERAR`               |

**Usuarios de prueba** (uno por rol, ya activos, solo para desarrollo local — no cumplen la política de contraseña y no deben usarse fuera de este entorno):

| Correo             | Rol                  | Contraseña |
| ------------------- | --------------------- | ---------- |
| `admin@ded.dev`      | Director              | `123123`   |
| `analista@ded.dev`   | Analista              | `123123`   |
| `archivo@ded.dev`    | Personal de Archivo   | `123123`   |
| `ventanilla@ded.dev` | Personal de Ventanilla | `123123`  |

**Primer Director (flujo real):** además, el seed crea el usuario de producción/staging a partir de `SEED_DIRECTOR_CORREO` / `SEED_DIRECTOR_NOMBRES` / `SEED_DIRECTOR_APELLIDOS` en `.env`. No tiene contraseña por defecto: se crea con una contraseña aleatoria inutilizable y un token de **activación** de un solo uso (vigente `HORAS_TOKEN_ACTIVACION`, 24 h por defecto). El enlace de activación se imprime en la consola y se envía por correo (Mailpit en desarrollo).

Cualquier otro usuario se crea desde la API (`POST /api/usuarios`, requiere `USUARIO_GESTIONAR`) y sigue el mismo flujo de activación.

## Cómo levantar el proyecto

### Requisitos

- Go 1.27+
- Docker / Docker Compose (Redis, Mailpit, Postgres de pruebas)
- [Atlas CLI](https://atlasgo.io/getting-started#installation) (para migraciones)
- Una base de datos Postgres accesible para desarrollo (Neon u otra)

### Pasos

1. **Variables de entorno**

   ```bash
   cp .env.example .env
   ```

   Completa `DATABASE_URL` / `DATABASE_URL_DIRECT` (tu Postgres de desarrollo), `JWT_SECRETO` (mínimo 32 bytes aleatorios) y `SEED_DIRECTOR_CORREO`/`NOMBRES`/`APELLIDOS`. `.env` está en `.gitignore`: nunca se versiona.

2. **Servicios de apoyo** (Redis, Mailpit, Postgres de pruebas)

   ```bash
   docker compose up -d
   ```

   (el `docker-compose.yml` está en la raíz del repositorio, un nivel arriba de `backend/`)

3. **Migraciones**

   ```bash
   make migrate-apply
   make migrate-status   # verifica que todo esté aplicado
   ```

4. **Datos iniciales**

   ```bash
   make seed
   ```

   Copia el enlace de activación que se imprime en consola (o revisa Mailpit en http://localhost:8025) para activar al primer Director.

5. **Levantar la API**

   ```bash
   make dev
   ```

   La API queda en `http://localhost:8080/api`. En `APP_ENV=development` la documentación Swagger está en `http://localhost:8080/api/docs`.

### Pruebas

```bash
make test   # go test ./... + verificación de arquitectura (make arch)
make lint   # go vet
```

Las pruebas automáticas usan el Postgres local de `docker-compose.yml` (`DATABASE_URL_TEST`), nunca la base de datos de Neon.

## Estado actual

Implementado (Fases 0 a 5 del plan): base del repositorio, base de datos y migraciones, dominio y casos de uso de HU-01/HU-02, adaptadores de infraestructura y API HTTP completa de ambas historias, con pruebas unitarias e integración.

Pendiente: frontend (Next.js) y verificación final end-to-end (Fases 6 a 8), fuera del alcance de este README.
