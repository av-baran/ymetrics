package templates

const (
	GetAllPageTemplate = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8">
		<title> ymetrics </title>
        <style>
            body {
                background-color: lightgrey;
            }
            h1 {
                font-family: monospace;
                font-size: 160%;
                font-weight: bold;
            }
            ul {
                font-family: monospace;
                font-size: 120%;
            }
        </style>
    </head>
    <body>
        <h1>Metric list</h1>
        <ul>
            {{ range $key, $value := . }}
               <li><strong>{{ $key }}</strong>: {{ $value }}</li>
            {{ end }}
        </ul>
    </body>
</html>`
)
