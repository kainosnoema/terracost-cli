package terraform

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

type PlanJSON struct {
	FormatVersion   string               `json:"format_version"`
	ResourceChanges []ResourceChangeJSON `json:"resource_changes"`
	Configuration   ConfigurationJSON    `json:"configuration"`
}

type ResourceChangeJSON struct {
	Address string
	Type    string
	Change  ChangeJSON
}

type ChangeJSON struct {
	Actions []string
	Before  map[string]interface{}
	After   map[string]interface{}
}

type ConfigurationJSON struct {
	ProviderConfig ProviderConfigJSON `json:"provider_config"`
}

type ProviderConfigJSON struct {
	AWS AWSJSON
}

type AWSJSON struct {
	Expressions RegionExpressionsJSON
}

type RegionExpressionsJSON struct {
	Region RegionConstantValueJSON
}

type RegionConstantValueJSON struct {
	ConstantValue string `json:"constant_value"`
}

func (p PlanJSON) Region() string {
	return p.Configuration.ProviderConfig.AWS.Expressions.Region.ConstantValue
}

// ExecPlan runs `terraform plan` in the current directory and returns the plan JSON
func ExecPlan() (*PlanJSON, error) {
	planFile, err := ioutil.TempFile("", "tc-plan")
	if err != nil {
		return nil, err
	}
	defer os.Remove(planFile.Name())

	planCmd := exec.Command("terraform", "plan", "-out", planFile.Name())
	// planCmd.Stdout = os.Stdout
	planCmd.Stderr = os.Stderr
	err = planCmd.Run()
	if err != nil {
		return nil, err
	}

	return ShowPlan(planFile.Name())
}

// ShowPlan runs `terraform show -json` against the provided plan file
func ShowPlan(planFileName string) (*PlanJSON, error) {
	showCmd := exec.Command("terraform", "show", "-json", planFileName)
	out, err := showCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var tfPlan PlanJSON
	err = json.Unmarshal(out, &tfPlan)
	if err != nil {
		return nil, err
	}

	return &tfPlan, nil
}
