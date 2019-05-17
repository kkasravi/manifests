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

func writeWebhookBase(th *KustTestHarness) {
  th.writeF("/manifests/admission-webhook/webhook/base/cluster-role-binding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cluster-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-role
subjects:
- kind: ServiceAccount
  name: service-account
`)
  th.writeF("/manifests/admission-webhook/webhook/base/cluster-role.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: cluster-role
rules:
- apiGroups:
  - kubeflow.org
  resources:
  - podpresets
  verbs:
  - get
  - watch
  - list
  - update
  - create
  - patch
  - delete
`)
  th.writeF("/manifests/admission-webhook/webhook/base/deployment.yaml", `
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: deployment
spec:
  template:
    spec:
      containers:
      - image: gcr.io/kubeflow-images-public/admission-webhook:v20190502-v0-88-gb5732ba0-dirty-2759ff
        name: admission-webhook
        volumeMounts:
        - mountPath: /etc/webhook/certs
          name: webhook-cert
          readOnly: true
      volumes:
      - name: webhook-cert
        secret:
          secretName: webhook-certs
      serviceAccountName: service-account    
`)
  th.writeF("/manifests/admission-webhook/webhook/base/mutating-webhook-configuration.yaml", `
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  name: mutating-webhook-configuration
webhooks:
- clientConfig:
    caBundle: ""
    service:
      name: $(serviceName)
      namespace: $(namespace)
      path: /apply-podpreset
  name: $(deploymentName).kubeflow.org
  rules:
  - apiGroups:
    - ""
    apiVersions:
    - v1
    operations:
    - CREATE
    resources:
    - pods
`)
  th.writeF("/manifests/admission-webhook/webhook/base/service-account.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: service-account
`)
  th.writeF("/manifests/admission-webhook/webhook/base/service.yaml", `
apiVersion: v1
kind: Service
metadata:
  name: service
spec:
  ports:
  - port: 443
    targetPort: 443
`)
  th.writeF("/manifests/admission-webhook/webhook/base/crd.yaml", `
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: podpresets.kubeflow.org
spec:
  group: kubeflow.org
  names:
    kind: PodPreset
    plural: podpresets
    singular: podpreset
  scope: Namespaced
  version: v1alpha1
  validation:
    openAPIV3Schema:
      properties:
        apiVersion:
          type: string
        kind:
          type: string
        metadata:
          type: object
        spec:
          properties:
            desc:
              type: string
            serviceAccountName:
              type: string
            env:
              items:
                type: object
              type: array
            envFrom:
              items:
                type: object
              type: array
            selector:
              type: object
            volumeMounts:
              items:
                type: object
              type: array
            volumes:
              items:
                type: object
              type: array
          required:
          - selector
          type: object
        status:
          type: object
      type: object
`)
  th.writeF("/manifests/admission-webhook/webhook/base/params.yaml", `
varReference:
- path: webhooks/clientConfig/service/namespace
  kind: MutatingWebhookConfiguration
- path: webhooks/clientConfig/service/name
  kind: MutatingWebhookConfiguration
- path: webhooks/name
  kind: MutatingWebhookConfiguration
`)
  th.writeF("/manifests/admission-webhook/webhook/base/params.env", `
namespace=kubeflow
`)
  th.writeK("/manifests/admission-webhook/webhook/base", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- cluster-role-binding.yaml
- cluster-role.yaml
- deployment.yaml
- mutating-webhook-configuration.yaml
- service-account.yaml
- service.yaml
- crd.yaml
commonLabels:
  kustomize.component: admission-webhook
  app: admission-webhook
namePrefix: admission-webhook- 
images:
  - name: gcr.io/kubeflow-images-public/admission-webhook
    newName: gcr.io/kubeflow-images-public/admission-webhook
    newTag: v20190502-v0-88-gb5732ba0-dirty-2759ff
namespace: kubeflow  
configMapGenerator:
- name: admission-webhook-parameters
  env: params.env
generatorOptions:
  disableNameSuffixHash: true
vars:
- name: namespace
  objref:
    kind: ConfigMap
    name: admission-webhook-parameters 
    apiVersion: v1
  fieldref:
    fieldpath: data.namespace	
- name: serviceName
  objref:
    kind: Service
    name: service
    apiVersion: v1
  fieldref:
    fieldpath: metadata.name
- name: deploymentName
  objref:
    kind: Deployment
    name: deployment
    apiVersion: extensions/v1beta1
  fieldref:
    fieldpath: metadata.name
configurations:
- params.yaml
`)
}

func TestWebhookBase(t *testing.T) {
  th := NewKustTestHarness(t, "/manifests/admission-webhook/webhook/base")
  writeWebhookBase(th)
  m, err := th.makeKustTarget().MakeCustomizedResMap()
  if err != nil {
    t.Fatalf("Err: %v", err)
  }
  targetPath := "../admission-webhook/webhook/base"
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
