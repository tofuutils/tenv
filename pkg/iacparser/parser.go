package iacparser

// ExtDescription represents a file extension and its description
type ExtDescription struct {
	Ext         string
	Description string
}

// DefaultExtensions returns the default list of infrastructure-as-code file extensions
func DefaultExtensions() []ExtDescription {
	return []ExtDescription{
		{Ext: ".tf", Description: "Terraform configuration file"},
		{Ext: ".tf.json", Description: "Terraform JSON configuration file"},
		{Ext: ".tfvars", Description: "Terraform variables file"},
		{Ext: ".tfvars.json", Description: "Terraform JSON variables file"},
		{Ext: ".hcl", Description: "HCL configuration file"},
		{Ext: ".json", Description: "JSON configuration file"},
		{Ext: ".yaml", Description: "YAML configuration file"},
		{Ext: ".yml", Description: "YAML configuration file"},
	}
}
