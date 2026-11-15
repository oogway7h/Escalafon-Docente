-- Correo único sin distinguir mayúsculas
CREATE UNIQUE INDEX ux_usuario_correo ON usuario (lower(correo));

-- Búsqueda rápida del token pendiente de un usuario
CREATE INDEX idx_token_recuperacion_pendiente ON token_recuperacion (usuario_id) WHERE usado_en IS NULL;

-- Auditoría de solo anexión
CREATE OR REPLACE FUNCTION fn_solo_anexion() RETURNS trigger AS $$
BEGIN
  RAISE EXCEPTION 'La tabla % es de solo anexión: no se permite %', TG_TABLE_NAME, TG_OP
    USING ERRCODE = 'insufficient_privilege';
END; $$ LANGUAGE plpgsql;

CREATE TRIGGER tg_auditoria_inmutable BEFORE UPDATE OR DELETE ON auditoria
  FOR EACH ROW EXECUTE FUNCTION fn_solo_anexion();
CREATE TRIGGER tg_auditoria_truncate BEFORE TRUNCATE ON auditoria
  FOR EACH STATEMENT EXECUTE FUNCTION fn_solo_anexion();

-- actualizado_en automático
CREATE OR REPLACE FUNCTION fn_actualizado_en() RETURNS trigger AS $$
BEGIN NEW.actualizado_en := now(); RETURN NEW; END; $$ LANGUAGE plpgsql;
CREATE TRIGGER tg_usuario_actualizado BEFORE UPDATE ON usuario
  FOR EACH ROW EXECUTE FUNCTION fn_actualizado_en();
