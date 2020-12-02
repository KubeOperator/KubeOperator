{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "helm-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "helm-operator.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "helm-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name of the service account to use.
*/}}
{{- define "helm-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "helm-operator.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Create the name of the cluster role to use.
*/}}
{{- define "helm-operator.clusterRoleName" -}}
{{- if .Values.clusterRole.create -}}
    {{ default (include "helm-operator.fullname" .) .Values.clusterRole.name }}
{{- else -}}
    {{ default "default" .Values.clusterRole.name }}
{{- end -}}
{{- end -}}

{{/*
Create a custom repositories.yaml for Helm.
*/}}
{{- define "helm-operator.customRepositories" -}}
apiVersion: v1
generated: 0001-01-01T00:00:00Z
repositories:
{{- range .Values.configureRepositories.repositories }}
- name: {{ required "Please specify a name for the Helm repo" .name }}
  url: {{ required "Please specify a URL for the Helm repo" .url }}
  cache: /var/fluxd/helm/repository/cache/{{ .name }}-index.yaml
  caFile: "{{ .caFile | default "" }}"
  certFile: "{{ .certFile | default "" }}"
  keyFile: "{{ .keyFile | default "" }}"
  password: "{{ .password | default "" }}"
  username: "{{ .username | default "" }}"
{{- end }}
{{- end -}}

{{/*
Create Helm plugin init containers.
*/}}
{{- define "helm-operator.initPlugins" -}}
{{- range $i, $v := .Values.initPlugins.plugins -}}
{{- $name := printf "helm-plugin-init-%02d" $i -}}
{{- $plugin := required "Please specify the plugin" $v.plugin -}}
{{- $helmVersion := required "Please specify the targeted Helm version" $v.helmVersion -}}
{{- if contains $helmVersion $.Values.helm.versions }}
- name: {{ $name }}
  image: "{{ $.Values.image.repository }}:{{ $.Values.image.tag }}"
  imagePullPolicy: {{ $.Values.image.pullPolicy }}
{{- if eq $helmVersion "v2" }}
  command: ['sh', '-c', 'helm2 plugin install {{ $plugin }}{{ if $v.version }} --version {{ $v.version }}{{ end }}']
  volumeMounts:
  - name: {{ $.Values.initPlugins.cacheVolumeName | quote }}
    mountPath: /var/fluxd/helm/cache/plugins/
    subPath: {{ $helmVersion }}
  - name: {{ $.Values.initPlugins.cacheVolumeName | quote }}
    mountPath: /var/fluxd/helm/plugins
    subPath: {{ $helmVersion }}-config
{{- end }}
{{- if eq $helmVersion "v3" }}
  command: ['sh', '-c', 'helm3 plugin install {{ $plugin }}{{ if $v.version }} --version {{ $v.version }}{{ end }}']
  volumeMounts:
  - name: {{ $.Values.initPlugins.cacheVolumeName | quote }}
    mountPath: /root/.cache/helm/plugins
    subPath: {{ $helmVersion }}
  - name: {{ $.Values.initPlugins.cacheVolumeName | quote }}
    mountPath: /root/.local/share/helm/plugins
    subPath: {{ $helmVersion }}-config
{{- end }}
{{- end }}
{{- end }}
{{- end -}}

{{/*
Create the name of the Git config Secret.
*/}}
{{- define "git.config.secretName" -}}
{{- if .Values.git.config.enabled }}
    {{- if .Values.git.config.secretName -}}
        {{ default "default" .Values.git.config.secretName }}
    {{- else -}}
        {{ default (printf "%s-git-config" (include "helm-operator.fullname" .)) }}
{{- end -}}
{{- end }}
{{- end }}
