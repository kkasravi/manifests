package tests_test

import (
	"sigs.k8s.io/kustomize/v3/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/v3/k8sdeps/transformer"
	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/loader"
	"sigs.k8s.io/kustomize/v3/pkg/plugins"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/target"
	"sigs.k8s.io/kustomize/v3/pkg/validators"
	"testing"
)

func writeCiPipelineRunOverlaysTestingRepoResource(th *KustTestHarness) {
	th.writeF("/manifests/ci/ci-pipeline-run/overlays/testing-repo-resource/pipeline-resource.yaml", `
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: testing
spec:
  type: git
  params:
    - name: revision
      value: $(testing_repo_revision)
    - name: url
      value: $(testing_repo_url)
`)
	th.writeF("/manifests/ci/ci-pipeline-run/overlays/testing-repo-resource/params.yaml", `
varReference:
- path: spec/params/value
  kind: PipelineResource
`)
	th.writeF("/manifests/ci/ci-pipeline-run/overlays/testing-repo-resource/pipeline-run-patch.yaml", `
- op: add
  path: /spec/resources/-
  value:
    name: testing
    resourceRef:
      name: testing
`)
	th.writeF("/manifests/ci/ci-pipeline-run/overlays/testing-repo-resource/params.env", `
testing_repo_revision=master
testing_repo_url=git@github.com:kubeflow/testing.git
`)
	th.writeK("/manifests/ci/ci-pipeline-run/overlays/testing-repo-resource", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
resources:
- pipeline-resource.yaml
patchesJson6902:
- target:
    group: tekton.dev
    version: v1alpha1
    kind: PipelineRun
    name: $(generateName)
  path: pipeline-run-patch.yaml
configMapGenerator:
- name: ci-pipeline-run-parameters
  behavior: merge
  env: params.env
vars:
- name: testing_repo_revision
  objref:
    kind: ConfigMap
    name: ci-pipeline-run-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.testing_repo_revision
- name: testing_repo_url
  objref:
    kind: ConfigMap
    name: ci-pipeline-run-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.testing_repo_url
configurations:
- params.yaml
`)
	th.writeF("/manifests/ci/ci-pipeline-run/base/persistent-volume-claim.yaml", `
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: ci-pipeline-run-persistent-volume-claim
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
`)
	th.writeF("/manifests/ci/ci-pipeline-run/base/service-account.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ci-pipeline-run-service-account
imagePullSecrets:
- name: docker-secret
secrets:
- name: github-secret
`)
	th.writeF("/manifests/ci/ci-pipeline-run/base/secrets.yaml", `
apiVersion: v1
data:
  .dockerconfigjson: eyJhdXRocyI6eyJodHRwczovL2djci5pbyI6eyJ1c2VybmFtZSI6Il9qc29uX2tleSIsInBhc3N3b3JkIjoie1xuICBcInR5cGVcIjogXCJzZXJ2aWNlX2FjY291bnRcIixcbiAgXCJwcm9qZWN0X2lkXCI6IFwiY29uc3RhbnQtY3ViaXN0LTE3MzEyM1wiLFxuICBcInByaXZhdGVfa2V5X2lkXCI6IFwiYWRiMzY3M2NiOTkzNzkyNjZiY2MxZDU1YmIxZTdiZDFlYzM5NGI1Y1wiLFxuICBcInByaXZhdGVfa2V5XCI6IFwiLS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tXFxuTUlJRXZRSUJBREFOQmdrcWhraUc5dzBCQVFFRkFBU0NCS2N3Z2dTakFnRUFBb0lCQVFDNm8zN0o4S2kxUWp3RVxcbnhNT3ROUVZaK2xsWUxIdlNXV2tDeXp1a3JwbHdZRU9KRk5VR00yQ3NySHpjM0pDUDhGYWo1RVRHMjlvT1pLVkJcXG5MSjU3eVdKSEpyekhIb2JyOHNsNytpcjRjYUovSzNiS2lybmZWYTZFeXk5azFIa0RMSlZ4T1lsaXFTbkdtRlZ5XFxuQ3lpYXltNTI1V3VqanZIQkRaZUdsYzlqb1RLMG9yQXYvUCthZzhleUUvY05DS0FwTkk4ZTFXYmlhMFNCdWEwblxcblVZbFB1RXRxdzJ3NDhJbkh6akVQY0VmdENzWjBOZGhkY3hTdVNuSVB5NW9ua2JuVXhZWnAzUjF3TmQ3eDdaQk5cXG5ESmFCWEJTMlVkR1M0ditzeWJVQlU1aXFBckRNbVNmWUwxN09TU3ZzdEZSdVEydkJaa0M3TU96RWd2MUlIZjBXXFxubzlOSzBFaHZBZ01CQUFFQ2dnRUFDUm1MbkZaSzJORHFrdU9kRkJ3dnIwYTdoY2NLeW5pOWhxWURaSERNTTduSVxcbmU5aUkxN2ZpNWgyMWdNeVM1OUcwc21KTGV0UDJwUmtCemFtdjdjMGwwNGp2VDFpM3IxZ0pFWU1Oc1V0VHZFRG1cXG42OUorWkRDTjc3K1FYS21DQ2tZKzRHUmVieHhjV0doNC9MUjZrd0Y5Qi9oV1JTL2xBdlZNc1ZmVjRyK3JTZVNjXFxubU1KOTRBUTROM3hyV0VRc3Vpd1ZIZldMdElMTWZGN1JoV3VzdjJiZ1gvRCs0ajdISHRoODVrYlcxSzR0MnFkN1xcbkIwaEJEcVlQTEtjYzJVNkNJR0NRZ1h3THNlYUUxRkptYWpsdnNVK0pXdmY2MmZTNk8wSlVMeVFLMzZkczZkRlJcXG5qaDg1TWJsZVlHMWdpaFpPTXJtcENvWklUazdFT01lU2pQZ0VaWG5pVVFLQmdRRDNtUHJrZmVKVWNXWmpNZndCXFxuYnJJNE1NRWl2R1JJVDR5RzZxZHZpZFIrd1U3djMxbG1Oa1A3S2s4K1hoQWtGMk1pQkRSakFGVm8rNHQ3a3paRFxcbk45Zk9NSlgwakZnUVpLajFuR1gxekZ2ZERqRTh3ZWRTN2ZMV3BOczJlZm5GWUpQRm5SRU16eG81VWpGNTZ4V1pcXG5ZQmI2VHlNaTNRa1lEblA0WTVjTUZCVUpiUUtCZ1FEQStPNDlUc2EyYWdWUnJYbC8xRDNzWHd0UjdKSGhuRURkXFxuNWlZM0FtOVQxV2pVVTE1T2lwTUxOaXpBb3lRWGlKRVlaMmNuRHZnbHdkNEsvMWFLbityc3hzWUdFZjhoWVBIclxcbkJoN3FueW44SzJseTJoakUxY0xpVFg4NEVnd1VMcFJjeGo3bkM0ZWFLOEdJeUdLNnZrR3NoNCs1bnJLVFlkaUtcXG5MeUhSMUc2cnl3S0JnUURnLzJqSGFNbmEySzRsYUUvTWNXNk05MmtiQ3IzS3BGZGNaeksrZmk3Vy9RMmhsNEtqXFxuQ3A4ZVNDVjQxSHV3Z0h3NmRqMncxYVhINEFheHhtWWlFVVlQL2tEVzJRNVIzMWRXMHNnbzVJdDZSeUpoUndmU1xcbmFaOHFoT2NjQ3gzNXlqaWU5SXVBNjFhMlRrWGR0ODZKOFRNUVJnZjA3NDRMQ1Y5RGtpUzUraW5meFFLQmdFMVdcXG5ObHlacXFmR203VWRPZmxSL1RNeThCMTRHd3I1RFVJaEQ2V3lNeDI5QkpNN2lpc2QvRXBjL3RpQlNXQ3BHY1ZYXFxuQTQ4eXY1NmFNTHZsa3pCaFlNeGQ2VlRiZDQxUUJnUXo0c1lTM2Nlek9rS09SNmp6Sm5SOXJJT3pMK1lTdU9EcFxcbmpxSVlDOU5zdjlacXdLNm91emRDNlFYeUpRMU9CSE4wNmkvbTNDZTdBb0dBU01wRStscDlxV2ZWYXlGV2tlWVBcXG5OOFhId2FNUWNkT0ZkbDZFdlF0ZWtQY0xiQ1F6UzRSdEhBT01NTDN5ci9DQUk5SmZkanhWMHdicW1oNlJ3WFAzXFxuKzhkOVJpNjhsMGV3NUhLMDJWRHFhZE8vOTJhaHNrNmYxV1ZOL0dMcFg4Yk9NZEZFdnJOS09zUVk0RW9DV0JTa1xcblF1ZmRBdFZueE1UZG9ydTNxY0N4RG1vPVxcbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cXG5cIixcbiAgXCJjbGllbnRfZW1haWxcIjogXCJrZi1hY2NvdW50QGNvbnN0YW50LWN1YmlzdC0xNzMxMjMuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb21cIixcbiAgXCJjbGllbnRfaWRcIjogXCIxMDkyODcyODAxMzE5ODQ2MTA2MTZcIixcbiAgXCJhdXRoX3VyaVwiOiBcImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoXCIsXG4gIFwidG9rZW5fdXJpXCI6IFwiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW5cIixcbiAgXCJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmxcIjogXCJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHNcIixcbiAgXCJjbGllbnRfeDUwOV9jZXJ0X3VybFwiOiBcImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkva2YtYWNjb3VudCU0MGNvbnN0YW50LWN1YmlzdC0xNzMxMjMuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb21cIlxufSIsImVtYWlsIjoia2YtYWNjb3VudEBjb25zdGFudC1jdWJpc3QtMTczMTIzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwiYXV0aCI6IlgycHpiMjVmYTJWNU9uc0tJQ0FpZEhsd1pTSTZJQ0p6WlhKMmFXTmxYMkZqWTI5MWJuUWlMQW9nSUNKd2NtOXFaV04wWDJsa0lqb2dJbU52Ym5OMFlXNTBMV04xWW1semRDMHhOek14TWpNaUxBb2dJQ0p3Y21sMllYUmxYMnRsZVY5cFpDSTZJQ0poWkdJek5qY3pZMkk1T1RNM09USTJObUpqWXpGa05UVmlZakZsTjJKa01XVmpNemswWWpWaklpd0tJQ0FpY0hKcGRtRjBaVjlyWlhraU9pQWlMUzB0TFMxQ1JVZEpUaUJRVWtsV1FWUkZJRXRGV1MwdExTMHRYRzVOU1VsRmRsRkpRa0ZFUVU1Q1oydHhhR3RwUnpsM01FSkJVVVZHUVVGVFEwSkxZM2RuWjFOcVFXZEZRVUZ2U1VKQlVVTTJiek0zU2poTGFURlJhbmRGWEc1NFRVOTBUbEZXV2l0c2JGbE1TSFpUVjFkclEzbDZkV3R5Y0d4M1dVVlBTa1pPVlVkTk1rTnpja2g2WXpOS1ExQTRSbUZxTlVWVVJ6STViMDlhUzFaQ1hHNU1TalUzZVZkS1NFcHlla2hJYjJKeU9ITnNOeXRwY2pSallVb3ZTek5pUzJseWJtWldZVFpGZVhrNWF6RklhMFJNU2xaNFQxbHNhWEZUYmtkdFJsWjVYRzVEZVdsaGVXMDFNalZYZFdwcWRraENSRnBsUjJ4ak9XcHZWRXN3YjNKQmRpOVFLMkZuT0dWNVJTOWpUa05MUVhCT1NUaGxNVmRpYVdFd1UwSjFZVEJ1WEc1VldXeFFkVVYwY1hjeWR6UTRTVzVJZW1wRlVHTkZablJEYzFvd1RtUm9aR040VTNWVGJrbFFlVFZ2Ym10aWJsVjRXVnB3TTFJeGQwNWtOM2czV2tKT1hHNUVTbUZDV0VKVE1sVmtSMU0wZGl0emVXSlZRbFUxYVhGQmNrUk5iVk5tV1V3eE4wOVRVM1p6ZEVaU2RWRXlka0phYTBNM1RVOTZSV2QyTVVsSVpqQlhYRzV2T1U1TE1FVm9ka0ZuVFVKQlFVVkRaMmRGUVVOU2JVeHVSbHBMTWs1RWNXdDFUMlJHUW5kMmNqQmhOMmhqWTB0NWJtazVhSEZaUkZwSVJFMU5OMjVKWEc1bE9XbEpNVGRtYVRWb01qRm5UWGxUTlRsSE1ITnRTa3hsZEZBeWNGSnJRbnBoYlhZM1l6QnNNRFJxZGxReGFUTnlNV2RLUlZsTlRuTlZkRlIyUlVSdFhHNDJPVW9yV2tSRFRqYzNLMUZZUzIxRFEydFpLelJIVW1WaWVIaGpWMGRvTkM5TVVqWnJkMFk1UWk5b1YxSlRMMnhCZGxaTmMxWm1WalJ5SzNKVFpWTmpYRzV0VFVvNU5FRlJORTR6ZUhKWFJWRnpkV2wzVmtobVYweDBTVXhOWmtZM1VtaFhkWE4yTW1KbldDOUVLelJxTjBoSWRHZzROV3RpVnpGTE5IUXljV1EzWEc1Q01HaENSSEZaVUV4TFkyTXlWVFpEU1VkRFVXZFlkMHh6WldGRk1VWktiV0ZxYkhaelZTdEtWM1ptTmpKbVV6WlBNRXBWVEhsUlN6TTJaSE0yWkVaU1hHNXFhRGcxVFdKc1pWbEhNV2RwYUZwUFRYSnRjRU52V2tsVWF6ZEZUMDFsVTJwUVowVmFXRzVwVlZGTFFtZFJSRE50VUhKclptVktWV05YV21wTlpuZENYRzVpY2trMFRVMUZhWFpIVWtsVU5IbEhObkZrZG1sa1VpdDNWVGQyTXpGc2JVNXJVRGRMYXpncldHaEJhMFl5VFdsQ1JGSnFRVVpXYnlzMGREZHJlbHBFWEc1T09XWlBUVXBZTUdwR1oxRmFTMm94YmtkWU1YcEdkbVJFYWtVNGQyVmtVemRtVEZkd1RuTXlaV1p1UmxsS1VFWnVVa1ZOZW5odk5WVnFSalUyZUZkYVhHNVpRbUkyVkhsTmFUTlJhMWxFYmxBMFdUVmpUVVpDVlVwaVVVdENaMUZFUVN0UE5EbFVjMkV5WVdkV1VuSlliQzh4UkROeldIZDBVamRLU0dodVJVUmtYRzQxYVZrelFXMDVWREZYYWxWVk1UVlBhWEJOVEU1cGVrRnZlVkZZYVVwRldWb3lZMjVFZG1kc2QyUTBTeTh4WVV0dUszSnplSE5aUjBWbU9HaFpVRWh5WEc1Q2FEZHhibmx1T0VzeWJIa3lhR3BGTVdOTWFWUllPRFJGWjNkVlRIQlNZM2hxTjI1RE5HVmhTemhIU1hsSFN6WjJhMGR6YURRck5XNXlTMVJaWkdsTFhHNU1lVWhTTVVjMmNubDNTMEpuVVVSbkx6SnFTR0ZOYm1FeVN6UnNZVVV2VFdOWE5rMDVNbXRpUTNJelMzQkdaR05hZWtzclptazNWeTlSTW1oc05FdHFYRzVEY0RobFUwTldOREZJZFhkblNIYzJaR295ZHpGaFdFZzBRV0Y0ZUcxWmFVVlZXVkF2YTBSWE1sRTFVak14WkZjd2MyZHZOVWwwTmxKNVNtaFNkMlpUWEc1aFdqaHhhRTlqWTBONE16VjVhbWxsT1VsMVFUWXhZVEpVYTFoa2REZzJTamhVVFZGU1oyWXdOelEwVEVOV09VUnJhVk0xSzJsdVpuaFJTMEpuUlRGWFhHNU9iSGxhY1hGbVIyMDNWV1JQWm14U0wxUk5lVGhDTVRSSGQzSTFSRlZKYUVRMlYzbE5lREk1UWtwTk4ybHBjMlF2UlhCakwzUnBRbE5YUTNCSFkxWllYRzVCTkRoNWRqVTJZVTFNZG14cmVrSm9XVTE0WkRaV1ZHSmtOREZSUW1kUmVqUnpXVk16WTJWNlQydExUMUkyYW5wS2JsSTVja2xQZWt3cldWTjFUMFJ3WEc1cWNVbFpRemxPYzNZNVduRjNTelp2ZFhwa1F6WlJXSGxLVVRGUFFraE9NRFpwTDIwelEyVTNRVzlIUVZOTmNFVXJiSEE1Y1ZkbVZtRjVSbGRyWlZsUVhHNU9PRmhJZDJGTlVXTmtUMFprYkRaRmRsRjBaV3RRWTB4aVExRjZVelJTZEVoQlQwMU5URE41Y2k5RFFVazVTbVprYW5oV01IZGljVzFvTmxKM1dGQXpYRzRyT0dRNVVtazJPR3d3WlhjMVNFc3dNbFpFY1dGa1R5ODVNbUZvYzJzMlpqRlhWazR2UjB4d1dEaGlUMDFrUmtWMmNrNUxUM05SV1RSRmIwTlhRbE5yWEc1UmRXWmtRWFJXYm5oTlZHUnZjblV6Y1dORGVFUnRiejFjYmkwdExTMHRSVTVFSUZCU1NWWkJWRVVnUzBWWkxTMHRMUzFjYmlJc0NpQWdJbU5zYVdWdWRGOWxiV0ZwYkNJNklDSnJaaTFoWTJOdmRXNTBRR052Ym5OMFlXNTBMV04xWW1semRDMHhOek14TWpNdWFXRnRMbWR6WlhKMmFXTmxZV05qYjNWdWRDNWpiMjBpTEFvZ0lDSmpiR2xsYm5SZmFXUWlPaUFpTVRBNU1qZzNNamd3TVRNeE9UZzBOakV3TmpFMklpd0tJQ0FpWVhWMGFGOTFjbWtpT2lBaWFIUjBjSE02THk5aFkyTnZkVzUwY3k1bmIyOW5iR1V1WTI5dEwyOHZiMkYxZEdneUwyRjFkR2dpTEFvZ0lDSjBiMnRsYmw5MWNta2lPaUFpYUhSMGNITTZMeTl2WVhWMGFESXVaMjl2WjJ4bFlYQnBjeTVqYjIwdmRHOXJaVzRpTEFvZ0lDSmhkWFJvWDNCeWIzWnBaR1Z5WDNnMU1EbGZZMlZ5ZEY5MWNtd2lPaUFpYUhSMGNITTZMeTkzZDNjdVoyOXZaMnhsWVhCcGN5NWpiMjB2YjJGMWRHZ3lMM1l4TDJObGNuUnpJaXdLSUNBaVkyeHBaVzUwWDNnMU1EbGZZMlZ5ZEY5MWNtd2lPaUFpYUhSMGNITTZMeTkzZDNjdVoyOXZaMnhsWVhCcGN5NWpiMjB2Y205aWIzUXZkakV2YldWMFlXUmhkR0V2ZURVd09TOXJaaTFoWTJOdmRXNTBKVFF3WTI5dWMzUmhiblF0WTNWaWFYTjBMVEUzTXpFeU15NXBZVzB1WjNObGNuWnBZMlZoWTJOdmRXNTBMbU52YlNJS2ZRPT0ifX19
kind: Secret
metadata:
  name: docker-secret
type: kubernetes.io/dockerconfigjson
---
apiVersion: v1
data:
  kaniko-secret.json: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiY29uc3RhbnQtY3ViaXN0LTE3MzEyMyIsCiAgInByaXZhdGVfa2V5X2lkIjogImFkYjM2NzNjYjk5Mzc5MjY2YmNjMWQ1NWJiMWU3YmQxZWMzOTRiNWMiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQzZvMzdKOEtpMVFqd0VcbnhNT3ROUVZaK2xsWUxIdlNXV2tDeXp1a3JwbHdZRU9KRk5VR00yQ3NySHpjM0pDUDhGYWo1RVRHMjlvT1pLVkJcbkxKNTd5V0pISnJ6SEhvYnI4c2w3K2lyNGNhSi9LM2JLaXJuZlZhNkV5eTlrMUhrRExKVnhPWWxpcVNuR21GVnlcbkN5aWF5bTUyNVd1amp2SEJEWmVHbGM5am9USzBvckF2L1ArYWc4ZXlFL2NOQ0tBcE5JOGUxV2JpYTBTQnVhMG5cblVZbFB1RXRxdzJ3NDhJbkh6akVQY0VmdENzWjBOZGhkY3hTdVNuSVB5NW9ua2JuVXhZWnAzUjF3TmQ3eDdaQk5cbkRKYUJYQlMyVWRHUzR2K3N5YlVCVTVpcUFyRE1tU2ZZTDE3T1NTdnN0RlJ1UTJ2QlprQzdNT3pFZ3YxSUhmMFdcbm85TkswRWh2QWdNQkFBRUNnZ0VBQ1JtTG5GWksyTkRxa3VPZEZCd3ZyMGE3aGNjS3luaTlocVlEWkhETU03bklcbmU5aUkxN2ZpNWgyMWdNeVM1OUcwc21KTGV0UDJwUmtCemFtdjdjMGwwNGp2VDFpM3IxZ0pFWU1Oc1V0VHZFRG1cbjY5SitaRENONzcrUVhLbUNDa1krNEdSZWJ4eGNXR2g0L0xSNmt3RjlCL2hXUlMvbEF2Vk1zVmZWNHIrclNlU2Ncbm1NSjk0QVE0TjN4cldFUXN1aXdWSGZXTHRJTE1mRjdSaFd1c3YyYmdYL0QrNGo3SEh0aDg1a2JXMUs0dDJxZDdcbkIwaEJEcVlQTEtjYzJVNkNJR0NRZ1h3THNlYUUxRkptYWpsdnNVK0pXdmY2MmZTNk8wSlVMeVFLMzZkczZkRlJcbmpoODVNYmxlWUcxZ2loWk9Ncm1wQ29aSVRrN0VPTWVTalBnRVpYbmlVUUtCZ1FEM21QcmtmZUpVY1daak1md0JcbmJySTRNTUVpdkdSSVQ0eUc2cWR2aWRSK3dVN3YzMWxtTmtQN0trOCtYaEFrRjJNaUJEUmpBRlZvKzR0N2t6WkRcbk45Zk9NSlgwakZnUVpLajFuR1gxekZ2ZERqRTh3ZWRTN2ZMV3BOczJlZm5GWUpQRm5SRU16eG81VWpGNTZ4V1pcbllCYjZUeU1pM1FrWURuUDRZNWNNRkJVSmJRS0JnUURBK080OVRzYTJhZ1ZSclhsLzFEM3NYd3RSN0pIaG5FRGRcbjVpWTNBbTlUMVdqVVUxNU9pcE1MTml6QW95UVhpSkVZWjJjbkR2Z2x3ZDRLLzFhS24rcnN4c1lHRWY4aFlQSHJcbkJoN3FueW44SzJseTJoakUxY0xpVFg4NEVnd1VMcFJjeGo3bkM0ZWFLOEdJeUdLNnZrR3NoNCs1bnJLVFlkaUtcbkx5SFIxRzZyeXdLQmdRRGcvMmpIYU1uYTJLNGxhRS9NY1c2TTkya2JDcjNLcEZkY1p6SytmaTdXL1EyaGw0S2pcbkNwOGVTQ1Y0MUh1d2dIdzZkajJ3MWFYSDRBYXh4bVlpRVVZUC9rRFcyUTVSMzFkVzBzZ281SXQ2UnlKaFJ3ZlNcbmFaOHFoT2NjQ3gzNXlqaWU5SXVBNjFhMlRrWGR0ODZKOFRNUVJnZjA3NDRMQ1Y5RGtpUzUraW5meFFLQmdFMVdcbk5seVpxcWZHbTdVZE9mbFIvVE15OEIxNEd3cjVEVUloRDZXeU14MjlCSk03aWlzZC9FcGMvdGlCU1dDcEdjVlhcbkE0OHl2NTZhTUx2bGt6QmhZTXhkNlZUYmQ0MVFCZ1F6NHNZUzNjZXpPa0tPUjZqekpuUjlySU96TCtZU3VPRHBcbmpxSVlDOU5zdjlacXdLNm91emRDNlFYeUpRMU9CSE4wNmkvbTNDZTdBb0dBU01wRStscDlxV2ZWYXlGV2tlWVBcbk44WEh3YU1RY2RPRmRsNkV2UXRla1BjTGJDUXpTNFJ0SEFPTU1MM3lyL0NBSTlKZmRqeFYwd2JxbWg2UndYUDNcbis4ZDlSaTY4bDBldzVISzAyVkRxYWRPLzkyYWhzazZmMVdWTi9HTHBYOGJPTWRGRXZyTktPc1FZNEVvQ1dCU2tcblF1ZmRBdFZueE1UZG9ydTNxY0N4RG1vPVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImtmLWFjY291bnRAY29uc3RhbnQtY3ViaXN0LTE3MzEyMy5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIsCiAgImNsaWVudF9pZCI6ICIxMDkyODcyODAxMzE5ODQ2MTA2MTYiLAogICJhdXRoX3VyaSI6ICJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20vby9vYXV0aDIvYXV0aCIsCiAgInRva2VuX3VyaSI6ICJodHRwczovL29hdXRoMi5nb29nbGVhcGlzLmNvbS90b2tlbiIsCiAgImF1dGhfcHJvdmlkZXJfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9vYXV0aDIvdjEvY2VydHMiLAogICJjbGllbnRfeDUwOV9jZXJ0X3VybCI6ICJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9yb2JvdC92MS9tZXRhZGF0YS94NTA5L2tmLWFjY291bnQlNDBjb25zdGFudC1jdWJpc3QtMTczMTIzLmlhbS5nc2VydmljZWFjY291bnQuY29tIgp9Cg==
kind: Secret
metadata:
  name: kaniko-secret
type: Opaque
---
apiVersion: v1
data:
  CLIENT_ID: MzM2MzM1NTQxOTkzLTJ0NXJsMWMydDQ1czZnb2MzbzQxNWdsNm9uYWpobWt0LmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29tIC1uCg==
  CLIENT_SECRET: ZlFlQjhyOFNUMk5kdXlyOE9UMjVUTE5hIC1uCg==
kind: Secret
metadata:
  name: client-secret
  namespace: kubeflow-ci
type: Opaque
---
apiVersion: v1
kind: Secret
metadata:
  name: github-secret
  annotations:
    tekton.dev/git-0: github.com
type: kubernetes.io/ssh-auth
data:
  known_hosts: Z2l0aHViLmNvbSBzc2gtcnNhIEFBQUFCM056YUMxeWMyRUFBQUFCSXdBQUFRRUFxMkE3aFJHbWRubTl0VURiTzlJRFN3Qks2VGJRYStQWFlQQ1B5NnJiVHJUdHc3UEhrY2NLcnBwMHlWaHA1SGRFSWNLcjZwTGxWREJmT0xYOVFVc3lDT1Ywd3pmaklKTmxHRVlzZGxMSml6SGhibjJtVWp2U0FIUXFaRVRZUDgxZUZ6TFFOblBIdDRFVlZVaDdWZkRFU1U4NEtlem1ENVFsV3BYTG12VTMxL3lNZitTZTh4aEhUdktTQ1pJRkltV3dvRzZtYlVvV2Y5bnpwSW9hU2pCK3dlcXFVVW1wYWFhc1hWYWw3MkorVVgyQisyUlBXM1JjVDBlT3pRZ3FsSkwzUktyVEp2ZHNqRTNKRUF2R3EzbEdIU1pYeTI4RzNza3VhMlNtVmkvdzR5Q0U2Z2JPRHFuVFdsZzcrd0M2MDR5ZEdYQThWSmlTNWFwNDNKWGlVRkZBYVE9PQo=
  ssh-privatekey: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb1FJQkFBS0NBUUVBMVRaNUFoWkNNU2hXUTg1V2dFOWo3d2ZYVzg0b0Y0TU1QUXpVQ1JyNTNSa1pZY2gwCjFaOUQ5bzUxT2FPcnQvMGRma3pSazdIek5sWFU2TmxwcUtESEczSm9RL0hlOUdkTzI0MnVaQ09Ebm5RWHRIOEsKOUFENExvSUF0WUxOcFk5dHEvdVQ0TEwzVDVobFE3VWNTWlR2eDNhUlBmZi9WNW9KL1J5VWpWSm8rcm1tTlZ1ago2dWF4ZlZzOGV4UXpPZnlMaFNlQ0dHN2NpYUtXNGk1MU1EWU1iMHgvc2MyeUZoTGFwLzBBVUJ2VUdZNElVdmNVCm5sYi9UNmNDUzZjVVpldGRwNHR2eXU2R1VSaGI5ZkQ4YUNEakNqL3Y1aFh1ZE8zN3BMbHlsSlhkZG1GeGUrVVoKb0xZS2NNWE9rcjF3Qm5wc29yeE1hU29senJGbTMxZVJoWHhsOFFJQkl3S0NBUUVBc0tsNk5PM2ZFc21tai9RYgp5V1phSGROTUNqMFo1NWgzMnRBZFpxRlNyK0dZc0JzZStpVGdpbllKVzZUZXZRVHNqVGhWNE1hazhvR2h5RDhkCkNCQXY4aXVZT0V3UWdWenBrVmdpeUFBck9pV3RQY2haaUZpVEh6aUVQcCtVY3p4VGpuLytMelV6VUpRMnBkQ1QKeC9CdTV4a1pRZmxlaW1td2dVT1JCMkdLTHRMUzZJTHFRS0p6S0Iwd282TERsbS8rajNvZzhLb09sdGVUbG9zQwpHRVUvOGluNUFZYStEa2hSS0EzODhGdnFjWlZ6RHpkUDEyUEYwODBHVmkyQnIwTFZNTFpSeU9YOGs4WjlFaWs0CjM0SFhaYzF1MXk1RnZkWlhuRjJWQWN2a3RmMzhKcWNod1F4T0VDNWZRUGtRRkNIMDJmeFdacnhDQjlQcE8rUE4KTVg4WlN3S0JnUUQxTGI4UENMVGQ3SUpkRXhuTGNZU0tBRzdKanZUckkvS2llNVg1S1U3TlV3K3J1QmhtblNEawo3eVA0dHJqYUwwZG5taEVOTnQ4SDNwbDZIYUxCWUNmNjhpVm1laEtid3lOTm1leVRjcmFRMXNXNlZRUGE1N1NKClFVQlBoQXRZOWNBYldmcm5pU2YrNUtnVVdUQmlTZFdsek5wakFQYlVxTTAzM1JFRC9BenNXUUtCZ1FEZW40dlMKY09DR1lHOEZ6cC9scVBtcSt0NjBGNUk0RUwxdWFrajk2WngvTE8rd1IxQjN2L3ZJdlhpekJJb2N3a3VqR1VyLwpVZHc5cG1XSlhCUUY5SEs1cjVMU3YzKzV1RXdmSXNuVTN0WnhqSGx5QUNEWjNSTnlsUm5aS05zUFhsMGx6L1JZCnpYWmQ0b1NvWGVUZjJVK29KeEo0VFZBYjhFNWNJK3ovMGh5eldRS0JnSDRYaHRTSUlvREN4ck9HS29YaW0rOHoKYkMwazh2eVdKUU1hL0t3R25ZNHF0NW9rS2M1ZmNBQzFnRGJFWHhFZm9SQzFvbDZSUDRCajJlY0g3VTJCNkt6eApqNVBKd0cxZEN0ZHp1MzhIeTVyeUhJdTJzWWFGeW8rN0tHTm9ldDFMTDVrSnREeVByaXRSQmZ2V0dPSVhWKzdlCllhOUNVdzVJTHdiSmRuY1Q2V09iQW9HQVRGUDFiTGo4aGRmcmpQWml0U1A5MURqQllsRTVjazdoM0xhNThLZmQKNG5YT2haVExJYjRxY0xYOWZ6VEJvM1hReWpTT3ZmQXVRUVhacTNBRzNYaGgyVU9DeStuaXBoTk5USUQwdmdOQwpOWTgvbE16WUVDNktVeXZPV1JWU2E2dFRJdWcyaE5pNjNseEt2V0lGVkEvK0RjUkEyTW9NMW1FRTdHUVd2aHduCkdPc0NnWUJrZVFGMFY0V1hLYXVqejM4QkI4NTA0WENsMUhYZjkxSjgvOUlKeW13M3NzdFA3aVVhVUNuZ1A4dEkKNDFJdnRMUmJGTUV6QU93WUVNTUVTSXh5QS9vOUhVSDN0djh6c0pVOTNBNlF1QmN5dU1reWsxNkxLT2lERkc2QgpQR0ZYQi9XaC9SVlhFTFpyNWRqb1lxQWg4ZUw5aW05WUVaOENXaTlVd1JXN2pvdW92Zz09Ci0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0tCg==
  ssh-publickey: c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBQkl3QUFBUUVBMVRaNUFoWkNNU2hXUTg1V2dFOWo3d2ZYVzg0b0Y0TU1QUXpVQ1JyNTNSa1pZY2gwMVo5RDlvNTFPYU9ydC8wZGZrelJrN0h6TmxYVTZObHBxS0RIRzNKb1EvSGU5R2RPMjQydVpDT0RublFYdEg4SzlBRDRMb0lBdFlMTnBZOXRxL3VUNExMM1Q1aGxRN1VjU1pUdngzYVJQZmYvVjVvSi9SeVVqVkpvK3JtbU5WdWo2dWF4ZlZzOGV4UXpPZnlMaFNlQ0dHN2NpYUtXNGk1MU1EWU1iMHgvc2MyeUZoTGFwLzBBVUJ2VUdZNElVdmNVbmxiL1Q2Y0NTNmNVWmV0ZHA0dHZ5dTZHVVJoYjlmRDhhQ0RqQ2ovdjVoWHVkTzM3cExseWxKWGRkbUZ4ZStVWm9MWUtjTVhPa3Ixd0JucHNvcnhNYVNvbHpyRm0zMWVSaFh4bDhRPT0ga2thc3JhdmlAMTkyLjE2OC4xLjIK
---
apiVersion: v1
kind: Secret
metadata:
  name: github-token-secret
type: Opaque
data:
  token: MDg5N2ZhMmI0N2Y3Y2NiNzJkMjdiMTBkMjAzODUxMDBhMjQ4YWM2Yw==
`)
	th.writeF("/manifests/ci/ci-pipeline-run/base/cluster-role-binding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: ci-pipeline-run-cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: ci-pipeline-run-service-account
`)
	th.writeF("/manifests/ci/ci-pipeline-run/base/pipeline-run.yaml", `
apiVersion: tekton.dev/v1alpha1
kind: PipelineRun
metadata:
  name: $(generateName)
  labels:
    scope: $(namespace)
spec:
  serviceAccount: ci-pipeline-run-service-account
  pipelineRef:
    name: $(pipeline)
  resources: []
`)
	th.writeF("/manifests/ci/ci-pipeline-run/base/params.yaml", `
varReference:
- path: metadata/name
  kind: PipelineRun
- path: metadata/labels/scope
  kind: PipelineRun
- path: metadata/namespace
  kind: PersistentVolumeClaim
- path: spec/pipelineRef/name
  kind: PipelineRun
`)
	th.writeF("/manifests/ci/ci-pipeline-run/base/params.env", `
namespace=
generateName=
pipeline=
`)
	th.writeK("/manifests/ci/ci-pipeline-run/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- persistent-volume-claim.yaml
- service-account.yaml
- secrets.yaml
- cluster-role-binding.yaml
- pipeline-run.yaml
namespace: $(namespace)
configMapGenerator:
- name: ci-pipeline-run-parameters
  env: params.env
vars:
- name: namespace
  objref:
    kind: ConfigMap
    name: ci-pipeline-run-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.namespace
- name: generateName
  objref:
    kind: ConfigMap
    name: ci-pipeline-run-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.generateName
- name: pipeline
  objref:
    kind: ConfigMap
    name: ci-pipeline-run-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.pipeline
configurations:
- params.yaml
`)
}

func TestCiPipelineRunOverlaysTestingRepoResource(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/ci/ci-pipeline-run/overlays/testing-repo-resource")
	writeCiPipelineRunOverlaysTestingRepoResource(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.AsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../ci/ci-pipeline-run/overlays/testing-repo-resource"
	fsys := fs.MakeRealFS()
	lrc := loader.RestrictionRootOnly
	_loader, loaderErr := loader.NewLoader(lrc, validators.MakeFakeValidator(), targetPath, fsys)
	if loaderErr != nil {
		t.Fatalf("could not load kustomize loader: %v", loaderErr)
	}
	rf := resmap.NewFactory(resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl()), transformer.NewFactoryImpl())
	pc := plugins.DefaultPluginConfig()
	kt, err := target.NewKustTarget(_loader, rf, transformer.NewFactoryImpl(), plugins.NewLoader(pc, rf))
	if err != nil {
		th.t.Fatalf("Unexpected construction error %v", err)
	}
	actual, err := kt.MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(actual, string(expected))
}
