## {{.type}} {{.time}}
{{ range $key, $value := . }}
* **{{ $key }}**: {{ $value }}
{{ end }}
---
