// +build plugin

//go:generate go run github/kubeflow/manifests/plugins/kustomize/plugin/testgenerator
package main

import (
	"apps.kubeflow.org/v1alpha1"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	ldr         ifc.Loader
	rf          *resmap.Factory
	kfdef       *v1alpha1.KfDef
	types.GeneratorOptions
}

var KustomizePlugin plugin

func (p *plugin) Config(ldr ifc.Loader, rf *resmap.Factory, buf []byte) error {
	p.kfdef = &v1alpha1.KfDef{}
	p.GeneratorOptions = types.GeneratorOptions{}
	p.ldr = ldr
	p.rf = rf
	return yaml.Unmarshal(buf, p.Application)
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	buf, err := yaml.Marshal(p.Application)
	if err != nil {
		return nil, err
	}
	return p.rf.NewResMapFromBytes(buf)
}

// generator functionality
func (p *plugin) genTargetStart() {

}

func (p *plugin) genTargetMiddle() {

}

func (p *plugin) genTargetEnd() {

}

func (p *plugin)genTarget() {

}

func (p *plugin) genTargetBase() {

}

func (p *plugin) genTargetKustomization() {

}

func (p *plugin) Transform(m resmap.ResMap) error {
        t, err := transformers.NewNamePrefixSuffixTransformer(
                p.Prefix,
                p.Suffix,
                p.FieldSpecs,
        )
        if err != nil {
                return err
        }
        return t.Transform(m)
}

