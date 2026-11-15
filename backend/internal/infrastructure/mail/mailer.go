// Package mail implementa el envío de correo transaccional por SMTP (Mailpit en desarrollo).
package mail

import (
	"fmt"
	"html"
	"net/smtp"
)

// MailerSMTP implementa ports.Mailer.
type MailerSMTP struct {
	Host      string
	Puerto    int
	Usuario   string
	Clave     string
	Remitente string
}

func NewMailerSMTP(host string, puerto int, usuario, clave, remitente string) *MailerSMTP {
	return &MailerSMTP{Host: host, Puerto: puerto, Usuario: usuario, Clave: clave, Remitente: remitente}
}

func (m *MailerSMTP) EnviarActivacion(destino, nombres, enlace string) error {
	asunto := "Activa tu cuenta — Escalafón Docente UAGRM"
	cuerpo := plantilla("Activa tu cuenta", nombres,
		"Se creó una cuenta para ti en el Sistema de Gestión de Flujos de Trabajo y Rastreo Documental. "+
			"Activa tu cuenta y define tu contraseña con el siguiente enlace (vigente por tiempo limitado):",
		enlace, "Activar cuenta")
	return m.enviar(destino, asunto, cuerpo)
}

func (m *MailerSMTP) EnviarRecuperacion(destino, nombres, enlace string) error {
	asunto := "Recupera tu contraseña — Escalafón Docente UAGRM"
	cuerpo := plantilla("Recupera tu contraseña", nombres,
		"Recibimos una solicitud para restablecer tu contraseña. Si fuiste tú, usa el siguiente enlace "+
			"(vigente por 30 minutos). Si no fuiste tú, ignora este correo:",
		enlace, "Restablecer contraseña")
	return m.enviar(destino, asunto, cuerpo)
}

var textoEstado = map[string]string{
	"OBSERVADO":   "tiene observaciones que debes atender",
	"EN_FACULTAD": "fue remitido a la facultad (Tribunal, Decano y Consejo Facultativo); en esta etapa no se puede rastrear su avance y te avisaremos cuando regrese",
	"REMITIDO":    "fue remitido al Rectorado",
	"ARCHIVADO":   "fue archivado",
}

func (m *MailerSMTP) EnviarCambioEstado(destino, nombres, codigo, estado, enlace string) error {
	asunto := "Novedades de tu trámite " + codigo + " — Escalafón Docente UAGRM"
	cuerpo := plantilla("Tu trámite cambió de estado", nombres,
		fmt.Sprintf("Tu trámite %s %s. Consulta el detalle en el portal con tu código de seguimiento:", html.EscapeString(codigo), textoEstado[estado]),
		enlace, "Ver seguimiento")
	return m.enviar(destino, asunto, cuerpo)
}

func (m *MailerSMTP) enviar(destino, asunto, cuerpoHTML string) error {
	addr := fmt.Sprintf("%s:%d", m.Host, m.Puerto)
	var auth smtp.Auth
	if m.Usuario != "" {
		auth = smtp.PlainAuth("", m.Usuario, m.Clave, m.Host)
	}
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		m.Remitente, destino, asunto, cuerpoHTML,
	)
	return smtp.SendMail(addr, auth, m.Remitente, []string{destino}, []byte(msg))
}

func plantilla(titulo, nombres, mensaje, enlace, textoBoton string) string {
	return fmt.Sprintf(`<!doctype html>
<html lang="es"><body style="font-family:Arial,sans-serif;background:#F5F6F8;padding:32px;margin:0;">
<div style="max-width:480px;margin:0 auto;background:#fff;border-radius:12px;padding:32px;">
<h1 style="color:#0F1B2E;font-size:20px;margin:0 0 16px;">%s</h1>
<p style="color:#0F1B2E;margin:0 0 12px;">Hola %s,</p>
<p style="color:#5C6478;line-height:1.5;margin:0 0 24px;">%s</p>
<a href="%s" style="display:inline-block;background:#C8102E;color:#fff;text-decoration:none;padding:12px 20px;border-radius:8px;font-weight:bold;">%s</a>
<p style="color:#9AA1B0;font-size:12px;margin-top:24px;">Departamento de Escalafón Docente — UAGRM</p>
</div></body></html>`, titulo, html.EscapeString(nombres), mensaje, html.EscapeString(enlace), textoBoton)
}
