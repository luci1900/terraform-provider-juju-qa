package utils

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/shell"
	"gopkg.in/yaml.v3"
)

type ControllerInfo struct {
	Name      string
	Addresses string
	Username  string
	Password  string
	CACert    string
	CloudType string
}

type whoamiOutput struct {
	Controller string `json:"controller"`
}

type ctrlDetails struct {
	APIEndpoints []string `json:"api-endpoints"`
	CACert       string   `json:"ca-cert"`
	Cloud        string   `json:"cloud"`
}

type ctrlInfo struct {
	Details ctrlDetails `json:"details"`
}

type cloudInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ctrlCredentials struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type accountsFile struct {
	Controllers map[string]ctrlCredentials `yaml:"controllers"`
}

func GetCurrentControllerName(t *testing.T) string {
	out := shell.RunCommandAndGetOutput(t, shell.Command{
		Command: "juju",
		Args:    []string{"whoami", "--format", "json"},
	})

	var whoami whoamiOutput
	err := json.Unmarshal([]byte(out), &whoami)
	if err != nil {
		t.Fatalf("failed to unmarshal whoami output: %s", err)
	}
	return whoami.Controller
}

func GetControllerInfo(t *testing.T, controllerName string) ControllerInfo {
	out := shell.RunCommandAndGetOutput(t, shell.Command{
		Command: "juju",
		Args:    []string{"show-controller", controllerName, "--format", "json"},
	})

	var showCtrl map[string]ctrlInfo
	err := json.Unmarshal([]byte(out), &showCtrl)
	if err != nil {
		t.Fatalf("failed to unmarshal show-controller output: %s", err)
	}

	out = shell.RunCommandAndGetOutput(t, shell.Command{
		Command: "juju",
		Args:    []string{"show-cloud", showCtrl[controllerName].Details.Cloud, "--format", "json"},
	})

	var showCloud []cloudInfo
	err = json.Unmarshal([]byte(out), &showCloud)
	if err != nil {
		t.Fatalf("failed to unmarshal show-cloud output: %s", err)
	}

	file, err := os.Open(os.ExpandEnv("$HOME/.local/share/juju/accounts.yaml"))
	if err != nil {
		t.Fatalf("failed to open accounts file: %s", err)
	}
	defer file.Close()

	var accounts accountsFile
	err = yaml.NewDecoder(file).Decode(&accounts)
	if err != nil {
		t.Fatalf("failed to unmarshal accounts file: %s", err)
	}

	addresses := strings.Join(showCtrl[controllerName].Details.APIEndpoints, ",")
	var cloudType string
	for _, c := range showCloud {
		if c.Name == showCtrl[controllerName].Details.Cloud {
			cloudType = c.Type
			break
		}
	}

	return ControllerInfo{
		Name:      controllerName,
		Addresses: addresses,
		Username:  accounts.Controllers[controllerName].User,
		Password:  accounts.Controllers[controllerName].Password,
		CACert:    showCtrl[controllerName].Details.CACert,
		CloudType: cloudType,
	}
}

func (i ControllerInfo) Env() map[string]string {
	return map[string]string{
		"JUJU_CONTROLLER_ADDRESSES": i.Addresses,
		"JUJU_USERNAME":             i.Username,
		"JUJU_PASSWORD":             i.Password,
		"JUJU_CA_CERT":              i.CACert,
	}
}

func (i ControllerInfo) OfferingVars() map[string]any {
	return map[string]any{
		"offering_controller_name":      i.Name,
		"offering_controller_addresses": i.Addresses,
		"offering_controller_username":  i.Username,
		"offering_controller_password":  i.Password,
		"offering_controller_ca_cert":   i.CACert,
	}
}
