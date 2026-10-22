package charttests

import (
	"os/exec"
	"strings"
	"testing"
)

func TestAccountModeFlag(t *testing.T) {
	if _, err := exec.LookPath("helm"); err != nil {
		t.Skip("Helm is required; CI installs Helm")
	}
	const chart = "../../charts/apeiroprovider-operator"
	for _, mode := range []string{"false", "true"} {
		output, err := exec.Command("helm", "template", "account-test", chart, "-f", chart+"/examples/garden-dev.yaml", "--set", "scpc.ironCoreKCPAccount="+mode).CombinedOutput()
		if err != nil {
			t.Fatalf("helm: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "--ironcore-kcp-account="+mode) {
			t.Fatalf("chart did not pass account mode %s", mode)
		}
	}
}
