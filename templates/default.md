## {{.action}} {{.collection}} {{.time}}
{{ range $key, $value := .record }}
* **{{ $key }}**: {{ $value }}
{{ end }}
---
