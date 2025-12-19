package qa

import (
	"os/exec"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	utils "github.com/juju/terraform-provider-juju-qa"
)

func TestQA_CrossController(t *testing.T) {
	// arrange
	consumingInfo := utils.GetMainControllerInfo(t)
	offeringInfo := utils.GetOfferingControllerInfo(t)

	offeringTfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./offering",
		EnvVars:      offeringInfo.Env(),
		Reconfigure:  true,
		NoColor:      true,
	})

	consumingTfOpts := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "./consuming",
		EnvVars:      consumingInfo.Env(),
		Vars:         offeringInfo.OfferingVars(),
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
		consumingInfo.Name+":"+modelName,
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
