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

func writeTektoncdPipelinesBase(th *KustTestHarness) {
	th.writeF("/manifests/tektoncd/tektoncd-pipelines/base/cluster-role-binding.yaml", `
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
  namespace: $(namespace)
`)
	th.writeF("/manifests/tektoncd/tektoncd-pipelines/base/pipeline-resource.yaml", `
---
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: kfctl-git
spec:
  type: git
  params:
  - name: revision
    value: $(pullrequest)
  - name: url
    value: https://github.com/kubeflow/kfctl.git
---
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: kfctl-image
spec:
  type: image
  params:
  - name: url
    value: gcr.io/$(project)/kfctl
`)
	th.writeF("/manifests/tektoncd/tektoncd-pipelines/base/pipeline.yaml", `
apiVersion: tekton.dev/v1alpha1
kind: Pipeline
metadata:
  name: kfctl-build-apply
spec:
  resources:
  - name: source-repo
    type: git
  - name: web-image
    type: image
  tasks:
  - name: kfctl-build-push
    taskRef:
      name: build-kfctl-image-from-git-source
    params:
    - name: pathToDockerFile
      value: /workspace/docker-source/Dockerfile
    - name: pathToContext
      value: /workspace/docker-source
    resources:
      inputs:
      - name: docker-source
        resource: source-repo
      outputs:
      - name: builtImage
        resource: web-image
  - name: kfctl-init-generate-apply
    taskRef:
      name: deploy-using-kfctl
    resources:
      inputs:
      - name: image
        resource: web-image
        from:
        - kfctl-build-push
    params:
    - name: app_dir
      value: $(app_dir)
    - name: project
      value: $(project)
    - name: configPath
      value: $(configPath)
    - name: zone
      value: $(zone)
`)
	th.writeF("/manifests/tektoncd/tektoncd-pipelines/base/pipeline-run.yaml", `
apiVersion: tekton.dev/v1alpha1
kind: PipelineRun
metadata:
  name: kfctl-build-apply-pipeline-run
spec:
  serviceAccount: tekton-pipelines
  pipelineRef:
    name: kfctl-build-apply
  resources:
    - name: source-repo
      resourceRef:
        name: kfctl-git
    - name: web-image
      resourceRef:
        name: kfctl-image
`)
	th.writeF("/manifests/tektoncd/tektoncd-pipelines/base/params.yaml", `
varReference:
- path: spec/tasks/params/value
  kind: Pipeline
- path: spec/params/value
  kind: PipelineResource
- path: subjects/namespace
  kind: ClusterRoleBinding
`)
	th.writeF("/manifests/tektoncd/tektoncd-pipelines/base/params.env", `
project=constant-cubist-173123
namespace=tekton-pipelines
pullrequest=refs/pull/10/head
app_dir=/kubeflow/kubeflow-e2e
zone=us-west1-a
configPath=https://raw.githubusercontent.com/kkasravi/kfctl/tektoncd_pipelines/config/kfctl_e2e.yaml
`)
	th.writeK("/manifests/tektoncd/tektoncd-pipelines/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- cluster-role-binding.yaml
- pipeline-resource.yaml
- pipeline.yaml
- pipeline-run.yaml
namespace: tekton-pipelines
configMapGenerator:
- name: kfctl-pipelines-parameters
  env: params.env
vars:
- name: namespace
  objref:
    kind: ConfigMap
    name: kfctl-pipelines-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.namespace
- name: project
  objref:
    kind: ConfigMap
    name: kfctl-pipelines-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.project
- name: configPath
  objref:
    kind: ConfigMap
    name: kfctl-pipelines-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.configPath
- name: pullrequest
  objref:
    kind: ConfigMap
    name: kfctl-pipelines-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.pullrequest
- name: app_dir
  objref:
    kind: ConfigMap
    name: kfctl-pipelines-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.app_dir
- name: zone
  objref:
    kind: ConfigMap
    name: kfctl-pipelines-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.zone
configurations:
- params.yaml
`)
}

func TestTektoncdPipelinesBase(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/tektoncd/tektoncd-pipelines/base")
	writeTektoncdPipelinesBase(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../tektoncd/tektoncd-pipelines/base"
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
