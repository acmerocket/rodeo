# Summary
{{ range $key, $value := . }}
* **{{ $key }}**: {{ $value }}
{{ end }}
---
