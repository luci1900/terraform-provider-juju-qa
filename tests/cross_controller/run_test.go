package qa

import (
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/juju/terraform-provider-juju-qa"
)

func TestQA_CrossController(t *testing.T) {
	// arrange
	main := utils.GetControllerInfo(t, utils.DefaultControllerName)
	offering := utils.GetControllerInfo(t, "tfqa-offering")

	offeringTfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./offering",
		EnvVars:      offering.Env(),
		Reconfigure:  true,
		NoColor:      true,
	})

	consumingTfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./consuming",
		EnvVars:      main.Env(),
		Vars:         offering.OfferingVars(),
		Reconfigure:  true,
		NoColor:      true,
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
		"application", "--timeout", "60m",
		"dummy-sink",
	)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed juju wait-for: %s", out)
	}
}
