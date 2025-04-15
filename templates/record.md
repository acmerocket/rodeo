{{ .command }} {{ .collection }}
{{ range $key, $value := .record }}
* **{{ $key }}**: {{ $value }}
{{ end }}
---
