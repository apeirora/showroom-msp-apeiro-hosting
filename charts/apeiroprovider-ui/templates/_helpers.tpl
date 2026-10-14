{{- define "apeiroprovider-ui.labels" -}}
app.kubernetes.io/name: apeiroprovider-ui
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: Helm
app.kubernetes.io/part-of: apeiroprovider
app.kubernetes.io/version: {{ .Chart.AppVersion | replace "+" "_" | trunc 63 | trimSuffix "-" | quote }}
{{- end }}

{{- define "apeiroprovider-ui.image" -}}
{{- $tag := default .Chart.AppVersion .Values.ui.image.tag -}}
{{- printf "%s:%s" .Values.ui.image.repository $tag -}}
{{- end -}}

{{- define "apeiroprovider-ui.externalUrl" -}}
{{- printf "%s://%s" .Values.PUBLIC_SCHEME .Values.PUBLIC_HOST -}}
{{- end -}}
