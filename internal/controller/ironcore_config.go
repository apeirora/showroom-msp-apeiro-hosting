package controller

import (
	"encoding/json"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// IronCoreConfig is the chart's scpc.ironcore contract. ShootSpec is the complete
// Gardener Shoot spec, so extension settings need no new operator flags.
type IronCoreConfig struct {
	Autoscaling      *bool                    `json:"autoscaling,omitempty"`
	ProjectNamespace string                   `json:"projectNamespace"`
	Profiles         []map[string]interface{} `json:"profiles"`
	ShootSpec        map[string]interface{}   `json:"shootSpec"`
}

// ConfigureIronCoreManifest validates configuration before the manager starts.
// The returned manifest feeds both the namespace and standalone Shoot callers.
func ConfigureIronCoreManifest(data []byte) ([]byte, error) {
	var cfg IronCoreConfig
	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return nil, fmt.Errorf("scpc.ironcore: %w", err)
	}
	worker, zones, err := cfg.validate()
	if err != nil {
		return nil, err
	}
	inputs, err := cfg.configureInputs(worker, zones)
	if err != nil {
		return nil, err
	}
	worker["minimum"] = "{{ .SC.workers.min }}"
	worker["maximum"] = "{{ .SC.workers.max }}"
	if err := unstructured.SetNestedSlice(cfg.ShootSpec, []interface{}{worker}, "provider", "workers"); err != nil {
		return nil, err
	}
	shootBytes, err := json.Marshal(map[string]interface{}{"apiVersion": "core.gardener.cloud/v1beta1", "kind": "Shoot", "spec": cfg.ShootSpec})
	if err != nil {
		return nil, err
	}
	shootText := string(shootBytes)
	for _, key := range []string{"min", "max"} {
		token := "{{ .SC.workers." + key + " }}"
		shootText = strings.ReplaceAll(shootText, fmt.Sprintf("%q", token), token)
	}
	var obj map[string]interface{}
	if err := yaml.Unmarshal(ironcoreSCPCManifest, &obj); err != nil {
		return nil, err
	}
	scpc := &unstructured.Unstructured{Object: obj}
	if cfg.Autoscaling != nil {
		if err := unstructured.SetNestedField(obj, *cfg.Autoscaling, "spec", "autoscaling"); err != nil {
			return nil, err
		}
	}
	if err := unstructured.SetNestedField(obj, cfg.ProjectNamespace, "spec", "target", "namespace"); err != nil {
		return nil, err
	}
	if err := unstructured.SetNestedSlice(obj, inputs, "spec", "inputs"); err != nil {
		return nil, err
	}
	profiles := make([]interface{}, len(cfg.Profiles))
	for i, p := range cfg.Profiles {
		profiles[i] = p
	}
	if err := unstructured.SetNestedSlice(obj, profiles, "spec", "profiles"); err != nil {
		return nil, err
	}
	resources, _, _ := unstructured.NestedSlice(obj, "spec", "resources")
	resource := resources[0].(map[string]interface{})
	var wrapper map[string]interface{}
	if err := yaml.Unmarshal([]byte(resource["template"].(string)), &wrapper); err != nil {
		return nil, err
	}
	wrapper["spec"] = shootText
	wrapperBytes, err := yaml.Marshal(wrapper)
	if err != nil {
		return nil, err
	}
	resource["template"] = string(wrapperBytes)
	if err := unstructured.SetNestedSlice(obj, resources, "spec", "resources"); err != nil {
		return nil, err
	}
	manifest, err := yaml.Marshal(scpc.Object)
	if err != nil {
		return nil, err
	}
	defaults, err := parseFlavorDefaults(manifest)
	if err != nil {
		return nil, err
	}
	if _, err = renderDefaultShootManifest(defaults, "validation", "validation"); err != nil {
		return nil, fmt.Errorf("scpc.ironcore template: %w", err)
	}
	return manifest, nil
}

