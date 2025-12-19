package qa

import (
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/juju/terraform-provider-juju-qa"
)

func TestQA_Minimal(t *testing.T) {
	// arrange
	info := utils.GetMainControllerInfo(t)

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
		info.Name+":"+modelName,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed juju switch: %s", out)
	}

	cmd = exec.Command(
		"juju", "wait-for",
		"application", "--timeout", "5m",
		"qa-test",
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed juju wait-for: %s", out)
	}
}
