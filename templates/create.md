{{.action}} {{.collection}} {{.time}}
{{ range $key, $value := . }}
* **{{ $key }}**: {{ $value }}
{{ end }}
