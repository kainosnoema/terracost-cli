package terraform

type PlanJSON struct {
	FormatVersion   string               `json:"format_version"`
	ResourceChanges []ResourceChangeJSON `json:"resource_changes"`
	Configuration   ConfigurationJSON    `json:"configuration"`
}

type ResourceChangeJSON struct {
	Address string
	Mode    string
	Type    string
	Name    string
	Index   int
	Change  ChangeJSON
}

type ChangeJSON struct {
	Before map[string]interface{}
	After  map[string]interface{}
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
