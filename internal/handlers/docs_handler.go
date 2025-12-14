package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func Docs(c *gin.Context) {
    html := `<!doctype html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Devices API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: 'BaseLayout'
      });
    </script>
  </body>
</html>`
    c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

