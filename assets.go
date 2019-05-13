package echoswagger

// CDN refer to https://www.jsdelivr.com/package/npm/swagger-ui-dist
const DefaultCDN = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@3.19.0"

const SwaggerUIContent = `{{define "swagger"}}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8">
    <title>{{.title}}</title>
    <link rel="stylesheet" href="{{.cdn}}/swagger-ui.css" crossorigin="anonymous" />
    <link rel="icon" type="image/png" href="{{.cdn}}/favicon-32x32.png" sizes="32x32" />
    <link rel="icon" type="image/png" href="{{.cdn}}/favicon-16x16.png" sizes="16x16" />
    <style>
      html
      {
        box-sizing: border-box;
        overflow: -moz-scrollbars-vertical;
        overflow-y: scroll;
      }

      *,
      *:before,
      *:after
      {
        box-sizing: inherit;
      }

      body
      {
        margin:0;
        background: #fafafa;
      }

      {{if .hideTop}}#swagger-ui>.swagger-container>.topbar
      {
        display: none;
      }{{end}}
    </style>
  </head>

  <body>
    <div id="swagger-ui"></div>

    <script src="{{.cdn}}/swagger-ui-bundle.js" crossorigin="anonymous"></script>
    <script src="{{.cdn}}/swagger-ui-standalone-preset.js" crossorigin="anonymous"></script>
    <script>
    window.onload = function() {
      var specPath = "{{.specName}}"
      if (!window.location.pathname.endsWith("/")) {
        specPath = "/" + specPath
      }
      var specStr = "{{.spec}}"
      var spec = specStr ? JSON.parse(specStr) : undefined
      if (spec) {
        spec.host = window.location.host
        var docPath = "{{.docPath}}"
        var basePath = window.location.pathname
        if (!docPath.endsWith("/")) { docPath += "/" }
        if (!basePath.endsWith("/")) { basePath += "/" }
        if (basePath.endsWith(docPath)) {
          basePath = basePath.slice(0, -docPath.length)
        }
        spec.basePath = basePath
      }
      // Build a system
      const ui = SwaggerUIBundle({
        url: window.location.origin+window.location.pathname+specPath,
        spec: spec,
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      })

      window.ui = ui
    }
  </script>
  </body>
</html>
{{end}}`
