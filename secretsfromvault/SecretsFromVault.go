package main

import (
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	rf        *resmap.Factory
	ldr       ifc.Loader
	Name      string   `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string   `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Keys      []string `json:"keys,omitempty" yaml:"keys,omitempty"`
}

//noinspection GoUnusedGlobalVariable
//nolint: golint
var KustomizePlugin plugin

func (p *plugin) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	p.rf = rf
	p.ldr = ldr
	return yaml.Unmarshal(c, p)
}

var database = map[string]string{
	"TREE":      "oak",
	"ROCKET":    "SaturnV",
	"FRUIT":     "apple",
	"VEGETABLE": "carrot",
	"SIMPSON":   "homer",
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace
	for _, k := range p.Keys {
		if v, ok := database[k]; ok {
			args.LiteralSources = append(
				args.LiteralSources, k+"="+v)
		}
	}
	return p.rf.FromSecretArgs(p.ldr, nil, args)
}