func (cfg *IronCoreConfig) validate() (map[string]interface{}, []string, error) {
	if !strings.HasPrefix(cfg.ProjectNamespace, "garden-") || len(cfg.ProjectNamespace) <= len("garden-") {
		return nil, nil, fmt.Errorf("scpc.ironcore.projectNamespace must be a Garden project namespace (garden-<project>)")
	}
	for _, path := range [][]string{{"cloudProfile", "name"}, {"credentialsBindingName"}, {"region"}, {"kubernetes", "version"}, {"networking", "type"}, {"networking", "nodes"}, {"networking", "pods"}, {"networking", "services"}} {
		value, _, err := unstructured.NestedString(cfg.ShootSpec, path...)
		if err != nil || strings.TrimSpace(value) == "" {
			return nil, nil, fmt.Errorf("scpc.ironcore.shootSpec.%s is required", strings.Join(path, "."))
		}
	}
	provider, _, _ := unstructured.NestedString(cfg.ShootSpec, "provider", "type")
	if provider != "ironcore" {
		return nil, nil, fmt.Errorf("scpc.ironcore.shootSpec.provider.type must be ironcore")
	}
	workers, _, _ := unstructured.NestedSlice(cfg.ShootSpec, "provider", "workers")
	if len(workers) != 1 {
		return nil, nil, fmt.Errorf("scpc.ironcore.shootSpec.provider.workers must contain one worker pool")
	}
	worker, ok := workers[0].(map[string]interface{})
	if !ok {
		return nil, nil, fmt.Errorf("scpc.ironcore.shootSpec.provider.workers[0] must be an object")
	}
	for _, path := range [][]string{{"name"}, {"machine", "image", "name"}, {"machine", "image", "version"}, {"volume", "type"}, {"volume", "size"}} {
		value, _, err := unstructured.NestedString(worker, path...)
		if err != nil || strings.TrimSpace(value) == "" {
			return nil, nil, fmt.Errorf("scpc.ironcore.shootSpec.provider.workers[0].%s is required", strings.Join(path, "."))
		}
	}
	zones, _, err := unstructured.NestedStringSlice(worker, "zones")
	if err != nil || len(zones) == 0 {
		return nil, nil, fmt.Errorf("scpc.ironcore.shootSpec.provider.workers[0].zones is required")
	}
	for _, zone := range zones {
		if strings.TrimSpace(zone) == "" {
			return nil, nil, fmt.Errorf("scpc.ironcore worker zones must not be empty")
		}
	}
	return worker, zones, nil
}

func (cfg *IronCoreConfig) configureInputs(worker map[string]interface{}, zones []string) ([]interface{}, error) {
	inputs := []interface{}{}
	for _, input := range []struct {
		name string
		path []string
	}{
		{"cloudProfile", []string{"cloudProfile", "name"}},
		{"region", []string{"region"}},
		{"k8sVersion", []string{"kubernetes", "version"}},
	} {
		value, _, _ := unstructured.NestedString(cfg.ShootSpec, input.path...)
		inputs = append(inputs, map[string]interface{}{"name": input.name, "type": "string", "default": value})
		if err := unstructured.SetNestedField(cfg.ShootSpec, "{{ .Inputs."+input.name+" }}", input.path...); err != nil {
			return nil, err
		}
	}
	imageName, _, _ := unstructured.NestedString(worker, "machine", "image", "name")
	inputs = append(inputs, map[string]interface{}{"name": "machineImage", "type": "string", "default": imageName})
	if err := unstructured.SetNestedField(worker, "{{ .Inputs.machineImage }}", "machine", "image", "name"); err != nil {
		return nil, err
	}
	if len(zones) == 1 {
		inputs = append(inputs, map[string]interface{}{"name": "zone", "type": "string", "default": zones[0]})
		if err := unstructured.SetNestedStringSlice(worker, []string{"{{ .Inputs.zone }}"}, "zones"); err != nil {
			return nil, err
		}
	}
	machineType, _, _ := unstructured.NestedString(worker, "machine", "type")
	if machineType == "" && len(cfg.Profiles) == 0 {
		return nil, fmt.Errorf("scpc.ironcore requires worker machine.type or profiles")
	}
	if len(cfg.Profiles) > 0 {
		for _, profile := range cfg.Profiles {
			for _, key := range []string{"name", "cpu", "memory"} {
				if v, ok := profile[key].(string); !ok || strings.TrimSpace(v) == "" {
					return nil, fmt.Errorf("scpc.ironcore.profiles requires string %s", key)
				}
			}
		}
		inputs = append(inputs, map[string]interface{}{"name": "machineType", "from": "profile"})
	} else {
		inputs = append(inputs, map[string]interface{}{"name": "machineType", "type": "string", "default": machineType})
	}
	if err := unstructured.SetNestedField(worker, "{{ .Inputs.machineType }}", "machine", "type"); err != nil {
		return nil, err
	}
	return inputs, nil
}
