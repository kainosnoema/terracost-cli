package mappings

import (
	"fmt"
)

var rdsEngineUsageOperationMap = map[string]string{
	"mysql:general-public-license":             "CreateDBInstance:0002",
	"oracle-se1:bring-your-own-license":        "CreateDBInstance:0003",
	"oracle-se:bring-your-own-license":         "CreateDBInstance:0004",
	"oracle-ee:bring-your-own-license":         "CreateDBInstance:0005",
	"oracle-se1:license-included":              "CreateDBInstance:0006",
	"sqlserver-se:bring-your-own-license":      "CreateDBInstance:0008",
	"sqlserver-ee:bring-your-own-license":      "CreateDBInstance:0009",
	"sqlserver-ex:license-included":            "CreateDBInstance:0010",
	"sqlserver-web:license-included":           "CreateDBInstance:0011",
	"sqlserver-se:license-included":            "CreateDBInstance:0012",
	"postgres:general-public-license":          "CreateDBInstance:0014",
	"sqlserver-ee:license-included":            "CreateDBInstance:0015",
	"aurora:general-public-license":            "CreateDBInstance:0016",
	"aurora-mysql:general-public-license":      "CreateDBInstance:0016",
	"mariadb:general-public-license":           "CreateDBInstance:0018",
	"oracle-se2:bring-your-own-license":        "CreateDBInstance:0019",
	"oracle-se2:license-included":              "CreateDBInstance:0020",
	"aurora-postgresql:general-public-license": "CreateDBInstance:0021",
}

var rdsInstanceClassMap = map[string]string{
	"db.m4.10xlarge":  "db.m4.10xl",
	"db.m4.16xlarge":  "db.m4.16xl",
	"db.m5.xlarge":    "db.m5.xl",
	"db.m5.2xlarge":   "db.m5.2xl",
	"db.m5.4xlarge":   "db.m5.4xl",
	"db.m5.8xlarge":   "db.m5.8xl",
	"db.m5.12xlarge":  "db.m5.12xl",
	"db.m5.16xlarge":  "db.m5.16xl",
	"db.m5.24xlarge":  "db.m5.24xl",
	"db.r4.16xlarge":  "db.r4.16xl",
	"db.r5.xlarge":    "db.r5.xl",
	"db.r5.2xlarge":   "db.r5.2xl",
	"db.r5.4xlarge":   "db.r5.4xl",
	"db.r5.8xlarge":   "db.r5.8xl",
	"db.r5.12xlarge":  "db.r5.12xl",
	"db.r5.16xlarge":  "db.r5.16xl",
	"db.r5.24xlarge":  "db.r5.24xl",
	"db.t3.xlarge":    "db.t3.xl",
	"db.t3.2xlarge":   "db.t3.2xl",
	"db.x1e.2xlarge":  "db.x1e.2xl",
	"db.x1e.4xlarge":  "db.x1e.4xl",
	"db.x1e.8xlarge":  "db.x1e.8xl",
	"db.x1e.16xlarge": "db.x1e.16xl",
	"db.x1e.32xlarge": "db.x1e.32xl",
	"db.x1.16xlarge":  "db.x1.16xl",
	"db.x1.32xlarge":  "db.x1.32xl",
	"db.z1d.xlarge":   "db.z1d.xl",
	"db.z1d.2xlarge":  "db.z1d.2xl",
	"db.z1d.3xlarge":  "db.z1d.3xl",
	"db.z1d.6xlarge":  "db.z1d.6xl",
	"db.z1d.12xlarge": "db.z1d.12xl",
}

func RDSInstance(region string, changeAttrs map[string]interface{}) string {
	instanceClass := changeAttrs["instance_class"].(string)
	if mappedClass, ok := rdsInstanceClassMap[instanceClass]; ok {
		instanceClass = mappedClass
	}

	licenseModel := ""
	if changeAttrs["license_model"] != nil {
		licenseModel = changeAttrs["license_model"].(string)
	}

	if licenseModel == "" {
		licenseModel = "general-public-license"
	}
	usageOperation := rdsEngineUsageOperationMap[changeAttrs["engine"].(string)+":"+licenseModel]

	return fmt.Sprintf(
		"%s-InstanceUsage:%s:%s",
		regionMap[region],
		instanceClass,
		usageOperation,
	)
}
