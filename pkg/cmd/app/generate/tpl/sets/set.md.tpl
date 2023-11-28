# {{ .Name }}

## Route

:::mermaid
flowchart LR;
{{- range $i, $e := .Environments -}}
{{if $i}} -.-> {{ end }}{{ $e }}
{{- end }}
:::
