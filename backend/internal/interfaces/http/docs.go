package http

import (
	"github.com/gofiber/fiber/v2"

	"escalafon/backend/api"
)

const swaggerHTML = `<!doctype html>
<html><head><title>Escalafón — API docs</title>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui.min.css">
</head><body>
<div id="swagger-ui"></div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/5.17.14/swagger-ui-bundle.min.js"></script>
<script>window.onload = () => SwaggerUIBundle({url: '/api/docs/openapi.yaml', dom_id: '#swagger-ui'});</script>
</body></html>`

// RegistrarDocs monta Swagger UI en /api/docs. Solo se llama en desarrollo (sección 7).
func RegistrarDocs(app *fiber.App) {
	app.Get("/api/docs", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTMLCharsetUTF8)
		return c.SendString(swaggerHTML)
	})
	app.Get("/api/docs/openapi.yaml", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, "application/yaml")
		return c.Send(api.OpenAPI)
	})
}
