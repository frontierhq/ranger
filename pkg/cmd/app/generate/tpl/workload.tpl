# Workloads

## Workloads available in this project

| Name | Instances | Latest version | Latest build |
|:-----------|:-----------|:-----------|:-----------|
{{- range .Workloads }}
| {{ .Name }} | {{ .Instances }} | {{ .Version }} | {{ .Build }} |
{{- end }}
