-- Foto de perfil del usuario. Tabla aparte (no una columna en "usuario"): si el código nuevo corre antes
-- que esta migración, solo fallan los endpoints de la foto; el inicio de sesión y el resto no la tocan.
CREATE TABLE "usuario_foto" (
  "usuario_id" uuid NOT NULL,
  "clave" character varying(200) NOT NULL,
  "tipo_mime" character varying(20) NOT NULL,
  "actualizado_en" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("usuario_id"),
  CONSTRAINT "fk_usuario_foto_usuario" FOREIGN KEY ("usuario_id") REFERENCES "usuario" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
