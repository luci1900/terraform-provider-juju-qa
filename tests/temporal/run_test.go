package qa

import (
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/juju/terraform-provider-juju-qa"
)

func TestQA_Temporal(t *testing.T) {
	info := utils.GetControllerInfo(t, utils.DefaultControllerName)
	if info.CloudType != "k8s" {
		t.Skip("Skipping private registry test on non-microk8s cloud")
	}

	// arrange
	tfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: ".",
		EnvVars:      info.Env(),
		Reconfigure:  true,
		NoColor:      true,
	})

	// act
	defer terraform.Destroy(t, tfOpts)
	terraform.InitAndApply(t, tfOpts)

	// assert
	modelName := terraform.Output(t, tfOpts, "model_name")
	cmd := exec.Command(
		"juju", "switch",
		modelName,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed juju switch: %s", out)
	}

	cmd = exec.Command(
		"juju", "wait-for",
		"application", "--timeout", "60m",
		"temporal",
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed juju wait-for: %s", out)
	}
}
