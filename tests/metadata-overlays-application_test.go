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

func writeMetadataOverlaysApplication(th *KustTestHarness) {
	th.writeF("/manifests/metadata/overlays/application/application.yaml", `
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  name: $(generateName)
spec:
  componentKinds:
  - group: core
    kind: Service
  - group: apps
    kind: Deployment
  - group: core
    kind: ConfigMap
  - group: core
    kind: ServiceAccount
  descriptor:
    type: "metadata"
    version: "alpha"
    description: "Tracking and managing metadata of machine learning workflows in Kubeflow."
    maintainers:
    - name: Zhenghui Wang
      email: zhenghui@google.com
    owners:
    - name: Ajay Gopinathan
      email: ajaygopinathan@google.com
    - name: Zhenghui Wang
      email: zhenghui@google.com
    keywords:
    - "metadata"
    links:
    - description: Docs
      url: "https://www.kubeflow.org/docs/components/misc/metadata/"
  addOwnerRef: true
`)
	th.writeF("/manifests/metadata/overlays/application/params.yaml", `
varReference:
- path: metadata/name
  kind: Application
`)
	th.writeF("/manifests/metadata/overlays/application/params.env", `
generateName=
`)
	th.writeK("/manifests/metadata/overlays/application", `
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
bases:
- ../../base
resources:
- application.yaml
configMapGenerator:
- name: metadata-parameters
  env: params.env
vars:
- name: generateName
  objref:
    kind: ConfigMap
    name: metadata-parameters
    apiVersion: v1
  fieldref:
    fieldpath: data.generateName
configurations:
- params.yaml
commonLabels:
  app.kubernetes.io/name: metadata
  app.kubernetes.io/instance: $(generateName)
  app.kubernetes.io/managed-by: kfctl
  app.kubernetes.io/component: metadata
  app.kubernetes.io/part-of: kubeflow
  app.kubernetes.io/version: v0.6
`)
	th.writeF("/manifests/metadata/base/metadata-db-pvc.yaml", `
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
`)
	th.writeF("/manifests/metadata/base/metadata-db-secret.yaml", `
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: db-secrets
data:
  MYSQL_ROOT_PASSWORD: dGVzdA== # "test"
`)
	th.writeF("/manifests/metadata/base/metadata-db-deployment.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: db
  labels:
    component: db
spec:
  replicas: 1
  template:
    metadata:
      name: db
      labels:
        component: db
    spec:
      containers:
      - name: db-container
        image: mysql:8.0.3
        args:
        - --datadir
        - /var/lib/mysql/datadir
        env:
          - name: MYSQL_ROOT_PASSWORD
            valueFrom:
              secretKeyRef:
                name: metadata-db-secrets
                key: MYSQL_ROOT_PASSWORD
          - name: MYSQL_ALLOW_EMPTY_PASSWORD
            value: "true"
          - name: MYSQL_DATABASE
            value: "metadb"
        ports:
        - name: dbapi
          containerPort: 3306
        readinessProbe:
          exec:
            command:
            - "/bin/bash"
            - "-c"
            - "mysql -D $$MYSQL_DATABASE -p$$MYSQL_ROOT_PASSWORD -e 'SELECT 1'"
          initialDelaySeconds: 5
          periodSeconds: 2
          timeoutSeconds: 1
        volumeMounts:
        - name: metadata-mysql
          mountPath: /var/lib/mysql
      volumes:
      - name: metadata-mysql
        persistentVolumeClaim:
          claimName: metadata-mysql
`)
	th.writeF("/manifests/metadata/base/metadata-db-service.yaml", `
apiVersion: v1
kind: Service
metadata:
  name: db
  labels:
    component: db
spec:
  type: ClusterIP
  ports:
    - port: 3306
      protocol: TCP
      name: dbapi
  selector:
    component: db
`)
	th.writeF("/manifests/metadata/base/metadata-deployment.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment
  labels:
    component: server
spec:
  replicas: 3
  selector:
    matchLabels:
      component: server
  template:
    metadata:
      labels:
        component: server
    spec:
      containers:
      - name: container
        image: gcr.io/kubeflow-images-public/metadata:v0.1.8
        env:
          - name: MYSQL_ROOT_PASSWORD
            valueFrom:
              secretKeyRef:
                name: metadata-db-secrets
                key: MYSQL_ROOT_PASSWORD
        command: ["./server/server",
                  "--http_port=8080",
                  "--mysql_service_host=metadata-db.kubeflow",
                  "--mysql_service_port=3306",
                  "--mysql_service_user=root",
                  "--mysql_service_password=$(MYSQL_ROOT_PASSWORD)",
                  "--mlmd_db_name=metadb"]
        ports:
        - name: backendapi
          containerPort: 8080
`)
	th.writeF("/manifests/metadata/base/metadata-service.yaml", `
kind: Service
apiVersion: v1
metadata:
  labels:
    app: metadata
  name: service
spec:
  selector:
    component: server
  type: ClusterIP
  ports:
  - port: 8080
    protocol: TCP
    name: backendapi
`)
	th.writeF("/manifests/metadata/base/metadata-ui-deployment.yaml", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ui
  labels:
    app: metadata-ui
spec:
  selector:
    matchLabels:
      app: metadata-ui
  template:
    metadata:
      name: ui
      labels:
        app: metadata-ui
    spec:
      containers:
      - image: gcr.io/kubeflow-images-public/metadata-frontend:v0.1.8
        imagePullPolicy: IfNotPresent
        name: metadata-ui
        ports:
        - containerPort: 3000
      serviceAccountName: ui

`)
	th.writeF("/manifests/metadata/base/metadata-ui-role.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: Role
metadata:
  labels:
    app: metadata-ui
  name: ui
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - pods/log
  verbs:
  - create
  - get
  - list
- apiGroups:
  - "kubeflow.org"
  resources:
  - viewers
  verbs:
  - create
  - get
  - list
  - watch
  - delete
`)
	th.writeF("/manifests/metadata/base/metadata-ui-rolebinding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  labels:
    app: metadata-ui
  name: ui
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: ui
subjects:
- kind: ServiceAccount
  name: ui
  namespace: kubeflow
`)
	th.writeF("/manifests/metadata/base/metadata-ui-sa.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: ui
`)
	th.writeF("/manifests/metadata/base/metadata-ui-service.yaml", `
apiVersion: v1
kind: Service
metadata:
  name: ui
  labels:
    app: metadata-ui
spec:
  ports:
  - port: 80
    targetPort: 3000
  selector:
    app: metadata-ui
`)
	th.writeF("/manifests/metadata/base/params.env", `
uiClusterDomain=cluster.local
`)
	th.writeK("/manifests/metadata/base", `
namePrefix: metadata-
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
commonLabels:
  kustomize.component: metadata
configMapGenerator:
- name: ui-parameters
  env: params.env
resources:
- metadata-db-pvc.yaml
- metadata-db-secret.yaml
- metadata-db-deployment.yaml
- metadata-db-service.yaml
- metadata-deployment.yaml
- metadata-service.yaml
- metadata-ui-deployment.yaml
- metadata-ui-role.yaml
- metadata-ui-rolebinding.yaml
- metadata-ui-sa.yaml
- metadata-ui-service.yaml
namespace: kubeflow
vars:
- name: ui-namespace
  objref:
    kind: Service
    name: ui
    apiVersion: v1
  fieldref:
    fieldpath: metadata.namespace
- name: ui-clusterDomain
  objref:
    kind: ConfigMap
    name: ui-parameters
    version: v1
  fieldref:
    fieldpath: data.uiClusterDomain
- name: service
  objref:
    kind: Service
    name: ui
    apiVersion: v1
  fieldref:
    fieldpath: metadata.name`)
}

func TestMetadataOverlaysApplication(t *testing.T) {
	th := NewKustTestHarness(t, "/manifests/metadata/overlays/application")
	writeMetadataOverlaysApplication(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	expected, err := m.AsYaml()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	targetPath := "../metadata/overlays/application"
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
