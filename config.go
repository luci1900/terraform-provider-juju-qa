package utils

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/shell"
	"gopkg.in/yaml.v3"
)

type ControllerConfig struct {
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

func GetControllerEnv(t *testing.T, name string) map[string]string {
	cfg := getControllerConfig(t, name)
	return map[string]string{
		"JUJU_CONTROLLER_ADDRESSES": cfg.Addresses,
		"JUJU_USERNAME":             cfg.Username,
		"JUJU_PASSWORD":             cfg.Password,
		"JUJU_CA_CERT":              cfg.CACert,
	}
}

func GetOfferingControllerVars(t *testing.T, name string) map[string]any {
	cfg := getControllerConfig(t, name)
	return map[string]any{
		"offering_controller_name":      name,
		"offering_controller_addresses": cfg.Addresses,
		"offering_controller_username":  cfg.Username,
		"offering_controller_password":  cfg.Password,
		"offering_controller_ca_cert":   cfg.CACert,
	}
}

func getControllerConfig(t *testing.T, controllerName string) ControllerConfig {
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

	return ControllerConfig{
		Name:      controllerName,
		Addresses: addresses,
		Username:  accounts.Controllers[controllerName].User,
		Password:  accounts.Controllers[controllerName].Password,
		CACert:    showCtrl[controllerName].Details.CACert,
		CloudType: cloudType,
	}
}
