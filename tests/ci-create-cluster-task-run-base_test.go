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

func writeCiCreateClusterTaskRunBase(th *KustTestHarness) {
	th.writeF("/manifests/ci/ci-create-cluster-task-run/base/persistent-volume-claim.yaml", `
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: persistent-volume-claim
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
`)
	th.writeF("/manifests/ci/ci-create-cluster-task-run/base/secret.yaml", `
---
apiVersion: v1
kind: Secret
metadata:
  name: kaniko-secret
type: Opaque
data:
  ci-create-cluster-kaniko-secret.json: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiY29uc3RhbnQtY3ViaXN0LTE3MzEyMyIsCiAgInByaXZhdGVfa2V5X2lkIjogImQxMjkzODQ1NWI2NGEyZDlhZWE1MDVjNjZkNzIyMjJmNmUyNDg0MzYiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQ2dJRHptZmFqMkdKSnVcbjlvOUxsa3dDYXIvVXYyalN3MEJTRXpNeWozQlNCc3BSdXVhOG1qU3MyZTdNRjVnYmMxVGdoL21HSkNqMitZS3JcbkZDOHdyMEt5bWx6NDdRWE5aeUMxb2I2TitxUnJncEJyeXNZWE1iSkFnQUFQck1rS2lhY2VjRzBHN3JkVGhjQVNcblh2Y2hBZG9RdTFTTWhNQys4aEtHMytpM1ZETXE2K3o4ZStwL1AzeEhhdUl3UTBYcWp4RGpONE5VRlMrbTF3WG9cbmpQM2xkbE15cWJUZHFvdFNHanIrTUxJZ1U3OEowdlp1NDF5WDl1ZjY2QWJyOE91N2E4bFdpVnpJUzZubHVqSTJcbm00T3NWYXlZaUN6dHhOam5iazVCVUhaYXZpekNaLzE1YzRlWFRuYXlsQXI0S2IrS3ZHSk5XbXVjSUdjQlNkZnlcbjhGZE1aNEdSQWdNQkFBRUNnZ0VBQWhXaFdpQ1Y4d3d4dmMvMHh1UmdQVzgzQlM4OGdveXVjcmw1OWpIOGgvd0lcbkpSN3VHWWJ3bE9nUXg3Ukg3UldPVWx0YXkzYVl2b1h0ZitPOW9JYXJqTnVRclBubUlUQW1QMDhQMzRpV3RnOWdcbjZRV0UvdkROSU10VFlhMm9NdStyVHRMNzhzTG5zQ2ExMWxkaTE0QUNPS3R3Ymk0UWIwaTJ3Q0VJK1Y1a0ZsVElcbmlkZkFiWGdiakt6dDR0dFg5YmR3OVBpUEcrTkRxQkVpQ0FiaXNMYmxyWlVOWUl1WVVsVWpoTmNKS29hRWNXQkdcbkMzcS9BMk1oTElTZ2FVQlYyZnludXEyaVdVZFZQTGwxNTRNNWtMREFSNk8rS0J6UWVZVnUraGxzdVcrSE1qcGFcbjI5U05FSWhuWmYySERPTi82bUw3R1VEcVlOblhISDVCK1AvT3NnSVQ5UUtCZ1FEV3R4SWNXSks3NFdadmhYeGNcbktOR1pLdHlxQWxielhSVFI3M1NGTVpSQjZNUjdEOFloenIwN3dmRXI5UEk5ekl3QjB3V3BzR2VGRVFQUFpTYXlcbnlKaDQxcTkvYlZ2bUtiM2svRndrMmVFSy9JRFBNRzdqTHdqbVBBbG5UYVpVRUZncHVKYUw5VEVvYlBSWXZtbGVcbmpyeXNmNVRpWXBIY3RkVlhkL3IwZzRqbDNRS0JnUUMrNmg4a28yNVkrS3BPRitxaGU3UnA1Z1VZdFhwQjh6RnpcbjB1UnZacE5DQXFMUmw3Vjk5RW9YcnBZWFF5TE5CY3BSUkpDaEpVenZDMkRtWVhhRXhLUzlWbmZYSzRncG5uVzZcbnAzSitPQ0NtOUJOdUMrYStxYXp0WDluK1JrRXpaem0waFdjTFRTejVpUlhrOWkxTkRKZnY4aXBKd1ptMnhPNDBcbkNXdU5pZ2R4UlFLQmdEbllhRkNxckIxaHhDOFhUMEdrM1pMZU1VUzhES0RUMnVBVUd0Z25XMEhHYStpYmYwMXNcblhSN1VTUjBHaUp5TmxzcUhCMmVIMXR2S2tiUTJGQTdtYSsxaUtUV3pTS2JoWi85ZzNaSXdBS2p0RGViRHJad1dcbjk5YlBKZGxtMmdDYnhxUzJ6aGcybmwrOXVyYU4xZVZibndqNTlpcG5VOVNhU0RlZ1kwT3NqQjBoQW9HQkFKUzRcbnN5d0tlRktzMjVaY1FUWXN0TDF1SjNnNUh4VXpDc29NZGxGbDJiOHBhSWJYcE5XS3NSRkR1cjVDV1dEWGF1VG1cbkFialc0dGl3eDNxUVlCQkxVMzMvVnZueWVtN1pkeUxCZ0lwYzFPclo1aXpxN29TR2p5U1hiNjBLTTQ2RWtrcFRcblJaTmpPbTdsWUgzdFhCclNmYVc0dzBLVG8xZmlqeUZRV1UxNFFoWDFBb0dBSnJrcHg5TjE5UExnQ3MrWDhKa1ZcblRRZXppeHpUd0VLekdRQVU1SnJVai9Cd1g5RDBMa2tkSVdXWXBTMzZqdnpqQUM2Nk93SHBOQlF2dXRRRC9CU3FcbjZUd0ovbVhQN1p0U1hsUFd3MGExVkhNNG5oTmcrbDRXR3BHNCtmdTduZUM5bHlJZkZ2am1jZHM1d0RNNXRDOVFcbjcwOUlYSDEzamdXbzlNYzQyTExkdUxFPVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImRvY2tlckBjb25zdGFudC1jdWJpc3QtMTczMTIzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAiY2xpZW50X2lkIjogIjExMzA4NDc2Nzg5NTE3MzQ1MTIwOCIsCiAgImF1dGhfdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwKICAidG9rZW5fdXJpIjogImh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwKICAiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0cyIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZG9ja2VyJTQwY29uc3RhbnQtY3ViaXN0LTE3MzEyMy5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIKfQo=
---
apiVersion: v1
kind: Secret
metadata:
  name: docker-secret
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: eyJhdXRocyI6eyJodHRwczovL2djci5pbyI6eyJ1c2VybmFtZSI6Il9qc29uX2tleSIsInBhc3N3b3JkIjoie1xuICBcInR5cGVcIjogXCJzZXJ2aWNlX2FjY291bnRcIixcbiAgXCJwcm9qZWN0X2lkXCI6IFwiY29uc3RhbnQtY3ViaXN0LTE3MzEyM1wiLFxuICBcInByaXZhdGVfa2V5X2lkXCI6IFwiZTMzZDhhZDQ4OWZkYTEzMTg0ZmQxZDYzZmVjMDhjY2RhZTlkZTE5MFwiLFxuICBcInByaXZhdGVfa2V5XCI6IFwiLS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tXFxuTUlJRXZnSUJBREFOQmdrcWhraUc5dzBCQVFFRkFBU0NCS2d3Z2dTa0FnRUFBb0lCQVFDV0xaWW5iU3BuZk1Uc1xcblRsSktZU2lUNWxvd1pveWRPZndMNEQxM0dCUnlneml2aWpDOUtValE4a2t2Y05NeEJweHRSUUtNYVd2UjNoaHZcXG5XejE4eC9xZW1DdlhuV1FMMngxNzJzeTVwcTRsVnlaUzJKbi9jNHhJQVNGNHhXb3lvQVpsNXpFZk8wem01M3BkXFxuNCtqd0NMcDdDRSs1UVJTclBjdXhIdkNvZUhKeXhxWGtmOGJoTVpXN1BJQUVNMWlZSEdNb01ua1VOdVFiU0VvTlxcbkw0SHJTcUgwZnV5RkVLNjFRYUJKR3RSdloyZVZTUmhkUEFaRTFEQW5OWFVldXI0NTE3bkhXQTl5WVRTQVdnRTZcXG41U01TSEVxWVVMSmdMMzd2SktiUUFwSjAwZ25nNnBNTFJEanB0aFdBa0lBL3BYamxNR2lmMnZsdnY5REVRN3p1XFxuSFRwU2dEZWJBZ01CQUFFQ2dnRUFTTTN0MHN4SjkrU1ZaUWZ0UGZEUEpyQlFQZEdoVHFHckxxaTFzNVJCYVdoelxcbkpTcWx5VGFIL2YvUGVnZkU0cW9WVUtYWmYrK2xuUmNDR280TmgzNDlZZ0JjbE1sUkZLeFRwVlVqMWNiWCt2TStcXG5lWUJYVysrTTdPVmJjRHlvYU1XS2hJRnBuMzMwb0taTWZOTDkvTXdHZDVuR2FJV0QreVpZcHRQY2tKZmZ5QU1HXFxuT1dBTTd6cXdrRHg3RDV4Umg2T0Nzb1kydlNwUGh2cmJzbUxOdTVlYm5TWUxPWFJlV0ZvR1JjN1JOQi9LOHRVUFxcbkxueTVNR25jVGFvRXhOUWR5ZkM0Y3B3M3prZHhqM2NUdXVGVFJzb2tBdHFSRElSUVM2MFl6bTZQMi9vdFNWKytcXG5QN3NMeVdMVUdRZ3ZqTkZNdkZQRGJzSVE5Rm14bUpHb25MVjRpQUFaOVFLQmdRRFFQckpsc2U5U1hZT1ZrQzdmXFxuWEtnK2xmb1luNjBUbTI2YmJKNWV5WC84TlF0Y1BlSFZTcXk4cGRZVnpiR25Xckd2Qzh3dzREdFJFdVVZcytYMFxcblFRZU5TNHNJUFBDZHNqNXBMRjl6d2lGVERReWl0S016RXg2NDlSRnFteVRmV3Z2Z0RueUFqWlJJeU9nS3NTWnJcXG5NOGc1eVY1NVlDdlRPTlovUWYvR2g4NmJkUUtCZ1FDNG5mOFBWWHFRRlI2YWczNHhqNURqc2hXZlhJT3ZaWTg3XFxuclNzK3BUSDJxdWxYYzhudEhZQVBPZktnVGJmZWs1cGRSNStnRzN5WWcrbUcrWE5WQlRheDJXbDVWT2tZU0VxSlxcbkhVZDFrQ25PbVlUZVp0Ri8vdTdzektwRmE4L0hkcUtMNFdETkZaeU9iZ1NYcWdXbGFVM2FWWkhkK3V2ZzM3SEVcXG56M05oUlpIMHp3S0JnUUNzRG1GUGJNaVRnUGdySnNuVGVyYjNudXJZVlhXbThaRmRrVXo0ZS92bTRkelZCYndGXFxuZ29GZURKYnB4TjIzckZPS2tYRFFJVFJoTS85ZGZhWE5QYjJEbkpydTM0cmVnRnJZZ3ZVS3E2YmsrNjhvNzU2M1xcbm9HQ042TTNQQ3doWUV0Qnd1d2RiSDU4WTFBWUViNEdTcVdJUmZMTTJEYU9vRFJvTVl2ZDFqTmZEMFFLQmdDQjNcXG4rUkd6VU5qaVBmMml2cURzeE9pbXUxTEpySWMrYjFCcGhqK0FRaWRGcThBcnB3bkN0SEQ1R2dqRFltRU15SXM3XFxuTzRHbkUrU20zbjFVaGNvZ0hweHN4allHanZBc1ZwK0N2THlhWEIvdnRBU0JSTHNrRk5Va3NaVi8vb3p2K21wclxcbmV1RFd1aS82ZldoSENMTXNyL3FFTGlGQ0xoWGdnWjFCZHVOV252TFZBb0dCQUxLcTZLcmdhNk5oWWpEY3hkN3FcXG5tU2FzRm9qdjN3SDUralBwVmU0akIyeW0xeDJGN2JqQmVWMWdpZDNPVWJqQkRSTllWSldZUkZNNzFxc0xyRUY5XFxuTEhidVNVRnorSVZNOXFnOUhzT3NqWDh3ekI1QVh2bG9Vclo0a1Y3RmJLOUp5WnRXU2RYMnpNNllmY1lJSkdiRlxcbjVHUU9HM3hKM1hZZEJKTm52T1YwcXIvS1xcbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cXG5cIixcbiAgXCJjbGllbnRfZW1haWxcIjogXCJkb2NrZXJAY29uc3RhbnQtY3ViaXN0LTE3MzEyMy5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbVwiLFxuICBcImNsaWVudF9pZFwiOiBcIjExMzA4NDc2Nzg5NTE3MzQ1MTIwOFwiLFxuICBcImF1dGhfdXJpXCI6IFwiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGhcIixcbiAgXCJ0b2tlbl91cmlcIjogXCJodHRwczovL29hdXRoMi5nb29nbGVhcGlzLmNvbS90b2tlblwiLFxuICBcImF1dGhfcHJvdmlkZXJfeDUwOV9jZXJ0X3VybFwiOiBcImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0c1wiLFxuICBcImNsaWVudF94NTA5X2NlcnRfdXJsXCI6IFwiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9kb2NrZXIlNDBjb25zdGFudC1jdWJpc3QtMTczMTIzLmlhbS5nc2VydmljZWFjY291bnQuY29tXCJcbn0iLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJhdXRoIjoiWDJwemIyNWZhMlY1T25zS0lDQWlkSGx3WlNJNklDSnpaWEoyYVdObFgyRmpZMjkxYm5RaUxBb2dJQ0p3Y205cVpXTjBYMmxrSWpvZ0ltTnZibk4wWVc1MExXTjFZbWx6ZEMweE56TXhNak1pTEFvZ0lDSndjbWwyWVhSbFgydGxlVjlwWkNJNklDSmxNek5rT0dGa05EZzVabVJoTVRNeE9EUm1aREZrTmpObVpXTXdPR05qWkdGbE9XUmxNVGt3SWl3S0lDQWljSEpwZG1GMFpWOXJaWGtpT2lBaUxTMHRMUzFDUlVkSlRpQlFVa2xXUVZSRklFdEZXUzB0TFMwdFhHNU5TVWxGZG1kSlFrRkVRVTVDWjJ0eGFHdHBSemwzTUVKQlVVVkdRVUZUUTBKTFozZG5aMU5yUVdkRlFVRnZTVUpCVVVOWFRGcFpibUpUY0c1bVRWUnpYRzVVYkVwTFdWTnBWRFZzYjNkYWIzbGtUMlozVERSRU1UTkhRbEo1WjNwcGRtbHFRemxMVldwUk9HdHJkbU5PVFhoQ2NIaDBVbEZMVFdGWGRsSXphR2gyWEc1WGVqRTRlQzl4WlcxRGRsaHVWMUZNTW5neE56SnplVFZ3Y1RSc1ZubGFVekpLYmk5ak5IaEpRVk5HTkhoWGIzbHZRVnBzTlhwRlprOHdlbTAxTTNCa1hHNDBLMnAzUTB4d04wTkZLelZSVWxOeVVHTjFlRWgyUTI5bFNFcDVlSEZZYTJZNFltaE5XbGMzVUVsQlJVMHhhVmxJUjAxdlRXNXJWVTUxVVdKVFJXOU9YRzVNTkVoeVUzRklNR1oxZVVaRlN6WXhVV0ZDU2tkMFVuWmFNbVZXVTFKb1pGQkJXa1V4UkVGdVRsaFZaWFZ5TkRVeE4yNUlWMEU1ZVZsVVUwRlhaMFUyWEc0MVUwMVRTRVZ4V1ZWTVNtZE1NemQyU2t0aVVVRndTakF3WjI1bk5uQk5URkpFYW5CMGFGZEJhMGxCTDNCWWFteE5SMmxtTW5ac2RuWTVSRVZSTjNwMVhHNUlWSEJUWjBSbFlrRm5UVUpCUVVWRFoyZEZRVk5OTTNRd2MzaEtPU3RUVmxwUlpuUlFaa1JRU25KQ1VWQmtSMmhVY1VkeVRIRnBNWE0xVWtKaFYyaDZYRzVLVTNGc2VWUmhTQzltTDFCbFoyWkZOSEZ2VmxWTFdGcG1LeXRzYmxKalEwZHZORTVvTXpRNVdXZENZMnhOYkZKR1MzaFVjRlpWYWpGallsZ3JkazByWEc1bFdVSllWeXNyVFRkUFZtSmpSSGx2WVUxWFMyaEpSbkJ1TXpNd2IwdGFUV1pPVERrdlRYZEhaRFZ1UjJGSlYwUXJlVnBaY0hSUVkydEtabVo1UVUxSFhHNVBWMEZOTjNweGQydEVlRGRFTlhoU2FEWlBRM052V1RKMlUzQlFhSFp5WW5OdFRFNTFOV1ZpYmxOWlRFOVlVbVZYUm05SFVtTTNVazVDTDBzNGRGVlFYRzVNYm5rMVRVZHVZMVJoYjBWNFRsRmtlV1pETkdOd2R6TjZhMlI0YWpOalZIVjFSbFJTYzI5clFYUnhVa1JKVWxGVE5qQlplbTAyVURJdmIzUlRWaXNyWEc1UU4zTk1lVmRNVlVkUlozWnFUa1pOZGtaUVJHSnpTVkU1Um0xNGJVcEhiMjVNVmpScFFVRmFPVkZMUW1kUlJGRlFja3BzYzJVNVUxaFpUMVpyUXpkbVhHNVlTMmNyYkdadldXNDJNRlJ0TWpaaVlrbzFaWGxZTHpoT1VYUmpVR1ZJVmxOeGVUaHdaRmxXZW1KSGJsZHlSM1pET0hkM05FUjBVa1YxVlZseksxZ3dYRzVSVVdWT1V6UnpTVkJRUTJSemFqVndURVk1ZW5kcFJsUkVVWGxwZEV0TmVrVjROalE1VWtaeGJYbFVabGQyZG1kRWJubEJhbHBTU1hsUFowdHpVMXB5WEc1Tk9HYzFlVlkxTlZsRGRsUlBUbG92VVdZdlIyZzRObUprVVV0Q1oxRkRORzVtT0ZCV1dIRlJSbEkyWVdjek5IaHFOVVJxYzJoWFpsaEpUM1phV1RnM1hHNXlVM01yY0ZSSU1uRjFiRmhqT0c1MFNGbEJVRTltUzJkVVltWmxhelZ3WkZJMUsyZEhNM2xaWnl0dFJ5dFlUbFpDVkdGNE1sZHNOVlpQYTFsVFJYRktYRzVJVldReGEwTnVUMjFaVkdWYWRFWXZMM1UzYzNwTGNFWmhPQzlJWkhGTFREUlhSRTVHV25sUFltZFRXSEZuVjJ4aFZUTmhWbHBJWkN0MWRtY3pOMGhGWEc1Nk0wNW9VbHBJTUhwM1MwSm5VVU56UkcxR1VHSk5hVlJuVUdkeVNuTnVWR1Z5WWpOdWRYSlpWbGhYYlRoYVJtUnJWWG8wWlM5MmJUUmtlbFpDWW5kR1hHNW5iMFpsUkVwaWNIaE9Nak55Ums5TGExaEVVVWxVVW1oTkx6bGtabUZZVGxCaU1rUnVTbkoxTXpSeVpXZEdjbGxuZGxWTGNUWmlheXMyT0c4M05UWXpYRzV2UjBOT05rMHpVRU4zYUZsRmRFSjNkWGRrWWtnMU9Ga3hRVmxGWWpSSFUzRlhTVkptVEUweVJHRlBiMFJTYjAxWmRtUXhhazVtUkRCUlMwSm5RMEl6WEc0clVrZDZWVTVxYVZCbU1tbDJjVVJ6ZUU5cGJYVXhURXB5U1dNcllqRkNjR2hxSzBGUmFXUkdjVGhCY25CM2JrTjBTRVExUjJkcVJGbHRSVTE1U1hNM1hHNVBORWR1UlN0VGJUTnVNVlZvWTI5blNIQjRjM2hxV1VkcWRrRnpWbkFyUTNaTWVXRllRaTkyZEVGVFFsSk1jMnRHVGxWcmMxcFdMeTl2ZW5ZcmJYQnlYRzVsZFVSWGRXa3ZObVpYYUVoRFRFMXpjaTl4UlV4cFJrTk1hRmhuWjFveFFtUjFUbGR1ZGt4V1FXOUhRa0ZNUzNFMlMzSm5ZVFpPYUZscVJHTjRaRGR4WEc1dFUyRnpSbTlxZGpOM1NEVXJhbEJ3Vm1VMGFrSXllVzB4ZURKR04ySnFRbVZXTVdkcFpETlBWV0pxUWtSU1RsbFdTbGRaVWtaTk56RnhjMHh5UlVZNVhHNU1TR0oxVTFWR2VpdEpWazA1Y1djNVNITlBjMnBZT0hkNlFqVkJXSFpzYjFWeVdqUnJWamRHWWtzNVNubGFkRmRUWkZneWVrMDJXV1pqV1VsS1IySkdYRzQxUjFGUFJ6TjRTak5ZV1dSQ1NrNXVkazlXTUhGeUwwdGNiaTB0TFMwdFJVNUVJRkJTU1ZaQlZFVWdTMFZaTFMwdExTMWNiaUlzQ2lBZ0ltTnNhV1Z1ZEY5bGJXRnBiQ0k2SUNKa2IyTnJaWEpBWTI5dWMzUmhiblF0WTNWaWFYTjBMVEUzTXpFeU15NXBZVzB1WjNObGNuWnBZMlZoWTJOdmRXNTBMbU52YlNJc0NpQWdJbU5zYVdWdWRGOXBaQ0k2SUNJeE1UTXdPRFEzTmpjNE9UVXhOek0wTlRFeU1EZ2lMQW9nSUNKaGRYUm9YM1Z5YVNJNklDSm9kSFJ3Y3pvdkwyRmpZMjkxYm5SekxtZHZiMmRzWlM1amIyMHZieTl2WVhWMGFESXZZWFYwYUNJc0NpQWdJblJ2YTJWdVgzVnlhU0k2SUNKb2RIUndjem92TDI5aGRYUm9NaTVuYjI5bmJHVmhjR2x6TG1OdmJTOTBiMnRsYmlJc0NpQWdJbUYxZEdoZmNISnZkbWxrWlhKZmVEVXdPVjlqWlhKMFgzVnliQ0k2SUNKb2RIUndjem92TDNkM2R5NW5iMjluYkdWaGNHbHpMbU52YlM5dllYVjBhREl2ZGpFdlkyVnlkSE1pTEFvZ0lDSmpiR2xsYm5SZmVEVXdPVjlqWlhKMFgzVnliQ0k2SUNKb2RIUndjem92TDNkM2R5NW5iMjluYkdWaGNHbHpMbU52YlM5eWIySnZkQzkyTVM5dFpYUmhaR0YwWVM5NE5UQTVMMlJ2WTJ0bGNpVTBNR052Ym5OMFlXNTBMV04xWW1semRDMHhOek14TWpNdWFXRnRMbWR6WlhKMmFXTmxZV05qYjNWdWRDNWpiMjBpQ24wPSJ9fX0=
---
apiVersion: v1
kind: Secret
metadata:
  name: client-secret
type: Opaque
data:
  CLIENT_ID: MzM2MzM1NTQxOTkzLWdzZTFyMnZvc3Q1Z2JiMTN0ZWpjYmk0M3UyY3NjYTRpLmFwcHMuZ29vZ2xldXNlcmNvbnRlbnQuY29t
  CLIENT_SECRET: ZFBIbmFvbUc3dUNodjFVTWY0bVFuX0tk
---
apiVersion: v1
kind: Secret
metadata:
  name: gcp-secret
type: Opaque
data:
  ci-create-cluster-gcp-secret.json: ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAiY29uc3RhbnQtY3ViaXN0LTE3MzEyMyIsCiAgInByaXZhdGVfa2V5X2lkIjogIjllYjBhN2I3YTI3NDJiYTc2ODJjNWM2OGZiYzllODVmZmViOTM3OTUiLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRRFpRVUlFTlFOK0FEZndcbjZBSDRlbmZTakZVcG1YY2Q3ZGJRQ1R2YVY5VVFMT2JhSDJJS1V4Q0t4WkVoNWZDcWZUeDVJWVpUdDZLckFabExcbklNU0JTWmRqd2NLa293SmpPcTZraStQbFJ5YlUrNVpMTk9PNXZOR200dEtNbXVDc2hkTmVHaXVlbThKZ3k2ZFhcbmZzSzR0K3JzVk5kK1JDSUNoMXJEa2VXb0ZPSXhLSzJuSGtIbjVWU0VsTmxjUGJuZUt4dVVNOTl6T1J4MldlRTFcbk9pVjNjTUgvRGhGNWpVNVAvRVJlUVU0YTk3WnhkN01DSmdCZlV4Ky96K0w5dWNJRkJjOWhZQlA4cGNPbGdFbktcblJMSytkZEhaMlJDZU5XaUFRdC9oeEl6UUczbG5PaVJLakViZ2RMU28yQVdvbXBLalI2SGp6NGh0ekdtNTM2R0dcbkdwTmlLN0pyQWdNQkFBRUNnZ0VBR3gvOU5KWkQ0ZHI0SVJGdWtZNEU0TnBabGJDT0FVUWRRbXNzdUdXbitmV0pcbk95bVk3WTRTYmlrZHBqeFYwSXVEWGVKVUthYXZYaWQ4Y3JkY0lZSkZMeFRWandXMU9odHRDNmxWb2w1QVdHNHpcbkJSL000UGRVdTcvdEp0WDlnRHpUTjVnUDR5VXlYekIrSzd2dFp1KzdtcGM3TW80aUt1dW9adXVUMzJrQUZyL3hcbkNCTERoQjAyd3VaTzFRc2JIdzZIdXdsK2RUL2NPOHgyVFp0Y3pBcXFST29ZdXVjZllsMCtIZ29iK3FXbG5ZYTdcbkJVQ0NudHUyNmVFUDhjV2p4dVM4UGh4R0NhOU94QlFpWE8xOVBWR1RXWFNTUGlPMmRGeHpQdWxUNExjaEdpWlBcbkhuWEJRZE1saXIvY2U3Rm5nVlpsNm1GQUdYSStnQVBYN09XWEp1L211UUtCZ1FEeTdHZStha0xjSkwzaURidGtcbmx1WVJscCtXSUtnMFoyMkxDeHdWK05Yb3dwbG9vRjdmVTBzN0c5ZWxneU5xYzdmNElqSlEvVWtTYmErYzNpRldcbjJYeGh6OFF5K0xZQnlEUFZaYkFCaWJPUjlSZkU2dGxFbXVPOTN4SHZqZklJOE96MGMzRzRjTmQrNU9WbXYxNWNcbldmK0Z1TmFpTzJjbndadGtJY1VwWVlRRUl3S0JnUURrOHlEbVcvYWlXMlVYcFV5WUFjY3lWUVJ5cU9QODI5T2FcbkFiUzgvdUU3cFFqSHFsVzBXV09LWm42T3BoT0VxUy9ieWR6aG9nSmMzUDJ4UkxyRStlNmFoQy9EdTVPQzM0SG9cbnVwOXJEWGdRaG5IdWpnVnBaV25vV3JBbU5ubnAwQmlPRU95N0dEd0orOGY0WTExZWkweDdjUWhlOUt3dmRMSnVcbmxRTzlYUzY1R1FLQmdBZWtwSWI3Tk90VVJKMHVMVzAyeWpwWGNPSDZXUkI5Q0pkTlhDN2N5MjR0WVVKSGVYU3hcblhEYVo0NmtUZlRQR1BFMlVWZHp5ZXpBWFAyVkNIKzVwblY0K2VUL1pUM0N5NmQ2VytuaXg2bkozTWE1Q2JWK1pcbk4vMHJYWmNaOGptUnl3TE45eEFFak9Nek5IeU5ITnp1Ly9rbkhhbXhFTWZSY2FBdTU5TXJmRW5kQW9HQkFOczRcbmNhZ2hKbWNQWEJ6b0NnOENwTmxzem5WN2dkSDhLd0NyNFlPV0NkUXlrZFdkSTdNc1pFT0JJRzAyV0RvT1JlVU5cbnhKSEhycnQ4WHUzK0FWZmFlTDA3RlFFMSttaTEybzRkSThnOWZWbFZZb0lwT3NWUWRiZ21IY1I1SlFMY1hxYXBcblRnTlhrU1YrRUZ1bHlTRmVBRDJ5WFhHT2xkQmF6UDlWYjk5Qitoc0JBb0dCQUxqVE44VFYveG5yb05IdjNSNDRcbnRTUHpkUjRLQ0dZV01QU1cvZ3NWN3c4T0xJcjRYWS9rSkFPVnRjL0RCWFN3TmtBdjFpcDQwS0pjeDMrdVJqdWxcbkVzalRjVGJlNHpobUR1SzZIejhYUGJRdFZmSmMwdDkzN0M0SFdxd1QyWW41LzZBeE13MExjcXdJa1VkYjRTL3RcbjhMcVUweG5lRk51aEZEUllkSlRzbTR5bVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImtmY3RsLWUyZUBjb25zdGFudC1jdWJpc3QtMTczMTIzLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwKICAiY2xpZW50X2lkIjogIjEwNDAxMDY4MzczNjg5NDk5NzkzNSIsCiAgImF1dGhfdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi9hdXRoIiwKICAidG9rZW5fdXJpIjogImh0dHBzOi8vb2F1dGgyLmdvb2dsZWFwaXMuY29tL3Rva2VuIiwKICAiYXV0aF9wcm92aWRlcl94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL29hdXRoMi92MS9jZXJ0cyIsCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkva2ZjdGwtZTJlJTQwY29uc3RhbnQtY3ViaXN0LTE3MzEyMy5pYW0uZ3NlcnZpY2VhY2NvdW50LmNvbSIKfQo=
`)
	th.writeF("/manifests/ci/ci-create-cluster-task-run/base/service-account.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: service-account
imagePullSecrets:
- name: docker-secret
`)
	th.writeF("/manifests/ci/ci-create-cluster-task-run/base/cluster-role-binding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: service-account
`)
	th.writeF("/manifests/ci/ci-create-cluster-task-run/base/task-run.yaml", `
apiVersion: tekton.dev/v1alpha1
kind: TaskRun
metadata:
  name: $(generateName)
spec:
  serviceAccount: ci-create-cluster-service-account
  inputs:
    params:
    - name: namespace
      value: $(namespace)
    - name: app_dir
      value: $(app_dir)
    - name: project
      value: $(project)
    - name: configPath
      value: $(configPath)
    - name: zone
      value: $(zone)
    - name: email
      value: $(email)
    - name: platform
      value: $(platform)
    - name: cluster
      value: $(cluster)
    - name: kfctl_image
      value: $(kfctl_image)
    - name: pvc_mount_path
      value: $(pvc_mount_path)
  taskSpec:
    inputs:
      params:
      - name: kfctl_image
        type: string
        description: the kfctl container image
      - name: namespace
        type: string
        description: the namespace to deploy kf 
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
      - name: platform
        type: string
        description: all | k8s
      - name: email
        type: string
        description: email for gcp
      - name: cluster
        type: string
        description: name of the cluster
      - name: pvc_mount_path
        type: string
        description: parent dir for kfctl
    steps:
    - name: kfctl-activate-service-account
      image: "${inputs.params.kfctl_image}"
      imagePullPolicy: IfNotPresent
      workingDir: "${inputs.params.pvc_mount_path}"
      command: ["/opt/google-cloud-sdk/bin/gcloud"]
      args:
      - "auth"
      - "activate-service-account"
      - "--key-file"
      - "/secret/ci-create-cluster-gcp-secret.json"
      env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/ci-create-cluster-gcp-secret.json
      - name: CLIENT_ID
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_ID
      - name: CLIENT_SECRET
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_SECRET
      volumeMounts:
      - name: ci-create-cluster-gcp-secret
        mountPath: /secret
      - name: kubeflow
        mountPath: /kubeflow
    - name: kfctl-set-account
      image: "${inputs.params.kfctl_image}"
      imagePullPolicy: IfNotPresent
      workingDir: "${inputs.params.pvc_mount_path}"
      command: ["/opt/google-cloud-sdk/bin/gcloud"]
      args:
      - "config"
      - "set"
      - "account"
      - "${inputs.params.email}"
      env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/ci-create-cluster-gcp-secret.json
      - name: CLIENT_ID
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_ID
      - name: CLIENT_SECRET
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_SECRET
      volumeMounts:
      - name: ci-create-cluster-gcp-secret
        mountPath: /secret
      - name: kubeflow
        mountPath: /kubeflow
    - name: kfctl-init
      image: "${inputs.params.kfctl_image}"
      workingDir: "${inputs.params.pvc_mount_path}"
      command: ["/usr/local/bin/kfctl"]
      args:
      - "init"
      - "--config"
      - "${inputs.params.configPath}"
      - "--project"
      - "${inputs.params.project}"
      - "--namespace"
      - "${inputs.params.namespace}"
      - "${inputs.params.app_dir}"
      env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/ci-create-cluster-kaniko-secret.json
      volumeMounts:
      - name: ci-create-cluster-kaniko-secret
        mountPath: /secret
      - name: kubeflow
        mountPath: "${inputs.params.pvc_mount_path}"
      imagePullPolicy: IfNotPresent
    - name: kfctl-generate
      image: "${inputs.params.kfctl_image}"
      imagePullPolicy: IfNotPresent
      workingDir: "${inputs.params.pvc_mount_path}/${inputs.params.app_dir}"
      command: ["/usr/local/bin/kfctl"]
      args:
      - "generate"
      - "${inputs.params.platform}"
      - "--zone"
      - "${inputs.params.zone}"
      - "--email"
      - "${inputs.params.email}"
      env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/ci-create-cluster-kaniko-secret.json
      - name: CLIENT_ID
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_ID
      - name: CLIENT_SECRET
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_SECRET
      volumeMounts:
      - name: ci-create-cluster-kaniko-secret
        mountPath: /secret
      - name: kubeflow
        mountPath: /kubeflow
    - name: kfctl-apply
      image: "${inputs.params.kfctl_image}"
      imagePullPolicy: IfNotPresent
      workingDir: "${inputs.params.pvc_mount_path}/${inputs.params.app_dir}"
  #    command: ["/bin/sleep", "infinity"]
      command: ["/usr/local/bin/kfctl"]
      args:
      - "apply"
      - "${inputs.params.platform}"
      - "--verbose"
      env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/ci-create-cluster-gcp-secret.json
      - name: CLIENT_ID
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_ID
      - name: CLIENT_SECRET
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_SECRET
      volumeMounts:
      - name: ci-create-cluster-gcp-secret
        mountPath: /secret
      - name: kubeflow
        mountPath: /kubeflow
    - name: kfctl-configure-kubectl
      image: "${inputs.params.kfctl_image}"
      imagePullPolicy: IfNotPresent
      workingDir: "${inputs.params.pvc_mount_path}"
      command: ["/opt/google-cloud-sdk/bin/gcloud"]
      args:
      - "--project"
      - "${inputs.params.project}"
      - "container"
      - "clusters"
      - "--zone"
      - "${inputs.params.zone}"
      - "get-credentials"
      - "${inputs.params.cluster}"
      env:
      - name: GOOGLE_APPLICATION_CREDENTIALS
        value: /secret/ci-create-cluster-gcp-secret.json
      - name: CLIENT_ID
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_ID
      - name: CLIENT_SECRET
        valueFrom:
          secretKeyRef:
            name: ci-create-cluster-client-secret
            key: CLIENT_SECRET
      volumeMounts:
      - name: ci-create-cluster-gcp-secret
        mountPath: /secret
      - name: kubeflow
        mountPath: /kubeflow
    volumes:
    - name: ci-create-cluster-kaniko-secret
      secret:
        secretName: ci-create-cluster-kaniko-secret
    - name: ci-create-cluster-docker-secret
      secret:
        secretName: ci-create-cluster-docker-secret
    - name: ci-create-cluster-gcp-secret
      secret:
        secretName: ci-create-cluster-gcp-secret
    - name: kubeflow
      persistentVolumeClaim:
        claimName: ci-create-cluster-persistent-volume-claim
`)
	th.writeF("/manifests/ci/ci-create-cluster-task-run/base/params.yaml", `
varReference:
- path: metadata/name
  kind: TaskRun
- path: spec/inputs/params/value
  kind: TaskRun
`)
	th.writeF("/manifests/ci/ci-create-cluster-task-run/base/params.env", `
namespace=kubeflow-ci
project=constant-cubist-173123
app_dir=/kubeflow/kubeflow-ci
zone=us-west1-a
email=foo@bar.com
configPath=https://raw.githubusercontent.com/kubeflow/kubeflow/master/bootstrap/config/ci-cluster.yaml
platform=all
cluster=kubeflow-ci
pvc_mount_path=/kubeflow
kfctl_image=gcr.io/constant-cubist-173123/kfctl@sha256:ab0c4986322e3e6a755056278c7270983b0f3bdc0751aefff075fb2b3d0c3254
`)
	th.writeK("/manifests/ci/ci-create-cluster-task-run/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- persistent-volume-claim.yaml
- secret.yaml
- service-account.yaml
- cluster-role-binding.yaml
- task-run.yaml
namespace: tekton-pipelines
namePrefix: ci-create-cluster-
configMapGenerator:
- name: parameters
  env: params.env
vars:
- name: namespace
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.namespace
- name: project
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.project
- name: configPath
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.configPath
- name: app_dir
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.app_dir
- name: zone
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.zone
- name: email
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.email
- name: platform
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.platform
- name: cluster
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.cluster
- name: kfctl_image
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.kfctl_image
- name: pvc_mount_path
  objref:
    kind: ConfigMap
    name: parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.pvc_mount_path
configurations:
- params.yaml
`)
}

func TestCiCreateClusterTaskRunBase(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/ci/ci-create-cluster-task-run/base")
	writeCiCreateClusterTaskRunBase(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.AsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../ci/ci-create-cluster-task-run/base"
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
