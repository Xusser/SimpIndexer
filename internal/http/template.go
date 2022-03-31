package http

import (
	"html/template"
)

type Item struct {
	Name    string
	Href    string
	ModDate string
	Size    string
}

var html = template.Must(template.New("page").Parse(`<html>
<head><title>Index of {{ .path }}</title></head>
<body>
<h1>Index of {{ .path }}</h1>
<hr>
<ul>
{{ range $item := .items }}
<li><a href="{{ .Href }}">{{ .Name }}</a></li>
{{ end }}
</ul>
<hr>
</body>
</html>`))

// var html = template.Must(template.New("page").Parse(`
// <!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 3.2 Final//EN">
// <html>
// <head>
// <title>Index of {{ .path }}</title>
// </head>
// <body>
// <h1>Index of {{ .path }}</h1>
// <ul>
// {{ range $item := .items }}
// <li><a href="{{ .Href }}"> {{ .Name }}</a></li>
// {{ end }}
// </ul>
// </body>
// </html>
// `))
