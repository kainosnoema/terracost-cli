package terraform

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

// ExecPlan runs `terraform plan` in the current directory and saves the JSON output
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

	showCmd := exec.Command("terraform", "show", "-json", planFile.Name())
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
