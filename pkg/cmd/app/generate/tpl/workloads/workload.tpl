# Workload
#
#
#
{{- range .Workloads }}
| {{ .Name }} | {{ .Instances }} | {{ .Version }} | {{ .Build }} |
{{- end }}
