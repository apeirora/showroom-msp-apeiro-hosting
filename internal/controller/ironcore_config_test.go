package controller

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func gardenDevConfig(t interface {
	Helper()
	Fatal(...interface{})
}) []byte {
	t.Helper()
	data, err := os.ReadFile("../../charts/sapprovider-operator/examples/garden-dev.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var values struct {
		SCPC struct {
			IronCore IronCoreConfig `json:"ironcore"`
		} `json:"scpc"`
	}
	if err := yaml.Unmarshal(data, &values); err != nil {
		t.Fatal(err)
	}
	data, err = yaml.Marshal(values.SCPC.IronCore)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestIronCoreConfigurationRequired(t *testing.T) {
	if _, err := ConfigureIronCoreManifest([]byte("{}")); err == nil || !strings.Contains(err.Error(), "projectNamespace") {
		t.Fatalf("expected project error, got %v", err)
	}
	if _, err := parseFlavorDefaults(ironcoreSCPCManifest); err == nil {
		t.Fatal("unconfigured skeleton must not have landscape defaults")
	}
	for _, old := range []string{"garden-pm-cc-d2", "wdf4-ora-dev", "1.34.5", "10Gi"} {
		if strings.Contains(string(ironcoreSCPCManifest), old) {
			t.Fatalf("old default remains: %s", old)
		}
	}
}

func TestIronCoreConfigurationMissingFields(t *testing.T) {
	for _, path := range [][]string{{"cloudProfile", "name"}, {"credentialsBindingName"}, {"kubernetes", "version"}, {"region"}, {"networking", "pods"}} {
		t.Run(strings.Join(path, "."), func(t *testing.T) {
			var config IronCoreConfig
			if err := yaml.Unmarshal(gardenDevConfig(t), &config); err != nil {
				t.Fatal(err)
			}
			unstructured.RemoveNestedField(config.ShootSpec, path...)
			data, err := yaml.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := ConfigureIronCoreManifest(data); err == nil || !strings.Contains(err.Error(), strings.Join(path, ".")) {
				t.Fatalf("expected missing field error, got %v", err)
			}
		})
	}
}

func checkGardenDevManifest(t *testing.T, data []byte) {
	t.Helper()
	manifest, err := ConfigureIronCoreManifest(data)
	if err != nil {
		t.Fatal(err)
	}
	defaults, err := parseFlavorDefaults(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if defaults.targetNamespace != "garden-dev" {
		t.Fatal(defaults.targetNamespace)
	}
	shoot, err := renderDefaultShootManifest(defaults, "account", "shoot")
	if err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{"spec.cloudProfile.name": "ironcore", "spec.credentialsBindingName": "ironcore-infra", "spec.kubernetes.version": "1.33.13", "spec.region": "wdf4", "spec.networking.type": "calico"} {
		got, _, err := unstructured.NestedString(shoot.Object, strings.Split(path, ".")...)
		if err != nil || got != want {
			t.Fatalf("%s: got %q, want %q (%v)", path, got, want, err)
		}
	}
	workers, _, _ := unstructured.NestedSlice(shoot.Object, "spec", "provider", "workers")
	worker := workers[0].(map[string]interface{})
	for _, key := range []string{"minimum", "maximum"} {
		got, _, err := unstructured.NestedInt64(worker, key)
		if err != nil || got != 1 {
			t.Fatalf("%s must be integer 1: %v (%v)", key, worker[key], err)
		}
	}
	machine, _, _ := unstructured.NestedString(worker, "machine", "type")
	if machine != "memory-intel1-1cpu" {
		t.Fatal(machine)
	}
	volume, _, _ := unstructured.NestedString(worker, "volume", "size")
	if volume != "50Gi" {
		t.Fatal(volume)
	}
	zones, _, _ := unstructured.NestedStringSlice(worker, "zones")
	if len(zones) != 3 {
		t.Fatal(zones)
	}
	scpc, err := renderSCPC(flavor{scpcName: ironcoreSCPCName, manifest: manifest}, "account")
	if err != nil {
		t.Fatal(err)
	}
	target, _, _ := unstructured.NestedString(scpc.Object, "spec", "target", "namespace")
	if target != "{{ .SC.namespace }}" {
		t.Fatal(target)
	}
	resources, _, _ := unstructured.NestedSlice(scpc.Object, "spec", "resources")
	if resources[0].(map[string]interface{})["template"] != defaults.template {
		t.Fatal("namespace and standalone callers use different templates")
	}
}

func TestGardenDevConfiguration(t *testing.T) { checkGardenDevManifest(t, gardenDevConfig(t)) }

func TestIronCoreChartRender(t *testing.T) {
	if _, err := exec.LookPath("helm"); err != nil {
		t.Skip("helm is required for chart render tests")
	}
	chart := "../../charts/sapprovider-operator"
	output, err := exec.Command("helm", "template", "test", chart).CombinedOutput()
	if err == nil || !strings.Contains(string(output), "scpc.ironcore.projectNamespace is required") {
		t.Fatalf("default values must fail: %s (%v)", output, err)
	}
	output, err = exec.Command("helm", "template", "test", chart, "-f", chart+"/examples/garden-dev.yaml").CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v", output, err)
	}
	var configData string
	for _, doc := range strings.Split(string(output), "\n---") {
		var obj map[string]interface{}
		if err := yaml.Unmarshal([]byte(doc), &obj); err != nil {
			t.Fatal(err)
		}
		if obj["kind"] == "ConfigMap" {
			configData, _, _ = unstructured.NestedString(obj, "data", "ironcore.yaml")
		}
		if obj["kind"] == "Deployment" {
			raw := string(output)
			for _, value := range []string{"--ironcore-config=/etc/sap-provider/ironcore.yaml", "checksum/ironcore-config:", "mountPath: /etc/sap-provider", "name: sap-provider-ironcore-config"} {
				if !strings.Contains(raw, value) {
					t.Fatalf("missing config connection: %s", value)
				}
			}
		}
	}
	if configData == "" {
		t.Fatal("chart has no IronCore config")
	}
	checkGardenDevManifest(t, []byte(configData))
}

func TestIronCoreExplicitMachineAndExtraSettings(t *testing.T) {
	var config IronCoreConfig
	if err := yaml.Unmarshal(gardenDevConfig(t), &config); err != nil {
		t.Fatal(err)
	}
	config.Profiles = nil
	autoscaling := false
	config.Autoscaling = &autoscaling
	workers, _, _ := unstructured.NestedSlice(config.ShootSpec, "provider", "workers")
	worker := workers[0].(map[string]interface{})
	if err := unstructured.SetNestedField(worker, "memory-intel1-8cpu", "machine", "type"); err != nil {
		t.Fatal(err)
	}
	if err := unstructured.SetNestedField(worker, "100Gi", "volume", "size"); err != nil {
		t.Fatal(err)
	}
	if err := unstructured.SetNestedStringSlice(worker, []string{"wdf4c"}, "zones"); err != nil {
		t.Fatal(err)
	}
	if err := unstructured.SetNestedSlice(config.ShootSpec, []interface{}{worker}, "provider", "workers"); err != nil {
		t.Fatal(err)
	}
	config.ShootSpec["purpose"] = "evaluation"
	config.ShootSpec["extensions"] = []interface{}{map[string]interface{}{"type": "shoot-cert-service"}}
	if err := unstructured.SetNestedField(config.ShootSpec, "10.0.0.0/16", "networking", "nodes"); err != nil {
		t.Fatal(err)
	}
	if err := unstructured.SetNestedField(config.ShootSpec, false, "networking", "providerConfig", "vxlan", "enabled"); err != nil {
		t.Fatal(err)
	}
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := ConfigureIronCoreManifest(data)
	if err != nil {
		t.Fatal(err)
	}
	defaults, err := parseFlavorDefaults(manifest)
	if err != nil {
		t.Fatal(err)
	}
	shoot, err := renderDefaultShootManifest(defaults, "account", "test")
	if err != nil {
		t.Fatal(err)
	}
	workers, _, _ = unstructured.NestedSlice(shoot.Object, "spec", "provider", "workers")
	worker = workers[0].(map[string]interface{})
	for path, want := range map[string]string{"machine.type": "memory-intel1-8cpu", "volume.size": "100Gi"} {
		got, _, err := unstructured.NestedString(worker, strings.Split(path, ".")...)
		if err != nil || got != want {
			t.Fatalf("%s = %q (%v)", path, got, err)
		}
	}
	scpc, err := renderSCPC(flavor{scpcName: ironcoreSCPCName, manifest: manifest}, "account")
	if err != nil {
		t.Fatal(err)
	}
	autoscale, found, err := unstructured.NestedBool(scpc.Object, "spec", "autoscaling")
	if err != nil || !found || autoscale {
		t.Fatalf("autoscaling override lost: %v %v %v", autoscale, found, err)
	}
	purpose, _, _ := unstructured.NestedString(shoot.Object, "spec", "purpose")
	if purpose != "evaluation" {
		t.Fatal(purpose)
	}
	nodes, _, _ := unstructured.NestedString(shoot.Object, "spec", "networking", "nodes")
	if nodes != "10.0.0.0/16" {
		t.Fatal(nodes)
	}
	vxlan, found, err := unstructured.NestedBool(shoot.Object, "spec", "networking", "providerConfig", "vxlan", "enabled")
	if err != nil || !found || vxlan {
		t.Fatalf("vxlan override lost: %v %v %v", vxlan, found, err)
	}
	zones, _, _ := unstructured.NestedStringSlice(worker, "zones")
	if len(zones) != 1 || zones[0] != "wdf4c" {
		t.Fatal(zones)
	}
	extensions, _, _ := unstructured.NestedSlice(shoot.Object, "spec", "extensions")
	if len(extensions) != 1 {
		t.Fatal(extensions)
	}
}
