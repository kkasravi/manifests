package tests_test

import (
	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/target"
	"testing"
)

func writeTektoncdTasksOverlaysApplication(th *KustTestHarness) {
	th.writeF("/manifests/tektoncd/tektoncd-tasks/overlays/application/application.yaml", `
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  name: tektoncd
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: 
      app.kubernetes.io/instance: tektoncd
      app.kubernetes.io/managed-by: kfctl
      app.kubernetes.io/component: tektoncd
      app.kubernetes.io/part-of: kubeflow
      app.kubernetes.io/version: v0.6
  componentKinds:
  - group: v1
    kind: ConfigMap
  - group: v1
    kind: Secret
  - group: tekton.dev/v1alpha1
    kind: PipelineResource
  - group: tekton.dev/v1alpha1
    kind: Task
  - group: tekton.dev/v1alpha1
    kind: TaskRun
  descriptor:
    type: tekton
    version: v1beta1
    description: Launches a build-and-push
    maintainers:
    owners:
    keywords:
    links:
    - description: About
      url: 
  addOwnerRef: true
`)
	th.writeK("/manifests/tektoncd/tektoncd-tasks/overlays/application", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
resources:
- application.yaml
commonLabels:
  app.kubernetes.io/name: centraldashboard
  app.kubernetes.io/instance: centraldashboard
  app.kubernetes.io/managed-by: kfctl
  app.kubernetes.io/component: centraldashboard
  app.kubernetes.io/part-of: kubeflow
  app.kubernetes.io/version: v0.6
`)
	th.writeF("/manifests/tektoncd/tektoncd-tasks/base/persistent-volume-claim.yaml", `
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: kubeflow-pvc
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
`)
	th.writeF("/manifests/tektoncd/tektoncd-tasks/base/secret.yaml", `
---
apiVersion: v1
kind: Secret
metadata:
  name: kaniko-secret
type: Opaque
data:
  kaniko-secret.json: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiY29uc3RhbnQtY3ViaXN0LTE3MzEyMyIsCiAgInByaXZhdGVfa2V5X2lkIjogImQxMjkzODQ1NWI2NGEyZDlhZWE1MDVjNjZkNzIyMjJmNmUyNDg0MzYiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ2dJRHptZmFqMkdKSnVcbjlvOUxsa3dDYXIvVXYyalN3MEJTRXpNeWozQlNCc3BSdXVhOG1qU3MyZTdNRjVnYmMxVGdoL21HSkNqMitZS3JcbkZDOHdyMEt5bWx6NDdRWE5aeUMxb2I2TitxUnJncEJyeXNZWE1iSkFnQUFQck1rS2lhY2VjRzBHN3JkVGhjQVNcblh2Y2hBZG9RdTFTTWhNQys4aEtHMytpM1ZETXE2K3o4ZStwL1AzeEhhdUl3UTBYcWp4RGpONE5VRlMrbTF3WG9cbmpQM2xkbE15cWJUZHFvdFNHanIrTUxJZ1U3OEowdlp1NDF5WDl1ZjY2QWJyOE91N2E4bFdpVnpJUzZubHVqSTJcbm00T3NWYXlZaUN6dHhOam5iazVCVUhaYXZpekNaLzE1YzRlWFRuYXlsQXI0S2IrS3ZHSk5XbXVjSUdjQlNkZnlcbjhGZE1aNEdSQWdNQkFBRUNnZ0VBQWhXaFdpQ1Y4d3d4dmMvMHh1UmdQVzgzQlM4OGdveXVjcmw1OWpIOGgvd0lcbkpSN3VHWWJ3bE9nUXg3Ukg3UldPVWx0YXkzYVl2b1h0ZitPOW9JYXJqTnVRclBubUlUQW1QMDhQMzRpV3RnOWdcbjZRV0UvdkROSU10VFlhMm9NdStyVHRMNzhzTG5zQ2ExMWxkaTE0QUNPS3R3Ymk0UWIwaTJ3Q0VJK1Y1a0ZsVElcbmlkZkFiWGdiakt6dDR0dFg5YmR3OVBpUEcrTkRxQkVpQ0FiaXNMYmxyWlVOWUl1WVVsVWpoTmNKS29hRWNXQkdcbkMzcS9BMk1oTElTZ2FVQlYyZnludXEyaVdVZFZQTGwxNTRNNWtMREFSNk8rS0J6UWVZVnUraGxzdVcrSE1qcGFcbjI5U05FSWhuWmYySERPTi82bUw3R1VEcVlOblhISDVCK1AvT3NnSVQ5UUtCZ1FEV3R4SWNXSks3NFdadmhYeGNcbktOR1pLdHlxQWxielhSVFI3M1NGTVpSQjZNUjdEOFloenIwN3dmRXI5UEk5ekl3QjB3V3BzR2VGRVFQUFpTYXlcbnlKaDQxcTkvYlZ2bUtiM2svRndrMmVFSy9JRFBNRzdqTHdqbVBBbG5UYVpVRUZncHVKYUw5VEVvYlBSWXZtbGVcbmpyeXNmNVRpWXBIY3RkVlhkL3IwZzRqbDNRS0JnUUMrNmg4a28yNVkrS3BPRitxaGU3UnA1Z1VZdFhwQjh6RnpcbjB1UnZacE5DQXFMUmw3Vjk5RW9YcnBZWFF5TE5CY3BSUkpDaEpVenZDMkRtWVhhRXhLUzlWbmZYSzRncG5uVzZcbnAzSitPQ0NtOUJOdUMrYStxYXp0WDluK1JrRXpaem0waFdjTFRTejVpUlhrOWkxTkRKZnY4aXBKd1ptMnhPNDBcbkNXdU5pZ2R4UlFLQmdEbllhRkNxckIxaHhDOFhUMEdrM1pMZU1VUzhES0RUMnVBVUd0Z25XMEhHYStpYmYwMXNcblhSN1VTUjBHaUp5TmxzcUhCMmVIMXR2S2tiUTJGQTdtYSsxaUtUV3pTS2JoWi85ZzNaSXdBS2p0RGViRHJad1dcbjk5YlBKZGxtMmdDYnhxUzJ6aGcybmwrOXVyYU4xZVZibndqNTlpcG5VOVNhU0RlZ1kwT3NqQjBoQW9HQkFKUzRcbnN5d0tlRktzMjVaY1FUWXN0TDF1SjNnNUh4VXpDc29NZGxGbDJiOHBhSWJYcE5XS3NSRkR1cjVDV1dEWGF1VG1cbkFialc0dGl3eDNxUVlCQkxVMzMvVnZueWVtN1pkeUxCZ0lwYzFPclo1aXpxN29TR2p5U1hiNjBLTTQ2RWtrcFRcblJaTmpPbTdsWUgzdFhCclNmYVc0dzBLVG8xZmlqeUZRV1UxNFFoWDFBb0dBSnJrcHg5TjE5UExnQ3MrWDhKa1ZcblRRZXppeHpUd0VLekdRQVU1SnJVai9Cd1g5RDBMa2tkSVdXWXBTMzZqdnpqQUM2Nk93SHBOQlF2dXRRRC9CU3FcbjZUd0ovbVhQN1p0U1hsUFd3MGExVkhNNG5oTmcrbDRXR3BHNCtmdTduZUM5bHlJZkZ2am1jZHM1d0RNNXRDOVFcbjcwOUlYSDEzamdXbzlNYzQyTExkdUxFPVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImRvY2tlckBjb25zdGFudC1jdWJpc3QtMTczMTIzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAiY2xpZW50X2lkIjogIjExMzA4NDc2Nzg5NTE3MzQ1MTIwOCIsCiAgImF1dGhfdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwKICAidG9rZW5fdXJpIjogImh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwKICAiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0cyIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZG9ja2VyJTQwY29uc3RhbnQtY3ViaXN0LTE3MzEyMy5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIKfQo=
---
apiVersion: v1
kind: Secret
metadata:
  name: gcr-secret
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: eyJhdXRocyI6eyJnY3IuaW8iOnsidXNlcm5hbWUiOiJfanNvbl9rZXkiLCJwYXNzd29yZCI6IntcbiAgXCJ0eXBlXCI6IFwic2VydmljZV9hY2NvdW50XCIsXG4gIFwicHJvamVjdF9pZFwiOiBcImNvbnN0YW50LWN1YmlzdC0xNzMxMjNcIixcbiAgXCJwcml2YXRlX2tleV9pZFwiOiBcImQxMjkzODQ1NWI2NGEyZDlhZWE1MDVjNjZkNzIyMjJmNmUyNDg0MzZcIixcbiAgXCJwcml2YXRlX2tleVwiOiBcIi0tLS0tQkVHSU4gUFJJVkFURSBLRVktLS0tLVxcbk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ2dJRHptZmFqMkdKSnVcXG45bzlMbGt3Q2FyL1V2MmpTdzBCU0V6TXlqM0JTQnNwUnV1YThtalNzMmU3TUY1Z2JjMVRnaC9tR0pDajIrWUtyXFxuRkM4d3IwS3ltbHo0N1FYTlp5QzFvYjZOK3FScmdwQnJ5c1lYTWJKQWdBQVByTWtLaWFjZWNHMEc3cmRUaGNBU1xcblh2Y2hBZG9RdTFTTWhNQys4aEtHMytpM1ZETXE2K3o4ZStwL1AzeEhhdUl3UTBYcWp4RGpONE5VRlMrbTF3WG9cXG5qUDNsZGxNeXFiVGRxb3RTR2pyK01MSWdVNzhKMHZadTQxeVg5dWY2NkFicjhPdTdhOGxXaVZ6SVM2bmx1akkyXFxubTRPc1ZheVlpQ3p0eE5qbmJrNUJVSFphdml6Q1ovMTVjNGVYVG5heWxBcjRLYitLdkdKTldtdWNJR2NCU2RmeVxcbjhGZE1aNEdSQWdNQkFBRUNnZ0VBQWhXaFdpQ1Y4d3d4dmMvMHh1UmdQVzgzQlM4OGdveXVjcmw1OWpIOGgvd0lcXG5KUjd1R1lid2xPZ1F4N1JIN1JXT1VsdGF5M2FZdm9YdGYrTzlvSWFyak51UXJQbm1JVEFtUDA4UDM0aVd0ZzlnXFxuNlFXRS92RE5JTXRUWWEyb011K3JUdEw3OHNMbnNDYTExbGRpMTRBQ09LdHdiaTRRYjBpMndDRUkrVjVrRmxUSVxcbmlkZkFiWGdiakt6dDR0dFg5YmR3OVBpUEcrTkRxQkVpQ0FiaXNMYmxyWlVOWUl1WVVsVWpoTmNKS29hRWNXQkdcXG5DM3EvQTJNaExJU2dhVUJWMmZ5bnVxMmlXVWRWUExsMTU0TTVrTERBUjZPK0tCelFlWVZ1K2hsc3VXK0hNanBhXFxuMjlTTkVJaG5aZjJIRE9OLzZtTDdHVURxWU5uWEhINUIrUC9Pc2dJVDlRS0JnUURXdHhJY1dKSzc0V1p2aFh4Y1xcbktOR1pLdHlxQWxielhSVFI3M1NGTVpSQjZNUjdEOFloenIwN3dmRXI5UEk5ekl3QjB3V3BzR2VGRVFQUFpTYXlcXG55Smg0MXE5L2JWdm1LYjNrL0Z3azJlRUsvSURQTUc3akx3am1QQWxuVGFaVUVGZ3B1SmFMOVRFb2JQUll2bWxlXFxuanJ5c2Y1VGlZcEhjdGRWWGQvcjBnNGpsM1FLQmdRQys2aDhrbzI1WStLcE9GK3FoZTdScDVnVVl0WHBCOHpGelxcbjB1UnZacE5DQXFMUmw3Vjk5RW9YcnBZWFF5TE5CY3BSUkpDaEpVenZDMkRtWVhhRXhLUzlWbmZYSzRncG5uVzZcXG5wM0orT0NDbTlCTnVDK2ErcWF6dFg5bitSa0V6WnptMGhXY0xUU3o1aVJYazlpMU5ESmZ2OGlwSndabTJ4TzQwXFxuQ1d1TmlnZHhSUUtCZ0RuWWFGQ3FyQjFoeEM4WFQwR2szWkxlTVVTOERLRFQydUFVR3RnblcwSEdhK2liZjAxc1xcblhSN1VTUjBHaUp5TmxzcUhCMmVIMXR2S2tiUTJGQTdtYSsxaUtUV3pTS2JoWi85ZzNaSXdBS2p0RGViRHJad1dcXG45OWJQSmRsbTJnQ2J4cVMyemhnMm5sKzl1cmFOMWVWYm53ajU5aXBuVTlTYVNEZWdZME9zakIwaEFvR0JBSlM0XFxuc3l3S2VGS3MyNVpjUVRZc3RMMXVKM2c1SHhVekNzb01kbEZsMmI4cGFJYlhwTldLc1JGRHVyNUNXV0RYYXVUbVxcbkFialc0dGl3eDNxUVlCQkxVMzMvVnZueWVtN1pkeUxCZ0lwYzFPclo1aXpxN29TR2p5U1hiNjBLTTQ2RWtrcFRcXG5SWk5qT203bFlIM3RYQnJTZmFXNHcwS1RvMWZpanlGUVdVMTRRaFgxQW9HQUpya3B4OU4xOVBMZ0NzK1g4SmtWXFxuVFFleml4elR3RUt6R1FBVTVKclVqL0J3WDlEMExra2RJV1dZcFMzNmp2empBQzY2T3dIcE5CUXZ1dFFEL0JTcVxcbjZUd0ovbVhQN1p0U1hsUFd3MGExVkhNNG5oTmcrbDRXR3BHNCtmdTduZUM5bHlJZkZ2am1jZHM1d0RNNXRDOVFcXG43MDlJWEgxM2pnV285TWM0MkxMZHVMRT1cXG4tLS0tLUVORCBQUklWQVRFIEtFWS0tLS0tXFxuXCIsXG4gIFwiY2xpZW50X2VtYWlsXCI6IFwiZG9ja2VyQGNvbnN0YW50LWN1YmlzdC0xNzMxMjMuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb21cIixcbiAgXCJjbGllbnRfaWRcIjogXCIxMTMwODQ3Njc4OTUxNzM0NTEyMDhcIixcbiAgXCJhdXRoX3VyaVwiOiBcImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoXCIsXG4gIFwidG9rZW5fdXJpXCI6IFwiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW5cIixcbiAgXCJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmxcIjogXCJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHNcIixcbiAgXCJjbGllbnRfeDUwOV9jZXJ0X3VybFwiOiBcImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZG9ja2VyJTQwY29uc3RhbnQtY3ViaXN0LTE3MzEyMy5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbVwiXG59IiwiZW1haWwiOiJhbnlAdmFsaWQuZW1haWwiLCJhdXRoIjoiWDJwemIyNWZhMlY1T25zS0lDQWlkSGx3WlNJNklDSnpaWEoyYVdObFgyRmpZMjkxYm5RaUxBb2dJQ0p3Y205cVpXTjBYMmxrSWpvZ0ltTnZibk4wWVc1MExXTjFZbWx6ZEMweE56TXhNak1pTEFvZ0lDSndjbWwyWVhSbFgydGxlVjlwWkNJNklDSmtNVEk1TXpnME5UVmlOalJoTW1RNVlXVmhOVEExWXpZMlpEY3lNakl5WmpabE1qUTRORE0ySWl3S0lDQWljSEpwZG1GMFpWOXJaWGtpT2lBaUxTMHRMUzFDUlVkSlRpQlFVa2xXUVZSRklFdEZXUzB0TFMwdFhHNU5TVWxGZGxGSlFrRkVRVTVDWjJ0eGFHdHBSemwzTUVKQlVVVkdRVUZUUTBKTFkzZG5aMU5xUVdkRlFVRnZTVUpCVVVOblNVUjZiV1poYWpKSFNrcDFYRzQ1YnpsTWJHdDNRMkZ5TDFWMk1tcFRkekJDVTBWNlRYbHFNMEpUUW5Od1VuVjFZVGh0YWxOek1tVTNUVVkxWjJKak1WUm5hQzl0UjBwRGFqSXJXVXR5WEc1R1F6aDNjakJMZVcxc2VqUTNVVmhPV25sRE1XOWlOazRyY1ZKeVozQkNjbmx6V1ZoTllrcEJaMEZCVUhKTmEwdHBZV05sWTBjd1J6ZHlaRlJvWTBGVFhHNVlkbU5vUVdSdlVYVXhVMDFvVFVNck9HaExSek1yYVROV1JFMXhOaXQ2T0dVcmNDOVFNM2hJWVhWSmQxRXdXSEZxZUVScVRqUk9WVVpUSzIweGQxaHZYRzVxVUROc1pHeE5lWEZpVkdSeGIzUlRSMnB5SzAxTVNXZFZOemhLTUhaYWRUUXhlVmc1ZFdZMk5rRmljamhQZFRkaE9HeFhhVlo2U1ZNMmJteDFha2t5WEc1dE5FOXpWbUY1V1dsRGVuUjRUbXB1WW1zMVFsVklXbUYyYVhwRFdpOHhOV00wWlZoVWJtRjViRUZ5TkV0aUswdDJSMHBPVjIxMVkwbEhZMEpUWkdaNVhHNDRSbVJOV2pSSFVrRm5UVUpCUVVWRFoyZEZRVUZvVjJoWGFVTldPSGQzZUhaakx6QjRkVkpuVUZjNE0wSlRPRGhuYjNsMVkzSnNOVGxxU0Rob0wzZEpYRzVLVWpkMVIxbGlkMnhQWjFGNE4xSklOMUpYVDFWc2RHRjVNMkZaZG05WWRHWXJUemx2U1dGeWFrNTFVWEpRYm0xSlZFRnRVREE0VURNMGFWZDBaemxuWEc0MlVWZEZMM1pFVGtsTmRGUlpZVEp2VFhVcmNsUjBURGM0YzB4dWMwTmhNVEZzWkdreE5FRkRUMHQwZDJKcE5GRmlNR2t5ZDBORlNTdFdOV3RHYkZSSlhHNXBaR1pCWWxoblltcExlblEwZEhSWU9XSmtkemxRYVZCSEswNUVjVUpGYVVOQlltbHpUR0pzY2xwVlRsbEpkVmxWYkZWcWFFNWpTa3R2WVVWalYwSkhYRzVETTNFdlFUSk5hRXhKVTJkaFZVSldNbVo1Ym5WeE1tbFhWV1JXVUV4c01UVTBUVFZyVEVSQlVqWlBLMHRDZWxGbFdWWjFLMmhzYzNWWEswaE5hbkJoWEc0eU9WTk9SVWxvYmxwbU1raEVUMDR2Tm0xTU4wZFZSSEZaVG01WVNFZzFRaXRRTDA5elowbFVPVkZMUW1kUlJGZDBlRWxqVjBwTE56UlhXblpvV0hoalhHNUxUa2RhUzNSNWNVRnNZbnBZVWxSU056TlRSazFhVWtJMlRWSTNSRGhaYUhweU1EZDNaa1Z5T1ZCSk9YcEpkMEl3ZDFkd2MwZGxSa1ZSVUZCYVUyRjVYRzU1U21nME1YRTVMMkpXZG0xTFlqTnJMMFozYXpKbFJVc3ZTVVJRVFVjM2FreDNhbTFRUVd4dVZHRmFWVVZHWjNCMVNtRk1PVlJGYjJKUVVsbDJiV3hsWEc1cWNubHpaalZVYVZsd1NHTjBaRlpZWkM5eU1HYzBhbXd6VVV0Q1oxRkRLelpvT0d0dk1qVlpLMHR3VDBZcmNXaGxOMUp3TldkVldYUlljRUk0ZWtaNlhHNHdkVkoyV25CT1EwRnhURkpzTjFZNU9VVnZXSEp3V1ZoUmVVeE9RbU53VWxKS1EyaEtWWHAyUXpKRWJWbFlZVVY0UzFNNVZtNW1XRXMwWjNCdWJsYzJYRzV3TTBvclQwTkRiVGxDVG5WREsyRXJjV0Y2ZEZnNWJpdFNhMFY2V25wdE1HaFhZMHhVVTNvMWFWSllhemxwTVU1RVNtWjJPR2x3U25kYWJUSjRUelF3WEc1RFYzVk9hV2RrZUZKUlMwSm5SRzVaWVVaRGNYSkNNV2g0UXpoWVZEQkhhek5hVEdWTlZWTTRSRXRFVkRKMVFWVkhkR2R1VnpCSVIyRXJhV0ptTURGelhHNVlVamRWVTFJd1IybEtlVTVzYzNGSVFqSmxTREYwZGt0cllsRXlSa0UzYldFck1XbExWRmQ2VTB0aWFGb3ZPV2N6V2tsM1FVdHFkRVJsWWtSeVduZFhYRzQ1T1dKUVNtUnNiVEpuUTJKNGNWTXllbWhuTW01c0t6bDFjbUZPTVdWV1ltNTNhalU1YVhCdVZUbFRZVk5FWldkWk1FOXpha0l3YUVGdlIwSkJTbE0wWEc1emVYZExaVVpMY3pJMVdtTlJWRmx6ZEV3eGRVb3paelZJZUZWNlEzTnZUV1JzUm13eVlqaHdZVWxpV0hCT1YwdHpVa1pFZFhJMVExZFhSRmhoZFZSdFhHNUJZbXBYTkhScGQzZ3pjVkZaUWtKTVZUTXpMMVoyYm5sbGJUZGFaSGxNUW1kSmNHTXhUM0phTldsNmNUZHZVMGRxZVZOWVlqWXdTMDAwTmtWcmEzQlVYRzVTV2s1cVQyMDNiRmxJTTNSWVFuSlRabUZYTkhjd1MxUnZNV1pwYW5sR1VWZFZNVFJSYUZneFFXOUhRVXB5YTNCNE9VNHhPVkJNWjBOeksxZzRTbXRXWEc1VVVXVjZhWGg2VkhkRlMzcEhVVUZWTlVweVZXb3ZRbmRZT1VRd1RHdHJaRWxYVjFsd1V6TTJhblo2YWtGRE5qWlBkMGh3VGtKUmRuVjBVVVF2UWxOeFhHNDJWSGRLTDIxWVVEZGFkRk5ZYkZCWGR6QmhNVlpJVFRSdWFFNW5LMncwVjBkd1J6UXJablUzYm1WRE9XeDVTV1pHZG1wdFkyUnpOWGRFVFRWMFF6bFJYRzQzTURsSldFZ3hNMnBuVjI4NVRXTTBNa3hNWkhWTVJUMWNiaTB0TFMwdFJVNUVJRkJTU1ZaQlZFVWdTMFZaTFMwdExTMWNiaUlzQ2lBZ0ltTnNhV1Z1ZEY5bGJXRnBiQ0k2SUNKa2IyTnJaWEpBWTI5dWMzUmhiblF0WTNWaWFYTjBMVEUzTXpFeU15NXBZVzB1WjNObGNuWnBZMlZoWTJOdmRXNTBMbU52YlNJc0NpQWdJbU5zYVdWdWRGOXBaQ0k2SUNJeE1UTXdPRFEzTmpjNE9UVXhOek0wTlRFeU1EZ2lMQW9nSUNKaGRYUm9YM1Z5YVNJNklDSm9kSFJ3Y3pvdkwyRmpZMjkxYm5SekxtZHZiMmRzWlM1amIyMHZieTl2WVhWMGFESXZZWFYwYUNJc0NpQWdJblJ2YTJWdVgzVnlhU0k2SUNKb2RIUndjem92TDI5aGRYUm9NaTVuYjI5bmJHVmhjR2x6TG1OdmJTOTBiMnRsYmlJc0NpQWdJbUYxZEdoZmNISnZkbWxrWlhKZmVEVXdPVjlqWlhKMFgzVnliQ0k2SUNKb2RIUndjem92TDNkM2R5NW5iMjluYkdWaGNHbHpMbU52YlM5dllYVjBhREl2ZGpFdlkyVnlkSE1pTEFvZ0lDSmpiR2xsYm5SZmVEVXdPVjlqWlhKMFgzVnliQ0k2SUNKb2RIUndjem92TDNkM2R5NW5iMjluYkdWaGNHbHpMbU52YlM5eWIySnZkQzkyTVM5dFpYUmhaR0YwWVM5NE5UQTVMMlJ2WTJ0bGNpVTBNR052Ym5OMFlXNTBMV04xWW1semRDMHhOek14TWpNdWFXRnRMbWR6WlhKMmFXTmxZV05qYjNWdWRDNWpiMjBpQ24wPSJ9fX0=
---
apiVersion: v1
kind: Secret
metadata:
  name: client-secret
type: Opaque
data:
  CLIENT_ID: MzM2MzM1NTQxOTkzLWdzZTFyMnZvc3Q1Z2JiMTN0ZWpjYmk0M3UyY3NjYTRpLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29t
  CLIENT_SECRET: ZFBIbmFvbUc3dUNodjFVTWY0bVFuX0tk
`)
	th.writeF("/manifests/tektoncd/tektoncd-tasks/base/service-account.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: tekton-pipelines
imagePullSecrets:
- name: gcr-secret
`)
	th.writeF("/manifests/tektoncd/tektoncd-tasks/base/cluster-role-binding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: tekton-pipelines-controller-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: tekton-pipelines-admin
subjects:
- kind: ServiceAccount
  name: tekton-pipelines
`)
	th.writeF("/manifests/tektoncd/tektoncd-tasks/base/task.yaml", `
---
apiVersion: tekton.dev/v1alpha1
kind: Task
metadata:
  name: build-kfctl-image-from-git-source
spec:
  inputs:
    resources:
    - name: docker-source
      type: git
    params:
    - name: pathToDockerFile
      type: string
      description: The path to the dockerfile to build
      default: /workspace/docker-source/Dockerfile
    - name: pathToContext
      type: string
      description:
        The build context used by Kaniko
        (https://github.com/GoogleContainerTools/kaniko#kaniko-build-contexts)
      default: /workspace/docker-source
  outputs:
    resources:
    - name: builtImage
      type: image
      outputImageDir: /workspace/builtImage
  steps:
  - name: build-and-push
    image: gcr.io/kaniko-project/executor:v0.10.0
    command:
    - /kaniko/executor
    env:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /secret/kaniko-secret.json
    args: ["--dockerfile=${inputs.params.pathToDockerFile}",
           "--destination=${outputs.resources.builtImage.url}",
           "--context=${inputs.params.pathToContext}",
           "--target=kfctl_base"]
    volumeMounts:
    - name: kaniko-secret
      mountPath: /secret
  volumes:
  - name: kaniko-secret
    secret:
      secretName: kaniko-secret
---
apiVersion: tekton.dev/v1alpha1
kind: Task
metadata:
  name: deploy-using-kfctl
spec:
  inputs:
    resources:
    - name: image
      type: image
    params:
    - name: app_dir
      type: string
      description: where to create the kf app
    - name: configPath
      type: string
      description: url for config arg
    - name: project
      type: string
      description: name of project
    - name: zone
      type: string
      description: zone of project
  steps:
  - name: kfctl-init
    image: "${inputs.resources.image.url}"
    command: ["/usr/local/bin/kfctl"]
    args:
    - "init"
    - "--config"
    - "${inputs.params.configPath}"
    - "--project"
    - "${inputs.params.project}"
    - "${inputs.params.app_dir}"
    env:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /secret/kaniko-secret.json
    volumeMounts:
    - name: kaniko-secret
      mountPath: /secret
    - name: kubeflow
      mountPath: /kubeflow
    imagePullPolicy: Always
  - name: kfctl-generate
    image: "${inputs.resources.image.url}"
    imagePullPolicy: Always
    workingDir: "${inputs.params.app_dir}"
    command: ["/usr/local/bin/kfctl"]
    args:
    - "generate"
    - "all"
    - "--zone"
    - "${inputs.params.zone}"
    env:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /secret/kaniko-secret.json
    - name: CLIENT_ID
      valueFrom:
        secretKeyRef:
          name: client-secret
          key: CLIENT_ID
    - name: CLIENT_SECRET
      valueFrom:
        secretKeyRef:
          name: client-secret
          key: CLIENT_SECRET
    volumeMounts:
    - name: kaniko-secret
      mountPath: /secret
    - name: kubeflow
      mountPath: /kubeflow
  - name: kfctl-apply
    image: "${inputs.resources.image.url}"
    imagePullPolicy: Always
    workingDir: "${inputs.params.app_dir}"
    command: ["/usr/local/bin/kfctl"]
    args:
    - "apply"
    - "all"
    - "--verbose"
    env:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /secret/kaniko-secret.json
    - name: CLIENT_ID
      valueFrom:
        secretKeyRef:
          name: client-secret
          key: CLIENT_ID
    - name: CLIENT_SECRET
      valueFrom:
        secretKeyRef:
          name: client-secret
          key: CLIENT_SECRET
    volumeMounts:
    - name: kaniko-secret
      mountPath: /secret
    - name: kubeflow
      mountPath: /kubeflow
  volumes:
  - name: kaniko-secret
    secret:
      secretName: kaniko-secret
  - name: kubeflow
    persistentVolumeClaim:
      claimName: kubeflow-pvc
---
`)
	th.writeK("/manifests/tektoncd/tektoncd-tasks/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- persistent-volume-claim.yaml
- secret.yaml
- service-account.yaml
- cluster-role-binding.yaml
- task.yaml
namespace: tekton-pipelines
`)
}

func TestTektoncdTasksOverlaysApplication(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/tektoncd/tektoncd-tasks/overlays/application")
	writeTektoncdTasksOverlaysApplication(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../tektoncd/tektoncd-tasks/overlays/application"
	fsys := fs.MakeRealFS()
	_loader, loaderErr := loader.NewLoader(targetPath, fsys)
	if loaderErr != nil {
		t.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()))
	kt, err := target.NewKustTarget(_loader, rf, transformer.NewFactoryImpl())
	if err != nil {
		th.t.Fatalf("Unexpected construction error %v", err)
	}
	n, err := kt.MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := n.EncodeAsYaml()
	th.assertActualEqualsExpected(m, string(expected))
}
