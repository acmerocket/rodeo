## {{.type}}
{{ range $key, $value := . }}
* **{{ $key }}**: {{ $value }}
{{ end }}
---
