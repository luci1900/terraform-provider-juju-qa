package qa

import (
	"encoding/base64"
	"encoding/json"
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/juju/terraform-provider-juju-qa"
)

func TestQA_PrivateRegistry(t *testing.T) {
	info := utils.GetControllerInfo(t, utils.DefaultControllerName)
	if info.CloudType != "k8s" {
		t.Skip("Skipping test on non-k8s cloud")
	}

	// arrange
	tfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
		EnvVars:      info.Env(),
		Reconfigure:  true,
		NoColor:      true,
	})

	//act
	defer terraform.Destroy(t, tfOpts)
	terraform.InitAndApply(t, tfOpts)

	// assert
	ns := terraform.Output(t, tfOpts, "model_name")
	secretName := "test-app-coredns-secret"

	cmd := exec.Command(
		"microk8s", "kubectl",
		"-n", ns,
		"wait", "--for=create", "secret", secretName,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed call to kubectl: %s", out)
	}

	cmd = exec.Command(
		"microk8s", "kubectl",
		"get", "secret",
		"-n", ns,
		"-o", "json", secretName,
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed call to kubectl: %s", out)
	}
	var secret k8sSecret
	if err := json.Unmarshal(out, &secret); err != nil {
		t.Fatalf("failed to unmarshal secret")
	}
	encoded := secret.Data[".dockerconfigjson"]
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64")
	}
	var creds k8sCreds
	if err := json.Unmarshal(decoded, &creds); err != nil {
		t.Fatalf("failed to unmarshal credentials")
	}
	if creds.Auths["ghcr.io"].Username != "token" {
		t.Fatalf("invalid username")
	}
	if creds.Auths["ghcr.io"].Password != "token" {
		t.Fatalf("invalid password")
	}
}

type k8sSecret struct {
	Data map[string]string `json:"data"`
}

type k8sCreds struct {
	Auths map[string]ociCreds `json:"auths"`
}

type ociCreds struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}
