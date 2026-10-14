{{- define "sapprovider-operator.labels" -}}
app.kubernetes.io/name: sapprovider-operator
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: Helm
app.kubernetes.io/part-of: sapprovider-operator
app.kubernetes.io/version: {{ .Chart.AppVersion | replace "+" "_" | trunc 63 | trimSuffix "-" | quote }}
{{- end }}

{{/* Build exceptions only from operator-owned configuration. */}}
{{- define "sapprovider-operator.secretReferenceAllowlist" -}}
{{- $entries := .Values.secretReferenceAllowlist | default list -}}
{{- $namespace := .Values.scpc.ironcore.projectNamespace -}}
{{- $spec := .Values.scpc.ironcore.shootSpec | default dict -}}
{{- if $spec.credentialsBindingName -}}
{{- $entries = append $entries (dict "namespace" $namespace "name" $spec.credentialsBindingName "purpose" "gardener-credential-binding") -}}
{{- end -}}
{{- if $spec.secretBindingName -}}
{{- $entries = append $entries (dict "namespace" $namespace "name" $spec.secretBindingName "purpose" "gardener-secret-binding") -}}
{{- end -}}
{{- $dns := $spec.dns | default dict -}}
{{- range $provider := $dns.providers | default list -}}
{{- if $provider.secretName -}}
{{- $entries = append $entries (dict "namespace" $namespace "name" $provider.secretName "purpose" "gardener-dns") -}}
{{- end -}}
{{- $ref := $provider.credentialsRef | default dict -}}
{{- if $ref.name -}}
{{- $entries = append $entries (dict "namespace" ($ref.namespace | default $namespace) "name" $ref.name "purpose" "gardener-dns") -}}
{{- end -}}
{{- end -}}
{{- $entries | toJson -}}
{{- end -}}
