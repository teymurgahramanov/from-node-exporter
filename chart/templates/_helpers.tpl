{{- define "resource.name" -}}
{{- $top := index . 0 -}}
{{- $arg := index . 1 -}}
{{- if $arg -}}
{{- /*
Used in range loops.
*/}}
{{- $top.Values.app.name }}-{{ $arg }}-{{ $top.Values.deploy.env | default "local" -}}
{{- else -}}
{{- $top.Values.app.name }}-{{ $top.Values.deploy.env | default "local" -}}
{{- end -}}
{{- end -}}

{{- define "resource.namespace" -}}
{{- $Chart := .Chart -}}
{{ .Values.app.namespace | default "default" }}
{{- end }}

{{- define "common.labels" -}}
{{- $Chart := .Chart -}}
app: {{ .Values.app.name }}
{{- if .Values.labels }}
{{- range $k, $v := .Values.labels }}
  {{ $k }}: {{ $v }}
{{- end }}
{{- end }}
{{- end }}

{{- define "common.annotations" -}}
{{- $Chart := .Chart -}}
{{- if .Values.annotations }}
{{- range $k, $v := .Values.annotations }}
  {{ $k }}: {{ $v }}
{{- end }}
{{- end }}
{{- end }}

{{- define "common.metadata" -}}
{{- $top := index . 0 -}}
{{- $arg := index . 1 -}}
name: {{ include "resource.name" ( list $top $arg ) }}
namespace: {{ include "resource.namespace" $top }}
labels:
  {{ include "common.labels" $top }}
annotations:
  {{ include "common.annotations" $top }}
{{- end }}