package qa

import (
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/juju/terraform-provider-juju-qa"
)

func TestQA_CrossController(t *testing.T) {
	// arrange
	offeringTfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./offering",
		EnvVars:      utils.GetControllerEnv(t, "tfqa-offering"),
		Reconfigure:  true,
	})

	consumingTfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./consuming",
		EnvVars:      utils.GetControllerEnv(t, utils.DefaultControllerName),
		Vars:         utils.GetOfferingControllerVars(t, "tfqa-offering"),
		Reconfigure:  true,
	})

	//act
	defer terraform.Destroy(t, offeringTfOpts)
	terraform.InitAndApply(t, offeringTfOpts)
	defer terraform.Destroy(t, consumingTfOpts)
	terraform.InitAndApply(t, consumingTfOpts)

	// assert
	modelName := terraform.Output(t, consumingTfOpts, "model_name")

	cmd := exec.Command(
		"juju", "switch",
		"tfqa:"+modelName,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed juju switch: %s", out)
	}

	cmd = exec.Command(
		"juju", "wait-for",
		"application", "dummy-sink",
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed juju wait-for: %s", out)
	}
}
